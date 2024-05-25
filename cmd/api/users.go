package main

import (
	"errors"
	"net/http"
	"time"

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

func (app *application) loginUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	v := validator.New()
	v.Check(input.Email != "", "email", "must be provided")
	v.Check(validator.ValidEmail(input.Email), "email", "must be a valid email address")
	v.Check(input.Password != "", "password", "must be provided")
	v.Check(len(input.Password) < 72, "password", "must not be greater than 72 characters long")
	v.Check(len(input.Password) > 10, "password", "must be greater than 10 characters long")
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	user, err := app.models.User.GetByEmail(input.Email)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrNoRows):
			app.invalidCredentialsResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	ok := user.Password.Compare(input.Password)
	if !ok {
		app.invalidCredentialsResponse(w, r)
		return
	}
	token, err := app.models.Token.New(user.Id, 1*time.Hour, data.ScopeAuthentication)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, token)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
