package jobs

import "github.com/google/uuid"

type CleanupJob struct {
	ID         uuid.UUID
	StorageKey string
	Attempts   int
}
