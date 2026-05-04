package postgres

import (
	"context"

	"github.com/antoniomiletta/fileman/internal/domain"
)

type FolderRepository struct {
	db *DB
}

func NewFolderRepository(db *DB) *FolderRepository {
	return &FolderRepository{db: db}
}

// sql queries
func (r *FolderRepository) Create(ctx context.Context, folder *domain.Folder) error
func (r *FolderRepository) ListChildren(ctx context.Context, folderID string) ([]*domain.Folder, error)
func (r *FolderRepository) GetByID(ctx context.Context, id string) (*domain.Folder, error)
func (r *FolderRepository) Move(ctx context.Context, id, newParentID string) error
func (r *FolderRepository) Delete(ctx context.Context, id string) error
