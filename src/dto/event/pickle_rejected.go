package event

import (
	"encoding/json"

	gherkin "github.com/cucumber/gherkin-go"
)

// PickleRejected is an event for when a pickle is rejected by the filters
type PickleRejected struct {
	pickleEvent *gherkin.PickleEvent
}

// MarshalJSON is the custom JSON marshalling to add the event type
func (p *PickleRejected) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		URI    string          `json:"uri"`
		Pickle *gherkin.Pickle `json:"pickle"`
		Type   string          `json:"type"`
	}{
		URI:    p.pickleEvent.URI,
		Pickle: p.pickleEvent.Pickle,
		Type:   "pickle-rejected",
	})
}

// NewPickleRejected creates a PickleRejected
func NewPickleRejected(pickleEvent *gherkin.PickleEvent) *PickleRejected {
	return &PickleRejected{pickleEvent: pickleEvent}
}
