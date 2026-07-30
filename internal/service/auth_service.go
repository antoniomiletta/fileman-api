package service

import (
	"context"
	"strings"

	"github.com/antoniomiletta/fileman/internal/adapter/db/postgres"
	"github.com/antoniomiletta/fileman/internal/domain/auth"
	"github.com/antoniomiletta/fileman/internal/pkg/authenticator"
	"github.com/antoniomiletta/fileman/internal/ports"
	"github.com/google/uuid"
)

type AuthService struct {
	authRepo      ports.AuthRepository
	txRunner      *postgres.TxRunner
	authenticator *authenticator.Authenticator
}

func NewAuthService(authRepo ports.AuthRepository, txRunner *postgres.TxRunner, authenticator *authenticator.Authenticator) *AuthService {
	return &AuthService{
		authRepo:      authRepo,
		txRunner:      txRunner,
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

	existing, err := s.authRepo.FindByEmail(ctx, input.Email)
	if err != nil {
		return err
	}
	if existing != nil {
		return auth.ErrEmailTaken
	}

	hashed, err := authenticator.Hash(input.Password)
	if err != nil {
		return err
	}

	user := auth.User{
		ID:       uuid.New(),
		Email:    strings.ToLower(strings.TrimSpace(input.Email)),
		Password: hashed,
	}

	return s.txRunner.Run(ctx, func(q postgres.Querier) error {
		authRepo := postgres.NewAuthRepository(q)
		if err := authRepo.SignUp(ctx, &user); err != nil {
			return err
		}

		folderRepo := postgres.NewFolderRepository(q)
		return folderRepo.CreateRoot(ctx, user.ID)
	})
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	email = strings.ToLower(strings.TrimSpace(email))

	if email == "" || password == "" {
		return "", auth.ErrCredentialsRequired
	}

	if !auth.IsValidEmail(email) {
		return "", auth.ErrInvalidCredentials
	}

	user, err := s.authRepo.FindByEmail(ctx, email)
	if err != nil || user == nil {
		return "", auth.ErrInvalidCredentials
	}

	if !authenticator.Compare(password, user.Password) {
		return "", auth.ErrInvalidCredentials
	}

	return s.authenticator.GenerateToken(user.ID.String())
}
