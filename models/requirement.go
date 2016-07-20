package models

import "errors"

// Status of requirements
const (
	ReqStatusActivate = 1 + iota
	ReqStatusInactivate
	ReqStatusDraft
)

// Requirement Model
type Requirement struct {
	BaseModel
	Name      string
	Status    int
	Priority  int
	SectionID int
}

// Validate check input value and return error map
//
// if second return value 'error' is not nil,
// you may consider invoke 500 internal error.
func (req *Requirement) Validate() (map[string]string, error) {
	errorMap := make(map[string]string)
	if !req.Required(req.Name) {
		errorMap["Name"] = "Name is required."
	}
	var err error

	if req.SectionID == 0 {
		err = errors.New("Invalid section id")
	}

	if req.Status < ReqStatusActivate || req.Status > ReqStatusInactivate {
		err = errors.New("Invalid requirement status")
	}
	return errorMap, err
}
