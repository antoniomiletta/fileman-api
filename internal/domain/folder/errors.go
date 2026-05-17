package folder

import "errors"

var (
	ErrFolderNotFound         = errors.New("folder not found")
	ErrFolderNameRequired     = errors.New("folder name is required")
	ErrFolderNameInvalid      = errors.New("invalid folder name")
	ErrFolderNameTooLong      = errors.New("folder name is too long")
	ErrFolderNameConflict     = errors.New("folder name already exists in target directory")
	ErrCannotMoveToDescendant = errors.New("folder cannot be moved into a descendant")
	ErrCannotMoveRootFolder   = errors.New("root folder cannot be moved")
)
