package dto

import (
	"github.com/google/uuid"
)

type CreateFileRequest struct {
	ParentID *uuid.UUID `json:"parent_id"`
	Name     string     `json:"name"`
}

type MoveFileRequest struct {
	NewParentID uuid.UUID `json:"new_parent_id"`
}

type DeleteFileRequest struct {
	ID uuid.UUID `json:"id"`
}
