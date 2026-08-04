package workers

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/antoniomiletta/fileman/internal/jobs"
	"github.com/antoniomiletta/fileman/internal/ports"
)

type CleanupConfig struct {
	PollInterval time.Duration
	BatchSize    int
	MaxInFlight  int
	JobTimeout   time.Duration
}

type CleanupWorker struct {
	jobRepo ports.CleanupJobRepository
	store   ports.StorageBackend
	cfg     CleanupConfig
}

func NewCleanupWorker(jobRepo ports.CleanupJobRepository, store ports.StorageBackend, cfg CleanupConfig) *CleanupWorker {
	return &CleanupWorker{
		jobRepo: jobRepo,
		store:   store,
		cfg:     cfg,
	}
}

// Run polls for job batches until ctx is cancelled.
func (w *CleanupWorker) Run(ctx context.Context) {
	ticker := time.NewTicker(w.cfg.PollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("cleanup worker: shutting down")
			return
		case <-ticker.C:
			w.processBatch(ctx)
		}
	}
}

// processBatch claims one batch and runs its jobs into completion.
// Bounded to maxInFlight concurrent deletes.
func (w *CleanupWorker) processBatch(ctx context.Context) {
	batch, err := w.jobRepo.ClaimBatch(ctx, w.cfg.BatchSize)
	if err != nil {
		log.Printf("cleanup worker: claim batch: %v", err)
		return
	}

	inFlight := make(chan struct{}, w.cfg.MaxInFlight)

	var wg sync.WaitGroup
	defer wg.Wait()

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
}

// runJob runs a single cleanup job and records whether it failed or succeeded.
// Cancels on job timeout.
func (w *CleanupWorker) runJob(ctx context.Context, job jobs.CleanupJob) {
	jobCtx, jobCancel := context.WithTimeout(ctx, w.cfg.JobTimeout)
	defer jobCancel()

	// Recording the outcome always uses a fresh context, since an expired
	// job context would fail to recordFailure and possibly fail to
	// recordSuccess (if the job context expires between Delete succeeding and recordSuccess running).
	if err := w.store.Delete(jobCtx, job.StorageKey); err != nil {
		w.recordFailure(context.Background(), job, err)
		return
	}

	w.recordSuccess(context.Background(), job)
}

func (w *CleanupWorker) recordSuccess(ctx context.Context, job jobs.CleanupJob) {
	if err := w.jobRepo.MarkDone(ctx, job.ID); err != nil {
		log.Printf("cleanup worker: mark job: %s done: %v", job.ID, err)
	}
}

func (w *CleanupWorker) recordFailure(ctx context.Context, job jobs.CleanupJob, cause error) {
	if err := w.jobRepo.MarkFailed(ctx, job.ID, cause); err != nil {
		log.Printf("cleanup worker: mark job: %s failed: %v", job.ID, err)
	}
}
