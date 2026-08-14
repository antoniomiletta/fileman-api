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
	authn *authenticator.Authenticator
	resp  *transport.Responder
}

func New(authn *authenticator.Authenticator, resp *transport.Responder) *Middleware {
	return &Middleware{
		authn: authn,
		resp:  resp,
	}
}

func (m *Middleware) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := extractToken(r)
		if err != nil {
			m.resp.WriteError(w, err)
			return
		}

		claims, err := m.authn.VerifyToken(token)
		if err != nil {
			m.resp.WriteError(w, auth.ErrUnauthenticated)
			return
		}

		callerID, err := uuid.Parse(claims.Subject)
		if err != nil {
			m.resp.WriteError(w, auth.ErrUnauthenticated)
			return
		}

		ctx := reqctx.WithCallerID(r.Context(), callerID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func extractToken(r *http.Request) (string, error) {
	h := r.Header.Get("Authorization")
	prefix := "Bearer "

	if h == "" {
		return "", auth.ErrUnauthenticated
	}

	if !strings.HasPrefix(h, prefix) {
		return "", fmt.Errorf("%w: expected Bearer {token}", authenticator.ErrInvalidToken)
	}

	return strings.TrimPrefix(h, prefix), nil

}
