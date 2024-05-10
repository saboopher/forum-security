package schemas

import (
	"forum/internal/models"

	"github.com/gofrs/uuid"
)

type CreateComment struct {
	ID       uuid.UUID
	Content  string
	PostID   uuid.UUID
	UserID   uuid.UUID
	Likes    int
	Dislikes int
}

type CreateCommentForm struct {
	TemplateCommentForm *TemplateCommentForm
	Session             *models.Session
}

type TemplateCommentForm struct {
	CommentErrors     string
	CommentDataForErr string
}

type UpdateComment struct {
	ID       uuid.UUID
	Likes    int
	Dislikes int
}

type ShowComment struct {
	ID       uuid.UUID
	Content  string
	PostID   uuid.UUID
	UserID   uuid.UUID
	Likes    int
	Dislikes int
}
