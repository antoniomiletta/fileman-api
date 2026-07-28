package auth

import "github.com/antoniomiletta/fileman/internal/domain"

var (
	ErrUnauthenticated = domain.ApplicationError{
		Msg: "not authenticated",
		Typ: domain.ErrorTypeAuthorization,
	}
	ErrForbidden = domain.ApplicationError{
		Msg: "not authorized",
		Typ: domain.ErrorTypeAuthorization,
	}
	ErrUserNotFound = domain.ApplicationError{
		Msg: "user not found",
		Typ: domain.ErrorTypeRetrieval,
	}
	ErrInvalidCredentials = domain.ApplicationError{
		Msg: "invalid credentials",
		Typ: domain.ErrorTypeAuthorization,
	}
	ErrCredentialsRequired = domain.ApplicationError{
		Msg: "email and password are required",
		Typ: domain.ErrorTypeValidation,
	}
	ErrEmailTaken = domain.ApplicationError{
		Msg: "email already taken",
		Typ: domain.ErrorTypeConflict,
	}
	ErrInvalidEmail = domain.ApplicationError{
		Msg: "invalid email format",
		Typ: domain.ErrorTypeValidation,
	}
	ErrPasswordTooWeak = domain.ApplicationError{
		Msg: "password is too weak",
		Typ: domain.ErrorTypeValidation,
	}
)
