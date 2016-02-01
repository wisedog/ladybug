package controllers

import (
	"encoding/json"
	"github.com/revel/revel"
	"github.com/wisedog/ladybug/app/models"
	"github.com/wisedog/ladybug/app/routes"
	"strconv"
)

type TestDesign struct {
	Application
}

func (c TestDesign) DesignIndex(project string) revel.Result {

	if user := c.connected(); user == nil {
		c.Flash.Error("Please log in first")
		return c.Redirect(routes.Application.Index())
	}

	// Project list is needed to show upper right menu for moving between projects
	var prj models.Project
	r := c.Tx.Where("name = ?", project).First(&prj)
	if r.Error != nil{
		revel.ERROR.Println("Project is not found in TestDesign.Index")
		
	}
	

	var sections []models.Section
	c.Tx.Where("project_id = ?", prj.ID).Where("for_test_case = ?", true).Find(&sections)

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

	return c.Render(project, treedata)
}

/*
	Get all testcase of the section what "id" is matching
*/
func (c TestDesign) GetAllTestCases(project string, id int) revel.Result {

	var testcases []models.TestCase
	// Check user or return 404 or ....
	if user := c.connected(); user == nil {
		c.Flash.Error("Please log in first")
		return c.RenderJson(testcases)
	}

	// if there is no result, return empty JSON array
	r := c.Tx.Where("section_id = ? ", id).Preload("Category").Find(&testcases)
	if r.Error != nil {
		revel.ERROR.Println("Error while select section operation in TestDesign.GetAllTestCases")
	}

	return c.RenderJson(testcases)
}
