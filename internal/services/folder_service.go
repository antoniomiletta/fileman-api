package services

import (
	"context"
	"strings"

	"github.com/antoniomiletta/fileman/internal/domain/auth"
	"github.com/antoniomiletta/fileman/internal/domain/folder"
	"github.com/antoniomiletta/fileman/internal/pkg/reqctx"
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

	callerID, err := reqctx.CallerIDFrom(ctx)
	if err != nil {
		return nil, err
	}

	if callerID != fol.OwnerID {
		return nil, auth.ErrForbidden
	}

	content, err := s.repo.ListChildren(ctx, folderID)
	if err != nil {
		return nil, err
	}

	return content, nil
}

func (s *FolderService) Move(ctx context.Context, id, newParentID uuid.UUID) error {
	fol, err := s.repo.FindByID(ctx, id)
	if err != nil || fol == nil {
		return folder.ErrFolderNotFound
	}

	if fol.ParentID == nil {
		return folder.ErrCannotMoveRootFolder
	}

	if *fol.ParentID == newParentID {
		return folder.ErrAlreadyInDestination
	}

	if id == newParentID {
	}

	return nil
}

func (s *FolderService) Delete(ctx context.Context, id uuid.UUID) error
