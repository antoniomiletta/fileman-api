package file

import "errors"

var (
	ErrFileNotFound     = errors.New("file not found")
	ErrFileNameRequired = errors.New("file name is required")
	ErrFileNameInvalid  = errors.New("invalid file name")
	ErrFileNameTooLong  = errors.New("file name is too long")
	ErrFileNameConflict = errors.New("file name already exists in target directory")
	ErrFileTooLarge     = errors.New("file is too large")
	ErrInvalidStatus    = errors.New("invalid status transition")
)
