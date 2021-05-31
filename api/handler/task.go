package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mpfen/Go-Todo-REST-API/api/model"
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
		sendJSONResponse(w, "No Project with that name exists", http.StatusNotFound)
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

// Handler for route GET /project/{projectName}/task
func GetAllProjectTasksHandler(p store.TodoStore, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectName := vars["projectName"]
	// Check if projects exists
	project := p.GetProject(projectName)

	if project.Name == "" {
		sendJSONResponse(w, "No project with that name exists", http.StatusNotFound)
		return
	}

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
	task := p.GetTask(projectName, taskName)

	if task.Name == "" {
		sendJSONResponse(w, "No task with that name exists", http.StatusNotFound)
		return
	}

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
	project := p.GetProject(projectName)

	if project.Name == "" {
		sendJSONResponse(w, "No project with that name exists", http.StatusNotFound)
		return
	}

	// Check if tasks exists
	task := p.GetTask(projectName, taskName)

	if task.Name == "" {
		sendJSONResponse(w, "No task with that name exists", http.StatusNotFound)
		return
	}

	// Decode task from request
	updatedTask := model.Task{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&updatedTask); err != nil {
		sendJSONResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

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

func CompleteTaskHandler(p store.TodoStore, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectName := vars["projectName"]
	taskName := vars["taskName"]

	// Check if projects exists
	project := p.GetProject(projectName)

	if project.Name == "" {
		sendJSONResponse(w, "No project with that name exists", http.StatusNotFound)
		return
	}

	// Check if tasks exists
	task := p.GetTask(projectName, taskName)

	if task.Name == "" {
		sendJSONResponse(w, "No task with that name exists", http.StatusNotFound)
		return
	}

	// Complete and update task
	task.CompleteTask()

	err := p.UpdateTask(task)
	if err != nil {
		sendJSONResponse(w, "Problem upating task", http.StatusInternalServerError)
		return
	} else {
		sendJSONResponse(w, "Task was completed", http.StatusOK)
	}
}
