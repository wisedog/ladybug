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

/*func (hotel *Hotel) Validate(v *revel.Validation) {
	v.Check(hotel.Name,
		revel.Required{},
		revel.MaxSize{50},
	)

	v.MaxSize(hotel.Address, 100)

	v.Check(hotel.City,
		revel.Required{},
		revel.MaxSize{40},
	)

	v.Check(hotel.State,
		revel.Required{},
		revel.MaxSize{6},
		revel.MinSize{2},
	)

	v.Check(hotel.Zip,
		revel.Required{},
		revel.MaxSize{6},
		revel.MinSize{5},
	)

	v.Check(hotel.Country,
		revel.Required{},
		revel.MaxSize{40},
		revel.MinSize{2},
	)
}*/
