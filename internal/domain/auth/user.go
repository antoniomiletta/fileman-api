package auth

import "github.com/google/uuid"

type User struct {
	ID       uuid.UUID
	Email    string
	Password string
}

type NewUserParams struct {
	ID       uuid.UUID
	Email    string
	Password string
}

func NewUser(p NewUserParams) User {
	return User{
		ID:       p.ID,
		Email:    p.Email,
		Password: p.Password,
	}
}
