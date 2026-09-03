package workers

import (
	"context"
	"log"
	"time"

	"github.com/antoniomiletta/fileman/internal/ports"
)

type UploadReconciler struct {
	fileRepo ports.FileRepository
	store    ports.StorageBackend
}

func NewUploadReconciler(fileRepo ports.FileRepository, store ports.StorageBackend) *UploadReconciler {
	return &UploadReconciler{
		fileRepo: fileRepo,
		store:    store,
	}
}

// Sweep reconciles inconsistent upload records older than stale time by marking them as
// uploaded if they exist in storage, or deleting the record if they do not.
//
// Inconsistent uploads are either records pointing to non existent objects in storage or
// records incorrectly marked as pending. Since they only exist after unexpected crashes or panics,
// this is not intended to be called in a background loop
func (r *UploadReconciler) Sweep(ctx context.Context, staleTime time.Duration) error {
	stale, err := r.fileRepo.ListStale(ctx, staleTime)
	if err != nil {
		return err
	}

	for _, f := range stale {
		exists, err := r.store.Exists(ctx, f.StorageKey)
		if err != nil {
			log.Printf("upload reconciler: check existence for file: %s exists: %v", f.ID, err)
			continue
		}

		if exists {
			if err := r.fileRepo.MarkUploaded(ctx, f.ID); err != nil {
				log.Printf("upload reconciler: mark file: %s uploaded; %v", f.ID, err)
			}
			continue
		}

		if err := r.fileRepo.Delete(ctx, f.ID); err != nil {
			log.Printf("upload reconciler: delete file: %s:%v", f.ID, err)
		}
	}

	return nil
}
