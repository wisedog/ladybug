package models

// Execution represents an execution test
type Execution struct {
	BaseModel

	Status        int // ready, in progress, N/A(if a tester is not assigned), Done
	ExecutionType int //manual, automation
	Plan          TestPlan
	PlanID        int //TODO add an index
	Project       Project
	ProjectID     int
	Executor      User
	ExecutorID    int
	TargetBuild   BuildItem
	TargetBuildID int
	//Message stores a reason of deny or comment of test execution
	Message string

	// Test result : Pass Fail if Status is Done
	Result bool

	TotalCaseNum int
	PassCaseNum  int
	FailCaseNum  int
	Progress     int `sql:"-"`
}
