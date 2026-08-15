package service

import (
	"context"
	"errors"
	"strings"

	"github.com/antoniomiletta/fileman/internal/adapters/db/pg"
	"github.com/antoniomiletta/fileman/internal/domain/auth"
	"github.com/antoniomiletta/fileman/internal/pkg/authenticator"
	"github.com/antoniomiletta/fileman/internal/ports"
	"github.com/google/uuid"
)

type AuthService struct {
	authRepo ports.AuthRepository
	txRunner ports.TxRunner
	authn    *authenticator.Authenticator
}

func NewAuthService(
	authRepo ports.AuthRepository,
	txRunner ports.TxRunner,
	authn *authenticator.Authenticator,
) *AuthService {
	return &AuthService{
		authRepo: authRepo,
		txRunner: txRunner,
		authn:    authn,
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
	if err != nil && !errors.Is(err, auth.ErrUserNotFound) {
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

	// txRunners use temporary instances of repos so the operations can
	// run against the passed transaction instead of the default pool.
	return s.txRunner.RunTx(ctx, func(q ports.Querier) error {
		authRepo := pg.NewAuthRepository(q)
		if err := authRepo.SignUp(ctx, &user); err != nil {
			return err
		}

		folderRepo := pg.NewFolderRepository(q)
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
	if err != nil && errors.Is(err, auth.ErrUserNotFound) {
		return "", auth.ErrInvalidCredentials
	}

	if !authenticator.Compare(password, user.Password) {
		return "", auth.ErrInvalidCredentials
	}

	return s.authn.GenerateToken(user.ID.String())
}
