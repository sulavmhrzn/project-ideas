package main

import (
	"net/http"
	"runtime"
)

func (app *application) logError(err error) {
	_, file, line, _ := runtime.Caller(2)
	app.errorLog.Printf("%s:%d %v", file, line, err)
}

func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, message any) {
	data := map[string]any{
		"error": message,
	}
	err := app.writeJSON(w, http.StatusInternalServerError, data)
	if err != nil {
		app.logError(err)
		http.Error(w, "", http.StatusInternalServerError)
	}

}

func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	message := "internal server error"
	app.logError(err)
	app.errorResponse(w, r, message)
}
