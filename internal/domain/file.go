package domain

import (
	"time"

	"github.com/google/uuid"
)

type UploadStatus string

const (
	UploadStatusPending  UploadStatus = "pending"
	UploadStatusComplete UploadStatus = "complete"
)

type File struct {
	ID          uuid.UUID
	OwnerID     uuid.UUID
	ParentID    *uuid.UUID
	Name        string
	Extension   string
	Size        int64
	CreatedAt   time.Time
	UpdatedAt   time.Time
	StoragePath string
	Status      UploadStatus
}
