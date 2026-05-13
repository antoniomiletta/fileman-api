package domain

import (
	"time"

	"github.com/google/uuid"
)

type Folder struct {
	ID        uuid.UUID
	OwnerID   uuid.UUID
	ParentID  *uuid.UUID
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type FolderContent struct {
	Subfolders []Folder
	Files      []File
}
