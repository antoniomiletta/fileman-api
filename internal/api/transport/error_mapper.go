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

type StatusCarrier interface {
	Error() string
	Type() string
}

func MapError(err error) APIError {
	if carrier, ok := AsType[StatusCarrier](err); ok {
		return APIError{
			Message: carrier.Error(),
			Status:  MapHTTPStatus(carrier.Type()),
		}
	}

	return APIError{
		Message: "internal server error",
		Status:  http.StatusInternalServerError,
	}
}

func MapHTTPStatus(errType string) int {
	switch errType {
	case domain.ErrorTypeLogical.String():
		return http.StatusUnprocessableEntity

	case domain.ErrorTypeAuthorization.String():
		return http.StatusUnauthorized

	case domain.ErrorTypeValidation.String():
		return http.StatusBadRequest

	case domain.ErrorTypeConflict.String():
		return http.StatusConflict

	case domain.ErrorTypeRetrieval.String():
		return http.StatusNotFound

	case ErrorTypeParsing.String():
		return http.StatusBadRequest

	default:
		return http.StatusInternalServerError
	}
}

func AsType[T any](err error) (T, bool) {
	var target T
	if errors.As(err, &target) {
		return target, true
	}
	return target, false
}
