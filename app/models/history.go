package models

import (
	"time"
)
	
const (
    HISTORY_TYPE_TC = 1 + iota
    )

/*
  A model for history of test cases or something stuff.
*/
type History struct {
	ID              int
	Category        int // Testcase or .... 
	// TODO find diff library and apply. the other fileds are depend on that.


	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
