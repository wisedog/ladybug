package controllers

import (
	"github.com/revel/revel"
	"github.com/wisedog/ladybug/app/models"
	"github.com/wisedog/ladybug/app/routes"
)

type Builds struct {
	Application
}

func (c Builds) checkUser() revel.Result {
	if user := c.connected(); user == nil {
		c.Flash.Error("Please log in first")
		return c.Redirect(routes.Application.Index())
	}
	return nil
}

/*
 A page to show testcase's information
*/
func (c Builds) Index(project string) revel.Result {
	var builds []models.Build

	r := c.Tx.Find(&builds)
	if r.Error != nil {
	}

	return c.Render(project, builds)
}

/*
Render a page to add
*/
func (c Builds) Add(project string) revel.Result {
	var build models.Build

	return c.Render(project, build)
}

/**
POST handler for save build
*/
func (c Builds) Save(project string, build models.Build) revel.Result {
	if user := c.connected(); user == nil {
		c.Flash.Error("Please log in first")
		return c.Redirect(routes.Application.Index())
	}

	build.Validate(c.Validation)

	if c.Validation.HasErrors() {
		c.Flash.Error("invalid!")
		c.Validation.Keep()
		c.FlashParams()

		return c.Redirect(routes.Builds.Add(project))
	}
	var prj models.Project
	r := c.Tx.Where("name = ?", project).First(&prj)
	if r.Error != nil {
		c.Response.Status = 500
		return c.Render()
	}

	build.Project_id = prj.ID

	c.Tx.NewRecord(build)
	r = c.Tx.Create(&build)

	if r.Error != nil {
		revel.ERROR.Println("Insert operation failed on Builds.Save")
	}

	return c.Redirect(routes.Builds.Index(project))

}

/**
Render a page to view
*/
func (c Builds) View(project string, id int) revel.Result {
	var prj models.Project
	r := c.Tx.Where("name = ?", project).First(&prj)
	if r.Error != nil {
		c.Response.Status = 500
		return c.Render()
	}
	var build models.Build
	r = c.Tx.Where("id = ?", id).First(&build)

	if r.Error != nil {
		c.Response.Status = 500
		return c.Render()
	}

	return c.Render(project, build)
}

/**
Render a page to edit
*/
func (c Builds) Edit(project string, id int) revel.Result {

	var build models.Build
	r := c.Tx.Where("id = ?", id).First(&build)

	if r.Error != nil {
		c.Response.Status = 500
		return c.Render()
	}

	return c.Render(project, build)
}
