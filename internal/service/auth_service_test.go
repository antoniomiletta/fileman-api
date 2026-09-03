package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/antoniomiletta/fileman/config"
	"github.com/antoniomiletta/fileman/internal/domain/auth"
	"github.com/antoniomiletta/fileman/internal/pkg/authenticator"
	"github.com/antoniomiletta/fileman/internal/service"
	"github.com/google/uuid"
)

type fakeAuthRepo struct {
	usersByEmail   map[string]*auth.User
	findByEmailErr error
}

func newFakeAuthRepo() *fakeAuthRepo {
	return &fakeAuthRepo{usersByEmail: make(map[string]*auth.User)}
}

func (f *fakeAuthRepo) seed(u *auth.User) {
	f.usersByEmail[u.Email] = u
}

func (f *fakeAuthRepo) SignUp(ctx context.Context, u *auth.User) error {
	f.usersByEmail[u.Email] = u
	return nil
}

func (f *fakeAuthRepo) FindByEmail(ctx context.Context, email string) (*auth.User, error) {
	if f.findByEmailErr != nil {
		return nil, f.findByEmailErr
	}
	u, ok := f.usersByEmail[email]
	if !ok {
		return nil, auth.ErrUserNotFound
	}
	return u, nil
}

func newTestAuthenticator(t *testing.T) *authenticator.Authenticator {
	return authenticator.NewAuthenticator(config.AuthConfig{
		JWTSecret:   "test-secret",
		TokenExpiry: time.Hour,
	})
}

func TestAuthService_CreateUser(t *testing.T) {
	setup := func() (*fakeAuthRepo, *authenticator.Authenticator) {
		return newFakeAuthRepo(), newTestAuthenticator(t)
	}

	t.Run("rejects empty credentials", func(t *testing.T) {
		repo, authn := setup()
		svc := service.NewAuthService(repo, newFakeTxRunner(t), authn)

		err := svc.CreateUser(context.Background(), service.CreateUserInput{Email: "", Password: ""})
		if !errors.Is(err, auth.ErrCredentialsRequired) {
			t.Fatalf("expected ErrCredentialsRequired, got: %v", err)
		}
	})

	t.Run("rejects invalid email", func(t *testing.T) {
		repo, authn := setup()
		svc := service.NewAuthService(repo, newFakeTxRunner(t), authn)

		err := svc.CreateUser(context.Background(), service.CreateUserInput{
			Email: "not-an-email", Password: "SuperStrongP@ss1",
		})
		if !errors.Is(err, auth.ErrInvalidEmail) {
			t.Fatalf("expected ErrInvalidEmail, got: %v", err)
		}
	})

	t.Run("rejects weak password", func(t *testing.T) {
		repo, authn := setup()
		svc := service.NewAuthService(repo, newFakeTxRunner(t), authn)

		err := svc.CreateUser(context.Background(), service.CreateUserInput{
			Email: "user@example.com", Password: "123",
		})
		if !errors.Is(err, auth.ErrPasswordTooWeak) {
			t.Fatalf("expected ErrPasswordTooWeak, got: %v", err)
		}
	})

	t.Run("rejects an email already taken", func(t *testing.T) {
		repo, authn := setup()
		repo.seed(&auth.User{ID: uuid.New(), Email: "user@example.com", Password: "irrelevant-hash"})
		svc := service.NewAuthService(repo, newFakeTxRunner(t), authn)

		err := svc.CreateUser(context.Background(), service.CreateUserInput{
			Email: "user@example.com", Password: "SuperStrongP@ss1",
		})
		if !errors.Is(err, auth.ErrEmailTaken) {
			t.Fatalf("expected ErrEmailTaken, got: %v", err)
		}
	})
}

func TestAuthService_Login(t *testing.T) {
	setup := func(t *testing.T) (*fakeAuthRepo, *authenticator.Authenticator) {
		return newFakeAuthRepo(), newTestAuthenticator(t)
	}

	t.Run("rejects empty credentials", func(t *testing.T) {
		repo, authn := setup(t)
		svc := service.NewAuthService(repo, newFakeTxRunner(t), authn)

		_, err := svc.Login(context.Background(), "", "")
		if !errors.Is(err, auth.ErrCredentialsRequired) {
			t.Fatalf("expected ErrCredentialsRequired, got: %v", err)
		}
	})

	t.Run("rejects invalid email format", func(t *testing.T) {
		repo, authn := setup(t)
		svc := service.NewAuthService(repo, newFakeTxRunner(t), authn)

		_, err := svc.Login(context.Background(), "not-an-email", "SuperStrongP@ss1")
		if !errors.Is(err, auth.ErrInvalidCredentials) {
			t.Fatalf("expected ErrInvalidCredentials, got: %v", err)
		}
	})

	t.Run("rejects a nonexistent user without leaking existence", func(t *testing.T) {
		repo, authn := setup(t)
		svc := service.NewAuthService(repo, newFakeTxRunner(t), authn)

		_, err := svc.Login(context.Background(), "ghost@example.com", "SuperStrongP@ss1")
		if !errors.Is(err, auth.ErrInvalidCredentials) {
			t.Fatalf("expected ErrInvalidCredentials, got: %v", err)
		}
	})

	t.Run("rejects an incorrect password", func(t *testing.T) {
		repo, authn := setup(t)
		hashed, err := authenticator.Hash("CorrectP@ss1")
		if err != nil {
			t.Fatalf("failed to hash test password: %v", err)
		}
		repo.seed(&auth.User{ID: uuid.New(), Email: "user@example.com", Password: hashed})
		svc := service.NewAuthService(repo, newFakeTxRunner(t), authn)

		_, err = svc.Login(context.Background(), "user@example.com", "WrongPassword")
		if !errors.Is(err, auth.ErrInvalidCredentials) {
			t.Fatalf("expected ErrInvalidCredentials, got: %v", err)
		}
	})

	t.Run("propagates a repository error instead of panicking", func(t *testing.T) {
		repo, authn := setup(t)
		repo.findByEmailErr = errors.New("simulated db failure")
		svc := service.NewAuthService(repo, newFakeTxRunner(t), authn)

		_, err := svc.Login(context.Background(), "user@example.com", "whatever")
		if err == nil || errors.Is(err, auth.ErrInvalidCredentials) {
			t.Fatalf("expected the underlying repository error to propagate, got: %v", err)
		}
	})

	t.Run("successful login returns a token", func(t *testing.T) {
		repo, authn := setup(t)
		hashed, err := authenticator.Hash("CorrectP@ss1")
		if err != nil {
			t.Fatalf("failed to hash test password: %v", err)
		}
		user := &auth.User{ID: uuid.New(), Email: "user@example.com", Password: hashed}
		repo.seed(user)
		svc := service.NewAuthService(repo, newFakeTxRunner(t), authn)

		token, err := svc.Login(context.Background(), "user@example.com", "CorrectP@ss1")
		if err != nil {
			t.Fatalf("expected success, got: %v", err)
		}
		if token == "" {
			t.Fatal("expected a non-empty token")
		}
	})
}
