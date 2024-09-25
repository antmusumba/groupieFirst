package groupie

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"time"
)

// Define a struct to match the structure of the API response
type Artist struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Members      []string `json:"members"`
	Locations    string   `json:"locations"`
	ConcertDates string   `json:"concertDates"`
	Relations    string   `json:"relations"`
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	var error []string
	// Create a custom HTTP client with a timeout
	client := &http.Client{
		Timeout: 20 * time.Second, // 20-second timeout
	}

	// Make the GET request with the custom client
	resp, err := client.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		fmt.Println(w, "Failed to get data from api", http.StatusInternalServerError)
		error = append(error, "Internal Server Error")
		ErrorHandler(w, r, http.StatusInternalServerError, error)
		return
	}
	defer resp.Body.Close()

	// Read and parse the JSON response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(w, "Failed to read response", http.StatusInternalServerError)
		error = append(error, "Internal Server Error")
		ErrorHandler(w, r, http.StatusInternalServerError, error)
		return
	}
	// fmt.Println(string(body))
	var artists []Artist
	err = json.Unmarshal(body, &artists)
	if err != nil {
		fmt.Println(w, "Failed to parse JSON", http.StatusInternalServerError)
		error = append(error, "Internal Server Error")
		ErrorHandler(w, r, http.StatusInternalServerError, error)
		return
	}

	// Load and parse the template
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		fmt.Println(w, err.Error(), http.StatusInternalServerError)
		error = append(error, "Internal Server Error")
		ErrorHandler(w, r, http.StatusInternalServerError, error)
		return
	}

	// Execute the template with the data
	err = tmpl.Execute(w, artists)
	if err != nil {
		fmt.Println(w, err.Error(), http.StatusInternalServerError)
		error = append(error, "Internal Server Error")
		ErrorHandler(w, r, http.StatusMethodNotAllowed, error)
		return
	}
}
