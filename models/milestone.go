package models

import (
	"errors"
	"time"
)

// Milestone status defined.
const (
	MilestoneStatusActive = 1 + iota
	MilestoneStatusDone
)

// Milestone represents a model for Milestone
type Milestone struct {
	BaseModel

	Name        string `sql:"not null;unique"`
	ProjectID   int
	DueDate     time.Time
	Status      int    //Active, Close
	Description string `sql:"size:1000"` // description of the milestone
}

// Validate check input value and return error map
//
// if second return value 'error' is not nil,
// you may consider invoke 500 internal error.
func (milestone *Milestone) Validate() (map[string]string, error) {
	errorMap := make(map[string]string)
	if !milestone.Required(milestone.Name) {
		errorMap["Name"] = "Name is required."
	}
	var err error

	if milestone.ProjectID == 0 {
		err = errors.New("Invalid project id")
	}

	if milestone.Status < MilestoneStatusActive || milestone.Status > MilestoneStatusDone {
		err = errors.New("Invalid milestone status")
	}
	return errorMap, err
}
