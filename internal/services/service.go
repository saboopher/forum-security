package service

import (
	repository "forum/internal/repositories"
	"forum/internal/services/authservice"
	"forum/internal/services/postservice"
)

type Service struct {
	AuthService authservice.AuthServiceI
	PostService postservice.PostServiceI
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		AuthService: authservice.NewAuthService(repo.AuthRepo),
		PostService: postservice.NewPostService(repo.PostRepo, repo.AuthRepo),
	}
}
