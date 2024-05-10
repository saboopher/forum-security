package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type Comment struct {
	ID          uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Description string
	PostID      uuid.UUID
	UserID      uuid.UUID
	UserName    string
	Likes       int
	Dislikes    int
}
