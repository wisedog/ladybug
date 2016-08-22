package models

// Review type
const (
	ReviewTestcase = 1 + iota
	ReviewBug
)

// Review status
const (
	ReviewStatusReady = 1 + iota
	ReviewStatusInProgress
	ReviewStatusReject
	ReviewStatusDone
)

// Review represents review of a test case
type Review struct {
	BaseModel

	Status      int    // Ready,In progress, Reject, Done
	Category    int    // from issue management system or inhouse testcase review
	Comment     string `sql:"size:500;"` // description of the issue
	Description string `sql:"size:500;"`
	ToUser      User
	ToUserID    int
	FromUser    User
	FromUserID  int
}
