package dto

import tagexpressions "github.com/cucumber/tag-expressions-go"

// TestCaseHookDefinitionConfig is hook that run before or after a test case
type TestCaseHookDefinitionConfig struct {
	ID            string `json:"id"`
	TagExpression string `json:"tagExpression"`
	URI           string `json:"uri"`
	Line          int    `json:"line"`
}

// TestCaseHookDefinition is a TestCaseHookDefinitionConfig where the
// tag expression has been converted to a TagExpression
type TestCaseHookDefinition struct {
	ID            string
	TagExpression tagexpressions.Evaluatable
	URI           string
	Line          int
}

// Location returns a Location for the test case hook definition
func (t *TestCaseHookDefinition) Location() *Location {
	return &Location{
		URI:  t.URI,
		Line: t.Line,
	}
}
