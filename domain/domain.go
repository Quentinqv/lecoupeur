package domain

import (
	"lecoupeur/database"
	"math/rand"
	"net/url"
)

const baseUrl = "http://localhost:8080/r/"

func ShortenURL(url string) string {
	if CheckURL(url) {
		id := GenerateUniqueID()
		// Store the URL and the ID in the database
		database.StoreURL(url, id)
		return baseUrl + id
	}
	return ""
}

func GenerateUniqueID() string {
	const idLength = 6
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	id := make([]byte, idLength)
	for i := range id {
		id[i] = charset[rand.Intn(len(charset))]
	}
	return string(id)
}

func CheckURL(urlInput string) bool {
	u, err := url.ParseRequestURI(urlInput)
	if err != nil {
		return false
	}
	return u.Scheme != "" && u.Host != ""
}
