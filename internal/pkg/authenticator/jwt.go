package authenticator

import "github.com/google/uuid"

func Hash(password string) string
func CompareHash(a, b string) bool
func GetIDFromToken(token string) uuid.UUID
