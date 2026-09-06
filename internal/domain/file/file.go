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
	ID           uuid.UUID    `json:"id"`
	OwnerID      uuid.UUID    `json:"owner_id"`
	ParentID     uuid.UUID    `json:"parent_id"`
	Name         string       `json:"name"`
	MIMEType     string       `json:"mime_type"`
	Size         int64        `json:"size"`
	StorageKey   string       `json:"storage_key"`
	UploadStatus UploadStatus `json:"upload_status"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
}
