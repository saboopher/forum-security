package models

import "github.com/gofrs/uuid"

type Category struct {
	ID   uuid.UUID
	Name string
}

type CreateCategoryPost struct {
	ID           uuid.UUID
	PostID       uuid.UUID
	CategoryName string
}
