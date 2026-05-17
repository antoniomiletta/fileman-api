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

type NewFileParams struct {
	OwnerID   uuid.UUID
	ParentID  *uuid.UUID
	Name      string
	Extension string
	Size      int64
}

func NewFile(p NewFileParams) File {
	return File{
		ID:          uuid.New(),
		OwnerID:     p.OwnerID,
		ParentID:    p.ParentID,
		Name:        p.Name,
		Extension:   p.Extension,
		Size:        p.Size,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		StoragePath: "",
		Status:      UploadStatusPending,
	}
}
