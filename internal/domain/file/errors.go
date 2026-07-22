package file

import "github.com/antoniomiletta/fileman/internal/domain"

var (
	ErrFileNotFound = domain.DomainError{
		Msg: "file not found",
		Typ: domain.ErrorTypeRetrieval,
	}

	ErrFileNameRequired = domain.DomainError{
		Msg: "file name is required",
		Typ: domain.ErrorTypeValidation,
	}

	ErrFileNameInvalid = domain.DomainError{
		Msg: "invalid file name",
		Typ: domain.ErrorTypeValidation,
	}

	ErrFileNameTooLong = domain.DomainError{
		Msg: "file name is too long",
		Typ: domain.ErrorTypeValidation,
	}

	ErrFileNameConflict = domain.DomainError{
		Msg: "file name already exists in target directory",
		Typ: domain.ErrorTypeConflict,
	}

	ErrFileTooLarge = domain.DomainError{
		Msg: "file is too large",
		Typ: domain.ErrorTypeValidation,
	}

	ErrFileTypeNotAllowed = domain.DomainError{
		Msg: "this file type is not allowed",
		Typ: domain.ErrorTypeValidation,
	}

	ErrInvalidStatus = domain.DomainError{
		Msg: "invalid status transition",
		Typ: domain.ErrorTypeLogical,
	}

	ErrInexistentParent = domain.DomainError{
		Msg: "target parent does not exist",
		Typ: domain.ErrorTypeLogical,
	}
)
