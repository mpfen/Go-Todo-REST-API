package api_test

import (
	"bytes"
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
			"school":   false,
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
		}, {Name: "physics",
			Priority:  "1",
			Deadline:  &time,
			Done:      false,
			ProjectID: "homework",
		}},
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

// Tests for route POST /projects/{projectName}/tasks/{taskName}
func TestPostTask(t *testing.T) {
	server, store := setupTaskTests()

	t.Run("Create a new task for project homework", func(t *testing.T) {
		requestBody := makeNewPostTaskBody(t, "biology", "homework")
		request, _ := http.NewRequest(http.MethodPost, "/projects/homework/task", requestBody)
		response := httptest.NewRecorder()

		server.Router.ServeHTTP(response, request)

		assertResponseStatus(t, response.Code, http.StatusCreated)
		assertTaskCreated(t, store, "biology")
	})

	t.Run("Try to create an already existing task", func(t *testing.T) {
		requestBody := makeNewPostTaskBody(t, "biology", "homework")
		request, _ := http.NewRequest(http.MethodPost, "/projects/homework/task", requestBody)
		response := httptest.NewRecorder()

		server.Router.ServeHTTP(response, request)

		assertResponseStatus(t, response.Code, http.StatusBadRequest)
	})

	t.Run("Try to create a task for a not existing project", func(t *testing.T) {
		requestBody := makeNewPostTaskBody(t, "biology", "school")
		request, _ := http.NewRequest(http.MethodPost, "/projects/school/task", requestBody)
		response := httptest.NewRecorder()

		server.Router.ServeHTTP(response, request)

		assertResponseStatus(t, response.Code, http.StatusBadRequest)
	})
}

// Test for getting all task of a project GET /projects/{projectName}/tasks
func TestGetAllTasksOfAProject(t *testing.T) {
	server, store := setupTaskTests()

	t.Run("Get all task from project homework", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/projects/homework/task", nil)
		response := httptest.NewRecorder()

		server.Router.ServeHTTP(response, request)

		assertResponseStatus(t, response.Code, http.StatusOK)
		want := []stubTask{}
		want = append(want, store.Tasks[0], store.Tasks[2])
		got := decodeMultipleTaskFromResponse(t, response.Body)
		t.Log(want)
		t.Log(got)

		assertTaskList(t, got, want)
	})

	t.Run("Try to get task from a project without tasks", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/projects/school/task", nil)
		response := httptest.NewRecorder()

		server.Router.ServeHTTP(response, request)

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

// Decodes the response body to a task struct
func decodeMultipleTaskFromResponse(t testing.TB, rdr io.Reader) []stubTask {
	t.Helper()

	var tasks []stubTask

	err := json.NewDecoder(rdr).Decode(&tasks)
	if err != nil {
		t.Errorf("problem parsing task, %v", err)
	}

	return tasks
}

func makeNewPostTaskBody(t *testing.T, taskName, projectName string) *bytes.Buffer {
	requestBody, err := json.Marshal(map[string]string{
		"name":      taskName,
		"ProjectID": projectName,
	})

	if err != nil {
		t.Errorf("Failed to make requestBody: %s", err)
	}

	return bytes.NewBuffer(requestBody)

}

// assert functions specific for tasks
func assertTaskCreated(t testing.TB, store StubTodoStore, name string) {
	t.Helper()
	for _, task := range store.Tasks {
		if task.Name == name {
			return
		}
	}
	t.Errorf("Task %v was not created", name)
}

func assertTaskList(t testing.TB, got, want []stubTask) {
	for i, task := range got {
		if task.Name != want[i].Name {
			t.Errorf("List of Tasks not matching: got %v want %v", task.Name, want[i].Name)
		}
	}
}
