package models

// JSTreeTestSuiteWrapper is a wrapper type for print out JSON format used in JSTree
type JSTreeTestSuiteWrapper struct {
	Text     string `json:"text"`
	ID       string `json:"id"`
	Type     string `json:"type"`
	Children bool   `json:"children"`
}

// JSTreeTestSuiteWrapperWithChildren is a wrapper type for print out JSON format
// used in children node of JTree
type JSTreeTestSuiteWrapperWithChildren struct {
	Text     string                  `json:"text"`
	ID       string                  `json:"id"`
	Type     string                  `json:"type"`
	Children []JSTreeTestCaseWrapper `json:"children"`
}

//JSTreeTestCaseWrapper is a wrapper type for print out JSON format used in JSTree
type JSTreeTestCaseWrapper struct {
	Text string `json:"text"`
	ID   string `json:"id"`
	Type string `json:"type"`
}

// JSTreeNode represents a node of JSTree
type JSTreeNode struct {
	Text   string `json:"text"`
	ID     string `json:"id"`
	Type   string `json:"type"`
	Parent string `json:"parent"`
}
