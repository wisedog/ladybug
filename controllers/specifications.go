package controllers

import (
	"strconv"
	"encoding/json"
  "net/http"

  "github.com/gorilla/mux"
  "github.com/wisedog/ladybug/models"
  "github.com/wisedog/ladybug/interfacer"
  "github.com/wisedog/ladybug/errors"

  log "gopkg.in/inconshreveable/log15.v2"
)

// SpecIndex renders Specification index page.
func SpecIndex(c *interfacer.AppContext, w http.ResponseWriter, r *http.Request) error{
	var user *models.User
  if user = connected(c, r); user == nil{
    log.Debug("Not found login information")
    http.Redirect(w, r, "/", http.StatusFound)
  }

  vars := mux.Vars(r)
  projectName := vars["projectName"]

	// Project list is needed to show upper right menu for moving between projects
	var prj models.Project
	
	if err := c.Db.Where("name = ?", projectName).First(&prj); err.Error != nil{
    log.Error("Specification", "type", "database", "msg ", err.Error )
    return errors.HttpError{http.StatusInternalServerError, "Project is not found in TestDesign.Index"}
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
    "User": user,
    "TreeData" : treedata,
    "ProjectName" : projectName,
    "Active_idx" : 6,
  }

  return Render(w, items, "views/base.tmpl", "views/specifications/spec.tmpl")
}


//SpecList returns a list of specifications which belongs to requested section
func SpecList(c *interfacer.AppContext, w http.ResponseWriter, r *http.Request) error{
	var specs []models.Specification

  vars := mux.Vars(r)
  projectName := vars["projectName"]
  section_id := vars["id"]
	
  var prj models.Project
  if err := c.Db.Where("name = ?", projectName).First(&prj); err.Error != nil{
    log.Error("Build", "type", "database", "msg ", err.Error )
    return errors.HttpError{http.StatusInternalServerError, "Not found project"}
  }
	
	if err := c.Db.Where("section_id = ?", section_id).Find(&specs); err.Error != nil{
    log.Error("Build", "type", "database", "msg ", err.Error )
    return errors.HttpError{http.StatusInternalServerError, "Not found section"}
	}
	
  return RenderJSON(w, specs)
}