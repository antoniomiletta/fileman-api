package auth

import "errors"

var (
	ErrForbidden           = errors.New("not authorized")
	ErrUserNotFound        = errors.New("user not found")
	ErrInvalidToken        = errors.New("invalid or expired authentication token")
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrCredentialsRequired = errors.New("email and password are required")
	ErrEmailTaken          = errors.New("email already taken")
	ErrInvalidEmail        = errors.New("invalid email format")
	ErrPasswordTooWeak     = errors.New("password is too weak")
)
