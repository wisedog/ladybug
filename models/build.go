package models

import (
	"time"
)

// BuildItem is unit of a build. A build may have several BuildItem.
type BuildItem struct {
	BaseModel
	BuildProject	Build		// TODO index.. or something else
	BuildProjectID	int			// id for build.  TODO index? 
	Toolname		string		// jenkins, .....
	IdByTool		string		// for example, jenkins builds id is string
	IdByToolInt		int			
	DisplayName		string		// for example in jenkins : "#3"
	
	/* for example in jenkins : "cJson #3". 
	  This is the text what user see in test execution, testplan pages */
	FullDisplayName	string		
	ItemUrl			string
	ArtifactsUrl	string
	ArtifactsName	string
	Result			string		// for example in jenkins "SUCCESS"
	Status			int			// 0: failed 1: successful
	Seq				int			// for adding manual build. start from 1
	TimeStamp		int64
	BuildAt			time.Time	
}

// Build is a set of several BuildItems.  
type Build struct {
	BaseModel
	Name         	string
	Description  	string
	Project      	Project
	Project_id   	int
	BuildUrl	 	string	// build url
	ToolName       	string 	// manual, jenkins, teamcity ....
	Status		 	int 	// 0 : unknown 1 : successful, 2 : failed and so on.... 
	
	BuildItemNum	int		// total build items of this build project
	BuildItems		[]BuildItem

}

func (build *Build) Validate() map[string]string {
	errorMap := make(map[string]string)
	if build.Required(build.Name) == false {
		errorMap["Name"] = "Name is required."
	}

	if build.MaxSize(build.Name, 200) == false{
		errorMap["Name"] = errorMap["Name"] + " Size of Name exceeds 200 characters."
	}
	return errorMap
}


func (builditem *BuildItem) Validate() map[string]string {
	errorMap := make(map[string]string)
	/*if builditem.Required(builditem.Name) == false {
		errorMap["Name"] = "Name is required."
	}

	if build.MaxSize(build.Name, 200) == false{
		errorMap["Name"] = errorMap["Name"] + " Size of Name exceeds 200 characters."
	}*/

	return errorMap
}

