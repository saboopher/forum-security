package handler

import (
	"net/http"
)

func (h *Handler) Routes() *http.ServeMux {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// handlers passed through middleware
	loginHandler := http.HandlerFunc(h.AuthHandler.SignIn)          // +
	logoutHandler := http.HandlerFunc(h.AuthHandler.LogOut)         // +
	postCreateHandler := http.HandlerFunc(h.PostHandler.PostCreate) // +
	postUpdateHandler := http.HandlerFunc(h.PostHandler.PostUpdate)
	errorsHandler := http.HandlerFunc(errorHandler)

	postGetHandler := http.HandlerFunc(h.PostHandler.PostGet)           // +
	postGetAllHandler := http.HandlerFunc(h.PostHandler.PostGetAll)     // +
	postGetMyPostsHandler := http.HandlerFunc(h.PostHandler.GetMyPosts) // +
	postCommentCreateHandler := http.HandlerFunc(h.PostHandler.CommentCreate)
	postCommentUpdateHandler := http.HandlerFunc(h.PostHandler.CommentUpdate)

	mux.HandleFunc("/register", h.AuthHandler.RegisterUser)
	mux.Handle("/signin", h.SessionMiddleware(loginHandler))
	mux.Handle("/logout", h.SessionMiddleware(logoutHandler))
	mux.Handle("/post/", h.SessionMiddleware(postGetAllHandler))
	mux.Handle("/post/get", h.SessionMiddleware(postGetHandler))
	mux.Handle("/post/create", h.SessionMiddleware(postCreateHandler))
	mux.Handle("/post/update", h.SessionMiddleware(postUpdateHandler))
	mux.Handle("/post/myposts", h.SessionMiddleware(postGetMyPostsHandler))
	mux.Handle("/comment/create", h.SessionMiddleware(postCommentCreateHandler))
	mux.Handle("/comment/update", h.SessionMiddleware(postCommentUpdateHandler))
	mux.Handle("/", h.ErrorMiddleware(errorsHandler))
	return mux
}

/*
+ /post/get/all (all posts, filters)
+ /post/update
+ /post/get
+ /post/myposts

/comment/create
/comment/delete
/like/create
/like/delete
/session/refresh

*/
