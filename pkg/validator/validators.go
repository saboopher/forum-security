package validator

import (
	"errors"
	"fmt"
	"forum/internal/schemas"
	"regexp"
	"strings"
	"unicode"

	"github.com/gofrs/uuid"
)

var (
	ErrorInvalidPassword      = errors.New("Invalid password")
	ErrorInvalidEmail         = errors.New("Invalid email adress")
	ErrorPasswordMatch        = errors.New("Password and confirmed password do not match")
	ErrorEmptyField           = errors.New("Empty field")
	ErrorInvalidPostTtitleLen = errors.New("the title should be 5 - 50 characters in length")
	ErrorInvalidPostBodyLen   = errors.New("the body  should be 20 - 2000 characters in length")
	ErrorInvalidUUID          = errors.New("invalid uuid")
)

func ValidateName(userName string) (bool, string) {
	if userName == "" {
		return false, "User name should not be empty"
	}

	if len(userName) < 5 {
		return false, "User name length should be at least 5 characters"
	}

	return true, "success"
}

func ValidateEmail(userEmail string) (string, bool, string) {
	if userEmail == "" {
		return "", false, "User email should not be empty"
	}

	if !validateEmail(userEmail) {
		return "", false, "User email contains incorrect characters"
	}

	email := ClearEmail(userEmail)

	return email, true, "success"
}

func validateEmail(e string) bool {
	emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return emailRegex.MatchString(e)
}

func ClearEmail(email string) string {
	email = strings.TrimSpace(email)
	email = strings.ToLower(email)
	return email
}

func ValidatePassword(password string) (bool, string) {
	if password == "" {
		return false, "User password should not be empty"
	}

	if len(password) < 8 {
		return false, "User password should contain at least 8 characters"
	}

	hasUpper := false
	hasLower := false
	hasDigit := false

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		}
	}

	if !hasUpper || !hasLower || !hasDigit || !validateSpecialChar(password) {
		return false, "User password should contain at least one upper, lower, digit and special characters"
	}

	return true, "success"
}

//func ValidateRegisterInput(user schemas.CreateUser) error {
//	if user.Username == "" || user.Email == "" || user.Password == "" || user.PasswordConfirm == "" {
//		return ErrorEmptyField
//	}
//	if !validateEmail(user.Email) {
//		return ErrorInvalidEmail
//	}
//
//	// if err := validatePassword(user.Password); err != nil {
//	// 	return ErrorInvalidPassword
//	// }
//
//	if !validatePasswordConfirmed(user.Password, user.PasswordConfirm) {
//		return ErrorPasswordMatch
//	}
//
//	return nil
//}

func ValidateSignInInput(user schemas.AuthUser) error {
	if user.Email == "" || user.Password == "" {
		return ErrorEmptyField
	}
	if !validateEmail(user.Email) {
		return ErrorInvalidEmail
	}

	return nil
}

func ValidateCreatePostInput(post schemas.CreatePost) error {
	title := strings.TrimSpace(post.Title)
	body := strings.TrimSpace(post.Body)

	if len(title) > 50 && len(title) < 5 {
		return ErrorInvalidPostTtitleLen
	}
	if len(body) < 20 && len(body) > 2000 {
		return ErrorInvalidPostBodyLen
	}
	return nil
}

func ValidatePostTitle(title string) (bool, string) {
	title = strings.TrimSpace(title)
	fmt.Println(len(title))
	if len(title) > 50 || len(title) < 5 {
		return false, "Post title should be at least 5 and at most 50 characters"
	}
	return true, "success"
}

func ValidatePostBody(body string) (bool, string) {
	body = strings.TrimSpace(body)

	if len(body) < 20 || len(body) > 2000 {
		return false, "Post content should be at least 20 and at most 2000 characters"
	}
	return true, "success"
}

func ValidatePostComment(body string) (bool, string) {
	body = strings.TrimSpace(body)

	if len(body) < 5 || len(body) > 250 {
		return false, "Post comment should be at least 5 and at most 250 characters"
	}
	return true, "success"
}

func ValidateCategoryLen(categoriesInput []string, categories []*schemas.Category) (bool, string) {
	if len(categoriesInput) < 1 {
		return false, "Post should has at least one category"
	}

	if !isUniqueCategory(categoriesInput) {
		return false, "Duplicate post category"
	}

	for _, categoryInput := range categoriesInput {
		var present bool
		for _, category := range categories {
			if categoryInput == category.Name {
				present = true
			}
		}

		if !present {
			return false, "Unknown post category"
		}
	}

	// for _, category := range categories {
	// 	for _, categoryInput := range categoriesInput {
	// 		if category.Name != categoryInput {
	// 			return false, "Unknown post category"
	// 		}
	// 	}
	// }
	return true, "success"
}

func isUniqueCategory(categoriesInput []string) bool {
	encountered := make(map[string]bool)

	for _, v := range categoriesInput {
		if encountered[v] {
			return false
		}
		encountered[v] = true
	}
	return true
}

func ValidateUpdatePostInput(post schemas.UpdatePost) error {
	title := strings.TrimSpace(post.CreatePost.Title)
	body := strings.TrimSpace(post.CreatePost.Body)

	if len(title) > 50 && len(title) < 5 {
		return ErrorInvalidPostTtitleLen
	}
	if len(body) < 20 && len(body) > 2000 {
		return ErrorInvalidPostBodyLen
	}
	return nil
}

func ValidatePasswordConfirmed(password, PasswordConfirm string) (bool, string) {
	if password != PasswordConfirm {
		return false, "Passwords do not match"
	}

	return true, "success"
}

func validateSpecialChar(password string) bool {
	specialCharsSet := "!@$%^&*()_-+"

	for _, char := range password {
		if strings.ContainsRune(specialCharsSet, char) {
			return true
		}
	}

	return false
}

func ValidateUUID(input string) bool {
	_, err := uuid.FromString(input)
	return err == nil
}
