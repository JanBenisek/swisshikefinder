package main //belongs to the main package

import (
	// embed static files in the binary
	"database/sql"
	"embed"

	"log"
	"net/http" // webserver
	"os"

	// internal stuff
	"internal/models"

	_ "github.com/marcboeker/go-duckdb"
)

// package level variables - means that it is available anywhere in this package
//
//go:embed all:static
var static embed.FS

func openDB() (*sql.DB, error) {
	db, err := sql.Open("duckdb", "./duck.db?autoinstall_known_extensions=1&autoload_known_extensions=1")
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}

	_, err = db.Exec(`IMPORT DATABASE './db/'`)
	if err != nil {
		log.Fatal(err)
	}

	return db, nil
}

type application struct {
	DebugLog *log.Logger
	InfoLog  *log.Logger
	ErrorLog *log.Logger
	Port     string
	Tours    *models.TourModels
}

func main() {

	DebugLog := log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime)
	InfoLog := log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)
	ErrorLog := log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime)

	// Get a database handle.
	db, err := openDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Get the application struct and set some env values
	app := &application{
		InfoLog:  InfoLog,
		DebugLog: DebugLog,
		ErrorLog: ErrorLog,
		Tours:    &models.TourModels{DB: db},
	}

	// Digital Ocean always listens on 8080 and has the env var set
	app.Port = os.Getenv("PORT") // will be available at http://localhost:8080
	if app.Port == "" {
		app.InfoLog.Printf("Port not found in .env, using default 8080")
		app.Port = "8080"
	}

	// my version of server, I can pass my own logger
	srv := &http.Server{
		Addr:     ":" + app.Port,
		ErrorLog: app.ErrorLog,
		Handler:  app.routes(), // giving it my routes
	}

	app.InfoLog.Printf("Starting server on %s", srv.Addr)
	srv.ListenAndServe() //start the service and listen to the port with the mux

}
