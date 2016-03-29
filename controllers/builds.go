package controllers

import (
	"fmt"
  "time"
  "net/http"
  "encoding/json"

  "github.com/gorilla/mux"
  "github.com/gorilla/sessions"
  "github.com/wisedog/ladybug/models"
  "github.com/wisedog/ladybug/interfacer"
  "github.com/wisedog/ladybug/errors"
	"github.com/wisedog/ladybug/controllers/buildtools"

  log "gopkg.in/inconshreveable/log15.v2"
)

const(
  BuildFlashKey = "LADYBUG_BUILD"
  ErrorMsg = "LADUBUG_ERROR_MSG"
)

// Resp struct is for response form 
// TODO will remove Status and make http header response to represent response status
type Resp struct{
  Status int  `json:"status"`
  Msg   string  `json:"msg"`
}

//TODO remove this line and user application context cookie
var store = sessions.NewCookieStore([]byte("lady"))

// BuildsIndex is to render a Builds page
func BuildsIndex(c *interfacer.AppContext, w http.ResponseWriter, r *http.Request) error{

	var builds []models.Build
	
	if err := c.Db.Order("id desc").Find(&builds); err.Error != nil {
		log.Error("Builds", "type", "database", "msg", "Error on Builds.Index")
    return errors.HttpError{http.StatusInternalServerError, "Error on Builds.BuildsIndex"}
	}

  items := map[string]interface{}{
    "Builds" : builds,
    "Active_idx" : 3,
  }

  return Render2(c, w, items, "views/base.tmpl","views/builds/build_index.tmpl")
}


// BuildsAddProject is to render an Add page
// But rendering routine is on handleAddEditPage
func BuildsAddProject(c *interfacer.AppContext, w http.ResponseWriter, r *http.Request) error{
  return handleAddEditPage(c, w, r, false)	
}

// BuildsEditProject function renders an EDIT page.
// But rendering routine is on handleAddEditPage
func BuildsEditProject(c *interfacer.AppContext, w http.ResponseWriter, r *http.Request) error{
  return handleAddEditPage(c, w, r, true)
}

// handleAddEditPage renders an Add or Edit build page.
func handleAddEditPage(c *interfacer.AppContext, w http.ResponseWriter, r *http.Request, isEdit bool) error{
  vars := mux.Vars(r)
  id := vars["id"]

  var build models.Build
  var errorMap map[string]string

  session, e := store.Get(r, "lady")

  if e != nil{
    log.Info("error ", "msg", e.Error())
  }

  // Check if there are invalid form values from SAVE/UPDATE
  if fm := session.Flashes(BuildFlashKey); fm != nil {
    b, ok := fm[0].(*models.Build)
    if ok{
      build = *b
    }else{
      log.Debug("Build", "msg", "flash type assertion failed")
    }    

    delete(session.Values, BuildFlashKey)
    errorMap = *getErrorMap(session)
    session.Save(r, w)

  }else{
    log.Debug("Not found flash")
    if isEdit{
      if err := c.Db.Where("id = ?", id).First(&build).Error; err != nil{
        log.Error("Builds", "Type", "database", "Error ", err )
        return errors.HttpError{http.StatusInternalServerError, "not found build project on edit"}
      }
    }
  }

  items := map[string]interface{}{
    "Build" : build,
    "ID" : id,
    "Active_idx" : 3,
    "IsEdit": isEdit,
    "ErrorMap" : errorMap,
  }

  return Render2(c, w, items, "views/base.tmpl","views/builds/build_add_project.tmpl")
}

//BuildsSaveProject handles a save request via POST
func BuildsSaveProject(c *interfacer.AppContext, w http.ResponseWriter, r *http.Request) error{
	var user *models.User
  if user = connected(c, r); user == nil{
    log.Debug("Not found login information")
    http.Redirect(w, r, "/", http.StatusFound)
  }
  
  vars := mux.Vars(r)
  projectName := vars["projectName"]

  if err := r.ParseForm(); err != nil {
    log.Error("Build", "Error ", err )
    return errors.HttpError{http.StatusInternalServerError, "ParseForm failed"}
  }

  var build models.Build
  
  if err := c.Decoder.Decode(&build, r.PostForm); err != nil {
    log.Warn("Build", "Error", err, "msg", "Decode failed but go ahead")
  }

  
  if errorMap := build.Validate(); len(errorMap) > 0{
    log.Debug("Validation failed")

    session, e := store.Get(r, "lady")
    if e != nil{
      log.Warn("error ", "msg", e.Error)
    }

    session.AddFlash(build, BuildFlashKey)
    session.AddFlash(errorMap, ErrorMsg)

    session.Save(r, w)
    http.Redirect(w, r, "/project/" + projectName + "/build/add", http.StatusFound)
    return nil
  }

	var prj models.Project
	if err := c.Db.Where("name = ?", projectName).First(&prj).Error; err != nil{
    log.Error("Builds", "Type", "database", "Error ", err )
    return errors.HttpError{http.StatusInternalServerError, "Not found database select project"}  
  }

	build.Project_id = prj.ID
	build.ToolName = "manual"

	c.Db.NewRecord(build)

	if err := c.Db.Create(&build).Error; err != nil {
		log.Error("Builds", "Type", "database", "Error", err)
    return errors.HttpError{http.StatusInternalServerError, "could not create a build model"}  
	}

  http.Redirect(w, r, "/project/" + projectName + "/build", http.StatusFound)
	return nil

}


