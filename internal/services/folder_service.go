package services

import (
	"context"

	"github.com/antoniomiletta/fileman/internal/domain"
	"github.com/antoniomiletta/fileman/internal/ports"
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
func (s *FolderService) CreateFolder(ctx context.Context, folder *domain.Folder) error
