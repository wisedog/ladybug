package controllers

import (
	"fmt"
	"strconv"

	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/wisedog/ladybug/errors"
	"github.com/wisedog/ladybug/interfacer"
	"github.com/wisedog/ladybug/models"
	log "gopkg.in/inconshreveable/log15.v2"
)

// makeMessage is private utility function for comparing and make
// messages for status changing
//func  makeMessage(historyUnit *[]models.HistoryTestCaseUnit){
func makeHistoryMessage(historyUnit *[]models.HistoryTestCaseUnit) {

	if historyUnit == nil {
		log.Error("Nil historyunit!")
		return
	}

	var msg string

	// TODO L10N with using GetI18nMessage format string
	for i := 0; i < len(*historyUnit); i++ {
		if (*historyUnit)[i].ChangeType == models.HistoryChangeTypeChanged {
			if (*historyUnit)[i].From == 0 && (*historyUnit)[i].To == 0 {
				msg = fmt.Sprintf(`"%s" is changed from "%s" to "%s".`,
					(*historyUnit)[i].What, (*historyUnit)[i].FromStr, (*historyUnit)[i].ToStr)
			} else {
				msg = fmt.Sprintf(`"%s" is changed from %d to %d.`,
					(*historyUnit)[i].What, (*historyUnit)[i].From, (*historyUnit)[i].To)
			}
		} else if (*historyUnit)[i].ChangeType == models.HistoryChangeTypeSet {
			msg = fmt.Sprintf(`"%s" is set to "%s".`,
				(*historyUnit)[i].What, (*historyUnit)[i].Set)
		} else if (*historyUnit)[i].ChangeType == models.HistoryChangeTypeDiff {
			msg = fmt.Sprintf(`"%s" is changed(diff).`,
				(*historyUnit)[i].What)
		} else if (*historyUnit)[i].ChangeType == models.HistoryChangeTypeNote {
			msg = fmt.Sprintf("%s added a note.", (*historyUnit)[i].What)
		} else {
			msg = ""
		}
		(*historyUnit)[i].Msg = msg
	}
}

// CaseView renders a page to show testcase's information
func CaseView(c *interfacer.AppContext, w http.ResponseWriter, r *http.Request) error {

	// get ID from /view/{id} request
	vars := mux.Vars(r)
	id := vars["id"]

	var tc models.TestCase
	if err := c.Db.Where("id = ?", id).Preload("Category").First(&tc); err.Error != nil {
		return LogAndHTTPError(http.StatusInternalServerError, "testcases", "db", err.Error.Error())
	}

	tc.PriorityStr = getPriorityI18n(tc.Priority)

	var histories []models.History
	c.Db.Where("category = ?", models.HISTORY_TYPE_TC).Where("target_id = ?", tc.ID).Preload("User").Find(&histories)

	// making status changement messages
	for i := 0; i < len(histories); i++ {
		var res []models.HistoryTestCaseUnit
		json.Unmarshal([]byte(histories[i].ChangesJson), &res)

		// make message
		makeHistoryMessage(&res)

		histories[i].Changes = res
		log.Debug("testcases", "res", histories[i].Changes)
	}

	// Find related Test Cases with this requirement
	var reqs []models.Requirement
	if err := c.Db.Model(&tc).Association("RelatedRequirements").Find(&reqs); err.Error != nil {
		log.Error("TestCase", "type", "database", "msg ", err.Error)
		return errors.HttpError{Status: http.StatusBadRequest, Desc: "Not found related testcases"}
	}

	items := map[string]interface{}{
		"TestCase":    tc,
		"Histories":   histories,
		"RelatedReqs": reqs,
		"Active_idx":  2,
	}

	return Render2(c, w, items, "views/base.tmpl", "views/testcases/caseindex.tmpl")
}

