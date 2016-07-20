package models

// Activity only belongs to user
type Activity struct {
	BaseModel

	User    User
	UserID  int
	Content string `sql:"size:1000"`
}

// Validate function checks all input values are validated
func (activity *Activity) Validate() error {
	if activity.UserID == 0 {
		//TODO make an error
	}

	return nil
}
