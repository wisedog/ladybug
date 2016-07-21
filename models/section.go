package models

import (
	"errors"
)

// Section represents section of test case or requirements
type Section struct {
	BaseModel
	Prefix      string
	Seq         int // may use...
	Title       string
	Description string // may not used
	Status      int    // may not use
	ParentsID   int
	ProjectID   int
	RootNode    bool
	ForTestCase bool //TestCase or Specification
}

// Validate check input value and return error map
//
// if second return value 'error' is not nil,
// you may consider invoke 500 internal error.
func (section *Section) Validate() (map[string]string, error) {
	errorMap := make(map[string]string)
	if !section.Required(section.Title) {
		errorMap["Title"] = "Title is required."
	}
	// TODO max size of title
	// TODO check prefix and size

	var err error
	if section.ProjectID == 0 {
		err = errors.New("Invalid project id")
	}

	return errorMap, err
}
