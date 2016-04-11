package controllers

import (
  "strings"
  "strconv"
  "net/http"

  "github.com/gorilla/mux"
  "github.com/wisedog/ladybug/models"
  "github.com/wisedog/ladybug/interfacer"
  "github.com/wisedog/ladybug/errors"
	
  log "gopkg.in/inconshreveable/log15.v2" 
)

//ExecIndex renders A page to show all test execution
func ExecIndex(c *interfacer.AppContext, w http.ResponseWriter, r *http.Request) error{

	var prj models.Project
	if err := c.Db.Where("name = ?", c.ProjectName).First(&prj); err.Error != nil {
    log.Error("TestPlan", "type", "database", "msg ", err.Error )
    return errors.HttpError{Status : http.StatusInternalServerError, 
      Desc : "An error while find project in ExecIndex.Index"}
	}

	var testexecs []models.Execution

	c.Db.Preload("Plan").Preload("Executor").Preload("TargetBuild").Where("project_id = ?", prj.ID).Find(&testexecs)
	
	for idx := 0; idx < len(testexecs); idx++ {
		f, p := calculatePassFail(c, &testexecs[idx])
		rv := calculateProgress(c, &testexecs[idx])
		testexecs[idx].FailCaseNum = f
		testexecs[idx].PassCaseNum = p
		testexecs[idx].Progress = rv
	}

  items := map[string]interface{}{
    "TestExecs" : testexecs,
    "Active_idx" : 5,
  }
  return Render2(c, w, items, "views/base.tmpl", "views/testexecs/execindex.tmpl")
}


// ExecDone is a POST Handler for Done button of test execution.
// This function ensures that all test cases are executed.
func ExecDone(c *interfacer.AppContext, w http.ResponseWriter, r *http.Request) error{

  if err := r.ParseForm(); err != nil {
    log.Error("TestExec", "type", "http", "msg ", err )
    return errors.HttpError{Status : http.StatusInternalServerError, Desc : "ParseForm failed"}
  }

  execID := r.FormValue("exec_id")
  //comment := r.FormValue("comment")

	var prj models.Project
	
	if err := c.Db.Where("name = ?", c.ProjectName).First(&prj); err.Error != nil {
    log.Error("TestExec", "type", "database", "msg ", err.Error )
		return RenderJSONWithStatus(w, Resp{Msg : "Find Project failed in TestExecs.Done"}, http.StatusInternalServerError)
	}

	var testexec models.Execution
	if err := c.Db.Preload("Plan").Where("id = ?", execID).First(&testexec); err.Error != nil{
    log.Error("TestExec", "type", "database", "msg ", err.Error )
		return RenderJSONWithStatus(w, Resp{Msg : "Find Test Execution entity failed in TestExecs.Done"},
        http.StatusInternalServerError)
	}
	
	// validation this execution. 
	// first, find all test results of this test execution
	// second, find all test cases belongs to this execution's testplan
	// third, check all test cases are tested
	if calculateProgress(c, &testexec) != 100{
		return RenderJSONWithStatus(w, Resp{Msg : "Not complete test execution."}, http.StatusInternalServerError)
	}

	f, _ := calculatePassFail(c, &testexec)
	
	if f > 0{
		testexec.Status = ExecStatusFail
	}else{
		testexec.Status = ExecStatusPass
	}
	
	//testexec.Status = ExecStatusDone
	//TODO add param pass_cnt, fail_cnt. 
	//TODO validation pass_cnt+fail_cnt == total count
	// if true, choose ExecStatusPass or ExecStatusFail
	// else ExecStatusDone
	
	if err := c.Db.Save(&testexec); err.Error != nil{
    log.Error("TestExec", "type", "database", "msg ", err.Error )
		return RenderJSONWithStatus(w, Resp{Msg : "Update Test Execution entity failed in TestExecs.Done"}, 
        http.StatusInternalServerError)
	}

	return RenderJSON(w, Resp{Msg:"OK"})
}

