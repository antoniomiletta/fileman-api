package services

import (
	"context"

	"github.com/antoniomiletta/fileman/internal/domain/folder"
	"github.com/antoniomiletta/fileman/internal/ports"
	"github.com/google/uuid"
)

type FolderService struct {
	repo ports.FolderRepository
}

func NewFolderService(repo ports.FolderRepository) *FolderService {
	return &FolderService{
		repo: repo,
	}
}

func (s *FolderService) Create(ctx context.Context, fol *folder.Folder) error {
	if fol.Name == "" {
		return folder.ErrFolderNameRequired
	}

	if len(fol.Name) > folder.MaxFolderNameLength {
		return folder.ErrFolderNameTooLong
	}

	return nil
}

func (s *FolderService) ListChildren(ctx context.Context, folderID uuid.UUID) ([]*folder.FolderContent, error)
func (s *FolderService) Move(ctx context.Context, id, newParentID uuid.UUID) error
func (s *FolderService) Delete(ctx context.Context, id uuid.UUID) error
