package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/wisedog/ladybug/errors"
	"github.com/wisedog/ladybug/interfacer"
	"github.com/wisedog/ladybug/models"

	log "gopkg.in/inconshreveable/log15.v2"
)

//ExecIndex renders A page to show all test executions
func ExecIndex(c *interfacer.AppContext, w http.ResponseWriter, r *http.Request) error {

	var prj models.Project
	if err := c.Db.Where("name = ?", c.ProjectName).First(&prj); err.Error != nil {
		log.Error("TestPlan", "type", "database", "msg ", err.Error)
		return errors.HttpError{Status: http.StatusInternalServerError,
			Desc: "An error while find project in ExecIndex.Index"}
	}

	var testexecs []models.Execution

	c.Db.Preload("Plan").Preload("Executor").Preload("TargetBuild").Where("project_id = ?", prj.ID).Find(&testexecs)

	//Update progress
	for idx := 0; idx < len(testexecs); idx++ {
		prg := float32(testexecs[idx].PassCaseNum+testexecs[idx].FailCaseNum) / float32(testexecs[idx].TotalCaseNum) * 100
		testexecs[idx].Progress = int(prg)
	}

	items := map[string]interface{}{
		"TestExecs":  testexecs,
		"Active_idx": 5,
	}
	return Render2(c, w, items, "views/base.tmpl", "views/testexecs/execindex.tmpl")
}

// ExecDone is a POST Handler for Done button of test execution.
// This function ensures that all test cases are executed.
func ExecDone(c *interfacer.AppContext, w http.ResponseWriter, r *http.Request) error {

	if err := r.ParseForm(); err != nil {
		log.Error("TestExec", "type", "http", "msg ", err)
		return errors.HttpError{Status: http.StatusInternalServerError, Desc: "ParseForm failed"}
	}

	execID := r.FormValue("exec_id")
	//comment := r.FormValue("comment")

	var prj models.Project

	if err := c.Db.Where("name = ?", c.ProjectName).First(&prj); err.Error != nil {
		log.Error("TestExec", "type", "database", "msg ", err.Error)
		return RenderJSONWithStatus(w, Resp{Msg: "Find Project failed in TestExecs.Done"}, http.StatusInternalServerError)
	}

	var testexec models.Execution
	if err := c.Db.Preload("Plan").Where("id = ?", execID).First(&testexec); err.Error != nil {
		log.Error("TestExec", "type", "database", "msg ", err.Error)
		return RenderJSONWithStatus(w, Resp{Msg: "Find Test Execution entity failed in TestExecs.Done"},
			http.StatusInternalServerError)
	}

	// validation this execution is finished.
	// first, just verify that (PassCaseNum + FailCaseNum == TotalCaseNum)
	if testexec.PassCaseNum+testexec.FailCaseNum < testexec.TotalCaseNum {
		return RenderJSONWithStatus(w, Resp{Msg: "Not complete test execution."}, http.StatusInternalServerError)
	}
	// TODO below second, find all test results of this test execution and count
	count := 0
	c.Db.Model(&models.TestCaseResult{}).Count(&count)
	if count != testexec.TotalCaseNum {
		return RenderJSONWithStatus(w, Resp{Msg: "Not complete test execution."}, http.StatusInternalServerError)
	}

	if testexec.FailCaseNum > 0 {
		testexec.Status = ExecStatusFail
	} else {
		testexec.Status = ExecStatusPass
	}

	if err := c.Db.Model(&testexec).Update("status", testexec.Status); err.Error != nil {
		log.Error("TestExec", "type", "database", "msg ", err.Error)
		return RenderJSONWithStatus(w, Resp{Msg: "Update Test Execution entity failed in TestExecs.Done"},
			http.StatusInternalServerError)
	}

	return RenderJSONWithStatus(w, Resp{Msg: "OK"}, http.StatusOK)
}

