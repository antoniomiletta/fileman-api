package domain

type ErrorType string

const (
	ErrorTypeLogical       ErrorType = "LOGICAL"       // core application logic was infringed, resulting in a logical impossibility
	ErrorTypeAuthorization ErrorType = "AUTHORIZATION" // not authorized to perform action
	ErrorTypeValidation    ErrorType = "VALIDATION"    // data or its format is invalid.
	ErrorTypeConflict      ErrorType = "CONFLICT"      // data conflicts with existing data
	ErrorTypeRetrieval     ErrorType = "RETRIEVAL"     // failed to retrieve requested data
)

// StatusCarrier should be implemented by all application-defined errors.
// Necessary for domain to http layer mapping.
type StatusCarrier interface {
	Error() string
	Type() ErrorType
}

type ApplicationError struct {
	Msg string
	Typ ErrorType
}

func (e ApplicationError) Error() string {
	return e.Msg
}

func (e ApplicationError) Type() ErrorType {
	return e.Typ
}
