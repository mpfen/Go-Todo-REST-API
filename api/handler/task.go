package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mpfen/Go-Todo-REST-API/api/model"
	"github.com/mpfen/Go-Todo-REST-API/api/store"
)

// Handler for GET /projects/{projectName}/tasks/{taskName}
func GetTaskHandler(p store.TodoStore, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectName := vars["projectName"]
	taskName := vars["taskName"]

	task := p.GetTask(projectName, taskName)

	if task.Name == "" {
		sendJSONResponse(w, fmt.Sprintf("No task %v in project %v found", taskName, projectName), http.StatusNotFound)
		return
	} else {
		w.Header().Set("content-type", jsonContentType)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(task)
	}

}

// Handler for POST /projects/{projectName}/tasks/{taskName}
func PostTaskHandler(p store.TodoStore, w http.ResponseWriter, r *http.Request) {
	// Decode task from request
	task := model.Task{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&task); err != nil {
		sendJSONResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	vars := mux.Vars(r)
	projectName := vars["projectName"]
	taskName := task.Name

	// Create new task
	err := p.PostTask(task)

	if err != nil {
		sendJSONResponse(w, "Task with the same name already exists", http.StatusBadRequest)
		return
	}

	sendJSONResponse(w, fmt.Sprintf("Task %v for project %v created", taskName, projectName), http.StatusCreated)
}