// ExecUpdateResult updates the result of each test case
func ExecUpdateResult(c *interfacer.AppContext, w http.ResponseWriter, r *http.Request) error {
	var rv models.TestCaseResult

	if err := r.ParseForm(); err != nil {
		log.Error("TestExec", "type", "http", "msg ", err)
		return errors.HttpError{Status: http.StatusInternalServerError, Desc: "ParseForm failed"}
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

	var exec models.Execution

	if err := c.Db.Where("id = ?", execID).First(&exec); err.Error != nil {
		log.Error("An error while on finding test result in TestExecs.calculateProgress")
		return RenderJSONWithStatus(w, Resp{Msg: "An error while on finding test result in TestExecs.calculateProgress"},
			http.StatusInternalServerError)
	}

	var rvTmp []models.TestCaseResult

	c.Db.Where("exec_id = ? and test_case_id = ?", execID, caseID).Find(&rvTmp)

	// if there is no exist test result
	if len(rvTmp) == 0 {
		rv.TestCaseID = caseID
		rv.ExecID = execID
		rv.Status = result
		rv.Actual = actual
		rv.TestCaseVer = caseVer

		//TODO validation
		c.Db.NewRecord(rv)

		if err := c.Db.Create(&rv); err.Error != nil {
			log.Error("Insert operation failed in TestExecs.UpdateResult")
			return RenderJSONWithStatus(w, Resp{Msg: "Insert operation failed in TestExecs.UpdateResult"},
				http.StatusInternalServerError)
		}

		//set exec's status ready to in progress
		exec.Status = ExecStatusInProgress
		if rv.Status == true {
			exec.PassCaseNum = exec.PassCaseNum + 1
		} else {
			exec.FailCaseNum = exec.FailCaseNum + 1
		}

		c.Db.Save(&exec)

	} else {
		if err := c.Db.Where("exec_id = ? and test_case_id = ?", execID, caseID).First(&rv); err.Error != nil {
			log.Error("Select operation failed in TestExecs.UpdateResult")
			return RenderJSONWithStatus(w, Resp{Msg: "Select operation failed in TestExecs.UpdateResult"},
				http.StatusInternalServerError)
		}

		if rv.Status == true {
			if result == true {
				// same value. not updated
			} else {
				// update. decrease execution.passcounter and increase failcounter
				exec.PassCaseNum = exec.PassCaseNum - 1
				exec.FailCaseNum = exec.FailCaseNum + 1
			}
		} else {
			if result == true {
				// update increase execution.passcounter and decrease failcounter
				exec.PassCaseNum = exec.PassCaseNum + 1
				exec.FailCaseNum = exec.FailCaseNum - 1
			} else {
				// same value. not updated
			}
		}
		rv.Status = result
		rv.Actual = actual
		rv.TestPlanID = exec.PlanID

		// update testcase result
		if err := c.Db.Save(&rv); err.Error != nil {
			log.Error("Update operation failed in TestExecs.UpdateResult")
			return RenderJSONWithStatus(w, Resp{Msg: "Update operation failed in TestExecs.UpdateResult"},
				http.StatusInternalServerError)
		}

		// update exec
		//db.Model(&user).Updates(map[string]interface{}{"name": "hello", "age": 18, "actived": false})
		if err := c.Db.Model(&exec).Updates(
			map[string]interface{}{
				"pass_case_num": exec.PassCaseNum,
				"fail_case_num": exec.FailCaseNum,
				"status":        exec.Status}); err.Error != nil {
			log.Error("Update operation failed in TestExecs.UpdateResult")
			return RenderJSONWithStatus(w, Resp{Msg: "Update operation failed in TestExecs.UpdateResult"},
				http.StatusInternalServerError)
		}
	}

	return RenderJSONWithStatus(w, Resp{Msg: "OK"}, http.StatusOK)
}

// ExecRun renders a test running page.
func ExecRun(c *interfacer.AppContext, w http.ResponseWriter, r *http.Request) error {

	vars := mux.Vars(r)
	id := vars["id"] // test execution id

	var testexec models.Execution
	var prj models.Project

	if err := c.Db.Where("name = ?", c.ProjectName).First(&prj); err.Error != nil {
		log.Error("TestExec", "type", "database", "msg ", err)
		return errors.HttpError{Status: http.StatusInternalServerError, Desc: "ParseForm failed"}
	}

	c.Db.Preload("Plan").Preload("TargetBuild").Where("id = ?", id).First(&testexec)

	arr := strings.Split(testexec.Plan.ExecuteCases, ",")

	// Question. Could the test execution have no testcase?
	var cases []models.TestCase
	c.Db.Where("id in (?)", arr).Find(&cases)

	var results []models.TestCaseResult
	c.Db.Order("test_case_id").Where("exec_id = ?", id).Find(&results)

	passCounter := 0
	failCounter := 0
	for _, k := range results {
		if k.Status == true {
			passCounter++
		} else {
			failCounter++
		}
	}

	items := map[string]interface{}{
		"TestExec":    testexec,
		"Cases":       cases,
		"Results":     results,
		"PassCounter": passCounter,
		"FailCounter": failCounter,
		"Active_idx":  5,
	}

	return Render2(c, w, items, "views/base.tmpl", "views/testexecs/exec_run.tmpl")
}

// ExecRemove is a POST handler and deletes test execution entity.
func ExecRemove(c *interfacer.AppContext, w http.ResponseWriter, r *http.Request) error {

	if err := r.ParseForm(); err != nil {
		log.Error("TestExec", "type", "http", "msg ", err)
		return errors.HttpError{Status: http.StatusInternalServerError, Desc: "ParseForm failed"}
	}

	id := r.FormValue("id")

	if err := c.Db.Where("id = ?", id).Delete(models.Execution{}); err.Error != nil {
		return RenderJSONWithStatus(w, Resp{Msg: "Problem while deleting execution"},
			http.StatusInternalServerError)
	}
	return RenderJSON(w, Resp{Msg: "ok"})
}

//ExecDeny makes the test execution denied
func ExecDeny(c *interfacer.AppContext, w http.ResponseWriter, r *http.Request) error {

	if err := r.ParseForm(); err != nil {
		log.Error("TestExec", "type", "http", "msg ", err)
		return errors.HttpError{Status: http.StatusInternalServerError, Desc: "ParseForm failed"}
	}

	id := r.FormValue("id")
	msg := r.FormValue("msg")

	var testexec models.Execution
	if err := c.Db.Where("id = ?", id).First(&testexec); err.Error != nil {
		return RenderJSONWithStatus(w, Resp{Msg: "Not found Test Execution"}, http.StatusInternalServerError)
	}

	testexec.Status = ExecStatusDeny
	testexec.Message = msg
	if err := c.Db.Save(testexec); err.Error != nil {
		return RenderJSONWithStatus(w, Resp{Msg: "Error while saving"}, http.StatusInternalServerError)
	}

	return RenderJSON(w, Resp{Msg: "ok"})
}
