package controllers

import (
	"fmt"
  "net/http"

	"github.com/wisedog/ladybug/models"
  "github.com/wisedog/ladybug/interfacer"
  log "gopkg.in/inconshreveable/log15.v2"
	)


// Welcome renders a page to select project
func Welcome(c *interfacer.AppContext, w http.ResponseWriter, r *http.Request) error {
	
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
  
  if err := c.Db.Table("user_project").Select("user_id, project_id").
            Where("user_id = ?", c.User.ID).Scan(&p);err.Error != nil{
    fmt.Println("Failed to select operation in Hello")
  }

  var ids []int
  for _, k := range p {
    ids = append(ids, k.ProjectID)
  }
  
  var projects []models.Project
  
  if err := c.Db.Where("id in (?)", ids).Find(&projects); err.Error != nil{
    fmt.Println("Failed to select operation in Hello2")
    log.Error("Hello")
  }
  
  // Find Test executions for the user
  var execs []models.Execution
  
  if err := c.Db.Where("executor_id = ? and status = 1", c.User.ID).
            Preload("Project").Preload("Plan").Find(&execs); err.Error != nil{
    fmt.Println("Fail to count operation in Hello")
  }
  
  execCount := len(execs)
  
  items := map[string]interface{}{
    "Projects" : projects,
    "Execs" : execs,
    "ExecCount" : execCount,
    "Active_idx" : 0,
  }

  // TODO Find review for the user

  return Render2(c, w, items, "views/base.tmpl", "views/hello/hello.tmpl")
}
