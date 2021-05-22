package api_test

import (
	"encoding/json"
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

// Tests for route GET /projects/{name}
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

	t.Run("return 404 on missing projects", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/projects/laundry", nil)
		response := httptest.NewRecorder()

		server.Router.ServeHTTP(response, request)

		got := response.Code
		want := http.StatusNotFound

		assertResponseStatus(t, got, want)
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
