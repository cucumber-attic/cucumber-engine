package runner

import (
	"regexp"

	cucumberexpressions "github.com/cucumber/cucumber-expressions-go"
	"github.com/cucumber/cucumber-pickle-runner/src/dto"
	tagexpressions "github.com/cucumber/tag-expressions-go"
)

// SupportCodeLibrary represents the support code for the test run
type SupportCodeLibrary struct {
	afterTestCaseHookDefinitions  []*dto.TestCaseHookDefinition
	beforeTestCaseHookDefinitions []*dto.TestCaseHookDefinition
	parameterTypeRegistry         *cucumberexpressions.ParameterTypeRegistry
	stepDefinitions               []*dto.StepDefinition
}

// NewSupportCodeLibrary returns a SupportCodeLibrary for the given config
func NewSupportCodeLibrary(config *dto.SupportCodeConfig) (*SupportCodeLibrary, error) {
	afterTestCaseHookDefinitions, err := createTestCaseHookDefinitions(config.AfterTestCaseHookDefinitionConfigs)
	if err != nil {
		return nil, err
	}
	beforeTestCaseHookDefinitions, err := createTestCaseHookDefinitions(config.BeforeTestCaseHookDefinitionConfigs)
	if err != nil {
		return nil, err
	}
	parameterTypeRegistry, err := createParameterTypeRegistry(config.ParameterTypeConfigs)
	if err != nil {
		return nil, err
	}
	stepDefinitions, err := createStepDefinitions(config.StepDefinitionConfigs, parameterTypeRegistry)
	if err != nil {
		return nil, err
	}
	return &SupportCodeLibrary{
		afterTestCaseHookDefinitions:  afterTestCaseHookDefinitions,
		beforeTestCaseHookDefinitions: beforeTestCaseHookDefinitions,
		parameterTypeRegistry:         parameterTypeRegistry,
		stepDefinitions:               stepDefinitions,
	}, nil
}

// GetMatchingAfterTestCaseHookDefinitions returns the TestCaseHookDefinition that match the given tag names
func (s *SupportCodeLibrary) GetMatchingAfterTestCaseHookDefinitions(tagNames []string) []*dto.TestCaseHookDefinition {
	return filterHookDefinitions(s.afterTestCaseHookDefinitions, tagNames)
}

// GetMatchingBeforeTestCaseHookDefinitions returns the TestCaseHookDefinition that match the given tag names
func (s *SupportCodeLibrary) GetMatchingBeforeTestCaseHookDefinitions(tagNames []string) []*dto.TestCaseHookDefinition {
	return filterHookDefinitions(s.beforeTestCaseHookDefinitions, tagNames)
}

// GetMatchingStepDefinitions returns the StepDefinitions that match the given text
//   the pattern matches are only returned if a single step definition matches
func (s *SupportCodeLibrary) GetMatchingStepDefinitions(text string) ([]*dto.StepDefinition, []*dto.PatternMatch, error) {
	stepDefinitions := []*dto.StepDefinition{}
	var patternMatches []*dto.PatternMatch
	for _, def := range s.stepDefinitions {
		args, err := def.Expression.Match(text)
		if err != nil {
			return nil, nil, err
		}
		if args == nil {
			continue
		}
		stepDefinitions = append(stepDefinitions, def)
		if len(stepDefinitions) == 1 {
			patternMatches = make([]*dto.PatternMatch, len(args))
			for i, arg := range args {
				patternMatches[i] = &dto.PatternMatch{
					Captures:          arg.Group().Values(),
					ParameterTypeName: arg.ParameterType().Name(),
				}
			}
		} else {
			patternMatches = nil
		}
	}
	return stepDefinitions, patternMatches, nil
}

// GenerateExpressions returns the generated expressions for an undefined step
func (s *SupportCodeLibrary) GenerateExpressions(text string) []*dto.GeneratedExpression {
	generator := cucumberexpressions.NewCucumberExpressionGenerator(s.parameterTypeRegistry)
	expressions := generator.GenerateExpressions(text)
	result := make([]*dto.GeneratedExpression, len(expressions))
	for i, expression := range expressions {
		parameterTypeNames := make([]string, len(expression.ParameterTypes()))
		for j, parameterType := range expression.ParameterTypes() {
			parameterTypeNames[j] = parameterType.Name()
		}
		result[i] = &dto.GeneratedExpression{
			Text:               expression.Source(),
			ParameterTypeNames: parameterTypeNames,
		}
	}
	return result
}

func filterHookDefinitions(hookDefinitions []*dto.TestCaseHookDefinition, tagNames []string) []*dto.TestCaseHookDefinition {
	result := []*dto.TestCaseHookDefinition{}
	for _, hookDefinition := range hookDefinitions {
		if hookDefinition.TagExpression.Evaluate(tagNames) {
			result = append(result, hookDefinition)
		}
	}
	return result
}

func createParameterTypeRegistry(parameterTypeConfigs []*dto.ParameterTypeConfig) (*cucumberexpressions.ParameterTypeRegistry, error) {
	parameterTypeRegistry := cucumberexpressions.NewParameterTypeRegistry()
	for _, parameterTypeConfig := range parameterTypeConfigs {
		regexps := make([]*regexp.Regexp, len(parameterTypeConfig.Regexps))
		for i, regexpSource := range parameterTypeConfig.Regexps {
			var err error
			regexps[i], err = regexp.Compile(regexpSource)
			if err != nil {
				// TODO wrap error with parameterType name
				return nil, err
			}
		}
		parameterType, err := cucumberexpressions.NewParameterType(
			parameterTypeConfig.Name,
			regexps,
			"",
			nil,
			parameterTypeConfig.UseForSnippets,
			parameterTypeConfig.PreferForRegexpMatch,
		)
		if err != nil {
			// TODO wrap error with parameterType name
			return nil, err
		}
		err = parameterTypeRegistry.DefineParameterType(parameterType)
		if err != nil {
			// TODO wrap error with parameterType name
			return nil, err
		}
	}
	return parameterTypeRegistry, nil
}

func createTestCaseHookDefinitions(configs []*dto.TestCaseHookDefinitionConfig) ([]*dto.TestCaseHookDefinition, error) {
	result := make([]*dto.TestCaseHookDefinition, len(configs))
	for i, config := range configs {
		tagExpression, err := tagexpressions.Parse(config.TagExpression)
		if err != nil {
			// TODO wrap error with tag expression and line / uri
			return nil, err
		}
		result[i] = &dto.TestCaseHookDefinition{
			ID:            config.ID,
			Line:          config.Line,
			TagExpression: tagExpression,
			URI:           config.URI,
		}
	}
	return result, nil
}

func createStepDefinitions(configs []*dto.StepDefinitionConfig, parameterTypeRegistry *cucumberexpressions.ParameterTypeRegistry) ([]*dto.StepDefinition, error) {
	result := make([]*dto.StepDefinition, len(configs))
	for i, config := range configs {
		expression, err := config.Pattern.Expression(parameterTypeRegistry)
		if err != nil {
			// TODO wrap error with pattern and line / uri
			return nil, err
		}
		result[i] = &dto.StepDefinition{
			Expression: expression,
			ID:         config.ID,
			Line:       config.Line,
			URI:        config.URI,
		}
	}
	return result, nil
}
