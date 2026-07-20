package service

import (
	"context"
	"fmt"
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
	input.Name = strings.ToLower(strings.TrimSpace(input.Name))

	if err := folder.ValidateFolderName(input.Name); err != nil {
		return err
	}

	if input.ParentID != nil {
		parent, err := s.repo.FindByID(ctx, *input.ParentID)
		if err != nil {
			return fmt.Errorf("%w: specified parent does not exist", err)
		}

		if parent.OwnerID != input.OwnerID {
			return fmt.Errorf("%w: cannot create on specified parent", auth.ErrForbidden)
		}
	}

	clash, _ := s.repo.FindByNameInParent(ctx, input.Name, *input.ParentID)
	if clash != nil {
		return folder.ErrFolderNameConflict
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
	if err != nil {
		return nil, err
	}

	if callerID != fol.OwnerID {
		return nil, fmt.Errorf("%w: cannot access contents of specified folder", auth.ErrForbidden)
	}

	content, err := s.repo.ListChildren(ctx, folderID)
	if err != nil {
		return nil, err
	}

	return content, nil
}

func (s *FolderService) MoveFolder(ctx context.Context, id, newParentID uuid.UUID) error {
	callerID, err := reqctx.CallerIDFrom(ctx)
	if err != nil {
		return err
	}

	fol, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if callerID != fol.OwnerID {
		return fmt.Errorf("%w: cannot move this folder", auth.ErrForbidden)
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
	if err != nil {
		return err
	}

	if callerID != dest.OwnerID {
		return fmt.Errorf("%w: cannot move into specified folder", auth.ErrForbidden)
	}

	ancestor := dest
	for ancestor.ParentID != nil {
		if *ancestor.ParentID == id {
			return folder.ErrCannotMoveIntoDescendant
		}

		ancestor, err = s.repo.FindByID(ctx, *ancestor.ParentID)
		if err != nil {
			return err
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

func (s *FolderService) DeleteFolder(ctx context.Context, id uuid.UUID) error {
	callerID, err := reqctx.CallerIDFrom(ctx)
	if err != nil {
		return err
	}

	fol, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if callerID != fol.OwnerID {
		return fmt.Errorf("%w: cannot delete this folder", auth.ErrForbidden)
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
