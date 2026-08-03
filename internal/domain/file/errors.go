package file

import "github.com/antoniomiletta/fileman/internal/domain"

var (
	ErrFileNotFound = domain.ApplicationError{
		Msg: "file not found",
		Typ: domain.ErrorTypeRetrieval,
	}

	ErrFileNameRequired = domain.ApplicationError{
		Msg: "file name is required",
		Typ: domain.ErrorTypeValidation,
	}

	ErrFileNameInvalid = domain.ApplicationError{
		Msg: "invalid file name",
		Typ: domain.ErrorTypeValidation,
	}

	ErrFileNameTooLong = domain.ApplicationError{
		Msg: "file name is too long",
		Typ: domain.ErrorTypeValidation,
	}

	ErrFileNameConflict = domain.ApplicationError{
		Msg: "file name already exists in target directory",
		Typ: domain.ErrorTypeConflict,
	}

	ErrFileTooLarge = domain.ApplicationError{
		Msg: "file is too large",
		Typ: domain.ErrorTypeValidation,
	}

	ErrFileTypeNotAllowed = domain.ApplicationError{
		Msg: "file type is not allowed",
		Typ: domain.ErrorTypeValidation,
	}

	ErrInvalidStatus = domain.ApplicationError{
		Msg: "invalid status transition",
		Typ: domain.ErrorTypeLogical,
	}

	ErrInexistentParent = domain.ApplicationError{
		Msg: "target parent does not exist",
		Typ: domain.ErrorTypeLogical,
	}
)
