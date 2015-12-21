package controllers

import (
	"github.com/revel/revel"
	"github.com/wisedog/ladybug/app/models"
)

type Projects struct {
	Application
}

func (c Projects) Index(project string) revel.Result {
	var prj models.Project
	r := c.Tx.Where("name = ?", project).First(&prj)
	
	if r.Error != nil{
		c.Response.Status = 500
		return c.Render()
	}
	
	// counting testcases
	tc_count := 0
	c.Tx.Model(models.TestCase{}).Where("project_id = ?", prj.ID).Count(&tc_count)
	//db.Model(User{}).Where("name = ?", "jinzhu").Count(&count)	
	
	exec_count := 0
	c.Tx.Model(models.Execution{}).Where("project_id = ? and status = 1", prj.ID).Count(&exec_count)

	// TODO will be removed. make it AJAX call (see this.List)
	var prjs []models.Project
	c.Tx.Find(&prjs)
	
	

	return c.Render(prj, project, prjs, tc_count, exec_count)
}

/*
 A handler for "Show All Projects"
*/
func (c Projects) List() revel.Result {
	var prjs []models.Project
	c.Tx.Find(&prjs)
	//TODO renderJson to cover "Show All Projects", Clicking all project in the right upper menu
	return c.Render(prjs)
}
