package models

// Activity only belongs to user
type Activity struct {
	BaseModel
	User            User
	UserID      	int
	
	Content         string  `sql:"size:1000"`
}

/*func (activity *Activity) Validate(v *revel.Validation) {
	v.Check(activity.UserID,
		revel.Required{},
	)
}
*/