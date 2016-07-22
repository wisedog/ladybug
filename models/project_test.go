package models

import (
	"testing"
)

// TestProjectValidationEmptyProjectName tests empty project name
func TestProjectValidationEmptyProjectName(t *testing.T) {

	// check empty project name
	prj := Project{Name: "", Prefix: "aaa"}
	errorMap := prj.Validate()

	if len(errorMap) == 0 {
		t.Errorf(`models.Project.Validate() error ({Name:"", Prefix:"aaa"})`)
	}

	if errorMap["Name"] != "Name is required." {
		t.Errorf(`({Name:"", Prefix:"aaa"}) errorMap['name'] = %q expected %q`, errorMap["Name"], "Name is required.")
	}
}

// TestProjectValidateNormalCondition tests normal condition
// No error should be
func TestProjectValidateNormalCondition(t *testing.T) {
	// normal condition
	prj := Project{Name: "Lübeck", Prefix: "HL"}
	errorMap := prj.Validate()

	if len(errorMap) > 0 {
		t.Errorf(`models.Project.Validate() error ({Name : "Lübeck", Prefix:"HL"})`)
	}

	if errorMap["Name"] != "" {
		t.Errorf(`({Name : "Lübeck", Prefix:"HL"}) errorMap['name'] = %q expected %q`, errorMap["Name"], "")
	}

}

// TestProjectValidateEmptyPrefix tests empty prefix
func TestProjectValidateEmptyPrefix(t *testing.T) {
	// check empty prefix
	prj := Project{Name: "Hamburg", Prefix: ""}
	errorMap := prj.Validate()

	if len(errorMap) == 0 {
		t.Errorf(`models.Project.Validate() error ({Name : "Hamburg", Prefix:""})`)
	}

	if errorMap["Prefix"] != "Prefix is required." {
		t.Errorf(`({Name : "Hamburg", Prefix:""}) errorMap["Prefix"] = %q expected %q`, errorMap["Prefix"], "Prefix is required.")
	}

}

// TestProjectValidatePrefix tests long prefix
func TestProjectValidatePrefix(t *testing.T) {
	// check max size of prefix
	// check empty prefix
	prj := Project{Name: "Kiel", Prefix: "kielkielkielkiel"}
	errorMap := prj.Validate()

	if len(errorMap) == 0 {
		t.Errorf(`models.Project.Validate() error ({Name : "Kiel", Prefix:"kielkielkielkiel"})`)
	}

	if errorMap["Prefix"] != " Size of Prefix exceeds 12 characters." {
		t.Errorf(`({Name : "Kiel", Prefix:"kielkielkielkiel"}) errorMap["Prefix"] = %q expected %q`, errorMap["Prefix"], " Size of Prefix exceeds 12 characters.")
	}

}
