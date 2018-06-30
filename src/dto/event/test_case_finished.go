package event

import (
	"encoding/json"

	"github.com/cucumber/cucumber-engine/src/dto"
	gherkin "github.com/cucumber/gherkin-go"
)

// TestCaseFinished is an event for when a test case finishes running
type TestCaseFinished struct {
	Result         *dto.TestResult
	SourceLocation *dto.Location
}

// MarshalJSON is the custom JSON marshalling to add the event type
func (t *TestCaseFinished) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Result         *dto.TestResult `json:"result"`
		SourceLocation *dto.Location   `json:"sourceLocation"`
		Type           string          `json:"type"`
	}{
		Result:         t.Result,
		SourceLocation: t.SourceLocation,
		Type:           "test-case-finished",
	})
}

// NewTestCaseFinishedOptions are the options for NewTestCaseFinished
type NewTestCaseFinishedOptions struct {
	Pickle *gherkin.Pickle
	Result *dto.TestResult
	URI    string
}

// NewTestCaseFinished creates a TestStepFinished
func NewTestCaseFinished(opts NewTestCaseFinishedOptions) *TestCaseFinished {
	return &TestCaseFinished{
		SourceLocation: dto.NewLocationForPickle(opts.Pickle, opts.URI),
		Result:         opts.Result,
	}
}
