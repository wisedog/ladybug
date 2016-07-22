package models

import (
	"testing"
)

// TestExecutionValidateNormalCondition tests normal condition of Execution
func TestExecutionValidateNormalCondition(t *testing.T) {

	exec := Execution{Status: ExecStatusReady, ProjectID: 1, PlanID: 2}

	// check empty name
	err := exec.Validate()

	if err != nil {
		t.Errorf(`models.Execution.Validate() error ({Status: ExecStatusReady, ProjectID: 1, PlanID: 2})`)
	}
}

// TestExecutionValidateAbsenseProjectID tests ProjectID is not set
func TestExecutionValidateAbsenseProjectID(t *testing.T) {

	exec := Execution{Status: ExecStatusReady, PlanID: 2}

	// check empty name
	err := exec.Validate()

	if err == nil {
		t.Errorf(`models.Execution.Validate() error ({Status: ExecStatusReady, PlanID: 2})`)
	}
}

// TestExecutionValidateAbsensePlanID tests PlanID is not set
func TestExecutionValidateAbsensePlanID(t *testing.T) {
	exec := Execution{Status: ExecStatusReady, ProjectID: 2}

	// check empty name
	err := exec.Validate()

	if err == nil {
		t.Errorf(`models.Execution.Validate() error ({Status: ExecStatusReady, ProjectID: 2})`)
	}
}
