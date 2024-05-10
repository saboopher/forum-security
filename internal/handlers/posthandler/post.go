package posthandler

import (
	"forum/internal/exceptions"
	. "forum/internal/models"
	"forum/internal/schemas"
	"forum/pkg/cust_encoders"
	"forum/pkg/validator"
	"html/template"
	"net/http"

	"github.com/gofrs/uuid"
)

func (ah *PostHandler) PostCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet || r.Method == http.MethodPost {
		sessionValue := r.Context().Value("session")

		if sessionValue == nil {
			dataErr := exceptions.NewInternalServerError()
			params := cust_encoders.EncodeParams(dataErr)
			http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
			return
		}

		var toOk Session
		_ = toOk

		session, ok := sessionValue.(Session)
		if !ok {
			dataErr := exceptions.NewInternalServerError()
			params := cust_encoders.EncodeParams(dataErr)
			http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
			return
		}
		categories, err := ah.PostService.GetAllCategories()
		if err != nil {
			params := cust_encoders.EncodeParams(err)
			http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
		}
		createPostForm := &schemas.CreatePostForm{}
		createPostForm.Session = &session
		createPostForm.Categories = categories
		if r.Method == http.MethodPost {
			if err := r.ParseForm(); err != nil {
				dataErr := exceptions.NewBadRequestError("Invalid form values")
				params := cust_encoders.EncodeParams(dataErr)
				http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
			} else {
				title := r.FormValue("title")
				body := r.FormValue("body")
				categoriesInput := r.PostForm["categories"]

				titleOk, msgTitle := validator.ValidatePostTitle(title)
				bodyOk, msgBody := validator.ValidatePostBody(body)
				categoryOk, msgCategory := validator.ValidateCategoryLen(categoriesInput, categories)

				if !titleOk || !bodyOk || !categoryOk {
					createPostForm.TemplatePostForm = &schemas.TemplatePostForm{}
					if !titleOk {
						createPostForm.TemplatePostForm.PostErrors.Title = msgTitle
					}

					if !bodyOk {
						createPostForm.TemplatePostForm.PostErrors.Body = msgBody
					}

					if !categoryOk {
						createPostForm.TemplatePostForm.PostErrors.Category = msgCategory
					}

					createPostForm.TemplatePostForm.PostDataForErr.Title = title
					createPostForm.TemplatePostForm.PostDataForErr.Body = body
				}
				if titleOk && bodyOk && categoryOk {
					post := schemas.CreatePost{
						Title:      title,
						Body:       body,
						Image:      "/dgf/dfg",
						Likes:      0,
						Dislikes:   0,
						Categories: categoriesInput,
					}

					err = ah.PostService.CreatePost(session.UserID, post)
					if err != nil {
						params := cust_encoders.EncodeParams(err)
						http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
						return
					}
					http.Redirect(w, r, "/post", http.StatusSeeOther)
					return
				}
			}
		}

		t, err := template.ParseFiles("ui/templates/create.html")
		if err != nil {
			dataErr := exceptions.NewInternalServerError()
			params := cust_encoders.EncodeParams(dataErr)
			http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
			return
		}
		t.Execute(w, createPostForm)
		return
	}
	dataErr := exceptions.NewStatusMethodNotAllowed()
	params := cust_encoders.EncodeParams(dataErr)
	http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
	return
}

