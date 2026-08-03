package authenticator

import "github.com/antoniomiletta/fileman/internal/domain"

var (
	ErrUnexpectedAlg = domain.ApplicationError{
		Msg: "unexpected signing algorithm",
		Typ: domain.ErrorTypeAuthentication,
	}
	ErrInvalidToken = domain.ApplicationError{
		Msg: "invalid authentication token",
		Typ: domain.ErrorTypeAuthentication,
	}
)
