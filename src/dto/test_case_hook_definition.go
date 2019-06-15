package dto

import (
	messages "github.com/cucumber/cucumber-messages-go/v3"
	tagexpressions "github.com/cucumber/tag-expressions-go"
)

// TestCaseHookDefinition wraps a TestCaseHookDefinitionConfig where the
// tag expression has been converted to a TagExpression
type TestCaseHookDefinition struct {
	Config        *messages.TestCaseHookDefinitionConfig
	TagExpression tagexpressions.Evaluatable
}
