package handler

import (
	"context"
	"fmt"
	"forum/internal/exceptions"
	"forum/pkg/cust_encoders"
	"net/http"
)

func (h *Handler) SessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		session, sessionErr := h.AuthHandler.AuthService.GetSession()
		cookie, cookieErr := r.Cookie("session")
		fmt.Println(r.URL.Path)
		if _, exclude := h.ExcludeSessionHandlersPath[r.URL.Path]; !exclude {
			if _, ok := h.ValidRoutes[r.URL.Path]; !ok {
				dataErr := exceptions.NewResourceNotFoundError("Page is not found")
				params := cust_encoders.EncodeParams(dataErr)
				http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
				return
			}
			fmt.Println("Not excluded url paths")
			if sessionErr != nil {
				http.Redirect(w, r, "/signin", http.StatusSeeOther)
				return
			}

			if cookieErr != nil {
				http.Redirect(w, r, "/signin", http.StatusSeeOther)
				return
			}

			if session.Token != cookie.Value {
				expiredCookie := &http.Cookie{
					Name:     "session",
					Value:    "",
					Path:     "/",
					HttpOnly: true,
					MaxAge:   -1,
				}

				h.AuthHandler.AuthService.DeleteSession()
				http.SetCookie(w, expiredCookie)

				dataErr := exceptions.NewAuthenticationError("Session is expired or invalid")
				params := cust_encoders.EncodeParams(dataErr)
				http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
				return
			}

			ctx := context.WithValue(r.Context(), "session", session)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			fmt.Println("SESSION IS EXCLUDED URL PATHS")
			if sessionErr == nil { // если сессия есть
				fmt.Println("SESSION IS PRESENT")
				if r.URL.Path == "/signin" {
					fmt.Println("URL PATH IS SIGNIN")
					if cookieErr != nil {
						fmt.Println("COOKIE IS NOT PRESENT")
						if r.Method == http.MethodPost {
							fmt.Println("METHOD IS POST")
							h.AuthHandler.AuthService.DeleteSession()
						}
						next.ServeHTTP(w, r)
						return
					} else {
						fmt.Println("COOKIE IS PRESENT")
						if session.Token == cookie.Value {
							fmt.Println("SESSION VALUE IS EQUAL TO COOKIE")
							http.Redirect(w, r, "/post/", http.StatusSeeOther)
							return
							//ctx := context.WithValue(r.Context(), "session", session)
							//next.ServeHTTP(w, r.WithContext(ctx))
							//return
						} else {
							fmt.Println("SESSION VALUE IS NOT EQUAL TO COOKIE")
							expiredCookie := &http.Cookie{
								Name:     "session",
								Value:    "",
								Path:     "/",
								HttpOnly: true,
								MaxAge:   -1,
							}
							http.SetCookie(w, expiredCookie)
							dataErr := exceptions.NewAuthenticationError("User session or cookie is expired")
							params := cust_encoders.EncodeParams(dataErr)
							http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
							return
						}
					}
				} else {
					fmt.Println("URL IS EITHER POST OR POST/GET")
					if cookieErr != nil {
						//http.Redirect(w, r, "/signin", http.StatusSeeOther)
						//return
						next.ServeHTTP(w, r)
						return
					} else {
						fmt.Println("")
						if session.Token == cookie.Value {
							ctx = context.WithValue(r.Context(), "session", session)
							next.ServeHTTP(w, r.WithContext(ctx))
							return
						} else {
							expiredCookie := &http.Cookie{
								Name:     "session",
								Value:    "",
								Path:     "/",
								HttpOnly: true,
								MaxAge:   -1,
							}
							http.SetCookie(w, expiredCookie)

							dataErr := exceptions.NewAuthenticationError("User cookie is expired or incorrect")
							params := cust_encoders.EncodeParams(dataErr)
							http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
							return
						}
					}
				}
			} else {
				expiredCookie := &http.Cookie{
					Name:     "session",
					Value:    "",
					Path:     "/",
					HttpOnly: true,
					MaxAge:   -1,
				}
				http.SetCookie(w, expiredCookie)
				next.ServeHTTP(w, r)
			}
		}
	})
}

func (h *Handler) ErrorMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			stringParams := r.URL.RawQuery
			if stringParams == "" {
				http.Redirect(w, r, "/signin", http.StatusSeeOther)
			} else {
				dataErr, err := cust_encoders.DecodeParams(stringParams)
				if err != nil {
					dataErr := exceptions.NewBadRequestError("Invalid error parameters are given")
					params := cust_encoders.EncodeParams(dataErr)
					http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
				}
				var (
					customErr  error
					ctx        context.Context
					httpStatus int
				)
				switch dataErr.(type) {
				case exceptions.AuthenticationError:
					customErr = dataErr
					httpStatus = dataErr.(exceptions.AuthenticationError).StatusCode
				case exceptions.ForbiddenError:
					customErr = dataErr
					httpStatus = dataErr.(exceptions.ForbiddenError).StatusCode
				case exceptions.ResourceNotFoundError:
					customErr = dataErr
					httpStatus = dataErr.(exceptions.ResourceNotFoundError).StatusCode
				case exceptions.StatusMethodNotAllowed:
					customErr = dataErr
					httpStatus = dataErr.(exceptions.StatusMethodNotAllowed).StatusCode
				case exceptions.StatusConflictError:
					customErr = dataErr
					httpStatus = dataErr.(exceptions.StatusConflictError).StatusCode
				case exceptions.ValidationError:
					customErr = dataErr
					httpStatus = dataErr.(exceptions.ValidationError).StatusCode
				case exceptions.InternalServerError:
					customErr = dataErr
					httpStatus = dataErr.(exceptions.InternalServerError).StatusCode
				case exceptions.BadRequestError:
					customErr = dataErr
					httpStatus = dataErr.(exceptions.BadRequestError).StatusCode
				default:
					http.Redirect(w, r, "/signin", http.StatusSeeOther)
					return
				}

				ctx = context.WithValue(r.Context(), "error", customErr)
				ctx = context.WithValue(ctx, "httpStatus", httpStatus)
				next.ServeHTTP(w, r.WithContext(ctx))
			}
		} else {
			dataErr := exceptions.NewResourceNotFoundError("Page is not found")
			params := cust_encoders.EncodeParams(dataErr)
			http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
		}
	})
}
