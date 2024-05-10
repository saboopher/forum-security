package authhandler

import (
	"fmt"
	"forum/internal/exceptions"
	"forum/internal/schemas"
	"forum/pkg/cust_encoders"
	"forum/pkg/validator"
	"html/template"
	"net/http"
)

func (ah *AuthHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost || r.Method == http.MethodGet {
		registerForm := &schemas.RegisterForm{}
		if r.Method == http.MethodPost {
			if err := r.ParseForm(); err != nil {
				dataErr := exceptions.NewBadRequestError("Invalid form values")
				params := cust_encoders.EncodeParams(dataErr)
				http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
			} else {
				userName := r.FormValue("username")
				email := r.FormValue("email")
				password := r.FormValue("password")
				confirmPassword := r.FormValue("password_confirm")

				nameOk, msgName := validator.ValidateName(userName)

				email, emailOk, msgEmail := validator.ValidateEmail(email)

				passwordOk, msgPassword := validator.ValidatePassword(password)

				confirmPasswordOk, msgConfirmedPassword := validator.ValidatePasswordConfirmed(password, confirmPassword)
				if !nameOk || !emailOk || !passwordOk || !confirmPasswordOk {
					registerForm.TemplateForm = &schemas.TemplateForm{}
					if !nameOk {
						registerForm.TemplateForm.RegisterErrors.Name = msgName
					}

					if !emailOk {
						registerForm.TemplateForm.RegisterErrors.Email = msgEmail
					}
					if !passwordOk {
						registerForm.TemplateForm.RegisterErrors.Password = msgPassword
					}

					if !confirmPasswordOk {
						registerForm.TemplateForm.RegisterErrors.ConfirmPassword = msgConfirmedPassword
					}
					registerForm.TemplateForm.RegisterDataForErr.Name = userName
					registerForm.TemplateForm.RegisterDataForErr.Email = email
					registerForm.TemplateForm.RegisterDataForErr.Password = password
					registerForm.TemplateForm.RegisterDataForErr.ConfirmPassword = confirmPassword
				}

				if nameOk && emailOk && passwordOk && confirmPasswordOk {
					user := schemas.CreateUser{
						UpdateUser: schemas.UpdateUser{
							Username: userName,
							Email:    email,
							Password: password,
						},
						PasswordConfirm: r.Form.Get("password_confirm"),
					}
					err := ah.AuthService.CreateUser(user)
					if err != nil {
						params := cust_encoders.EncodeParams(err)
						http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
					} else {
						fmt.Println("asd")
						http.Redirect(w, r, "/signin", http.StatusSeeOther)
					}
				}
			}
		}

		t, err := template.ParseFiles("ui/templates/register.html")
		if err != nil {
			dataErr := exceptions.NewInternalServerError()
			params := cust_encoders.EncodeParams(dataErr)
			http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
		} else {
			err := t.Execute(w, registerForm)
			if err != nil {
				dataErr := exceptions.NewInternalServerError()
				params := cust_encoders.EncodeParams(dataErr)
				http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
				return
			}
		}
	} else {
		dataErr := exceptions.NewStatusMethodNotAllowed()
		params := cust_encoders.EncodeParams(dataErr)
		http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
		return
	}
}

func (ah *AuthHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		t, err := template.ParseFiles("ui/templates/signin.html")
		if err != nil {
			dataErr := exceptions.NewInternalServerError()
			params := cust_encoders.EncodeParams(dataErr)
			http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
		} else {
			err := t.Execute(w, nil)
			if err != nil {
				dataErr := exceptions.NewInternalServerError()
				params := cust_encoders.EncodeParams(dataErr)
				http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
			}
		}

	} else if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			dataErr := exceptions.NewBadRequestError("Invalid form values")
			params := cust_encoders.EncodeParams(dataErr)
			http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
		} else {
			email := r.Form.Get("email")
			email = validator.ClearEmail(email)

			user := schemas.AuthUser{
				Email:    email,
				Password: r.Form.Get("password"),
			}
			session, err := ah.AuthService.CreateSession(user)
			if err != nil {
				params := cust_encoders.EncodeParams(err)
				http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
			} else {
				cookie := &http.Cookie{
					Name:     "session",
					Value:    session.Token,
					Path:     "/",
					Expires:  session.ExpireTime,
					HttpOnly: true,
					MaxAge:   7200,
				}

				http.SetCookie(w, cookie)
				http.Redirect(w, r, "/post/", http.StatusSeeOther)
			}
		}
	} else {
		dataErr := exceptions.NewStatusMethodNotAllowed()
		params := cust_encoders.EncodeParams(dataErr)
		http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
	}
}

func (ah *AuthHandler) LogOut(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		dataErr := exceptions.NewStatusMethodNotAllowed()
		params := cust_encoders.EncodeParams(dataErr)
		http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
		return
	}

	ah.AuthService.DeleteSession()
	expiredCookie := &http.Cookie{
		Name:     "session",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	}
	http.SetCookie(w, expiredCookie)

	http.Redirect(w, r, "/post/", http.StatusSeeOther)
}
