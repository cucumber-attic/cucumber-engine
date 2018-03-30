package dto

// TestCaseHookDefinition is hook that run before or after a test case
type TestCaseHookDefinition struct {
	ID            string `json:"id"`
	TagExpression string `json:"tagExpression"`
	URI           string `json:"uri"`
	Line          int    `json:"line"`
}

// StepDefinition is the implementation of a step
type StepDefinition struct {
	ID          string `json:"id"`
	Pattern     string `json:"pattern"`
	PatternType string `json:"patternType"`
	URI         string `json:"uri"`
	Line        int    `json:"line"`
}

// ParameterType is the a configuration for cucumber expressions
type ParameterType struct {
	Name    string   `json:"name"`
	Regexps []string `json:"regexps"`
}

// RuntimeConfig is the configuration for the run
type RuntimeConfig struct {
	IsFailFast                    bool                      `json:"isFailFast"`
	IsDryRun                      bool                      `json:"isDryRun"`
	IsStrict                      bool                      `json:"isStrict"`
	BeforeTestCaseHookDefinitions []*TestCaseHookDefinition `json:"beforeTestCaseHookDefinitions"`
	AfterTestCaseHookDefinitions  []*TestCaseHookDefinition `json:"afterTestCaseHookDefinitions"`
	StepDefinitions               []*StepDefinition         `json:"stepDefinitions"`
	ParameterTypes                []*ParameterType          `json:"parameterTypes"`
}
