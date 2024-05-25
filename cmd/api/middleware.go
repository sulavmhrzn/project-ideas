package main

import "net/http"

func (app *application) logRequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s %s %s", r.Method, r.RemoteAddr, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}
