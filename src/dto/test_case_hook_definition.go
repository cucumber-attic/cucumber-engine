package dto

// TestCaseHookDefinition is hook that run before or after a test case
type TestCaseHookDefinition struct {
	ID            string `json:"id"`
	TagExpression string `json:"tagExpression"`
	URI           string `json:"uri"`
	Line          int    `json:"line"`
}
