package runner

import (
	"github.com/cucumber/cucumber-engine/src/dto"
	"github.com/cucumber/gherkin-go"
	"github.com/satori/go.uuid"
)

type runTestCasesOptions struct {
	baseDirectory               string
	pickleEvents                []*gherkin.PickleEvent
	runtimeConfig               *dto.RuntimeConfig
	sendCommand                 func(*dto.Command)
	sendCommandAndAwaitResponse func(*dto.Command) *dto.Command
	supportCodeLibrary          *SupportCodeLibrary
}

// RunTestCasesInParallel runs the given tests cases in parallel
func RunTestCasesInParallel(opts *runTestCasesOptions) (*dto.TestRunResult, error) {
	master := newParallelTestCaseRunnerMaster(opts)
	return master.run()
}

// RunTestCasesSequentially runs the given tests cases sequentially
func RunTestCasesSequentially(opts *runTestCasesOptions) (*dto.TestRunResult, error) {
	testRunResult := dto.NewTestRunResult()
	isSkipped := opts.runtimeConfig.IsDryRun
	for _, pickleEvent := range opts.pickleEvents {
		testCaseRunner, err := NewTestCaseRunner(&NewTestCaseRunnerOptions{
			BaseDirectory:               opts.baseDirectory,
			ID:                          uuid.NewV4().String(),
			IsSkipped:                   isSkipped,
			Pickle:                      pickleEvent.Pickle,
			SendCommand:                 opts.sendCommand,
			SendCommandAndAwaitResponse: opts.sendCommandAndAwaitResponse,
			SupportCodeLibrary:          opts.supportCodeLibrary,
			URI:                         pickleEvent.URI,
		})
		if err != nil {
			return nil, err
		}
		testCaseResult := testCaseRunner.Run()
		testRunResult.Update(testCaseResult, opts.runtimeConfig.IsStrict)
		if !isSkipped && !testRunResult.Success && opts.runtimeConfig.IsFailFast {
			isSkipped = true
		}
	}
	return testRunResult, nil
}
