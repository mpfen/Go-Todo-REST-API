package api

import (
	"encoding/json"
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
	PostProject(name string) error
}

type ProjectServer struct {
	Router *mux.Router
	Store  ProjectStore
}

// Initalize ProjectServer and create a gorilla/mux Router
func NewProjectServer(store ProjectStore) *ProjectServer {
	p := new(ProjectServer)
	p.Store = store

	p.Router = mux.NewRouter()
	p.Router.HandleFunc("/projects/{name}", p.GetProject).Methods("GET")
	p.Router.HandleFunc("/projects/", p.PostProject).Methods("POST")

	return p
}

func (p *ProjectServer) GetProject(w http.ResponseWriter, r *http.Request) {
	getProjectHandler(p.Store, w, r)
}

func (p *ProjectServer) PostProject(w http.ResponseWriter, r *http.Request) {
	postProjectHandler(p.Store, w, r)
}

// Handlerfunction for GET /project/{name}
func getProjectHandler(p ProjectStore, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectName := vars["name"]

	project := p.GetProject(projectName)

	if project.Name == "" {
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("content-type", jsonContentType)
		json.NewEncoder(w).Encode(map[string]string{"message": "No project with this name found"})
		return
	}

	w.Header().Set("content-type", jsonContentType)
	json.NewEncoder(w).Encode(project)
	w.WriteHeader(http.StatusOK)
}

func postProjectHandler(p ProjectStore, w http.ResponseWriter, r *http.Request) {
	project := model.Project{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&project); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("content-type", jsonContentType)
		json.NewEncoder(w).Encode(err.Error())
		return
	}
	defer r.Body.Close()

	err := p.PostProject(project.Name)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("content-type", jsonContentType)
		json.NewEncoder(w).Encode(map[string]string{"message": "Project with same name already exists"})
		return
	}

	w.WriteHeader(http.StatusCreated)
}
