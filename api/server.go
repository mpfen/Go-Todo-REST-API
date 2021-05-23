package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mpfen/Go-Todo-REST-API/api/handler"
	"github.com/mpfen/Go-Todo-REST-API/api/store"
)

type ProjectServer struct {
	Router *mux.Router
	Store  store.ProjectStore
}

// Initalize ProjectServer and create a gorilla/mux Router
func NewProjectServer(store store.ProjectStore) *ProjectServer {
	p := new(ProjectServer)
	p.Store = store

	p.Router = mux.NewRouter()
	p.Router.HandleFunc("/projects/{name}", p.GetProject).Methods("GET")
	p.Router.HandleFunc("/projects/", p.PostProject).Methods("POST")
	p.Router.HandleFunc("/projects", p.GetAllProjects).Methods("GET")
	p.Router.HandleFunc("/projects/{name}", p.DeleteProject).Methods("DELETE")

	return p
}

func (p *ProjectServer) GetProject(w http.ResponseWriter, r *http.Request) {
	handler.GetProjectHandler(p.Store, w, r)
}

func (p *ProjectServer) PostProject(w http.ResponseWriter, r *http.Request) {
	handler.PostProjectHandler(p.Store, w, r)
}

func (p *ProjectServer) GetAllProjects(w http.ResponseWriter, r *http.Request) {
	handler.GetAllProjectsHandler(p.Store, w, r)
}

func (p *ProjectServer) DeleteProject(w http.ResponseWriter, r *http.Request) {
	handler.DeleteProjectHandler(p.Store, w, r)
}
