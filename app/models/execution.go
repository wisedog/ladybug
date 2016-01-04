package models

import (
	"github.com/revel/revel"
	"time"
)

type Execution struct {
	ID            int
	Status        int    // ready, in progress, N/A(if a tester is not assigned)
	ExecutionType int    //manual, automation
	Plan          TestPlan
	PlanId        int //TODO add an index
	Project       Project
	ProjectId     int
	Executor      User
	ExecutorId    int
	TargetBuild   BuildItem
	TargetBuildId int
	Message			string	// for store a reason of deny
	
	PassCaseNum		int `sql:"-"`
	FailCaseNum		int	`sql:"-"`
	Progress		int `sql:"-"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

func (execution *Execution) Validate(v *revel.Validation) {

}
