package runner

import (
	"fmt"

	"github.com/cucumber/cucumber-engine/src/dto"
	"github.com/cucumber/cucumber-engine/src/dto/event"
	messages "github.com/cucumber/cucumber-messages-go/v3"
)

// NewTestCaseRunnerOptions are the options for NewTestCaseRunner
type NewTestCaseRunnerOptions struct {
	BaseDirectory               string
	ID                          string
	Pickle                      *messages.Pickle
	SendCommand                 func(*messages.Envelope)
	SendCommandAndAwaitResponse func(*messages.Envelope) *messages.Envelope
	SupportCodeLibrary          *SupportCodeLibrary
	IsSkipped                   bool
}

// TestCaseRunner runs a test case
type TestCaseRunner struct {
	afterTestCaseHookDefinitions  []*dto.TestCaseHookDefinition
	baseDirectory                 string
	beforeTestCaseHookDefinitions []*dto.TestCaseHookDefinition
	id                            string
	isSkipped                     bool
	pickle                        *messages.Pickle
	sendCommand                   func(*messages.Envelope)
	sendCommandAndAwaitResponse   func(*messages.Envelope) *messages.Envelope
	stepIndexToStepDefinitions    [][]*dto.StepDefinition
	stepIndexToPatternMatches     [][]*messages.PatternMatch
	supportCodeLibrary            *SupportCodeLibrary

	result *messages.TestResult
}

// NewTestCaseRunner returns a TestCaseRunner
func NewTestCaseRunner(opts *NewTestCaseRunnerOptions) (*TestCaseRunner, error) {
	stepIndexToStepDefinitions := make([][]*dto.StepDefinition, len(opts.Pickle.Steps))
	stepIndexToPatternMatches := make([][]*messages.PatternMatch, len(opts.Pickle.Steps))
	for i, step := range opts.Pickle.Steps {
		var err error
		stepDefinitions, patternMatches, err := opts.SupportCodeLibrary.GetMatchingStepDefinitions(step.Text)
		if err != nil {
			return nil, err
		}
		stepIndexToStepDefinitions[i] = stepDefinitions
		stepIndexToPatternMatches[i] = patternMatches
	}
	tagNames := make([]string, len(opts.Pickle.Tags))
	for i, tag := range opts.Pickle.Tags {
		tagNames[i] = tag.Name
	}
	initialStatus := messages.TestResult_PASSED
	if opts.IsSkipped {
		initialStatus = messages.TestResult_SKIPPED
	}
	return &TestCaseRunner{
		afterTestCaseHookDefinitions:  opts.SupportCodeLibrary.GetMatchingAfterTestCaseHookDefinitions(tagNames),
		baseDirectory:                 opts.BaseDirectory,
		beforeTestCaseHookDefinitions: opts.SupportCodeLibrary.GetMatchingBeforeTestCaseHookDefinitions(tagNames),
		id:                            opts.ID,
		isSkipped:                     opts.IsSkipped,
		pickle:                        opts.Pickle,
		result: &messages.TestResult{
			DurationNanoseconds: 0,
			Status:              initialStatus,
		},
		sendCommand:                 opts.SendCommand,
		sendCommandAndAwaitResponse: opts.SendCommandAndAwaitResponse,
		stepIndexToStepDefinitions:  stepIndexToStepDefinitions,
		stepIndexToPatternMatches:   stepIndexToPatternMatches,
		supportCodeLibrary:          opts.SupportCodeLibrary,
	}, nil
}

// Run runs a test case
func (t *TestCaseRunner) Run() *messages.TestResult {
	t.sendTestCasePreparedEvent()
	t.sendTestCaseStartedEvent()
	if !t.isSkipped {
		t.sendCommandAndAwaitResponse(&messages.Envelope{
			Message: &messages.Envelope_CommandInitializeTestCase{
				CommandInitializeTestCase: &messages.CommandInitializeTestCase{
					Pickle: t.pickle,
				},
			},
		})
	}
	for index, runHookOrStepFunc := range t.getRunHookAndStepFuncs() {
		t.sendTestStepStartedEvent(index)
		hookOrStepResult := runHookOrStepFunc()
		t.sendTestStepFinishedEvent(index, hookOrStepResult)
		t.updateResult(hookOrStepResult)
	}
	t.sendTestCaseFinishedEvent()
	return t.result
}

