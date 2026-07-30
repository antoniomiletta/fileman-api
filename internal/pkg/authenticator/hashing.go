package authenticator

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// Hash returns a bcrypt hash of the provided password.
func Hash(password string) (string, error) {
	cost := bcrypt.DefaultCost
	hash, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", fmt.Errorf("authenticator: hash generation: %w", err)
	}

	return string(hash), nil
}

// Compare reports whether a given password string
// is the plain text equivalent to a given hash string.
func Compare(password, hashed string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))

	return err == nil
}
