package models

import (
	"fmt"
	"github.com/revel/revel"
	"regexp"
	"time"
)

/*
type Role int

const (
	USER
	MODERATOR
	MANAGER
	ADMIN
)*/

type User struct {
	ID             int
	Name           string
	Email          string `sql:"not null;unique"`
	Password       string
	HashedPassword []byte
	Language       string
	Region         string
	Projects       []Project
	Message        string
	Location       string
	Photo          string
	//TODO Roles          []Role
	//TODO link of email, homepage, FB, TW, G+ ....

	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         time.Time
	LastLoginAt       time.Time
	PasswordUpdatedAt time.Time
}

func (u *User) String() string {
	return fmt.Sprintf("User(%s)", u.Email)
}

var userRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)

func (user *User) Validate(v *revel.Validation) {
	v.Email(user.Email)
	/*	v.Check(user.Email,
		revel.Required{},
		revel.MaxSize{60},
		revel.MinSize{6},
		revel.Match{userRegex},
	)*/

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
