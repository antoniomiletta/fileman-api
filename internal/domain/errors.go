package domain

type ErrorType string

const (
	ErrorTypeLogical       ErrorType = "LOGICAL"
	ErrorTypeAuthorization ErrorType = "AUTHORIZATION"
	ErrorTypeValidation    ErrorType = "VALIDATION"
	ErrorTypeConflict      ErrorType = "CONFLICT"
	ErrorTypeRetrieval     ErrorType = "RETRIEVAL"
)

func (t ErrorType) String() string { return string(t) }

type DomainError struct {
	Msg string
	Typ ErrorType
}

func (e DomainError) Error() string {
	return e.Msg
}

func (e DomainError) Type() ErrorType {
	return e.Typ
}
