package event

import (
	"encoding/json"

	gherkin "github.com/cucumber/gherkin-go"
)

// TestCaseStarted is an event for when a test case starts running
type TestCaseStarted struct {
	SourceLocation Location `json:"source_location"`
}

// MarshalJSON is the custom JSON marshalling to add the event type
func (t *TestCaseStarted) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		SourceLocation Location `json:"sourceLocation"`
		Type           string   `json:"type"`
	}{
		SourceLocation: t.SourceLocation,
		Type:           "test-case-finished",
	})
}

// NewTestCaseStartedOptions are the options for NewTestCaseStarted
type NewTestCaseStartedOptions struct {
	Pickle *gherkin.Pickle
	URI    string
}

// NewTestCaseStarted creates a TestCaseStarted
func NewTestCaseStarted(opts NewTestCaseStartedOptions) *TestCaseStarted {
	return &TestCaseStarted{
		SourceLocation: pickleToLocation(opts.Pickle, opts.URI),
	}
}
