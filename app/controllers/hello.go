package controllers

import (
	"github.com/revel/revel"
	"github.com/wisedog/ladybug/app/models"
	"github.com/wisedog/ladybug/app/routes"
	)

type Hello struct {
	Application
}

func (c Hello) Index() revel.Result {
	var user *models.User

	if user = c.connected(); user == nil {
		c.Flash.Error("Please log in first")
		return c.Redirect(routes.Application.Index())
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
	r := c.Tx.Table("user_project").Select("user_id, project_id").Where("user_id = ?", user.ID).Scan(&p)
	
	if r.Error != nil{
		revel.ERROR.Println("Failed to select operation in Hello")
	}
	
	var ids []int
	for _, k := range p {
		ids = append(ids, k.ProjectID)
	}
	
	var projects []models.Project
	
	r = c.Tx.Where("id in (?)", ids).Find(&projects)
	
	if r.Error != nil{
		revel.ERROR.Println("Failed to select operation in Hello2")
	}
	
	// Find Test executions for the user
	var execs []models.Execution
	r = c.Tx.Where("executor_id = ? and status = 1", user.ID).Preload("Project").Preload("Plan").Find(&execs)
	
	if r.Error != nil{
		revel.ERROR.Println("Fail to count operation in Hello")
	}
	
	exec_count := len(execs)
	
	// TODO Find review for the user
	
	return c.Render(user, projects, execs, exec_count)

}
