package models

import (
	"testing"
)

// TestMilestoneValidateEmptyName tests empty milestone name
func TestMilestoneValidateEmptyName(t *testing.T) {

	milestone := Milestone{Name: ""}

	// check empty name
	errorMap, _ := milestone.Validate()

	if len(errorMap) == 0 {
		t.Errorf(`models.Milestone.Validate() error ({Name:""})`)
	}
}
