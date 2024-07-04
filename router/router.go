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

func Router() *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("/shorten", shortenHandler)
	router.HandleFunc("/r/", redirectHandler)
	router.Handle("/", http.FileServer(http.Dir("static")))
	router.HandleFunc("/stats", statsHandler)
	return router
}

func StartServer() {
	router := Router()
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		panic(err)
	}
}

func shortenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_ = json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid request method"})
	}

	var req URLRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.URL == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid request payload"})
	}

	var shortenUrl = domain.ShortenURL(req.URL)
	if shortenUrl == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid URL"})
	}

	response := URLResponse{ShortURL: shortenUrl}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
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

func statsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_ = json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid request method"})
		return
	}

	urls, err := database.GetAllURLs()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to retrieve stats"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(urls)
	if err != nil {
		return
	}
}
