package auth

import "github.com/antoniomiletta/fileman/internal/domain"

var (
	ErrUnauthenticated = domain.DomainError{
		Msg: "not authenticated",
		Typ: domain.ErrorTypeAuthorization,
	}
	ErrForbidden = domain.DomainError{
		Msg: "not authorized",
		Typ: domain.ErrorTypeAuthorization,
	}
	ErrUserNotFound = domain.DomainError{
		Msg: "user not found",
		Typ: domain.ErrorTypeRetrieval,
	}
	ErrInvalidToken = domain.DomainError{
		Msg: "invalid or expired authentication token",
		Typ: domain.ErrorTypeAuthorization,
	}
	ErrInvalidCredentials = domain.DomainError{
		Msg: "invalid credentials",
		Typ: domain.ErrorTypeAuthorization,
	}
	ErrCredentialsRequired = domain.DomainError{
		Msg: "email and password are required",
		Typ: domain.ErrorTypeValidation,
	}
	ErrEmailTaken = domain.DomainError{
		Msg: "email already taken",
		Typ: domain.ErrorTypeConflict,
	}
	ErrInvalidEmail = domain.DomainError{
		Msg: "invalid email format",
		Typ: domain.ErrorTypeValidation,
	}
	ErrPasswordTooWeak = domain.DomainError{
		Msg: "password is too weak",
		Typ: domain.ErrorTypeValidation,
	}
)
