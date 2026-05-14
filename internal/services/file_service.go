package services

import (
	"context"

	"github.com/antoniomiletta/fileman/internal/domain"
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
func (s *FileService) Create(ctx context.Context, file *domain.File) error
func (s *FileService) Move(ctx context.Context, id, newFolderID uuid.UUID) error
func (s *FileService) Delete(ctx context.Context, id uuid.UUID) error
