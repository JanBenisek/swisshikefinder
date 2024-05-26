package main

import (
	"fmt"
	"bytes"
	"html/template"
	"internal/models"
	"io/fs"
	"net/http"
	"path/filepath"
	"runtime/debug"
)

// The serverError helper writes an error message and stack trace to the errorLog,
// then sends a generic 500 Internal Server Error response to the user.
func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.ErrorLog.Output(2, trace) // this shows where error occured
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// The clientError helper sends a specific status code and corresponding description
// to the user. We'll use this later in the book to send responses like 400 "Bad
// Request" when there's a problem with the request that the user sent.
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// For consistency, we'll also implement a notFound helper. This is simply a
// convenience wrapper around clientError which sends a 404 Not Found response to
// the user.
func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

// TODO: Better place to put these?
type Search struct {
	Query        string
	NextPage     int
	TotalPages   int
	TotalResults int
	Results      []*models.Tour // this will be a pointer
}

type Home struct {
	Results []*models.TourPicture // this will be a pointer
}

// Struct that holds all data passed to the template
// Let's revisit later if I want a generic struct
// type TemplateData struct {
// 	Search *Search
// }

func (s *Search) IsLastPage() bool {
	// Operate on the struct Search,
	// returns bool (if last page)
	return s.NextPage > s.TotalPages
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

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	// Use fs.Glob() to get a slice of all filepaths in the Files
	// This essentially gives us a slice of all the 'page' templates for the application, just
	pages, err := fs.Glob(Files, "static/templates/pages/*.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		// Create a slice containing the filepath patterns for the templates we
		// want to parse.
		patterns := []string{
			"static/templates/base.html",
			"static/templates/partials/header.html",
			page,
		}

		// ts, err := template.New(name).Funcs(functions).ParseFS(Files, patterns...)
		ts, err := template.ParseFS(Files, patterns...)
		if err != nil {
			return nil, err
		}

		cache[name] = ts

	}
	return cache, nil
}

func (app *application) render(w http.ResponseWriter, status int, page string, data interface{}) {
	// Retrieve the appropriate template set from the cache based on the page
	// name (like 'home.tmpl'). If no entry exists in the cache with the
	// provided name, then create a new error and call the serverError() helper
	// method that we made earlier and return.
	//
	// data is interface, meaning it accepts any struct
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.serverError(w, err)
		return
	}
	// write first to buffer to catch parsing errors
	buf := new(bytes.Buffer)
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Write out the provided HTTP status code to the response header.
	w.WriteHeader(status)

	// Write the contents of the buffer to the http.ResponseWriter.
	buf.WriteTo(w)
}
