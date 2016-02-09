package controllers

import (
	"strconv"
	"fmt"
	"encoding/json"
	
	"github.com/revel/revel"
	"github.com/wisedog/ladybug/app/models"
	"github.com/wisedog/ladybug/app/routes"
)

// TestCases is inherited from Application. 
// See app.go
type TestCases struct {
	Application
}

// checkUser is private utilitiy function for checking 
// this user id now on connected or not.
func (c TestCases) checkUser() revel.Result {
	if user := c.connected(); user == nil {
		c.Flash.Error("Please log in first")
		return c.Redirect(routes.Application.Index())
	}
	return nil
}

// makeMessage is private utility function for comparing and make
// messages for status changing
func (c TestCases) makeMessage(historyUnit *[]models.HistoryTestCaseUnit){
	if historyUnit == nil{
		revel.ERROR.Println("Nil historyunit!")
		return
	}
	
	var msg string
	
	// TODO L10N with using c.Message format string
	for i := 0; i < len(*historyUnit); i++{
		if (*historyUnit)[i].ChangeType == models.HISTORY_CHANGE_TYPE_CHANGED{
			if((*historyUnit)[i].From == 0 && (*historyUnit)[i].To == 0 ){
				msg = fmt.Sprintf(`"%s" is changed from "%s" to "%s".`, 
				(*historyUnit)[i].What, (*historyUnit)[i].FromStr, (*historyUnit)[i].ToStr)
			}else{
				msg = fmt.Sprintf(`"%s" is changed from %d to %d.`, 
				(*historyUnit)[i].What, (*historyUnit)[i].From, (*historyUnit)[i].To)
			}
		}else if(*historyUnit)[i].ChangeType == models.HISTORY_CHANGE_TYPE_SET{
			msg = fmt.Sprintf(`"%s" is set to "%s".`, 
				(*historyUnit)[i].What, (*historyUnit)[i].Set)
		}else if(*historyUnit)[i].ChangeType == models.HISTORY_CHANGE_TYPE_DIFF{
			msg = fmt.Sprintf(`"%s" is changed(diff).`, 
				(*historyUnit)[i].What)
		}else if(*historyUnit)[i].ChangeType == models.HISTORY_CHANGE_TYPE_NOTE{
			revel.INFO.Println("INFO : ", (*historyUnit)[i])
			msg = fmt.Sprintf("%s added a note.", (*historyUnit)[i].What)
		}else{
			msg = ""
		}
		(*historyUnit)[i].Msg = msg
	}
}


// CaseIndex renders a page to show testcase's information
func (c TestCases) CaseIndex(project string, id int) revel.Result {
	var tc models.TestCase
	r := c.Tx.Where("id = ?", id).First(&tc)

	if r.Error != nil {
		revel.ERROR.Println("Error on Select-ID query in TestCases.Index", r.Error)
		c.Response.Status = 500
		return c.Render()
	}
	
	
	tc.PriorityStr = c.getPriorityL10N(tc.Priority)
	
	var histories []models.History
	c.Tx.Where("category = ?", models.HISTORY_TYPE_TC).Where("target_id = ?", tc.ID).Preload("User").Find(&histories)
	
	
	// making status changement messages
	for i := 0; i < len(histories); i++{
		var res []models.HistoryTestCaseUnit
		json.Unmarshal([]byte(histories[i].ChangesJson), &res)
		
		// make message
		c.makeMessage(&res)
		
		histories[i].Changes = res
		revel.INFO.Println("res", histories[i].Changes)
	}
	

	return c.Render(project, tc, histories)

}

//	Save is a POST handler from Add page.
func (c TestCases) Save(project string, testcase models.TestCase, reviewerID int) revel.Result {

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
	
	display_id := testcase.Prefix + "-" + strconv.Itoa(testcase.ID)
	testcase.DisplayID = display_id
	r = c.Tx.Save(&testcase)
	
	if r.Error != nil{
		revel.ERROR.Println("Update operation failed in TestCase.Save", r.Error)
	}

	return c.Redirect(routes.TestDesign.DesignIndex(project))
}

