package services

import (
	"context"
	"io"
	"strings"

	"github.com/antoniomiletta/fileman/internal/domain/file"
	"github.com/antoniomiletta/fileman/internal/ports"
	"github.com/google/uuid"
)

type FileService struct {
	repo    ports.FileRepository
	storage ports.StorageBackend
}

func NewFileService(repo ports.FileRepository, storage ports.StorageBackend) *FileService {
	return &FileService{
		repo:    repo,
		storage: storage,
	}
}

// File content to storage with s.storage.Upload(), metadata to db with s.repo.Create()
func (s *FileService) Create(ctx context.Context, fil *file.File, content io.Reader) error {
	fil.Name = strings.TrimSpace(fil.Name)

	if err := file.ValidateFileName(fil.Name); err != nil {
		return err
	}

	if !file.IsAllowedMIME(fil.MIMEType) {
		return file.ErrFileTypeNotAllowed
	}

	s.repo.Create(ctx, fil)

	return nil
}

func (s *FileService) Move(ctx context.Context, id, newFolderID uuid.UUID) error
func (s *FileService) Delete(ctx context.Context, id uuid.UUID) error
