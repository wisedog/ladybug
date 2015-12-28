package models

import (
	"github.com/revel/revel"
	"time"
)

type TestPlan struct {
	ID            int
	Title         string
	Description   string
	Status        int		//Active , In Review 
	ExecutionType int
	Project       Project
	Project_id    int
	Creator       User
	CreatorId     int
	Executor      User
	ExecutorId    int
	ExecCaseNum		int
	ExecuteCases  string	// ',' joined string like 1,2,3.... 
	TargetBuild   BuildItem
	TargetBuildId int

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

func (testplan *TestPlan) Validate(v *revel.Validation) {
	v.Check(testplan.Title,
		revel.Required{},
		revel.MaxSize{400},
	)

	v.Check(testplan.ExecuteCases,
		revel.Required{},
	)

}
