package main

import (
	"net/http"
	"time"

	"github.com/sulavmhrzn/projectideas/internal/data"
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
		Username:  input.Username,
		Email:     input.Email,
		CreatedAt: time.Now(),
	}
	err = user.Password.Set(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJSON(w, http.StatusOK, map[string]any{"user": user})
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
