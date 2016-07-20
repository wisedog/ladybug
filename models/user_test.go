package models

import (
	"testing"
)

// TestUserValidateAbsenseName tests when user name is blank
func TestUserValidateAbsenseName(t *testing.T) {
	user := User{Name: "", Email: "aaa@bbb.com", Password: "Destory_Rome"}
	errorMap := user.Validate()

	if len(errorMap) == 0 {
		t.Errorf(`User{Name:"", Email:"aaa@bbb.com"}`)
	}

	if errorMap["Name"] != "Name is required." {
		t.Errorf(`(User{Name:"", Email:"aaa@bbb.com"}) errorMap['name'] = %q expected %q`, errorMap["Name"], "Name is required.")
	}
}

// TestUserValidateAbsenseName tests absense of email address
func TestUserValidateAbsenseEmail(t *testing.T) {
	user := User{Name: "Brünhild", Email: "", Password: "Destory_Rome"}
	errorMap := user.Validate()

	if len(errorMap) == 0 {
		t.Error(`User{Name : "Brünhild", Email:""}`)
	}

	if errorMap["Email"] != "Email is required." {
		t.Errorf(`User{Name : "Brünhild", Email:""} errorMap['Email'] = %q expected %q `,
			errorMap["Email"], "Email is required.")
	}
}

// TestUserValidateInformalEmail tests informal email input
func TestUserValidateInformalEmail(t *testing.T) {
	user := User{Name: "Brünhild", Email: "asdadsasda", Password: "Destory_Rome"}
	errorMap := user.Validate()
	if len(errorMap) == 0 {
		t.Error(`User{Name : "Brünhild", Email:"asdasdad"})`)
	}

	if errorMap["Email"] != "Invalid email address" {
		t.Errorf(`User{Name : "Brünhild", Email:"asdadsasda"} errorMap['Email'] = %q expected %q `,
			errorMap["Email"], "Invalid email address")
	}
}

// TestUserValidateAbsensePassword tests absense password
func TestUserValidateAbsensePassword(t *testing.T) {
	user := User{Name: "Attilla", Email: "attilla@hun.com"}
	errorMap := user.Validate()

	if len(errorMap) == 0 {
		t.Errorf(`User{Name: "Attilla", Email: "attilla@hun.com"`)
	}

	if errorMap["Password"] != "Password is required" {
		t.Errorf(`User{Name: "Attilla", Email: "attilla@hun.com"} errorMap['Password'] = %q expected %q `,
			errorMap["Password"], "Password is required")
	}
}

// TestUserValidateNormalCondition tests normal condition of User
func TestUserValidateNormalCondition(t *testing.T) {
	user := User{Name: "Attilla", Email: "attilla@hun.com", Password: "Destory_Rome"}
	errorMap := user.Validate()

	if len(errorMap) > 0 {
		t.Errorf(`User{Name: "Attilla", Email: "attilla@hun.com", Password: "Destory_Rome"}`)
	}
}