// ExecUpdateResult update execution results 
func ExecUpdateResult(c *interfacer.AppContext, w http.ResponseWriter, r *http.Request) error{
	var rv models.TestResult
	var rvTmp []models.TestResult
	var count int64

  if err := r.ParseForm(); err != nil {
    log.Error("TestExec", "type", "http", "msg ", err )
    return errors.HttpError{Status : http.StatusInternalServerError, Desc : "ParseForm failed"}
  }

  execIDStr := r.FormValue("exec_id")
  execID, _ := strconv.Atoi(execIDStr)
  caseIDStr := r.FormValue("case_id")
  caseID, _ := strconv.Atoi(caseIDStr)
  resultStr := r.FormValue("result")
  result, _ := strconv.ParseBool(resultStr)
  actual := r.FormValue("actual")
  caseVerStr := r.FormValue("case_ver")
  caseVer, _ := strconv.Atoi(caseVerStr)

	c.Db.Where("exec_id = ? and test_case_id = ?", execID, caseID).Find(&rvTmp).Count(&count)
	
	if count == 0 {
		rv.TestCaseId = caseID
		rv.ExecId = execID
		rv.Status = result
		rv.Actual = actual
		rv.TestCaseVer = caseVer
		
		//TODO validation
		c.Db.NewRecord(rv)
		
	
		if err := c.Db.Create(&rv); err.Error != nil {
			log.Error("Insert operation failed in TestExecs.UpdateResult")
			return RenderJSONWithStatus(w, Resp{Msg : "Insert operation failed in TestExecs.UpdateResult"},
            http.StatusInternalServerError)
		}
	} else{
		if err := c.Db.Where("exec_id = ? and test_case_id = ?", execID, caseID).First(&rv); err.Error != nil{
			log.Error("Select operation failed in TestExecs.UpdateResult")
			return RenderJSONWithStatus(w, Resp{Msg : "Select operation failed in TestExecs.UpdateResult"}, 
            http.StatusInternalServerError)
		}
		rv.Status = result
		rv.Actual = actual
		
		// update
		if err := c.Db.Save(&rv); err.Error != nil{
			log.Error("Update operation failed in TestExecs.UpdateResult")
			return RenderJSONWithStatus(w, Resp{Msg : "Update operation failed in TestExecs.UpdateResult"}, 
            http.StatusInternalServerError)
		}
	}
	
	//set exec's status ready to in progress
	var exec models.Execution
	
	if err := c.Db.Where("id = ?", execID).First(&exec); err.Error == nil {
		exec.Status = ExecStatusInProgress
		c.Db.Save(&exec)
	}else{
		log.Error("An error while on finding test result in TestExecs.calculateProgress")
		return RenderJSONWithStatus(w, Resp{Msg : "An error while on finding test result in TestExecs.calculateProgress"},
          http.StatusInternalServerError)
	}
	
	return RenderJSON(w, Resp{Msg : "OK"})
}

// calculatePassFail calculates pass fail of the excution
func calculatePassFail(c *interfacer.AppContext, exec *models.Execution) (int,int){
	var results []models.TestResult
	c.Db.Where("exec_id = ?", exec.ID).Find(&results)
	passCounter := 0
	failCounter := 0
	for _, k := range results{
		if k.Status == true{
			passCounter++
		}else{
			failCounter++
		}
	}
	
	return failCounter, passCounter
}


