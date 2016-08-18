package controllers

import (
	"net/http"
	"strconv"

	"github.com/wisedog/ladybug/interfacer"
	"github.com/wisedog/ladybug/models"

	log "gopkg.in/inconshreveable/log15.v2"
)

// SectionEdit (Not Implemented) handles POST request from edit. Unlink the other EDIT functions,
// this function does not render templates, just JSON.
func SectionEdit(c *interfacer.AppContext, w http.ResponseWriter, r *http.Request) error {
	// get "section_id" from post form
	if err := r.ParseForm(); err != nil {
		return logAndRenderJSONWithStatus(w, Resp{Msg: "Parse form is not valid"}, http.StatusBadRequest,
			logTypeErr, "Section", "app", err.Error())
	}

	title := r.FormValue("title")
	if len(title) == 0 {
		return logAndRenderJSONWithStatus(w, Resp{Msg: "Empty title"}, http.StatusBadRequest,
			logTypeErr, "Section", "app", "Empty title")
	}

	// get selected node
	sectionIDStr := r.FormValue("parent_id")

	sectionID, err := strconv.Atoi(sectionIDStr)
	if err != nil {
		return logAndRenderJSONWithStatus(w, Resp{Msg: err.Error()}, http.StatusBadRequest,
			logTypeErr, "Section", "app", err.Error())
	}
	var sec models.Section
	if err := c.Db.Where("id = ?", sectionID).First(&sec).Error; err != nil {
		return logAndRenderJSONWithStatus(w, Resp{Msg: err.Error()}, http.StatusInternalServerError,
			logTypeErr, "Section", "Database", err.Error())
	}

	if err := c.Db.Model(&sec).Update("title", title).Error; err != nil {
		return logAndRenderJSONWithStatus(w, Resp{Msg: err.Error()}, http.StatusInternalServerError,
			logTypeErr, "Section", "Database", err.Error())
	}

	return RenderJSONWithStatus(w, Resp{Msg: "OK"}, http.StatusOK)
}