// renderAddEdit renders ADD and EDIT pages
func renderAddEdit(c *interfacer.AppContext, w http.ResponseWriter, r *http.Request, isEdit bool) error {

	// find [add|save]/{id}
	vars := mux.Vars(r)
	id := vars["id"]

	// valid when rendering Add
	sectionID := r.URL.Query().Get("sectionid")
	if sectionID == "" && isEdit == false {
		return LogAndHTTPError(http.StatusBadRequest, "testcases", "http",
			"invalid condition. Add rendering requires section id")
	}

	var category []models.Category
	c.Db.Find(&category)

	var testcase models.TestCase
	if isEdit {
		if err := c.Db.Where("id = ?", id).First(&testcase); err.Error != nil {
			log.Error("testcase", "type", "database", "msg", err.Error)
			return errors.HttpError{Status: http.StatusInternalServerError,
				Desc: "An Error while SELECT operation for TestCase.Edit"}
		}
	} else {
		var prj models.Project
		if err := c.Db.Where("name = ?", c.ProjectName).First(&prj); err.Error != nil {
			log.Error("testcase", "type", "database", "msg", err.Error)
			return errors.HttpError{Status: http.StatusInternalServerError,
				Desc: "An Error while SELECT operation for TestCase.Edit"}
		}

		testcase.SectionID, _ = strconv.Atoi(sectionID)
		testcase.Prefix = prj.Prefix
	}

	//TODO change section here.

	items := map[string]interface{}{
		"ID":         id,
		"SectionID":  sectionID,
		"TestCase":   testcase,
		"Category":   category,
		"IsEdit":     isEdit,
		"Active_idx": 2,
	}

	return Render2(c, w, items, "views/base.tmpl", "views/testcases/caseaddedit.tmpl")
}

// CaseEdit just renders a EDIT page. Add and Edit pages are integrated. see renderAddEdit
func CaseEdit(c *interfacer.AppContext, w http.ResponseWriter, r *http.Request) error {
	return renderAddEdit(c, w, r, true)
}

// CaseAdd just renders a Add page. Add and Edit pages are integrated. see renderAddEdit
func CaseAdd(c *interfacer.AppContext, w http.ResponseWriter, r *http.Request) error {

	return renderAddEdit(c, w, r, false)
}

// handleSaveUpdate handles POST request from save and update
func handleSaveUpdate(c *interfacer.AppContext, w http.ResponseWriter, r *http.Request, isUpdate bool) error {

	var testcase models.TestCase
	vars := mux.Vars(r)
	projectName := vars["projectName"]
	id := vars["id"]

	if err := r.ParseForm(); err != nil {
		log.Error("Testcase", "Error ", err)
		return errors.HttpError{Status: http.StatusInternalServerError, Desc: "ParseForm failed"}
	}

	if err := c.Decoder.Decode(&testcase, r.PostForm); err != nil {
		log.Warn("Testcase", "Error", err, "msg", "Decode failed but go ahead")
	}

	// TODO : Validate input testcase
	//testcase.Validate()

	redirectionTarget := fmt.Sprintf("/project/%s/design", projectName)

	if isUpdate == false {

		var largestSeqTc models.TestCase

		if err := c.Db.Where("prefix=?", testcase.Prefix).Order("seq desc").First(&largestSeqTc); err.Error != nil {
			log.Error("An error while SELECT operation to find largest seq")
		} else {
			testcase.Seq = largestSeqTc.Seq + 1
			log.Debug("testcase", "Largest number is : ", largestSeqTc.Seq+1)
		}

		c.Db.NewRecord(testcase)

		if err := c.Db.Create(&testcase); err.Error != nil {
			log.Error("testcase", "type", "database", "msg", err.Error)
			return errors.HttpError{Status: http.StatusInternalServerError,
				Desc: "Insert operation failed in TestCase.Save"}
		}

		displayID := testcase.Prefix + "-" + strconv.Itoa(testcase.ID)
		testcase.DisplayID = displayID

		if err := c.Db.Save(&testcase); err.Error != nil {
			log.Error("testcase", "type", "database", "msg", err.Error)
			return errors.HttpError{Status: http.StatusInternalServerError,
				Desc: "Update operation failed in TestCase.Save"}
		}

		http.Redirect(w, r, redirectionTarget, http.StatusFound)
	}

	var existCase models.TestCase

	if err := c.Db.Where("id = ?", id).First(&existCase); err.Error != nil {
		log.Error("testcase", "type", "database", "msg", err.Error)
		return errors.HttpError{Status: http.StatusInternalServerError,
			Desc: "An error while select exist testcase operation"}
	}
	note := r.FormValue("Note")
	findDiff(c, &existCase, &testcase, note, c.User)

	existCase.Title = testcase.Title
	existCase.Description = testcase.Description
	existCase.Seq = testcase.Seq
	existCase.Status = testcase.Status
	existCase.Prefix = testcase.Prefix
	existCase.Precondition = testcase.Precondition
	existCase.Steps = testcase.Steps
	existCase.Expected = testcase.Expected
	existCase.Priority = testcase.Priority

	if err := c.Db.Save(&existCase); err.Error != nil {
		log.Error("testcase", "type", "database", "msg", err.Error)
		return errors.HttpError{Status: http.StatusInternalServerError,
			Desc: "An error while SAVE operation on TestCases.Update"}
	}

	//c.Flash.Success("Update Success!")

	http.Redirect(w, r, redirectionTarget, http.StatusFound)
	return nil
}

