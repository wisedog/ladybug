package controllers

import (
  "fmt"
  "strconv"

  "net/http"
	"encoding/json"
  "html/template"

	"github.com/gorilla/mux"
  "github.com/wisedog/ladybug/models"
  "github.com/wisedog/ladybug/interfacer"
	"github.com/wisedog/ladybug/errors"
  log "gopkg.in/inconshreveable/log15.v2"
)

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
	if rv.Error != nil{
		fmt.Println("Project is not found in TestDesign.Index")
		
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

  items := struct {
    User *models.User
    Project models.Project
    TreeData string
    Active_idx  int
  }{
    User: user,
    Project : prj,
    TreeData : treedata,
    Active_idx : 2,
  }

  t, er := template.ParseFiles(
    "views/base.tmpl",
    "views/designindex.tmpl",
    )

  if er != nil{
    log.Error("Error ", "type", "Template ParseFiles error", "msg", er )
    return errors.HttpError{http.StatusInternalServerError, "Template ParseFiles error"}
  }
  
  er = t.Execute(w, items)
  if er != nil{
    log.Error("Template Execution Error ", "type", "Template Exection Error", "msg", er )
    return errors.HttpError{http.StatusInternalServerError, "Template Exection Error"}
  }

  return nil
	//return c.Render(project, treedata)
}

/*
	Get all testcase of the section what "id" is matching
*/
func GetAllTestCases(c *interfacer.AppContext, w http.ResponseWriter, r *http.Request) error {
  log.Info("in TestDesign - GetAllTestCase")
  vars := mux.Vars(r)
  //project := vars["projectName"]
  sectionId := vars["sectionID"]
	var testcases []models.TestCase
	// Check user or return 404 or ....
	if user := connected(c, r); user == nil {
		log.Info("Not found login information")
    http.Redirect(w, r, "/", http.StatusFound)
	}

	// if there is no result, return empty JSON array
	if err := c.Db.Where("section_id = ? ", sectionId).Preload("Category").Find(&testcases); err.Error != nil {
    log.Error("[TestDesign]", "msg", "Error while select section operation in TestDesign.GetAllTestCases")
    return errors.HttpError{http.StatusInternalServerError, "Template Exection Error"}
	}

  js, err := json.Marshal(testcases)
  if err != nil {
    return errors.HttpError{http.StatusInternalServerError, "Json Marshalling failed"}
  }

  w.Header().Set("Content-Type", "application/json")
  w.Write(js)

  return nil
	//return c.RenderJson(testcases)
}
