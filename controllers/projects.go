package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/wisedog/ladybug/errors"
	"github.com/wisedog/ladybug/interfacer"
	"github.com/wisedog/ladybug/models"
	log "gopkg.in/inconshreveable/log15.v2"
)

const (
	ProjectFlashKey = "LADYBUG_PROJECT"
)

// ProjectCreate renders a page to create project
func ProjectCreate(c *interfacer.AppContext, w http.ResponseWriter, r *http.Request) error {
	var project models.Project
	var errorMap map[string]string

	session, e := c.Store.Get(r, "ladybug")

	if e != nil {
		log.Info("error ", "msg", e.Error())
	}

	// Check if there are invalid form values from SAVE/UPDATE
	if fm := session.Flashes(ProjectFlashKey); fm != nil {
		p, ok := fm[0].(*models.Project)
		if ok {
			project = *p
		} else {
			log.Debug("Build", "msg", "flash type assertion failed")
		}

		delete(session.Values, ProjectFlashKey)
		errorMap = getErrorMap(session)
		session.Save(r, w)
	} else {
		log.Debug("Project", "msg", "no flashes")
	}

	items := map[string]interface{}{
		"Project":    project,
		"ErrorMap":   errorMap,
		"Active_idx": 0,
	}

	return Render2(c, w, items, "views/base.tmpl", "views/projects/create_project.tmpl")
}

// ProjectSave validates user input and save Project in database
func ProjectSave(c *interfacer.AppContext, w http.ResponseWriter, r *http.Request) error {

	if err := r.ParseForm(); err != nil {
		log.Error("Projects", "type", "http", "msg ", err)
	}

	userListStr := r.FormValue("UserList")

	var project models.Project
	if err := c.Decoder.Decode(&project, r.PostForm); err != nil {
		log.Warn("Projects", "Error", err, "msg", "Decode failed but go ahead")
	}

	/*if c != nil{
	  http.Redirect(w, r, "/hello", http.StatusFound)
	  return nil
	}*/

	errorMap := project.Validate()
	var tmpProject models.Project

	if project.Name != "" {
		// check name duplication. models.Validate() can not detect duplication.
		if err := c.Db.Where("name = ?", project.Name).First(&tmpProject); err.Error == nil {
			if tmpProject.ID != 0 {
				errorMap["Name"] += " Duplicated Project Name"
			}
		}
	}

	tmpProject.ID = 0
	if project.Prefix != "" {
		// check prefix duplication
		if err := c.Db.Where("prefix = ?", project.Prefix).First(&tmpProject); err.Error == nil {
			if tmpProject.ID != 0 {
				errorMap["Prefix"] += " Duplicated Project Prefix"
			}
		}
	}

	if len(errorMap) > 0 {

		session, e := c.Store.Get(r, "ladybug")
		if e != nil {
			log.Warn("error ", "msg", e.Error)
		}

		log.Info("Projects", "asdf", errorMap)

		session.AddFlash(project, ProjectFlashKey)
		session.AddFlash(errorMap, ErrorMsg)

		session.Save(r, w)
		http.Redirect(w, r, "/project/create", http.StatusFound)
		return nil
	}

	project.Status = models.TcStatusActivate

	// Split on comma
	ids := strings.Split(userListStr, ",")

	var userList []models.User
	if len(ids) > 0 {
		if err := c.Db.Where("id in (?)", ids).Find(&userList).Error; err != nil {
			log.Error("Projects", "msg", err)
			return errors.HttpError{Status: http.StatusInternalServerError, Desc: "multiple users selection is failed"}
		}

		project.Users = userList
	}

	if rv := c.Db.NewRecord(project); rv == false {
		log.Error("Projects", "type", "database", "msg", "duplicated primary key")
		return errors.HttpError{Status: http.StatusInternalServerError, Desc: "duplicated primary key"}
	}

	if err := c.Db.Create(&project).Error; err != nil {
		log.Error("Projects", "Type", "database", "Error", err)
		return errors.HttpError{Status: http.StatusInternalServerError, Desc: "could not create a build model"}
	}

	http.Redirect(w, r, "/project/"+project.Name, http.StatusFound)
	return nil
}

