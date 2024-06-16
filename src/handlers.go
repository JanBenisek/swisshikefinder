package main //belongs to the main package

import (
	// embed static files in the binary
	"errors"
	"math"
	"net/http" // webserver
	"net/url"  // access os stuff
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/julienschmidt/httprouter"

	"internal/models"

	_ "github.com/marcboeker/go-duckdb"
)

func (app *application) indexHandler(w http.ResponseWriter, r *http.Request) {
	// Handles HTTP request to main page
	// Params:
	// w - send responses to HTTP request (from net/http)
	// r - request received, we access the data (from net/http)

	app.InfoLog.Printf("Serving / endpoint")

	results, err := app.Tours.RandomTourPics()
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
	// 	app.DebugLog.Printf("URL: %s, ID: %s\n", pic.ID, pic.PictureURL)
	// }

	data := app.newTemplateData(r)
	data.Home = &Home{Results: results}

	app.render(w, http.StatusOK, "home.html", data)
}

func (app *application) searchHandler(w http.ResponseWriter, r *http.Request) {
	// This handles the search endpoint

	app.InfoLog.Printf("Serving /search endpoint")

	var pageSize int = 3

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
	// app.DebugLog.Printf("Search struct: %+v", search)

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

func (app *application) tourView(w http.ResponseWriter, r *http.Request) {
	// Serve Detailed Tour View

	app.InfoLog.Printf("Serving /tour endpoint")

	params := httprouter.ParamsFromContext(r.Context())

	app.DebugLog.Printf("Tour View params: %s\n", params)

	id := params.ByName("id")
	if id == "" {
		app.notFound(w)
		return
	}

	app.DebugLog.Printf("Tour View ID: %s\n", id)

	tour, err := app.Tours.TourBasicInfo(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	data := app.newTemplateData(r)
	data.Tour = &Tour{Result: tour}

	// render the app with the search results
	app.render(w, http.StatusOK, "tour.html", data)

	app.InfoLog.Printf("Tour request finished")
}

func (app *application) recommendView(w http.ResponseWriter, r *http.Request) {
	// Let user recommend a hike or tour or whatever

	app.InfoLog.Printf("Serving /recommend GET endpoint")

	recom, err := app.Recoms.GetAll()
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}
	data := app.newTemplateData(r)
	data.Recom = &Recom{Results: recom}

	app.render(w, http.StatusOK, "recommend.html", data)

	app.InfoLog.Printf("Recommend request finished")
}

func (app *application) recommendPost(w http.ResponseWriter, r *http.Request) {
	// Let user recommend a hike or tour or whatever

	app.InfoLog.Printf("Serving /recommend POST endpoint")

	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := &models.Recommendation{
		ID:          0,
		Title:       r.PostForm.Get("title"),
		Description: r.PostForm.Get("description"),
		FieldErrors: map[string]string{},
	}

	// Check that the title value is not blank and is not more than 100
	// characters long. If it fails either of those checks, add a message to the // errors map using the field name as the key.
	if strings.TrimSpace(form.Title) == "" {
		form.FieldErrors["title"] = "This field cannot be blank"
	} else if utf8.RuneCountInString(form.Title) > 100 {
		form.FieldErrors["title"] = "This field cannot be more than 100 characters long"
	}

	// Check that the Content value isn't blank.
	if strings.TrimSpace(form.Description) == "" {
		form.FieldErrors["description"] = "This field cannot be blank"
	}

	// If there are any validation errors re-display the create.tmpl template,
	// passing in the snippetCreateForm instance as dynamic data in the Form
	// field. Note that we use the HTTP status code 422 Unprocessable Entity
	// when sending the response to indicate that there was a validation error.
	if len(form.FieldErrors) > 0 {
		data := app.newTemplateData(r)
		data.RecomForm = &RecomForm{Results: form}
		app.render(w, http.StatusUnprocessableEntity, "recommend.html", data)
		return
	}

	id, err := app.Recoms.Insert(form.Title, form.Description)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.InfoLog.Printf("Recommend request finished with ID: %d", id)

	http.Redirect(w, r, "/recommend", http.StatusSeeOther)
}
