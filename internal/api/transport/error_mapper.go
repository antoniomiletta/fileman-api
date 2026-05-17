package transport

import (
	"errors"
	"net/http"

	"github.com/antoniomiletta/fileman/internal/domain/auth"
	"github.com/antoniomiletta/fileman/internal/domain/file"
	"github.com/antoniomiletta/fileman/internal/domain/folder"
)

type APIError struct {
	Status  int
	Message string
}

var errorMap = map[error]int{
	auth.ErrForbidden:           http.StatusUnauthorized,
	auth.ErrInvalidToken:        http.StatusUnauthorized,
	auth.ErrUserNotFound:        http.StatusNotFound,
	auth.ErrInvalidCredentials:  http.StatusUnauthorized,
	auth.ErrCredentialsRequired: http.StatusBadRequest,
	auth.ErrInvalidEmail:        http.StatusBadRequest,
	auth.ErrEmailTaken:          http.StatusConflict,
	auth.ErrPasswordTooWeak:     http.StatusBadRequest,

	file.ErrFileNotFound:     http.StatusNotFound,
	file.ErrFileNameRequired: http.StatusBadRequest,
	file.ErrFileNameInvalid:  http.StatusBadRequest,
	file.ErrFileNameConflict: http.StatusConflict,
	file.ErrFileTooLarge:     http.StatusBadRequest,
	file.ErrInvalidStatus:    http.StatusBadRequest,

	folder.ErrFolderNotFound:         http.StatusNotFound,
	folder.ErrFolderNameRequired:     http.StatusBadRequest,
	folder.ErrFolderNameInvalid:      http.StatusBadRequest,
	folder.ErrFolderNameConflict:     http.StatusConflict,
	folder.ErrCannotMoveToDescendant: http.StatusUnprocessableEntity,
	folder.ErrCannotMoveRootFolder:   http.StatusUnprocessableEntity,

	ErrMalformedJSON: http.StatusBadRequest,
	ErrInvalidQuery:  http.StatusBadRequest,
}

func MapError(err error) APIError {
	if status, exists := errorMap[err]; exists {
		return APIError{
			Message: err.Error(),
			Status:  status,
		}
	}

	// error.Is() for wrapped errors
	for error, status := range errorMap {
		if errors.Is(err, error) {
			return APIError{
				Message: error.Error(),
				Status:  status,
			}
		}
	}

	return APIError{
		Message: "internal server error",
		Status:  http.StatusInternalServerError,
	}
}
