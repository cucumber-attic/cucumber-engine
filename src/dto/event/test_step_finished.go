package event

import (
	"encoding/json"

	"github.com/cucumber/cucumber-pickle-runner/src/dto"
	gherkin "github.com/cucumber/gherkin-go"
)

// TestStepFinished is an event for when a test step finishes running
type TestStepFinished struct {
	Index    int
	Result   *dto.TestResult
	TestCase TestCase
}

// MarshalJSON is the custom JSON marshalling to add the event type
func (t *TestStepFinished) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Index    int             `json:"index"`
		Result   *dto.TestResult `json:"result"`
		TestCase TestCase        `json:"testCase"`
		Type     string          `json:"type"`
	}{
		Index:    t.Index,
		TestCase: t.TestCase,
		Result:   t.Result,
		Type:     "test-step-finished",
	})
}

// NewTestStepFinishedOptions are the options for NewTestStepFinished
type NewTestStepFinishedOptions struct {
	Index  int
	Pickle *gherkin.Pickle
	Result *dto.TestResult
	URI    string
}

// NewTestStepFinished creates a TestStepFinished
func NewTestStepFinished(opts NewTestStepFinishedOptions) *TestStepFinished {
	return &TestStepFinished{
		Index: opts.Index,
		TestCase: TestCase{
			SourceLocation: pickleToLocation(opts.Pickle, opts.URI),
		},
		Result: opts.Result,
	}
}
