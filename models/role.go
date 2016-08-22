package models

// Constants for Role
const (
	RoleAdmin   = 1 + iota
	RoleManager // Manager role for one or more projects
	RoleUser    // Just a user
	RoleGuest   // Not used now
)

// Role is a model represents a user
// This table contains user's role of each project
type Role struct {
	BaseModel

	ProjectID int
	Project   Project
	UserID    int
	User      User
	UserRole  int
}
