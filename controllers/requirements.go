package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/wisedog/ladybug/errors"
	"github.com/wisedog/ladybug/interfacer"
	"github.com/wisedog/ladybug/models"

	log "gopkg.in/inconshreveable/log15.v2"
)

// RequirementIndex renders requirement index page.
func RequirementIndex(c *interfacer.AppContext, w http.ResponseWriter, r *http.Request) error {
	var prj models.Project

	if err := c.Db.Where("name = ?", c.ProjectName).First(&prj); err.Error != nil {
		log.Error("Requirement", "type", "database", "msg ", err.Error)
		return errors.HttpError{Status: http.StatusInternalServerError,
			Desc: "Project is not found in TestDesign.Index"}
	}

	var sections []models.Section
	c.Db.Where("project_id = ?", prj.ID).Where("for_test_case = ?", false).Find(&sections)

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

	items := map[string]interface{}{
		"TreeData":   treedata,
		"Active_idx": 6,
	}

	return Render2(c, w, items, "views/base.tmpl", "views/requirements/req.tmpl")
}

// RequirementList returns a list of specifications which belongs to requested section
func RequirementList(c *interfacer.AppContext, w http.ResponseWriter, r *http.Request) error {
	var reqs []models.Requirement

	vars := mux.Vars(r)
	sectionID := vars["id"]

	var prj models.Project
	if err := c.Db.Where("name = ?", c.ProjectName).First(&prj); err.Error != nil {
		log.Error("Requirment", "type", "database", "msg ", err.Error)
		return errors.HttpError{Status: http.StatusInternalServerError, Desc: "Not found project"}
	}

	if err := c.Db.Where("section_id = ?", sectionID).Find(&reqs); err.Error != nil {
		log.Error("Requirement", "type", "database", "msg ", err.Error)
		return errors.HttpError{Status: http.StatusInternalServerError, Desc: "Not found section"}
	}

	return RenderJSONWithStatus(w, reqs, http.StatusOK)
}

// ViewRequirement renders a page of information of the requirements.
// This page should be a detail of requested requirement and related testcases
func ViewRequirement(c *interfacer.AppContext, w http.ResponseWriter, r *http.Request) error {
	var req models.Requirement

	vars := mux.Vars(r)
	id := vars["id"]

	var prj models.Project
	if err := c.Db.Where("name = ?", c.ProjectName).First(&prj); err.Error != nil {
		log.Error("Requirement", "type", "database", "msg ", err.Error)
		return errors.HttpError{Status: http.StatusBadRequest, Desc: "Not found project"}
	}

	if err := c.Db.Where("project_id = ?", prj.ID).Where("id = ?", id).Preload("ReqType").First(&req); err.Error != nil {
		log.Error("Requirement", "type", "database", "msg ", err.Error)
		return errors.HttpError{Status: http.StatusBadRequest, Desc: "Not found requirements"}
	}

	req.StatusStr = getReqStatusI18n(req.Status)

	// Find related Test Cases with this requirement
	var testcases []models.TestCase
	if err := c.Db.Model(&req).Association("RelatedTestCases").Find(&testcases); err.Error != nil {
		log.Error("Requirement", "type", "database", "msg ", err.Error)
		return errors.HttpError{Status: http.StatusBadRequest, Desc: "Not found related testcases"}
	}

	for i := 0; i < len(testcases); i++ {
		testcases[i].PriorityStr = getPriorityI18n(testcases[i].Priority)
	}

	items := map[string]interface{}{
		"Requirement":      req,
		"RelatedTestCases": testcases,
		"Active_idx":       6,
	}

	return Render2(c, w, items, "views/base.tmpl", "views/requirements/reqview.tmpl")
}

// AddRequirement renders just a page that a user can add a requirement
func AddRequirement(c *interfacer.AppContext, w http.ResponseWriter, r *http.Request) error {
	return addEditRequirement(c, w, r, false)
}

// EditRequirement just renders a page that user can edit the requirement.
func EditRequirement(c *interfacer.AppContext, w http.ResponseWriter, r *http.Request) error {
	return addEditRequirement(c, w, r, true)
}

