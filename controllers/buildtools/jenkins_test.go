package buildtools

import(
	"testing"

  "github.com/wisedog/ladybug/models"
)

func TestJenkinsConnectionTest(t *testing.T) {
  var j Jenkins

  if _, err := j.ConnectionTest("https://builds.apache.org/view/All/job/Abdera2-trunk/"); err != nil{
    t.Error("ConnectionTest for Jenkins has failed", err.Error)
  }

}

func TestGetJenkinsJobInfo(t *testing.T) {
  var j Jenkins
  if _, err := j.getJenkinsJobInfo("https://builds.apache.org/view/All/job/Abdera2-trunk/api/json"); err != nil{
    t.Error("GetJenkinsJobInfo has failed.", err.Error())
  }
}



func TestAddJenkinsBuilds(t *testing.T){

  var jenkins Jenkins
  // project id 1 depends on createDummy() script. 
  if err := jenkins.AddJenkinsBuilds("https://builds.apache.org/view/All/job/Abdera2-trunk/", 1, Database ); err != nil{
    t.Error("AddTravisBuilds returns error : ", err.Error())
  }
  var items []models.BuildItem
  if err := Database.Find(&items); err.Error != nil{
    t.Error("Nothing inserted")
  }
  
}