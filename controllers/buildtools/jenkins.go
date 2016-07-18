package buildtools

import (
  "strconv"
  "strings"
  "errors"
  "time"
  
  "net/http"
  "io/ioutil"
  "encoding/json"
    
  "github.com/jinzhu/gorm"
	"github.com/wisedog/ladybug/models"
  log "gopkg.in/inconshreveable/log15.v2"
)

const (
	BuildFail = iota
	BuildSuccess
	
	)

type Jenkins struct {
    
}

func (j Jenkins) AddJenkinsBuilds(url string, projectID int, db *gorm.DB) error{
  if db == nil{
    return errors.New("Wrong database handler!")
  }
  
	if strings.HasSuffix(url, "/api/json") == false{
    	url = url + "/api/json"
  }
  
  body, err := j.getJenkinsJobInfo(url)
  if err != nil {
    return errors.New("Fail to get Jenkins job informations")
  }
  var dat map[string]interface{}
  if err := json.Unmarshal(body, &dat); err != nil {
    return errors.New("Fail to unmarshall response")
  }
  
  name := dat["name"].(string)
  nextBuildNum := int(dat["nextBuildNumber"].(float64))
  lastSuccessfulBuild := dat["lastSuccessfulBuild"].(map[string]interface{})
  lastSucessfulBuildNum := int(lastSuccessfulBuild["number"].(float64))
  
  
  // get status for building. it may be successful or failed
  status := 0
  if nextBuildNum -1 == lastSucessfulBuildNum {
  	status = BuildSuccess
  } else{
  	status = BuildFail
  }

  builds := dat["builds"].([]interface{})
  
  job := models.Build{
  	Name : name,
  	Description : dat["description"].(string),
  	Project_id : projectID,
  	BuildUrl : dat["url"].(string),
  	Status : status,
  	ToolName : "jenkins",
  	BuildItemNum : len(builds),
  }
  r := db.Save(&job)
  if r.Error != nil{
  	return r.Error
  }
  
  
  for idx, b := range builds {
  	if idx > 10 {
  		break
  	}
  	
  	k := b.(map[string]interface{})
  	
  	targetURL := k["url"].(string) + "/api/json"
  	info, err := j.getJenkinsJobInfo(targetURL)
  	if err != nil{
  		continue
  	}
  	
  	var data map[string]interface{}
    if err := json.Unmarshal(info, &data); err != nil {
    	return nil
    }
  	
  	timestamp := int64(data["timestamp"].(float64))
  	displayname := data["displayName"].(string)
  	idbytool := data["id"].(string)
  	result := data["result"].(string)
  	url := data["url"].(string)
  	
  	
  	artifacts := data["artifacts"].([]interface{})
  	
  	num := len(artifacts)
  	var artifactsname string
  	var artifactsurl	string
  	
  	if num > 1 {
  		// TODO It is not supported to link multiple artifacts now
  		artifactsurl = url
  		artifactsname = "Multiple"
  	}else if num == 1 {
  		a := artifacts[0].(map[string]interface{})
  		artifactsname = a["fileName"].(string)
  		artifactsurl = url + "artifact/" + a["relativePath"].(string)
  	}else {
  		artifactsurl = ""
  		artifactsname = ""
  	}
  	
  	// timestamp of jenkins build item represents in millisecond.
  	// so divide by 1000 (1 second = 1000 milliseconds)
  	buildat := time.Unix(int64(timestamp/1000),0)
  	
  	var rv int 
  	if result == "SUCCESS"{
  		rv = BuildSuccess
  	}else {
  		rv = BuildFail
  	}
  	
  	
  	elem := models.BuildItem{
  		BuildProjectID : job.ID,
  		IdByTool : idbytool,
  		DisplayName : displayname,
  		FullDisplayName : name + " " + displayname,
  		ItemUrl : url,
  		ArtifactsUrl : artifactsurl,
  		ArtifactsName : artifactsname,
  		Result : result,
  		Toolname : "jenkins",
  		TimeStamp : timestamp,
  		BuildAt : buildat,
  		Status : rv,
  	}
  	
  	// save to BuildItem
  	r = db.Save(&elem)
  	if r.Error != nil{
  	  return errors.New("Fail to save Jenkins job information")
  	}
  }
  return nil
}

func (j Jenkins) ConnectionTest(url string) (string, error) {
  if strings.HasSuffix(url, "/api/json") == false{
  	url = url + "/api/json"
  }
  body, err := j.getJenkinsJobInfo(url)
	if err != nil {
	  return "", errors.New("Fail to get Jenkins job information")
	}
	
	var dat map[string]interface{}
	var msg string
	
  if err := json.Unmarshal(body, &dat); err != nil {
    return "", errors.New("Json Unmarshalling is failed")
  }
  
  msg += "Job name : "+ dat["name"].(string) + "\n"
  msg += "URL : " + dat["url"].(string) + "\n"
  
  build := dat["lastBuild"].(map[string]interface{})
  
  // FIXME runtime error build["number"].(int).... why? 3 is float?
  k := strconv.Itoa(int(build["number"].(float64)))
  msg += "LastBuild Number : " + k + "\n"
  
  build = dat["lastSuccessfulBuild"].(map[string]interface{})
  k = strconv.Itoa(int(build["number"].(float64)))
  msg += "Last Successful Build : " + k + "\n"
	
	return msg, nil
}

// getJenkinsJobInfo verifies given url.
func (j Jenkins) getJenkinsJobInfo(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
    log.Error("Jenkins", "msg", "An error while GET", "raw", err)
		return nil, err
	}
	defer resp.Body.Close()
	
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil{
    log.Error("Jenkins", "msg", "An error while readall", "raw", err)
		return nil, err
	}
	return body, nil
}