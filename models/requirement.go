package models

import "errors"

// Status of requirements
const (
	ReqStatusActivate = 1 + iota
	ReqStatusInactivate
	ReqStatusDraft
)

// Requirement Model
type Requirement struct {
	BaseModel

	Title       string
	Description string
	// Status of this requirement. Status may be one of below items
	// Draft, In Review, Not Testable, Deprecated ....
	Status   int
	Priority int

	// Each requirements are belongs to section
	SectionID int

	// ProjectID represents Project ID of this requirement. Each Requirement are belongs to project
	ProjectID int

	// Version or history of this Requirement
	Version int

	// ReqType is type of this requirement
	// Requirement type may be one of below
	// Use Case, Information, Feature, User Interface, Non Functional, Constraint, System Function...
	ReqType int

	// RelatedTestCase stores relationship between this requirement and related testcases
	// The relationship has many to many
	RelatedTestCases []TestCase `gorm:"many2many:testcases_reqs;"`
}

// Validate check input value and return error map
//
// if second return value 'error' is not nil,
// you may consider invoke 500 internal error.
func (req *Requirement) Validate() (map[string]string, error) {
	errorMap := make(map[string]string)
	if !req.Required(req.Title) {
		errorMap["Title"] = "Title is required."
	}
	var err error

	if req.SectionID == 0 {
		err = errors.New("Invalid section id")
	}

	if req.Status < ReqStatusActivate || req.Status > ReqStatusInactivate {
		err = errors.New("Invalid requirement status")
	}
	return errorMap, err
}
