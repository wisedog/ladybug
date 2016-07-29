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

// RequirementIndex renders requirement index page.
func RequirementIndex(c *interfacer.AppContext, w http.ResponseWriter, r *http.Request) error {
	var prj models.Project

	if err := c.Db.Where("name = ?", c.ProjectName).First(&prj); err.Error != nil {
		log.Error("Requirement", "type", "database", "msg ", err.Error)
		return errors.HttpError{Status: http.StatusInternalServerError,
			Desc: "Project is not found in TestDesign.Index"}
	}

	var sections []models.Section
	c.Db.Where("project_id = ?", prj.ID).Where("for_test_case = ?", false).Find(&sections)

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
		"TreeData":   treedata,
		"Active_idx": 6,
	}

	return Render2(c, w, items, "views/base.tmpl", "views/requirements/req.tmpl")
}

// RequirementList returns a list of specifications which belongs to requested section
func RequirementList(c *interfacer.AppContext, w http.ResponseWriter, r *http.Request) error {
	var reqs []models.Requirement

	vars := mux.Vars(r)
	sectionID := vars["id"]

	var prj models.Project
	if err := c.Db.Where("name = ?", c.ProjectName).First(&prj); err.Error != nil {
		log.Error("Requirment", "type", "database", "msg ", err.Error)
		return errors.HttpError{Status: http.StatusInternalServerError, Desc: "Not found project"}
	}

	if err := c.Db.Where("section_id = ?", sectionID).Find(&reqs); err.Error != nil {
		log.Error("Requirement", "type", "database", "msg ", err.Error)
		return errors.HttpError{Status: http.StatusInternalServerError, Desc: "Not found section"}
	}

	return RenderJSONWithStatus(w, reqs, http.StatusOK)
}

// AddRequirement renders just a page that a user can add a requirement
func AddRequirement(c *interfacer.AppContext, w http.ResponseWriter, r *http.Request) error {
	strSectionID := r.URL.Query().Get("section_id")

	sectionID, err := strconv.Atoi(strSectionID)
	if err != nil {
		return LogAndHTTPError(http.StatusBadRequest, "Requirement", "http", "can not convert section id")
	}

	var req models.Requirement
	req.SectionID = sectionID

	// get ReqTypes values like Use Case, Information, Feature....
	var reqTypes []models.ReqType
	if err := c.Db.Find(&reqTypes); err.Error != nil {
		return LogAndHTTPError(http.StatusInternalServerError, "Requirement", "db", "empty reqtype")
	}

	items := map[string]interface{}{
		"Requirement": req,
		"SectionID":   sectionID,
		"ReqType":     reqTypes,
		"IsEdit":      false,
		"Active_idx":  6,
	}

	return Render2(c, w, items, "views/base.tmpl", "views/requirements/reqadd.tmpl")
}

// ViewRequirement renders a page of information of the requirements.
// This page should be a detail of requested requirement and related testcases
func ViewRequirement(c *interfacer.AppContext, w http.ResponseWriter, r *http.Request) error {
	var req models.Requirement

	vars := mux.Vars(r)
	id := vars["id"]

	var prj models.Project
	if err := c.Db.Where("name = ?", c.ProjectName).First(&prj); err.Error != nil {
		log.Error("Requirement", "type", "database", "msg ", err.Error)
		return errors.HttpError{Status: http.StatusBadRequest, Desc: "Not found project"}
	}

	if err := c.Db.Where("project_id = ?", prj.ID).Where("id = ?", id).Preload("ReqType").First(&req); err.Error != nil {
		log.Error("Requirement", "type", "database", "msg ", err.Error)
		return errors.HttpError{Status: http.StatusBadRequest, Desc: "Not found requirements"}
	}

	req.StatusStr = getReqStatusI18n(req.Status)

	// Find related Test Cases with this requirement
	var testcases []models.TestCase
	if err := c.Db.Model(&req).Association("RelatedTestCases").Find(&testcases); err.Error != nil {
		log.Error("Requirement", "type", "database", "msg ", err.Error)
		return errors.HttpError{Status: http.StatusBadRequest, Desc: "Not found related testcases"}
	}

	for i := 0; i < len(testcases); i++ {
		testcases[i].PriorityStr = getPriorityI18n(testcases[i].Priority)
	}

	items := map[string]interface{}{
		"Requirement":      req,
		"RelatedTestCases": testcases,
		"Active_idx":       6,
	}

	return Render2(c, w, items, "views/base.tmpl", "views/requirements/reqview.tmpl")
}

// EditRequirement just renders a page that user can edit the requirement.
func EditRequirement(c *interfacer.AppContext, w http.ResponseWriter, r *http.Request) error {

	return nil
}

// SaveRequirement just renders a page that user can edit the requirement.
func SaveRequirement(c *interfacer.AppContext, w http.ResponseWriter, r *http.Request) error {

	return nil
}

// DeleteRequirement just renders a page that user can edit the requirement.
func DeleteRequirement(c *interfacer.AppContext, w http.ResponseWriter, r *http.Request) error {

	return nil
}

// UnlinkTestcaseRelation unlink a requirement and a related testcase
func UnlinkTestcaseRelation(c *interfacer.AppContext, w http.ResponseWriter, r *http.Request) error {

	return nil
}
