package middlewares

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

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := extractToken(r)
		if err != nil {
			transport.WriteError(w, auth.ErrForbidden)
			return
		}

		claims, err := authenticator.ValidateToken(token)
		if err != nil {
			transport.WriteError(w, auth.ErrInvalidToken)
			return
		}

		ctx := context.WithValue(r.Context(), userClaimsKey, claims)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func extractToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")

	if authHeader == "" {
		return "", auth.ErrUnauthenticated
	}

	split := strings.Split(authHeader, " ")
	if len(split) != 2 || !strings.EqualFold(split[0], "Bearer") {
		return "", fmt.Errorf("%w: expected Bearer {token}", transport.ErrMalformedToken)
	}

	return split[1], nil
}
