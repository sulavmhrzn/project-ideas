package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) router() http.Handler {
	mux := httprouter.New()
	mux.HandlerFunc(http.MethodGet, "/v1/ping", app.pingHandler)
	mux.HandlerFunc(http.MethodPost, "/v1/users", app.createUserHandler)
	return app.logRequestMiddleware(mux)
}