//	Add just renders a page too add a testcase
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
	
	var category []models.Category
	c.Tx.Find(&category)
	
	var reviewerID int

	return c.Render(testcase, project, category, reviewerID)
}

// Delete is a POST handler for DELETE request
func (c TestCases) Delete(project string, id int) revel.Result {

	// Delete the testcase  permanently for sequence
	r := c.Tx.Unscoped().Where("id = ?", id).Delete(&models.TestCase{})
	if r.Error != nil {
		revel.ERROR.Println("An error while delete testcase ", r.Error)
		c.Response.Status = 500
		return c.Render()
	}

	return c.Redirect(routes.TestDesign.DesignIndex(project))
}

// Edit just renders a EDIT page.
func (c TestCases) Edit(project string, id int) revel.Result {
	
	var category []models.Category
	c.Tx.Find(&category)
	
	if c.Validation.HasErrors(){
		//c.Flash.Success("Update Success!")
		c.Validation.Keep()
		c.FlashParams()
	
		return c.Render(project, id, category)
	}
	
	testcase := models.TestCase{}

	r := c.Tx.Where("id = ?", id).First(&testcase)

	if r.Error != nil {
		revel.ERROR.Println("An Error while SELECT operation for TestCase.Edit", r.Error)
		c.Response.Status = 500
		return c.Render(routes.TestDesign.DesignIndex(project))
	}
	
	flash := map[string]string{
	    "testcase.Priority": strconv.Itoa(testcase.Priority),
	    "testcase.Category": strconv.Itoa(testcase.CategoryID),
	    "testcase.Status" : strconv.Itoa(testcase.Status),
	}
	
	//TODO change section here.
	
	var note string

	return c.Render(project, id, testcase, category, flash, note)
}

/*
 Update POST handler for Testcase Edit
*/
func (c TestCases) Update(project string, id int, testcase models.TestCase, note string) revel.Result {
	//Validate input testcase
	testcase.Validate(c.Validation)

	if c.Validation.HasErrors() {
		c.Flash.Error("invalid!")
		c.Validation.Keep()
		c.FlashParams()

		return c.Redirect(routes.TestCases.Edit(project, id))
	}
	
	existCase := models.TestCase{}
	r := c.Tx.Where("id = ?", testcase.ID).First(&existCase)

	if r.Error != nil {
		revel.ERROR.Println("An error while select exist testcase operation", r.Error)
		c.Flash.Error("Invalid!")
		c.Response.Status = 400
		return c.Redirect(TestCases.Edit)

	}
	
	c.findDiff(&existCase, &testcase, note)
	
	existCase.Title = testcase.Title
	existCase.Description = testcase.Description
	existCase.Seq = testcase.Seq
	existCase.Status = testcase.Status
	existCase.Prefix = testcase.Prefix
	existCase.Precondition = testcase.Precondition
	existCase.Steps = testcase.Steps
	existCase.Expected = testcase.Expected
	existCase.Priority = testcase.Priority

	r = c.Tx.Save(&existCase)

	if r.Error != nil {
		revel.ERROR.Println("An error while SAVE operation on TestCases.Update")
		c.Response.Status = 500
		return c.Render()
	}
	
	c.Flash.Success("Update Success!")


	return c.Redirect(routes.TestDesign.DesignIndex(project))
}

