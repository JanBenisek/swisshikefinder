package main //belongs to the main package

import (
	"bytes" // embed static files in the binary
	"embed"
	"html/template"
	"log"
	"net/http" // webserver
	"net/url"
	"os" // access os stuff
	"strconv"
	"time"

	"github.com/janbenisek/swiss-hike-finder/hikes"
)

// package level variables - means that it is available anywhere in this package

//go:embed all:static
var static embed.FS

var tpl = template.Must(template.ParseFS(static, "static/templates/index.html"))

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

func indexHandler(w http.ResponseWriter, r *http.Request) {
	// Handles HTTP requests
	// Params:
	// w - send responses to HTTP request (from net/http)
	// r - request received, we access the data (from net/http)

	// buf is a pointer (&) which is nice thing to pass around, rather than copying the entire content
	buf := &bytes.Buffer{}       //the buffer stores results of executing a template in a memory
	err := tpl.Execute(buf, nil) // executes the template and writes to buffer, (writer, data to pass into the template)
	if err != nil {
		// if there is an error, we return 500
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// write content of the buffer to the writer
	// this sends the generated HTML template as a response
	buf.WriteTo(w)
}

func searchHandler(hikesapi *hikes.Client) http.HandlerFunc {
	// This handles the search endpoint
	// it uses closure which actually servers the request
	// Params:
	// pointer to hikesapi
	// Returns HandlerFunc function
	return func(w http.ResponseWriter, r *http.Request) {
		u, err := url.Parse(r.URL.String()) // we parse the URL from the request
		if err != nil {
			// if error, we return 500
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		params := u.Query()            // extract params from the query
		searchQuery := params.Get("q") // get value of the q param
		page := params.Get("page")     //get value of the page param
		if page == "" {
			// if page if not set, set to 1
			page = "1"
		}

		// here we call the API with the params from the request
		results, err := hikesapi.FetchEverything(searchQuery, page)
		if err != nil {
			// as always, return 500 if error
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// to debug
		// resultString := fmt.Sprintf("%+v", results)
		// fmt.Println("RESULT STRING: ", resultString)

		nextPage, err := strconv.Atoi(page) //strconv is a package, Atoi is ASCII to integer
		if err != nil {
			// return 500 if error
			http.Error(w, err.Error(), http.StatusInternalServerError)
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

		// again create a buffer in memory and write the HTML into it
		buf := &bytes.Buffer{}
		err = tpl.Execute(buf, search) // this time we pass search data into the template
		if err != nil {
			// 500 if error
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// write content of the buffer into to the HTTP response writer
		buf.WriteTo(w)

		// just for debugging
		// fmt.Println("Result: ", results)
		// fmt.Println("Result.Metadata: ", results.Meta)
		// fmt.Println("Result.Link: ", results.Links)
		// fmt.Println("Search Query is: ", searchQuery)
		// fmt.Println("Page is: ", page)
		// fmt.Println("Next page is: ", search.NextPage)
		// fmt.Println("Total pages is: ", search.TotalPages)
	}
}

func main() {

	port := os.Getenv("PORT") // will be available at http://localhost:8080
	if port == "" {
		port = "8080" //nasty
	}

	apiKey := os.Getenv("HIKE_API_KEY")
	if apiKey == "" {
		log.Fatal("Env: apiKey must be set")
	}

	// better to pass pointer to a client, than passing the whole client around, plus can modify it
	myClient := &http.Client{Timeout: 10 * time.Second} // create a new HTTP client with 10s timeout
	// not a pointer because the function returns a pointer
	hikesapi := hikes.NewClient(myClient, apiKey, 3) // inits new client for the API with page size

	// register file server, so we can serve them in this request
	// we need to serve the assets to the client, this is a nice way, also they are close and cached

	// var staticFS = http.FS(static)
	// fs := http.FileServer(staticFS)

	// creates new HTTP server multiplexer
	// checks each requests and routes it to appropriate function
	mux := http.NewServeMux()

	// in index.html another endpoint is /assets, we need to serve that ... I THINK???
	//mux.Handle("/assets/", http.StripPrefix("/static/", fs)) // use this file server with all requests with assets/
	// mux.Handle("/assets/", fs)
	// mux.Handle("/assets/", fs)
	mux.Handle("/static/", http.FileServer(http.FS(static)))

	mux.HandleFunc("/search", searchHandler(hikesapi)) // with /search, use the searchHandler
	mux.HandleFunc("/", indexHandler)                  // handles request to the root
	http.ListenAndServe(":"+port, mux)                 //start the service and listen to the port with the mux
}
