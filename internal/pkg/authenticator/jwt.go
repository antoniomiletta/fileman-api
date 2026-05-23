package authenticator

import "github.com/google/uuid"

func Hash(password string) string
func CompareHash(a, b string) bool

// accept pure tokens and bearer strings
func GetIDFromToken(token string) uuid.UUID
