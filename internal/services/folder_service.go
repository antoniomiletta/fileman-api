package services

import (
	"context"
	"time"

	"github.com/antoniomiletta/fileman/internal/domain"
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

// Store metadata to db with s.repo.Create()
func (s *FolderService) Create(ctx context.Context, folder *domain.Folder) error {
	folder.ID = uuid.New()
	folder.CreatedAt = time.Now()
	folder.UpdatedAt = time.Now()

	s.repo.Create(ctx, folder)

	return nil
}

func (s *FolderService) ListChildren(ctx context.Context, folderID uuid.UUID) ([]*domain.FolderContent, error) {
	s.repo.ListChildren(ctx, folderID)

	return []*domain.FolderContent{}, nil
}

func (s *FolderService) Move(ctx context.Context, id, newParentID uuid.UUID) error

func (s *FolderService) Delete(ctx context.Context, id uuid.UUID) error
