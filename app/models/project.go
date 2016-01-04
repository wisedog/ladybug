package models

import (
	"github.com/revel/revel"
	"time"
)

type Project struct {
	ID          int
	Name        string `sql:"not null;unique"`
	Status      int    // Activation, Deactivation, closed
	Description string `sql:"size:4000"` // description of the issue
	Users       []User `gorm:"many2many:user_project;"`
	Prefix      string

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

func (project *Project) Validate(v *revel.Validation) {
	v.Check(project.Name,
		revel.Required{},
		revel.MaxSize{48},
	)
}
