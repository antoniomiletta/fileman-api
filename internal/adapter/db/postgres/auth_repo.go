package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/antoniomiletta/fileman/internal/domain/auth"
	"github.com/jackc/pgx/v5"
)

type AuthRepository struct {
	db Querier
}

func NewAuthRepository(q Querier) *AuthRepository {
	return &AuthRepository{db: q}
}

func (r *AuthRepository) SignUp(ctx context.Context, user *auth.User) error {
	const query = `
		INSERT INTO users (id, email, password_hash, created_at)
		VALUES ($1, $2, $3, NOW())
		`

	_, err := r.db.Exec(ctx, query, user.ID, user.Email, user.Password)
	if err != nil {
		if isUniqueViolation(err) {
			return auth.ErrEmailTaken
		}
	}

	return nil
}

func (r *AuthRepository) FindByEmail(ctx context.Context, email string) (*auth.User, error) {
	const query = `
		SELECT id, email, password_hash, created_at
		FROM users
		WHERE email = $1
		`

	u := &auth.User{}
	err := r.db.QueryRow(ctx, query, email).Scan(u.ID, u.Email, u.Password, u.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, auth.ErrUserNotFound
		}

		return nil, fmt.Errorf("postgres: find user by email: %w", err)
	}

	return u, nil
}
