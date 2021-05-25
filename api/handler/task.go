package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mpfen/Go-Todo-REST-API/api/store"
)

func GetTaskHandler(p store.TodoStore, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectName := vars["projectName"]
	taskName := vars["taskName"]

	task := p.GetTask(projectName, taskName)

	if task.Name == "" {
		sendJSONResponse(w, fmt.Sprintf("No task %v in project %v found", projectName, taskName), http.StatusNotFound)
		return
	} else {
		w.Header().Set("content-type", jsonContentType)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(task)
	}

}
