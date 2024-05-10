package posthandler

import "forum/internal/services/postservice"

type PostHandler struct {
	PostService postservice.PostServiceI
}

func NewPostHandler(postService postservice.PostServiceI) *PostHandler {
	return &PostHandler{
		PostService: postService,
	}
}
