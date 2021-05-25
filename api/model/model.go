package model

import (
	"time"

	"gorm.io/gorm"
)

type Project struct {
	gorm.Model `json:"id" gorm:"unique"`
	Name       string `json:"name" gorm:"unique"`
	Archived   bool   `json:"archived"`
	Tasks      []Task `gorm:"ForeignKey:ProjectID" json:"tasks"`
}

func DbMigrate(db *gorm.DB) *gorm.DB {
	db.AutoMigrate(&Project{})
	return db
}

func (p *Project) ArchiveProject() {
	p.Archived = true
}

func (p *Project) UnArchiveProject() {
	p.Archived = false
}

type Task struct {
	gorm.Model
	Name      string     `json:"name"`
	Priority  string     `gorm:"type:ENUM('0', '1', '2', '3');default:'0'" json:"priority"`
	Deadline  *time.Time `gorm:"default:null" json:"deadline"`
	Done      bool       `json:"done"`
	ProjectID uint       `json:"project_id"`
}
