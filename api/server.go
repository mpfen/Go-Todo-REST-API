package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mpfen/Go-Todo-REST-API/api/model"
)

const jsonContentType = "application/json"

// ProjectStore interface for testing
// Tests use own implementation with
// StubStore instead of a real database
type ProjectStore interface {
	GetProject(name string) model.Project
}

type ProjectServer struct {
	Router *mux.Router
	Store  ProjectStore
}

func NewProjectServer(store ProjectStore) *ProjectServer {
	p := new(ProjectServer)
	p.Store = store

	p.Router = mux.NewRouter()
	p.Router.HandleFunc("/projects/{name}", p.GetProject).Methods("GET")

	return p
}

func (p *ProjectServer) GetProject(w http.ResponseWriter, r *http.Request) {
	getProjectHandler(p.Store, w, r)
}

func getProjectHandler(p ProjectStore, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectName := vars["name"]
	log.Print("Handlername:" + projectName)

	project := p.GetProject(projectName)

	if project.Name == "" {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "")
		return
	}

	w.Header().Set("content-type", jsonContentType)
	json.NewEncoder(w).Encode(project)
	w.WriteHeader(http.StatusOK)
}
