package transport

import "github.com/antoniomiletta/fileman/internal/domain"

var (
	ErrInvalidJSON = domain.ApplicationError{
		Msg: "invalid JSON payload",
		Typ: domain.ErrorTypeValidation,
	}
	ErrInvalidPathParam = domain.ApplicationError{
		Msg: "invalid path parameter",
		Typ: domain.ErrorTypeValidation,
	}
)
