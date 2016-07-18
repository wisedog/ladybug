package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/wisedog/ladybug/errors"
	"github.com/wisedog/ladybug/interfacer"
	"github.com/wisedog/ladybug/models"

	log "gopkg.in/inconshreveable/log15.v2"
)

//PlanIndex renders index page of TestPlan
func PlanIndex(c *interfacer.AppContext, w http.ResponseWriter, r *http.Request) error {
	var plans []models.TestPlan

	if err := c.Db.Find(&plans); err.Error != nil {
		log.Error("TestPlan", "type", "database", "msg ", err.Error)
		return errors.HttpError{Status: http.StatusInternalServerError, Desc: "An error while select all Testplans in PlanIndex"}
	}

	items := map[string]interface{}{
		"Plans":      plans,
		"Active_idx": 4,
	}
	return Render2(c, w, items, "views/base.tmpl", "views/testplans/planindex.tmpl")
}

// PlanAdd is a Handler for rendering add testplan page
func PlanAdd(c *interfacer.AppContext, w http.ResponseWriter, r *http.Request) error {
	testplan := new(models.TestPlan)
	var prj models.Project

	if err := c.Db.Preload("Users").Where("name = ?", c.ProjectName).Find(&prj); err.Error != nil {
		log.Error("TestPlan", "type", "database", "msg ", err.Error)
		return errors.HttpError{Status: http.StatusInternalServerError, Desc: "An error while select all Testplans in PlanAdd"}
	}

	var builds []models.Build
	if err := c.Db.Where("project_id = ?", prj.ID).Find(&builds); err.Error != nil {
		log.Error("TestPlan", "type", "database", "msg ", err.Error)
		return errors.HttpError{Status: http.StatusInternalServerError, Desc: "Project is not found in PlanAdd"}
	}

	var sections []models.Section
	c.Db.Where("project_id = ?", prj.ID).Where("for_test_case = ?", true).Find(&sections)

	var nodes []models.JSTreeNode

	for _, n := range sections {
		var nodeType string
		var parent string
		if n.RootNode == true {
			nodeType = "root"
			parent = "#"
		} else {
			nodeType = "default"
			parent = strconv.Itoa(n.ParentsID)
		}
		childNode := models.JSTreeNode{Id: strconv.Itoa(n.ID), Text: n.Title, Type: nodeType, Parent: parent}
		nodes = append(nodes, childNode)
	}
	treedataByte, _ := json.Marshal(nodes)
	treedata := string(treedataByte)

	items := map[string]interface{}{
		"TestPlan":   testplan,
		"TreeData":   treedata,
		"Project":    prj,
		"Builds":     builds,
		"Active_idx": 4,
	}
	return Render2(c, w, items, "views/base.tmpl", "views/testplans/planadd.tmpl")
}

// PlanSave is a POST handler for save testplan
//func (c TestPlans) Save(project string, testplan models.TestPlan, execs string) revel.Result {
func PlanSave(c *interfacer.AppContext, w http.ResponseWriter, r *http.Request) error {
	var prj models.Project

	if err := c.Db.Where("name = ?", c.ProjectName).First(&prj); err.Error != nil {
		log.Error("TestPlan", "type", "database", "msg ", err.Error)
		return errors.HttpError{Status: http.StatusInternalServerError, Desc: "Project is not found in PlanSave"}
	}

	var testplan models.TestPlan

	if err := r.ParseForm(); err != nil {
		log.Error("TestPlan", "type", "http", "msg ", err)
		return errors.HttpError{Status: http.StatusInternalServerError, Desc: "ParseForm failed"}
	}

	execs := r.FormValue("Execs")

	if err := c.Decoder.Decode(&testplan, r.PostForm); err != nil {
		log.Warn("Build", "Error", err, "msg", "Decode failed but go ahead")
	}

	var t string
	rv := handleSelected(execs, &t)

	testplan.ExecuteCases = t
	testplan.ExecCaseNum = rv
	testplan.ProjectID = prj.ID
	testplan.CreatorID = c.User.ID
	//testplan.Validate(c.Validation)

	/*if c.Validation.HasErrors() {
		c.Flash.Error("Invalid")
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(routes.TestPlans.Add(project))
	}*/

	c.Db.NewRecord(testplan)

	if err := c.Db.Create(&testplan); err.Error != nil {
		log.Error("TestPlan", "type", "database", "msg ", err.Error)
		return errors.HttpError{Status: http.StatusInternalServerError, Desc: "Insert operation failed in TestPlans.Save"}
	}

	http.Redirect(w, r, "/project/"+c.ProjectName+"/testplan", http.StatusFound)
	return nil
}

/*
//Render an edit page for test plans
func (c TestPlans) Edit(project string, id int) revel.Result {
	// TODO merge edit to add
	testplan := models.TestPlan{}

	r := c.Db.Where("id = ?", id).First(&testplan)
	if r.Error != nil {
		revel.ERROR.Println("An Error while SELECT operation for TestCase.Edit", r.Error)
		c.Response.Status = 500
		return c.Render(routes.TestPlans.Index(project))
	}

	return c.Render(project, id, testplan)
}*/

