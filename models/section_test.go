package models

import (
	"testing"
)

// TestSectionNormal tests that Section.Validate() function
// returns error map
func TestSectionNormal(t *testing.T) {

	section := Section{Title: "Functional Requirements", ProjectID: 1}
	errorMap, err := section.Validate()

	if len(errorMap) > 0 {
		t.Errorf(`models.Section.Validate() error ({Title: "Functional Requirements", ProjectID: 1})`)
	}

	if err != nil {
		t.Errorf(`unexpected error on models.Section.Validate()`)
	}

}

// TestSectionEmptyTitle tests that Section.Validate() function
// returns error map
func TestSectionEmptyTitle(t *testing.T) {

	section := Section{Title: "", ProjectID: 1}
	errorMap, err := section.Validate()

	if len(errorMap) == 0 {
		t.Errorf(`models.Section.Validate() error ({Title: "", ProjectID: 1})`)
	}

	if err != nil {
		t.Errorf(`unexpected error on models.Section.Validate()`)
	}

}

// TestSectionInvalidProjectID tests that Section.Validate() function
// returns error map
func TestSectionInvalidProjectID(t *testing.T) {

	section := Section{Title: "System"}
	errorMap, err := section.Validate()

	if len(errorMap) > 0 {
		t.Errorf(`models.Section.Validate() error ({Title: "System"})`)
	}

	if err == nil {
		t.Errorf(`missing error on models.Section.Validate()`)
	}

}
