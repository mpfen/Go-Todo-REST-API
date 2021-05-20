package main

import (
	"log"
	"net/http"

	api "github.com/mpfen/Go-Todo-REST-API/api"
)

func main() {
	server := &api.ProjectServer{}

	err := http.ListenAndServe(":5000", server)

	if err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}
