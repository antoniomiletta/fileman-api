package dto

import (
	"github.com/google/uuid"
)

type MoveFileRequest struct {
	NewParentID uuid.UUID `json:"new_parent_id"`
}

type DeleteFileRequest struct {
	ID uuid.UUID `json:"id"`
}

const (
	MultipartFieldFile     = "file"
	MultipartFieldParentID = "parent_id"
)
