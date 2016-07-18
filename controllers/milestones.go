package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/wisedog/ladybug/errors"
	"github.com/wisedog/ladybug/interfacer"
	"github.com/wisedog/ladybug/models"

	log "gopkg.in/inconshreveable/log15.v2"
)

//MilestoneIndex renders index page of Milestone
func MilestoneIndex(c *interfacer.AppContext, w http.ResponseWriter, r *http.Request) error {
	var milestones []models.Milestone

	if err := c.Db.Find(&milestones); err.Error != nil {
		log.Error("TestPlan", "type", "database", "msg ", err.Error)
		return errors.HttpError{Status: http.StatusInternalServerError, Desc: "An error while select all Milestone in MilestoneIndex"}
	}

	items := map[string]interface{}{
		"Milestones": milestones,
		"Active_idx": 7,
	}
	return Render2(c, w, items, "views/base.tmpl", "views/milestone/milestone_index.tmpl")
}

// PlanAdd is a Handler for rendering add testplan page
func MilestoneAdd(c *interfacer.AppContext, w http.ResponseWriter, r *http.Request) error {
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
func MilestoneSave(c *interfacer.AppContext, w http.ResponseWriter, r *http.Request) error {
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

	//testplan.ExecuteCases = t
	testplan.ExecCaseNum = rv
	testplan.ProjectID = prj.ID
	testplan.CreatorID = c.User.ID

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
func MilestoneView(c *interfacer.AppContext, w http.ResponseWriter, r *http.Request) error {

	vars := mux.Vars(r)
	id := vars["id"]

	var testplan models.TestPlan
	if err := c.Db.Preload("Creator").Preload("Executor").Where("id= ?", id).Find(&testplan); err.Error != nil {
		log.Error("TestPlan", "type", "database", "msg", err.Error)
		return errors.HttpError{Status: http.StatusInternalServerError, Desc: "An error while select operation in TestPlans.PlanView"}
	}

	// convert "1,2" like string to data used in JSTree
	var arr []string //strings.Split(testplan.ExecuteCases, ",")
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
func MilestoneDelete(c *interfacer.AppContext, w http.ResponseWriter, r *http.Request) error {
	if err := r.ParseForm(); err != nil {
		log.Error("TestPlan", "type", "http", "msg ", err)
		return errors.HttpError{Status: http.StatusInternalServerError, Desc: "ParseForm failed"}
	}

	planId := r.FormValue("id")

	var plan models.TestPlan

	if err := c.Db.Where("id = ?", planId).First(&plan); err.Error != nil {
		log.Error("TestPlan", "type", "database", "msg", err.Error)
		return errors.HttpError{Status: http.StatusInternalServerError, Desc: "An error while select testplan operation in TestPlans.PlanDelete"}
	}

	if err := c.Db.Delete(plan); err.Error != nil {
		log.Error("TestPlan", "type", "database", "msg", err.Error)
		return errors.HttpError{Status: http.StatusInternalServerError, Desc: "An error while delete operation in TestPlans.PlanDelete"}
	}

	return RenderJSON(w, Resp{Msg: "OK"})
}
