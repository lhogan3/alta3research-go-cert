package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
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

// Opens index.html in a browser depending on the OS.
func openwebpage() {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", "index.html").Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", "index.html").Start()
	case "darwin":
		err = exec.Command("open", "index.html").Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}

}

// Formats date string to go from YYYY-MM-DD -> MM/DD/YYYY.
func formatdate(date string) string {
	//split the date by the dashes.
	s := strings.Split(date, "-")
	//reorder to be in the "MM/DD/YYYY" format, abd then return
	result := fmt.Sprintf("%s/%s/%s", s[1], s[2], s[0])
	return result
}

// Fetches the picture of the day information from the NASA endpoint, marshalls the respoonse into a Picture
// struct and returns it.
func fetchpictureoftheday() Picture {
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

	// Format the date to be "MM/DD/YYYY"
	picture.Date = formatdate(picture.Date)

	return picture
}

// Generates an index.html file leveraging the populated picture struct and
// index.html.template file.
func generateindexfileviatemplate(picture Picture) {
	// check our template for potential errors with Must
	tmpl := template.Must(template.ParseFiles("index.html.template"))

	//Create the index.html file and then write to the newly created file the template and sruct.
	file, err := os.Create("index.html")

	//Check and handle error response from file creation.
	if err != nil {
		log.Fatal(err)
	}

	//Defer closing the file if an eror occurs.
	defer file.Close()

	//Write the filled out template to the index.html file.
	tmpl.Execute(file, picture)

	//Write the filled out template to the console.
	fmt.Print("The filled out index.html file is printed below displaying the Astronomy Picture of the Day:\n\n")
	tmpl.Execute(os.Stdout, picture)
}

func main() {
	//Fetch a populated picture struct of the Picture of the Day.
	var picture Picture = fetchpictureoftheday()

	//Generate an index.html file via the index.html.template.
	generateindexfileviatemplate(picture)

	//Open the newly-generated file in the appropriate browser.
	openwebpage()
}
