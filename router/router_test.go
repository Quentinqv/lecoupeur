package router

import (
	"bytes"
	"encoding/json"
	"lecoupeur/database"
	"lecoupeur/domain"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func setup() {
	database.Connect()
	database.FlushAll()
}

func TestShortenHandler(t *testing.T) {
	setup()

	type args struct {
		url    string
		method string
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
				url:    "https://example.com",
				method: "POST",
			},
			wantStatusCode: http.StatusOK,
			wantShortURL:   true,
		},
		{
			name: "Shorten invalid URL",
			args: args{
				url:    "example.com",
				method: "POST",
			},
			wantStatusCode: http.StatusBadRequest,
			wantShortURL:   false,
		},
		{
			name: "Invalid request method",
			args: args{
				url:    "https://example.com",
				method: "GET",
			},
			wantStatusCode: http.StatusMethodNotAllowed,
			wantShortURL:   false,
		},
		{
			name: "Invalid request payload",
			args: args{
				url:    "",
				method: "POST",
			},
			wantStatusCode: http.StatusBadRequest,
			wantShortURL:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload := URLRequest{URL: tt.args.url}
			body, _ := json.Marshal(payload)
			req, err := http.NewRequest(tt.args.method, "/shorten", bytes.NewBuffer(body))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			router := Router()
			router.ServeHTTP(rr, req)

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
			router := Router()
			router.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.wantStatusCode {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.wantStatusCode)
			}

			if location := rr.Header().Get("Location"); location != tt.wantLocation {
				t.Errorf("handler returned wrong location header: got %v want %v", location, tt.wantLocation)
			}
		})
	}
}

func Test_statsHandler(t *testing.T) {
	setup()
	database.StoreURL("https://example.com", "123456")
	database.StoreURL("https://example.fr", "abced")

	tests := []struct {
		name       string
		want       string
		method     string
		statusCode int
	}{
		{
			name:       "Get all URLs",
			want:       `{"123456":{"url":"https://example.com","counter":0},"abced":{"url":"https://example.fr","counter":0}}`,
			method:     "GET",
			statusCode: http.StatusOK,
		},
		{
			name:       "Bad request method",
			want:       `{"error":"Invalid request method"}`,
			method:     "POST",
			statusCode: http.StatusMethodNotAllowed,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, "/stats", nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			router := Router()
			router.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.statusCode {
				t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
			}

			if strings.Replace(rr.Body.String(), "\n", "", -1) != tt.want {
				t.Errorf("handler returned wrong body: got %v want %v", rr.Body.String(), tt.want)
			}
		})
	}
}
