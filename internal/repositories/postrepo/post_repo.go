package postrepo

import (
	"database/sql"
	"forum/internal/models"

	"github.com/gofrs/uuid"
)

type PostRepo struct {
	db *sql.DB
}

func NewPostRepo(db *sql.DB) *PostRepo {
	return &PostRepo{
		db: db,
	}
}

type PostRepoI interface {
	CreatePost(post models.Post) error
	CreatePostCategories(input models.CreateCategoryPost) error
	UpdatePost(post models.Post) error
	GetPost(postID uuid.UUID) (models.Post, error)
	GetPostsAll() ([]models.Post, error)
	GetMyPosts(userID uuid.UUID) ([]models.Post, error)

	GetCategoriesByPostID(postID uuid.UUID) ([]string, error)
	GetAllCategories() ([]*models.Category, error)

	GetVoteOfPost(postID uuid.UUID, userID uuid.UUID) (models.Vote, error)
	DeleteVoteOfPost(voteID uuid.UUID) error
	CreateVote(vote models.Vote) error
	GetVotes() ([]models.Vote, error)

	CreateComment(comment models.Comment) error
	GetCommentsByPostID(postID uuid.UUID) ([]*models.Comment, error)
	GetVoteOfComment(commentID uuid.UUID, userID uuid.UUID) (models.Vote, error)
	GetComment(commentID uuid.UUID) (models.Comment, error)
	UpdateComment(comment models.Comment) error
}

func (ar *PostRepo) CreatePostCategories(input models.CreateCategoryPost) error {
	stmt := `
		INSERT INTO categories_posts_association (association_id, category_id, post_id) 
		VALUES (?, (select category_id from categories where name = ?), ?);
	`
	if _, err := ar.db.Exec(stmt, input.ID, input.CategoryName, input.PostID); err != nil {
		return err
	}
	return nil
}

func (ar *PostRepo) CreatePost(post models.Post) error {
	stmt := `
		INSERT INTO posts (post_id, user_id, title, body, likes, dislikes, image) 
		VALUES (?, ?, ?, ?, ?, ?, ?);
	`
	if _, err := ar.db.Exec(stmt, post.ID, post.UserId, post.Title, post.Body, post.Likes, post.Dislikes, post.Image); err != nil {
		return err
	}
	return nil
}

func (ar *PostRepo) UpdatePost(post models.Post) error {
	stmt := `
		UPDATE posts 
		SET updated_at = CURRENT_TIMESTAMP, title = ?, body = ?, image = ?, likes = ?, dislikes = ? 
		WHERE post_id = ?;
	`
	if _, err := ar.db.Exec(stmt, post.Title, post.Body, post.Image, post.Likes, post.Dislikes, post.ID); err != nil {
		return err
	}
	return nil
}

func (ar *PostRepo) GetPost(postID uuid.UUID) (models.Post, error) {
	var post models.Post
	stmt := `
		SELECT
			post_id,
			created_at,
			updated_at,
			user_id,
			title,
			body,
			likes,
			dislikes,
			image
		FROM posts
		WHERE post_id = ?;
	`
	raw := ar.db.QueryRow(stmt, postID)

	if err := raw.Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt, &post.UserId, &post.Title, &post.Body, &post.Likes, &post.Dislikes, &post.Image); err != nil {
		return post, err
	}
	return post, nil
}

func (ar *PostRepo) GetPostsAll() ([]models.Post, error) {
	var posts []models.Post
	stmt := `
		SELECT
			post_id,
			created_at,
			updated_at,
			user_id,
			title,
			body,
			likes,
			dislikes,
			image
		FROM posts;
	`

	rows, err := ar.db.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt, &post.UserId, &post.Title, &post.Body, &post.Likes, &post.Dislikes, &post.Image); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (ar *PostRepo) GetCategoriesByPostID(postID uuid.UUID) ([]string, error) {
	var categories []string
	stmt := `
		SELECT c.name
		FROM categories c
		JOIN categories_posts_association a ON c.category_id = a.category_id
		JOIN posts p ON a.post_id = p.post_id
		WHERE p.post_id = ?;
	`

	rows, err := ar.db.Query(stmt, postID)
	if err != nil {
		return categories, err
	}
	defer rows.Close()

	for rows.Next() {
		var categoryName string
		if err := rows.Scan(&categoryName); err != nil {
			return categories, err
		}
		categories = append(categories, categoryName)
	}

	if err := rows.Err(); err != nil {
		return categories, err
	}

	return categories, nil
}

func (ar *PostRepo) GetAllCategories() ([]*models.Category, error) {
	var categories []*models.Category
	stmt := `
		SELECT name
		FROM categories 
	`

	rows, err := ar.db.Query(stmt)
	if err != nil {
		return categories, err
	}
	defer rows.Close()

	for rows.Next() {
		var categoryName string
		if err := rows.Scan(&categoryName); err != nil {
			return categories, err
		}
		category := &models.Category{Name: categoryName}
		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		return categories, err
	}

	return categories, nil
}

