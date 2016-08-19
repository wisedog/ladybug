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
	log.Debug("in Design Index")
	vars := mux.Vars(r)
	project := vars["projectName"]

	var user *models.User
	if user = connected(c, r); user == nil {
		log.Info("Not found login information")
		http.Redirect(w, r, "/", http.StatusFound)
	}

	// Project list is needed to show upper right menu for moving between projects
	var prj models.Project
	rv := c.Db.Where("name = ?", project).First(&prj)
	if rv.Error != nil {
		log.Error("TestDesign", "type", "database", "msg", "Project is not found in TestDesign.Index")
		return errors.HttpError{Status: http.StatusInternalServerError, Desc: "Project is not found in TestDesign.Index"}
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
		"Active_idx": 2,
		"TreeData":   treedata,
		"Project":    prj,
	}

	return Render2(c, w, items, "views/base.tmpl", "views/testdesign/designindex.tmpl")
}

// GetAllTestCases returns all testcase of the section what "id" is matching
// Results in JSON format
func GetAllTestCases(c *interfacer.AppContext, w http.ResponseWriter, r *http.Request) error {
	log.Info("in TestDesign - GetAllTestCase")
	vars := mux.Vars(r)
	//project := vars["projectName"]
	sectionID := vars["sectionID"]
	var testcases []models.TestCase
	// Check user or return 404 or ....
	if user := connected(c, r); user == nil {
		log.Info("Not found login information")
		// TODO When extract API server from this app,
		// auth API token and return status 500 403? 400?
		http.Redirect(w, r, "/", http.StatusFound)
	}

	// if there is no result, return empty JSON array
	if err := c.Db.Where("section_id = ? ", sectionID).Preload("Category").Find(&testcases); err.Error != nil {
		log.Error("[TestDesign]", "msg", "Error while select section operation in TestDesign.GetAllTestCases")
		return errors.HttpError{Status: http.StatusInternalServerError, Desc: "Template Exection Error"}
	}

	js, err := json.Marshal(testcases)
	if err != nil {
		return errors.HttpError{Status: http.StatusInternalServerError, Desc: "Json Marshalling failed"}
	}
	// TODO define status and content struct and serialized
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)

	return nil
}
