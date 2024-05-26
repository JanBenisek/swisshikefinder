package main //belongs to the main package

import (
	// embed static files in the binary
	"errors"
	"math"
	"net/http" // webserver
	"net/url"  // access os stuff
	"strconv"

	"internal/models"

	_ "github.com/marcboeker/go-duckdb"
)

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

		app.InfoLog.Printf("Serving / endpoint")

		results, err := app.Tours.RandomTourPics(3)
		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				app.notFound(w)
			} else {
				app.serverError(w, err)
			}
			return
		}

		// debug
		// for _, pic := range results {
		// 	app.DebugLog.Printf("URL: %s\n", pic.PictureURL)
		// }

		data := app.newTemplateData(r)
		data.Home = &Home{Results: results,}

		// home := &Home{Results:results,}

		// app.render(w, http.StatusOK, "home.html", &Home{Results: results,})
		app.render(w, http.StatusOK, "home.html", data)
	}
}

func (app *application) searchHandler(pageSize int) http.HandlerFunc {
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

		app.InfoLog.Printf("Serving /search endpoint")

		u, err := url.Parse(r.URL.String()) // we parse the URL from the request
		if err != nil {
			app.serverError(w, err)
			return
		}

		app.InfoLog.Printf("Parsed URL: %s", u)

		params := u.Query()            // extract params from the query
		searchQuery := params.Get("q") // get value of the q param
		page := params.Get("page")     //get value of the page param
		if page == "" {
			page = "1"
		}
		app.InfoLog.Printf("Parsed parameters search: %s, page: %s", searchQuery, page)

		// all this seems very cumbersome
		pageInt, err := strconv.Atoi(page) //strconv is a package, Atoi is ASCII to integer
		if err != nil {
			app.serverError(w, err)
		}

		// assuming search has 11 records, page size=3
		// page 1: offset 0
		// page 2: offset 3
		// page 3: offset 6
		// page 4: offset 9 (will have only two records)
		offset := (pageInt - 1) * pageSize

		// here we call the API with the params from the request
		results, err := app.Tours.SearchTour(searchQuery, pageSize, offset)
		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				app.notFound(w)
			} else {
				app.serverError(w, err)
			}
			return
		}

		app.InfoLog.Printf("Obtained %d result(s) from DB", len(results))

		// to debug
		// for _, tour := range results {
		// 	app.DebugLog.Printf("ID: %s, RecordType: %s, Name: %s, Abstract: %s, Logo: %s, Count: %d\n", tour.ID, tour.RecordType, tour.Name, tour.Abstract, tour.LogoURL, tour.RecordCount)
		// }

		// this also looks iffy, no?
		var totalResults int
		var totalPages int

		if len(results) > 0 {
			totalResults = results[0].RecordCount
			totalPages = int(math.Ceil(float64(totalResults) / 3))
		} else {
			totalPages = 0
			totalResults = 0
		}

		// we create an instance of struct Search
		// we use pointer to avoid copying
		// if I want mutability outside, than pointer also makes sense
		search := &Search{
			Query:        searchQuery,
			NextPage:     pageInt,
			TotalPages:   totalPages,
			TotalResults: totalResults,
			Results:      results,
		}

		// to debug
		// resultStringB := fmt.Sprintf("%+v", search)
		// fmt.Println("BEFORE: ", resultStringB)
		app.DebugLog.Printf("Search struct: %+v", search)

		// increment page if page is not the last page
		// this is if with initialiser
		if ok := !search.IsLastPage(); ok {
			search.NextPage++
			app.InfoLog.Printf("Incremented next page to %d", search.NextPage)
		}

		data := app.newTemplateData(r)
		data.Search = search

		// render the app with the search results
		app.render(w, http.StatusOK, "search.html", data)

		app.InfoLog.Printf("Search request finished")
	}
}
