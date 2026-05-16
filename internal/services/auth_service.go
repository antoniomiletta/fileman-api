package services

import (
	"context"
	"strings"

	"github.com/antoniomiletta/fileman/internal/api/dto"
	"github.com/antoniomiletta/fileman/internal/domain"
	"github.com/antoniomiletta/fileman/internal/pkg/authenticator"
	"github.com/antoniomiletta/fileman/internal/ports"
)

type AuthService struct {
	repo ports.AuthRepository
}

func NewAuthService(repo ports.AuthRepository) *AuthService {
	return &AuthService{
		repo: repo,
	}
}

func (s *AuthService) Register(ctx context.Context, user *domain.User) error {
	user.Email = strings.ToLower(strings.TrimSpace(user.Email))

	if user.Email == "" || user.Password == "" {
		return domain.ErrCredentialsRequired
	}

	if !dto.IsValidEmail(user.Email) {
		return domain.ErrInvalidEmail
	}

	if !dto.IsPasswordStrong(user.Password) {
		return domain.ErrPasswordTooWeak
	}

	existing, _ := s.repo.FindByEmail(ctx, user.Email)
	if existing != nil {
		return domain.ErrEmailTaken
	}

	user.Password = authenticator.Hash(user.Password)

	if err := s.repo.Register(ctx, user); err != nil {
		return err
	}

	return nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	email = strings.ToLower(strings.TrimSpace(email))

	if email == "" || password == "" {
		return "", domain.ErrCredentialsRequired
	}

	if !dto.IsValidEmail(email) {
		return "", domain.ErrInvalidCredentials
	}

	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil || user == nil {
		return "", domain.ErrInvalidCredentials
	}

	if !authenticator.CompareHash(user.Password, password) {
		return "", domain.ErrInvalidCredentials
	}

	token, err := s.repo.Login(ctx, email, password)
	if err != nil {
		return "", err
	}

	return token, nil
}
