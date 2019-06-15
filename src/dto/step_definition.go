package dto

import (
	cucumberexpressions "github.com/cucumber/cucumber-expressions-go"
	messages "github.com/cucumber/cucumber-messages-go/v3"
)

// StepDefinition wraps a StepDefinitionConfig where the patten has been
// converted to a cucumber expression
type StepDefinition struct {
	Config     *messages.StepDefinitionConfig
	Expression cucumberexpressions.Expression
}
