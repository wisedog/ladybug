package models

import (
	"github.com/revel/revel"
	"time"
)

type Execution struct {
	ID            int
	Title         string //???
	Status        int    // ready, in progress, N/A(if a tester is not assigned)
	ExecutionType int    //manual, automation
	Plan          TestPlan
	PlanId        int //TODO add an index
	Project       Project
	ProjectId     int
	Executor      User
	ExecutorId    int
	ExecuteCases  string	// , separated string
	ExecResult    string	// , separated string
	TargetBuild   Build
	TargetBuildId int

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

func (execution *Execution) Validate(v *revel.Validation) {

}
