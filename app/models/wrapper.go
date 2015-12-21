package models

import ()

// Wrapper type for print out JSON format used in JSTree
type JSTreeTestSuiteWrapper struct {
	Text     string `json:"text"`
	Id       string `json:"id"`
	Type     string `json:"type"`
	Children bool   `json:"children"`
}

type JSTreeTestSuiteWrapperWithChildren struct {
	Text     string                  `json:"text"`
	Id       string                  `json:"id"`
	Type     string                  `json:"type"`
	Children []JSTreeTestCaseWrapper `json:"children"`
}

// Wrapper type for print out JSON format used in JSTree
type JSTreeTestCaseWrapper struct {
	Text string `json:"text"`
	Id   string `json:"id"`
	Type string `json:"type"`
}

type JSTreeNode struct {
	Text   string `json:"text"`
	Id     string `json:"id"`
	Type   string `json:"type"`
	Parent string `json:"parent"`
}
