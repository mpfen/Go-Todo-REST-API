package model

import (
	"gorm.io/gorm"
)

type Project struct {
	gorm.Model `json:"id" gorm:"unique"`
	Name       string `json:"name" gorm:"unique"`
	Archived   bool   `json:"archived"`
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
