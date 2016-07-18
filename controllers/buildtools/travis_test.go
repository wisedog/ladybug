package buildtools

import (
  "fmt"
  "testing"
  "os"

  "github.com/jinzhu/gorm"
  "github.com/wisedog/ladybug/interfacer"
  "github.com/wisedog/ladybug/database"
  "github.com/wisedog/ladybug/models"

  )

var Database *gorm.DB

func setup(){
  fmt.Println("Setup Testing for package buildtools...")
  cf := interfacer.LoadConfig()

  var err error
  Database, err = database.InitDB(cf) 
  if err != nil{
    fmt.Println("Database initialization is failed.")
  }
  //defer Database.Close()
}

func tearDown(){
  fmt.Println("Tear Down Testing for package buildtools...")
}

func TestMain(m *testing.M) { 
  fmt.Println("TestMain for package buildtools...")
  setup()
  retCode := m.Run()
  tearDown()
  os.Exit(retCode)
}


func TestGetApiUrlTravis(t *testing.T) {
  cases := []struct{
    in, expected string
  }{
    {"", ""},
    {"https://travis-ci.org/wisedog/ladybug", "https://api.travis-ci.org/repos/wisedog/ladybug"},
  }

  var travis Travis

  for _, c := range cases {
    got := travis.getApiUrl(c.in)
    if got != c.expected {
      t.Fail()
      t.Logf("getApiUrl(%q) == %q, expected %q", c.in, got, c.expected)
    }
  }
}

func TestTravisConnectionTest(t *testing.T){
  var travis Travis

  _, _ , err:= travis.ConnectionTest("https://travis-ci.org/wisedog/ladybug")
  if err != nil{
    t.Fail()
    t.Log("ConnectionTest Failed with URL : ", "https://travis-ci.org/wisedog/ladybug", " Msg : ", err.Error())
  }
}


func TestConnectionTest_Failed(t *testing.T){
  var travis Travis

  _, _, err := travis.ConnectionTest("")
  if err == nil{
    t.Fail()
    t.Log("ConnectionTest Failed with URL : \" \" should be failed")
  }
}

func TestGetTravisRepoInfo(t *testing.T){
  var travis Travis

  rv, err := travis.getTravisRepoInfo("https://travis-ci.org/wisedog/ladybug")

  if err != nil{
    t.Fail()
    t.Log("getTravisRepoInfo Failed with URL https://api.travis-ci.org/repos/wisedog/ladybug", "msg", err.Error())    
  }

  if len(rv) == 0{
    t.Fail()
    t.Log("getTravisRepoInfo received 0 byte")
  }
}

func TestAddTravisBuilds(t *testing.T){

  var travis Travis
  // project id 1 depends on createDummy() script. 
  if err := travis.AddTravisBuilds("https://travis-ci.org/wisedog/ladybug", 1, Database ); err != nil{
    t.Fail()
    t.Log("AddTravisBuilds returns error : ", err.Error())
  }
  var items []models.BuildItem
  if err := Database.Find(&items); err.Error != nil{
    t.Error("Nothing inserted", err.Error)
  }
}