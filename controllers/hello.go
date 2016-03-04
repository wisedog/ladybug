package controllers

import (
	"fmt"

  "html/template"
  "net/http"

	"github.com/wisedog/ladybug/models"
  "github.com/wisedog/ladybug/interfacer"
  "github.com/wisedog/ladybug/errors"
  log "gopkg.in/inconshreveable/log15.v2"
	)


// Welcome renders a page to select project
func Welcome(c *interfacer.AppContext, w http.ResponseWriter, r *http.Request) error {
	var user *models.User
  log.Debug("In Hello")
	if user = connected(c, r); user == nil {
    log.Info("Not found login information.")
		//c.Flash.Error("Please log in first")

    http.Redirect(w, r, "/", http.StatusFound)

    return nil
	}

  // Find Project the user is associated
  // if the project is over 10 -> make a link to show all projects
  
  
  // GORM many-to-many relationship does not work. 
  /*var projects []models.Project
  r := c.Tx.Model(&projects).Related(&user)
  */

  // Do it like below
  type user_project struct{
    UserID int
    ProjectID int
  }
  
  var p []user_project
  err := c.Db.Table("user_project").Select("user_id, project_id").Where("user_id = ?", user.ID).Scan(&p)
  
  if err.Error != nil{
    fmt.Println("Failed to select operation in Hello")
  }

  var ids []int
  for _, k := range p {
    ids = append(ids, k.ProjectID)
  }
  
  var projects []models.Project
  
  err = c.Db.Where("id in (?)", ids).Find(&projects)
  
  if err.Error != nil{
    fmt.Println("Failed to select operation in Hello2")
  }
  
  // Find Test executions for the user
  var execs []models.Execution
  err = c.Db.Where("executor_id = ? and status = 1", user.ID).Preload("Project").Preload("Plan").Find(&execs)
  
  if err.Error != nil{
    fmt.Println("Fail to count operation in Hello")
  }
  
  exec_count := len(execs)

  t, er := template.ParseFiles(
    "views/base.tmpl",
    "views/hello.tmpl",
    )

  if er != nil{
    log.Error("Error ", er )
    return errors.HttpError{http.StatusInternalServerError, "Template ParseFiles error"}
  }

  items := struct {
    User models.User
    Projects  []models.Project
    Execs   []models.Execution
    Exec_count  int
    Active_idx  int
  }{
    User: *user,
    Projects : projects,
    Execs : execs,
    Exec_count : exec_count,
    Active_idx : 0,
  }

  
  er = t.Execute(w, items)
  if er != nil{
    log.Error("Template Execution Error ", er )
    return errors.HttpError{http.StatusInternalServerError, "Template Exection Error"}
  }
  // TODO Find review for the user

  return nil
}