func (ah *PostHandler) PostUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		dataErr := exceptions.NewStatusMethodNotAllowed()
		params := cust_encoders.EncodeParams(dataErr)
		http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
		return
	}

	sessionValue := r.Context().Value("session")

	if sessionValue == nil {
		dataErr := exceptions.NewInternalServerError()
		params := cust_encoders.EncodeParams(dataErr)
		http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
		return
	}

	session, ok := sessionValue.(Session)
	if !ok {
		dataErr := exceptions.NewInternalServerError()
		params := cust_encoders.EncodeParams(dataErr)
		http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
		return
	}
	userID := session.UserID
	var postID string
	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			dataErr := exceptions.NewBadRequestError("Invalid form values")
			params := cust_encoders.EncodeParams(dataErr)
			http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
			return
		} else {
			postID = r.FormValue("post_id")
			if postID == "" {
				dataErr := exceptions.NewBadRequestError("Post ID is not given")
				params := cust_encoders.EncodeParams(dataErr)
				http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
				return
			}

			values := r.URL.Query()
			like := values.Get("like")
			dislike := values.Get("dislike")
			var isLike bool
			if like != "" && dislike != "" {
				dataErr := exceptions.NewBadRequestError("Like or dislike should take place")
				params := cust_encoders.EncodeParams(dataErr)
				http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
				return
			}
			if like != "" {
				isLike = true
			}
			postIDArg, err := uuid.FromString(postID)
			if err != nil {
				dataErr := exceptions.NewBadRequestError("Invalid post ID")
				params := cust_encoders.EncodeParams(dataErr)
				http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
				return
			}
			getPostResponse, err := ah.PostService.GetPost(postIDArg)
			if err != nil {
				params := cust_encoders.EncodeParams(err)
				http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
				return
			}

			postUpdate := schemas.UpdatePost{
				PostID: postIDArg,
				CreatePost: schemas.CreatePost{
					Title:    getPostResponse.PostTitle,
					Body:     getPostResponse.PostBody,
					Image:    getPostResponse.PostImage,
					Likes:    getPostResponse.Likes,
					Dislikes: getPostResponse.Dislikes,
				},
			}
			vote, err := ah.PostService.GetVote(postIDArg, userID, "post")

			voteCreate := schemas.CreateVote{
				ShowVote: schemas.ShowVote{
					VoteID: uuid.Must(uuid.NewV4()),
					UserID: userID,
					PostID: postIDArg,
					Binary: -1,
				},
			}
			var binary int
			if err != nil {
				if isLike {
					postUpdate.Likes++
					binary = 1
				} else {
					postUpdate.Dislikes++
				}
				err := ah.PostService.UpdatePost(userID, postUpdate)
				if err != nil {
					params := cust_encoders.EncodeParams(err)
					http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
					return
				}
				voteCreate.Binary = binary
				err = ah.PostService.CreateVote(voteCreate, "post")
				if err != nil {
					params := cust_encoders.EncodeParams(err)
					http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
					return
				}
			} else {
				var isDelete bool
				var unvote bool
				if isLike && vote.Binary == 0 {
					postUpdate.Likes++
					if postUpdate.Dislikes != 0 {
						postUpdate.Dislikes--
					}
					isDelete = true
					binary = 1
				} else if isLike && vote.Binary == 1 {
					postUpdate.Likes--
					isDelete = true
					unvote = true
				} else if !isLike && vote.Binary == 0 {
					postUpdate.Dislikes--
					isDelete = true
					unvote = true
				} else if !isLike && vote.Binary == 1 {
					postUpdate.Dislikes++
					if postUpdate.Likes != 0 {
						postUpdate.Likes--
					}
					isDelete = true
				}
				if isDelete {
					err := ah.PostService.DeleteVote(vote.VoteID)
					if err != nil {
						params := cust_encoders.EncodeParams(err)
						http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
						return
					}
					err = ah.PostService.UpdatePost(userID, postUpdate)
					if err != nil {
						params := cust_encoders.EncodeParams(err)
						http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
						return
					}
					if !unvote {
						voteCreate.Binary = binary
						voteCreate.VoteID = uuid.Must(uuid.NewV4())
						err = ah.PostService.CreateVote(voteCreate, "post")
						if err != nil {
							params := cust_encoders.EncodeParams(err)
							http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
							return
						}
					}
				}
			}
		}

		http.Redirect(w, r, "/post/get?post_id="+postID, http.StatusSeeOther)
	}
}

func (ah *PostHandler) PostGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		dataErr := exceptions.NewStatusMethodNotAllowed()
		params := cust_encoders.EncodeParams(dataErr)
		http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
		return
	}

	var session *Session
	sessionValue := r.Context().Value("session")
	if sessionValue != nil {
		sessionVal, ok := sessionValue.(Session)
		if ok {
			session = &sessionVal
		}
	}

	postIDStr := r.URL.Query().Get("post_id")
	postID, err := uuid.FromString(postIDStr)
	if err != nil {
		dataErr := exceptions.NewBadRequestError("Invalid Post ID")
		params := cust_encoders.EncodeParams(dataErr)
		http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
		return
	}

	getPostResponse, err := ah.PostService.GetPost(postID)
	if err != nil {
		params := cust_encoders.EncodeParams(err)
		http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
		return
	}
	resp := &schemas.Data{
		Session: session,
		Post:    getPostResponse,
	}

	t, err := template.ParseFiles("ui/templates/post.html") // different html
	if err != nil {
		dataErr := exceptions.NewInternalServerError()
		params := cust_encoders.EncodeParams(dataErr)
		http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
		return
	}
	t.Execute(w, resp)
}

