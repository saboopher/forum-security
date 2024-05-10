package models

import (
	"github.com/gofrs/uuid"
)

type Vote struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	PostID    uuid.UUID
	CommentID uuid.UUID
	Binary    int
}
