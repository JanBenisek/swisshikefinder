package main

import (
	"fmt"
	conf "internal/config"
	"internal/hikes"
	"net/http"
	"runtime/debug"
)

// The serverError helper writes an error message and stack trace to the errorLog,
// then sends a generic 500 Internal Server Error response to the user.
func serverError(app *conf.Application, w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.ErrorLog.Output(2, trace) // this shows where error occured
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// The clientError helper sends a specific status code and corresponding description
// to the user. We'll use this later in the book to send responses like 400 "Bad
// Request" when there's a problem with the request that the user sent.
func clientError(app *conf.Application, w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// For consistency, we'll also implement a notFound helper. This is simply a
// convenience wrapper around clientError which sends a 404 Not Found response to
// the user.
func notFound(app *conf.Application, w http.ResponseWriter) {
	clientError(app, w, http.StatusNotFound)
}

// TODO: Better place to put these?
type Search struct {
	Query      string
	NextPage   int
	TotalPages int
	Results    *hikes.Results // this will be a pointer
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
