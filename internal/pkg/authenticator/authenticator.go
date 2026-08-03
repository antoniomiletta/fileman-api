package authenticator

import (
	"fmt"
	"time"

	"github.com/antoniomiletta/fileman/config"
	"github.com/golang-jwt/jwt/v5"
)

type Authenticator struct {
	cfg config.AuthConfig
}

type Claims struct {
	jwt.RegisteredClaims
}

func NewAuthenticator(cfg config.AuthConfig) *Authenticator {
	return &Authenticator{cfg: cfg}
}

// GenerateToken returns a signed JWT token.
// The provided userID is registered as the 'subject' claim.
func (a *Authenticator) GenerateToken(userID string) (string, error) {
	now := time.Now()
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			ExpiresAt: jwt.NewNumericDate(now.Add(a.cfg.TokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed, err := token.SignedString([]byte(a.cfg.JWTSecret))
	if err != nil {
		return "", fmt.Errorf("authenticator: sign token: %w", err)
	}

	return signed, nil
}

// VerifyToken parses and validates the given token string, returning its claims.
// If the token is falsified, expired or uses an unexpected signing algorithm, it returns an error.
func (a *Authenticator) VerifyToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, ErrUnexpectedAlg
		}

		return []byte(a.cfg.JWTSecret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("authenticator: verify: parse token: %w", err)
	}
	if !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}
