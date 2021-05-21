package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/mpfen/Go-Todo-REST-API/api/model"
)

const jsonContentType = "application/json"

// ProjectStore interface for testing
// Tests use own implementation
type ProjectStore interface {
	GetProject(name string) model.Project
}

type ProjectServer struct {
	Store ProjectStore
}

func (p *ProjectServer) GetProject(name string) model.Project {
	project := p.Store.GetProject(name)
	return project
}

// GET /projects/{name}
// response is json
func (p *ProjectServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	projectName := strings.TrimPrefix(r.URL.Path, "/projects/")

	project := p.GetProject(projectName)

	if project.Name == "" {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "")
		return
	}

	w.Header().Set("content-type", jsonContentType)
	json.NewEncoder(w).Encode(project)
	w.WriteHeader(http.StatusOK)
}
