package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) router() *httprouter.Router {
	mux := httprouter.New()
	mux.HandlerFunc(http.MethodGet, "/v1/ping", app.pingHandler)
	return mux
}
