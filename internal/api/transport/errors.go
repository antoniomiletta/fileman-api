package transport

import (
	"errors"
)

var (
	ErrMalformedJSON = errors.New("invalid JSON payload")
	ErrInvalidQuery  = errors.New("invalid query parameters")
)
