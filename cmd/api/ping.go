package main

import (
	"log"
	"net/http"
)

func (app *application) pingHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]any{
		"ping": "pong",
	}
	err := app.writeJSON(w, http.StatusOK, data)
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
