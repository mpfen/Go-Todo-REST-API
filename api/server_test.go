package api_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	api "github.com/mpfen/Go-Todo-REST-API/api"
	"github.com/mpfen/Go-Todo-REST-API/api/model"
)

// DB store stub for testing
type StubProjectStore struct {
	projects map[string]string
}

// creates a makeshift project struct to comply with ProjectStore interface
func (s *StubProjectStore) GetProject(name string) model.Project {
	project := model.Project{}
	project.Name = s.projects[name]
	return project
}

// creates a new project
func (s *StubProjectStore) PostProject(name string) error {
	if duplicate := s.projects[name]; duplicate != "" {
		return errors.New("project already created")
	}
	s.projects[name] = name
	return nil
}

// Returns an array of all projects
func (s *StubProjectStore) GetAllProjects() []model.Project {
	var projects []model.Project

	for key, _ := range s.projects {
		projects = append(projects, model.Project{Name: key})
	}

	return projects
}

// Deletes a project from store
func (s *StubProjectStore) DeleteProject(name string) error {
	delete(s.projects, name)
	return nil
}

// Tests for route GET /projects/{name}
// todo update map to struct or array
func TestGetProject(t *testing.T) {
	store := StubProjectStore{
		map[string]string{
			"homework": "homework",
			"cleaning": "cleaning",
		},
	}

	// Uses the ProjectServer with our StubProjectStore
	server := api.NewProjectServer(&store)

	t.Run("returns project homework", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/projects/homework", nil)
		response := httptest.NewRecorder()

		server.Router.ServeHTTP(response, request)

		got := decodeProjectFromResponse(t, response.Body).Name
		want := "homework"

		assertResponseStatus(t, response.Code, 200)
		assertResponseBody(t, got, want)
	})

	t.Run("returns project cleaning", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/projects/cleaning", nil)
		response := httptest.NewRecorder()

		server.Router.ServeHTTP(response, request)

		got := decodeProjectFromResponse(t, response.Body).Name
		want := "cleaning"

		assertResponseStatus(t, response.Code, 200)
		assertResponseBody(t, got, want)
	})

	t.Run("return 404 on not existing projects", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/projects/laundry", nil)
		response := httptest.NewRecorder()

		server.Router.ServeHTTP(response, request)

		got := response.Code
		want := http.StatusNotFound

		assertResponseStatus(t, got, want)
	})
}

// Test for route POST /projects/
func TestPostProject(t *testing.T) {
	store := StubProjectStore{
		map[string]string{
			"homework": "homework",
			"cleaning": "cleaning",
		},
	}

	// Uses the ProjectServer with our StubProjectStore
	server := api.NewProjectServer(&store)

	t.Run("Creates new Project laundry", func(t *testing.T) {
		requestBody := makeNewPostProjectBody(t, "laundry")
		request, _ := http.NewRequest(http.MethodPost, "/projects/", requestBody)
		response := httptest.NewRecorder()

		server.Router.ServeHTTP(response, request)

		assertResponseStatus(t, response.Code, http.StatusCreated)
		assertProjectCreated(t, store, "laundry")
	})

	t.Run("Try to create a project that already exists", func(t *testing.T) {
		requestBody := makeNewPostProjectBody(t, "homework")
		request, _ := http.NewRequest(http.MethodPost, "/projects/", requestBody)
		response := httptest.NewRecorder()

		server.Router.ServeHTTP(response, request)

		assertResponseStatus(t, response.Code, http.StatusBadRequest)
	})
}

// Test for Route GET /projects
func TestGetAllProjects(t *testing.T) {
	store := StubProjectStore{
		map[string]string{
			"homework": "homework",
			"cleaning": "cleaning",
		},
	}

	// Uses the ProjectServer with our StubProjectStore
	server := api.NewProjectServer(&store)

	t.Run("Get all stored projects", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/projects", nil)
		response := httptest.NewRecorder()

		server.Router.ServeHTTP(response, request)

		gotProjects := decodeAllProjectsFromResponse(t, response.Body)

		var got [2]string
		for i, project := range gotProjects {
			got[i] = project.Name
		}

		want := [2]string{"homework", "cleaning"}

		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}

		assertResponseStatus(t, response.Code, 200)
	})
}

// Test for Route DEL /Project/{name}
func TestDeleteProject(t *testing.T) {
	store := StubProjectStore{
		map[string]string{
			"homework": "homework",
			"cleaning": "cleaning",
		},
	}

	// Uses the ProjectServer with our StubProjectStore
	server := api.NewProjectServer(&store)

	t.Run("Delete project homework", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodDelete, "/projects/homework", nil)
		response := httptest.NewRecorder()

		server.Router.ServeHTTP(response, request)

		assertProjectDeleted(t, store, "homework")
		assertResponseStatus(t, response.Code, 200)
	})

}

func assertResponseBody(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func assertResponseStatus(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("did not get correct status code, got %d, want %d", got, want)
	}
}

func assertProjectCreated(t testing.TB, store StubProjectStore, name string) {
	t.Helper()
	if project := store.projects[name]; project == "" {
		t.Fatalf("project was not created")
	}
}

func assertProjectDeleted(t testing.TB, store StubProjectStore, name string) {
	t.Helper()
	if project := store.projects[name]; project != "" {
		t.Fatalf("project was not deleted")
	}
}

func assertError(t testing.TB, context string, err error) {
	if err != nil {
		t.Fatalf("error - %v: %v", context, err)
	}
}

// Decodes the response body to a project struct
func decodeProjectFromResponse(t testing.TB, rdr io.Reader) model.Project {
	t.Helper()

	var project model.Project

	err := json.NewDecoder(rdr).Decode(&project)
	if err != nil {
		t.Errorf("problem parsing project, %v", err)
	}

	return project
}

// Decodes the response body to a project struct
func decodeAllProjectsFromResponse(t testing.TB, rdr io.Reader) []model.Project {
	t.Helper()

	var projects []model.Project

	err := json.NewDecoder(rdr).Decode(&projects)
	if err != nil {
		t.Errorf("problem parsing project, %v", err)
	}

	return projects
}

// makes a new json request body for POST /projects/
func makeNewPostProjectBody(t *testing.T, name string) *bytes.Buffer {
	requestBody, err := json.Marshal(map[string]string{
		"name": name,
	})

	if err != nil {
		t.Fatalf("Failed to make requestBody: %s", err)
	}

	return bytes.NewBuffer(requestBody)
}
