package service

import (
	"context"
	"errors"
	"fmt"
	"io"

	"strings"

	"github.com/antoniomiletta/fileman/internal/adapters/db/pg"
	"github.com/antoniomiletta/fileman/internal/domain/auth"
	"github.com/antoniomiletta/fileman/internal/domain/file"
	"github.com/antoniomiletta/fileman/internal/domain/folder"
	"github.com/antoniomiletta/fileman/internal/pkg/reqctx"
	"github.com/antoniomiletta/fileman/internal/pkg/storage"
	"github.com/antoniomiletta/fileman/internal/ports"
	"github.com/google/uuid"
)

type FileService struct {
	fileRepo    ports.FileRepository
	folderRepo  ports.FolderRepository
	store       ports.StorageBackend
	txRunner    ports.TxRunner
	cleanupRepo ports.CleanupJobRepository
}

func NewFileService(
	fileRepo ports.FileRepository,
	folderRepo ports.FolderRepository,
	store ports.StorageBackend,
	txRunner ports.TxRunner,
	cleanupRepo ports.CleanupJobRepository,
) *FileService {
	return &FileService{
		fileRepo:    fileRepo,
		folderRepo:  folderRepo,
		store:       store,
		txRunner:    txRunner,
		cleanupRepo: cleanupRepo,
	}
}

type CreateFileInput struct {
	ParentID uuid.UUID
	Name     string
	MIMEType string
	Size     int64
}

func (s *FileService) CreateFile(ctx context.Context, input CreateFileInput, fileContent io.Reader) error {
	if err := file.ValidateFileName(input.Name); err != nil {
		return err
	}

	callerID, err := reqctx.CallerIDFrom(ctx)
	if err != nil {
		return err
	}

	parent, err := s.folderRepo.FindByID(ctx, input.ParentID)
	if err != nil {
		return err
	}

	if parent.OwnerID != callerID {
		return fmt.Errorf("%w: cannot create on specified parent", auth.ErrForbidden)
	}

	clash, err := s.fileRepo.FindByNameInParent(ctx, input.Name, input.ParentID)
	if err != nil && !errors.Is(err, file.ErrFileNotFound) {
		return err
	}
	if clash != nil {
		return file.ErrFileNameConflict
	}

	resourceType, ok := storage.Classify(input.MIMEType)
	if !ok {
		return fmt.Errorf("%w: %s", file.ErrFileTypeNotAllowed, input.MIMEType)
	}

	id := uuid.New()

	f := file.File{
		ID:       id,
		OwnerID:  callerID,
		ParentID: input.ParentID,
		Name:     strings.TrimSpace(input.Name),
		MIMEType: input.MIMEType,
		Size:     input.Size,
		StorageKey: storage.GenerateFileKey(storage.FileKeyParams{
			OwnerID:      callerID,
			ResourceType: resourceType,
			FileID:       id,
		}),
		UploadStatus: file.UploadStatusPending,
	}

	if err := s.fileRepo.Create(ctx, &f); err != nil {
		return err
	}

	if err := s.store.Upload(ctx, f.StorageKey, fileContent, f.Size); err != nil {
		return err
	}

	return s.fileRepo.MarkUploaded(ctx, f.ID)
}

func (s *FileService) MoveFile(ctx context.Context, id, newParentID uuid.UUID) error {
	callerID, err := reqctx.CallerIDFrom(ctx)
	if err != nil {
		return err
	}

	f, err := s.fileRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if callerID != f.OwnerID {
		return fmt.Errorf("%w: cannot move this file", auth.ErrForbidden)
	}

	if f.ParentID == newParentID {
		return folder.ErrAlreadyInDestination
	}

	dest, err := s.folderRepo.FindByID(ctx, newParentID)
	if err != nil {
		return err
	}

	if callerID != dest.OwnerID {
		return fmt.Errorf("%w: cannot move into specified folder", auth.ErrForbidden)
	}

	clash, err := s.fileRepo.FindByNameInParent(ctx, f.Name, newParentID)
	if err != nil && !errors.Is(err, file.ErrFileNotFound) {
		return err
	}
	if clash != nil {
		return file.ErrFileNameConflict
	}

	return s.fileRepo.Move(ctx, id, newParentID)
}

func (s *FileService) RenameFile(ctx context.Context, id uuid.UUID, newName string) error {
	callerID, err := reqctx.CallerIDFrom(ctx)
	if err != nil {
		return err
	}

	f, err := s.fileRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if callerID != f.OwnerID {
		return fmt.Errorf("%w: cannot rename specified file", auth.ErrForbidden)
	}

	if err := file.ValidateFileName(newName); err != nil {
		return err
	}

	return s.fileRepo.Rename(ctx, id, newName)
}

func (s *FileService) DeleteFile(ctx context.Context, id uuid.UUID) error {
	callerID, err := reqctx.CallerIDFrom(ctx)
	if err != nil {
		return err
	}

	f, err := s.fileRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if callerID != f.OwnerID {
		return fmt.Errorf("%w: cannot delete this file", auth.ErrForbidden)
	}

	return s.txRunner.RunTx(ctx, func(q ports.Querier) error {
		fileRepo := pg.NewFileRepository(q)
		if err := fileRepo.Delete(ctx, id); err != nil {
			return err
		}

		cleanupRepo := pg.NewCleanupJobRepository(q)
		if err := cleanupRepo.Enqueue(ctx, f.StorageKey); err != nil {
			return err
		}
		return nil
	})

}
