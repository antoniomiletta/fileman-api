package domain

import "github.com/google/uuid"

type User struct {
	ID       uuid.UUID
	Email    string
	Password string
}

type NewUserParams struct {
	Email    string
	Password string
}

func NewUser(p NewUserParams) User {
	return User{
		ID:       uuid.New(),
		Email:    p.Email,
		Password: p.Password,
	}
}
