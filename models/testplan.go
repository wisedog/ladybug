package models

import (
	"github.com/revel/revel"
)

type TestPlan struct {
	BaseModel
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
