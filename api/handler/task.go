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

	// Check if project exists and get its id
	project := p.GetProject(projectName)

	if project.Name == "" {
		sendJSONResponse(w, "No Project with that name exists", http.StatusBadRequest)
		return
	}

	task.ProjectID = project.ID

	// Check if the task already exists
	duplicateTask := p.GetTask(projectName, taskName)

	if duplicateTask.Name != "" {
		sendJSONResponse(w, "A Task with that name already exists for this project", http.StatusBadRequest)
		return
	}

	// Create new task
	err := p.PostTask(task)

	if err != nil {
		sendJSONResponse(w, "Task with the same name already exists", http.StatusBadRequest)
		return
	}

	sendJSONResponse(w, fmt.Sprintf("Task %v for project %v created", taskName, projectName), http.StatusCreated)
}
