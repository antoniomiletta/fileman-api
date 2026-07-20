package authenticator

import (
	"time"

	"github.com/antoniomiletta/fileman/config"
	"github.com/antoniomiletta/fileman/internal/domain"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Authenticator struct {
	cfg config.AuthConfig
}

type Claims struct {
	jwt.RegisteredClaims
}

var (
	ErrUnexpectedAlg = domain.DomainError{
		Msg: "unexpected signing algorithm",
		Typ: domain.ErrorTypeValidation,
	}
	ErrInvalidToken = domain.DomainError{
		Msg: "invalid authentication token",
		Typ: domain.ErrorTypeValidation,
	}
)

func NewAuthenticator(cfg config.AuthConfig) *Authenticator {
	return &Authenticator{cfg: cfg}
}

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

	return token.SignedString(a.cfg.JWTSecret)
}

func (a *Authenticator) VerifyToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, ErrUnexpectedAlg
		}

		return a.cfg.JWTSecret, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

// accept pure tokens and bearer strings
func SubFromToken(token string) uuid.UUID
