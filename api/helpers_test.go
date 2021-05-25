package api_test

import (
	"errors"
	"testing"

	"github.com/mpfen/Go-Todo-REST-API/api/model"
)

// DB store stub for testing
type StubProjectStore struct {
	projects map[string]bool
}

// Creates a makeshift project struct to comply with ProjectStore interface
func (s *StubProjectStore) GetProject(name string) model.Project {
	project := model.Project{}
	if _, exists := s.projects[name]; exists {
		project.Name = name
		return project
	} else {
		return project
	}
}

// Creates a new project
func (s *StubProjectStore) PostProject(name string) error {
	if _, exists := s.projects[name]; exists {
		return errors.New("project already created")
	} else {
		s.projects[name] = false
		return nil
	}
}

// Returns an array of all projects
func (s *StubProjectStore) GetAllProjects() []model.Project {
	var projects []model.Project

	for key := range s.projects {
		projects = append(projects, model.Project{Name: key})
	}

	return projects
}

// Deletes a project from store
func (s *StubProjectStore) DeleteProject(name string) error {
	delete(s.projects, name)
	return nil
}

// "Updates" a project in store
func (s *StubProjectStore) UpdateProject(project model.Project) error {
	s.projects[project.Name] = project.Archived
	return nil
}

// common assert functions
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

func assertError(t testing.TB, context string, err error) {
	if err != nil {
		t.Errorf("error - %v: %v", context, err)
	}
}
