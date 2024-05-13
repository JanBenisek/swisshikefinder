package main //belongs to the main package

import (
	// embed static files in the binary
	"html/template"
	"net/http" // webserver
	"net/url"  // access os stuff
	"strconv"

	"internal/hikes"

	_ "github.com/marcboeker/go-duckdb"
)

var tpl = template.Must(template.ParseFS(static, "static/templates/index.html"))

func (app *application) indexHandler() http.HandlerFunc {
	// Handles HTTP requests
	// Params:
	// w - send responses to HTTP request (from net/http)
	// r - request received, we access the data (from net/http)
	return func(w http.ResponseWriter, r *http.Request) {
		// buf is a pointer (&) which is nice thing to pass around, rather than copying the entire content
		if r.Method != http.MethodGet {
			w.Header().Set("Allow", http.MethodGet)
			app.clientError(w, http.StatusMethodNotAllowed)
			return
		}

		// write to template using HTTP writer and return it
		err := tpl.Execute(w, nil)
		if err != nil {
			app.serverError(w, err)
			return
		}
	}
}

func (app *application) searchHandler(hikesapi *hikes.Client) http.HandlerFunc {
	// This handles the search endpoint
	// it uses closure which actually servers the request
	// Params:
	// pointer to hikesapi
	// Returns HandlerFunc function
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.Header().Set("Allow", http.MethodGet)
			app.clientError(w, http.StatusMethodNotAllowed)
			return
		}

		u, err := url.Parse(r.URL.String()) // we parse the URL from the request
		if err != nil {
			app.serverError(w, err)
			return
		}

		params := u.Query()            // extract params from the query
		searchQuery := params.Get("q") // get value of the q param
		page := params.Get("page")     //get value of the page param
		if page == "" {
			page = "1"
		}

		// here we call the API with the params from the request
		results, err := hikesapi.FetchEverything(searchQuery, page)
		if err != nil {
			app.serverError(w, err)
			return
		}

		// to debug
		// resultString := fmt.Sprintf("%+v", results)
		// fmt.Println("RESULT STRING: ", resultString)

		nextPage, err := strconv.Atoi(page) //strconv is a package, Atoi is ASCII to integer
		if err != nil {
			app.notFound(w)
			return
		}

		// we create an instance of struct Search
		// we use pointer to avoid copying
		// if I want mutability outside, than pointer also makes sense
		search := &Search{
			Query:      searchQuery,
			NextPage:   nextPage,
			TotalPages: results.Meta.Page.TotalPages,
			Results:    results,
		}

		// increment page if page is not the last page
		// this is if with initialiser
		if ok := !search.IsLastPage(); ok {
			search.NextPage++
		}

		// this time we pass search data into the template and write into the HTTP response writer
		err = tpl.Execute(w, search)
		if err != nil {
			app.serverError(w, err)
			return
		}
	}
}
