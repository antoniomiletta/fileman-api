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

type NewFolderParams struct {
	OwnerID  uuid.UUID
	ParentID *uuid.UUID
	Name     string
}

func NewFolder(p NewFolderParams) Folder {
	return Folder{
		ID:        uuid.New(),
		OwnerID:   p.OwnerID,
		ParentID:  p.ParentID,
		Name:      p.Name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
