package handler

import (
	"encoding/json"
	"net/http"

	"github.com/mpfen/Go-Todo-REST-API/api/model"
	"github.com/mpfen/Go-Todo-REST-API/api/store"
)

const jsonContentType = "application/json"

// Sends a json response with specified message and httpStatusCode
func sendJSONResponse(w http.ResponseWriter, message string, code int) {
	w.Header().Set("content-type", jsonContentType)
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"message": message})

}

// Checks if a project with that name exists and returns the project or sends 404 message
func checkIfProjectExistsOr404(p store.TodoStore, w http.ResponseWriter, projectName string) model.Project {
	project := p.GetProject(projectName)

	if project.Name == "" {
		sendJSONResponse(w, "No project with this name found", http.StatusNotFound)
		return project
	}
	return project
}

// Checks if a task with that name exists in that project and returns the task or sends 404 message
func checkIfTasksExistsOr404(p store.TodoStore, w http.ResponseWriter, taskName, projectName string) model.Task {
	task := p.GetTask(projectName, taskName)

	if task.Name == "" {
		sendJSONResponse(w, "No task with that name exists", http.StatusNotFound)
		return task
	}
	return task
}

// Decodes a project struct from the request body. Returns it if successfull or send a http.StatusBadRequest
func decodeProjectFromRequestOr400(w http.ResponseWriter, r *http.Request) model.Project {
	project := model.Project{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&project); err != nil {
		sendJSONResponse(w, err.Error(), http.StatusBadRequest)
		return project
	}
	return project
}

// Decodes a task struct from the request body. Returns it if successfull or send a http.StatusBadRequest
func decodeTaskFromRequestOr400(w http.ResponseWriter, r *http.Request) model.Task {
	task := model.Task{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&task); err != nil {
		sendJSONResponse(w, err.Error(), http.StatusBadRequest)
		return task
	}
	return task
}
