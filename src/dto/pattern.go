package dto

import (
	"fmt"
	"regexp"

	cucumberexpressions "github.com/cucumber/cucumber-expressions-go"
	messages "github.com/cucumber/cucumber-messages-go/v2"
)

// GetExpression returns the cucumber expression this pattern defines
func GetExpression(pattern *messages.StepDefinitionPattern, parameterTypeRegistry *cucumberexpressions.ParameterTypeRegistry) (cucumberexpressions.Expression, error) {
	switch pattern.Type {
	case messages.StepDefinitionPatternType_CUCUMBER_EXPRESSION:
		return cucumberexpressions.NewCucumberExpression(pattern.Source, parameterTypeRegistry)
	case messages.StepDefinitionPatternType_REGULAR_EXPRESSION:
		r, err := regexp.Compile(pattern.Source)
		if err != nil {
			return nil, err
		}
		return cucumberexpressions.NewRegularExpression(r, parameterTypeRegistry), nil
	default:
		return nil, fmt.Errorf(
			"Unexpected pattern type: `%s`. Should be `%s` or `%s`",
			pattern.Type,
			messages.StepDefinitionPatternType_CUCUMBER_EXPRESSION,
			messages.StepDefinitionPatternType_REGULAR_EXPRESSION)
	}
}
