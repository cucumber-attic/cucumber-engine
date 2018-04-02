package event

import (
	"github.com/cucumber/cucumber-pickle-runner/src/dto"
	gherkin "github.com/cucumber/gherkin-go"
)

// TestCase is the location information for a test case
type TestCase struct {
	SourceLocation *Location `json:"sourceLocation"`
}

// Location is location information within a file
type Location struct {
	Line int    `json:"line"`
	URI  string `json:"uri"`
}

func testCaseHookDefinitionToLocation(def *dto.TestCaseHookDefinition) *Location {
	return &Location{
		URI:  def.URI,
		Line: def.Line,
	}
}

func stepDefinitionToLocation(def *dto.StepDefinition) *Location {
	return &Location{
		URI:  def.URI,
		Line: def.Line,
	}
}

func pickleToLocation(pickle *gherkin.Pickle, uri string) *Location {
	return &Location{
		URI:  uri,
		Line: pickle.Locations[0].Line,
	}
}

func pickleStepToLocation(step *gherkin.PickleStep, uri string) *Location {
	return &Location{
		URI:  uri,
		Line: step.Locations[len(step.Locations)-1].Line,
	}
}
