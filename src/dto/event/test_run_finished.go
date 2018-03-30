package event

import "encoding/json"

// TestRunFinished is an event for when the test run finishes
type TestRunFinished struct {
	Success bool `json:"success"`
}

// MarshalJSON is the custom JSON marshalling to add the event type
func (t *TestRunFinished) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Success bool   `json:"success"`
		Type    string `json:"type"`
	}{
		Success: t.Success,
		Type:    "test-run-finished",
	})
}

// NewTestRunFinished creates a TestRunFinished
func NewTestRunFinished(success bool) *TestRunFinished {
	return &TestRunFinished{
		Success: success,
	}
}
