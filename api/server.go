package api

import (
	"fmt"
	"net/http"
	"strings"
)

type Project struct {
	Name string
}

// Projects only have a name right now
type ProjectStore interface {
	GetProjectInfo(name string) string
}

type ProjectServer struct {
	Store ProjectStore
}

func (p *ProjectServer) GetProjectInfo(name string) string {

	return name
}

func (p *ProjectServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	project := strings.TrimPrefix(r.URL.Path, "/projects/")
	fmt.Fprint(w, p.Store.GetProjectInfo(project))
}