// SectionDelete handles POST request from delete.
// Hold on. This function is pretty long to read.
func SectionDelete(c *interfacer.AppContext, w http.ResponseWriter, r *http.Request) error {
	// get "section_id" from post form
	if err := r.ParseForm(); err != nil {
		return logAndRenderJSONWithStatus(w, Resp{Msg: "Parse form is not valid"}, http.StatusBadRequest,
			logTypeErr, "Section", "app", err.Error())
	}

	sectionIDStr := r.FormValue("section_id")
	isRequirementStr := r.FormValue("is_requirement")
	isRequirement := false
	if isRequirementStr == "true" {
		isRequirement = true
	}
	var targetModel interface{}
	if isRequirement {
		targetModel = &models.Requirement{}
	} else {
		targetModel = &models.TestCase{}
	}

	sectionID, err := strconv.Atoi(sectionIDStr)
	if err != nil {
		log.Error("Section", "type", "app", "msg ", err.Error())
		return RenderJSONWithStatus(w, Resp{Msg: err.Error()}, http.StatusBadRequest)
	}
	var prj models.Project
	if err := c.Db.Where("name = ?", c.ProjectName).First(&prj); err.Error != nil {
		log.Error("Section", "type", "app", "msg ", err.Error.Error(), "additional", "project is not found")
		return RenderJSONWithStatus(w, Resp{Msg: err.Error.Error()}, http.StatusInternalServerError)
	}

	// if Uncategorized section, delete it's decendant testcases.
	// There should not be decendant section
	var sec models.Section
	if err := c.Db.Where("id = ?", sectionID).First(&sec); err.Error != nil {
		log.Error("Section", "type", "app", "msg ", err.Error.Error())
		return RenderJSONWithStatus(w, Resp{Msg: err.Error.Error()}, http.StatusInternalServerError)
	}

	if sec.SpecialNode == true {
		//var testCasesToDelete []models.TestCase
		//var objectToDelete []interface{}
		table := "test_cases"
		if isRequirement {
			table = "requirements"
		}

		// data structure for get section ID
		type Result struct {
			ID int
		}

		var result []Result
		if err := c.Db.Table(table).Select("id").Where("section_id = ?", sectionID).Scan(&result); err.Error != nil {
			log.Error("Section", "type", "app", "msg ", err.Error.Error(), "additional", "not found target section")
			return RenderJSONWithStatus(w, Resp{Msg: err.Error.Error()}, http.StatusInternalServerError)
		}

		// Technical debut. I tried to use GORM's association model, failed I have..
		type testcasesReqs struct {
			RequirementID int
			TestCaseID    int
		}

		var objectToDeleteIDs []int

		for _, v := range result {
			objectToDeleteIDs = append(objectToDeleteIDs, v.ID)
		}

		var p []testcasesReqs
		var whereText string
		if isRequirement {
			whereText = "requirement_id in (?)"
		} else {
			whereText = "test_case_id in (?)"
		}

		if err := c.Db.Table("testcases_reqs").Select("requirement_id, test_case_id").
			Where(whereText, objectToDeleteIDs).Scan(&p); err.Error != nil {
			log.Error("Section", "type", "app", "msg", err.Error.Error(), "additional", "not found in join table")
		}

		var historyList []models.TcReqRelationHistory
		for _, v := range p {
			history := models.TcReqRelationHistory{
				Kind:          models.TcReqRelationHistoryUnlink,
				RequirementID: v.RequirementID,
				TestCaseID:    v.TestCaseID,
				ProjectID:     prj.ID,
			}
			historyList = append(historyList, history)
		}

		for _, v := range historyList {
			c.Db.NewRecord(v)
			c.Db.Create(&v)
		}

		if err := c.Db.Where("section_id = ?", sectionID).Delete(targetModel); err.Error != nil {
			log.Error("Section", "type", "app", "msg ", err.Error.Error())
			return RenderJSONWithStatus(w, Resp{Msg: err.Error.Error()}, http.StatusInternalServerError)
		}
	}

	// delete all sections belongs to the section and it's descendant,
	// move all testcases belongs to the section and it's descendant to 'Uncategorized'

	// first find all sections belongs to the section and it's descendant
	var targetSections []models.Section
	if err := c.Db.Where("parents_id = ?", sectionID).Find(&targetSections); err.Error != nil {
		log.Error("Section", "type", "app", "msg ", err.Error.Error())
		return RenderJSONWithStatus(w, Resp{Msg: err.Error.Error()}, http.StatusInternalServerError)
	}

	var targetIDs []int
	targetIDs = append(targetIDs, sectionID)

	for _, v := range targetSections {
		targetIDs = append(targetIDs, v.ID)
	}

	// section has only three-level depth
	var lastTargetSections []models.Section
	if err := c.Db.Where("parents_id in (?)", targetIDs).Find(&lastTargetSections); err.Error != nil {
		log.Error("Section", "type", "database", "msg ", err.Error.Error())
		return RenderJSONWithStatus(w, Resp{Msg: err.Error.Error()}, http.StatusInternalServerError)
	}

	for _, v := range lastTargetSections {
		targetIDs = append(targetIDs, v.ID)
	}

	// second, create 'Uncategorized' root section and special flag on
	// Be sure that it is already there, load it
	var uncategorizedSection models.Section
	if err := c.Db.Where("project_id = ?", prj.ID).Where("special_node = ?", true).
		Where("for_test_case = ?", !isRequirement).
		First(&uncategorizedSection); err.Error != nil {
		if err.RecordNotFound() {
			uncategorizedSection = models.Section{Seq: 1, Title: "Uncategorized", Status: 0, RootNode: true,
				Prefix: prj.Prefix, ProjectID: prj.ID,
				ForTestCase: !isRequirement, SpecialNode: true}
			c.Db.NewRecord(uncategorizedSection)
			c.Db.Create(&uncategorizedSection)
		} else {
			// really error
			log.Debug("WWWW", "asda", "asdasdsa")
		}
	}

	// third, find all testcases belong to the sections
	// and then update all testcases' section_id to 'Uncategorized' section
	if err := c.Db.Model(targetModel).Where("section_id in (?)", targetIDs).Update("section_id", uncategorizedSection.ID); err.Error != nil {
		log.Error("Section", "type", "database", "msg ", err.Error.Error(), "additional", "update")
		return RenderJSONWithStatus(w, Resp{Msg: err.Error.Error()}, http.StatusInternalServerError)
	}

	// fourth, delete all decendant sections and target sections
	if err := c.Db.Where("parents_id in (?)", targetIDs).Delete(models.Section{}); err.Error != nil {
		log.Error("Section", "type", "database", "msg ", err.Error.Error())
		return RenderJSONWithStatus(w, Resp{Msg: err.Error.Error()}, http.StatusInternalServerError)
	}

	if err := c.Db.Where("id = ?", sectionID).Delete(models.Section{}); err.Error != nil {
		log.Error("Section", "type", "database", "msg ", err.Error.Error())
		return RenderJSONWithStatus(w, Resp{Msg: err.Error.Error()}, http.StatusInternalServerError)
	}

	type ResponseData struct {
		Msg  string `json:"msg"`
		Data string `json:"data"`
	}

	var sections []models.Section
	if err := c.Db.Where("project_id = ?", prj.ID).Find(&sections); err.Error != nil {
		log.Error("TestCase", "type", "database", "msg ", err.Error)
		return RenderJSONWithStatus(w, ResponseData{Msg: "Error while delete linking"},
			http.StatusInternalServerError)
	}

	treeData := getJSTreeNodeData(sections)

	// this should return jstree data and make client to refresh jstree
	return RenderJSONWithStatus(w, ResponseData{Msg: "OK", Data: *treeData}, http.StatusOK)
}

