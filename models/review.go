package models

import (
	"time"
)

const (
    ReviewTestcase = 1 + iota
    ReviewBug
    )

const (
    ReviewStatusReady = 1 + iota
    ReviewStatusInProgress
    ReviewStatusReject
    ReviewStatusDone
    )

type Review struct {
	ID          int
	Status      int     // Ready,In progress, Reject, Done
	Category    int     // from issue management system or inhouse testcase review
	Comment     string `sql:"size:500;"` // description of the issue
	Description string  `sql:"size:500;"`
	ToUser      User
	ToUserID    int
	FromUser    User
	FromUserID  int

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
