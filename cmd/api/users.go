package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/sulavmhrzn/projectideas/internal/data"
	"github.com/sulavmhrzn/projectideas/internal/mailer"
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

func (app *application) generateTokenHandler(w http.ResponseWriter, r *http.Request) {
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

func (app *application) sendResetPasswordTokenHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email string `json:"email"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	user, err := app.models.User.GetByEmail(input.Email)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrNoRows):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	token, err := app.models.Token.New(user.Id, 24*time.Hour, data.ScopePasswordReset)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	go func() {
		dialer := mailer.NewDialer(app.cfg.mailer.host, app.cfg.mailer.port, app.cfg.mailer.username, app.cfg.mailer.password)
		err = mailer.SendMail(dialer,
			app.cfg.mailer.EmailFrom,
			user.Email,
			"Reset Password",
			fmt.Sprintf("Your reset token %q. Please make an request to /v1/users/resetPassword with this token", token.Token),
		)
		if err != nil {
			app.logError(err)
		}
	}()
	err = app.writeJSON(w, http.StatusOK, map[string]any{"message": "reset password token sent"})
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) resetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Token    string `json:"token"`
		Password string `json:"password"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	v := validator.New()
	v.Check(input.Token != "", "token", "must be provided")
	v.Check(input.Password != "", "password", "must be provided")
	v.Check(len(input.Password) > 10, "password", "must be greater than 10 characters long")
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	user, err := app.models.User.GetForToken(input.Token, data.ScopePasswordReset)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrNoRows):
			app.invalidTokenResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = user.Password.Set(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.models.User.UpdatePassword(user.Id, user.Password.HashedPassword)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.models.Token.DeleteForUser(user.Id)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJSON(w, http.StatusOK, map[string]string{"message": "password reset successfull"})
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
