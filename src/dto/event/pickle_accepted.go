package event

import (
	"encoding/json"

	gherkin "github.com/cucumber/gherkin-go"
)

// PickleAccepted is an event for when a pickle is accepted by the filters
type PickleAccepted struct {
	pickleEvent *gherkin.PickleEvent
}

// MarshalJSON is the custom JSON marshalling to add the event type
func (p *PickleAccepted) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		URI    string          `json:"uri"`
		Pickle *gherkin.Pickle `json:"pickle"`
		Type   string          `json:"type"`
	}{
		URI:    p.pickleEvent.URI,
		Pickle: p.pickleEvent.Pickle,
		Type:   "pickle-accepted",
	})
}

// NewPickleAccepted creates a PickleAccepted
func NewPickleAccepted(pickleEvent *gherkin.PickleEvent) *PickleAccepted {
	return &PickleAccepted{pickleEvent: pickleEvent}
}
