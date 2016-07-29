package buildtools

import (
	"fmt"
	"os"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/wisedog/ladybug/database"
	"github.com/wisedog/ladybug/interfacer"
	"github.com/wisedog/ladybug/models"
)

// get Database instance
var Database *gorm.DB

func setup() {
	fmt.Println("Setup Testing for package buildtools...")
	cf := interfacer.LoadConfig()

	var err error
	Database, err = database.InitDB(cf)
	if err != nil {
		fmt.Println("Database initialization is failed.")
	}
	//defer Database.Close()
}

func tearDown() {
	fmt.Println("Tear Down Testing for package buildtools...")
}

func TestMain(m *testing.M) {
	fmt.Println("TestMain for package buildtools...")
	setup()
	retCode := m.Run()
	tearDown()
	os.Exit(retCode)
}

// TestGetAPIURLTravis tests getAPIURL() function.
func TestGetAPIURLTravis(t *testing.T) {
	cases := []struct {
		in, expected string
	}{
		{"", ""},
		{"https://travis-ci.org/wisedog/ladybug", "https://api.travis-ci.org/repos/wisedog/ladybug"},
	}

	var travis Travis

	for _, c := range cases {
		got := travis.getAPIURL(c.in)
		if got != c.expected {
			t.Fail()
			t.Logf("getAPIURL(%q) == %q, expected %q", c.in, got, c.expected)
		}
	}
}

// TestTravisConnectionTest tests TravisConnectionTest() function
func TestTravisConnectionTest(t *testing.T) {
	var travis Travis

	_, _, err := travis.ConnectionTest("https://travis-ci.org/wisedog/ladybug")
	if err != nil {
		t.Fail()
		t.Log("ConnectionTest Failed with URL : ", "https://travis-ci.org/wisedog/ladybug", " Msg : ", err.Error())
	}
}

// TestConnectionTestFailed tests ConnectionTest() on failure condition
func TestConnectionTestFailed(t *testing.T) {
	var travis Travis

	_, _, err := travis.ConnectionTest("")
	if err == nil {
		t.Fail()
		t.Log("ConnectionTest Failed with URL : \" \" should be failed")
	}
}

// TestGetTravisRepoInfo tests getTravisRepoInfo() function
func TestGetTravisRepoInfo(t *testing.T) {
	var travis Travis

	rv, err := travis.getTravisRepoInfo("https://travis-ci.org/wisedog/ladybug")

	if err != nil {
		t.Fail()
		t.Log("getTravisRepoInfo Failed with URL https://api.travis-ci.org/repos/wisedog/ladybug", "msg", err.Error())
	}

	if len(rv) == 0 {
		t.Fail()
		t.Log("getTravisRepoInfo received 0 byte")
	}
}

// TestAddTravisBuilds tests AddTravisBuilds() function
func TestAddTravisBuilds(t *testing.T) {
	var travis Travis
	// project id 1 depends on createDummy() script.
	if err := travis.AddTravisBuilds("https://travis-ci.org/wisedog/ladybug", 1, Database); err != nil {
		t.Fail()
		t.Log("AddTravisBuilds returns error : ", err.Error())
	}
	var items []models.BuildItem
	if err := Database.Find(&items); err.Error != nil {
		t.Error("Nothing inserted", err.Error)
	}
}
