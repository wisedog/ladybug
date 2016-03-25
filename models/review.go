package models

import (
	"time"
)

const (
    REVIEW_TESTCASE = 1 + iota
    REVIEW_BUG
    )

const (
    REVIEW_STATUS_READY = 1 + iota
    REVIEW_STATUS_IN_PROGRESS
    REVIEW_STATUS_REJECT
    REVIEW_STATUS_DONE
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
