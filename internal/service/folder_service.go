package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/antoniomiletta/fileman/internal/adapters/db/pg"
	"github.com/antoniomiletta/fileman/internal/domain/auth"
	"github.com/antoniomiletta/fileman/internal/domain/folder"
	"github.com/antoniomiletta/fileman/internal/pkg/reqctx"
	"github.com/antoniomiletta/fileman/internal/ports"
	"github.com/google/uuid"
)

type FolderService struct {
	folderRepo ports.FolderRepository
	txRunner   ports.TxRunner
}

func NewFolderService(repo ports.FolderRepository, txRunner ports.TxRunner) *FolderService {
	return &FolderService{
		folderRepo: repo,
		txRunner:   txRunner,
	}
}

func (s *FolderService) CreateFolder(ctx context.Context, parentID *uuid.UUID, name string) error {
	if err := folder.ValidateFolderName(name); err != nil {
		return err
	}

	if parentID == nil {
		return folder.ErrCannotCreateNewRootFolder
	}

	callerID, err := reqctx.CallerIDFrom(ctx)
	if err != nil {
		return err
	}

	parent, err := s.folderRepo.FindByID(ctx, *parentID)
	if err != nil {
		return err
	}

	if parent.OwnerID != callerID {
		return fmt.Errorf("%w: cannot create on specified parent", auth.ErrForbidden)
	}

	clash, err := s.folderRepo.FindByNameInParent(ctx, name, *parentID)
	if err != nil && !errors.Is(err, folder.ErrFolderNotFound) {
		return err
	}
	if clash != nil {
		return folder.ErrFolderNameConflict
	}

	f := folder.Folder{
		ID:       uuid.New(),
		OwnerID:  callerID,
		ParentID: parentID,
		Name:     strings.TrimSpace(name),
	}

	return s.folderRepo.Create(ctx, &f)
}

func (s *FolderService) ListChildren(ctx context.Context, folderID uuid.UUID) (*folder.FolderContent, error) {
	callerID, err := reqctx.CallerIDFrom(ctx)
	if err != nil {
		return nil, err
	}

	f, err := s.folderRepo.FindByID(ctx, folderID)
	if err != nil {
		return nil, err
	}

	if callerID != f.OwnerID {
		return nil, fmt.Errorf("%w: cannot access contents of specified folder", auth.ErrForbidden)
	}

	return s.folderRepo.ListChildren(ctx, folderID)
}

func (s *FolderService) MoveFolder(ctx context.Context, id, newParentID uuid.UUID) error {
	callerID, err := reqctx.CallerIDFrom(ctx)
	if err != nil {
		return err
	}

	f, err := s.folderRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if callerID != f.OwnerID {
		return fmt.Errorf("%w: cannot move this folder", auth.ErrForbidden)
	}

	if f.ParentID == nil {
		return folder.ErrCannotMoveRootFolder
	}

	if *f.ParentID == newParentID {
		return folder.ErrAlreadyInDestination
	}

	if id == newParentID {
		return folder.ErrCannotMoveIntoItself
	}

	dest, err := s.folderRepo.FindByID(ctx, newParentID)
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

		ancestor, err = s.folderRepo.FindByID(ctx, *ancestor.ParentID)
		if err != nil {
			return err
		}
	}

	clash, err := s.folderRepo.FindByNameInParent(ctx, f.Name, newParentID)
	if err != nil && !errors.Is(err, folder.ErrFolderNotFound) {
		return err
	}
	if clash != nil {
		return folder.ErrFolderNameConflict
	}

	return s.folderRepo.Move(ctx, id, newParentID)
}

func (s *FolderService) RenameFolder(ctx context.Context, id uuid.UUID, newName string) error {
	callerID, err := reqctx.CallerIDFrom(ctx)
	if err != nil {
		return err
	}

	f, err := s.folderRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if callerID != f.OwnerID {
		return fmt.Errorf("%w: cannot rename specified folder", auth.ErrForbidden)
	}

	if err := folder.ValidateFolderName(newName); err != nil {
		return err
	}

	return s.folderRepo.Rename(ctx, id, newName)
}

func (s *FolderService) DeleteFolder(ctx context.Context, id uuid.UUID) error {
	callerID, err := reqctx.CallerIDFrom(ctx)
	if err != nil {
		return err
	}

	f, err := s.folderRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if callerID != f.OwnerID {
		return fmt.Errorf("%w: cannot delete this folder", auth.ErrForbidden)
	}

	if f.ParentID == nil {
		return folder.ErrCannotDeleteRootFolder
	}

	return s.txRunner.RunTx(ctx, func(q ports.Querier) error {
		folderRepo := pg.NewFolderRepository(q)
		cleanupRepo := pg.NewCleanupJobRepository(q)

		keys, err := folderRepo.ListDescendantFileKeys(ctx, id)
		if err != nil {
			return err
		}

		for _, key := range keys {
			if err := cleanupRepo.Enqueue(ctx, key); err != nil {
				return err
			}
		}

		return folderRepo.Delete(ctx, id)
	})
}