func (t *TestCaseRunner) updateResult(hookOrStepResult *messages.TestResult) {
	t.result.DurationNanoseconds += hookOrStepResult.DurationNanoseconds
	if t.shouldUpdateResultStatus(hookOrStepResult) {
		t.result.Status = hookOrStepResult.Status
	}
	if hookOrStepResult.Message != "" && t.result.Message == "" {
		t.result.Message = hookOrStepResult.Message
	}
}

func (t *TestCaseRunner) shouldUpdateResultStatus(hookOrStepResult *messages.TestResult) bool {
	switch hookOrStepResult.Status {
	case messages.TestResult_FAILED, messages.TestResult_AMBIGUOUS:
		return t.result.Status != messages.TestResult_FAILED && t.result.Status != messages.TestResult_AMBIGUOUS
	default:
		return t.result.Status == messages.TestResult_PASSED || t.result.Status == messages.TestResult_SKIPPED
	}
}

func (t *TestCaseRunner) sendTestStepStartedEvent(index int) {
	t.sendCommand(&messages.Envelope{
		Message: &messages.Envelope_TestStepStarted{
			TestStepStarted: &messages.TestStepStarted{
				PickleId: t.pickle.Id,
				Index:    uint32(index),
			},
		},
	})
}

func (t *TestCaseRunner) sendTestStepFinishedEvent(index int, result *messages.TestResult) {
	t.sendCommand(&messages.Envelope{
		Message: &messages.Envelope_TestStepFinished{
			TestStepFinished: &messages.TestStepFinished{
				PickleId:   t.pickle.Id,
				Index:      uint32(index),
				TestResult: result,
			},
		},
	})
}

func (t *TestCaseRunner) sendTestCaseFinishedEvent() {
	t.sendCommand(&messages.Envelope{
		Message: &messages.Envelope_TestCaseFinished{
			TestCaseFinished: &messages.TestCaseFinished{
				PickleId:   t.pickle.Id,
				TestResult: t.result,
			},
		},
	})
}

func (t *TestCaseRunner) sendTestCasePreparedEvent() {
	t.sendCommand(&messages.Envelope{
		Message: &messages.Envelope_TestCasePrepared{
			TestCasePrepared: event.NewTestCasePrepared(event.NewTestCasePreparedOptions{
				AfterTestCaseHookDefinitions:  t.afterTestCaseHookDefinitions,
				BeforeTestCaseHookDefinitions: t.beforeTestCaseHookDefinitions,
				Pickle:                        t.pickle,
				StepIndexToStepDefinitions:    t.stepIndexToStepDefinitions,
			}),
		},
	})
}

func (t *TestCaseRunner) sendTestCaseStartedEvent() {
	t.sendCommand(&messages.Envelope{
		Message: &messages.Envelope_TestCaseStarted{
			TestCaseStarted: &messages.TestCaseStarted{
				PickleId: t.pickle.Id,
			},
		},
	})
}

func (t *TestCaseRunner) getRunHookAndStepFuncs() []func() *messages.TestResult {
	var result []func() *messages.TestResult
	for _, beforeTestCaseHook := range t.beforeTestCaseHookDefinitions {
		result = append(result, t.runHookFunc(beforeTestCaseHook, true))
	}
	for index, step := range t.pickle.Steps {
		result = append(result, t.runStepFunc(index, step))
	}
	for _, afterTestCaseHook := range t.afterTestCaseHookDefinitions {
		result = append(result, t.runHookFunc(afterTestCaseHook, false))
	}
	return result
}

func (t *TestCaseRunner) runHookFunc(hook *dto.TestCaseHookDefinition, isBeforeHook bool) func() *messages.TestResult {
	return func() *messages.TestResult {
		if t.isSkipped || (isBeforeHook && t.result.Status != messages.TestResult_PASSED) {
			return &messages.TestResult{Status: messages.TestResult_SKIPPED}
		}
		command := &messages.Envelope{
			Message: &messages.Envelope_CommandRunBeforeTestCaseHook{
				CommandRunBeforeTestCaseHook: &messages.CommandRunBeforeTestCaseHook{
					PickleId:                 t.pickle.Id,
					TestCaseHookDefinitionId: hook.Config.Id,
				},
			},
		}
		if !isBeforeHook {
			command = &messages.Envelope{
				Message: &messages.Envelope_CommandRunAfterTestCaseHook{
					CommandRunAfterTestCaseHook: &messages.CommandRunAfterTestCaseHook{
						PickleId:                 t.pickle.Id,
						TestCaseHookDefinitionId: hook.Config.Id,
					},
				},
			}
		}
		response := t.sendCommandAndAwaitResponse(command)
		switch x := response.Message.(type) {
		case *messages.Envelope_CommandActionComplete:
			switch y := x.CommandActionComplete.Result.(type) {
			case *messages.CommandActionComplete_TestResult:
				return y.TestResult
			}
		}
		panic(fmt.Sprintf("Received unexpected response (%v) to command (%v)", response, command))
	}
}

