package groupie

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

type Locations struct {
	Index []struct {
		ID        int      `json:"id"`
		Locations []string `json:"locations"`
		Dates     string   `json:"dates"`
	} `json:"index"`
}

func LocationsHandler(w http.ResponseWriter, r *http.Request) {
	var error []string
	// Get the artist ID from the query parameters
	artistID := r.URL.Query().Get("id")
	if artistID == "" {
		fmt.Println(w, "Missing artist ID", http.StatusBadRequest)
		error = append(error, "Missing artist ID")
		ErrorHandler(w, r, http.StatusBadRequest, error)
		return
	}

	// Create a custom HTTP client with a timeout
	client := &http.Client{
		Timeout: 20 * time.Second, // 20-second timeout
	}

	// Make the GET request to fetch location data
	resp, err := client.Get("https://groupietrackers.herokuapp.com/api/locations") // Update with correct URL
	if err != nil {
		fmt.Println(w, "Failed to fetch data: "+err.Error(), http.StatusInternalServerError)
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

	var locations Locations
	err = json.Unmarshal(body, &locations)
	if err != nil {
		fmt.Println(w, "Failed to parse JSON", http.StatusInternalServerError)
		error = append(error, "Internal Server Error")
		ErrorHandler(w, r, http.StatusInternalServerError, error)
		return
	}

	// Find the location data for the requested artist ID
	var locationData struct {
		ID        int      `json:"id"`
		Locations []string `json:"locations"`
		Dates     string   `json:"dates"`
	}
	found := false
	for _, loc := range locations.Index {
		id, err := strconv.Atoi(artistID)
		if err != nil {
			fmt.Println(w, "Invalid artist ID", http.StatusBadRequest)
			error = append(error, "Invalid artist ID")
			ErrorHandler(w, r, http.StatusBadRequest, error)
			return
		}
		if loc.ID == id {
			locationData = loc
			found = true
			break
		}
	}
	if !found {
		fmt.Println(w, "Artist ID not found", http.StatusBadRequest)
		error = append(error, "Artist ID not found")
		ErrorHandler(w, r, http.StatusBadRequest, error)
		return
	}
	// Return the location data as JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(locationData); err != nil {
		fmt.Println(w, "Failed to encode JSON", http.StatusInternalServerError)
		error = append(error, "Internal Server Error")
		ErrorHandler(w, r, http.StatusInternalServerError, error)
		return
	}
}
