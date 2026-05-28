package authenticator

import "github.com/google/uuid"

type JWTSignPayload struct {
	ID       string
	Email    string
	Password string
}

type UserClaims struct {
	ID string
}

func Hash(password string) string
func CompareHash(a, b string) bool

// accept pure tokens and bearer strings
func GetIDFromToken(token string) uuid.UUID

func Sign(payload JWTSignPayload) string
func ValidateToken(token string) (*UserClaims, error)
