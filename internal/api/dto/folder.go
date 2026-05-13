package dto

import (
	"github.com/antoniomiletta/fileman/internal/domain"
	"github.com/google/uuid"
)

type CreateFolderRequest struct {
	ParentID *uuid.UUID `json:"parent_id"`
	Name     string     `json:"name"`
}

type MoveFolderRequest struct {
	NewParentID uuid.UUID `json:"new_parent_id"`
}

type ListChildrenResponse struct {
	content domain.FolderContent
}

func (r *CreateFolderRequest) ToDomain() *domain.Folder {
	return &domain.Folder{
		ParentID: r.ParentID,
		Name:     r.Name,
	}
}
