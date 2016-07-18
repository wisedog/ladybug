package models

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

/*
func (testplan *TestPlan) Validate(v *revel.Validation) {
	v.Check(testplan.Title,
		revel.Required{},
		revel.MaxSize{400},
	)

	v.Check(testplan.ExecuteCases,
		revel.Required{},
	)

}
*/