//BuildsView renders a page to show build detail
func BuildsView(c *interfacer.AppContext, w http.ResponseWriter, r *http.Request) error{
  vars := mux.Vars(r)
  id := vars["id"]

	var prj models.Project
	if err := c.Db.Where("name = ?", c.ProjectName).First(&prj).Error; err != nil{
    log.Error("Builds", "type", "database", "msg" , err )
    return errors.HttpError{http.StatusInternalServerError, "Not found project"}
  }

	var build models.Build
	if err := c.Db.Where("id = ?", id).First(&build).Error; err != nil{
    log.Error("Builds", "type", "database", "msg" , err )
    return errors.HttpError{http.StatusInternalServerError, "Not found build project"} 
  }
	
	var builds []models.BuildItem
	if err := c.Db.Order("id_by_tool_int desc").Where("build_project_id = ?", build.ID).Find(&builds).Error; err != nil{
    log.Error("Builds", "type", "database", "msg" , err )
    return errors.HttpError{http.StatusInternalServerError, "select builds error"}  
  }

  items := map[string]interface{}{
    "Build" : build,
    "BuildItems" : builds,
    "Active_idx" : 3,
  }

  return Render2(c, w, items, "views/base.tmpl", "views/builds/build_view_project.tmpl")
}

// BuildsAddItem function renders build item page.
// This page is used for only adding manual build items.
func BuildsAddItem(c *interfacer.AppContext, w http.ResponseWriter, r *http.Request) error{
  vars := mux.Vars(r)
  id := vars["id"]

  var build models.Build
  if err := c.Db.Where("id = ?", id).First(&build).Error; err != nil{
    log.Error("Build", "Type", "database", "msg", err)
    return errors.HttpError{http.StatusInternalServerError, "not found build project on add item"}
  }
  
  items := map[string]interface{}{
    "Build" : build,
    "Active_idx" : 3,
  }

  return Render2(c, w, items, "views/base.tmpl", "views/builds/build_add_item.tmpl")
}

// BuildsSaveItem function handles a POST save request. 
// This handler is for creating manual build item. If you want find creating jenkins or travis build items, 
// see ..... 
func BuildsSaveItem(c *interfacer.AppContext, w http.ResponseWriter, r *http.Request) error{
  if err := r.ParseForm(); err != nil {
    log.Error("Build", "type", "http", "msg ", err )
    return errors.HttpError{http.StatusInternalServerError, "ParseForm failed"}
  }

  var builditem models.BuildItem
  
  if err := c.Decoder.Decode(&builditem, r.PostForm); err != nil {
    log.Warn("Build", "Error", err, "msg", "Decode failed but go ahead")
  }


  var build models.Build
  if err := c.Db.Where("id = ?", builditem.BuildProjectID).First(&build).Error; err != nil{
    log.Error("Build", "type", "database", "msg ", err )
    return errors.HttpError{http.StatusInternalServerError, "Could not find build project"}
  }
  
  //get largest number
  var largest models.BuildItem
  c.Db.Where("build_project_id = ?", build.ID).Order("seq desc").First(&largest)
  
  idByTool := fmt.Sprintf("%d", largest.Seq + 1)
  displayName := "#" + idByTool
  fullDisplayName := build.Name + " " + displayName
  
  bi := models.BuildItem{Toolname : "manual", BuildProjectID : build.ID, BuildAt : time.Now(), 
      Seq : largest.Seq +1, TimeStamp : 0, Status : 1, IdByTool : idByTool, 
      DisplayName : displayName , FullDisplayName : fullDisplayName,
  }
  
  c.Db.NewRecord(bi)
  c.Db.Create(&bi)
  
  build.BuildItemNum++
  c.Db.Save(&build)
  
  target := fmt.Sprintf("/project/%s/build/view/%d", c.ProjectName, build.ID)

  http.Redirect(w, r, target, http.StatusFound)
  return nil
}



//Integrate function renders only a a page helps integration with CI tools
func Integrate(c *interfacer.AppContext, w http.ResponseWriter, r *http.Request) error{
  items := map[string]interface{}{
    "Active_idx" : 3,
  }

  return Render2(c, w, items, "views/base.tmpl", "views/builds/build_integrate.tmpl")
}


