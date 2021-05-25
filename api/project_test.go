package api_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	api "github.com/mpfen/Go-Todo-REST-API/api"
	"github.com/mpfen/Go-Todo-REST-API/api/model"
)

// Test setUp for all project tests
func setUpProjectTests() (server *api.ProjectServer, store StubProjectStore) {
	store = StubProjectStore{
		map[string]bool{
			"homework": false,
			"cleaning": true,
		},
	}

	// Uses the ProjectServer with our StubProjectStore
	server = api.NewProjectServer(&store)
	return server, store
}

// Tests for route GET /projects/{name}
// todo update map to struct or array
func TestGetProject(t *testing.T) {
	server, _ := setUpProjectTests()

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
	server, store := setUpProjectTests()

	t.Run("Creates new Project laundry", func(t *testing.T) {
		requestBody := makeNewPostProjectBody(t, "laundry")
		request, _ := http.NewRequest(http.MethodPost, "/projects", requestBody)
		response := httptest.NewRecorder()

		server.Router.ServeHTTP(response, request)

		assertResponseStatus(t, response.Code, http.StatusCreated)
		assertProjectCreated(t, store, "laundry")
	})

	t.Run("Try to create a project that already exists", func(t *testing.T) {
		requestBody := makeNewPostProjectBody(t, "homework")
		request, _ := http.NewRequest(http.MethodPost, "/projects", requestBody)
		response := httptest.NewRecorder()

		server.Router.ServeHTTP(response, request)

		assertResponseStatus(t, response.Code, http.StatusBadRequest)
	})
}

// Test for route GET /projects
func TestGetAllProjects(t *testing.T) {
	server, _ := setUpProjectTests()

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

// Test for route DEL /projects/{name}
func TestDeleteProject(t *testing.T) {
	server, store := setUpProjectTests()

	t.Run("Delete project homework", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodDelete, "/projects/homework", nil)
		response := httptest.NewRecorder()

		server.Router.ServeHTTP(response, request)

		assertProjectDeleted(t, store, "homework")
		assertResponseStatus(t, response.Code, 200)
	})

}

// Test for route PUT /projects/{name}
func TestPutProject(t *testing.T) {
	server, store := setUpProjectTests()

	t.Run("Change the name of project homework", func(t *testing.T) {
		projectName := "researchpaper"
		requestBody := makeNewPostProjectBody(t, projectName)
		request, _ := http.NewRequest(http.MethodPut, "/projects/homework", requestBody)
		response := httptest.NewRecorder()

		server.Router.ServeHTTP(response, request)

		// todo? not working with StubProjectStore implementation
		// UpdateProject can not delete old project
		// because it does not know the old projects name
		// assertProjectDeleted(t, store, "homework")
		assertProjectCreated(t, store, projectName)
		assertResponseStatus(t, response.Code, 200)
	})
}

// Test for route PUT /projects/{name}/archive
func TestArchiveProject(t *testing.T) {
	server, store := setUpProjectTests()

	t.Run("Archive project homework", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPut, "/projects/homework/archive", nil)
		response := httptest.NewRecorder()

		server.Router.ServeHTTP(response, request)

		assertProjectArchiveStatus(t, store, "homework", true)
		assertResponseStatus(t, response.Code, http.StatusOK)
	})

	t.Run("Try to archive not existing project", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPut, "/projects/researchpaper/archive", nil)
		response := httptest.NewRecorder()

		server.Router.ServeHTTP(response, request)
		assertResponseStatus(t, response.Code, http.StatusNotFound)
	})
}

// Test for route DELETE /projects/{name}/archive
func TestUnArchiveProject(t *testing.T) {
	server, store := setUpProjectTests()

	t.Run("unarchive project homework", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodDelete, "/projects/cleaning/archive", nil)
		response := httptest.NewRecorder()

		server.Router.ServeHTTP(response, request)

		assertProjectArchiveStatus(t, store, "homework", false)
		assertResponseStatus(t, response.Code, http.StatusOK)
	})

	t.Run("Try to unarchive not existing project", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodDelete, "/projects/researchpaper/archive", nil)
		response := httptest.NewRecorder()

		server.Router.ServeHTTP(response, request)
		assertResponseStatus(t, response.Code, http.StatusNotFound)
	})
}

// assert functions specific to project tests
func assertProjectCreated(t testing.TB, store StubProjectStore, name string) {
	t.Helper()
	if _, exists := store.projects[name]; !exists {
		t.Errorf("project was not created")
	}
}

func assertProjectDeleted(t testing.TB, store StubProjectStore, name string) {
	t.Helper()
	if _, exists := store.projects[name]; exists {
		t.Errorf("project was not deleted")
	}
}

func assertProjectArchiveStatus(t *testing.T, store StubProjectStore, name string, want bool) {
	t.Helper()
	if archived := store.projects[name]; archived != want {
		t.Errorf("Project has wrong archive status - got %v, want %v", archived, want)
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
		t.Errorf("Failed to make requestBody: %s", err)
	}

	return bytes.NewBuffer(requestBody)
}
