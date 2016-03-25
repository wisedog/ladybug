package models


// TestResult is a model for test result.
type TestResult struct {
  BaseModel

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
}
