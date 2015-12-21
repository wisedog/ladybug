package models

import (
	"github.com/revel/revel"
	"time"
)

type Build struct {
	ID           int
	Name         string
	BuildNumber  int
	Description  string
	Project      Project
	Project_id   int
	DownloadLink string
	From         string // manual, jenkins, teamcity ....
	ReleaseAt    time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

func (build *Build) Validate(v *revel.Validation) {
	v.Check(build.Name,
		revel.Required{},
		revel.MaxSize{200},
	)

}
