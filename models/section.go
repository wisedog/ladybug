package models

import (
	"github.com/revel/revel"
	"time"
)

type Section struct {
	ID          int
	Prefix      string
	Seq         int // may use...
	Title       string
	Description string // may not used
	Status      int	// may not use
	ParentsID   int
	ProjectID   int
	RootNode    bool
	ForTestCase	bool	//TestCase or Specification

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

func (section *Section) Validate(v *revel.Validation) {
	v.Check(section.Title,
		revel.Required{},
		revel.MaxSize{400},
	)

	v.Check(section.Prefix,
		revel.Required{},
		revel.MaxSize{16},
	)
}
