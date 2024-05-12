package main

import (
	"net/http"
	"time"

	"internal/hikes"

	_ "github.com/marcboeker/go-duckdb"
)

func (app *application) routes() *http.ServeMux {

	// better to pass pointer to a client, than passing the whole client around, plus can modify it
	myClient := &http.Client{Timeout: 10 * time.Second} // create a new HTTP client with 10s timeout
	// not a pointer because the function returns a pointer
	hikesapi := hikes.NewClient(myClient, app.API_key, 3) // inits new client for the API with page size

	// creates new HTTP server multiplexer
	// checks each requests and routes it to appropriate function
	mux := http.NewServeMux()

	// in index.html another endpoint is /static, we need to serve that ... I THINK???
	// we are giving it a file server (we need to serve static files), from which it serves the request
	// TODO: disable access to static files (through middleware)
	mux.Handle("/static/", http.FileServer(http.FS(static))) //they are close and cached

	mux.HandleFunc("/search", app.searchHandler(hikesapi)) // with /search, use the searchHandler
	mux.HandleFunc("/", app.indexHandler())                // handles request to the root

	return mux
}
