package models

import (
	"errors"
)

// Kind of this object
const (
	TcReqRelationHistoryLink = 1 + iota
	TcReqRelationHistoryUnlink
)

// TcReqRelationHistory represents the history of relationship between TestCase and Requirement
type TcReqRelationHistory struct {
	BaseModel

	// Indicates this object kind. The kind may be link or unlink
	Kind          int
	RequirementID int
	TestCaseID    int
	ProjectID     int
}

// Validate check input value and return error
//
// if returned value 'error' is not nil,
// you may consider invoke 500 internal error.
func (relation *TcReqRelationHistory) Validate() error {
	var err error
	if relation.RequirementID == 0 {
		err = errors.New("Invalid requirement id")
	}

	if relation.TestCaseID == 0 {
		err = errors.New("Invalid testcase id")
	}

	return err
}
