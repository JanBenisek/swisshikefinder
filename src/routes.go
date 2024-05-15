package main

import (
	"net/http"

	_ "github.com/marcboeker/go-duckdb"
)

func (app *application) routes() *http.ServeMux {

	// creates new HTTP server multiplexer
	// checks each requests and routes it to appropriate function
	mux := http.NewServeMux()

	// in index.html another endpoint is /static, we need to serve that ... I THINK???
	// we are giving it a file server (we need to serve static files), from which it serves the request
	// TODO: disable access to static files (through middleware)
	mux.Handle("/static/", http.FileServer(http.FS(static))) //they are close and cached

	mux.HandleFunc("/search", app.searchHandler(3)) // with /search, use the searchHandler
	mux.HandleFunc("/", app.indexHandler())         // handles request to the root

	return mux
}
