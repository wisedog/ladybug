package models

// TestCaseResult is a model for test result of each test case
// A test execution consists of one or more test cases.
// Test Case(s) <--> TestResult
// Test Execution <--> TestExecution
// Test Case(s) ---> Test Execution
type TestCaseResult struct {
	BaseModel

	// Status represents pass or fail of this test execution
	Status bool

	Plan        TestPlan
	TestPlanID  int
	Exec        Execution
	ExecID      int
	Case        TestCase
	TestCaseID  int
	TestCaseVer int

	// Actual describes actual result of this test case.
	// It will be if the Status is failed
	Actual string `sql:"size:800"`
}
