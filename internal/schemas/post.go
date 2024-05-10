package schemas

import (
	"forum/internal/models"
	"time"

	"github.com/gofrs/uuid"
)

type CreatePost struct {
	ID         uuid.UUID
	Title      string
	Body       string
	Image      string
	Likes      int
	Dislikes   int
	Categories []string
}

type UpdatePost struct {
	PostID uuid.UUID
	CreatePost
}

type GetPostResponse struct {
	Username   string
	PostID     uuid.UUID
	CreatedAt  time.Time
	UpdatedAt  time.Time
	PostTitle  string
	PostBody   string
	PostImage  string
	Categories []string
	Likes      int
	Dislikes   int
	Comments   []*models.Comment
}

type Data struct {
	Session             *models.Session
	Post                *GetPostResponse
	Posts               []*GetPostResponse
	Categories          []*Category
	TemplateCommentForm *TemplateCommentForm
}

type PostDataForErr struct {
	Title string
	Body  string
}

type PostErrors struct {
	Title    string
	Body     string
	Category string
}

type TemplatePostForm struct {
	PostErrors     PostErrors
	PostDataForErr PostDataForErr
}

type CreatePostForm struct {
	TemplatePostForm *TemplatePostForm
	Session          *models.Session
	Categories       []*Category
}

type UpdatePostForm struct{}