// findDiff compares between two models.TestCase and create 
// HistoryTestCaseUnit to database. Used in Update
func (c TestCases) findDiff(existCase, newCase *models.TestCase, note string){
	user := c.connected()
	if user == nil {
		c.Flash.Error("Please log in first")
		return
	}
	
	var changes []models.HistoryTestCaseUnit
	his := models.History{Category : models.HISTORY_TYPE_TC,
			TargetID : existCase.ID, UserID : user.ID,
	}
	
	// check note
	if note != ""{
		// this means user add a note on this testcase
		unit := models.HistoryTestCaseUnit{
			ChangeType : models.HISTORY_CHANGE_TYPE_NOTE, What : user.Name,
		}
		
		his.Note = note
		changes = append(changes, unit)
	}
	
	// check title
	if existCase.Title != newCase.Title {
		unit := models.HistoryTestCaseUnit{
			ChangeType : models.HISTORY_CHANGE_TYPE_CHANGED, What : "Title",
			FromStr : existCase.Title, ToStr : newCase.Title,
		}
		changes = append(changes, unit)
	}
	
	// check priority
	if existCase.Priority != newCase.Priority {
		unit := models.HistoryTestCaseUnit{
			ChangeType : models.HISTORY_CHANGE_TYPE_CHANGED, 
			What : c.Message("testcase.priority"),
			FromStr : c.getStatusL10N(existCase.Priority),
			ToStr : c.getStatusL10N(newCase.Priority),
		}
		changes = append(changes, unit)
	}
	
	// check execution type
	if existCase.ExecutionType != newCase.ExecutionType {
		arr := [2]string{"Manual", "Automated"}
		unit := models.HistoryTestCaseUnit{
			ChangeType : models.HISTORY_CHANGE_TYPE_CHANGED, What : "Execution type",
			FromStr : arr[existCase.ExecutionType], ToStr : arr[newCase.ExecutionType],
		}
		changes = append(changes, unit)
	}
	
	// TODO check Status
	
	// check Description
	if existCase.Description != newCase.Description {
		unit := models.HistoryTestCaseUnit{
			ChangeType : models.HISTORY_CHANGE_TYPE_DIFF, 
			What : c.Message("description"),
			DiffID : 2,
		}
		changes = append(changes, unit)
	}
	
	// check Precondition
	if existCase.Precondition != newCase.Precondition {
		unit := models.HistoryTestCaseUnit{
			ChangeType : models.HISTORY_CHANGE_TYPE_DIFF, 
			What : c.Message("priority.precondition"),
			DiffID : 2,//TODO should be implemnted DIFF
		}
		changes = append(changes, unit)
	}
	
	// check Estimated
	if existCase.Estimated != newCase.Estimated {
		unit := models.HistoryTestCaseUnit{
			ChangeType : models.HISTORY_CHANGE_TYPE_DIFF, What : "Estimated Time",
			DiffID : 2,	//TODO should be implemnted DIFF
		}
		changes = append(changes, unit)
	}
	
	// check Steps
	if existCase.Steps != newCase.Steps {
		unit := models.HistoryTestCaseUnit{
			ChangeType : models.HISTORY_CHANGE_TYPE_DIFF, What : "Steps",
			DiffID : 2,	//TODO should be implemnted DIFF
		}
		changes = append(changes, unit)
	}
	
	// check Expected
	if existCase.Expected != newCase.Expected {
		unit := models.HistoryTestCaseUnit{
			ChangeType : models.HISTORY_CHANGE_TYPE_DIFF, 
			What : c.Message("priority.expected"),
			DiffID : 2,	//TODO should be implemnted DIFF
		}
		changes = append(changes, unit)
	}
	
	var existCategory models.Category
	var testcaseCategory models.Category
	
	c.Tx.Where("id = ?", existCase.CategoryID).Find(&existCategory)
	c.Tx.Where("id = ?", newCase.CategoryID).Find(&testcaseCategory)
	
	// check CategoryID
	if existCase.CategoryID != newCase.CategoryID {
		unit := models.HistoryTestCaseUnit{
			ChangeType : models.HISTORY_CHANGE_TYPE_CHANGED, 
			What : c.Message("priority.category"),
			FromStr : "",
			ToStr : "",
		}
		changes = append(changes, unit)
	}
	
	result, _ := json.Marshal(changes)
	his.ChangesJson = string(result)
	
	c.Tx.NewRecord(his)
	c.Tx.Create(&his)
}


// getPriorityL10N function returns localization string 
func (c TestCases) getPriorityL10N(priority int) string{
	s := ""
	switch priority{
		case 1: s = c.Message("priority.highest")
		case 2: s = c.Message("priority.high")
		case 3: s = c.Message("priority.medium")
		case 4: s = c.Message("priority.low")
		case 5: s = c.Message("priority.lowest")
		default : s = "Unknown"
	}
	
	return s
}