func (ar *PostRepo) GetMyPosts(userID uuid.UUID) ([]models.Post, error) {
	var posts []models.Post
	stmt := `
		SELECT
			post_id,
			created_at,
			updated_at,
			user_id,
			title,
			body,
			likes,
			dislikes,
			image
		FROM posts
		WHERE user_id = ?;
	`

	rows, err := ar.db.Query(stmt, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt, &post.UserId, &post.Title, &post.Body, &post.Likes, &post.Dislikes, &post.Image); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (ar *PostRepo) GetVoteOfPost(postID uuid.UUID, userID uuid.UUID) (models.Vote, error) {
	var vote models.Vote
	stmt := `
		SELECT
			vote_id,
			user_id,
			post_id,
			binary
		FROM votes
		WHERE post_id = ?
		  AND user_id = ?;
	`
	raw := ar.db.QueryRow(stmt, postID, userID)

	if err := raw.Scan(&vote.ID, &vote.UserID, &vote.PostID, &vote.Binary); err != nil {
		return vote, err
	}
	return vote, nil
}

func (ar *PostRepo) DeleteVoteOfPost(voteID uuid.UUID) error {
	stmt := `
		DELETE FROM votes
		WHERE vote_id = ?;
	`

	if _, err := ar.db.Exec(stmt, voteID); err != nil {
		return err
	}
	return nil
}

func (ar *PostRepo) CreateVote(vote models.Vote) error {
	stmt := `
		INSERT INTO votes (vote_id, user_id, post_id, comment_id, binary) 
		VALUES (?, ?, ?, ?, ?);
	`
	if _, err := ar.db.Exec(stmt, vote.ID, vote.UserID, vote.PostID, vote.CommentID, vote.Binary); err != nil {
		return err
	}
	return nil
}

func (ar *PostRepo) GetVotes() ([]models.Vote, error) {
	var votes []models.Vote
	stmt := `
		SELECT
			vote_id,
			user_id,
			post_id,
			binary
		FROM votes;
	`

	rows, err := ar.db.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var vote models.Vote
		if err := rows.Scan(&vote.ID, &vote.UserID, &vote.PostID, &vote.Binary); err != nil {
			return nil, err
		}
		votes = append(votes, vote)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return votes, nil
}

func (ar *PostRepo) CreateComment(comment models.Comment) error {
	stmt := `
		INSERT INTO comments (comment_id, content, post_id, user_id, likes, dislikes) 
		VALUES (?, ?, ?, ?, ?, ?);
	`
	if _, err := ar.db.Exec(stmt, comment.ID, comment.Description, comment.PostID, comment.UserID, comment.Likes, comment.Dislikes); err != nil {
		return err
	}
	return nil
}

func (ar *PostRepo) GetCommentsByPostID(postID uuid.UUID) ([]*models.Comment, error) {
	var comments []*models.Comment
	stmt := `
		SELECT
			c.comment_id,
			c.created_at,
			c.updated_at,
			c.content,
			c.post_id,
			c.user_id,
			c.likes,
			c.dislikes,
			u.username
		FROM comments c
		JOIN users u ON u.user_id = c.user_id 
		WHERE c.post_id = ?
		ORDER BY c.created_at asc;
	`

	rows, err := ar.db.Query(stmt, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var comment models.Comment
		if err := rows.Scan(&comment.ID, &comment.CreatedAt, &comment.UpdatedAt, &comment.Description, &comment.PostID, &comment.UserID, &comment.Likes, &comment.Dislikes, &comment.UserName); err != nil {
			return nil, err
		}
		comments = append(comments, &comment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}

func (ar *PostRepo) GetVoteOfComment(commentID uuid.UUID, userID uuid.UUID) (models.Vote, error) {
	var vote models.Vote
	stmt := `
		SELECT
			vote_id,
			user_id,
			comment_id,
			binary
		FROM votes
		WHERE comment_id = ?
		  AND user_id = ?;
	`
	raw := ar.db.QueryRow(stmt, commentID, userID)

	if err := raw.Scan(&vote.ID, &vote.UserID, &vote.CommentID, &vote.Binary); err != nil {
		return vote, err
	}
	return vote, nil
}

func (ar *PostRepo) GetComment(commentID uuid.UUID) (models.Comment, error) {
	var comment models.Comment
	stmt := `
		SELECT
			comment_id,
			content,	
			user_id,
			post_id,
			likes,
			dislikes,
			created_at,
			updated_at
		FROM comments
		WHERE comment_id = ?;
	`
	raw := ar.db.QueryRow(stmt, commentID)

	if err := raw.Scan(&comment.ID, &comment.Description, &comment.UserID, &comment.PostID, &comment.Likes, &comment.Dislikes, &comment.CreatedAt, &comment.UpdatedAt); err != nil {
		return comment, err
	}
	return comment, nil
}

func (ar *PostRepo) UpdateComment(comment models.Comment) error {
	stmt := `
		UPDATE comments 
		SET updated_at = CURRENT_TIMESTAMP, likes = ?, dislikes = ? 
		WHERE comment_id = ?;
	`
	if _, err := ar.db.Exec(stmt, comment.Likes, comment.Dislikes, comment.ID); err != nil {
		return err
	}
	return nil
}
