package main

import (
	"encoding/json"
	//"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
)

// Definition of the Picture struct which is the response object of the JSON return from the endpoint.
type Picture struct {
	Copyright      string `json:"copyright"`
	Date           string `json:"date"`
	Explanation    string `json:"explanation"`
	Hdurl          string `json:"hdurl"`
	MediaType      string `json:"media_type"`
	ServiceVersion string `json:"service_version"`
	Title          string `json:"title"`
	URL            string `json:"url"`
}

func main() {
	//This is the URL for the Astronomy Picutre of the Day. This will use the Demo_Key as the api_key.
	var URL string = "https://api.nasa.gov/planetary/apod?api_key=DEMO_KEY"

	//Get the Astronomy Picture of the Day via the API and using the DEMO_KEY as the api_key query param.
	resp, err := http.Get(URL)

	//Check and handle the situation where there is an error with the GET request.
	if err != nil {
		log.Fatal(err)
	}

	//Defer closing the response body in case there is an error in the process.
	defer resp.Body.Close()

	//Read the content of the body of the response.
	body, err := io.ReadAll(resp.Body)

	//Check and handle error response from the ReadAll.
	if err != nil {
		log.Fatal(err)
	}

	//Declare a picture variable to store the struct in, and then unmarshall the
	//the JSON and store the result in this new variable.
	var picture Picture
	json.Unmarshal(body, &picture)

	// check our template for potential errors with Must
    tmpl := template.Must(template.ParseFiles("index.html"))

	tmpl.Execute(os.Stdout, picture)
}
