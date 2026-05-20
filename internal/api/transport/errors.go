package transport

type ErrorType string

var (
	ErrorTypeParsing ErrorType = "PARSING"
)

func (t ErrorType) String() string { return string(t) }

type TransportError struct {
	msg string
	typ ErrorType
}

func (e TransportError) Error() string {
	return e.msg
}

func (e TransportError) Type() ErrorType {
	return e.typ
}

var (
	ErrMalformedJSON = TransportError{
		msg: "invalid JSON payload",
		typ: ErrorTypeParsing,
	}
	ErrInvalidPathParam = TransportError{
		msg: "invalid path parameter",
		typ: ErrorTypeParsing,
	}
	ErrMalformedToken = TransportError{
		msg: "invalid authentication token format",
		typ: ErrorTypeParsing,
	}
)
