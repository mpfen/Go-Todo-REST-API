package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mpfen/Go-Todo-REST-API/api/handler"
	"github.com/mpfen/Go-Todo-REST-API/api/store"
)

type TodoStore struct {
	Router *mux.Router
	Store  store.TodoStore
}

// Initalize TodoStore and create a gorilla/mux Router
func NewTodoStore(store store.TodoStore) *TodoStore {
	p := new(TodoStore)
	p.Store = store

	p.Router = mux.NewRouter()

	// Project routes
	p.Router.HandleFunc("/projects", p.PostProject).Methods("POST")
	p.Router.HandleFunc("/projects", p.GetAllProjects).Methods("GET")
	p.Router.HandleFunc("/projects/{name}", p.GetProject).Methods("GET")
	p.Router.HandleFunc("/projects/{name}", p.DeleteProject).Methods("DELETE")
	p.Router.HandleFunc("/projects/{name}", p.UpdateProject).Methods("PUT")
	p.Router.HandleFunc("/projects/{name}/archive", p.ArchiveProject).Methods("PUT")
	p.Router.HandleFunc("/projects/{name}/archive", p.UnArchiveProject).Methods("DELETE")

	// Task routes
	p.Router.HandleFunc("/projects/{projectName}/task/{taskName}", p.GetTask).Methods("GET")
	p.Router.HandleFunc("/projects/{projectName}/task", p.PostTask).Methods("POST")
	p.Router.HandleFunc("/projects/{projectName}/task", p.GetAllProjectTasks).Methods("GET")
	p.Router.HandleFunc("/projects/{projectName}/task/{taskName}", p.DeleteTask).Methods("DELETE")
	p.Router.HandleFunc("/projects/{projectName}/task/{taskName}", p.UpdateTask).Methods("PUT")
	return p
}

// Project Handler
func (p *TodoStore) GetProject(w http.ResponseWriter, r *http.Request) {
	handler.GetProjectHandler(p.Store, w, r)
}

func (p *TodoStore) PostProject(w http.ResponseWriter, r *http.Request) {
	handler.PostProjectHandler(p.Store, w, r)
}

func (p *TodoStore) GetAllProjects(w http.ResponseWriter, r *http.Request) {
	handler.GetAllProjectsHandler(p.Store, w, r)
}

func (p *TodoStore) DeleteProject(w http.ResponseWriter, r *http.Request) {
	handler.DeleteProjectHandler(p.Store, w, r)
}

func (p *TodoStore) UpdateProject(w http.ResponseWriter, r *http.Request) {
	handler.UpdateProjectHandler(p.Store, w, r)
}

func (p *TodoStore) ArchiveProject(w http.ResponseWriter, r *http.Request) {
	handler.ArchiveProjectHandler(p.Store, w, r)
}

func (p *TodoStore) UnArchiveProject(w http.ResponseWriter, r *http.Request) {
	handler.UnArchiveProjectHandler(p.Store, w, r)
}

// Task Handler

func (p *TodoStore) GetTask(w http.ResponseWriter, r *http.Request) {
	handler.GetTaskHandler(p.Store, w, r)
}

func (p *TodoStore) PostTask(w http.ResponseWriter, r *http.Request) {
	handler.PostTaskHandler(p.Store, w, r)
}

func (p *TodoStore) GetAllProjectTasks(w http.ResponseWriter, r *http.Request) {
	handler.GetAllProjectTasksHandler(p.Store, w, r)
}

func (p *TodoStore) DeleteTask(w http.ResponseWriter, r *http.Request) {
	handler.DeleteTaskHandler(p.Store, w, r)
}

func (p *TodoStore) UpdateTask(w http.ResponseWriter, r *http.Request) {
	handler.UpdateTaskHandler(p.Store, w, r)
}
