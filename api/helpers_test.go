package api_test

import (
	"errors"

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

// Update a project in store
func (s *StubProjectStore) UpdateProject(oldName, newName string) error {
	delete(s.projects, oldName)
	s.projects[newName] = newName
	return nil
}
