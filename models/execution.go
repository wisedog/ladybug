package models

import "errors"

//status for test execution
const (
	ExecStatusReady = 1 + iota
	ExecStatusInProgress
	ExecStatusNotAvailable
	ExecStatusDone
	ExecStatusDeny
	ExecStatusPass
	ExecStatusFail
)

// Execution represents an execution test
type Execution struct {
	BaseModel

	Status        int // ready, in progress, N/A(if a tester is not assigned), Done
	ExecutionType int // manual, automation
	Plan          TestPlan
	PlanID        int // TODO add an index
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

// Validate check input value and return error map
//
// if second return value 'error' is not nil,
// you may consider invoke 500 internal error.
func (testexec *Execution) Validate() error {
	var err error
	if testexec.ProjectID == 0 {
		err = errors.New("Invalid project id")
	}

	if testexec.PlanID == 0 {
		err = errors.New("Invalid test plan id")
	}

	if testexec.Status < ExecStatusReady || testexec.Status > ExecStatusFail {
		err = errors.New("Invalid status")
	}
	return err
}
