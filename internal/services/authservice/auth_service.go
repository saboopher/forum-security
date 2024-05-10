package authservice

import (
	"errors"
	"forum/internal/exceptions"
	"forum/internal/models"
	"forum/internal/repositories/authrepo"
	"forum/internal/schemas"
	"forum/pkg/hashbcrypt"
	"forum/pkg/validator"
	"time"

	"github.com/gofrs/uuid"
)

type AuthService struct {
	AuthRepo authrepo.AuthRepoI
}

func NewAuthService(authRepo authrepo.AuthRepoI) *AuthService {
	return &AuthService{
		AuthRepo: authRepo,
	}
}

type AuthServiceI interface {
	CreateUser(user schemas.CreateUser) error
	CreateSession(user schemas.AuthUser) (models.Session, error)

	CheckUserPassword(email, password string) error
	GetSession() (models.Session, error)

	DeleteSession() error
}

func (as *AuthService) CreateUser(user schemas.CreateUser) error {
	hashedPassword, err := hashbcrypt.GenerateHashedPassword(user.Password)
	if err != nil {
		return exceptions.NewInternalServerError()
	}
	userFromDb, _ := as.AuthRepo.GetUserByEmail(user.Email)
	if userFromDb.Email == user.Email {
		return exceptions.NewStatusConflicError("User with this email already exists")
	}
	userFromDB2, _ := as.AuthRepo.GetUserByUsername(user.Username)
	if userFromDB2.Username == user.Username {
		return exceptions.NewStatusConflicError("User with this username already exists")
	}
	userModel := models.User{
		ID:             uuid.Must(uuid.NewV4()),
		Username:       user.Username,
		Email:          user.Email,
		HashedPassword: hashedPassword,
	}
	if err := as.AuthRepo.CreateUser(userModel); err != nil {
		return exceptions.NewInternalServerError()
	}

	return nil
}

func (as *AuthService) CreateSession(user schemas.AuthUser) (models.Session, error) {
	session := models.Session{}
	if err := validator.ValidateSignInInput(user); err != nil {
		return session, exceptions.NewValidationError("Email invalid")
	}

	userDB, err := as.AuthRepo.GetUserByEmail(user.Email)
	if err != nil {
		return session, exceptions.NewStatusConflicError("User not found")
	}

	err = as.CheckUserPassword(userDB.HashedPassword, user.Password)
	if err != nil {
		return session, exceptions.NewAuthenticationError("Password is incorrect")
	}

	as.AuthRepo.DeleteSession()

	session = models.Session{
		ID:         uuid.Must(uuid.NewV4()),
		UserID:     userDB.ID,
		Token:      uuid.Must(uuid.NewV4()).String(),
		ExpireTime: time.Now().Add(time.Hour * 2),
	}

	err = as.AuthRepo.CreateSession(session)
	if err != nil {
		return models.Session{}, exceptions.NewInternalServerError()
	}

	return session, nil
}

func (as *AuthService) CheckUserPassword(passwordDB, password string) error {
	if !hashbcrypt.CheckHashedPassword(password, passwordDB) {
		return errors.New("wrong password") // где хранить все ошибки?
	}

	return nil
}

func (as *AuthService) DeleteSession() error {
	err := as.AuthRepo.DeleteSession()
	if err != nil {
		return exceptions.NewInternalServerError()
	}
	return nil
}

func (as *AuthService) GetSession() (models.Session, error) {
	session, err := as.AuthRepo.GetSession()
	if err != nil {
		return models.Session{}, exceptions.NewInternalServerError()
	}
	return session, nil
}
