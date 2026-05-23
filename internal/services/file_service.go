package services

import (
	"context"
	"io"
	"strings"

	"github.com/antoniomiletta/fileman/internal/domain/auth"
	"github.com/antoniomiletta/fileman/internal/domain/file"
	"github.com/antoniomiletta/fileman/internal/domain/folder"
	"github.com/antoniomiletta/fileman/internal/pkg/storage"
	"github.com/antoniomiletta/fileman/internal/ports"
	"github.com/google/uuid"
)

type FileService struct {
	fileRepo   ports.FileRepository
	folderRepo ports.FolderRepository
	storage    ports.StorageBackend
}

func NewFileService(fileRepo ports.FileRepository, folderRepo ports.FolderRepository, storage ports.StorageBackend) *FileService {
	return &FileService{
		fileRepo:   fileRepo,
		folderRepo: folderRepo,
		storage:    storage,
	}
}

type CreateFileInput struct {
	OwnerID  uuid.UUID
	ParentID *uuid.UUID
	Name     string
	MIMEType string
	Size     int64
}

func (s *FileService) CreateFile(ctx context.Context, input CreateFileInput, content io.Reader) error {
	input.Name = strings.TrimSpace(input.Name)

	if err := file.ValidateFileName(input.Name); err != nil {
		return err
	}

	if !file.IsAllowedMIME(input.MIMEType) {
		return file.ErrFileTypeNotAllowed
	}

	parent, err := s.folderRepo.FindByID(ctx, *input.ParentID)
	if err != nil || parent == nil {
		return folder.ErrFolderNotFound
	}

	if parent.OwnerID != input.OwnerID {
		return auth.ErrForbidden
	}

	clash, _ := s.fileRepo.FindByNameInParent(ctx, input.Name, *input.ParentID)
	if clash != nil {
		return file.ErrFileNameConflict
	}

	fil := file.File{
		OwnerID:  input.OwnerID,
		ParentID: input.ParentID,
		Name:     input.Name,
		MIMEType: input.MIMEType,
		Size:     input.Size,
		StorageKey: storage.GenerateFileKey(storage.FileKeyParams{
			OwnerID:      input.OwnerID,
			ResourceType: storage.ResourceTypeFrom(input.MIMEType),
			Filename:     input.Name,
		}),
	}

	if err := s.storage.Upload(ctx, fil.StorageKey, content, fil.Size); err != nil {
		return err
	}

	if err := s.fileRepo.Create(ctx, &fil); err != nil {
		return err
	}

	return nil
}

func (s *FileService) MoveFile(ctx context.Context, id, newFolderID uuid.UUID) error
func (s *FileService) DeleteFile(ctx context.Context, id uuid.UUID) error