func (ah *PostHandler) PostGetAll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		dataErr := exceptions.NewStatusMethodNotAllowed()
		params := cust_encoders.EncodeParams(dataErr)
		http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
		return
	}
	var session *Session
	var categories []*schemas.Category
	sessionValue := r.Context().Value("session")
	if sessionValue != nil {
		sessionVal, ok := sessionValue.(Session)
		if ok {
			session = &sessionVal
		}
	}

	category := r.URL.Query().Get("category")
	likedFilter := r.URL.Query().Get("liked")

	if likedFilter != "" && likedFilter != "1" {
		dataErr := exceptions.NewBadRequestError("Bad Parameter Request")
		params := cust_encoders.EncodeParams(dataErr)
		http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
		return
	}

	posts, err := ah.PostService.GetPostsAll(category)
	if err != nil {
		params := cust_encoders.EncodeParams(err)
		http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
		return
	}
	if likedFilter != "" {
		postsResp, err := ah.PostService.GetLikedPosts(session.UserID, posts)
		if err != nil {
			params := cust_encoders.EncodeParams(err)
			http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
			return
		}
		posts = postsResp
	} else {
		categoriesResp, err := ah.PostService.GetAllCategories()
		if err != nil {
			params := cust_encoders.EncodeParams(err)
			http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
			return
		}
		categories = categoriesResp
	}
	resp := &schemas.Data{
		Session:    session,
		Posts:      posts,
		Categories: categories,
	}

	t, err := template.ParseFiles("ui/templates/home.html")
	if err != nil {
		dataErr := exceptions.NewInternalServerError()
		params := cust_encoders.EncodeParams(dataErr)
		http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
		return
	}
	t.Execute(w, resp)
}

func (ah *PostHandler) GetMyPosts(w http.ResponseWriter, r *http.Request) {
	sessionValue := r.Context().Value("session")

	if sessionValue == nil {
		dataErr := exceptions.NewInternalServerError()
		params := cust_encoders.EncodeParams(dataErr)
		http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
		return
	}

	session, ok := sessionValue.(Session)
	if !ok {
		dataErr := exceptions.NewInternalServerError()
		params := cust_encoders.EncodeParams(dataErr)
		http.Redirect(w, r, "/?p"+params, http.StatusSeeOther)
		return
	}

	getPostAllResponse, err := ah.PostService.GetMyPosts(session.UserID)
	if err != nil {
		params := cust_encoders.EncodeParams(err)
		http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
		return
	}

	resp := &schemas.Data{
		Session: &session,
		Posts:   getPostAllResponse,
	}

	t, err := template.ParseFiles("ui/templates/home.html")
	if err != nil {
		dataErr := exceptions.NewInternalServerError()
		params := cust_encoders.EncodeParams(dataErr)
		http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
		return
	}
	t.Execute(w, resp)
}

// comments are going to be here

func (ah *PostHandler) CommentCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		dataErr := exceptions.NewStatusMethodNotAllowed()
		params := cust_encoders.EncodeParams(dataErr)
		http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
		return
	}
	sessionValue := r.Context().Value("session")

	if sessionValue == nil {
		dataErr := exceptions.NewInternalServerError()
		params := cust_encoders.EncodeParams(dataErr)
		http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
		return
	}

	var toOk Session
	_ = toOk

	session, ok := sessionValue.(Session)
	if !ok {
		dataErr := exceptions.NewInternalServerError()
		params := cust_encoders.EncodeParams(dataErr)
		http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
		return
	}
	createCommentForm := &schemas.CreateCommentForm{}
	createCommentForm.Session = &session
	if err := r.ParseForm(); err != nil {
		dataErr := exceptions.NewBadRequestError("Invalid form values")
		params := cust_encoders.EncodeParams(dataErr)
		http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
	} else {
		content := r.FormValue("content")
		postIDStr := r.Form.Get("post_id")
		postID, err := uuid.FromString(postIDStr)
		if err != nil {
			dataErr := exceptions.NewBadRequestError("Invalid Post ID")
			params := cust_encoders.EncodeParams(dataErr)
			http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
			return
		}
		commentOk, msgComment := validator.ValidatePostComment(content)
		if !commentOk {
			dataErr := exceptions.NewValidationError(msgComment)
			params := cust_encoders.EncodeParams(dataErr)
			http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
			return
		}
		comment := schemas.CreateComment{
			Content:  content,
			PostID:   postID,
			UserID:   session.UserID,
			Likes:    0,
			Dislikes: 0,
		}

		err = ah.PostService.CreateComment(comment)
		if err != nil {
			params := cust_encoders.EncodeParams(err)
			http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
			return
		}
		http.Redirect(w, r, "/post/get?post_id="+postIDStr, http.StatusSeeOther)
		return
	}

	t, err := template.ParseFiles("ui/templates/post.html")
	if err != nil {
		dataErr := exceptions.NewInternalServerError()
		params := cust_encoders.EncodeParams(dataErr)
		http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
		return
	}
	t.Execute(w, createCommentForm)
}

