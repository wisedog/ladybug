package models

import (
	"github.com/revel/revel"
	"time"
)

// Comment model
type Comment struct {
	ID				int
	EntityType		int		// This field indicates which page will this comment be attached. maybe test case, specification...
	EntityID		int
	UserID      	int
	User            User
	Content         string  `sql:"size:1000"`
	
	Spec			Specification
	
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

func (comment *Comment) Validate(v *revel.Validation) {
	v.Check(comment.Content,
		revel.Required{},
	)
}
