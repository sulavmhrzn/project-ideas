package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) router() http.Handler {
	mux := httprouter.New()
	mux.HandlerFunc(http.MethodGet, "/v1/ping", app.requireLoginMiddleware(app.pingHandler))
	mux.HandlerFunc(http.MethodPost, "/v1/users/register", app.createUserHandler)
	mux.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.generateTokenHandler)
	return app.logRequestMiddleware(mux)
}
