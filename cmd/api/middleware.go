package main

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/sulavmhrzn/projectideas/internal/data"
)

func (app *application) logRequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s %s %s", r.Method, r.RemoteAddr, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func (app *application) requireLoginMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorizationHeader := r.Header.Get("Authorization")
		if authorizationHeader == "" {
			app.invalidTokenResponse(w, r)
			return
		}
		format := strings.Split(authorizationHeader, " ")
		if len(format) != 2 || format[0] != "Bearer" || format[1] == "" {
			app.invalidTokenResponse(w, r)
			return
		}
		token := format[1]
		user, err := app.models.User.GetForToken(token)
		if err != nil {
			switch {
			case errors.Is(err, data.ErrNoRows):
				app.invalidTokenResponse(w, r)
			default:
				app.serverErrorResponse(w, r, err)
			}
			return
		}
		// TODO: add value to the request context
		fmt.Println("logged in as ", user.Email)
		next.ServeHTTP(w, r)
	})
}
