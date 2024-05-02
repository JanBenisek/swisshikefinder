package hikes

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Client struct {
	http     *http.Client // this will be a pointer, it is HTTP client to send HTTP request
	key      string
	PageSize int
}

func (c *Client) FetchEverything(query, page string) (*Results, error) {
	// Function fetches everything from an endpoint
	// operates on a struct Client, a pointer (is a method of that struct)
	// params:
	// query: string, query to search for
	// page: how many pages to return
	// returns: tuple, pointer to result and error

	// Define the endpoint URL
	endpoint := "https://opendata.myswitzerland.io/v1/destinations"

	// Construct the complete URL with query parameters
	url := fmt.Sprintf("%s?query=%s&facets=%s&lang=%s&hitsPerPage=%d&striphtml=%s&top=%s&page=%s", endpoint, url.QueryEscape(query), "*", "en", c.PageSize, "true", "top", page)

	fmt.Println("URL: ", url)

	// Create a new HTTP GET request object
	req, err := http.NewRequest("GET", url, nil) // third param is request body
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	// Set request headers
	req.Header.Set("accept", "application/json")
	req.Header.Set("x-api-key", c.key)

	// Send the HTTP request
	client := http.DefaultClient //default instance, holds a reference to HTTP client instance
	resp, err := client.Do(req)  // sends the actual request
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}
	// after HTTP request has been processed, close it. We need to release resources linked to the request
	defer resp.Body.Close()

	// Check if the response status code is OK
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non-OK status code received: %d", resp.StatusCode)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body) // returns byte slice containing the data
	if err != nil {
		return nil, fmt.Errorf("error reading response: %v", err)
	}
	// to debug
	// bodyString := string(body)
	// fmt.Println("BODY: ", bodyString)

	// Unmarshal JSON response
	res := &Results{}
	err = json.Unmarshal(body, res) // converts JSON into struct, (data, pointer to struct)
	if err != nil {
		// Handle unmarshal error
		fmt.Println(err.Error())
		return nil, err
	}

	return res, nil
}

func NewClient(httpClient *http.Client, key string, pageSize int) *Client {
	// Create a new instance of a client
	// params:
	// httpClient: pointer to a Client instance
	// key: key to the API
	// pageSize: size of page to return
	// returns pointer to a client
	if pageSize > 50 {
		pageSize = 50
	}

	return &Client{httpClient, key, pageSize}
}

// STRUCTS
type Results struct {
	Meta  Metadata `json:"meta"`
	Links Link     `json:"links"`
	Data  []Data   `json:"data"`
}

