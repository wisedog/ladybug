package buildtools

import(
	"testing"


  "github.com/wisedog/ladybug/database"
  "github.com/wisedog/ladybug/models"
)

func TestJenkinsConnectionTest(t *testing.T) {
  var j Jenkins

  if _, err := j.ConnectionTest("https://builds.apache.org/job/beam_PreCommit"); err != nil{
    t.Error("ConnectionTest for Jenkins has failed", err.Error)
  }

}

func TestGetJenkinsJobInfo(t *testing.T) {
  var j Jenkins
  if _, err := j.getJenkinsJobInfo("https://builds.apache.org/job/beam_PreCommit/api/json"); err != nil{
    t.Error("GetJenkinsJobInfo has failed.", err.Error())
  }
}



func TestAddJenkinsBuilds(t *testing.T){
  if err := database.InitDB(); err != nil{
    t.Log("Database initialization is failed.")
    return
  }
  //AddTravisBuilds(url string, projectID int, db *gorm.DB) error{

  var jenkins Jenkins
  db := &database.Database
  // project id 1 depends on createDummy() script. 
  if err := jenkins.AddJenkinsBuilds("https://builds.apache.org/job/beam_PreCommit", 1, db ); err != nil{
    t.Error("AddTravisBuilds returns error : ", err.Error())
  }
  var items []models.BuildItem
  if err := db.Find(&items); err.Error != nil{
    t.Error("Nothing inserted")
  }
  
}