package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mpfen/Go-Todo-REST-API/api/store"
)

// Handler for GET /projects/{projectName}/task/{taskName}
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

// Handler for POST /projects/{projectName}/task/{taskName}
func PostTaskHandler(p store.TodoStore, w http.ResponseWriter, r *http.Request) {
	// Decode task from request
	task := decodeTaskFromRequestOr400(w, r)

	vars := mux.Vars(r)
	projectName := vars["projectName"]
	taskName := task.Name

	// Check if project exists and get its id
	project := checkIfProjectExistsOr404(p, w, projectName)

	task.ProjectID = project.ID

	// Check if task already exists
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

// Handler for route GET /project/{projectName}/task
func GetAllProjectTasksHandler(p store.TodoStore, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectName := vars["projectName"]

	// Check if projects exists
	project := checkIfProjectExistsOr404(p, w, projectName)

	// Get all tasks
	tasks := p.GetAllProjectTasks(project)

	if len(tasks) == 0 {
		sendJSONResponse(w, fmt.Sprintf("No tasks in project %v found", projectName), http.StatusNotFound)
		return
	} else {
		w.Header().Set("content-type", jsonContentType)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(tasks)
	}
}

// Handler for route DELETE /projects/{projectName}/task/{taskName}
func DeleteTaskHandler(p store.TodoStore, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectName := vars["projectName"]
	taskName := vars["taskName"]

	// Check if tasks exists
	task := checkIfTasksExistsOr404(p, w, taskName, projectName)

	// Delete task
	err := p.DeleteTask(task)

	if err != nil {
		sendJSONResponse(w, fmt.Sprintf("Problem deleting Task: %v", err), http.StatusInternalServerError)
	} else {
		sendJSONResponse(w, "Task was successfully deleted", http.StatusOK)
	}
}

// Handler for route PUT /projects/{projectName}/task/{taskName}
func UpdateTaskHandler(p store.TodoStore, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectName := vars["projectName"]
	taskName := vars["taskName"]

	// Check if projects exists
	checkIfProjectExistsOr404(p, w, projectName)

	// Check if tasks exists
	task := checkIfTasksExistsOr404(p, w, taskName, projectName)

	// Decode task from request
	updatedTask := decodeTaskFromRequestOr400(w, r)

	// Update task
	task.Name = updatedTask.Name
	err := p.UpdateTask(task)

	if err != nil {
		sendJSONResponse(w, "Problem updating task", http.StatusInternalServerError)
		return
	} else {
		sendJSONResponse(w, "Task successfully updated", http.StatusOK)
		return
	}
}

// ComepleteTaskHandler PUT DELETE /projects/{projectName}/task/{taskName}/complete
// PUT completes task - DELETE reopens task
func CompleteTaskHandler(p store.TodoStore, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectName := vars["projectName"]
	taskName := vars["taskName"]

	// Check if projects exists
	checkIfProjectExistsOr404(p, w, projectName)

	// Check if tasks exists
	task := checkIfTasksExistsOr404(p, w, taskName, projectName)

	// Complete or reopen and update task
	var responseText string
	if r.Method == "PUT" {
		task.CompleteTask()
		responseText = "Task successfully completed"
	} else {
		task.ReopenTask()
		responseText = "Task successfully reopened"
	}

	// Update task
	err := p.UpdateTask(task)
	if err != nil {
		sendJSONResponse(w, "Problem upating task", http.StatusInternalServerError)
		return
	} else {
		sendJSONResponse(w, responseText, http.StatusOK)
	}
}