//PlanView renders a test plan page for viewing
func PlanView(c *interfacer.AppContext, w http.ResponseWriter, r *http.Request) error {

	vars := mux.Vars(r)
	id := vars["id"]

	var testplan models.TestPlan
	if err := c.Db.Preload("Creator").Preload("Executor").Where("id= ?", id).Find(&testplan); err.Error != nil {
		log.Error("TestPlan", "type", "database", "msg", err.Error)
		return errors.HttpError{Status: http.StatusInternalServerError, Desc: "An error while select operation in TestPlans.PlanView"}
	}

	// convert "1,2" like string to data used in JSTree
	arr := strings.Split(testplan.ExecuteCases, ",")
	var cases []models.TestCase
	if err := c.Db.Where("id in (?)", arr).Find(&cases); err.Error != nil {
		log.Error("TestPlan", "type", "database", "msg", err.Error)
		return errors.HttpError{Status: http.StatusInternalServerError, Desc: "An error while finding testcases SELECT operation"}
	}

	var build models.BuildItem
	if testplan.TargetBuildID != 0 {
		if err := c.Db.Where("id = ?", testplan.TargetBuildID).First(&build); err.Error != nil {
			log.Error("TestPlan", "type", "database", "msg", err.Error)
			return errors.HttpError{Status: http.StatusInternalServerError, Desc: "An error while finding Build information in TestPlans"}
		}
	}

	items := map[string]interface{}{
		"TestPlan":   testplan,
		"Cases":      cases,
		"Build":      build,
		"Active_idx": 4,
	}

	return Render2(c, w, items, "views/base.tmpl", "views/testplans/planview.tmpl")
}

//PlanDelete is a POST handler for DELETE operation for testplan
func PlanDelete(c *interfacer.AppContext, w http.ResponseWriter, r *http.Request) error {
	if err := r.ParseForm(); err != nil {
		log.Error("TestPlan", "type", "http", "msg ", err)
		return errors.HttpError{Status: http.StatusInternalServerError, Desc: "ParseForm failed"}
	}

	planID := r.FormValue("id")

	var plan models.TestPlan

	if err := c.Db.Where("id = ?", planID).First(&plan); err.Error != nil {
		log.Error("TestPlan", "type", "database", "msg", err.Error)
		return errors.HttpError{Status: http.StatusInternalServerError, Desc: "An error while select testplan operation in TestPlans.PlanDelete"}
	}

	if err := c.Db.Delete(plan); err.Error != nil {
		log.Error("TestPlan", "type", "database", "msg", err.Error)
		return errors.HttpError{Status: http.StatusInternalServerError, Desc: "An error while delete operation in TestPlans.PlanDelete"}
	}

	return RenderJSON(w, Resp{Msg: "OK"})
}

//PlanRun adds test execution with selected test plan.
func PlanRun(c *interfacer.AppContext, w http.ResponseWriter, r *http.Request) error {

	vars := mux.Vars(r)
	idStr := vars["id"]
	id, _ := strconv.Atoi(idStr)

	var prj models.Project
	if err := c.Db.Where("name = ?", c.ProjectName).First(&prj); err.Error != nil {
		log.Error("TestPlan", "type", "database", "msg", err.Error)
		return errors.HttpError{Status: http.StatusInternalServerError,
			Desc: "An error while select project operation in TestPlans.PlanFire"}
	}

	var plan models.TestPlan
	if err := c.Db.Where("id = ?", id).First(&plan); err.Error != nil {
		log.Error("TestPlan", "type", "database", "msg", err.Error)
		return errors.HttpError{Status: http.StatusInternalServerError,
			Desc: "Not found such testplan in TestPlans.PlanFire"}
	}

	// mendatory : testplan ID,project ID
	// advisory : if executor id is zero, make status N/A
	exec := models.Execution{Status: ExecStatusReady,
		ProjectID:     prj.ID,
		PlanID:        id,
		ExecutorID:    plan.ExecutorID,
		TargetBuildID: plan.TargetBuildID,
		TotalCaseNum:  plan.ExecCaseNum,
		PassCaseNum:   0,
		FailCaseNum:   0,
	}

	if plan.ExecutorID == 0 {
		exec.Status = ExecStatusNotAvailable
	}

	c.Db.NewRecord(exec)
	if err := c.Db.Create(&exec); err.Error != nil {
		log.Error("TestPlan", "type", "database", "msg", err.Error)
		return errors.HttpError{Status: http.StatusInternalServerError,
			Desc: "An error while create record in TestPlans.Fire"}
	}

	// well...
	http.Redirect(w, r, "/project/"+c.ProjectName+"/exec", http.StatusFound)
	return nil
}

//handleSelected seperates Testcases ID joined by ","
func handleSelected(content string, buf *string) int {
	arr := strings.Split(content, ",")
	var caseArray []string
	cnt := 0
	for _, s := range arr {
		if strings.HasPrefix(s, "cb") {
			cnt++
			k := strings.TrimPrefix(s, "cb")
			caseArray = append(caseArray, k)
		}
	}

	*buf = strings.Join(caseArray[:], ",")

	return cnt
}
