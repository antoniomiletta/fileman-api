package service

import (
	"context"
	"strings"

	"github.com/antoniomiletta/fileman/internal/domain/auth"
	"github.com/antoniomiletta/fileman/internal/pkg/authenticator"
	"github.com/antoniomiletta/fileman/internal/ports"
	"github.com/google/uuid"
)

type AuthService struct {
	repo          ports.AuthRepository
	authenticator *authenticator.Authenticator
}

func NewAuthService(repo ports.AuthRepository, authenticator *authenticator.Authenticator) *AuthService {
	return &AuthService{
		repo:          repo,
		authenticator: authenticator,
	}
}

type CreateUserInput struct {
	Email    string
	Password string
}

func (s *AuthService) CreateUser(ctx context.Context, input CreateUserInput) error {
	if input.Email == "" || input.Password == "" {
		return auth.ErrCredentialsRequired
	}

	if !auth.IsValidEmail(input.Email) {
		return auth.ErrInvalidEmail
	}

	if !auth.IsPasswordStrong(input.Password) {
		return auth.ErrPasswordTooWeak
	}

	existing, _ := s.repo.FindByEmail(ctx, input.Email)
	if existing != nil {
		return auth.ErrEmailTaken
	}

	user := auth.User{
		ID:       uuid.New(),
		Email:    strings.ToLower(strings.TrimSpace(input.Email)),
		Password: authenticator.Hash(input.Password),
	}

	if err := s.repo.Register(ctx, &user); err != nil {
		return err
	}

	return nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	email = strings.ToLower(strings.TrimSpace(email))

	if email == "" || password == "" {
		return "", auth.ErrCredentialsRequired
	}

	if !auth.IsValidEmail(email) {
		return "", auth.ErrInvalidCredentials
	}

	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil || user == nil {
		return "", auth.ErrInvalidCredentials
	}

	if !authenticator.CompareHash(user.Password, password) {
		return "", auth.ErrInvalidCredentials
	}

	token, err := s.authenticator.GenerateToken(user.ID.String())
	if err != nil {
		return "", err
	}

	return token, nil
}
