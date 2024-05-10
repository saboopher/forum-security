package authrepo

import (
	"database/sql"
	"forum/internal/models"

	"github.com/gofrs/uuid"
)

type AuthRepo struct {
	db *sql.DB
}

func NewAuthRepo(db *sql.DB) *AuthRepo {
	return &AuthRepo{
		db: db,
	}
}

type AuthRepoI interface {
	CreateUser(user models.User) error
	CreateSession(session models.Session) error

	GetUserByEmail(email string) (models.User, error)
	GetUserByUsername(username string) (models.User, error)
	GetUserByToken(token string) (models.User, error)
	GetSession() (models.Session, error)
	GetUserByUserID(userID uuid.UUID) (models.User, error)

	DeleteSession() error
}

func (ar *AuthRepo) CreateUser(user models.User) error {
	stmt := `
		INSERT INTO users (user_id, username, email, password) 
		VALUES (?, ?, ?, ?);
	`
	if _, err := ar.db.Exec(stmt, user.ID, user.Username, user.Email, user.HashedPassword); err != nil {
		return err
	}
	return nil
}

func (ar *AuthRepo) CreateSession(session models.Session) error {
	stmt := `
		INSERT INTO sessions (session_id, user_id, token, expire_time) 
		VALUES (?, ?, ?, ?);
	`
	if _, err := ar.db.Exec(stmt, session.ID, session.UserID, session.Token, session.ExpireTime); err != nil {
		return err
	}

	return nil
}

func (ar *AuthRepo) GetUserByEmail(email string) (models.User, error) {
	var user models.User
	stmt := `
		SELECT user_id, username, email, password FROM users 
		WHERE email = ?;
	`
	if err := ar.db.QueryRow(stmt, email).Scan(&user.ID, &user.Username, &user.Email, &user.HashedPassword); err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (ar *AuthRepo) GetUserByUsername(username string) (models.User, error) {
	var user models.User
	stmt := `
		SELECT user_id, username, email, password FROM users
		WHERE username = ?
	`

	if err := ar.db.QueryRow(stmt, username).Scan(&user.ID, &user.Username, &user.Email, &user.HashedPassword); err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (ar *AuthRepo) GetUserByToken(email string) (models.User, error) {
	var user models.User
	stmt := `
		SELECT user_id, username, email, password FROM users 
		WHERE email = ?;
	`
	if err := ar.db.QueryRow(stmt, email).Scan(&user.ID, &user.Username, &user.Email, &user.HashedPassword); err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (ar *AuthRepo) DeleteSession() error {
	stmt := `
		DELETE FROM sessions;
	`
	if _, err := ar.db.Exec(stmt); err != nil {
		return err
	}

	return nil
}

func (ar *AuthRepo) GetSession() (models.Session, error) {
	var session models.Session
	stmt := `
		SELECT session_id, user_id, token, expire_time FROM sessions;
	`
	if err := ar.db.QueryRow(stmt).Scan(&session.ID, &session.UserID, &session.Token, &session.ExpireTime); err != nil {
		return models.Session{}, err
	}

	return session, nil
}

func (ar *AuthRepo) GetUserByUserID(userID uuid.UUID) (models.User, error) {
	var user models.User
	stmt := `
		SELECT user_id, username, email FROM users 
		WHERE user_id = ?;
	`
	if err := ar.db.QueryRow(stmt, userID).Scan(&user.ID, &user.Username, &user.Email); err != nil {
		return models.User{}, err
	}

	return user, nil
}
