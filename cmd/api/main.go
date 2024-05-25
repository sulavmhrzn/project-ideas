package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

type config struct {
	port int
}
type application struct {
	cfg config
}

func main() {
	cfg := config{}
	flag.IntVar(&cfg.port, "port", 4000, "port to listen")
	flag.Parse()

	app := &application{
		cfg: cfg,
	}
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", app.cfg.port),
		Handler: app.router(),
	}
	log.Printf("server running on port %s", server.Addr)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
