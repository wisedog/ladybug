package models

const (
	SpecStatusActivate = 1 + iota
	SpecStatusInactivate
	SpecStatusDraft
	)

// Specification Model
type Specification struct {
  BaseModel
	Name            string
	Status			int
	Priority		int
	SectionID		int
}
