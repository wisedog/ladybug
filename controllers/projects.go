package controllers

import (
  "net/http"
  "html/template"

  "github.com/gorilla/mux"
  "github.com/wisedog/ladybug/models"
  "github.com/wisedog/ladybug/interfacer"
  "github.com/wisedog/ladybug/errors"
  log "gopkg.in/inconshreveable/log15.v2"
)


//status for test execution
const(
  EXEC_STATUS_READY = 1 + iota
  EXEC_STATUS_IN_PROGRESS
  EXEC_STATUS_NOT_AVAILABLE
  EXEC_STATUS_DONE
  EXEC_STATUS_DENY
  EXEC_STATUS_PASS
  EXEC_STATUS_FAIL
  )


func Dashboard(c *interfacer.AppContext, w http.ResponseWriter, r *http.Request) error {
  var user *models.User

  if user = connected(c, r); user == nil {
    log.Info("Not found login information")
    http.Redirect(w, r, "/", http.StatusFound)
  }

  vars := mux.Vars(r)
  name := vars["projectName"]
  log.Debug("In Project : ", "Name", name)

  var prj models.Project
  err := c.Db.Where("name = ?", name).First(&prj)
  
  if err.Error != nil{
    return errors.HttpError{http.StatusInternalServerError, "Not found project"}
  }
  
  // counting testcases
  tc_count := 0
  c.Db.Model(models.TestCase{}).Where("project_id = ?", prj.ID).Count(&tc_count)
  
  
  var execs []models.Execution
  
  c.Db.Where("project_id = ?", prj.ID).Find(&execs)
  
  exec_count := 0
  task_count := 0
  for _, k := range execs{
    if k.Status == EXEC_STATUS_READY{
      task_count++
    } else if k.Status == EXEC_STATUS_IN_PROGRESS{
      exec_count++
    }
  }

  items := struct {
    User *models.User
    Project *models.Project
    TestCaseCount int
    ExecCount   int
    TaskCount   int
    Active_idx  int
  }{
    User: user,
    Project : &prj,
    TestCaseCount : tc_count,
    ExecCount : exec_count,
    TaskCount : task_count,
    Active_idx : 1,
  }

  t, er := template.ParseFiles(
    "views/base.tmpl",
    "views/dashboard.tmpl",
    )

  if er != nil{
    log.Error("Error ", er )
    return errors.HttpError{http.StatusInternalServerError, "Template ParseFiles error"}
  }
  
  er = t.Execute(w, items)
  if er != nil{
    log.Error("Template Execution Error ", er )
    return errors.HttpError{http.StatusInternalServerError, "Template Exection Error"}
  }
  
  //return c.Render(user, prj, project, tc_count, exec_count, task_count)
  return nil
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
