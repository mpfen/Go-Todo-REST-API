package main

import (
	"log"
	"net/http"

	"github.com/mpfen/Go-Todo-REST-API/api"
)

func main() {
	db := api.NewDatabaseConnection()
	server := &api.ProjectServer{db}

	err := http.ListenAndServe(":5000", server)

	if err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}
