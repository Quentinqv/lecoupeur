package domain

import (
	"math/rand"
	"net/url"
)

// Le but de ce package est de gérer la logique métier de l'application.
// Il s'agit d'un URL Shortener, donc la logique métier est de raccourcir les URL.
// Pour ce faire, nous avons besoin de stocker les URL dans une base de données.
// Nous avons besoin de générer des URL raccourcies.
// Nous avons besoin de rediriger les utilisateurs vers les URL originales.
// Nous avons besoin de gérer les erreurs.
// Nous avons besoin de gérer les statistiques.

const baseUrl = "http://localhost:8080/"

func ShortenURL(url string) string {
	if CheckURL(url) {
		id := GenerateUniqueID()
		// store the URL and the ID in the database
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
