package event

import (
	"github.com/cucumber/cucumber-engine/src/dto"
	messages "github.com/cucumber/cucumber-messages-go/v3"
)

// NewTestCasePreparedOptions are the options for NewTestCasePrepared
type NewTestCasePreparedOptions struct {
	Pickle                        *messages.Pickle
	BeforeTestCaseHookDefinitions []*dto.TestCaseHookDefinition
	AfterTestCaseHookDefinitions  []*dto.TestCaseHookDefinition
	StepIndexToStepDefinitions    [][]*dto.StepDefinition
}

// NewTestCasePrepared creates a TestCasePrepared
func NewTestCasePrepared(opts NewTestCasePreparedOptions) *messages.TestCasePrepared {
	var steps []*messages.TestCasePreparedStep
	for _, def := range opts.BeforeTestCaseHookDefinitions {
		steps = append(steps, &messages.TestCasePreparedStep{
			ActionLocation: def.Config.Location,
		})
	}
	for stepIndex, step := range opts.Pickle.Steps {
		eventStep := &messages.TestCasePreparedStep{
			SourceLocation: &messages.SourceReference{
				Uri:      opts.Pickle.Uri,
				Location: step.Locations[len(step.Locations)-1],
			},
		}
		if len(opts.StepIndexToStepDefinitions[stepIndex]) == 1 {
			eventStep.ActionLocation = opts.StepIndexToStepDefinitions[stepIndex][0].Config.Location
		}
		steps = append(steps, eventStep)
	}
	for _, def := range opts.AfterTestCaseHookDefinitions {
		steps = append(steps, &messages.TestCasePreparedStep{
			ActionLocation: def.Config.Location,
		})
	}
	return &messages.TestCasePrepared{
		PickleId: opts.Pickle.Name, // TODO fix
		Steps:    steps,
	}
}
