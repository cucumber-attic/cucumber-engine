package runner

import (
	"fmt"
	"strings"

	"github.com/cucumber/cucumber-pickle-runner/src/dto"
	"github.com/cucumber/cucumber-pickle-runner/src/dto/event"
	gherkin "github.com/cucumber/gherkin-go"
)

// NewTestCaseRunnerOptions are the options for NewTestCaseRunner
type NewTestCaseRunnerOptions struct {
	BaseDirectory               string
	ID                          string
	Pickle                      *gherkin.Pickle
	URI                         string
	SendCommand                 func(*dto.Command)
	SendCommandAndAwaitResponse func(*dto.Command) *dto.Command
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
	pickle                        *gherkin.Pickle
	sendCommand                   func(*dto.Command)
	sendCommandAndAwaitResponse   func(*dto.Command) *dto.Command
	stepIndexToStepDefinitions    [][]*dto.StepDefinition
	stepIndexToPatternMatches     [][]*dto.PatternMatch
	supportCodeLibrary            *SupportCodeLibrary
	uri                           string

	result *dto.TestResult
}

// NewTestCaseRunner returns a TestCaseRunner
func NewTestCaseRunner(opts *NewTestCaseRunnerOptions) (*TestCaseRunner, error) {
	stepIndexToStepDefinitions := make([][]*dto.StepDefinition, len(opts.Pickle.Steps))
	stepIndexToPatternMatches := make([][]*dto.PatternMatch, len(opts.Pickle.Steps))
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
	initialStatus := dto.StatusPassed
	if opts.IsSkipped {
		initialStatus = dto.StatusSkipped
	}
	return &TestCaseRunner{
		afterTestCaseHookDefinitions:  opts.SupportCodeLibrary.GetMatchingAfterTestCaseHookDefinitions(tagNames),
		baseDirectory:                 opts.BaseDirectory,
		beforeTestCaseHookDefinitions: opts.SupportCodeLibrary.GetMatchingBeforeTestCaseHookDefinitions(tagNames),
		id:        opts.ID,
		isSkipped: opts.IsSkipped,
		pickle:    opts.Pickle,
		result: &dto.TestResult{
			Duration: 0,
			Status:   initialStatus,
		},
		sendCommand:                 opts.SendCommand,
		sendCommandAndAwaitResponse: opts.SendCommandAndAwaitResponse,
		stepIndexToStepDefinitions:  stepIndexToStepDefinitions,
		stepIndexToPatternMatches:   stepIndexToPatternMatches,
		supportCodeLibrary:          opts.SupportCodeLibrary,
		uri:                         opts.URI,
	}, nil
}

