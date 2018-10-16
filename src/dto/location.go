package dto

import "github.com/cucumber/gherkin-go"

// Location is location information within a file
type Location struct {
	Line int    `json:"line"`
	URI  string `json:"uri"`
}

// NewLocationForPickle returns a Location for the given pickle and uri
func NewLocationForPickle(pickle *gherkin.Pickle, uri string) *Location {
	return &Location{
		URI:  uri,
		Line: pickle.Locations[0].Line,
	}
}

// NewLocationForPickleStep returns a Location for the given pickle step and uri
func NewLocationForPickleStep(step *gherkin.PickleStep, uri string) *Location {
	return &Location{
		URI:  uri,
		Line: step.Locations[len(step.Locations)-1].Line,
	}
}
