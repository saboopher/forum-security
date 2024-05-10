package repository

import (
	"database/sql"
	"forum/internal/repositories/authrepo"
	"forum/internal/repositories/postrepo"
)

type Repository struct {
	AuthRepo authrepo.AuthRepoI
	PostRepo postrepo.PostRepoI
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		AuthRepo: authrepo.NewAuthRepo(db),
		PostRepo: postrepo.NewPostRepo(db),
	}
}
