package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mpfen/Go-Todo-REST-API/api/store"
)

// Handler for GET /project/{name}
func GetProjectHandler(p store.TodoStore, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectName := vars["name"]

	project := p.GetProject(projectName)

	if project.Name == "" {
		sendJSONResponse(w, "No project with this name found", http.StatusNotFound)
		return
	} else {
		w.Header().Set("content-type", jsonContentType)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(project)
	}
}

// Handler for POST /projects/
func PostProjectHandler(p store.TodoStore, w http.ResponseWriter, r *http.Request) {
	// Decode project from request
	project := decodeProjectFromRequestOr400(w, r)

	// Create new project
	err := p.PostProject(project.Name)

	if err != nil {
		sendJSONResponse(w, "Project with the same name already exists", http.StatusBadRequest)
		return
	}

	sendJSONResponse(w, "Project successfully created", http.StatusCreated)
}

// Handler for GET /projects/
func GetAllProjectsHandler(p store.TodoStore, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", jsonContentType)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(p.GetAllProjects())

}

// Handler for DELETE /projects/{name}
func DeleteProjectHandler(p store.TodoStore, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectName := vars["name"]

	// Check if project exists
	checkIfProjectExistsOr404(p, w, projectName)

	// Delete project if project exists
	err := p.DeleteProject(projectName)

	if err == nil {
		sendJSONResponse(w, "Project deleted", http.StatusOK)
		return
	} else {
		sendJSONResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Handler for PUT /projects/{name}
func UpdateProjectHandler(p store.TodoStore, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectName := vars["name"]

	// Check if project exists
	project := checkIfProjectExistsOr404(p, w, projectName)

	// Get new project name from request
	newProject := decodeProjectFromRequestOr400(w, r)

	project.Name = newProject.Name

	// Update project
	err := p.UpdateProject(project)

	if err != nil {
		sendJSONResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sendJSONResponse(w, "Project successfully updated", http.StatusOK)
}

// Handler for PUT DELETE /projects/{name}/archive
// PUT archived project - DELETE unarchives Project
func ArchiveProjectHandler(p store.TodoStore, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectName := vars["name"]

	// Check if project exists
	project := checkIfProjectExistsOr404(p, w, projectName)

	// Archive or unarchive project
	var responseText string
	if r.Method == "PUT" {
		project.ArchiveProject()
		responseText = "Project successfully archived"
	} else {
		project.UnArchiveProject()
		responseText = "Project successfully unarchived"
	}

	// Update project
	err := p.UpdateProject(project)

	if err != nil {
		sendJSONResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sendJSONResponse(w, responseText, http.StatusOK)
}
