package workers

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/antoniomiletta/fileman/config"
	"github.com/antoniomiletta/fileman/internal/jobs"
	"github.com/antoniomiletta/fileman/internal/ports"
)

type CleanupWorker struct {
	jobRepo ports.CleanupJobRepository
	store   ports.StorageBackend
	cfg     config.CleanupConfig
}

func NewCleanupWorker(jobRepo ports.CleanupJobRepository, store ports.StorageBackend, cfg config.CleanupConfig) *CleanupWorker {
	return &CleanupWorker{
		jobRepo: jobRepo,
		store:   store,
		cfg:     cfg,
	}
}

// Run polls for job batches and reclaims stale jobs until ctx is cancelled.
func (w *CleanupWorker) Run(ctx context.Context) {
	pollTicker := time.NewTicker(w.cfg.PollInterval)
	defer pollTicker.Stop()

	reclaimTicker := time.NewTicker(w.cfg.ReclaimInterval)
	defer reclaimTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("cleanup worker: shutting down")
			return
		case <-pollTicker.C:
			w.processBatch(ctx)
		case <-reclaimTicker.C:
			w.reclaimStaleJobs(ctx)
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

// reclaimStaleJobs is a safety net for crashes and other cases that can't be gracefully handled.
// For cases like context cancellation, dispatched jobs will finish before the worker returns.
func (w *CleanupWorker) reclaimStaleJobs(ctx context.Context) {
	reclaimed, err := w.jobRepo.Reclaim(ctx, w.cfg.StaleAfter)
	if err != nil {
		log.Printf("cleanup worker: %v", err)
	}

	log.Printf("cleanup worker: jobs reclaimed: %d", reclaimed)
}

func (w *CleanupWorker) recordSuccess(ctx context.Context, job jobs.CleanupJob) {
	if err := w.jobRepo.MarkDone(ctx, job.ID); err != nil {
		log.Printf("cleanup worker: mark job done: %s: %v", job.ID, err)
	}
}

func (w *CleanupWorker) recordFailure(ctx context.Context, job jobs.CleanupJob, cause error) {
	log.Printf("cleanup worker: job: %s failed (attempt %d/%d): %v", job.ID, job.Attempts+1, w.cfg.MaxAttempts, cause)

	if err := w.jobRepo.MarkFailed(ctx, job.ID, cause, w.cfg.MaxAttempts); err != nil {
		log.Printf("cleanup worker: mark job failed: %s: %v", job.ID, err)
	}
}
