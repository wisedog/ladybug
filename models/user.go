package models

import (
	"time"
)

const (
	RoleAdmin = 1 + iota
	RoleManager
	RoleUser
	RoleGuest
)

type User struct {
	BaseModel
	Name           string
	Email          string `sql:"not null;unique"`
	Password       string
	HashedPassword []byte
	Language       string
	Region         string
	Projects       []Project
	Message        string
	Location       string
	Photo          string	// TODO URL
	Roles			int
	Notes			string
	//TODO Roles          []Role
	//TODO link of email, homepage, FB, TW, G+ ....

	LastLoginAt       time.Time
	PasswordUpdatedAt time.Time
}

/*func (u *User) String() string {
	return fmt.Sprintf("User(%s)", u.Email)
}

var userRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)

func (user *User) Validate(v *revel.Validation) {
	v.Email(user.Email)

	ValidatePassword(v, user.Password).
		Key("user.Password")

	v.Check(user.Name,
		revel.Required{},
		revel.MaxSize{25},
	)
}

func ValidatePassword(v *revel.Validation, password string) *revel.ValidationResult {
	return v.Check(password,
		revel.Required{},
		revel.MaxSize{20},
		revel.MinSize{2},
	)
}
*/