// AddTool has responsibility for handling a POST request adding CI tool pages.
func AddTool(c *interfacer.AppContext, w http.ResponseWriter, r *http.Request) error{
  var user *models.User
  if user = connected(c, r); user == nil{
    log.Debug("Not found login information")
    http.Redirect(w, r, "/", http.StatusFound)
  }

  if err := r.ParseForm(); err != nil {
    log.Error("Build", "type", "http", "msg ", err )
    return errors.HttpError{http.StatusInternalServerError, "ParseForm failed"}
  }
  url := r.FormValue("url")
  toolname := r.FormValue("toolname")

  vars := mux.Vars(r)
  projectName := vars["projectName"]

	
	var prj models.Project
	
	if err := c.Db.Where("name = ?", projectName).First(&prj); err.Error != nil {
    return RenderJson(w, Resp{Status:501, Msg:"problem"})
	}
	
	
	//1. Validate input
	//2. fetch Jenkins API(http://52.192.120.218:8080/job/cJson/api/json)
	//   if there are too much builds only get last 10 builds? 
	//3. Iterate each builds and get information of artifacts, status(p/f)
	
	if toolname == "jenkins"{
		log.Debug("Entering Jenkins routine....\n",)
		var j buildtools.Jenkins
		_, err := j.ConnectionTest(url)
		if err != nil{
			return RenderJson(w, Resp{Status:500, Msg:err.Error()})
		}
		
		err = j.AddJenkinsBuilds(url, prj.ID, c.Db)
		if err != nil{
			return RenderJson(w, Resp{Status:500, Msg : err.Error()})
		}
		
	}else if toolname == "travis"{
		log.Debug("Entering Travis routine....\n",)
		var t buildtools.Travis
		err := t.AddTravisBuilds(url, prj.ID, c.Db)
		if err != nil{
			return RenderJson(w, Resp{Status:500, Msg : err.Error()})
		}
		
	}else{
		return RenderJson(w, Resp{Status:501, Msg : "Not support CI tool"})
	}
	
	return RenderJson(w, Resp{Status:200, Msg:"OK"})
}


// ValidationTool function checks the given url is valid.
// TODO Now this checks only jenkins connection without auth, we will add more
// next target is travis CI
func  ValidateTool(c *interfacer.AppContext, w http.ResponseWriter, r *http.Request) error{
  var user *models.User
  if user = connected(c, r); user == nil{
    log.Debug("Not found login information")
    http.Redirect(w, r, "/", http.StatusFound)
  }

  if err := r.ParseForm(); err != nil {
    log.Error("Build", "type", "http", "msg ", err )
    return errors.HttpError{http.StatusInternalServerError, "ParseForm failed"}
  }
  url := r.FormValue("url")
  toolname := r.FormValue("toolname")

 
  var msg string
  var status int

	if toolname == "jenkins"{
		// TODO now only support non-authorized type jenkins. 
		// another certification method should be supported
		log.Info("Builds", "msg", "Entering Jenkins routine....")
		var j buildtools.Jenkins
		m, err := j.ConnectionTest(url)
		if err != nil{
      status = 500
      msg = err.Error()
		}else{
      status = 200
      msg = m
    }
	}else if toolname == "travis"{
		log.Info("Entering Travis routine....\n")
		var t buildtools.Travis
		resp, _, err := t.ConnectionTest(url)
		if err != nil{
      status = 500
			msg = err.Error()
		}else{
      status = 200
      msg += "Slug : " + resp.Repo.Slug + "\n"
      msg += "Description : " + resp.Repo.Description + "\n"
      msg += "Last Build Number : " + resp.Repo.LastBuildNumber + "\n"
      msg += "Last Build State : " + resp.Repo.LastBuildStarted + "\n"
      msg += "Last Build Finished at : " + resp.Repo.LastBuildFinish + "\n"
    }
	}else{
    status = 500
		msg = "Not supported tool"
	}

  return RenderJson(w, Resp{Status : status, Msg : msg})
}

// GetBuildItems returns BuildItem matched with given "id"
func  GetBuildItems(c *interfacer.AppContext, w http.ResponseWriter, r *http.Request) error{
  var user *models.User
  if user = connected(c, r); user == nil{
    log.Debug("Not found login information")
    http.Redirect(w, r, "/", http.StatusFound)
  }
  
  vars := mux.Vars(r)
  id := vars["id"]

	var items []models.BuildItem
	
	c.Db.Where("build_project_id = ?", id).Order("id_by_tool_int desc").Find(&items)

  js, err := json.Marshal(items)
  if err != nil {
    log.Error("Builds", "msg", "Json Marshalling failed in GetBuildItems")
    return errors.HttpError{http.StatusInternalServerError, "Json Marshalling failed"}
  }

  w.Header().Set("Content-Type", "application/json")
  w.Write(js)
  return nil
}