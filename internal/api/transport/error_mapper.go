package transport

import (
	"errors"
	"net/http"

	"github.com/antoniomiletta/fileman/internal/domain"
)

type HttpError struct {
	Status  int
	Message string
}

// MapError maps application-defined errors that implement domain.StatusCarrier
// to an HttpError which can be presented to the HTTP layer.
// It extracts wrapped errors so implementation details don't leak to the client.
//
// If the given err does not implement domain.StatusCarrier, defaults to 500 Internal Server Error.
func MapError(err error) HttpError {
	if carrier, ok := errors.AsType[domain.StatusCarrier](err); ok {
		return HttpError{
			Message: carrier.Error(),
			Status:  MapHTTPStatus(carrier.Type()),
		}
	}

	return HttpError{
		Message: "internal server error",
		Status:  http.StatusInternalServerError,
	}
}

func MapHTTPStatus(errType domain.ErrorType) int {
	switch errType {
	case domain.ErrorTypeLogical:
		return http.StatusUnprocessableEntity

	case domain.ErrorTypeAuthentication:
		return http.StatusUnauthorized

	case domain.ErrorTypeAuthorization:
		return http.StatusForbidden

	case domain.ErrorTypeValidation:
		return http.StatusBadRequest

	case domain.ErrorTypeConflict:
		return http.StatusConflict

	case domain.ErrorTypeRetrieval:
		return http.StatusNotFound

	default:
		return http.StatusInternalServerError
	}
}
