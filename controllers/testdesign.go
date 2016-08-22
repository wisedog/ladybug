package controllers

import (
	"strconv"

	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/wisedog/ladybug/errors"
	"github.com/wisedog/ladybug/interfacer"
	"github.com/wisedog/ladybug/models"
	log "gopkg.in/inconshreveable/log15.v2"
)

// DesignIndex renders index page of /project/{projectName}/design
func DesignIndex(c *interfacer.AppContext, w http.ResponseWriter, r *http.Request) error {
	// Project list is needed to show upper right menu for moving between projects
	var prj models.Project
	if err := c.Db.Where("name = ?", c.ProjectName).First(&prj); err.Error != nil {
		log.Error("TestDesign", "type", "database", "msg", "Project is not found in TestDesign.Index")
		return errors.HttpError{Status: http.StatusInternalServerError, Desc: "Project is not found in TestDesign.Index"}
	}

	var sections []models.Section
	if err := c.Db.Where("project_id = ?", prj.ID).Where("for_test_case = ?", true).Find(&sections); err.Error != nil {
		log.Error("TestDesign", "type", "database", "msg", err.Error.Error())
		return errors.HttpError{Status: http.StatusInternalServerError, Desc: "Project is not found in TestDesign.Index"}
	}
	treedata := getJSTreeNodeData(sections)

	items := map[string]interface{}{
		"Active_idx": 2,
		"TreeData":   treedata,
		"Project":    prj,
	}

	return Render2(c, w, items, "views/base.tmpl", "views/testdesign/designindex.tmpl")
}

// getJSTreeNodeData converts sections to JSTree data structure in JSON format
func getJSTreeNodeData(sections []models.Section) *string {
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
		childNode := models.JSTreeNode{ID: strconv.Itoa(n.ID), Text: n.Title, Type: nodeType, Parent: parent}
		nodes = append(nodes, childNode)
	}
	treedataByte, _ := json.Marshal(nodes)
	treedata := string(treedataByte)
	return &treedata
}

// GetAllTestCases returns all testcase of the section what "id" is matching
// Results in JSON format
func GetAllTestCases(c *interfacer.AppContext, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	sectionID := vars["sectionID"]
	var testcases []models.TestCase

	// if there is no result, return empty JSON array
	if err := c.Db.Where("section_id = ? ", sectionID).Preload("Category").Find(&testcases); err.Error != nil {
		log.Error("[TestDesign]", "msg", "Error while select section operation in TestDesign.GetAllTestCases")
		return errors.HttpError{Status: http.StatusInternalServerError, Desc: "Template Exection Error"}
	}

	js, err := json.Marshal(testcases)
	if err != nil {
		return errors.HttpError{Status: http.StatusInternalServerError, Desc: "Json Marshalling failed"}
	}
	return RenderJSONWithStatus(w, string(js), http.StatusOK)
}
