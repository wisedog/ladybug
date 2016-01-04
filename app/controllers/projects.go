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
	
	
	var execs []models.Execution
	
	c.Tx.Where("projecT_id = ?", prj.ID).Find(&execs)
	
	exec_count := 0
	task_count := 0
	for _, k := range execs{
		if k.Status == EXEC_STATUS_READY{
			task_count++
		} else if k.Status == EXEC_STATUS_IN_PROGRESS{
			exec_count++
		}
	}
	
	return c.Render(prj, project, tc_count, exec_count, task_count)
}

func (c Projects) GetProjectList(limit int) revel.Result{
	var prjs []models.Project
	lm := 0
	if limit == 0 {
		lm = -1
	}else{
		lm = limit
	}
	c.Tx.Find(&prjs).Limit(lm)
	//TODO renderJson to cover "Show All Projects", Clicking all project in the right upper menu
	return c.RenderJson(prjs)
}

/*
 A handler for "Show All Projects"
*/
func (c Projects) List(limit int) revel.Result {
	var prjs []models.Project
	c.Tx.Find(&prjs)
	//TODO renderJson to cover "Show All Projects", Clicking all project in the right upper menu
	return c.Render(prjs)
}
