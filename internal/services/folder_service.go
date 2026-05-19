package services

import (
	"context"
	"strings"

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
	fol.Name = strings.TrimSpace(fol.Name)

	if err := folder.ValidateFolderName(fol.Name); err != nil {
		return err
	}

	if err := s.repo.Create(ctx, fol); err != nil {
		return err
	}

	return nil
}

func (s *FolderService) ListChildren(ctx context.Context, folderID uuid.UUID) (*folder.FolderContent, error) {
	fol, err := s.repo.FindByID(ctx, folderID)
	if err != nil || fol == nil {
		return nil, folder.ErrFolderNotFound
	}

	// TODO: ownderID

	return &folder.FolderContent{}, nil
}

func (s *FolderService) Move(ctx context.Context, id, newParentID uuid.UUID) error
func (s *FolderService) Delete(ctx context.Context, id uuid.UUID) error
