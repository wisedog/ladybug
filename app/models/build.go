package models

import (
	"github.com/revel/revel"
	"time"
)

// A build may have several BuildItem
type BuildItem struct {
	ID				int
	BuildID			int			// id for build.  TODO index? 
	Toolname		string		//jenkins, .....
	IdByTool		string		//for example, jenkins builds id is string
	DisplayName		string		//for example in jenkins : "cJson #3"
	Url				string
	ArtifactsUrl	string
	ArtifactsName	string
	Result			string		//for example in jenkins "SUCCESS"
	
	TimeStamp		int64
}

type Build struct {
	ID           	int
	Name         	string
	Description  	string
	Project      	Project
	Project_id   	int
	BuildUrl	 	string	// build url
	From         	string 	// manual, jenkins, teamcity ....
	Status		 	int 	// 0 : failed 1 : successful and so on.... 
	
	BuildItemNum	int
	BuildItems		[]BuildItem

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
