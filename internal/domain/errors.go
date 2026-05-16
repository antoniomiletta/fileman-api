package domain

import "errors"

var (
	ErrForbidden    = errors.New("not authorized")
	ErrInvalidToken = errors.New("invalid authentication token")
	ErrUserNotFound = errors.New("user not found")

	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrCredentialsRequired = errors.New("email and password are required")
	ErrEmailTaken          = errors.New("email already taken")
	ErrInvalidEmail        = errors.New("invalid email format")
	ErrPasswordTooWeak     = errors.New("password is too weak")

	ErrFileNotFound     = errors.New("file not found")
	ErrFileNameRequired = errors.New("file name is required")
	ErrFileNameInvalid  = errors.New("invalid file name")
	ErrFileNameConflict = errors.New("file name already exists in target directory")
	ErrFileTooLarge     = errors.New("file is too large")
	ErrInvalidStatus    = errors.New("invalid status transition")

	ErrFolderNotFound         = errors.New("folder not found")
	ErrFolderNameRequired     = errors.New("folder name is required")
	ErrFolderNameInvalid      = errors.New("invalid folder name")
	ErrFolderNameConflict     = errors.New("folder name already exists in target directory")
	ErrCannotMoveToDescendant = errors.New("folder cannot be moved into a descendant")
	ErrCannotMoveRootFolder   = errors.New("root folder cannot be moved")
)
