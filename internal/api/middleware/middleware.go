package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/antoniomiletta/fileman/internal/api/transport"
	"github.com/antoniomiletta/fileman/internal/domain/auth"
	"github.com/antoniomiletta/fileman/internal/pkg/authenticator"
)

type ctxKey int

const userClaimsKey ctxKey = iota

type Middleware struct {
	authenticator *authenticator.Authenticator
}

func New(a *authenticator.Authenticator) *Middleware {
	return &Middleware{
		authenticator: a,
	}
}

func (m *Middleware) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := extractToken(r)
		if err != nil {
			transport.WriteError(w, auth.ErrUnauthenticated)
			return
		}

		claims, err := m.authenticator.VerifyToken(token)
		if err != nil {
			transport.WriteError(w, authenticator.ErrInvalidToken)
			return
		}

		ctx := context.WithValue(r.Context(), userClaimsKey, claims)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func extractToken(r *http.Request) (string, error) {
	h := r.Header.Get("Authorization")
	prefix := "Bearer"

	if h == "" {
		return "", auth.ErrUnauthenticated
	}

	if !strings.HasPrefix(h, prefix) {
		return "", fmt.Errorf("%w: expected Bearer {token}", transport.ErrMalformedToken)
	}

	return strings.TrimPrefix(h, prefix), nil

}
