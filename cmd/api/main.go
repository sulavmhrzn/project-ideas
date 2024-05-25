package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

const (
	reset = "\033[0m"
	red   = "\033[31m"
	green = "\033[32m"
)

type config struct {
	port int
	dsn  string
}
type application struct {
	cfg      config
	infoLog  *log.Logger
	errorLog *log.Logger
}

func main() {
	cfg := config{}
	flag.IntVar(&cfg.port, "port", 4000, "port to listen")
	flag.StringVar(&cfg.dsn, "dsn", "postgres://projectideasuser:sulavpostgres@localhost:5432/projectideas", "Database dsn")
	flag.Parse()

	app := &application{
		cfg:      cfg,
		infoLog:  log.New(os.Stdout, fmt.Sprintf("%sINFO: %s", green, reset), log.Ldate|log.Ltime|log.Lshortfile),
		errorLog: log.New(os.Stderr, fmt.Sprintf("%sERROR: %s", red, reset), log.Ldate|log.Ltime|log.Lshortfile),
	}

	_, err := openDB(cfg)
	if err != nil {
		log.Fatal(err)
	}
	app.infoLog.Println("database connection successful")
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", app.cfg.port),
		Handler: app.router(),
	}
	app.infoLog.Printf("server running on port %s", server.Addr)
	err = server.ListenAndServe()
	if err != nil {
		app.errorLog.Fatal(err)
	}
}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
