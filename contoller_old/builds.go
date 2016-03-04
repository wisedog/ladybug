package controllers

import (
	"time"
	"fmt"
	
	"github.com/revel/revel"
	"github.com/wisedog/ladybug/app/models"
	"github.com/wisedog/ladybug/app/routes"
	"github.com/wisedog/ladybug/app/controllers/buildtools"
)


type Builds struct {
	Application
}


// checkUser is private utilitiy function for checking 
// this user id now on connected or not.
func (c Builds) checkUser() revel.Result {
	if user := c.connected(); user == nil {
		c.Flash.Error("Please log in first")
		return c.Redirect(routes.Application.Index())
	}
	return nil
}


// Index is to render a Builds page
func (c Builds) Index(project string) revel.Result {
	var builds []models.Build

	r := c.Tx.Order("id desc").Find(&builds)
	if r.Error != nil {
		revel.ERROR.Println("Error on Builds.Index")
	}

	return c.Render(project, builds)
}


// Add is to render a Add page
func (c Builds) Add(project string) revel.Result {
	var build models.Build

	return c.Render(project, build)
}

// AddBuildItem function renders build item page.
// This page is used for only adding manual build items.
func (c Builds) AddBuildItem(project string, id int) revel.Result{
	var build models.Build
	r := c.Tx.Where("id = ?", id).First(&build)
	
	if r.Error != nil{
		revel.ERROR.Println("Could not find build project")
		c.Response.Status = 500
		return c.Render()
	}
	
	return c.Render(project, build)
}

// SaveBuildItem function handles a POST save request. 
func (c Builds) SaveBuildItem(project string, id int) revel.Result{
	var build models.Build
	r := c.Tx.Where("id = ?", id).First(&build)
	
	if r.Error != nil{
		revel.ERROR.Println("Could not find build project")
		c.Response.Status = 500
		return c.Render()
	}
	
	//get largest number
	var largest models.BuildItem
	r = c.Tx.Where("build_project_id = ?", build.ID).Order("seq desc").First(&largest)
	
	idByTool := fmt.Sprintf("%d", largest.Seq + 1)
	displayName := "#" + idByTool
	fullDisplayName := build.Name + " " + displayName
	
	bi := models.BuildItem{Toolname : "manual", BuildProjectID : build.ID, BuildAt : time.Now(), 
			Seq : largest.Seq +1, TimeStamp : 0, Status : 1, IdByTool : idByTool, 
			DisplayName : displayName , FullDisplayName : fullDisplayName,
	}
	
	c.Tx.NewRecord(bi)
	c.Tx.Create(&bi)
	
	build.BuildItemNum++
	c.Tx.Save(&build)
	
	return c.Redirect(routes.Builds.View(project, id))
}


//Save handles a saving request via POST
func (c Builds) Save(project string, build models.Build) revel.Result {
	if user := c.connected(); user == nil {
		c.Flash.Error("Please log in first")
		return c.Redirect(routes.Application.Index())
	}

	build.Validate(c.Validation)

	if c.Validation.HasErrors() {
		c.Flash.Error("invalid!")
		c.Validation.Keep()
		c.FlashParams()

		return c.Redirect(routes.Builds.Add(project))
	}
	var prj models.Project
	r := c.Tx.Where("name = ?", project).First(&prj)
	if r.Error != nil {
		c.Response.Status = 500
		return c.Render()
	}

	build.Project_id = prj.ID
	build.ToolName = "manual"

	c.Tx.NewRecord(build)
	r = c.Tx.Create(&build)

	if r.Error != nil {
		revel.ERROR.Println("Insert operation failed on Builds.Save")
	}

	return c.Redirect(routes.Builds.Index(project))

}


//View renders a page to show build detail
func (c Builds) View(project string, id int) revel.Result {
	var prj models.Project
	r := c.Tx.Where("name = ?", project).First(&prj)
	if r.Error != nil {
		c.Response.Status = 500
		return c.Render()
	}
	var build models.Build
	r = c.Tx.Where("id = ?", id).First(&build)

	if r.Error != nil {
		c.Response.Status = 500
		return c.Render()
	}
	
	var builds []models.BuildItem
	r = c.Tx.Order("id_by_tool_int desc").Where("build_project_id = ?", build.ID).Find(&builds)
	if r.Error != nil{
		c.Response.Status = 500
		return c.Render()
	}
	
	revel.INFO.Println("itemsA : ", builds)

	return c.Render(project, build, builds)
}


