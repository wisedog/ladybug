package models

import (
	"testing"
)

func TestMilestoneValidation(t *testing.T) {

	milestone := Milestone{Name: ""}

	// check empty name
	errorMap, _ := milestone.Validate()

	if len(errorMap) == 0 {
		t.Errorf(`models.Milestone.Validate() error ({Name:""})`)
	}
}
