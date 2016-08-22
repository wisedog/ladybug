package buildtools

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/wisedog/ladybug/models"
	log "gopkg.in/inconshreveable/log15.v2"
)

// Travis for interfaces
type Travis struct {
}

// belows are data structure for fetching JSON from Travis CI
type travisBuildInfo struct {
	State      string `json:"state"`
	ID         int    `json:"id"`
	Number     string `json:"number"`
	FinishedAt string `json:"finished_at"`
}

type travisBuildBunch struct {
	Builds []travisBuildInfo `json:"builds"`
}

type travisRepoInfo struct {
	ID               int    `json:"id"`
	Slug             string `json:"slug"`
	Active           bool   `json:"active"`
	Description      string `json:"description"`
	LastBuildNumber  string `json:"last_build_number"`
	LastBuildState   string `json:"last_build_state"`
	LastBuildStarted string `json:"last_build_started_at"`
	LastBuildFinish  string `json:"last_build_finished_at"`
}

type travisRepo struct {
	Repo travisRepoInfo `json:"repo"`
}

// AddTravisBuilds is main routine of adding travis build items.
func (t Travis) AddTravisBuilds(url string, projectID int, db *gorm.DB) error {

	repo, u, err := t.ConnectionTest(url)
	if err != nil {
		return errors.New("Fail to get Travis repo information")
	}

	i, err := strconv.Atoi(repo.Repo.LastBuildNumber)
	if err != nil {
		return errors.New("Fail to convert string to integer")
	}

	// get only lastest 3 builds
	var buildsNum []int
	for k := i; k > i-3; k-- {
		if k < 0 {
			break
		}
		buildsNum = append(buildsNum, k)
	}

	log.Debug("Travis", "buildsNum:", buildsNum)

	status := 0
	if repo.Repo.LastBuildState == "passed" {
		status = BuildSuccess
	} else {
		status = BuildFail
	}

	// save Build project
	job := models.Build{
		Name:         repo.Repo.Slug,
		Description:  repo.Repo.Description,
		ProjectID:    projectID,
		BuildURL:     url,
		Status:       status,
		ToolName:     "Travis",
		BuildItemNum: len(buildsNum),
	}

	r := db.Save(&job)
	if r.Error != nil {
		return r.Error
	}

	for _, build := range buildsNum {
		reqURL := u + "/builds?number=" + strconv.Itoa(build)
		body, err := t.getTravisRepoInfo(reqURL)
		if err != nil {
			return errors.New("Fail to get builds")
		}

		b := travisBuildBunch{}
		if err := json.Unmarshal(body, &b); err != nil {
			return errors.New("Json Unmarshalling is failed")
		}

		var rv int
		if b.Builds[0].State == "passed" {
			rv = BuildSuccess
		} else {
			rv = BuildFail
		}

		buildid := strconv.Itoa(b.Builds[0].ID)
		displayName := "#" + b.Builds[0].Number
		elem := models.BuildItem{
			BuildProjectID:  job.ID,
			IDByTool:        buildid,
			IDByToolInt:     b.Builds[0].ID,
			DisplayName:     displayName,
			FullDisplayName: job.Name + " " + displayName,
			ItemURL:         url + "/builds/" + buildid,
			ArtifactsURL:    "",
			ArtifactsName:   "",
			Result:          b.Builds[0].State,
			Toolname:        "Travis",
			TimeStamp:       0,
			Status:          rv,
		}

		r = db.Save(&elem)
		if r.Error != nil {
			return errors.New("Fail to save Jenkins job information")
		}
	}

	return nil
}

// ConnectionTest verifies the given repo's URL
func (t Travis) ConnectionTest(url string) (*travisRepo, string, error) {
	u := t.getAPIURL(url)

	res := travisRepo{}
	body, err := t.getTravisRepoInfo(u)
	if err != nil {
		return &res, u, errors.New("Fail to get Travis repo information")
	}

	if err := json.Unmarshal(body, &res); err != nil {
		return &res, u, errors.New("Json Unmarshalling is failed")
	}

	return &res, u, nil
}

// getAPIURL manipulates from repo URL to API-url
// For example, the given argument is "https://travis-ci.org/wisedog/ladybug"
// this function change it to "https://api-travis-ci.org/repos/wisedog/ladybug"
// TODO now Ladybug only support open-source Travis CI, should support commercial plan
func (t Travis) getAPIURL(url string) string {

	var k string
	if strings.HasPrefix(url, "https://travis-ci.org/") {
		k = "https://api.travis-ci.org/repos/" + url[len("https://travis-ci.org/"):]
	}
	return k
}

// getTravisRepoInfo gets travis repo's information
func (t Travis) getTravisRepoInfo(url string) ([]byte, error) {

	// Travis should be set below:
	// user-agent: "something"
	// Accept: application/vnd.travis-ci.2+ json
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Error("Travis", "Request error", err.Error())
		return nil, err
	}

	// Travis only accepts http request which has the headers below
	req.Header.Set("User-Agent", "Ladybug")
	req.Header.Set("Accept", "application/vnd.travis-ci.2+ json")

	resp, err := client.Do(req)
	if err != nil {
		log.Error("Travis", "Request error", err.Error())
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("Travis", "An error while readall", err)
		return nil, err
	}

	return body, nil
}
