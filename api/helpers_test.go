package api_test

import (
	"errors"
	"testing"
	"time"

	"github.com/mpfen/Go-Todo-REST-API/api/model"
)

type stubTask struct {
	Name      string
	Priority  string
	Deadline  *time.Time
	Done      bool
	ProjectID string
}

// DB store stub for testing
type StubTodoStore struct {
	Projects map[string]bool
	Tasks    []stubTask
}

// Creates a makeshift project struct to comply with TodoStore interface
func (s *StubTodoStore) GetProject(name string) model.Project {
	project := model.Project{}
	if _, exists := s.Projects[name]; exists {
		project.Name = name
		return project
	} else {
		return project
	}
}

// Creates a new project
func (s *StubTodoStore) PostProject(name string) error {
	if _, exists := s.Projects[name]; exists {
		return errors.New("project already created")
	} else {
		s.Projects[name] = false
		return nil
	}
}

// Returns an array of all projects
func (s *StubTodoStore) GetAllProjects() []model.Project {
	var projects []model.Project

	for key := range s.Projects {
		projects = append(projects, model.Project{Name: key})
	}

	return projects
}

// Deletes a project from store
func (s *StubTodoStore) DeleteProject(name string) error {
	delete(s.Projects, name)
	return nil
}

// "Updates" a project in store
func (s *StubTodoStore) UpdateProject(project model.Project) error {
	s.Projects[project.Name] = project.Archived
	return nil
}

// Gets Task from store
func (s *StubTodoStore) GetTask(projectID, taskName string) model.Task {
	for _, t := range s.Tasks {
		if t.Name == taskName && t.ProjectID == projectID {
			return wrapStubTask(taskName)
		}
	}
	return model.Task{}
}

// Create task in store
func (s *StubTodoStore) PostTask(task model.Task) error {
	for _, t := range s.Tasks {
		if t.Name == task.Name {
			return errors.New("Task already exists")
		}
	}

	// todo - append to []Tasks does not work
	// newTask := stubTask{Name: task.Name}
	// s.Tasks = append(s.Tasks, newTask)
	s.Tasks[1].Name = "biology"

	return nil
}

// Return all tasks of a project
func (s *StubTodoStore) GetAllProjectTasks(project model.Project) []model.Task {
	tasks := []model.Task{}

	for _, t := range s.Tasks {
		if t.ProjectID == project.Name {
			tasks = append(tasks, wrapStubTask(t.Name))
		}
	}
	return tasks
}

// Delete a tasks from the store
func (s *StubTodoStore) DeleteTask(task model.Task) error {
	for i, storeTask := range s.Tasks {
		if storeTask.Name == task.Name {
			s.Tasks[i].Name = ""
			return nil
		}
	}
	return nil
}

// to comply with interface
func wrapStubTask(taskName string) model.Task {
	modelTask := model.Task{}
	modelTask.Name = taskName
	return modelTask
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