// Dashboard renders a dashboard page
func Dashboard(c *interfacer.AppContext, w http.ResponseWriter, r *http.Request) error {

	var prj models.Project
	err := c.Db.Where("name = ?", c.ProjectName).First(&prj)

	if err.Error != nil {
		log.Error("Projects", "msg", err.Error)
		return errors.HttpError{Status: http.StatusInternalServerError, Desc: "Not found project"}
	}

	// counting testcases
	var testcases []models.TestCase
	c.Db.Where("project_id = ?", prj.ID).Order("created_at desc").Find(&testcases)
	tcCount := len(testcases)

	// generate graph data only before 4 weeks
	now := time.Now()
	n1 := now.AddDate(0, 0, -7)
	n2 := now.AddDate(0, 0, -14)
	n3 := now.AddDate(0, 0, -21)
	n4 := now.AddDate(0, 0, -28)
	var dataset [4]int
	for _, v := range testcases {
		if v.CreatedAt.After(n1) {
			dataset[3]++
		} else if v.CreatedAt.After(n2) {
			dataset[3]++
			dataset[2]++
		} else if v.CreatedAt.After(n3) {
			dataset[3]++
			dataset[2]++
			dataset[1]++
		} else if v.CreatedAt.After(n4) {
			dataset[3]++
			dataset[2]++
			dataset[1]++
			dataset[0]++
		}
	}

	// make string be parsed by GraphJS
	dataSetString := []string{}
	for i := range dataset {
		num := dataset[i]
		txt := strconv.Itoa(num)
		dataSetString = append(dataSetString, txt)
	}
	dataArray := strings.Join(dataSetString, ",")

	// deal with test executions
	var execs []models.Execution
	c.Db.Where("project_id = ?", prj.ID).Preload("Executor").Preload("Plan").Find(&execs)

	execCount := 0
	taskCount := 0
	for idx := 0; idx < len(execs); idx++ {
		element := execs[idx]
		if element.Status == models.ExecStatusReady {
			taskCount++
		} else if element.Status == models.ExecStatusInProgress {
			prg := float32(element.PassCaseNum+element.FailCaseNum) / float32(element.TotalCaseNum) * 100
			execs[idx].Progress = int(prg)
			taskCount++
		}
	}

	buildCount := 0
	c.Db.Model(models.Build{}).Where("project_id=?", prj.ID).Count(&buildCount)

	testPlanCount := 0
	c.Db.Model(models.TestPlan{}).Where("project_id=?", prj.ID).Count(&testPlanCount)

	var milestone models.Milestone
	c.Db.Where("project_id=?", prj.ID).Where("due_date >= ?", time.Now()).Order("due_date").First(&milestone)

	var daysLeft time.Duration
	if milestone.Name != "" {
		daysLeft = milestone.DueDate.Sub(time.Now())
	}

	// TODO show Requirements graph. Now shows dummy
	reqDataSet := "0,1,4,7"

	// to show Requirements coverage.
	// Showing coverage graph is difficult. Because relationship between test cases and requirements
	// changes every time. So we need to consider more how to get it.
	currentCoverage, _ := getCurrentReqTestCaseCoverage(c)

	periodCoverage, _ := getPeriodReqTestCaseCoverage(c, prj.ID, currentCoverage)
	fmt.Println("AA", periodCoverage)

	items := map[string]interface{}{
		"Project":         prj,
		"ExecCount":       execCount,
		"BuildCount":      buildCount,
		"TaskCount":       taskCount,
		"Tasks":           execs,
		"TestPlanCount":   testPlanCount,
		"TestCaseNum":     tcCount,
		"Milestone":       milestone,
		"DaysLeft":        int(daysLeft.Hours() / 24),
		"TestCaseDataSet": dataArray,
		"ReqDataSet":      reqDataSet,
		"Coverage":        currentCoverage,
		"CoverageDataSet": periodCoverage,
		"Active_idx":      1,
	}

	return Render2(c, w, items, "views/base.tmpl", "views/projects/dashboard.tmpl")
}

