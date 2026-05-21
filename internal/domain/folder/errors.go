package folder

import "github.com/antoniomiletta/fileman/internal/domain"

var (
	ErrFolderNotFound = domain.DomainError{
		Msg: "folder not found",
		Typ: domain.ErrorTypeRetrieval,
	}

	ErrFolderNameRequired = domain.DomainError{
		Msg: "folder name is required",
		Typ: domain.ErrorTypeValidation,
	}

	ErrFolderNameInvalid = domain.DomainError{
		Msg: "invalid folder name",
		Typ: domain.ErrorTypeValidation,
	}

	ErrFolderNameTooLong = domain.DomainError{
		Msg: "folder name is too long",
		Typ: domain.ErrorTypeValidation,
	}

	ErrFolderNameConflict = domain.DomainError{
		Msg: "folder name already exists in target directory",
		Typ: domain.ErrorTypeConflict,
	}

	ErrAlreadyInDestination = domain.DomainError{
		Msg: "folder is already in the destination",
		Typ: domain.ErrorTypeLogical,
	}

	ErrCannotMoveIntoItself = domain.DomainError{
		Msg: "folder cannot be moved into itself",
		Typ: domain.ErrorTypeLogical,
	}

	ErrCannotMoveIntoDescendant = domain.DomainError{
		Msg: "folder cannot be moved into a descendant",
		Typ: domain.ErrorTypeLogical,
	}

	ErrCannotMoveRootFolder = domain.DomainError{
		Msg: "root folder cannot be moved",
		Typ: domain.ErrorTypeLogical,
	}

	ErrCannotDeleteRootFolder = domain.DomainError{
		Msg: "root folder cannot be deleted",
		Typ: domain.ErrorTypeLogical,
	}
)
