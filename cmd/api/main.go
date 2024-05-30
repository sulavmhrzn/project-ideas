package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sulavmhrzn/projectideas/internal/data"
)

const (
	reset = "\033[0m"
	red   = "\033[31m"
	green = "\033[32m"
)

type config struct {
	port   int
	dsn    string
	mailer struct {
		host      string
		port      int
		username  string
		password  string
		EmailFrom string
	}
}
type application struct {
	cfg      config
	infoLog  *log.Logger
	errorLog *log.Logger
	models   data.Model
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	cfg := config{}
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		log.Fatal(err)
	}
	mailerPort, err := strconv.Atoi(os.Getenv("MAILER_PORT"))
	if err != nil {
		log.Fatal(err)
	}

	flag.IntVar(&cfg.port, "port", port, "port to listen")
	flag.StringVar(&cfg.dsn, "dsn", os.Getenv("DSN"), "Database dsn")
	flag.StringVar(&cfg.mailer.host, "mailer-host", os.Getenv("MAILER_HOST"), "mailer host")
	flag.IntVar(&cfg.mailer.port, "mailer-port", mailerPort, "mailer port")
	flag.StringVar(&cfg.mailer.username, "mailer-username", os.Getenv("MAILER_USERNAME"), "mailer username")
	flag.StringVar(&cfg.mailer.password, "mailer-password", os.Getenv("MAILER_PASSWORD"), "mailer password")
	flag.StringVar(&cfg.mailer.EmailFrom, "mailer-email-from", os.Getenv("MAILER_EMAIL_FROM"), "mailer email from")
	flag.Parse()

	db, err := openDB(cfg)
	if err != nil {
		log.Fatal(err)
	}
	app := &application{
		cfg:      cfg,
		infoLog:  log.New(os.Stdout, fmt.Sprintf("%sINFO: %s", green, reset), log.Ldate|log.Ltime|log.Lshortfile),
		errorLog: log.New(os.Stderr, fmt.Sprintf("%sERROR: %s", red, reset), log.Ldate|log.Ltime|log.Lshortfile),
		models:   data.NewModel(db),
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