// GetProjectList returns project list in JSON format
func GetProjectList(c *interfacer.AppContext, w http.ResponseWriter, r *http.Request) error {
	var prjs []models.Project
	limitStr := r.URL.Query().Get("limit")
	limit, e := strconv.Atoi(limitStr)
	if e != nil {
		resp := map[string]interface{}{
			"ErrorMsg": "Invalid limit value",
		}
		return RenderJSONWithStatus(w, resp, http.StatusInternalServerError)
	}
	if err := c.Db.Limit(limit).Find(&prjs).Error; err != nil {
		resp := map[string]interface{}{
			"ErrorMsg": "Database Error",
		}
		return RenderJSONWithStatus(w, resp, http.StatusInternalServerError)
	}

	//TODO renderJson to cover "Show All Projects", Clicking all project in the right upper menu
	return RenderJSON(w, prjs)
}

//getPeriodReqTestCaseCoverage calculates coverage(%) of requirement by testcases in particular period
func getPeriodReqTestCaseCoverage(c *interfacer.AppContext, projectID, currentCov int) (string, error) {
	// req-testcase coverage is calculated by dealing with TcReqRelationHistory, Requirement table
	// This is pretty complex, if you have good idea to solve it : VERY WELCOME : )

	// first get all requirements in this project
	var reqs []models.Requirement
	if err := c.Db.Where("project_id = ?", projectID).Order("created_at").Find(&reqs); err.Error != nil {
		return "", err.Error
	}
	// second get all TcReqRelationHistory in this project
	var relations []models.TcReqRelationHistory
	if err := c.Db.Where("project_id = ?", projectID).Order("created_at").Find(&relations); err.Error != nil {
		return "", err.Error
	}

	// iteration and count pre-4weeks requirement count and relationship history
	// 1. get target requirements by period
	// 2. put their IDs to map (int, int)
	// 3. get target relationship by period
	// 4. find map item with requirement ID and increase value

	now := time.Now()
	n1 := now.AddDate(0, 0, -7)
	n2 := now.AddDate(0, 0, -14)
	n3 := now.AddDate(0, 0, -21)

	timeArray := [3]time.Time{n3, n2, n1}

	periodMap := make(map[int]int)

	covArray := []int{}

	// Question : querying is faster? or fetch all and three-dimension looping is faster?
	for _, t := range timeArray {
		for _, req := range reqs {
			if req.CreatedAt.Before(t) {
				periodMap[req.ID] = 0
			} else {
				continue
			}

			for _, rel := range relations {
				if rel.CreatedAt.Before(t) && req.ID == rel.RequirementID {
					count := periodMap[req.ID]
					if rel.Kind == models.TcReqRelationHistoryLink {
						periodMap[req.ID] = count + 1
					} else if rel.Kind == models.TcReqRelationHistoryUnlink {
						periodMap[req.ID] = count - 1
					}
				}
			}
		}
		var totalReq int
		var relatedReq int
		// dump map
		for _, v := range periodMap {
			totalReq++
			if v > 0 {
				relatedReq++
			}
		}

		fmt.Println("Map", periodMap)
		if totalReq == 0 {
			covArray = append(covArray, 0)
		} else {
			covArray = append(covArray, int(relatedReq*100/totalReq))
		}

	}

	covArray = append(covArray, currentCov)

	// make string be parsed by GraphJS
	dataSetString := []string{}
	for i := range covArray {
		num := covArray[i]
		txt := strconv.Itoa(num)
		dataSetString = append(dataSetString, txt)
	}
	dataArray := strings.Join(dataSetString, ",")
	return dataArray, nil
}

// getCurrentReqTestCaseCoverage calculates coverage(%) of requirement by testcases
func getCurrentReqTestCaseCoverage(c *interfacer.AppContext) (int, error) {
	type testcasesReqs struct {
		RequirementID int
		TestCAseID    int
	}

	var p []testcasesReqs

	// First get req-case join table. using half-raw query to fwtch
	if err := c.Db.Table("testcases_reqs").Select("requirement_id, test_case_id").
		Scan(&p); err.Error != nil {
		fmt.Println("Failed to select operation in Hello")
	}

	// Second, get requirements table
	var reqs []models.Requirement
	if err := c.Db.Find(&reqs); err.Error != nil {
		log.Error("Dashboard", "type", "database", "msg", "not found requirement")
		return 0, err.Error
	}

	// match, just counting unique requirement IDs
	matched := make(map[int]int)
	for _, req := range reqs {
		for _, record := range p {
			if req.ID == record.RequirementID {
				matched[req.ID] = 1
			}
		}
	}

	percent := int(len(matched) * 100 / len(reqs))

	return percent, nil
}
