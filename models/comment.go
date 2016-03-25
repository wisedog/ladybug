package models

// Comment model
type Comment struct {
	BaseModel

	EntityType		int		// This field indicates which page will this comment be attached. maybe test case, specification...
	EntityID		int
	UserID      	int
	User            User
	Content         string  `sql:"size:1000"`
	
	Spec			Specification
}

/*
func (comment *Comment) Validate(v *revel.Validation) {
	v.Check(comment.Content,
		revel.Required{},
	)
}
*/