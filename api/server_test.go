package api_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	api "github.com/mpfen/Go-Todo-REST-API/api"
)

// DB stub for testing
type StubProjectStore struct {
	projects map[string]string
}

func (s *StubProjectStore) GetProjectInfo(name string) string {
	return name
}

// GET /projects/{name}
func TestGetProject(t *testing.T) {
	store := StubProjectStore{
		map[string]string{
			"homework": "homework",
			"cleaning": "cleaning",
		},
	}

	server := &api.ProjectServer{&store}

	t.Run("returns project homework", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/projects/homework", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := response.Body.String()
		want := "homework"

		assertResponseStatus(t, response.Code, 200)
		assertResponseBody(t, got, want)
	})

	t.Run("returns project cleaning", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/projects/cleaning", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := response.Body.String()
		want := "cleaning"

		assertResponseStatus(t, response.Code, 200)
		assertResponseBody(t, got, want)
	})

	t.Run("return 404 on missing projects", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/projects/laundry", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

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
