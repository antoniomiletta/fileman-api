package postgres

import (
	"context"

	"github.com/antoniomiletta/fileman/internal/domain"
	"github.com/google/uuid"
)

type FolderRepository struct {
	db *DB
}

func NewFolderRepository(db *DB) *FolderRepository {
	return &FolderRepository{db: db}
}

// sql queries
func (r *FolderRepository) Create(ctx context.Context, folder *domain.Folder) error
func (r *FolderRepository) ListChildren(ctx context.Context, folderID uuid.UUID) ([]*domain.Folder, error)
func (r *FolderRepository) Move(ctx context.Context, id, newParentID uuid.UUID) error
func (r *FolderRepository) Delete(ctx context.Context, id uuid.UUID) error
