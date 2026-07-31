package transport

import "github.com/antoniomiletta/fileman/internal/domain"

var (
	ErrMalformedJSON = domain.ApplicationError{
		Msg: "invalid JSON payload",
		Typ: domain.ErrorTypeValidation,
	}
	ErrInvalidPathParam = domain.ApplicationError{
		Msg: "invalid path parameter",
		Typ: domain.ErrorTypeValidation,
	}
	ErrMalformedToken = domain.ApplicationError{
		Msg: "invalid authentication token format",
		Typ: domain.ErrorTypeValidation,
	}
)