//calculateProgress calculates progress of test execution. 
//parameter "exec" should contain TestPlan information.
func calculateProgress(c *interfacer.AppContext, exec *models.Execution) int{
	var plan *models.TestPlan
	if exec.Plan.ID != 0 {
		plan = &exec.Plan
	}else{
		if exec.PlanId != 0{
			
			if err := c.Db.Where("id = ?", exec.PlanId).First(plan); err.Error != nil{
        		log.Error("An error while on finding testplain in TestExecs.calculateProgress")
				return -1
			}
		}else{
			return -1
		}
	}
	
	arr := strings.Split(plan.ExecuteCases, ",")
	var results []models.TestResult
	
	if err := c.Db.Where("exec_id = ?", exec.ID).Find(&results); err.Error != nil{
    	log.Error("An error while on finding test result in TestExecs.calculateProgress")
		return -1
	}
	
	totalCounter := len(arr)
	if totalCounter == 0 {
		log.Warn("empty array in testplan.executecases")
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



// ExecRun renders a test running page.
func ExecRun(c *interfacer.AppContext, w http.ResponseWriter, r *http.Request) error{

  vars := mux.Vars(r)
  id := vars["id"]  // test execution id

	var testexec models.Execution
	var prj models.Project
	
	if err := c.Db.Where("name = ?", c.ProjectName).First(&prj); err.Error != nil {
		log.Error("TestExec", "type", "database", "msg ", err )
    return errors.HttpError{Status : http.StatusInternalServerError, Desc : "ParseForm failed"}
	}
	
	c.Db.Preload("Plan").Preload("TargetBuild").Where("id = ?", id).First(&testexec)
	
	arr := strings.Split(testexec.Plan.ExecuteCases, ",")

  // Question. Could the test execution have no testcase?
	var cases []models.TestCase
	c.Db.Where("id in (?)", arr).Find(&cases)
	
	var results []models.TestResult
	c.Db.Order("test_case_id").Where("exec_id = ?", id).Find(&results)
	
	passCounter := 0
	failCounter := 0
	for _, k := range results{
		if k.Status == true{
			passCounter++
		}else{
			failCounter++
		}
	}

  items := map[string]interface{}{
    "TestExec" : testexec,
    "Cases" : cases,
    "Results" : results,
    "PassCounter" : passCounter,
    "FailCounter" : failCounter,
    "Active_idx" : 5,
  }
	
	return Render2(c, w, items, "views/base.tmpl", "views/testexecs/exec_run.tmpl")
}


// ExecRemove is a POST handler and deletes test execution entity.
func ExecRemove(c *interfacer.AppContext, w http.ResponseWriter, r *http.Request) error{

  if err := r.ParseForm(); err != nil {
    log.Error("TestExec", "type", "http", "msg ", err )
    return errors.HttpError{Status : http.StatusInternalServerError, Desc : "ParseForm failed"}
  }

  id := r.FormValue("id")

	//revel.INFO.Println("REMOVE : id, project", id, project)
	if err := c.Db.Where("id = ?", id).Delete(models.Execution{}); err.Error != nil {
		return RenderJSONWithStatus(w, Resp{Msg:"Problem while deleting execution"}, 
          http.StatusInternalServerError)
	}
	return RenderJSON(w, Resp{Msg : "ok"})
}


//ExecDeny makes the test execution denied
func ExecDeny(c *interfacer.AppContext, w http.ResponseWriter, r *http.Request) error{

  if err := r.ParseForm(); err != nil {
    log.Error("TestExec", "type", "http", "msg ", err )
    return errors.HttpError{Status : http.StatusInternalServerError, Desc : "ParseForm failed"}
  }

  id := r.FormValue("id")
  msg := r.FormValue("msg")

	var testexec models.Execution
	if err := c.Db.Where("id = ?", id).First(&testexec); err.Error != nil{
		return RenderJSONWithStatus(w, Resp{Msg : "Not found Test Execution"}, http.StatusInternalServerError)
	}

	testexec.Status = ExecStatusDeny
	testexec.Message = msg
	if err := c.Db.Save(testexec); err.Error != nil{
		return RenderJSONWithStatus(w, Resp{Msg : "Error while saving"}, http.StatusInternalServerError)
	}
	
	return RenderJSON(w, Resp{Msg : "ok"})
}