func (ah *PostHandler) CommentUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		dataErr := exceptions.NewStatusMethodNotAllowed()
		params := cust_encoders.EncodeParams(dataErr)
		http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
		return
	}

	sessionValue := r.Context().Value("session")

	if sessionValue == nil {
		dataErr := exceptions.NewInternalServerError()
		params := cust_encoders.EncodeParams(dataErr)
		http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
		return
	}

	session, ok := sessionValue.(Session)
	if !ok {
		dataErr := exceptions.NewInternalServerError()
		params := cust_encoders.EncodeParams(dataErr)
		http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
		return
	}
	userID := session.UserID
	var commentID string
	var postID string
	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			dataErr := exceptions.NewBadRequestError("Invalid form values")
			params := cust_encoders.EncodeParams(dataErr)
			http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
			return
		} else {
			postID = r.FormValue("post_id")
			commentID = r.FormValue("comment_id")
			if commentID == "" || postID == "" {
				dataErr := exceptions.NewBadRequestError("Post ID and Comment ID should be present")
				params := cust_encoders.EncodeParams(dataErr)
				http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
				return
			}
			values := r.URL.Query()
			like := values.Get("like")
			dislike := values.Get("dislike")
			var isLike bool
			if like != "" && dislike != "" {
				dataErr := exceptions.NewBadRequestError("Like or Dislike should take place")
				params := cust_encoders.EncodeParams(dataErr)
				http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
				return
			}
			if like != "" {
				isLike = true
			}
			commentIDArg, err := uuid.FromString(commentID)
			if err != nil {
				dataErr := exceptions.NewBadRequestError("Invalid Comment ID")
				params := cust_encoders.EncodeParams(dataErr)
				http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
				return
			}
			comment, err := ah.PostService.GetComment(commentIDArg)
			if err != nil {
				params := cust_encoders.EncodeParams(err)
				http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
				return
			}
			commentUpdate := schemas.UpdateComment{
				ID:       commentIDArg,
				Likes:    comment.Likes,
				Dislikes: comment.Dislikes,
			}
			vote, err := ah.PostService.GetVote(commentIDArg, userID, "comment")

			voteCreate := schemas.CreateVote{
				ShowVote: schemas.ShowVote{
					VoteID:    uuid.Must(uuid.NewV4()),
					UserID:    userID,
					CommentID: commentIDArg,
					Binary:    -1,
				},
			}
			var binary int
			if err != nil {
				if isLike {
					commentUpdate.Likes++
					binary = 1
				} else {
					commentUpdate.Dislikes++
				}
				err := ah.PostService.UpdateComment(userID, commentUpdate)
				if err != nil {
					params := cust_encoders.EncodeParams(err)
					http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
					return
				}
				voteCreate.Binary = binary
				err = ah.PostService.CreateVote(voteCreate, "comment")
				if err != nil {
					params := cust_encoders.EncodeParams(err)
					http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
					return
				}
			} else {
				var isDelete bool
				var unvote bool
				if isLike && vote.Binary == 0 {
					commentUpdate.Likes++
					if commentUpdate.Dislikes != 0 {
						commentUpdate.Dislikes--
					}
					isDelete = true
					binary = 1
				} else if isLike && vote.Binary == 1 {
					commentUpdate.Likes--
					isDelete = true
					unvote = true
				} else if !isLike && vote.Binary == 0 {
					commentUpdate.Dislikes--
					isDelete = true
					unvote = true
				} else if !isLike && vote.Binary == 1 {
					commentUpdate.Dislikes++
					if commentUpdate.Likes != 0 {
						commentUpdate.Likes--
					}
					isDelete = true
				}
				if isDelete {
					err := ah.PostService.DeleteVote(vote.VoteID)
					if err != nil {
						params := cust_encoders.EncodeParams(err)
						http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
						return
					}
					err = ah.PostService.UpdateComment(userID, commentUpdate)
					if err != nil {
						params := cust_encoders.EncodeParams(err)
						http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
						return
					}
					if !unvote {
						voteCreate.Binary = binary
						voteCreate.VoteID = uuid.Must(uuid.NewV4())
						err = ah.PostService.CreateVote(voteCreate, "comment")
						if err != nil {
							params := cust_encoders.EncodeParams(err)
							http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
							return
						}
					}

				}
			}
		}

		http.Redirect(w, r, "/post/get?post_id="+postID, http.StatusSeeOther)
	}
}
