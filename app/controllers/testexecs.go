package controllers

import (
	"github.com/revel/revel"
	"github.com/wisedog/ladybug/app/models"
	"github.com/wisedog/ladybug/app/routes"
	"strings"
)

//status for test execution
const(
	READY = iota
	IN_PROGRESS
	NOT_AVAILABLE
	DONE
	)

type TestExecs struct {
	Application
}

func (c TestExecs) checkUser() revel.Result {
	if user := c.connected(); user == nil {
		c.Flash.Error("Please log in first")
		return c.Redirect(routes.Application.Index())
	}
	return nil
}

/*
 A page to show all test plans
*/
func (c TestExecs) Index(project string) revel.Result {

	var prj models.Project
	r := c.Tx.Where("name = ?", project).First(&prj)
	if r.Error != nil {
		revel.ERROR.Println("An error while find project in TestExecs.Index")
		c.Response.Status = 500
		return c.Render()
	}
	var testexecs []models.Execution
	// TODO need to pre-loading.
	// Plan, Executor, TargetBuild
	c.Tx.Preload("Plan").Preload("Executor").Preload("TargetBuild").Where("project_id = ?", prj.ID).Find(&testexecs)

	

	return c.Render(project, testexecs)

}

func (c TestExecs) Done(project string, exec_id int) revel.Result{
	
	//TODO examine all testcases are excuted
	return c.Render(project)
}

/*
id is testcase's id
result is pass/fail or something else (blocked, ......)
*/
func (c TestExecs) UpdateResult(case_id int, exec_id int, result bool, actual string, case_ver int) revel.Result{
	var rv models.TestResult
	var rv_tmp []models.TestResult
	var count int64
	r := c.Tx.Find(&rv_tmp).Count(&count)
	
	type res struct{
		Status int `json:"status"`
		Msg		string `json:"msg"`
	}
	
	if count == 0 {
		rv.TestCaseId = case_id
		rv.ExecId = exec_id
		rv.Status = result
		rv.Actual = actual
		rv.TestCaseVer = case_ver
		
		//TODO validation
		c.Tx.NewRecord(rv)
		r = c.Tx.Create(&rv)
	
		if r.Error != nil {
			revel.ERROR.Println("Insert operation failed in TestExecs.UpdateResult")
			k := res{Status:500, Msg : "Insert operation failed in TestExecs.UpdateResult"}
			return c.RenderJson(k)
		}
	} else{
		revel.INFO.Println("INFO : ", exec_id, case_id)
		r = c.Tx.Where("exec_id = ? and test_case_id = ?", exec_id, case_id).First(&rv)
		if r.Error != nil{
			revel.ERROR.Println("Select operation failed in TestExecs.UpdateResult")
			k := res{Status:500, Msg : "Select operation failed in TestExecs.UpdateResult"}
			return c.RenderJson(k)
		}
		rv.Status = result
		rv.Actual = actual
		
		// update
		r = c.Tx.Save(&rv)
		if r.Error != nil{
			revel.ERROR.Println("Update operation failed in TestExecs.UpdateResult")
			k := res{Status:500, Msg : "Update operation failed in TestExecs.UpdateResult"}
			return c.RenderJson(k)
		}
	}
	
	//set exec's status ready to in progress
	var exec models.Execution
	r = c.Tx.Where("id = ?", exec_id).First(&exec)
	
	if r.Error == nil {
		exec.Status = IN_PROGRESS
		revel.INFO.Println("Exec : ", exec)
		k := c.Tx.Save(&exec)
		revel.INFO.Println("Error : ", k.Error)
	}else{
		revel.ERROR.Println("Select execution operation failed in TestExecs.UpdateResult")
	}
	
	k := res{Status : 200, Msg : "OK"}
	
	return c.RenderJson(k)
}

//TODO make testexec table, archive test history
/*
ID is test exec's. 
*/
func (c TestExecs) Run(project string, id int) revel.Result {
	var testexec models.Execution
	var prj models.Project
	r := c.Tx.Where("name = ?", project).First(&prj)
	if r.Error != nil {
		revel.ERROR.Println("An error while find project in TestExecs.Run")
		c.Response.Status = 500
		return c.Render()
	}
	
	c.Tx.Preload("Plan").Preload("TargetBuild").Where("id = ?", id).First(&testexec)
	
	revel.INFO.Println("data : ", testexec)
	
	arr := strings.Split(testexec.Plan.ExecuteCases, ",")
	var cases []models.TestCase
	r = c.Tx.Where("id in (?)", arr).Find(&cases)
	
	revel.INFO.Println("data : ", cases)
	return c.Render(project, testexec, cases)
}
