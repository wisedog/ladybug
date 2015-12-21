package models

import (
	"github.com/revel/revel"
	"time"
)

type Issue struct {
	ID          int
	ProjectId   int    //Project id which belongs to this issue
	Status      int    // Status. Open, Close, In Progress, In QA Review, Waiting for QA, Resolved, Todo
	Summary     string `sql:"size:400"`  // summary of the issue
	Description string `sql:"size:4000"` // description of the issue
	Assignee    User
	AssigneeId  int // TODO add foreign key
	Reporter    User
	ReporterId  int // TODO add foreign key
	Resolution  int // fixed or unresolved
	Created     time.Time
	Updated     time.Time
	Due         time.Time
	Priority    int // Critical, Urgent, High, Low, misc
	Progress    int // progress. 0~100
	Roadmap     int // TODO Change to Roadmap object when a Roadmap table is created
	Started     time.Time
	Ended       time.Time
	Estimated   time.Time
	Attachment  int // TODO change to ....

}

func (issue *Issue) Validate(v *revel.Validation) {
	v.Check(issue.Summary,
		revel.Required{},
		revel.MaxSize{400},
	)

	v.Check(issue.Description,
		revel.Required{},
		revel.MaxSize{4000},
	)

}

/*func (hotel *Hotel) Validate(v *revel.Validation) {
	v.Check(hotel.Name,
		revel.Required{},
		revel.MaxSize{50},
	)

	v.MaxSize(hotel.Address, 100)

	v.Check(hotel.City,
		revel.Required{},
		revel.MaxSize{40},
	)

	v.Check(hotel.State,
		revel.Required{},
		revel.MaxSize{6},
		revel.MinSize{2},
	)

	v.Check(hotel.Zip,
		revel.Required{},
		revel.MaxSize{6},
		revel.MinSize{5},
	)

	v.Check(hotel.Country,
		revel.Required{},
		revel.MaxSize{40},
		revel.MinSize{2},
	)
}*/
