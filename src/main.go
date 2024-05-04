package main //belongs to the main package

import (
	// embed static files in the binary
	"database/sql"
	"embed"
	"fmt"
	"log"
	"net/http" // webserver
	"os"

	// access os stuff
	conf "internal/config"

	_ "github.com/marcboeker/go-duckdb"
)

// package level variables - means that it is available anywhere in this package
//
//go:embed all:static
var static embed.FS

type Tour struct {
	ID   string
	Name string
}

func getOneRow(n_rows int64) (Tour, error) {

	// Get a database handle.
	db, err := sql.Open("duckdb", "./duck.db?autoinstall_known_extensions=1&autoload_known_extensions=1")
	if err != nil {
		log.Fatal(err)
	}

	// An album to hold data from the returned row.
	var tr Tour

	row := db.QueryRow("select id, name from './data/tours.parquet' limit ?", n_rows)
	if err := row.Scan(&tr.ID, &tr.Name); err != nil {
		if err == sql.ErrNoRows {
			return tr, fmt.Errorf("id %d: no rows", n_rows)
		}
		return tr, fmt.Errorf("n_rows: %d: %v", n_rows, err)
	}
	return tr, nil
}

func main() {

	// Get the application struct and set some env values
	app := conf.AppLog

	app.Port = os.Getenv("PORT") // will be available at http://localhost:8080
	if app.Port == "" {
		app.Port = ":8080" //nasty
	}

	app.API_key = os.Getenv("HIKE_API_KEY") // maybe get rid of it?
	if app.API_key == "" {
		app.ErrorLog.Fatal("Env: apiKey must be set")
	}

	// Testing the DuckDB
	tour_sample, err := getOneRow(1)
	if err != nil {
		app.ErrorLog.Fatal(err)
	}
	app.DebugLog.Printf("ID: %s, Name: %s\n", tour_sample.ID, tour_sample.Name)

	// my version of server, I can pass my own logger
	srv := &http.Server{
		Addr:     app.Port,
		ErrorLog: app.ErrorLog,
		Handler:  routes(app), // giving it my routes
	}

	app.InfoLog.Printf("Starting server on %s", app.Port)
	srv.ListenAndServe() //start the service and listen to the port with the mux

}
