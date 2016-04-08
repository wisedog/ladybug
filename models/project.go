package models


// Project Status enumulations
const (
	ProjectStatusActivate = 1 + iota
	ProjectStatusInactivate
	ProjectStatusClose
	)

// Project represents a model for Project
type Project struct {
	BaseModel

	Name        string `sql:"not null;unique"`
	Status      int    // Activation, Deactivation, closed
	Description string `sql:"size:4000"` // description of the issue
	Users       []User `gorm:"many2many:user_project;"`
	Prefix      string `sql:"size:12"`

	// belows are used for form input, not store in database
	userOption 	string `sql:"-"`
	//publicFlag 	
}

// Validate check input value and return error map
func (project *Project) Validate() map[string]string {
	errorMap := make(map[string]string)
	if !project.Required(project.Name){
		errorMap["Name"] = "Name is required."
	}

	if !project.Required(project.Prefix){
		errorMap["Prefix"] = "Prefix is required."
	}

	if !project.MaxSize(project.Prefix, 12){
		errorMap["Prefix"] += " Size of Prefix exceeds 12 characters."
	}
	return errorMap
}