//Edit function renders a edit page.
func (c Builds) Edit(project string, id int) revel.Result {

	var build models.Build
	r := c.Tx.Where("id = ?", id).First(&build)

	if r.Error != nil {
		c.Response.Status = 500
		return c.Render()
	}

	return c.Render(project, build)
}



//Integrate function renders only a a page helps integration with CI tools
func (c Builds) Integrate(project string) revel.Result{
	
	return c.Render(project)
}

// AddTool has responsibility for handling a POST request adding CI tool pages.
func (c Builds) AddTool(url string, toolname string, project string) revel.Result{
	
	var prj models.Project
	r := c.Tx.Where("name = ?", project).First(&prj)
	if r.Error != nil {
		return c.RenderJson(res{Status:501, Msg:"problem"})//TEMP
	}
	
	/*
	1. Validate input
	2. fetch Jenkins API(http://52.192.120.218:8080/job/cJson/api/json)
	   if there are too much builds only get last 10 builds? 
	3. Iterate each builds and get information of artifacts, status(p/f)
	*/
	
	if toolname == "jenkins"{
		revel.INFO.Printf("Entering Jenkins routine....\n",)
		var j buildtools.Jenkins
		_, err := j.ConnectionTest(url)
		if err != nil{
			return c.RenderJson(res{Status:500, Msg:err.Error()})
		}
		
		err = j.AddJenkinsBuilds(url, prj.ID, c.Tx)
		if err != nil{
			return c.RenderJson(res{Status:500, Msg : err.Error()})
		}
		
	}else if toolname == "travis"{
		revel.INFO.Printf("Entering Travis routine....\n",)
		var t buildtools.Travis
		err := t.AddTravisBuilds(url, prj.ID, c.Tx)
		if err != nil{
			return c.RenderJson(res{Status:500, Msg : err.Error()})
		}
		
	}else{
		return c.RenderJson(res{Status:501, Msg : "Not support CI tool"})
	}
	
	return c.RenderJson(res{Status:200, Msg:"OK"})
}


// ValidationTool function checks the given url is valid.
// TODO Now this checks only jenkins connection without auth, we will add more
// next target is travis CI
func (c Builds) ValidateTool(url string, toolname string) revel.Result{
	
	if toolname == "jenkins"{
		// TODO now only support non-authorized type jenkins. 
		// another certification method should be supported
		revel.INFO.Printf("Entering Jenkins routine....\n")
		var j buildtools.Jenkins
		msg, err := j.ConnectionTest(url)
		if err != nil{
			return c.RenderJson(res{Status:500, Msg:err.Error()})
		}
		return c.RenderJson(res{Status:200, Msg:msg})
	}else if toolname == "travis"{
		revel.INFO.Printf("Entering Travis routine....\n")
		var t buildtools.Travis
		resp, err, _ := t.ConnectionTest(url)
		if err != nil{
			return c.RenderJson(res{Status:500, Msg:err.Error()})
		}
		
		var msg string
		msg += "Slug : " + resp.Repo.Slug + "\n"
		msg += "Description : " + resp.Repo.Description + "\n"
		msg += "Last Build Number : " + resp.Repo.LastBuildNumber + "\n"
		msg += "Last Build State : " + resp.Repo.LastBuildStarted + "\n"
		msg += "Last Build Finished at : " + resp.Repo.LastBuildFinish + "\n"
		
		return c.RenderJson(res{Status:200, Msg:msg})
	}else{
		// fall through
	}
	
	return c.RenderJson(res{Status:501, Msg:"Not supported tool"})
}

// GetBuildItems returns BuildItem matched with given "id"
func (c Builds) GetBuildItems(id int) revel.Result {
	var items []models.BuildItem
	
	c.Tx.Where("build_project_id = ?", id).Order("id_by_tool_int desc").Find(&items)
	
	revel.INFO.Println("items", items)
	
	return c.RenderJson(items)
}