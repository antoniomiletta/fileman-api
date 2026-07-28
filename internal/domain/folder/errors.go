package folder

import "github.com/antoniomiletta/fileman/internal/domain"

var (
	ErrFolderNotFound = domain.ApplicationError{
		Msg: "folder not found",
		Typ: domain.ErrorTypeRetrieval,
	}

	ErrFolderNameRequired = domain.ApplicationError{
		Msg: "folder name is required",
		Typ: domain.ErrorTypeValidation,
	}

	ErrFolderNameInvalid = domain.ApplicationError{
		Msg: "invalid folder name",
		Typ: domain.ErrorTypeValidation,
	}

	ErrFolderNameTooLong = domain.ApplicationError{
		Msg: "folder name is too long",
		Typ: domain.ErrorTypeValidation,
	}

	ErrFolderNameConflict = domain.ApplicationError{
		Msg: "folder name already exists in target directory",
		Typ: domain.ErrorTypeConflict,
	}

	ErrAlreadyInDestination = domain.ApplicationError{
		Msg: "folder is already in the destination",
		Typ: domain.ErrorTypeLogical,
	}

	ErrCannotMoveIntoItself = domain.ApplicationError{
		Msg: "folder cannot be moved into itself",
		Typ: domain.ErrorTypeLogical,
	}

	ErrCannotMoveIntoDescendant = domain.ApplicationError{
		Msg: "folder cannot be moved into a descendant",
		Typ: domain.ErrorTypeLogical,
	}

	ErrCannotCreateNewRootFolder = domain.ApplicationError{
		Msg: "root folder cannot be duplicated",
		Typ: domain.ErrorTypeLogical,
	}

	ErrCannotMoveRootFolder = domain.ApplicationError{
		Msg: "root folder cannot be moved",
		Typ: domain.ErrorTypeLogical,
	}

	ErrCannotDeleteRootFolder = domain.ApplicationError{
		Msg: "root folder cannot be deleted",
		Typ: domain.ErrorTypeLogical,
	}
)
