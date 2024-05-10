package schemas

import "github.com/gofrs/uuid"

type UpdateUser struct {
	Username string
	Email    string
	Password string
}

type CreateUser struct {
	UpdateUser
	PasswordConfirm string
}

type ShowUser struct {
	ID       uuid.UUID
	Username string
	Email    string
}

type AuthUser struct {
	Email    string
	Password string
}
