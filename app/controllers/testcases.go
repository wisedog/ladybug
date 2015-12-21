package controllers

import (
	"github.com/revel/revel"
	"github.com/wisedog/ladybug/app/models"
	"github.com/wisedog/ladybug/app/routes"
)

type TestCases struct {
	Application
}

func (c TestCases) checkUser() revel.Result {
	if user := c.connected(); user == nil {
		c.Flash.Error("Please log in first")
		return c.Redirect(routes.Application.Index())
	}
	return nil
}

/*
 A page to show testcase's information
*/
func (c TestCases) Index(project string, id int) revel.Result {
	var tc models.TestCase
	r := c.Tx.Where("id = ?", id).First(&tc)

	if r.Error != nil {
		revel.ERROR.Println("Error on Select-ID query in TestCases.Index", r.Error)
		c.Response.Status = 500
		return c.Render()
	}

	return c.Render(project, tc)

}

/*
	A POST handler from Add
*/
func (c TestCases) Save(project string, testcase models.TestCase) revel.Result {

	if user := c.connected(); user == nil {
		c.Flash.Error("Please log in first")
		return c.Redirect(routes.Application.Index())
	}

	//Validate input testcase
	testcase.Validate(c.Validation)

	if c.Validation.HasErrors() {
		c.Flash.Error("invalid!")
		c.Validation.Keep()
		c.FlashParams()

		return c.Redirect(routes.TestCases.Add(project, testcase.SectionID))
	}

	var largest_seq_tc models.TestCase
	r := c.Tx.Where("prefix=?", testcase.Prefix).Order("seq desc").First(&largest_seq_tc)

	if r.Error != nil {
		revel.ERROR.Println("An error while SELECT operation to find largest seq")
	} else {
		testcase.Seq = largest_seq_tc.Seq + 1
		revel.INFO.Println("Largest number is : ", largest_seq_tc.Seq+1)
	}

	c.Tx.NewRecord(testcase)
	r = c.Tx.Create(&testcase)

	if r.Error != nil {
		revel.ERROR.Println("Insert operation failed in TestCase.Save")
	}

	return c.Redirect(routes.TestDesign.Index(project))
}

/*
	Render a page to add a testcase
*/
func (c TestCases) Add(project string, section_id int) revel.Result {
	var testcase models.TestCase
	var section models.Section

	// validate test suite
	r := c.Tx.Where("id = ?", section_id).First(&section)

	if r.Error != nil {
		revel.ERROR.Println("An error while SELECT testsuite operation in TestCases.Add. section_id is ", section_id)
		c.Response.Status = 500
		return c.Render()
	}

	testcase.SectionID = section_id
	testcase.Prefix = section.Prefix

	return c.Render(testcase, project)
}

/**
A POST handler for delete testcase
*/
func (c TestCases) Delete(project string, id int) revel.Result {

	// Delete the testcase  permanently for sequence
	r := c.Tx.Unscoped().Where("id = ?", id).Delete(&models.TestCase{})
	if r.Error != nil {
		revel.ERROR.Println("An error while delete testcase ", r.Error)
		c.Response.Status = 500
		return c.Render()
	}

	return c.Redirect(routes.TestDesign.Index(project))
}

/*
	A handler for EDIT GET and render edit page
*/
func (c TestCases) Edit(project string, id int) revel.Result {
	testcase := models.TestCase{}

	r := c.Tx.Where("id = ?", id).First(&testcase)

	if r.Error != nil {
		revel.ERROR.Println("An Error while SELECT operation for TestCase.Edit", r.Error)
		c.Response.Status = 500
		return c.Render(routes.TestDesign.Index(project))
	}

	return c.Render(project, id, testcase)
}

/*
 Update POST handler for Testcase Edit
*/
func (c TestCases) Update(project string, id int, testcase models.TestCase) revel.Result {
	c.Validation.Required(testcase)

	if c.Validation.HasErrors() {
		c.Flash.Error("invalid!")
		c.Validation.Keep()
		c.FlashParams()

		return c.Redirect(TestCases.Edit)
	}

	revel.INFO.Println("In TC update : ", testcase)

	exist_case := models.TestCase{}
	r := c.Tx.Where("id = ?", testcase.ID).First(&exist_case)

	if r.Error != nil {
		revel.ERROR.Println("An error while select exist testcase operation", r.Error)
		c.Flash.Error("Invalid!")
		c.Response.Status = 400
		return c.Redirect(TestCases.Edit)

	}

	exist_case.Title = testcase.Title
	exist_case.Description = testcase.Description
	exist_case.Seq = testcase.Seq
	exist_case.Status = testcase.Status
	exist_case.Prefix = testcase.Prefix
	exist_case.Precondition = testcase.Precondition
	exist_case.Steps = testcase.Steps
	exist_case.Expected = testcase.Expected

	r = c.Tx.Save(&exist_case)

	if r.Error != nil {
		revel.ERROR.Println("An error while SAVE operation on TestCases.Update")
		c.Response.Status = 500
		return c.Render()
	}
	
	c.Flash.Success("Update Success!")


	return c.Redirect(routes.TestDesign.Index(project))
}
