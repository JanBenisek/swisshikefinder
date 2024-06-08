package main

import (
	"net/http"

	"github.com/justinas/alice" // New import

	_ "github.com/marcboeker/go-duckdb"
)

// func (app *application) routes() *http.ServeMux {
func (app *application) routes() http.Handler {
	// because we put middleware before, we just return the handler, not mux

	// creates new HTTP server multiplexer
	// checks each requests and routes it to appropriate function
	mux := http.NewServeMux()

	// in index.html another endpoint is /static, we need to serve that ... I THINK???
	// we are giving it a file server (we need to serve static files), from which it serves the request
	// TODO: disable access to static files (through middleware)
	mux.Handle("/static/", http.FileServer(http.FS(Files))) //they are close and cached

	mux.HandleFunc("/search", app.searchHandler(3)) // with /search, use the searchHandler
	mux.HandleFunc("/", app.indexHandler())         // handles request to the root

	// using middleware here for every request
	// Recover panic is first to handle Panics in all subsequent middlewares and handlers
	// Create a middleware chain containing our 'standard' middleware
	// which will be used for every request our application receives.
	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	// Return the 'standard' middleware chain followed by the servemux.
	return standard.Then(mux)
}
