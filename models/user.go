package models

import (
	"time"
)

// Constants for Role
const (
	RoleAdmin   = 1 + iota
	RoleManager // Manager role for one or more projects
	RoleUser    // Just a user
	RoleGuest   // Not used now
)

// User is a model represents a user
type User struct {
	BaseModel

	Name           string
	Email          string `sql:"not null;unique"`
	Password       string `json:"-"` // to be removed
	HashedPassword []byte `json:"-"`
	Language       string
	Region         string
	Projects       []Project
	Message        string
	Location       string
	Photo          string // TODO URL
	Role           int
	Notes          string
	//TODO Roles          []Role
	//TODO link of email, homepage, FB, TW, G+ ....

	LastLoginAt       time.Time
	PasswordUpdatedAt time.Time `json:"-"`
}

// Validate check input value and return error map
func (user *User) Validate() map[string]string {
	errorMap := make(map[string]string)

	if !user.Required(user.Name) {
		errorMap["Name"] = "Name is required."
	} else {
		if !user.MaxSize(user.Name, 30) {
			errorMap["Name"] = "Too long name"
		}
	}

	if !user.Required(user.Email) {
		errorMap["Email"] = "Email is required."
	} else {
		if !user.ValidateEmail(user.Email) {
			errorMap["Email"] = "Invalid email address"
		}
	}

	if !user.Required(user.Password) {
		errorMap["Password"] = "Password is required"
	}

	return errorMap

}
