package main

import "lecoupeur/router"
import "lecoupeur/database"

func main() {
	database.Connect()
	router.StartServer()
}
