package models

type Project struct {
	BaseModel

	Name        string `sql:"not null;unique"`
	Status      int    // Activation, Deactivation, closed
	Description string `sql:"size:4000"` // description of the issue
	Users       []User `gorm:"many2many:user_project;"`
	Prefix      string
}
/*
func (project *Project) Validate(v *revel.Validation) {
	v.Check(project.Name,
		revel.Required{},
		revel.MaxSize{48},
	)
}
*/