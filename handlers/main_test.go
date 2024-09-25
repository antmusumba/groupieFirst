package groupie

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDatesHandler(t *testing.T) {
	tests := []struct {
		name           string
		query          string
		mockResponse   string
		mockStatusCode int
		expectedStatus int
	}{
		{
			name:           "Missing artist ID",
			query:          "",
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "Invalid artist ID",
			query:          "?id=abc",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Artist ID found",
			query:          "?id=1",
			mockResponse:   `{"index": [{"id": 1, "dates": ["2023-09-12", "2023-10-01"]}]}`,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Artist ID not found",
			query:          "?id=999",
			mockResponse:   `{"index": []}`,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/dates"+tt.query, nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			rr := httptest.NewRecorder()

			// Mock HTTP server to return the appropriate response
			mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.mockStatusCode)
				w.Write([]byte(tt.mockResponse))
			}))
			defer mockServer.Close()

			// Replace the external API call with a call to the mock server
			http.DefaultClient = mockServer.Client()

			DatesHandler(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("Handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}
		})
	}
}

func TestLocationsHandler(t *testing.T) {
	tests := []struct {
		name           string
		query          string
		mockResponse   string
		mockStatusCode int
		expectedStatus int
	}{
		{
			name:           "Missing artist ID",
			query:          "",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Invalid artist ID",
			query:          "?id=abc",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Artist ID found",
			query:          "?id=1",
			mockResponse:   `{"index": [{"id": 1, "locations": ["New York", "Los Angeles"], "dates": "2023-09-12"}]}`,
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/locations"+tt.query, nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			rr := httptest.NewRecorder()

			// Mock HTTP server to return the appropriate response
			mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.mockStatusCode)
				w.Write([]byte(tt.mockResponse))
			}))
			defer mockServer.Close()

			// Replace the external API call with a call to the mock server
			http.DefaultClient = mockServer.Client()

			LocationsHandler(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("Handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}
		})
	}
}

func TestRelationHandler(t *testing.T) {
	tests := []struct {
		name           string
		query          string
		mockResponse   string
		mockStatusCode int
		expectedStatus int
	}{
		{
			name:           "Missing artist ID",
			query:          "",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Invalid artist ID",
			query:          "?id=abc",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Artist ID found",
			query:          "?id=1",
			mockResponse:   `{"index": [{"id": 1, "datesLocations": {"New York": ["2023-09-12"], "Los Angeles": ["2023-09-15"]}}]}`,
			mockStatusCode: http.StatusOK,
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new HTTP request
			req, err := http.NewRequest("GET", "/relation"+tt.query, nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			// Create a response recorder to capture the handler's response
			rr := httptest.NewRecorder()

			// Mock HTTP server to return the appropriate mock response
			mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.mockStatusCode)
				w.Write([]byte(tt.mockResponse))
			}))
			defer mockServer.Close()

			// Replace the external API call with a call to the mock server
			// To simulate calling the real API endpoint but with a mock response
			http.DefaultClient = mockServer.Client()

			// Call the handler
			RelationHandler(rr, req)

			// Check if the status code is what we expect
			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("Handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}
		})
	}
}
