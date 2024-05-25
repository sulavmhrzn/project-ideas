package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

const (
	reset = "\033[0m"
	red   = "\033[31m"
	green = "\033[32m"
)

type config struct {
	port int
}
type application struct {
	cfg      config
	infoLog  *log.Logger
	errorLog *log.Logger
}

func main() {
	cfg := config{}
	flag.IntVar(&cfg.port, "port", 4000, "port to listen")
	flag.Parse()

	app := &application{
		cfg:      cfg,
		infoLog:  log.New(os.Stdout, fmt.Sprintf("%sINFO: %s", green, reset), log.Ldate|log.Ltime|log.Lshortfile),
		errorLog: log.New(os.Stderr, fmt.Sprintf("%sERROR: %s", red, reset), log.Ldate|log.Ltime|log.Lshortfile),
	}
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", app.cfg.port),
		Handler: app.router(),
	}
	app.infoLog.Printf("server running on port %s", server.Addr)
	err := server.ListenAndServe()
	if err != nil {
		app.errorLog.Fatal(err)
	}
}
