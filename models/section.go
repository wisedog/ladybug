package models

type Section struct {
	BaseModel
	Prefix      string
	Seq         int // may use...
	Title       string
	Description string // may not used
	Status      int	// may not use
	ParentsID   int
	ProjectID   int
	RootNode    bool
	ForTestCase	bool	//TestCase or Specification
}
/*
func (section *Section) Validate(v *revel.Validation) {
	v.Check(section.Title,
		revel.Required{},
		revel.MaxSize{400},
	)

	v.Check(section.Prefix,
		revel.Required{},
		revel.MaxSize{16},
	)
}
*/