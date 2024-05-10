package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type Post struct {
	ID         uuid.UUID
	CreatedAt  time.Time
	UpdatedAt  time.Time
	UserId     uuid.UUID
	Title      string
	Body       string
	Image      string
	Likes      int
	Dislikes   int
	Categories []*Category
}
