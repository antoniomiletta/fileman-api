package domain

import (
	"time"

	"github.com/google/uuid"
)

type File struct {
	ID          uuid.UUID  `json:"id"`
	ParentID    *uuid.UUID `json:"parentId"`
	Name        string     `json:"name"`
	Extension   string     `json:"extension"`
	Size        int64      `json:"size"`
	StoragePath string     `json:"storagePath"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
}
