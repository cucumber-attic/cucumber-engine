package event

import (
	"encoding/json"

	"github.com/cucumber/cucumber-pickle-runner/src/dto"
)

// TestRunFinished is an event for when the test run finishes
type TestRunFinished struct {
	Result *dto.TestRunResult
}

// MarshalJSON is the custom JSON marshalling to add the event type
func (t *TestRunFinished) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Result *dto.TestRunResult `json:"result"`
		Type   string             `json:"type"`
	}{
		Result: t.Result,
		Type:   "test-run-finished",
	})
}

// NewTestRunFinished creates a TestRunFinished
func NewTestRunFinished(result *dto.TestRunResult) *TestRunFinished {
	return &TestRunFinished{
		Result: result,
	}
}
