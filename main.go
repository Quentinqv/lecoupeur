package main

import "lecoupeur/router"
import "lecoupeur/database"

func main() {
	// use the router package in router/router.go
	router.Router()
	database.Connect()
}
