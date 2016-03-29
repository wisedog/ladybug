package controllers

import (
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

// Dashboard renders a dashboard page
func Dashboard(c *interfacer.AppContext, w http.ResponseWriter, r *http.Request) error {

  var prj models.Project
  err := c.Db.Where("name = ?", c.ProjectName).First(&prj)
  
  if err.Error != nil{
    log.Error("Projects", "msg", err.Error)
    return errors.HttpError{http.StatusInternalServerError, "Not found project"}
  }
  
  // counting testcases
  tcCount := 0
  c.Db.Model(models.TestCase{}).Where("project_id = ?", prj.ID).Count(&tcCount)
  
  
  var execs []models.Execution
  
  c.Db.Where("project_id = ?", prj.ID).Find(&execs)
  
  execCount := 0
  taskCount := 0
  for _, k := range execs{
    if k.Status == ExecStatusReady{
      taskCount++
    } else if k.Status == ExecStatusInProgress{
      execCount++
    }
  }

  items := map[string]interface{}{
    "Project" : prj,
    "TestCaseCount" : tcCount,
    "ExecCount" : execCount,
    "TaskCount"  : taskCount,
    "Active_idx" : 1,
  }

  return Render2(c, w, items, "views/base.tmpl","views/dashboard.tmpl")
}

//(c *interfacer.AppContext, w http.ResponseWriter, r *http.Request) error {
func GetProjectList(limit int){
  /*var prjs []models.Project
  lm := 0
  if limit == 0 {
    lm = -1
  }else{
    lm = limit
  }
  c.Tx.Find(&prjs).Limit(lm)*/
  //TODO renderJson to cover "Show All Projects", Clicking all project in the right upper menu
  //return c.RenderJson(prjs)
}

/*
 A handler for "Show All Projects"
*/
func List(limit int){
  /*var prjs []models.Project
  c.Tx.Find(&prjs)*/
  //TODO renderJson to cover "Show All Projects", Clicking all project in the right upper menu
  //return c.Render(prjs)
}
