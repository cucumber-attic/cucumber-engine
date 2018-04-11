package event

import (
	"encoding/json"

	"github.com/cucumber/cucumber-pickle-runner/src/dto"
	gherkin "github.com/cucumber/gherkin-go"
)

// TestCasePreparedStep is the location information for a step / hook
type TestCasePreparedStep struct {
	SourceLocation *Location `json:"sourceLocation"`
	ActionLocation *Location `json:"actionLocation"`
}

// TestCasePrepared is an event for when a test case has computed what steps and hooks will run
type TestCasePrepared struct {
	SourceLocation *Location
	Steps          []*TestCasePreparedStep
}

// MarshalJSON is the custom JSON marshalling to add the event type
func (t *TestCasePrepared) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		SourceLocation *Location               `json:"sourceLocation"`
		Steps          []*TestCasePreparedStep `json:"steps"`
		Type           string                  `json:"type"`
	}{
		SourceLocation: t.SourceLocation,
		Steps:          t.Steps,
		Type:           "test-case-prepared",
	})
}

// NewTestCasePreparedOptions are the options for NewTestCasePrepared
type NewTestCasePreparedOptions struct {
	Pickle                        *gherkin.Pickle
	URI                           string
	BeforeTestCaseHookDefinitions []*dto.TestCaseHookDefinition
	AfterTestCaseHookDefinitions  []*dto.TestCaseHookDefinition
	StepIndexToStepDefinitions    [][]*dto.StepDefinition
}

// NewTestCasePrepared creates a TestCasePrepared
func NewTestCasePrepared(opts NewTestCasePreparedOptions) *TestCasePrepared {
	var steps []*TestCasePreparedStep
	for _, def := range opts.BeforeTestCaseHookDefinitions {
		steps = append(steps, &TestCasePreparedStep{
			ActionLocation: testCaseHookDefinitionToLocation(def),
		})
	}
	for stepIndex, step := range opts.Pickle.Steps {
		eventStep := &TestCasePreparedStep{
			SourceLocation: pickleStepToLocation(step, opts.URI),
		}
		if len(opts.StepIndexToStepDefinitions[stepIndex]) == 1 {
			eventStep.ActionLocation = stepDefinitionToLocation(opts.StepIndexToStepDefinitions[stepIndex][0])
		}
		steps = append(steps, eventStep)
	}
	for _, def := range opts.AfterTestCaseHookDefinitions {
		steps = append(steps, &TestCasePreparedStep{
			ActionLocation: testCaseHookDefinitionToLocation(def),
		})
	}
	return &TestCasePrepared{
		SourceLocation: pickleToLocation(opts.Pickle, opts.URI),
		Steps:          steps,
	}
}
