package dto

import cucumberexpressions "github.com/cucumber/cucumber-expressions-go"

// StepDefinitionConfig is the implementation of a step
type StepDefinitionConfig struct {
	ID      string  `json:"id"`
	Pattern Pattern `json:"pattern"`
	URI     string  `json:"uri"`
	Line    int     `json:"line"`
}

// StepDefinition is a StepDefinitionConfig where the patten has been
// converted to a cucumber expression
type StepDefinition struct {
	ID         string
	Expression cucumberexpressions.Expression
	URI        string
	Line       int
}
