package controllers

import (
	"github.com/revel/revel"
	"github.com/wisedog/ladybug/app/models"
	"github.com/wisedog/ladybug/app/routes"
	"strings"
	"strconv"
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

// struct for response in JSON 
type res struct{
	Status int `json:"status"`
	Msg		string `json:"msg"`
}

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

	c.Tx.Preload("Plan").Preload("Executor").Preload("TargetBuild").Where("project_id = ?", prj.ID).Find(&testexecs)
	
	for idx := 0; idx < len(testexecs); idx++ {
		var results []models.TestResult
		r = c.Tx.Where("exec_id = ?", testexecs[idx].ID).Find(&results)
		pass_counter := 0
		fail_counter := 0
		for _, k := range results{
			if k.Status == true{
				pass_counter++
			}else{
				fail_counter++
			}
		}
		rv := c.calculateProgress(&testexecs[idx])
		testexecs[idx].FailCaseNum = fail_counter
		testexecs[idx].PassCaseNum = pass_counter
		testexecs[idx].Progress = rv
	}
	
	return c.Render(project, testexecs)
}

/*
Handler for Done button of test execution.
This function ensures that all test cases are executed.
*/
func (c TestExecs) Done(project string, exec_id int, comment string) revel.Result{
	var prj models.Project
	r := c.Tx.Where("name = ?", project).First(&prj)
	if r.Error != nil {
		revel.ERROR.Println("An error while find project in TestExecs.Done")
		k := res{Status:500, Msg : "Find Project failed in TestExecs.Done"}
		return c.RenderJson(k)
	}
	var testexec models.Execution
	r = c.Tx.Preload("Plan").Where("id = ?", exec_id).First(&testexec)
	if r.Error != nil{
		revel.ERROR.Println("An error while find project in TestExecs.Done", exec_id)
		k := res{Status:500, Msg : "Find Test Execution entity failed in TestExecs.Done"}
		return c.RenderJson(k)
	}
	
	/* validation this execution. 
	 first, find all test results of this test execution
	 second, find all test cases belongs to this execution's testplan
	 third, check all test cases are tested*/
	progress := c.calculateProgress(&testexec)
	revel.INFO.Println("progress:", progress)
	if progress != 100{
		k := res{Status:500, Msg : "Not complete test execution."}
		return c.RenderJson(k)
	}
	
	testexec.Status = EXEC_STATUS_DONE
	//TODO add param pass_cnt, fail_cnt. 
	//TODO validation pass_cnt+fail_cnt == total count
	// if true, choose EXEC_STATUS_PASS or EXEC_STATUS_FAIL
	// else EXEC_STATUS_DONE
	
	r = c.Tx.Save(&testexec)
	if r.Error != nil{
		revel.ERROR.Println("An error while update in TestExecs.Done")
		k := res{Status:500, Msg : "Update Test Execution entity failed in TestExecs.Done"}
		return c.RenderJson(k)
	}
	
	k := res{Status:200, Msg:"OK"}
	return c.RenderJson(k)
}

/*
id is testcase's id
result is pass/fail or something else (blocked, ......)
*/
func (c TestExecs) UpdateResult(case_id int, exec_id int, result bool, actual string, case_ver int) revel.Result{
	var rv models.TestResult
	var rv_tmp []models.TestResult
	var count int64
	r := c.Tx.Where("exec_id = ? and test_case_id = ?", exec_id, case_id).Find(&rv_tmp).Count(&count)
	
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
		exec.Status = EXEC_STATUS_IN_PROGRESS
		c.Tx.Save(&exec)
	}else{
		revel.ERROR.Println("Select execution operation failed in TestExecs.UpdateResult")
	}
	
	k := res{Status : 200, Msg : "OK"}
	
	return c.RenderJson(k)
}

/*
calculate progress of test execution. 
parameter "exec" should contain TestPlan information.
*/
func (c TestExecs) calculateProgress(exec *models.Execution) int{
	var plan *models.TestPlan
	if exec.Plan.ID != 0 {
		plan = &exec.Plan
	}else{
		if exec.PlanId != 0{
			r := c.Tx.Where("id = ?", exec.PlanId).First(plan)
			if r.Error != nil{
				revel.ERROR.Println("An error while on finding testplain in TestExecs.calculateProgress ")
				return -1
			}
		}else{
			return -1
		}
	}
	
	arr := strings.Split(plan.ExecuteCases, ",")
	var results []models.TestResult
	r := c.Tx.Where("exec_id = ?", exec.ID).Find(&results)
	if r.Error != nil{
		revel.ERROR.Println("An error while on finding test result in TestExecs.calculateProgress")
		return -1
	}
	
	totalCounter := len(arr)
	if totalCounter == 0 {
		revel.WARN.Println("empty array in testplan.executecases")
		return -1;
	}
	hit := 0
	for _, k := range results{
		for _, j := range arr{
			converted, _ := strconv.Atoi(j)
			if k.TestCaseId == converted{
				hit++
			}
		}
	}
	
	return int((hit*100)/totalCounter)
}



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
	
	arr := strings.Split(testexec.Plan.ExecuteCases, ",")
	var cases []models.TestCase
	r = c.Tx.Where("id in (?)", arr).Find(&cases)
	
	var results []models.TestResult
	r = c.Tx.Where("exec_id = ?", id).Find(&results)
	
	pass_counter := 0
	fail_counter := 0
	for _, k := range results{
		if k.Status == true{
			pass_counter++
		}else{
			fail_counter++
		}
	}
	
	return c.Render(project, testexec, cases, results, pass_counter, fail_counter)
}


func (c TestExecs) Remove(project string, id int) revel.Result {
	revel.INFO.Println("REMOVE : id, project", id, project)
	r := c.Tx.Where("id = ?", id).Delete(models.Execution{})
	if r.Error != nil {
		return c.RenderJson(res{Status:500, Msg:"Problem while deleting execution"})
	}
	return c.RenderJson(res{Status : 200, Msg : "ok"})
}

func (c TestExecs) Deny(project string, id int, msg string) revel.Result {
	var testexec models.Execution
	r := c.Tx.Where("id = ?", id).First(&testexec)
	if r.Error != nil{
		return c.RenderJson(res{Status:500, Msg : "Not found Test Execution"})
	}
	testexec.Status = EXEC_STATUS_DENY
	testexec.Message = msg
	r = c.Tx.Save(testexec)
	if r.Error != nil{
		return c.RenderJson(res{Status:500, Msg : "Error while saving"})
	}
	
	return c.RenderJson(res{Status : 200, Msg : "ok"})
}