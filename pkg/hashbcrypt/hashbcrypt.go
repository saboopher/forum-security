package hashbcrypt

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

var ErrHashPass = errors.New("could not generate password hash")

func GenerateHashedPassword(password string) (string, error) {
	// add salt(uuid) to the password
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckHashedPassword(passwordInput, passwordDb string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(passwordDb), []byte(passwordInput))
	return err == nil
}