// SectionAdd is not implemented.
func SectionAdd(c *interfacer.AppContext, w http.ResponseWriter, r *http.Request) error {
	// get "section_id" from post form
	if err := r.ParseForm(); err != nil {
		log.Error("Section", "type", "app", "msg ", err.Error())
		return RenderJSONWithStatus(w, Resp{Msg: "Parse form is not valid"}, http.StatusBadRequest)
	}

	// get parent node
	parentIDStr := r.FormValue("parent_id")
	isRootNode := false
	var parentID int
	if parentIDStr == "#" {
		isRootNode = true
	} else {
		var err error
		parentID, err = strconv.Atoi(parentIDStr)
		if err != nil {
			log.Error("Section", "type", "app", "msg ", err.Error())
			return RenderJSONWithStatus(w, Resp{Msg: err.Error()}, http.StatusBadRequest)
		}
	}

	title := r.FormValue("title")
	isTestCase := true
	if r.FormValue("is_test_case") == "false" {
		isTestCase = false
	}

	var prj models.Project
	if err := c.Db.Where("name = ?", c.ProjectName).First(&prj); err.Error != nil {
		log.Error("Section", "type", "app", "msg ", err.Error.Error())
		return RenderJSONWithStatus(w, Resp{Msg: err.Error.Error()}, http.StatusInternalServerError)
	}

	section := models.Section{ProjectID: prj.ID, Title: title, Status: 1,
		SpecialNode: false, ForTestCase: isTestCase, RootNode: isRootNode,
		ParentsID: parentID,
	}

	if errorMap, err := section.Validate(); err != nil {
		return logAndRenderJSONWithStatus(w, Resp{Msg: err.Error()}, http.StatusInternalServerError,
			logTypeErr, "Section", "app", err.Error())

	} else {
		if len(errorMap) > 0 {
			return logAndRenderJSONWithStatus(w, Resp{Msg: err.Error()}, http.StatusInternalServerError,
				logTypeErr, "Section", "app", "Validation failed")
		}
	}

	c.Db.NewRecord(section)
	if err := c.Db.Create(&section).Error; err != nil {
		return logAndRenderJSONWithStatus(w, Resp{Msg: err.Error()}, http.StatusInternalServerError,
			logTypeErr, "Section", "database", err.Error())
	}

	return RenderJSONWithStatus(w, Resp{Msg: "ok"}, http.StatusOK)
}
