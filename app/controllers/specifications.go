package controllers

import (
	"github.com/revel/revel"
	"encoding/json"
	"github.com/wisedog/ladybug/app/models"
	"github.com/wisedog/ladybug/app/routes"
	"strconv"
	)


type Specifications struct {
	Application
}

func (c Specifications) Spec(project string) revel.Result {
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
	c.Tx.Where("project_id = ?", prj.ID).Where("for_test_case = ?", false).Find(&sections)

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


func (c Specifications) List(project string, section_id int) revel.Result {
	var specs []models.Specification
	r := c.Tx.Where("section_id = ?", section_id).Find(&specs)
	
	if r.Error != nil{
		revel.ERROR.Println("err")
	}
	
	return c.RenderJson(specs)
}