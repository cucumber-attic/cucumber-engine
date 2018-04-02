package dto

import (
	gherkin "github.com/cucumber/gherkin-go"
)

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
	FeaturesConfig *FeaturesConfig `json:"featuresConfig"`
	RuntimeConfig  *RuntimeConfig  `json:"runtimeConfig"`

	// Used for type "run before/after test case hook"
	TestCaseHookID string `json:"testCaseHookId"`

	// Used for type "run test step"
	TestStep *gherkin.PickleStep `json:"testStep"`

	// Used for type "event"
	Event interface{} `json:"event"`
}
