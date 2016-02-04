package models

import (
	"github.com/revel/revel"
	"time"
)

const (
	PRIORITY_HIGHEST = 1 + iota
	PRIORITY_HIGH
	PRIORITY_MEDIUM
	PRIORITY_LOW
	PRIORITY_LOWEST
	)

const (
	TC_STATUS_ACTIVATE = 1 + iota
	TC_STATUS_INACTIVE
	TC_STATUS_DRAFT
	)



/*
  A model for test case.
*/
type TestCase struct {
	ID            int
	Prefix        string `sql:"not null"`
	DisplayID		string	`sql:"not null"`
	Seq           int    `sql:"not null"`
	Title         string `sql:"size:400"`
	Section       Section
	SectionID     int    `sql:"index"`
	ExecutionType int    //Manual, Automated
	Status        int    // Activation, Deactivation, closed
	Description   string `sql:"size:2000"` // description of the issue
	Precondition  string `sql:"size:1000"`
	Priority      int    // 1 to 5. 1 is highest priority
	Estimated     int    // unit : min(s)
	RelatedReq    int    // TODO
	Version       int
	Steps         string `sql:"size:1000"`
	Expected      string `sql:"size:1000"`
	Project		  Project	
	ProjectID	  int
	Category		Category
	CategoryID		int
	
	// need many-to-many releationship between specification and testcase. 
	// Specifications       []Specification `gorm:"many2many:spec_cases;"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

func (testcase *TestCase) Validate(v *revel.Validation) {
	v.Required(testcase.Title)
	v.Check(testcase.Title,
		revel.Required{},
		revel.MinSize{1},
		revel.MaxSize{400},
	)
}