// CaseSave handles POST request from CaseAdd. Save and Update handlers are integrated. see handleSaveUpdate
func CaseSave(c *interfacer.AppContext, w http.ResponseWriter, r *http.Request) error {
	return handleSaveUpdate(c, w, r, false)
}

// CaseUpdate POST handler for Testcase Edit
func CaseUpdate(c *interfacer.AppContext, w http.ResponseWriter, r *http.Request) error {
	return handleSaveUpdate(c, w, r, true)
}

// CaseDelete is a POST handler for DELETE request
func CaseDelete(c *interfacer.AppContext, w http.ResponseWriter, r *http.Request) error {
	var user *models.User
	if user = connected(c, r); user == nil {
		log.Debug("Not found login information")
		http.Redirect(w, r, "/", http.StatusFound)
	}

	if err := r.ParseForm(); err != nil {
		log.Error("Testcase", "Error ", err)
		return errors.HttpError{Status: http.StatusInternalServerError, Desc: "ParseForm failed"}
	}

	id := r.FormValue("id")

	// Delete the testcase  permanently for sequence
	if err := c.Db.Unscoped().Where("id = ?", id).Delete(&models.TestCase{}); err.Error != nil {
		log.Error("testcase", "An error while delete testcase ", err.Error)
		return errors.HttpError{Status: http.StatusInternalServerError, Desc: "testcase delete failed"}
	}

	// the client do redirect or refresh.
	return nil
}

// findDiff compares between two models.TestCase and create
// HistoryTestCaseUnit to database. Used in Update
func findDiff(c *interfacer.AppContext, existCase, newCase *models.TestCase, note string, user *models.User) {

	var changes []models.HistoryTestCaseUnit
	his := models.History{Category: models.HISTORY_TYPE_TC,
		TargetID: existCase.ID, UserID: user.ID,
	}

	// check note
	if note != "" {
		// this means user add a note on this testcase
		unit := models.HistoryTestCaseUnit{
			ChangeType: models.HistoryChangeTypeNote, What: user.Name,
		}

		his.Note = note
		changes = append(changes, unit)
	}

	// check title
	if existCase.Title != newCase.Title {
		unit := models.HistoryTestCaseUnit{
			ChangeType: models.HistoryChangeTypeChanged, What: "Title",
			FromStr: existCase.Title, ToStr: newCase.Title,
		}
		changes = append(changes, unit)
	}

	// check priority
	if existCase.Priority != newCase.Priority {
		unit := models.HistoryTestCaseUnit{
			ChangeType: models.HistoryChangeTypeChanged,
			What:       GetI18nMessage("testcase.priority"),
			FromStr:    getPriorityI18n(existCase.Priority),
			ToStr:      getPriorityI18n(newCase.Priority),
		}
		changes = append(changes, unit)
	}

	// check execution type
	if existCase.ExecutionType != newCase.ExecutionType {
		arr := [2]string{"Manual", "Automated"}
		unit := models.HistoryTestCaseUnit{
			ChangeType: models.HistoryChangeTypeChanged, What: "Execution type",
			FromStr: arr[existCase.ExecutionType], ToStr: arr[newCase.ExecutionType],
		}
		changes = append(changes, unit)
	}

	// TODO check Status

	// check Description
	if existCase.Description != newCase.Description {
		unit := models.HistoryTestCaseUnit{
			ChangeType: models.HistoryChangeTypeDiff,
			What:       GetI18nMessage("description"),
			DiffID:     2,
		}
		changes = append(changes, unit)
	}

	// check Precondition
	if existCase.Precondition != newCase.Precondition {
		unit := models.HistoryTestCaseUnit{
			ChangeType: models.HistoryChangeTypeDiff,
			What:       GetI18nMessage("priority.precondition"),
			DiffID:     2, //TODO should be implemnted DIFF
		}
		changes = append(changes, unit)
	}

	// check Estimated
	if existCase.Estimated != newCase.Estimated {
		unit := models.HistoryTestCaseUnit{
			ChangeType: models.HistoryChangeTypeDiff, What: "Estimated Time",
			DiffID: 2, //TODO should be implemnted DIFF
		}
		changes = append(changes, unit)
	}

	// check Steps
	if existCase.Steps != newCase.Steps {
		// FIXME A bug : those strings are same, but,,,,
		unit := models.HistoryTestCaseUnit{
			ChangeType: models.HistoryChangeTypeDiff, What: "Steps",
			DiffID: 2, //TODO should be implemnted DIFF
		}
		changes = append(changes, unit)
	}

	// check Expected
	if existCase.Expected != newCase.Expected {
		unit := models.HistoryTestCaseUnit{
			ChangeType: models.HistoryChangeTypeDiff,
			What:       GetI18nMessage("priority.expected"),
			DiffID:     2, //TODO should be implemnted DIFF
		}
		changes = append(changes, unit)
	}

	var existCategory models.Category
	var testcaseCategory models.Category

	c.Db.Where("id = ?", existCase.CategoryID).Find(&existCategory)
	c.Db.Where("id = ?", newCase.CategoryID).Find(&testcaseCategory)

	// check CategoryID
	if existCase.CategoryID != newCase.CategoryID {
		unit := models.HistoryTestCaseUnit{
			ChangeType: models.HistoryChangeTypeChanged,
			What:       GetI18nMessage("priority.category"),
			FromStr:    "",
			ToStr:      "",
		}
		changes = append(changes, unit)
	}

	result, _ := json.Marshal(changes)
	his.ChangesJson = string(result)

	c.Db.NewRecord(his)
	c.Db.Create(&his)
}

