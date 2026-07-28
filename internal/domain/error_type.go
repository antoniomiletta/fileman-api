package domain

type ErrorType string

const (
	ErrorTypeLogical       ErrorType = "LOGICAL"       // when core application logic is infringed, resulting in a logical impossibility
	ErrorTypeAuthorization ErrorType = "AUTHORIZATION" // when not authorized to perform action
	ErrorTypeValidation    ErrorType = "VALIDATION"    // when data or its format is invalid.
	ErrorTypeConflict      ErrorType = "CONFLICT"      // when data conflicts with existing data
	ErrorTypeRetrieval     ErrorType = "RETRIEVAL"     // when failing to retrieve requested data
)

func (t ErrorType) String() string { return string(t) }

type ApplicationError struct {
	Msg string
	Typ ErrorType
}

func (e ApplicationError) Error() string {
	return e.Msg
}

func (e ApplicationError) Type() string {
	return e.Typ.String()
}
