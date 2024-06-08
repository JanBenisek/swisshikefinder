package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"

	_ "github.com/marcboeker/go-duckdb"
)

// func (app *application) routes() *http.ServeMux {
func (app *application) routes() http.Handler {
	// because we put middleware before, we just return the handler, not mux

	// Initialise the router
	// checks each requests and routes it to appropriate function
	// hence I do not need to check in each Handler if request is GET/POST/...
	router := httprouter.New()

	// Show Not found for wrong 404
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	// in index.html another endpoint is /static, we need to serve that ... I THINK???
	// we are giving it a file server (we need to serve static files), from which it serves the request
	// TODO: disable access to static files (through middleware)
	fileServer := http.FileServer(http.FS(Files))
	router.Handler(http.MethodGet, "/static/*filepath", fileServer)

	router.HandlerFunc(http.MethodGet, "/", app.indexHandler())         // handles request to the root
	router.HandlerFunc(http.MethodGet, "/search", app.searchHandler(3)) // with /search, use the searchHandler

	// using middleware here for every request
	// Recover panic is first to handle Panics in all subsequent middlewares and handlers
	// Create a middleware chain containing our 'standard' middleware
	// which will be used for every request our application receives.
	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	// Return the 'standard' middleware chain followed by the servemux.
	return standard.Then(router)
}
