package controllers

import (
	"encoding/json"
	"github.com/revel/revel"
	"github.com/wisedog/ladybug/app/models"
	"github.com/wisedog/ladybug/app/routes"
	"strconv"
	"strings"
)

type TestPlans struct {
	Application
}

func (c TestPlans) checkUser() revel.Result {
	if user := c.connected(); user == nil {
		c.Flash.Error("Please log in first")
		return c.Redirect(routes.Application.Index())
	}
	return nil
}

func (c TestPlans) Index(project string) revel.Result {
	var plans []models.TestPlan
	r := c.Tx.Find(&plans)

	if r.Error != nil {
		c.Response.Status = 500
		return c.Render()
	}

	return c.Render(project, plans)
}

/**
A Handler for rendering add testplan page
*/
func (c TestPlans) Add(project string) revel.Result {
	testplan := new(models.TestPlan)
	var prj models.Project
	r := c.Tx.Preload("Users").Where("name = ?", project).Find(&prj)

	if r.Error != nil {
		revel.ERROR.Println("An error while find user")
		c.Response.Status = 500
		return c.Render()
	}

	var builds []models.Build
	r = c.Tx.Where("project_id = ?", prj.ID).Find(&builds)

	var sections []models.Section
	c.Tx.Find(&sections)

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

	var execs string

	return c.Render(project, testplan, prj, execs, builds, treedata)
}

/**
A POST handler for save testplan
*/
func (c TestPlans) Save(project string, testplan models.TestPlan, execs string) revel.Result {
	var user *models.User
	if user = c.connected(); user == nil {
		c.Flash.Error("Please log in first")
		return c.Redirect(routes.Application.Index())
	}

	var prj models.Project
	r := c.Tx.Where("name = ?", project).First(&prj)

	if r.Error != nil {
		c.Response.Status = 500
		return c.Render()
	}
	
	var t string
	rv := handleSelected(execs, &t)
	
	revel.INFO.Println("buf : ", t)
	revel.INFO.Println("rv : ", rv)
	
	testplan.ExecuteCases = t
	testplan.ExecCaseNum = rv
	testplan.Project_id = prj.ID
	testplan.CreatorId = user.ID
	testplan.Validate(c.Validation)

	if c.Validation.HasErrors() {
		c.Flash.Error("Invalid")
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(routes.TestPlans.Add(project))
	}

	c.Tx.NewRecord(testplan)
	r = c.Tx.Create(&testplan)

	if r.Error != nil {
		revel.ERROR.Println("Insert operation failed in TestPlans.Save")
		//TODO status = 500 or something else
	}
	return c.Redirect(routes.TestPlans.Index(project))
}

/**
Render an edit page for test plans
*/
func (c TestPlans) Edit(project string, id int) revel.Result {
	// TODO merge edit to add
	testplan := models.TestPlan{}

	r := c.Tx.Where("id = ?", id).First(&testplan)
	if r.Error != nil {
		revel.ERROR.Println("An Error while SELECT operation for TestCase.Edit", r.Error)
		c.Response.Status = 500
		return c.Render(routes.TestPlans.Index(project))
	}

	return c.Render(project, id, testplan)
}

/**
Render a test plan page for viewing
*/
func (c TestPlans) View(project string, id int) revel.Result {
	var testplan models.TestPlan
	r := c.Tx.Preload("Creator").Preload("Executor").Where("id= ?", id).Find(&testplan)

	if r.Error != nil {
		revel.ERROR.Println("An error while find user")
		c.Response.Status = 500
		return c.Render()
	}

	// convert "1,2" like string to data used in JSTree
	arr := strings.Split(testplan.ExecuteCases, ",")
	var cases []models.TestCase
	r = c.Tx.Where("id in (?)", arr).Find(&cases)

	if r.Error != nil {
		revel.ERROR.Println("An error while finding testcases SELECT operation")
	}

	var build models.BuildItem
	if testplan.TargetBuildId != 0 {
		r = c.Tx.Where("id = ?", testplan.TargetBuildId).First(&build)

		if r.Error != nil {
			revel.ERROR.Println("An error while finding Build information in TestPlans", r.Error)
		}
	}

	return c.Render(project, testplan, cases, build)
}

/**
POST handler for DELETE operation for testplan
*/
func (c TestPlans) Delete(project string) revel.Result {

	//TODO fill
	return c.Render(project)
}

func (c TestPlans) Fire(project string, id int) revel.Result {
	var user *models.User
	if user = c.connected(); user == nil {
		c.Flash.Error("Please log in first")
		return c.Redirect(routes.Application.Index())
	}

	var prj models.Project
	r := c.Tx.Where("name = ?", project).First(&prj)

	if r.Error != nil {
		revel.ERROR.Println("An error while select project operation in TestPlans.Fire", r.Error)
		c.Response.Status = 500
		return c.Render()
	}

	var plan models.TestPlan
	r = c.Tx.Where("id = ?", id).First(&plan)

	if r.Error != nil {
		revel.ERROR.Println("An error while select testplan in TestPlans.Fire", r.Error)
		c.Response.Status = 500
		return c.Render()
	}

	// mendatory : testplan ID,project ID
	// advisory : if executor id is zero, make status N/A
	exec := models.Execution{ProjectId: prj.ID, PlanId: id, ExecutorId: plan.ExecutorId, TargetBuildId: plan.TargetBuildId}

	if plan.ExecutorId == 0 {
		exec.Status = 2
	}

	c.Tx.NewRecord(exec)
	r = c.Tx.Create(&exec)
	if r.Error != nil {
		revel.ERROR.Println("An error while create record in TestPlans.Fire")
		c.Response.Status = 500
		return c.Render()
	}

	// well...
	return c.Redirect(routes.TestExecs.Index(project))
}

/*
	Seperate Testcases ID joined by ","
*/
func handleSelected(content string, buf *string) int {
	arr := strings.Split(content, ",")
	var case_arr []string
	cnt := 0
	for _, s := range arr {
		if strings.HasPrefix(s, "cb") {
			cnt++
			k := strings.TrimPrefix(s, "cb")
			case_arr = append(case_arr, k)
		}
	}

	*buf = strings.Join(case_arr[:], ",")
	
	return cnt
}
