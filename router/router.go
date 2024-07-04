package router

import (
	"encoding/json"
	"lecoupeur/database"
	"lecoupeur/domain"
	"net/http"
)

type URLRequest struct {
	URL string `json:"url"`
}

type URLResponse struct {
	ShortURL string `json:"short_url"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func Router() {
	http.HandleFunc("/shorten", shortenHandler)
	http.HandleFunc("/r/", redirectHandler)
	http.Handle("/", http.FileServer(http.Dir("static")))

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

func shortenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var req URLRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid request payload"})
		if err != nil {
			return
		}
		return
	}

	var shortenUrl = domain.ShortenURL(req.URL)
	if shortenUrl == "" {
		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid URL"})
		if err != nil {
			return
		}
		return
	}

	response := URLResponse{ShortURL: shortenUrl}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		return
	}
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	// Get the ID from the URL path, so removing the first 3 characters (/r/)
	id := r.URL.Path[3:]
	url, err := database.GetURL(id)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	http.Redirect(w, r, url, http.StatusFound)
}
