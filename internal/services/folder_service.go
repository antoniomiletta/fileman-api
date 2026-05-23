package services

import (
	"context"
	"strings"
	"time"

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

type CreateFolderInput struct {
	OwnerID  uuid.UUID
	ParentID *uuid.UUID
	Name     string
}

func (s *FolderService) CreateFolder(ctx context.Context, input CreateFolderInput) error {

	if err := folder.ValidateFolderName(input.Name); err != nil {
		return err
	}

	fol := folder.Folder{
		ID:        uuid.New(),
		OwnerID:   input.OwnerID,
		ParentID:  input.ParentID,
		Name:      strings.TrimSpace(input.Name),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.repo.Create(ctx, &fol); err != nil {
		return err
	}

	return nil
}

func (s *FolderService) ListChildren(ctx context.Context, folderID uuid.UUID) (*folder.FolderContent, error) {
	callerID, err := reqctx.CallerIDFrom(ctx)
	if err != nil {
		return nil, err
	}

	fol, err := s.repo.FindByID(ctx, folderID)
	if err != nil || fol == nil {
		return nil, folder.ErrFolderNotFound
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
	callerID, err := reqctx.CallerIDFrom(ctx)
	if err != nil {
		return err
	}

	fol, err := s.repo.FindByID(ctx, id)
	if err != nil || fol == nil {
		return folder.ErrFolderNotFound
	}

	if callerID != fol.OwnerID {
		return auth.ErrForbidden
	}

	if *fol.ParentID == newParentID {
		return folder.ErrAlreadyInDestination
	}

	if fol.ParentID == nil {
		return folder.ErrCannotMoveRootFolder
	}

	if id == newParentID {
		return folder.ErrCannotMoveIntoItself
	}

	dest, err := s.repo.FindByID(ctx, newParentID)
	if err != nil || dest == nil {
		return folder.ErrFolderNotFound
	}

	if callerID != dest.OwnerID {
		return auth.ErrForbidden
	}

	ancestor := dest
	for ancestor.ParentID != nil {
		if *ancestor.ParentID == id {
			return folder.ErrCannotMoveIntoDescendant
		}

		ancestor, err = s.repo.FindByID(ctx, *ancestor.ParentID)
		if err != nil {
			return folder.ErrFolderNotFound
		}
	}

	clash, _ := s.repo.FindByNameInParent(ctx, fol.Name, newParentID)
	if clash != nil {
		return folder.ErrFolderNameConflict
	}

	if err := s.repo.Move(ctx, id, newParentID); err != nil {
		return err
	}

	return nil
}

func (s *FolderService) Delete(ctx context.Context, id uuid.UUID) error {
	callerID, err := reqctx.CallerIDFrom(ctx)
	if err != nil {
		return err
	}

	fol, err := s.repo.FindByID(ctx, id)
	if err != nil || fol == nil {
		return folder.ErrFolderNotFound
	}

	if callerID != fol.OwnerID {
		return auth.ErrForbidden
	}

	if fol.ParentID == nil {
		return folder.ErrCannotDeleteRootFolder
	}

	// TODO: handle non empty folders (delete recursevely)

	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	return nil
}
