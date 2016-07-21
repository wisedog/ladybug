package models

import (
	"testing"
)

// TestPlanValidateEmptyTitle tests that TestPlan.Validate() function
// returns error map
func TestPlanValidateEmptyTitle(t *testing.T) {

	plan := TestPlan{}

	// check empty name
	errorMap, _ := plan.Validate()

	if len(errorMap) == 0 {
		t.Errorf(`models.TestPlan.Validate() error ({})`)
	}

}

// TestPlanValidateMaxLengthTitle tests that TestPlan.Validate() function
func TestPlanValidateMaxLengthTitle(t *testing.T) {
	// check non-empty, but invalid section id ( = 0)
	plan := TestPlan{Title: "1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890"}
	errorMap, _ := plan.Validate()

	if len(errorMap) == 0 {
		t.Errorf(`models.TestPlan.Validate() error ({Title : Too Long to print})`)
	}
}

// TestPlanValidateInvalidProjectID tests that TestPlan.Validate() function
// returns error on invalid project id
func TestPlanValidateInvalidProjectID(t *testing.T) {
	// check status
	plan := TestPlan{Title: "test", CreatorID: 1}
	_, err := plan.Validate()

	if err == nil {
		t.Errorf(`models.TestPlan.Validate() error ({Name: "test",CreatorID : 1})`)
	}

}

// TestPlanValidateInvalidProjectID tests that TestPlan.Validate() function
// returns error on invalid creator id
func TestPlanValidateInvalidCreatorID(t *testing.T) {
	// check status
	plan := TestPlan{Title: "test", ProjectID: 1}
	_, err := plan.Validate()

	if err == nil {
		t.Errorf(`models.TestPlan.Validate() error ({Name: "test", ProjectID : 1})`)
	}

}

// TestPlanValidateEmptyTestCases tests that TestPlan.Validate() function
// returns error on empty test case
func TestPlanValidateEmptyTestCases(t *testing.T) {
	plan := TestPlan{Title: "test", ProjectID: 1, CreatorID: 1}
	_, err := plan.Validate()

	if err == nil {
		t.Errorf(`models.TestPlan.Validate() error ({Title: "test", ProjectID: 1, CreatorID : 1})`)
	}

}

// TestPlanValidateNormal Condition tests that TestPlan.Validate() function on normal condition
func TestPlanValidateNormal(t *testing.T) {
	plan := TestPlan{Title: "test", ProjectID: 1, CreatorID: 1, ExecuteCases: "1,2,3"}
	errorMap, err := plan.Validate()

	if err != nil {
		t.Errorf(`models.TestPlan.Validate() error ({Title: "test", ProjectID: 1, CreatorID : 1, ExecuteCases : "1,2,3"})`)
	}
	if len(errorMap) > 0 {
		t.Errorf(`models.TestPlan.Validate() error ({Title: "test", ProjectID: 1, CreatorID : 1, ExecuteCases : "1,2,3"})`)
	}

}
