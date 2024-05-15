package main //belongs to the main package

import (
	// embed static files in the binary
	"errors"
	"html/template"
	"math"
	"net/http" // webserver
	"net/url"  // access os stuff
	"strconv"

	"internal/models"

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

func (app *application) searchHandler(page_size int) http.HandlerFunc {
	// This handles the search endpoint
	// it uses closure which actually servers the request
	// Params:
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

		// all this seems very cumbersome
		page_int, err := strconv.Atoi(page) //strconv is a package, Atoi is ASCII to integer
		if err != nil {
			app.serverError(w, err)
		}

		// assuming search has 11 records, page size=3
		// page 1: offset 0
		// page 2: offset 3
		// page 3: offset 6
		// page 4: offset 9 (will have only two records)
		offset := (page_int - 1) * page_size

		// here we call the API with the params from the request
		results, err := app.Tours.SearchTour(searchQuery, page_size, offset)
		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				app.notFound(w)
			} else {
				app.serverError(w, err)
			}
			return
		}

		// to debug
		// for _, tour := range results {
		// 	app.DebugLog.Printf("ID: %s, Record_type: %s, Name: %s, Abstract: %s, Logo: %s, Count: %d\n", tour.ID, tour.Record_type, tour.Name, tour.Abstract, tour.Logo_url, tour.Record_count)
		// }

		totalPages := int(math.Ceil(float64(results[0].Record_count) / 3))

		// we create an instance of struct Search
		// we use pointer to avoid copying
		// if I want mutability outside, than pointer also makes sense
		search := &Search{
			Query:        searchQuery,
			NextPage:     page_int,
			TotalPages:   totalPages,
			TotalResults: results[0].Record_count,
			Results:      results,
		}

		// debugging
		// resultStringB := fmt.Sprintf("%+v", search)
		// fmt.Println("BEFORE: ", resultStringB)

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
