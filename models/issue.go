package models

import (
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