// addEditRequirement renders a page that user can edit or add the requirement
// Adding and Editing are somewhat similar, this function can handle both
func addEditRequirement(c *interfacer.AppContext, w http.ResponseWriter, r *http.Request, isEdit bool) error {
	var errorMap map[string]string
	var req models.Requirement

	session, e := c.Store.Get(r, "ladybug")

	if e != nil {
		log.Info("error ", "msg", e.Error())
	}

	// Check if there are invalid form values from SAVE/UPDATE
	if fm := session.Flashes(BuildFlashKey); fm != nil {
		b, ok := fm[0].(*models.Requirement)
		if ok {
			req = *b
		} else {
			log.Debug("Build", "msg", "flash type assertion failed")
		}

		delete(session.Values, BuildFlashKey)
		errorMap = getErrorMap(session)
		session.Save(r, w)

	} else {
		var id int
		var sectionID int
		var err error

		// only addRequirement has get a parameter
		if isEdit {
			// find [add|save]/{id}
			vars := mux.Vars(r)
			idStr := vars["id"]
			id, err = strconv.Atoi(idStr)
			if err != nil {
				return LogAndHTTPError(http.StatusBadRequest, "Requirement", "http", "can not convert req id")
			}
		} else {
			strSectionID := r.URL.Query().Get("section_id")

			sectionID, err = strconv.Atoi(strSectionID)
			if err != nil {
				return LogAndHTTPError(http.StatusBadRequest, "Requirement", "http", "can not convert section id")
			}
		}

		if isEdit {
			if err := c.Db.Where("id = ?", id).First(&req); err.Error != nil {
				var msg string
				var status int
				if err.RecordNotFound() {
					msg = "Bad request. Record not found"
					status = http.StatusBadRequest
				} else {
					msg = "An Error while SELECT operation for Requirement.Edit"
					status = http.StatusInternalServerError
				}
				log.Error("req", "type", "database", "msg", err.Error)
				return LogAndHTTPError(status, "requirement", "database", msg)
			}
		} else {
			var prj models.Project
			if err := c.Db.Where("name = ?", c.ProjectName).First(&prj); err.Error != nil {
				log.Error("requirement", "type", "database", "msg", err.Error)
				return errors.HttpError{Status: http.StatusInternalServerError,
					Desc: "An Error while SELECT project operation for Requirement.Add"}
			}
			req.SectionID = sectionID
			req.ProjectID = prj.ID
		}

	}

	// get ReqTypes values like Use Case, Information, Feature .. to render
	var reqTypes []models.ReqType
	if err := c.Db.Find(&reqTypes); err.Error != nil {
		return LogAndHTTPError(http.StatusInternalServerError, "Requirement", "db", "empty reqtype")
	}

	items := map[string]interface{}{
		"Requirement": req,
		"ReqType":     reqTypes,
		"IsEdit":      isEdit,
		"Active_idx":  6,
		"ErrorMap":    errorMap,
	}

	return Render2(c, w, items, "views/base.tmpl", "views/requirements/reqaddedit.tmpl")

}

// saveUpdateRequirement handles POST add or edit request, return in JSON format
func saveUpdateRequirement(c *interfacer.AppContext, w http.ResponseWriter, r *http.Request, isUpdate bool) error {
	var req models.Requirement
	vars := mux.Vars(r)
	idStr := vars["id"]

	if err := r.ParseForm(); err != nil {
		log.Error("Requirement", "Error ", err)
		return errors.HttpError{Status: http.StatusInternalServerError, Desc: "ParseForm failed"}
	}

	if err := c.Decoder.Decode(&req, r.PostForm); err != nil {
		log.Warn("Requirement", "Error", err, "msg", "Decode failed but go ahead")
	}

	//redirectionTarget := fmt.Sprintf("/project/%s/req", c.ProjectName)

	// from UPDATE page
	if isUpdate {
		if errorMap, err := req.Validate(); err != nil {
			return LogAndHTTPError(http.StatusInternalServerError, "Requirement", "app", "Failed to validate")
		} else {
			if len(errorMap) > 0 {
				ValidationFailAndRedirect(c, w, r, errorMap, "/project/"+c.ProjectName+"/req/edit/"+idStr, req)
				return nil
			}
		}

		if err := c.Db.Save(&req); err.Error != nil {
			return LogAndHTTPError(http.StatusInternalServerError, "Requirement", "database", "update failed")
		}

		// TODO Get Note and add to history
		//note := r.FormValue("Note")

	} else {
		// from ADD page

		// section ID is from form value (string type)
		sectionIDStr := r.FormValue("SectionID")

		// convert into integer
		sectionID, err := strconv.Atoi(sectionIDStr)
		if err != nil {
			return LogAndHTTPError(http.StatusBadRequest, "Requirement", "http",
				"Failed to add the requirement(Coverting Error)")
		}

		req.SectionID = sectionID
		// get project ID
		var prj models.Project
		if err := c.Db.Where("name = ?", c.ProjectName).First(&prj); err.Error != nil {
			return LogAndHTTPError(http.StatusBadRequest, "Requirement", "database",
				"Failed to add the requirement")
		}

		req.ProjectID = prj.ID
		if errorMap, err := req.Validate(); len(errorMap) > 0 {
			if err != nil {
				return LogAndHTTPError(http.StatusInternalServerError, "Requirement", "validation",
					err.Error())
			}
			ValidationFailAndRedirect(c, w, r, errorMap, "/project/"+c.ProjectName+"/req/add?section_id="+sectionIDStr, req)
			return nil
		}

		// Create!
		c.Db.NewRecord(req)

		if err := c.Db.Create(&req); err.Error != nil {
			return LogAndHTTPError(http.StatusInternalServerError, "Requirement", "database",
				"Failed to add the requirement")
		}
	}

	http.Redirect(w, r, "/project/"+c.ProjectName+"/req", http.StatusFound)
	return nil
}

