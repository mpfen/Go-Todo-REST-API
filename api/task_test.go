package api_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	api "github.com/mpfen/Go-Todo-REST-API/api"
)

func setupTaskTests() (server *api.TodoStore, store StubTodoStore) {
	time := time.Now()
	store = StubTodoStore{
		map[string]bool{
			"homework": false,
			"cleaning": true,
		},
		[]stubTask{{Name: "math",
			Priority:  "1",
			Deadline:  &time,
			Done:      false,
			ProjectID: "homework",
		}, {Name: "kitchen",
			Priority:  "1",
			Deadline:  &time,
			Done:      false,
			ProjectID: "cleaning",
		},
		},
	}
	// Uses the TodoStore with our StubTodoStore
	server = api.NewTodoStore(&store)
	return server, store
}

// Tests for route GET /projects/{projectname}/tasks/{taskname}
func TestGetTask(t *testing.T) {
	server, _ := setupTaskTests()

	t.Run("Get task Math from project homework", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/projects/homework/task/math", nil)
		response := httptest.NewRecorder()

		server.Router.ServeHTTP(response, request)

		want := "math"
		got := decodeTaskFromResponse(t, response.Body)

		assertResponseBody(t, got.Name, want)
		assertResponseStatus(t, response.Code, http.StatusOK)
	})

	t.Run("Try to get task from wrong project", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/projects/homework/task/kitchen", nil)
		response := httptest.NewRecorder()

		server.Router.ServeHTTP(response, request)

		want := ""
		got := decodeTaskFromResponse(t, response.Body)

		assertResponseBody(t, got.Name, want)
		assertResponseStatus(t, response.Code, http.StatusNotFound)
	})

	t.Run("Try to get not existing Task", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/projects/homework/task/biology", nil)
		response := httptest.NewRecorder()

		server.Router.ServeHTTP(response, request)

		want := ""
		got := decodeTaskFromResponse(t, response.Body)

		assertResponseBody(t, got.Name, want)
		assertResponseStatus(t, response.Code, http.StatusNotFound)
	})

}

// Decodes the response body to a task struct
func decodeTaskFromResponse(t testing.TB, rdr io.Reader) stubTask {
	t.Helper()

	var task stubTask

	err := json.NewDecoder(rdr).Decode(&task)
	if err != nil {
		t.Errorf("problem parsing task, %v", err)
	}

	return task
}
