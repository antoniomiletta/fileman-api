package postgres

import (
	"context"

	"github.com/antoniomiletta/fileman/internal/domain"
)

type AuthRepository struct {
	db *DB
}

func NewAuthRepository(db *DB) *AuthRepository {
	return &AuthRepository{db: db}
}

// sql queries
func (r *AuthRepository) Register(ctx context.Context, user *domain.User) error {
	r.db.conn.WithContext(ctx).Create(&user)
	return nil
}

func (r *AuthRepository) Login(ctx context.Context, email, password string) (string, error)
