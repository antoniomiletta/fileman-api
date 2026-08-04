package ports

import (
	"context"

	"github.com/antoniomiletta/fileman/internal/jobs"
	"github.com/google/uuid"
)

type CleanupJobRepository interface {
	Enqueue(ctx context.Context, storageKey string) error
	ClaimBatch(ctx context.Context, limit int) ([]jobs.CleanupJob, error)
	MarkDone(ctx context.Context, id uuid.UUID) error
	MarkFailed(ctx context.Context, id uuid.UUID, cause error) error
}
