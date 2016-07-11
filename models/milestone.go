package models

import (
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
func (milestone *Milestone) Validate() map[string]string {
	errorMap := make(map[string]string)
	if !milestone.Required(milestone.Name){
		errorMap["Name"] = "Name is required."
	}
	return errorMap
}