// UnlinkRequirementRelation unlinks a relationship between Test Case and Requirement
func UnlinkRequirementRelation(c *interfacer.AppContext, w http.ResponseWriter, r *http.Request) error {
	// get requirement, testcase ID from post
	if err := r.ParseForm(); err != nil {
		log.Error("TestCase", "type", "app", "msg ", err.Error())
		return RenderJSONWithStatus(w, Resp{Msg: "Parse form is not valid"}, http.StatusBadRequest)
	}

	reqIDStr := r.FormValue("req_id")
	tcIDStr := r.FormValue("tc_id")

	reqID, err := strconv.Atoi(reqIDStr)
	if err != nil {
		log.Error("TestCase", "type", "app", "msg ", err.Error())
		return RenderJSONWithStatus(w, Resp{Msg: err.Error()}, http.StatusBadRequest)
	}

	tcID, err := strconv.Atoi(tcIDStr)
	if err != nil {
		log.Error("TestCase", "type", "app", "msg ", err.Error())
		return RenderJSONWithStatus(w, Resp{Msg: err.Error()}, http.StatusBadRequest)
	}

	var targetReq models.Requirement
	if err := c.Db.Where("id = ?", reqID).First(&targetReq); err.Error != nil {
		log.Error("TestCase", "type", "database", "msg ", err.Error.Error())
		return RenderJSONWithStatus(w,
			Resp{Msg: "Not found or something wrong while unlinking testcase and related requirement"},
			http.StatusInternalServerError)
	}

	var tc models.TestCase
	tc.ID = tcID

	// Find related Test Cases with this requirement
	//var reqs []models.Requirement
	if err := c.Db.Model(&tc).Association("RelatedRequirements").Delete(targetReq); err.Error != nil {
		log.Error("TestCase", "type", "database", "msg ", err.Error)
		return RenderJSONWithStatus(w, Resp{Msg: "Error while delete linking"}, http.StatusInternalServerError)
	}

	return RenderJSONWithStatus(w, Resp{Msg: "Success to unlink"}, http.StatusOK)
}

// LinkRequirementRelation links a relationship between Test Case and Requirement
func LinkRequirementRelation(c *interfacer.AppContext, w http.ResponseWriter, r *http.Request) error {

	return nil
}
