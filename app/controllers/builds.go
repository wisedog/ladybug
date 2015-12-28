package controllers

import (
	"github.com/revel/revel"
	"github.com/wisedog/ladybug/app/models"
	"github.com/wisedog/ladybug/app/routes"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"strconv"
	"strings"
	"time"
)

const (
	BUILD_FAIL = iota
	BUILD_SUCCESS
	
	)

type Builds struct {
	Application
}

func (c Builds) checkUser() revel.Result {
	if user := c.connected(); user == nil {
		c.Flash.Error("Please log in first")
		return c.Redirect(routes.Application.Index())
	}
	return nil
}

/*
 A page to show testcase's information
*/
func (c Builds) Index(project string) revel.Result {
	var builds []models.Build

	r := c.Tx.Find(&builds)
	if r.Error != nil {
	}

	return c.Render(project, builds)
}

/*
Render a page to add
*/
func (c Builds) Add(project string) revel.Result {
	var build models.Build

	return c.Render(project, build)
}

/**
POST handler for save build
*/
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

	c.Tx.NewRecord(build)
	r = c.Tx.Create(&build)

	if r.Error != nil {
		revel.ERROR.Println("Insert operation failed on Builds.Save")
	}

	return c.Redirect(routes.Builds.Index(project))

}

/**
Render a page to view
*/
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
	r = c.Tx.Where("build_project_id = ?", build.ID).Find(&builds)
	if r.Error != nil{
		c.Response.Status = 500
		return c.Render()
	}

	return c.Render(project, build, builds)
}

/**
Render a page to edit
*/
func (c Builds) Edit(project string, id int) revel.Result {

	var build models.Build
	r := c.Tx.Where("id = ?", id).First(&build)

	if r.Error != nil {
		c.Response.Status = 500
		return c.Render()
	}

	return c.Render(project, build)
}


/*
Render a page helps integration with CI tools
*/
func (c Builds) Integrate(project string) revel.Result{
	
	return c.Render(project)
}

/*
POST Handler for adding CI tool
*/
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
	// TODO the code below only handle Jenkins
	
	if strings.HasSuffix(url, "/api/json") == false{
    	url = url + "/api/json"
    }
    
    body, err := c.getJenkinsJobInfo(url)
    if err != nil {
    	return c.RenderJson(res{Status:501, Msg:"problem"})//TEMP
    }
    var dat map[string]interface{}
    if err := json.Unmarshal(body, &dat); err != nil {
      return c.RenderJson(res{Status:502, Msg:"Problem"})	//TEMP
    }
    
    name := dat["name"].(string)
    nextBuildNum := int(dat["nextBuildNumber"].(float64))
    lastSuccessfulBuild := dat["lastSuccessfulBuild"].(map[string]interface{})
    lastSucessfulBuildNum := int(lastSuccessfulBuild["number"].(float64))
    
    
    // get status for building. it may be successful or failed
    status := 0
    if nextBuildNum -1 == lastSucessfulBuildNum {
    	status = BUILD_SUCCESS
    } else{
    	status = BUILD_FAIL
    }

    builds := dat["builds"].([]interface{})
    
    job := models.Build{
    	Name : name,
    	Description : dat["description"].(string),
    	Project_id : prj.ID,
    	BuildUrl : dat["url"].(string),
    	Status : status,
    	ToolName : "jenkins",
    	BuildItemNum : len(builds),
    }
    r = c.Tx.Save(&job)
    if r.Error != nil{
    	return c.RenderJson(res{Status:503, Msg:"Error while Saving "})	//TEMP
    }
    
    
    for idx, b := range builds {
    	if idx > 10 {
    		break
    	}
    	
    	k := b.(map[string]interface{})
    	
    	targetUrl := k["url"].(string) + "/api/json"
    	info, err := c.getJenkinsJobInfo(targetUrl)
    	if err != nil{
    		continue
    	}
    	
    	var data map[string]interface{}
	    if err := json.Unmarshal(info, &data); err != nil {
	      return c.RenderJson(res{Status:502, Msg:"Problem"})	//TEMP
	    }
    	
    	timestamp := int64(data["timestamp"].(float64))
    	displayname := data["displayName"].(string)
    	idbytool := data["id"].(string)
    	result := data["result"].(string)
    	url := data["url"].(string)
    	
    	
    	artifacts := data["artifacts"].([]interface{})
    	
    	num := len(artifacts)
    	var artifactsname string
    	var artifactsurl	string
    	
    	if num > 1 {
    		// TODO It is not supported to link multiple artifacts now
    		artifactsurl = url
    		artifactsname = "Multiple"
    	}else if num == 1 {
    		a := artifacts[0].(map[string]interface{})
    		artifactsname = a["fileName"].(string)
    		artifactsurl = url + "artifact/" + a["relativePath"].(string)
    	}else {
    		artifactsurl = ""
    		artifactsname = ""
    	}
    	
    	// timestamp of jenkins build item represents in millisecond.
    	// so divide by 1000 (1 second = 1000 milliseconds)
    	buildat := time.Unix(int64(timestamp/1000),0)
    	
    	var rv int 
    	if result == "SUCCESS"{
    		rv = BUILD_SUCCESS
    	}else {
    		rv = BUILD_FAIL
    	}
    	
    	
    	elem := models.BuildItem{
    		BuildProjectID : job.ID,
    		IdByTool : idbytool,
    		DisplayName : displayname,
    		FullDisplayName : name + " " + displayname,
    		Url : url,
    		ArtifactsUrl : artifactsurl,
    		ArtifactsName : artifactsname,
    		Result : result,
    		Toolname : "jenkins",
    		TimeStamp : timestamp,
    		BuildAt : buildat,
    		Status : rv,
    	}
    	
    	revel.INFO.Println("element", elem)
    	
    	// save to BuildItem
    	r = c.Tx.Save(&elem)
    	if r.Error != nil{
    		return c.RenderJson(res{Status:504, Msg:"Error while saving"})
    	}
    }
	
	// redirect index
	return c.RenderJson(res{Status:200, Msg:"OK"})
}

func (c Builds) ValidateTool(url string, toolname string) revel.Result{
	
    // TODO the code below only handle Jenkins
    if strings.HasSuffix(url, "/api/json") == false{
    	url = url + "/api/json"
    }
	
	body, err := c.getJenkinsJobInfo(url)
	if err != nil {
		return c.RenderJson(res{Status:501, Msg:"Internal server error"})
	}
	
	var dat map[string]interface{}
	var msg string
	
    if err := json.Unmarshal(body, &dat); err != nil {
      return c.RenderJson(res{Status:502, Msg:"Json Unmarshalling is failed"})
    }
    
    msg += "Job name : "+ dat["name"].(string) + "\n"
    msg += "URL : " + dat["url"].(string) + "\n"
    
    build := dat["lastBuild"].(map[string]interface{})
    
    // FIXME runtime error build["number"].(int).... why? 3 is float?
    k := strconv.Itoa(int(build["number"].(float64)))
    msg += "LastBuild Number : " + k + "\n"
    
    build = dat["lastSuccessfulBuild"].(map[string]interface{})
    k = strconv.Itoa(int(build["number"].(float64)))
    msg += "Last Successful Build : " + k + "\n"
	
	r := res{Status:200, Msg:msg}
	
	return c.RenderJson(r)
}

func (c Builds) getJenkinsJobInfo(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		revel.ERROR.Println("An error while GET", err)
		return nil, err
	}
	defer resp.Body.Close()
	
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil{
		revel.ERROR.Println("An error while readall", err)
		return nil, err
	}
	return body, nil
}

