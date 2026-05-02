package postgres

import (
	"context"

	"github.com/antoniomiletta/fileman/internal/domain"
)

type UserRepository struct {
	db *DB
}

func NewUserRepository(db *DB) *UserRepository {
	return &UserRepository{db: db}
}

// sql queries
func (r *UserRepository) Register(ctx context.Context, folder *domain.Folder) error
func (r *UserRepository) Login(ctx context.Context, folderID string) ([]*domain.Folder, error)
