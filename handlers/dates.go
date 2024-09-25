package groupie

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

// Struct to hold the dates data
type Dates struct {
	Index []struct {
		ID    int      `json:"id"`
		Dates []string `json:"dates"`
	} `json:"index"`
}

func DatesHandler(w http.ResponseWriter, r *http.Request) {
	var error []string
	// Get the artist ID from the query parameters
	artistID := r.URL.Query().Get("id")
	if artistID == "" {
		fmt.Println(w, "Missing artist ID", http.StatusMethodNotAllowed)
		error = append(error, "Missing artist ID")
		ErrorHandler(w, r, http.StatusMethodNotAllowed, error)
		return
	}

	// Create a custom HTTP client with a timeout
	client := &http.Client{
		Timeout: 20 * time.Second, // 20-second timeout
	}

	// Make the GET request to fetch dates data
	resp, err := client.Get("https://groupietrackers.herokuapp.com/api/dates") // Update with the correct URL
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

	var dates Dates
	err = json.Unmarshal(body, &dates)
	if err != nil {
		http.Error(w, "Failed to parse JSON", http.StatusInternalServerError)
		error = append(error, "Internal Server Error")
		ErrorHandler(w, r, http.StatusInternalServerError, error)
		return
	}

	// Find the dates data for the requested artist ID
	var datesData struct {
		ID    int      `json:"id"`
		Dates []string `json:"dates"`
	}
	found := false
	for _, date := range dates.Index {
		id, err := strconv.Atoi(artistID)
		if err != nil {
			http.Error(w, "Invalid artist ID", http.StatusBadRequest)
			error = append(error, "Invalid artists ID")
			ErrorHandler(w, r, http.StatusBadRequest, error)
			return
		}
		if date.ID == id {
			datesData = date
			found = true
			break
		}
	}

	// If the artist ID is not found, return an error
	if !found {
		fmt.Println(w, "Artist ID not found", http.StatusBadRequest)
		error = append(error, "Artist ID not found")
		ErrorHandler(w, r, http.StatusBadRequest, error)
		return
	}

	// Return the dates data as JSON
	w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:8080")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(datesData); err != nil {
		fmt.Println(w, "Failed to encode JSON", http.StatusInternalServerError)
		error = append(error, "Internal Server Error")
		ErrorHandler(w, r, http.StatusInternalServerError, error)
		return
	}
}
