# Fileman
A file-manager backend in Go, built to learn and exercise backend engineering.

Go, PostgreSQL (`pgx`), a pluggable storage backend (`local disk`) and a background worker for async cleanup and reconciliation.


#### [Setup](#setup)

#### [Notes](#design-notes)
- *[Architecting](#architecting-ports-and-adapters)*
- *[Database](#database)*
- *[Async cleanup](#asynchronous-storage-cleanup)*
- *[Upload consistency](#upload-consistency-the-same-problem-in-reverse)*
---

# Setup

To run it locally, you can just copy the example env file:
```bash
cp .env.example .env
```

Then:

1\. Start Postgres:
   ```bash
   make docker-up
   ```

2\. Run migrations:
   ```bash
   make migrate-up
   ```

3\. Run the server:
   ```bash
   make dev
   ```


# Design Notes

A walkthrough of the project and some decisions that shaped it.

## Architecting: ports and adapters

```
handlers →  services →  ports (interfaces) ←  postgres/local/s3 (adapters)
                                    ↑
                        domain (pure types and rules)
```

Every dependency that touches external state is an
interface in `ports`, services never depend on concrete
implementations. The entire service layer is unit-tested
against in-memory fakes.

One exception: `Querier` (lets a repository run against either a
pool or an open transaction) references `pgx` types directly.
`pgx` is the chosen driver, you can't actually swap it out.

---

## Database

### Schema and data integrity

Each `user` has exactly one root folder (`parent_id IS NULL`), enforced by a
partial unique index. Every other folder or file always has a real parent.

Parent/child owner integrity is enforced with **composite foreign keys** — 
`(parent_id, owner_id) →  (id, owner_id)` which prevents
a file/folder from nesting under another user's tree, 
regardless of which code path tries to write to the table.

### Transactions

Every repository takes a `Querier`, satisfied by both a pool and a
transaction — repositories never own a connection, callers decide what
they run against.

`TxRunner` lets a service compose many repositories into one atomic operation:

```go
// User signup and root folder creation in one transaction.
	return s.txRunner.RunTx(ctx, func(q ports.Querier) error {
		authRepo := pg.NewAuthRepository(q)
		if err := authRepo.SignUp(ctx, &user); err != nil {
			return err
		}

		folderRepo := pg.NewFolderRepository(q)
		return folderRepo.CreateRoot(ctx, user.ID)
	})

```

I tried a generic `UnitOfWork[T]` first (a typed repo bundle)
but that meant one injected dependency per transaction shape.
`TxRunner` collapsed that to a single dependency where
any method inlines exactly the repos it needs inside `RunTx(fn())`.

---

## Asynchronous storage cleanup

Deleting a file means removing a DB row **and** the bytes behind it — two
systems that can't share a transaction. Doing both inline risks a partial
failure, leaking storage with no way of knowing something went wrong.

The fix: Postgres holds *intent*, the actual storage delete is an
async background job.

```sql
  BEGIN TX
    DELETE FROM files WHERE id = ...
    INSERT INTO storage_cleanup_jobs (...)
  COMMIT
```

How it actually looks with `TxRunner`:
```go
		return s.txRunner.RunTx(ctx, func(q ports.Querier) error {
		fileRepo := pg.NewFileRepository(q)
		if err := fileRepo.Delete(ctx, id); err != nil {
			return err
		}

		cleanupRepo := pg.NewCleanupJobRepository(q)
		return cleanupRepo.Enqueue(ctx, f.StorageKey)
	})
```

- *Folder deletion cascades through **FKs**, meaning deleting a folder would leak it's descendant files in storage.
To fix that, a recursive CTE captures every descendant file's
storage key **before** the cascade runs and enqueues their deletion inside a transaction:*

```go
	return s.txRunner.RunTx(ctx, func(q ports.Querier) error {
		folderRepo := pg.NewFolderRepository(q)
		cleanupRepo := pg.NewCleanupJobRepository(q)

		keys, err := folderRepo.ListDescendantFileKeys(ctx, id)
		if err != nil {
			return err
		}

		for _, key := range keys {
			if err := cleanupRepo.Enqueue(ctx, key); err != nil {
				return err
			}
		}

		return folderRepo.Delete(ctx, id)
	})
```

### Job batches and bounded concurrency

A worker polls for jobs, processing batches concurrently, capped at a
configurable max in-flight limit:

```go
		for _, job := range batch {
		select {
		case inFlight <- struct{}{}:
		case <-ctx.Done():
			return
		}
		wg.Go(func() {
			defer func() { <-inFlight }()
			w.runJob(ctx, job)
		})
	}
```

### Claiming work across concurrent workers

```sql
		UPDATE storage_cleanup_jobs
		SET status = 'processing', updated_at = NOW()
		WHERE id IN (
		SELECT id
		FROM storage_cleanup_jobs
		WHERE status = 'pending'
		ORDER BY created_at
		LIMIT $1
		FOR UPDATE SKIP LOCKED
		)
		RETURNING id, storage_key, attempts

```

`FOR UPDATE` locks the rows, `SKIP LOCKED` makes concurrent
callers skip them. Multiple workers can safely process
the queue with no extra coordination needed.

### Self-healing

A job claimed but never finished (crash mid-job) stays at `processing`
forever, invisible to future claims. A periodic sweep resets anything
stuck past a stale threshold back to `pending`:

```sql
		UPDATE storage_cleanup_jobs
		SET status = 'pending', updated_at = NOW()
		WHERE status = 'processing' 
		AND updated_at < NOW() - $1::interval

```

The same idea i kept returning to: don't try to prevent every
inconsistency, make it detectable and self-healing instead.

---

## Upload consistency, the same problem in reverse

File creation inserts metadata, then uploads bytes — a crash between the
two leaves one of two states:

| DB status | Storage     | Conclusion                |
| :-------- | :------- | :------------------------- |
| `pending` | `missing` | upload never finished, just delete the row |
| `pending` | `uploaded` | only the confirmation write was lost, mark as complete |


Currently, the reconciliation runs at startup, since its trigger (a
crash) and its remedy (a process restart) are the same event.

A single request panicking without crashing the whole process (which `net/http` does
by default) is a gap that allows inconsistencies until the next restart. The way i would like to solve this (not implemented yet) is a protected endpoint
that runs the same code path as the startup sweep, only behind a manual curl or cron job. no idle loop for a rare event.
