package postgres

import (
	"context"

	"github.com/antoniomiletta/fileman/internal/domain/file"
	"github.com/google/uuid"
)

type FileRepository struct {
	db Querier
}

func NewFileRepository(q Querier) *FileRepository {
	return &FileRepository{db: q}
}

// sql queries
func (r *FileRepository) Create(ctx context.Context, file *file.File) error
func (r *FileRepository) Move(ctx context.Context, id, newParentID uuid.UUID) error
func (r *FileRepository) Delete(ctx context.Context, id uuid.UUID) error
func (r *FileRepository) FindByNameInParent(ctx context.Context, name string, parentID uuid.UUID) (*file.File, error)
func (r *FileRepository) FindByID(ctx context.Context, id uuid.UUID) (*file.File, error)
