package models

import (
	"github.com/revel/revel"
	"time"
)


// Activity only belongs to user
type Activity struct {
	ID				int
	User            User
	UserID      	int
	
	Content         string  `sql:"size:1000"`
	
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

func (activity *Activity) Validate(v *revel.Validation) {
	v.Check(activity.UserID,
		revel.Required{},
	)
}
