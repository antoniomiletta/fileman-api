package postgres

import (
	"context"

	"github.com/antoniomiletta/fileman/internal/domain"
)

type FileRepository struct {
	db *DB
}

func NewFileRepository(db *DB) *FileRepository {
	return &FileRepository{db: db}
}

// sql queries
func (r *FileRepository) Create(ctx context.Context, file *domain.File) error
func (r *FileRepository) ListFromFolder(ctx context.Context, folderID string) ([]*domain.File, error)
