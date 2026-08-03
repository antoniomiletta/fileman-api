package dto

import (
	"github.com/google/uuid"
)

type MoveFileRequest struct {
	NewParentID uuid.UUID `json:"new_parent_id"`
}

type RenameFileRequest struct {
	NewName string `json:"new_name"`
}

const (
	MultipartFieldFile     = "file"
	MultipartFieldParentID = "parent_id"
)