// UpdateRequirement is gateway of POST Update request from EDIT page
func UpdateRequirement(c *interfacer.AppContext, w http.ResponseWriter, r *http.Request) error {
	return saveUpdateRequirement(c, w, r, true)
}

// SaveRequirement is gateway of POST Save request from ADD page
func SaveRequirement(c *interfacer.AppContext, w http.ResponseWriter, r *http.Request) error {
	return saveUpdateRequirement(c, w, r, false)
}

// DeleteRequirement handles Delete request and return JSON
func DeleteRequirement(c *interfacer.AppContext, w http.ResponseWriter, r *http.Request) error {
	// TODO check auth, if not satisfied, return http.StatusForbidden

	if err := r.ParseForm(); err != nil {
		log.Error("Requirement", "Error ", err)
		return RenderJSONWithStatus(w, Resp{Msg: "Not found requested requirement"}, http.StatusBadRequest)
	}

	id := r.FormValue("id")
	// Delete the testcase  permanently for sequence
	if err := c.Db.Where("id = ?", id).Delete(&models.Requirement{}); err.Error != nil {
		var errType int
		var errMsg string

		if err.RecordNotFound() {
			errMsg = "Not found requirement to delete"
			errType = http.StatusBadRequest
		} else {
			errMsg = "An error while delete requirement"
			errType = http.StatusInternalServerError
		}
		log.Error("Requirement", "msg", errMsg, "id", id, "raw", err.Error)
		return RenderJSONWithStatus(w, Resp{Msg: "Fail to delete"}, errType)
	}

	// TODO unlink testcase relation, remove history related this requirement

	return RenderJSONWithStatus(w, Resp{Msg: "OK"}, http.StatusOK)
}

// UnlinkTestcaseRelation unlink a requirement and a related testcase
// Return in JSON
func UnlinkTestcaseRelation(c *interfacer.AppContext, w http.ResponseWriter, r *http.Request) error {
	// get requirement, testcase ID from post
	if err := r.ParseForm(); err != nil {
		log.Error("Requirement", "type", "app", "msg ", err.Error())
		return RenderJSONWithStatus(w, Resp{Msg: "Parse form is not valid"}, http.StatusBadRequest)
	}

	reqIDStr := r.FormValue("req_id")
	tcIDStr := r.FormValue("tc_id")

	reqID, err := strconv.Atoi(reqIDStr)
	if err != nil {
		log.Error("Requirement", "type", "app", "msg ", err.Error())
		return RenderJSONWithStatus(w, Resp{Msg: err.Error()}, http.StatusBadRequest)
	}

	tcID, err := strconv.Atoi(tcIDStr)
	if err != nil {
		log.Error("Requirement", "type", "app", "msg ", err.Error())
		return RenderJSONWithStatus(w, Resp{Msg: err.Error()}, http.StatusBadRequest)
	}

	var targetTestCase models.TestCase
	if err := c.Db.Where("id = ?", tcID).First(&targetTestCase); err.Error != nil {
		log.Error("Requirement", "type", "database", "msg ", err.Error.Error())
		return RenderJSONWithStatus(w,
			Resp{Msg: "Not found or something wrong while unlinking requirement and related testcase"},
			http.StatusInternalServerError)
	}

	var req models.Requirement
	req.ID = reqID

	// Find related Test Cases with this requirement
	//var reqs []models.Requirement
	if err := c.Db.Model(&req).Association("RelatedTestCases").Delete(targetTestCase); err.Error != nil {
		log.Error("Requirement", "type", "database", "msg ", err.Error.Error())
		return RenderJSONWithStatus(w, Resp{Msg: "Error while delete linking"}, http.StatusInternalServerError)
	}

	// TODO tc_req_relation
	// TODO add to history

	return RenderJSONWithStatus(w, Resp{Msg: "OK"}, http.StatusOK)
}

// LinkTestcaseRelation links a requirement and a related testcase
// Return in JSON
func LinkTestcaseRelation(c *interfacer.AppContext, w http.ResponseWriter, r *http.Request) error {
	// add tc_req_relation
	return nil
}
