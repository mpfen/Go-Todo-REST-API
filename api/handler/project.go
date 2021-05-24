package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mpfen/Go-Todo-REST-API/api/model"
	"github.com/mpfen/Go-Todo-REST-API/api/store"
)

const jsonContentType = "application/json"

// Handler for GET /project/{name}
func GetProjectHandler(p store.ProjectStore, w http.ResponseWriter, r *http.Request) {
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
func PostProjectHandler(p store.ProjectStore, w http.ResponseWriter, r *http.Request) {
	project := model.Project{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&project); err != nil {
		sendJSONResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	err := p.PostProject(project.Name)

	if err != nil {
		sendJSONResponse(w, "Project with the same name already exists", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// Handler for GET /projects/
func GetAllProjectsHandler(p store.ProjectStore, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", jsonContentType)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(p.GetAllProjects())

}

// Handler for DELETE /projects/{name}
func DeleteProjectHandler(p store.ProjectStore, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectName := vars["name"]

	// Check if project exists
	project := p.GetProject(projectName)

	if project.Name == "" {
		sendJSONResponse(w, "No project with this name found", http.StatusNotFound)
		return
	}

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
func UpdateProjectHandler(p store.ProjectStore, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	oldProjectName := vars["name"]

	// Check if project exists
	oldProject := p.GetProject(oldProjectName)

	if oldProject.Name == "" {
		sendJSONResponse(w, "No project with this name found", http.StatusNotFound)
		return
	}

	// Update Project
	newProject := model.Project{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&newProject); err != nil {
		sendJSONResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	err := p.UpdateProject(oldProject.Name, newProject.Name)

	if err != nil {
		sendJSONResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sendJSONResponse(w, "Project successfully updated", http.StatusOK)
}

func sendJSONResponse(w http.ResponseWriter, message string, code int) {
	w.Header().Set("content-type", jsonContentType)
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"message": message})

}
