package file

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
	ID         uuid.UUID
	OwnerID    uuid.UUID
	ParentID   *uuid.UUID
	Name       string
	MIMEType   string
	Size       int64
	StorageKey string
	Status     UploadStatus
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
