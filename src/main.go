package main //belongs to the main package

import (
	"bytes"
	"html/template"
	"log"
	"net/http" // webserver
	"net/url"
	"os" // access os stuff
	"strconv"
	"time"

	"github.com/janbenisek/swiss-hike-finder/hikes"
)

// package level variable, we point it to template and parse it
// Must - panics if we fail to parse
var tpl = template.Must(template.ParseFiles("index.html"))

type Search struct {
	Query      string
	NextPage   int
	TotalPages int
	Results    *hikes.Results
}

func (s *Search) IsLastPage() bool {
	return s.NextPage >= s.TotalPages
}

func (s *Search) CurrentPage() int {
	if s.NextPage == 1 {
		return s.NextPage
	}
	return s.NextPage - 1
}

func (s *Search) PreviousPage() int {
	return s.CurrentPage() - 1
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	// has same params as HandleFunc
	// w - send responses to HTTP request
	// r - request received, we access the data
	buf := &bytes.Buffer{}
	err := tpl.Execute(buf, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	buf.WriteTo(w)
}

func searchHandler(hikesapi *hikes.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, err := url.Parse(r.URL.String())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		params := u.Query()
		searchQuery := params.Get("q")
		page := params.Get("page")
		if page == "" {
			page = "1"
		}

		results, err := hikesapi.FetchEverything(searchQuery, page)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// to debug
		// resultString := fmt.Sprintf("%+v", results)
		// fmt.Println("RESULT STRING: ", resultString)

		nextPage, err := strconv.Atoi(page)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		search := &Search{
			Query:      searchQuery,
			NextPage:   nextPage,
			TotalPages: results.Meta.Page.TotalPages,
			Results:    results,
		}

		if ok := !search.IsLastPage(); ok {
			search.NextPage++
		}

		buf := &bytes.Buffer{}
		err = tpl.Execute(buf, search)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		buf.WriteTo(w)

		// just for debugging
		// fmt.Println("Result: ", results)
		// fmt.Println("Result.Datas: ", results.Data)
		// fmt.Println("Result.Metadata: ", results.Meta)
		// fmt.Println("Result.Link: ", results.Links)
		// fmt.Println("Search Query is: ", searchQuery)
		// fmt.Println("Page is: ", page)
		// fmt.Println("Next page is: ", search.NextPage)
		// fmt.Println("Total pages is: ", search.TotalPages)
	}
}

func main() {
	// I do not think I need this, given the setup
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Println("Error loading .env file")
	// }

	port := os.Getenv("PORT") // will be available at http://localhost:3000
	if port == "" {
		port = "3000"
	}

	apiKey := os.Getenv("HIKE_API_KEY")
	if apiKey == "" {
		log.Fatal("Env: apiKey must be set")
	}

	myClient := &http.Client{Timeout: 10 * time.Second}
	hikesapi := hikes.NewClient(myClient, apiKey, 3) //how many hikes to show

	// register file server, so we can serve them in this request
	fs := http.FileServer(http.Dir("assets"))

	mux := http.NewServeMux()                                // creates new HTTP server multiplexer
	mux.Handle("/assets/", http.StripPrefix("/assets/", fs)) // use this file server with all requests with assets/

	mux.HandleFunc("/search", searchHandler(hikesapi))
	mux.HandleFunc("/", indexHandler)  // we match it with incoming request and call the associated handler
	http.ListenAndServe(":"+port, mux) //start the service and listen
}