func (t *TestCaseRunner) runStepFunc(stepIndex int, step *messages.Pickle_PickleStep) func() *messages.TestResult {
	return func() *messages.TestResult {
		if len(t.stepIndexToStepDefinitions[stepIndex]) == 0 {
			return t.getSnippetTestResult(step)
		}
		if len(t.stepIndexToStepDefinitions[stepIndex]) > 1 {
			message, err := getAmbiguousStepDefinitionsMessage(t.stepIndexToStepDefinitions[stepIndex], t.baseDirectory)
			if err != nil {
				t.sendCommand(&messages.Envelope{
					Message: &messages.Envelope_CommandError{
						CommandError: err.Error(),
					},
				})
			}
			return &messages.TestResult{
				Status:  messages.TestResult_AMBIGUOUS,
				Message: message,
			}
		}
		if t.result.Status != messages.TestResult_PASSED {
			return &messages.TestResult{Status: messages.TestResult_SKIPPED}
		}
		return t.getRunStepTestResult(stepIndex, step)
	}
}

func (t *TestCaseRunner) getRunStepTestResult(stepIndex int, step *messages.Pickle_PickleStep) *messages.TestResult {
	command := t.getRunStepCommand(stepIndex, step)
	response := t.sendCommandAndAwaitResponse(command)
	switch x := response.Message.(type) {
	case *messages.Envelope_CommandActionComplete:
		switch y := x.CommandActionComplete.Result.(type) {
		case *messages.CommandActionComplete_TestResult:
			return y.TestResult
		}
	}
	panic(fmt.Sprintf("Received unexpected response (%v) to generate snippe command (%v)", response, command))
}

func (t *TestCaseRunner) getRunStepCommand(stepIndex int, step *messages.Pickle_PickleStep) *messages.Envelope {
	commandRunTestStep := &messages.CommandRunTestStep{
		PickleId:         t.pickle.Id,
		StepDefinitionId: t.stepIndexToStepDefinitions[stepIndex][0].Config.Id,
		PatternMatches:   t.stepIndexToPatternMatches[stepIndex],
	}
	return &messages.Envelope{
		Message: &messages.Envelope_CommandRunTestStep{
			CommandRunTestStep: commandRunTestStep,
		},
	}
}

func (t *TestCaseRunner) getSnippetTestResult(step *messages.Pickle_PickleStep) *messages.TestResult {
	command := t.getGenerateSnippetCommand(step)
	response := t.sendCommandAndAwaitResponse(command)
	switch x := response.Message.(type) {
	case *messages.Envelope_CommandActionComplete:
		switch y := x.CommandActionComplete.Result.(type) {
		case *messages.CommandActionComplete_Snippet:
			return &messages.TestResult{
				Status:  messages.TestResult_UNDEFINED,
				Message: y.Snippet,
			}
		}
	}
	panic(fmt.Sprintf("Received unexpected response (%v) to generate snippe command (%v)", response, command))
}

func (t *TestCaseRunner) getGenerateSnippetCommand(step *messages.Pickle_PickleStep) *messages.Envelope {
	commandGenerateSnippet := &messages.CommandGenerateSnippet{
		GeneratedExpressions: t.supportCodeLibrary.GenerateExpressions(step.Text),
	}
	if step.Argument != nil {
		switch x := step.Argument.Message.(type) {
		case *messages.PickleStepArgument_DataTable:
			commandGenerateSnippet.PickleStepArgument = &messages.PickleStepArgument{
				Message: &messages.PickleStepArgument_DataTable{
					DataTable: x.DataTable,
				},
			}
		case *messages.PickleStepArgument_DocString:
			commandGenerateSnippet.PickleStepArgument = &messages.PickleStepArgument{
				Message: &messages.PickleStepArgument_DocString{
					DocString: x.DocString,
				},
			}
		}
	}
	return &messages.Envelope{
		Message: &messages.Envelope_CommandGenerateSnippet{
			CommandGenerateSnippet: commandGenerateSnippet,
		},
	}
}
