package main

import (
	"net/http"
)

func (app *application) pingHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]any{
		"ping": "pong",
	}
	err := app.writeJSON(w, http.StatusOK, data)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
