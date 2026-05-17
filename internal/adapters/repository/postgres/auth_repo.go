package postgres

import (
	"context"

	"github.com/antoniomiletta/fileman/internal/domain/auth"
)

type AuthRepository struct {
	db *DB
}

func NewAuthRepository(db *DB) *AuthRepository {
	return &AuthRepository{db: db}
}

// sql queries
func (r *AuthRepository) Register(ctx context.Context, user *auth.User) error
func (r *AuthRepository) Login(ctx context.Context, email, password string) (string, error)
func (r *AuthRepository) FindByEmail(ctx context.Context, email string) (*auth.User, error)
