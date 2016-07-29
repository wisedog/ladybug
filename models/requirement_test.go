package models

import (
	"testing"
)

// TestRequirementValidateErrorMap tests that Requirement.Validate() function
// returns error map
func TestRequirementValidateErrorMap(t *testing.T) {

	req := Requirement{Title: ""}

	// check empty name
	errorMap, _ := req.Validate()

	if len(errorMap) == 0 {
		t.Errorf(`models.Requirement.Validate() error ({Title:""})`)
	}

}

// TestRequirementValidateSectionID tests that Requirement.Validate() function
// returns error on invalid section ID
func TestRequirementValidateSectionID(t *testing.T) {
	// check non-empty, but invalid section id ( = 0)
	req := Requirement{Title: "test", SectionID: 0}
	_, err := req.Validate()

	if err == nil {
		t.Errorf(`models.Requirement.Validate() error ({NaTitleme : "test", SectionID : 0})`)
	}
}

// TestRequirementValidateStatus tests that Requirement.Validate() function
// returns error on invalid status
func TestRequirementValidateStatus(t *testing.T) {
	// check status
	req := Requirement{Title: "test", SectionID: 1, Status: 10}
	_, err := req.Validate()

	if err == nil {
		t.Errorf(`models.Requirement.Validate() error ({TitleName: "test", SectionID: 1, Status : 6})`)
	}

}
