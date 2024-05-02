package main //belongs to the main package

import (
	// embed static files in the binary
	"database/sql"
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http" // webserver
	"os"       // access os stuff
	"time"

	"github.com/janbenisek/swiss-hike-finder/hikes"
	_ "github.com/marcboeker/go-duckdb"
)

// package level variables - means that it is available anywhere in this package

//go:embed all:static
var static embed.FS

var tpl = template.Must(template.ParseFS(static, "static/templates/index.html"))

type Search struct {
	Query      string
	NextPage   int
	TotalPages int
	Results    *hikes.Results // this will be a pointer
}

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

func (s *Search) IsLastPage() bool {
	// Operate on the struct Search,
	// returns bool (if last page)
	return s.NextPage >= s.TotalPages
}

func (s *Search) CurrentPage() int {
	// Operates on the struct Search
	// returns int (current page number)
	if s.NextPage == 1 {
		return s.NextPage
	}
	return s.NextPage - 1
}

func (s *Search) PreviousPage() int {
	// Operates on the struct Search
	// returns int (previous page number)
	return s.CurrentPage() - 1
}

func main() {

	// Testing the DuckDB
	tour_sample, err := getOneRow(1)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("ID: %s, Name: %s\n", tour_sample.ID, tour_sample.Name)

	port := os.Getenv("PORT") // will be available at http://localhost:8080
	if port == "" {
		port = "8080" //nasty
	}

	apiKey := os.Getenv("HIKE_API_KEY")
	if apiKey == "" {
		log.Fatal("Env: apiKey must be set")
	}

	// better to pass pointer to a client, than passing the whole client around, plus can modify it
	myClient := &http.Client{Timeout: 10 * time.Second} // create a new HTTP client with 10s timeout
	// not a pointer because the function returns a pointer
	hikesapi := hikes.NewClient(myClient, apiKey, 3) // inits new client for the API with page size

	// creates new HTTP server multiplexer
	// checks each requests and routes it to appropriate function
	mux := http.NewServeMux()

	// in index.html another endpoint is /static, we need to serve that ... I THINK???
	// we are giving it a file server (we need to serve static files), from which it serves the request
	mux.Handle("/static/", http.FileServer(http.FS(static))) //they are close and cached

	mux.HandleFunc("/search", searchHandler(hikesapi)) // with /search, use the searchHandler
	mux.HandleFunc("/", indexHandler)                  // handles request to the root
	http.ListenAndServe(":"+port, mux)                 //start the service and listen to the port with the mux

}
