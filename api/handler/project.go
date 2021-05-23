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
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("content-type", jsonContentType)
		json.NewEncoder(w).Encode(map[string]string{"message": "No project with this name found"})
		return
	}

	w.Header().Set("content-type", jsonContentType)
	json.NewEncoder(w).Encode(project)
	w.WriteHeader(http.StatusOK)
}

// Handler for POST /projects/
func PostProjectHandler(p store.ProjectStore, w http.ResponseWriter, r *http.Request) {
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
