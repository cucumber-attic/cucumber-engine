package dto

import (
	"fmt"
	"regexp"

	cucumberexpressions "github.com/cucumber/cucumber-expressions-go"
)

// Pattern is how the step definition matches text
type Pattern struct {
	Source string      `json:"source"`
	Type   PatternType `json:"type"`
}

// Expression returns the cucumber expression this pattern defines
func (p Pattern) Expression(parameterTypeRegistry *cucumberexpressions.ParameterTypeRegistry) (cucumberexpressions.Expression, error) {
	switch p.Type {
	case PatternTypeCucumberExpression:
		return cucumberexpressions.NewCucumberExpression(p.Source, parameterTypeRegistry)
	case PatternTypeRegularExpression:
		r, err := regexp.Compile(p.Source)
		if err != nil {
			return nil, err
		}
		return cucumberexpressions.NewRegularExpression(r, parameterTypeRegistry), nil
	default:
		return nil, fmt.Errorf("Unexpected pattern type: `%s`. Should be `%s` or `%s`", p.Type, PatternTypeCucumberExpression, PatternTypeRegularExpression)
	}
}
