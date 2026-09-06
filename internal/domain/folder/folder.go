package folder

import (
	"time"

	"github.com/antoniomiletta/fileman/internal/domain/file"
	"github.com/google/uuid"
)

type Folder struct {
	ID        uuid.UUID  `json:"id"`
	OwnerID   uuid.UUID  `json:"owner_id"`
	ParentID  *uuid.UUID `json:"parent_id"`
	Name      string     `json:"name"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type FolderContent struct {
	Subfolders []*Folder    `json:"subfolders"`
	Files      []*file.File `json:"files"`
}
