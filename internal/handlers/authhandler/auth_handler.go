package authhandler

import "forum/internal/services/authservice"

type AuthHandler struct {
	AuthService authservice.AuthServiceI
}

func NewAuthHandler(authService authservice.AuthServiceI) *AuthHandler {
	return &AuthHandler{
		AuthService: authService,
	}
}
