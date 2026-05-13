package dto

import (
	"path/filepath"

	"github.com/antoniomiletta/fileman/internal/domain"
	"github.com/google/uuid"
)

type CreateFileRequest struct {
	ParentID *uuid.UUID `json:"parent_id"`
	Name     string     `json:"name"`
}

type MoveFileRequest struct {
	NewParentID uuid.UUID `json:"new_parent_id"`
}

func (r *CreateFileRequest) ToDomain() *domain.File {
	ext := filepath.Ext(r.Name)
	name := r.Name[:len(r.Name)-len(ext)]

	return &domain.File{
		ParentID:  r.ParentID,
		Name:      name,
		Extension: ext,
	}
}
