package services

import (
	"context"
	"io"

	"github.com/antoniomiletta/fileman/internal/domain"
	"github.com/antoniomiletta/fileman/internal/ports"
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
func (s *FileService) Upload(ctx context.Context, fileKey string, r io.Reader, size int64) (*domain.File, error)
