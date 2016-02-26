package models

import (
	"github.com/revel/revel"
	"time"
)

// A build may have several BuildItem
type BuildItem struct {
	ID				int
	BuildProject	Build		// TODO index.. or something else
	BuildProjectID	int			// id for build.  TODO index? 
	Toolname		string		// jenkins, .....
	IdByTool		string		// for example, jenkins builds id is string
	IdByToolInt		int			
	DisplayName		string		// for example in jenkins : "#3"
	
	/* for example in jenkins : "cJson #3". 
	  This is the text what user see in test execution, testplan pages */
	FullDisplayName	string		
	Url				string
	ArtifactsUrl	string
	ArtifactsName	string
	Result			string		// for example in jenkins "SUCCESS"
	Status			int			// 0: failed 1: successful
	Seq				int			// for adding manual build. start from 1
	TimeStamp		int64
	BuildAt			time.Time
	
}

type Build struct {
	ID           	int
	Name         	string
	Description  	string
	Project      	Project
	Project_id   	int
	BuildUrl	 	string	// build url
	ToolName       	string 	// manual, jenkins, teamcity ....
	Status		 	int 	// 0 : unknown 1 : successful, 2 : failed and so on.... 
	
	BuildItemNum	int		// total build items of this build project
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
