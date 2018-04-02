package event

import "encoding/json"

// TestRunStarted is an event for when the test run starts
type TestRunStarted struct{}

// MarshalJSON is the custom JSON marshalling to add the event type
func (t *TestRunStarted) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Type string `json:"type"`
	}{
		Type: "test-run-started",
	})
}

// NewTestRunStarted creates a TestRunStarted
func NewTestRunStarted() *TestRunStarted {
	return &TestRunStarted{}
}
