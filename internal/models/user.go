package models

import "github.com/gofrs/uuid"

type User struct {
	ID             uuid.UUID
	Username       string
	Email          string
	HashedPassword string
}
