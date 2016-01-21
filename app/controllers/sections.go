package controllers

import (
	"github.com/revel/revel"
	"github.com/wisedog/ladybug/app/models"
	"github.com/wisedog/ladybug/app/routes"
)

type Sections struct {
	Application
}

//TODO template first
func (c Sections) Save(project string, section models.Section) revel.Result {
	if user := c.connected(); user == nil {
		c.Flash.Error("Please log in first")
		return c.Redirect(routes.Application.Index())
	}

	section.Validate(c.Validation)

	if c.Validation.HasErrors() {
		c.Flash.Error("invalid!")
		c.Validation.Keep()
		c.FlashParams()

		return c.Render(project)
		//return c.Redirect(routes.Sections.Add(project, section.ParentsID))
	}

	revel.INFO.Println("Section : ", section)

	c.Tx.NewRecord(section)
	r := c.Tx.Create(&section)

	if r.Error != nil {
		revel.ERROR.Println("An error while insert operation Sections.Save", r.Error)
	}

	return c.Render(project) //c.Redirect(routes.TestDesign.Index(project))

}

func (c Sections) Insert(project string, id int, parent_id int, title string, edit bool) revel.Result {

	type Reply struct {
		Status int    `json:"status"`
		Msg    string `json:"msg"`
	}

	var s models.Section
	var e Reply
	if user := c.connected(); user == nil {
		e = Reply{Status: 500, Msg: "Please log in first"}
		return c.RenderJson(e)
	}

	var prj models.Project
	r := c.Tx.Where("name = ?", project).First(&prj)
	if r.Error != nil {
		e = Reply{Status: 500, Msg: "Invalid project"}
		return c.RenderJson(e)
	}

	if parent_id == 0 {
		s = models.Section{Prefix: "Temp", Title: title, RootNode: true, ProjectID: prj.ID}
	} else {
		s = models.Section{Prefix: "Temp", Title: title, RootNode: false, ParentsID: parent_id, ProjectID: prj.ID}
	}

	c.Tx.NewRecord(s)
	r = c.Tx.Create(&s)

	if r.Error != nil {
		revel.ERROR.Println("An error while Insert opreation in Sections.Insert", r.Error)
		e = Reply{Status: 500, Msg: "An error while Insert operation in Sections.Insert"}
	} else {
		e = Reply{Status: 200, Msg: "OK"}
	}

	return c.RenderJson(e)
}