// Run runs a test case
func (t *TestCaseRunner) Run() *dto.TestResult {
	t.sendTestCasePreparedEvent()
	t.sendTestCaseStartedEvent()
	if !t.isSkipped {
		t.sendCommandAndAwaitResponse(&dto.Command{
			Type:       dto.CommandTypeInitializeTestCase,
			TestCaseID: t.id,
			TestCase: &dto.TestCase{
				SourceLocation: dto.NewLocationForPickle(t.pickle, t.uri),
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

func (t *TestCaseRunner) updateResult(hookOrStepResult *dto.TestResult) {
	t.result.Duration += hookOrStepResult.Duration
	if t.shouldUpdateResultStatus(hookOrStepResult) {
		t.result.Status = hookOrStepResult.Status
	}
	if hookOrStepResult.Message != "" && t.result.Message == "" {
		t.result.Message = hookOrStepResult.Message
	}
}

func (t *TestCaseRunner) shouldUpdateResultStatus(hookOrStepResult *dto.TestResult) bool {
	switch hookOrStepResult.Status {
	case dto.StatusFailed, dto.StatusAmbiguous:
		return t.result.Status != dto.StatusFailed || t.result.Status != dto.StatusAmbiguous
	default:
		return t.result.Status == dto.StatusPassed || t.result.Status == dto.StatusSkipped
	}
}

func (t *TestCaseRunner) sendTestStepStartedEvent(index int) {
	t.sendCommand(&dto.Command{
		Type: "event",
		Event: event.NewTestStepStarted(event.NewTestStepStartedOptions{
			Index:  index,
			Pickle: t.pickle,
			URI:    t.uri,
		}),
	})
}

func (t *TestCaseRunner) sendTestStepFinishedEvent(index int, result *dto.TestResult) {
	t.sendCommand(&dto.Command{
		Type: "event",
		Event: event.NewTestStepFinished(event.NewTestStepFinishedOptions{
			Index:  index,
			Pickle: t.pickle,
			Result: result,
			URI:    t.uri,
		}),
	})
}

func (t *TestCaseRunner) sendTestCaseFinishedEvent() {
	t.sendCommand(&dto.Command{
		Type: "event",
		Event: event.NewTestCaseFinished(event.NewTestCaseFinishedOptions{
			Pickle: t.pickle,
			Result: t.result,
			URI:    t.uri,
		}),
	})
}

func (t *TestCaseRunner) sendTestCasePreparedEvent() {
	t.sendCommand(&dto.Command{
		Type: "event",
		Event: event.NewTestCasePrepared(event.NewTestCasePreparedOptions{
			AfterTestCaseHookDefinitions:  t.afterTestCaseHookDefinitions,
			BeforeTestCaseHookDefinitions: t.beforeTestCaseHookDefinitions,
			Pickle: t.pickle,
			StepIndexToStepDefinitions: t.stepIndexToStepDefinitions,
			URI: t.uri,
		}),
	})
}

func (t *TestCaseRunner) sendTestCaseStartedEvent() {
	t.sendCommand(&dto.Command{
		Type: "event",
		Event: event.NewTestCaseStarted(event.NewTestCaseStartedOptions{
			Pickle: t.pickle,
			URI:    t.uri,
		}),
	})
}

func (t *TestCaseRunner) getRunHookAndStepFuncs() []func() *dto.TestResult {
	var result []func() *dto.TestResult
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

func (t *TestCaseRunner) runHookFunc(hook *dto.TestCaseHookDefinition, isBeforeHook bool) func() *dto.TestResult {
	return func() *dto.TestResult {
		if t.isSkipped || (isBeforeHook && t.result.Status != dto.StatusPassed) {
			return &dto.TestResult{Status: dto.StatusSkipped}
		}
		commandType := dto.CommandTypeRunAfterTestCaseHook
		if isBeforeHook {
			commandType = dto.CommandTypeRunBeforeTestCaseHook
		}
		// TODO also send the test case result
		response := t.sendCommandAndAwaitResponse(&dto.Command{
			Type:                     commandType,
			TestCaseID:               t.id,
			TestCaseHookDefinitionID: hook.ID,
		})
		return response.HookOrStepResult
	}
}

func (t *TestCaseRunner) runStepFunc(stepIndex int, step *gherkin.PickleStep) func() *dto.TestResult {
	return func() *dto.TestResult {
		if len(t.stepIndexToStepDefinitions[stepIndex]) == 0 {
			response := t.sendCommandAndAwaitResponse(&dto.Command{
				Type:                 dto.CommandTypeGenerateSnippet,
				GeneratedExpressions: t.supportCodeLibrary.GenerateExpressions(step.Text),
				PickleArguments:      step.Arguments,
			})
			return &dto.TestResult{
				Status:  dto.StatusUndefined,
				Message: fmt.Sprintf("Undefined. Implement with the following snippet:\n\n  %s", strings.Replace(response.Snippet, "\n", "\n  ", -1)),
			}
		}
		if len(t.stepIndexToStepDefinitions[stepIndex]) > 1 {
			message, err := getAmbiguousStepDefinitionsMessage(t.stepIndexToStepDefinitions[stepIndex], t.baseDirectory)
			if err != nil {
				t.sendCommand(&dto.Command{
					Type:  dto.CommandTypeError,
					Error: err.Error(),
				})
			}
			return &dto.TestResult{
				Status:  dto.StatusAmbiguous,
				Message: message,
			}
		}
		if t.result.Status != dto.StatusPassed {
			return &dto.TestResult{Status: dto.StatusSkipped}
		}
		response := t.sendCommandAndAwaitResponse(&dto.Command{
			Type:             dto.CommandTypeRunTestStep,
			TestCaseID:       t.id,
			StepDefinitionID: t.stepIndexToStepDefinitions[stepIndex][0].ID,
			PatternMatches:   t.stepIndexToPatternMatches[stepIndex],
			PickleArguments:  step.Arguments,
		})
		return response.HookOrStepResult
	}
}
