package models


import (
	"time"
)


const (
	SPEC_STATUS_ACTIVATE = 1 + iota
	SPEC_STATUS_INACTIVE
	SPEC_STATUS_DRAFT
	)


/*
Specification
*/
type Specification struct {
	ID              int
	Name            string
	Status			int
	Priority		int
	SectionID		int
	
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