type Data struct {
	Context   string `json:"@context"`
	Type      string `json:"@type"`
	SubjectOf struct {
		Type    string `json:"@type"`
		License string `json:"license"`
	} `json:"subjectOf"`
	Links struct {
		Self string `json:"self"`
	} `json:"links"`
	Identifier string `json:"identifier"`
	Name       string `json:"name"`
	Abstract   string `json:"abstract"`
	URL        string `json:"url"`
	Photo      string `json:"photo"`
	Geo        struct {
		Type      string  `json:"@type"`
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"geo"`
	Classification []struct {
		Context string `json:"@context"`
		Type    string `json:"@type"`
		Name    string `json:"name"`
		Values  []struct {
			Name  string `json:"name"`
			Title string `json:"title"`
		} `json:"values"`
	} `json:"classification"`
}

type Link struct {
	Self  string `json:"self"`
	First string `json:"first"`
	Last  string `json:"last"`
	Next  string `json:"next"`
}

type Metadata struct {
	Language   string `json:"language"`
	APIVersion string `json:"apiVersion"`
	Page       struct {
		Size          int `json:"size"`
		TotalElements int `json:"totalElements"`
		TotalPages    int `json:"totalPages"`
		Number        int `json:"number"`
	} `json:"page"`
	Facets struct {
		Stv struct {
			Familydestination int `json:"familydestination"`
		} `json:"stv"`
		Guest struct {
			Carfreeplace int `json:"carfreeplace"`
		} `json:"guest"`
		Views struct {
			Mountainview int `json:"mountainview"`
			Panorama     int `json:"panorama"`
			Flatlandview int `json:"flatlandview"`
		} `json:"views"`
		Seasons struct {
			Summer int `json:"summer"`
			Winter int `json:"winter"`
			Autumn int `json:"autumn"`
			Spring int `json:"spring"`
		} `json:"seasons"`
		Placetypes struct {
			Villages       int `json:"villages"`
			Regions        int `json:"regions"`
			Mountains      int `json:"mountains"`
			Mountainlakes  int `json:"mountainlakes"`
			Cities         int `json:"cities"`
			Natureparks    int `json:"natureparks"`
			Valleys        int `json:"valleys"`
			Biglakes       int `json:"biglakes"`
			Rivers         int `json:"rivers"`
			Mountainpasses int `json:"mountainpasses"`
		} `json:"placetypes"`
		Reachability struct {
			Reachablebycar      int `json:"reachablebycar"`
			Reachablebybus      int `json:"reachablebybus"`
			Reachablebylocalbus int `json:"reachablebylocalbus"`
			Reachablebytrain    int `json:"reachablebytrain"`
			Reachablebyboat     int `json:"reachablebyboat"`
		} `json:"reachability"`
		Specialevents struct {
			OneAugust int `json:"1august"`
		} `json:"specialevents"`
		Naturspectacle struct {
			Sunset    int `json:"sunset"`
			Moonshine int `json:"moonshine"`
		} `json:"naturspectacle"`
		Altitudinalbelt struct {
			Alps       int `json:"alps"`
			Flatland   int `json:"flatland"`
			Beforealps int `json:"beforealps"`
		} `json:"altitudinalbelt"`
		Distancetoairport struct {
			Max2H         int `json:"max2h"`
			Max3H         int `json:"max3h"`
			Max1H         int `json:"max1h"`
			Max1H30Min    int `json:"max1h30min"`
			Morethan3H    int `json:"morethan3h"`
			Lessthan30Min int `json:"lessthan30min"`
		} `json:"distancetoairport"`
		Reachabilitylocation struct {
			Closetopublictransport int `json:"closetopublictransport"`
			Bycar                  int `json:"bycar"`
			Nexttobikepath         int `json:"nexttobikepath"`
		} `json:"reachabilitylocation"`
		Regionalspecialities struct {
			Family                  int `json:"family"`
			Skiingandsnowboarding   int `json:"skiingandsnowboarding"`
			Crosscountryskiing      int `json:"crosscountryskiing"`
			Fish                    int `json:"fish"`
			Snowshoeandwinterhiking int `json:"snowshoeandwinterhiking"`
			Wellness                int `json:"wellness"`
			Tobogganing             int `json:"tobogganing"`
			Wine                    int `json:"wine"`
			Meeting                 int `json:"meeting"`
		} `json:"regionalspecialities"`
		Geographicallocations struct {
			Alonggrandtour       int `json:"alonggrandtour"`
			Inthemountains       int `json:"inthemountains"`
			Inthecountryside     int `json:"inthecountryside"`
			Atthelake            int `json:"atthelake"`
			Inthecity            int `json:"inthecity"`
			Inthealpinemountains int `json:"inthealpinemountains"`
			Bytheriver           int `json:"bytheriver"`
		} `json:"geographicallocations"`
		Geographicalsituation struct {
			Westswitzerland     int `json:"westswitzerland"`
			Swisscentralplateau int `json:"swisscentralplateau"`
			Southernswitzerland int `json:"southernswitzerland"`
			Easternswitzerland  int `json:"easternswitzerland"`
		} `json:"geographicalsituation"`
		Geographicalspecialties struct {
			Karstlandscape int `json:"karstlandscape"`
		} `json:"geographicalspecialties"`
	} `json:"facets"`
}
