package controllers

import (
  "time"
  "strings"
  "strconv"
  "net/http"

  "github.com/wisedog/ladybug/models"
  "github.com/wisedog/ladybug/interfacer"
  "github.com/wisedog/ladybug/errors"
  log "gopkg.in/inconshreveable/log15.v2"
)


//status for test execution
const(
  ExecStatusReady = 1 + iota
  ExecStatusInProgress
  ExecStatusNotAvailable
  ExecStatusDone
  ExecStatusDeny
  ExecStatusPass
  ExecStatusFail
  )

const(
  ProjectFlashKey = "LADYBUG_PROJECT"
)

// ProjectCreate renders a page to create project
func ProjectCreate(c *interfacer.AppContext, w http.ResponseWriter, r *http.Request) error {
  var project models.Project
  var errorMap map[string]string

  session, e := c.Store.Get(r, "ladybug")

  if e != nil{
    log.Info("error ", "msg", e.Error())
  }

  // Check if there are invalid form values from SAVE/UPDATE
  if fm := session.Flashes(ProjectFlashKey); fm != nil {
    p, ok := fm[0].(*models.Project)
    if ok{
      project = *p
    }else{
      log.Debug("Build", "msg", "flash type assertion failed")
    }    

    delete(session.Values, ProjectFlashKey)
    errorMap = getErrorMap(session)
    session.Save(r, w)
  }else{
    log.Debug("Project", "msg", "no flashes")
  }

  items := map[string]interface{}{
    "Project" : project,
    "ErrorMap" : errorMap,
    "Active_idx" : 0,
  }

  return Render2(c, w, items, "views/base.tmpl","views/projects/create_project.tmpl")
}

// ProjectSave validates user input and save Project in database
func ProjectSave(c *interfacer.AppContext, w http.ResponseWriter, r *http.Request) error {

  if err := r.ParseForm(); err != nil {
    log.Error("Projects", "type", "http", "msg ", err )
  }
  
  userListStr := r.FormValue("UserList")

  var project models.Project
  if err := c.Decoder.Decode(&project, r.PostForm); err != nil {
    log.Warn("Projects", "Error", err, "msg", "Decode failed but go ahead")
  }
  
  /*if c != nil{
    http.Redirect(w, r, "/hello", http.StatusFound)
    return nil
  }*/
  
  
  
  errorMap := project.Validate();
  var tmpProject models.Project
  
  if project.Name != ""{
    // check name duplication. models.Validate() can not detect duplication.
    if err := c.Db.Where("name = ?", project.Name).First(&tmpProject); err.Error == nil{
      if tmpProject.ID != 0{
        errorMap["Name"] += " Duplicated Project Name"
      }
    }
  }

  tmpProject.ID = 0
  if project.Prefix != ""{
    // check prefix duplication
    if err := c.Db.Where("prefix = ?", project.Prefix).First(&tmpProject); err.Error == nil{
      if tmpProject.ID != 0{
        errorMap["Prefix"] += " Duplicated Project Prefix"
      }
    }
  }  

  if len(errorMap) > 0{

    session, e := c.Store.Get(r, "ladybug")
    if e != nil{
      log.Warn("error ", "msg", e.Error)
    }

    log.Info("Projects", "asdf", errorMap)

    session.AddFlash(project, ProjectFlashKey)
    session.AddFlash(errorMap, ErrorMsg)

    session.Save(r, w)
    http.Redirect(w, r, "/project/create", http.StatusFound)
    return nil
  }
  
  project.Status = models.TcStatusActivate
  
  // Split on comma
  ids := strings.Split(userListStr, ",")
  
  var userList []models.User 
  if len(ids) > 0{
    if err := c.Db.Where("id in (?)", ids).Find(&userList).Error; err != nil{
      log.Error("Projects", "msg", err)
      return errors.HttpError{Status : http.StatusInternalServerError, Desc : "multiple users selection is failed"}
    }
    
    project.Users = userList
  }
  
  if rv := c.Db.NewRecord(project); rv == false {
    log.Error("Projects", "type", "database", "msg", "duplicated primary key")
    return errors.HttpError{Status : http.StatusInternalServerError, Desc : "duplicated primary key"}
  }
  
  if err := c.Db.Create(&project).Error; err != nil {
		log.Error("Projects", "Type", "database", "Error", err)
    return errors.HttpError{Status : http.StatusInternalServerError, Desc : "could not create a build model"}  
	}
  
  http.Redirect(w, r, "/project/" + project.Name, http.StatusFound)
  return nil
}


// Dashboard renders a dashboard page
func Dashboard(c *interfacer.AppContext, w http.ResponseWriter, r *http.Request) error {

  var prj models.Project
  err := c.Db.Where("name = ?", c.ProjectName).First(&prj)
  
  if err.Error != nil{
    log.Error("Projects", "msg", err.Error)
    return errors.HttpError{Status : http.StatusInternalServerError, Desc : "Not found project"}
  }
  
  // counting testcases
  /*tcCount := 0
  c.Db.Model(models.TestCase{}).Where("project_id = ?", prj.ID).Count(&tcCount)*/
  
  var execs []models.Execution
  
  c.Db.Where("project_id = ?", prj.ID).Preload("Executor").Preload("Plan").Find(&execs)
  
  execCount := 0
  taskCount := 0
  for _, k := range execs{
    if k.Status == ExecStatusReady{
      taskCount++
    } else if k.Status == ExecStatusInProgress{
      execCount++
    }
  }

  buildCount := 0
  c.Db.Model(models.Build{}).Where("project_id=?", prj.ID).Count(&buildCount)

  testPlanCount := 0
  c.Db.Model(models.TestPlan{}).Where("project_id=?", prj.ID).Count(&testPlanCount)

  var milestone models.Milestone
  c.Db.Where("project_id=?", prj.ID).Where("due_date >= ?", time.Now()).Order("due_date").First(&milestone)

  var daysLeft time.Duration
  if milestone.Name != ""{
    daysLeft = milestone.DueDate.Sub(time.Now()) 
  }

  items := map[string]interface{}{
    "Project" : prj,
    "ExecCount" : execCount,
    "BuildCount" : buildCount,
    "TaskCount"  : taskCount,
    "Tasks" : execs,
    "TestPlanCount" : testPlanCount,
    "Milestone" : milestone,
    "DaysLeft" : int(daysLeft.Hours() / 24), 
    "Active_idx" : 1,
  }

  return Render2(c, w, items, "views/base.tmpl","views/projects/dashboard.tmpl")
}

// GetProjectList returns project list in JSON format
func GetProjectList(c *interfacer.AppContext, w http.ResponseWriter, r *http.Request) error {
  var prjs []models.Project
  limitStr := r.URL.Query().Get("limit")
  limit, e := strconv.Atoi(limitStr)
  if e != nil{
    resp := map[string]interface{}{
      "ErrorMsg" : "Invalid limit value",
    }
    return RenderJSONWithStatus(w, resp, http.StatusInternalServerError)
  }
  if err := c.Db.Limit(limit).Find(&prjs).Error; err != nil{
    resp := map[string]interface{}{
      "ErrorMsg" : "Database Error",
    }
    return RenderJSONWithStatus(w, resp, http.StatusInternalServerError)
  }
  
  //TODO renderJson to cover "Show All Projects", Clicking all project in the right upper menu
  return RenderJSON(w, prjs)
}
