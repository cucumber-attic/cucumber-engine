package dto

import gherkin "github.com/cucumber/gherkin-go"

// Command is a the struct used to communicate between this process and the calling process
type Command struct {
	ID   string      `json:"id"`
	Type CommandType `json:"type"`

	// Used for type "action complete"
	ResponseTo string `json:"responseTo"`

	// Used for type "action complete" when action was
	//   "run before/after test case hook" or "run test step"
	HookOrStepResult *TestResult `json:"hookOrStepResult"`

	// Used for type "start"
	FeaturesConfig    *FeaturesConfig    `json:"featuresConfig"`
	RuntimeConfig     *RuntimeConfig     `json:"runtimeConfig"`
	SupportCodeConfig *SupportCodeConfig `json:"supportCodeConfig"`

	// Used for type "initialize_test_case", "run before/after test case hook",
	// and "run test step"
	TestCaseID string `json:"testCaseId"`

	// Used for type "run before/after test case hook"
	TestCaseHookDefinitionID string `json:"testCaseHookDefinitionId"`

	// Used for type "run test step"
	StepDefinitionID string          `json:"stepDefinitionId"`
	PatternMatches   []*PatternMatch `json:"patternMatches"`

	// Used for type "run test step" and "generate snippet"
	PickleArguments []gherkin.Argument `json:"pickleArguments"`

	// Used for type "generate snippet"
	GeneratedExpressions []*GeneratedExpression `json:"generateExpression"`

	// Used for type "action complete" when action was "generate_snippet"
	Snippet string `json:"string"`

	// Used for type "event"
	Event interface{} `json:"event"`

	// Used for type "error"
	Error string `json:"error"`
}
