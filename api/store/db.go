package store

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	model "github.com/mpfen/Go-Todo-REST-API/api/model"
)

// ProjectStore interface for testing
// Tests use own implementation with
// StubStore instead of a real database
type ProjectStore interface {
	GetProject(name string) model.Project
	PostProject(name string) error
}

type Database struct {
	DB *gorm.DB
}

// Gets project by name
func (d *Database) GetProject(name string) model.Project {
	project := model.Project{}
	err := d.DB.Find(&project, "Name = ?", name).Error

	if err != nil {
		return model.Project{}
	}

	return project
}

// Creates a new project
func (d *Database) PostProject(name string) error {
	project := model.Project{}
	project.Name = name
	project.Archived = false

	err := d.DB.Create(&project).Error

	return err
}

// creates database struct and runs automigrate
func NewDatabaseConnection(name string) *Database {
	db, err := gorm.Open(sqlite.Open(name), &gorm.Config{})

	if err != nil {
		log.Fatalf("Can not open Database %s", err)
	}

	db = model.DbMigrate(db)

	return &Database{DB: db}
}
