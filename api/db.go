package api

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	model "github.com/mpfen/Go-Todo-REST-API/api/model"
)

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
func NewDatabaseConnection() *Database {
	db, err := gorm.Open(sqlite.Open("database.db"), &gorm.Config{})

	if err != nil {
		log.Fatalf("Can not open Database %s", err)
	}

	db = model.DbMigrate(db)

	return &Database{DB: db}
}
