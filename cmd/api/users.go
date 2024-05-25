package main

import (
	"net/http"
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
	err = app.writeJSON(w, http.StatusOK, map[string]any{"input": input})
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
