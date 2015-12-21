package models

import (
	"github.com/revel/revel"
	"time"
)

/*
  A model for test result.
*/


	// make table : testresult
	// ID
	// actual result
	// note
	// status : pass/fail
	// testcase id
	// testcase version
	// testplan id
	// testexecution id 
	// created update .... 
type TestResult struct {
	ID            int
	Status          bool
	Note            string  `sql:"size:800"`
	Actual          string  `sql:"size:800"`
	Case            TestCase
	TestCaseId      int
	TestCaseVer     int
	Plan            TestPlan
	TestPlanId      int
	Exec            Execution
	ExecId          int

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

func (testresult *TestResult) Validate(v *revel.Validation) {
	
}
