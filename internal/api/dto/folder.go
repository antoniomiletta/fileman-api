package dto

import (
	"github.com/antoniomiletta/fileman/internal/domain/folder"
	"github.com/google/uuid"
)

type CreateFolderRequest struct {
	ParentID *uuid.UUID `json:"parent_id"`
	Name     string     `json:"name"`
}

type MoveFolderRequest struct {
	NewParentID uuid.UUID `json:"new_parent_id"`
}

type RenameFolderRequest struct {
	NewName string `json:"new_name"`
}

type ListChildrenResponse struct {
	Content folder.FolderContent `json:"content"`
}
