package main

import (
	"log"
	"net/http"

	"github.com/mpfen/Go-Todo-REST-API/api"
)

func main() {
	db := api.NewDatabaseConnection("database.db")
	server := api.NewProjectServer(db)

	err := http.ListenAndServe(":5000", server.Router)

	if err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}
