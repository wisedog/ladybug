package models
	
const (
    HISTORY_TYPE_TC = 1 + iota
    )

const (
	HISTORY_CHANGE_TYPE_CHANGED = 1 + iota
	HISTORY_CHANGE_TYPE_SET
	HISTORY_CHANGE_TYPE_NOTE
	HISTORY_CHANGE_TYPE_DIFF
	)

// to store multiple changement contents via JSON string
type HistoryTestCaseUnit struct{
	ChangeType	int
	What		string		// thing to be changed
	From		int
	FromStr		string
	To			int
	ToStr		string
	Set			string
	DiffID		int
	Msg			string
}

// History is a model for history of test cases or something stuff.
// 
type History struct {
	BaseModel

	ChangesJson		string	`sql:"size:1000"`
	Changes			[]HistoryTestCaseUnit
	User            User
	UserID      	int
	Category        int // Testcase or .... 
	TargetID		int
	ChangeType		int		// changed or set 
	What			string
	From			int
	FromStr			string
	To				int
	ToStr			string
	Set				string
	// TODO find diff library and apply. the other fileds are depend on that.
	DiffID			int
	Note			string	`sql:"size:512"`

}
