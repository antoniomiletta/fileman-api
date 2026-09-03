package pg

import (
	"context"
	"fmt"
	"time"

	"github.com/antoniomiletta/fileman/internal/jobs"
	"github.com/antoniomiletta/fileman/internal/ports"
	"github.com/google/uuid"
)

type CleanupJobRepository struct {
	db ports.Querier
}

func NewCleanupJobRepository(q ports.Querier) *CleanupJobRepository {
	return &CleanupJobRepository{
		db: q,
	}
}

// Enqueue records a cleanup job to be processed.
func (c *CleanupJobRepository) Enqueue(ctx context.Context, storageKey string) error {
	const query = `
		INSERT INTO storage_cleanup_jobs (id, storage_key, status, created_at, updated_at)
		VALUES ($1, $2, 'pending',NOW(), NOW())
		`

	_, err := c.db.Exec(ctx, query, uuid.New(), storageKey)
	if err != nil {
		return fmt.Errorf("postgres: enqueue cleanup job: %w", err)
	}

	return nil
}

// ClaimBatch selects up to limit pending jobs, locks and marks them as processing.
// Concurrent callers can never claim the same job because of SKIP LOCKED.
func (c *CleanupJobRepository) ClaimBatch(ctx context.Context, limit int) ([]jobs.CleanupJob, error) {
	const query = `
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
		`

	rows, err := c.db.Query(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("postgres: claim cleanup job batch: %w", err)
	}
	defer rows.Close()

	var batch []jobs.CleanupJob

	for rows.Next() {
		var j jobs.CleanupJob

		if err := rows.Scan(&j.ID, &j.StorageKey, &j.Attempts); err != nil {
			return nil, fmt.Errorf("postgres: scan cleanup job: %w", err)
		}

		batch = append(batch, j)

	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("postgres: iterate cleanup jobs: %w", err)
	}

	return batch, nil
}

// MarkDone marks a cleanup job as done.
func (c *CleanupJobRepository) MarkDone(ctx context.Context, id uuid.UUID) error {
	const query = `
		UPDATE storage_cleanup_jobs
		SET status = 'done', updated_at = NOW()
		WHERE id = $1
		`

	_, err := c.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("postgres: mark cleanup job done: %w", err)
	}

	return nil
}

// MarkFailed marks a cleanup job as failed (if max retries are exhausted) and records the error cause.
// If retries are still allowed, it marks the job back as pending.
func (c *CleanupJobRepository) MarkFailed(ctx context.Context, id uuid.UUID, cause error, maxRetries int) error {
	const query = `
		UPDATE storage_cleanup_jobs
		SET status = CASE WHEN attempts + 1 >= $1 THEN 'failed' ELSE 'pending' END, 
		attempts = attempts + 1,
		last_error = $2,
		updated_at = NOW()
		WHERE id = $3
		`

	_, err := c.db.Exec(ctx, query, maxRetries, cause.Error(), id)
	if err != nil {
		return fmt.Errorf("postgres: mark cleanup job failed: %w", err)
	}

	return nil
}

// Reclaim resets to pending every processing job that has not succeeded within staleTime.
func (c *CleanupJobRepository) Reclaim(ctx context.Context, staleTime time.Duration) (int64, error) {
	const query = `
		UPDATE storage_cleanup_jobs
		SET status = 'pending', updated_at = NOW()
		WHERE status = 'processing' 
		AND updated_at < NOW() - $1::interval
		`

	tag, err := c.db.Exec(ctx, query, staleTime.String())
	if err != nil {
		return 0, fmt.Errorf("postgres: reclaim stale jobs: %v", err)
	}

	return tag.RowsAffected(), nil
}
