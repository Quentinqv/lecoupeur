package router

import (
	"bytes"
	"encoding/json"
	"lecoupeur/database"
	"lecoupeur/domain"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setup() {
	database.Connect()
}

func TestShortenHandler(t *testing.T) {
	setup()

	type args struct {
		url string
	}
	tests := []struct {
		name           string
		args           args
		wantStatusCode int
		wantShortURL   bool
	}{
		{
			name: "Shorten valid URL",
			args: args{
				url: "https://example.com",
			},
			wantStatusCode: http.StatusOK,
			wantShortURL:   true,
		},
		{
			name: "Shorten invalid URL",
			args: args{
				url: "example.com",
			},
			wantStatusCode: http.StatusBadRequest,
			wantShortURL:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload := URLRequest{URL: tt.args.url}
			body, _ := json.Marshal(payload)
			req, err := http.NewRequest("POST", "/shorten", bytes.NewBuffer(body))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(shortenHandler)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.wantStatusCode {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.wantStatusCode)
			}

			var response URLResponse
			err = json.NewDecoder(rr.Body).Decode(&response)
			if err != nil {
				t.Errorf("handler returned an invalid response: %v", err)
			}

			if tt.wantShortURL && response.ShortURL == "" {
				t.Errorf("handler returned an empty short URL")
			}
		})
	}
}

func TestRedirectHandler(t *testing.T) {
	setup()

	type args struct {
		id string
	}
	tests := []struct {
		name           string
		args           args
		wantStatusCode int
		wantLocation   string
	}{
		{
			name: "Redirect to valid URL",
			args: args{
				id: domain.GenerateUniqueID(),
			},
			wantStatusCode: http.StatusFound,
			wantLocation:   "https://example.com",
		},
		{
			name: "Redirect with invalid ID",
			args: args{
				id: "nonexistent",
			},
			wantStatusCode: http.StatusNotFound,
			wantLocation:   "",
		},
	}

	// Store a URL for the valid test case
	database.StoreURL("https://example.com", tests[0].args.id)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/r/"+tt.args.id, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(redirectHandler)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.wantStatusCode {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.wantStatusCode)
			}

			if location := rr.Header().Get("Location"); location != tt.wantLocation {
				t.Errorf("handler returned wrong location header: got %v want %v", location, tt.wantLocation)
			}
		})
	}
}
