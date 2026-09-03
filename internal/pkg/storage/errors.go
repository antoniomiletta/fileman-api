package storage

import "github.com/antoniomiletta/fileman/internal/domain"

var (
	ErrFileSizeMismatch = domain.ApplicationError{
		Msg: "file size mismatch",
		Typ: domain.ErrorTypeValidation,
	}

	ErrObjectNotFound = domain.ApplicationError{
		Msg: "file object not found",
		Typ: domain.ErrorTypeRetrieval,
	}
)
