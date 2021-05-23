package api_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/mpfen/Go-Todo-REST-API/api/store"
)

const testdbfile = "testdb.db"

// create a testdb file
func createTestDB(t *testing.T) {
	testDB := []byte("")
	err := ioutil.WriteFile(testdbfile, testDB, 0644)

	if err != nil {
		t.Fatalf("could not create test database %v", err)
	}
}

// delete the testdb file
func deleteTestDB(t *testing.T) {
	err := os.Remove(testdbfile)

	if err != nil {
		t.Fatalf("Could not remove test databse %v", err)
	}
}

// Integration tests for database
// uses own database file
func TestDatabase(t *testing.T) {
	createTestDB(t)
	defer deleteTestDB(t)

	db := store.NewDatabaseConnection(testdbfile)

	t.Run("create a new project in database", func(t *testing.T) {
		want := "TestDatabase"
		err := db.PostProject(want)

		assertError(t, "Create new project in db", err)

		got := db.GetProject(want).Name

		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}

	})

	t.Run("try to create an already existing project", func(t *testing.T) {
		err := db.PostProject("TestDatabase")

		if err == nil {
			t.Fatalf("Project should not have been created")
		}
	})

}
