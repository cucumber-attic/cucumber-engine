package event

import (
	"encoding/json"

	"github.com/cucumber/cucumber-pickle-runner/src/dto"
	gherkin "github.com/cucumber/gherkin-go"
)

// TestStepStarted is an event for when a test step starts running
type TestStepStarted struct {
	Index    int
	TestCase *dto.TestCase
}

// MarshalJSON is the custom JSON marshalling to add the event type
func (t *TestStepStarted) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Index    int           `json:"index"`
		TestCase *dto.TestCase `json:"testCase"`
		Type     string        `json:"type"`
	}{
		Index:    t.Index,
		TestCase: t.TestCase,
		Type:     "test-step-started",
	})
}

// NewTestStepStartedOptions are the options for NewTestStepStarted
type NewTestStepStartedOptions struct {
	Index  int
	Pickle *gherkin.Pickle
	URI    string
}

// NewTestStepStarted creates a TestStepStarted
func NewTestStepStarted(opts NewTestStepStartedOptions) *TestStepStarted {
	return &TestStepStarted{
		Index: opts.Index,
		TestCase: &dto.TestCase{
			SourceLocation: dto.NewLocationForPickle(opts.Pickle, opts.URI),
		},
	}
}
