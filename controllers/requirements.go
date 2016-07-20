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
		log.Error("Build", "type", "database", "msg ", err.Error)
		return errors.HttpError{Status: http.StatusInternalServerError, Desc: "Not found project"}
	}

	if err := c.Db.Where("section_id = ?", sectionID).Find(&reqs); err.Error != nil {
		log.Error("Build", "type", "database", "msg ", err.Error)
		return errors.HttpError{Status: http.StatusInternalServerError, Desc: "Not found section"}
	}

	return RenderJSONWithStatus(w, reqs, http.StatusOK)
}
