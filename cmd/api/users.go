package main

import (
	"errors"
	"net/http"

	"github.com/sulavmhrzn/projectideas/internal/data"
	"github.com/sulavmhrzn/projectideas/internal/validator"
)

func (app *application) createUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := &data.User{
		Username: input.Username,
		Email:    input.Email,
	}
	user.Password.PlainPassword = input.Password

	v := validator.New()
	if data.ValidateUser(v, user); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = user.Password.Set(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	user, err = app.models.User.Insert(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateEmail) || errors.Is(err, data.ErrDuplicateUsername):
			app.badRequestResponse(w, r, err)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, map[string]any{"user": user})
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
