package models

import "errors"

// TestPlan represents a plan of test
type TestPlan struct {
	BaseModel
	Title         string
	Description   string
	Status        int //Active , In Review
	ExecutionType int
	Project       Project
	ProjectID     int
	Creator       User
	CreatorID     int
	Executor      User
	ExecutorID    int
	ExecCaseNum   int
	ExecuteCases  string // ',' joined string like 1,2,3....
	TargetBuild   BuildItem
	TargetBuildID int
}

// Validate check input value and return error map
//
// if second return value 'error' is not nil,
// you may consider invoke 500 internal error.
func (plan *TestPlan) Validate() (map[string]string, error) {
	errorMap := make(map[string]string)
	if !plan.Required(plan.Title) {
		errorMap["Title"] = "Title is required."
	}

	if !plan.MaxSize(plan.Title, 400) {
		errorMap["Title"] = "Too long title. Max length is 400"
	}

	var err error
	if plan.ProjectID == 0 {
		err = errors.New("Invalid project id")
	}

	// creator should be
	if plan.CreatorID == 0 {
		err = errors.New("Invalid project id")
	}

	// test cases should be
	if !plan.Required(plan.ExecuteCases) {
		err = errors.New("No test cases to be executed")
	}

	// but executor, targetbuild is not essential

	// TODO status
	// TODO manual or ...
	return errorMap, err
}
