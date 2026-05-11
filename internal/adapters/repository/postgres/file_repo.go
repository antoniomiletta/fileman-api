package postgres

import (
	"context"

	"github.com/antoniomiletta/fileman/internal/domain"
	"github.com/google/uuid"
)

type FileRepository struct {
	db *DB
}

func NewFileRepository(db *DB) *FileRepository {
	return &FileRepository{db: db}
}

// sql queries
func (r *FileRepository) Create(ctx context.Context, file *domain.File) error
func (r *FileRepository) Move(ctx context.Context, id, newFolderID uuid.UUID) error
func (r *FileRepository) Delete(ctx context.Context, id uuid.UUID) error
