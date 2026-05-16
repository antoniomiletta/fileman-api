package transport

import (
	"errors"
	"net/http"

	"github.com/antoniomiletta/fileman/internal/domain"
)

type APIError struct {
	Status  int
	Message string
}

var errorMap = map[error]int{
	domain.ErrForbidden:           http.StatusUnauthorized,
	domain.ErrInvalidToken:        http.StatusUnauthorized,
	domain.ErrUserNotFound:        http.StatusNotFound,
	domain.ErrInvalidCredentials:  http.StatusUnauthorized,
	domain.ErrCredentialsRequired: http.StatusBadRequest,
	domain.ErrInvalidEmail:        http.StatusBadRequest,
	domain.ErrEmailTaken:          http.StatusConflict,
	domain.ErrPasswordTooWeak:     http.StatusBadRequest,

	domain.ErrFileNotFound:     http.StatusNotFound,
	domain.ErrFileNameRequired: http.StatusBadRequest,
	domain.ErrFileNameInvalid:  http.StatusBadRequest,
	domain.ErrFileNameConflict: http.StatusConflict,
	domain.ErrFileTooLarge:     http.StatusBadRequest,
	domain.ErrInvalidStatus:    http.StatusBadRequest,

	domain.ErrFolderNotFound:         http.StatusNotFound,
	domain.ErrFolderNameRequired:     http.StatusBadRequest,
	domain.ErrFolderNameInvalid:      http.StatusBadRequest,
	domain.ErrFolderNameConflict:     http.StatusConflict,
	domain.ErrCannotMoveToDescendant: http.StatusUnprocessableEntity,
	domain.ErrCannotMoveRootFolder:   http.StatusUnprocessableEntity,

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
