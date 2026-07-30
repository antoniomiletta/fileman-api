package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/antoniomiletta/fileman/internal/api/transport"
	"github.com/antoniomiletta/fileman/internal/domain/auth"
	"github.com/antoniomiletta/fileman/internal/pkg/authenticator"
	"github.com/antoniomiletta/fileman/internal/pkg/reqctx"
	"github.com/google/uuid"
)

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
			transport.WriteError(w, err)
			return
		}

		claims, err := m.authenticator.VerifyToken(token)
		if err != nil {
			transport.WriteError(w, err)
			return
		}

		callerID, err := uuid.Parse(claims.Subject)
		if err != nil {
			transport.WriteError(w, auth.ErrUnauthenticated)
		}

		ctx := reqctx.WithCallerID(r.Context(), callerID)

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
