package transport

import (
	"errors"
	"net/http"

	"github.com/antoniomiletta/fileman/internal/domain"
)

type ErrorMapping struct {
	Err    error
	Status int
}

var errorMap = []ErrorMapping{
	{domain.ErrForbidden, http.StatusUnauthorized},
	{domain.ErrInvalidCredentials, http.StatusUnauthorized},
	{domain.ErrCredentialsRequired, http.StatusBadRequest},
	{domain.ErrInvalidEmail, http.StatusBadRequest},
	{domain.ErrEmailTaken, http.StatusConflict},
	{domain.ErrPasswordTooWeak, http.StatusBadRequest},

	{domain.ErrFileNotFound, http.StatusNotFound},
	{domain.ErrFileNameRequired, http.StatusBadRequest},
	{domain.ErrFileNameInvalid, http.StatusBadRequest},
	{domain.ErrFileNameConflict, http.StatusConflict},
	{domain.ErrFileTooLarge, http.StatusBadRequest},
	{domain.ErrInvalidStatus, http.StatusBadRequest},

	{domain.ErrFolderNotFound, http.StatusNotFound},
	{domain.ErrFolderNameRequired, http.StatusBadRequest},
	{domain.ErrFolderNameInvalid, http.StatusBadRequest},
	{domain.ErrFolderNameConflict, http.StatusConflict},
	{domain.ErrCannotMoveToDescendant, http.StatusUnprocessableEntity},
	{domain.ErrCannotMoveRootFolder, http.StatusUnprocessableEntity},
}

type HttpError struct {
	Message string
	Status  int
}

func MapError(err error) HttpError {
	for _, mapping := range errorMap {
		if errors.Is(err, mapping.Err) {
			return HttpError{
				Message: mapping.Err.Error(),
				Status:  mapping.Status,
			}
		}
	}

	return HttpError{
		Message: "internal server error",
		Status:  http.StatusInternalServerError,
	}
}
