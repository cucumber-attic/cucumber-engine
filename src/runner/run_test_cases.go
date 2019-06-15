package runner

import (
	"github.com/cucumber/cucumber-engine/src/dto"

	messages "github.com/cucumber/cucumber-messages-go/v3"
)

type runTestCasesOptions struct {
	baseDirectory               string
	pickles                     []*messages.Pickle
	runtimeConfig               *messages.RuntimeConfig
	sendCommand                 func(*messages.Envelope)
	sendCommandAndAwaitResponse func(*messages.Envelope) *messages.Envelope
	supportCodeLibrary          *SupportCodeLibrary
}

// RunTestCasesInParallel runs the given tests cases in parallel
func RunTestCasesInParallel(opts *runTestCasesOptions) (bool, error) {
	master := newParallelTestCaseRunnerMaster(opts)
	return master.run()
}

// RunTestCasesSequentially runs the given tests cases sequentially
func RunTestCasesSequentially(opts *runTestCasesOptions) (bool, error) {
	testRunResult := dto.NewTestRunResult()
	isSkipped := opts.runtimeConfig.IsDryRun
	for _, pickle := range opts.pickles {
		testCaseRunner, err := NewTestCaseRunner(&NewTestCaseRunnerOptions{
			BaseDirectory:               opts.baseDirectory,
			IsSkipped:                   isSkipped,
			Pickle:                      pickle,
			SendCommand:                 opts.sendCommand,
			SendCommandAndAwaitResponse: opts.sendCommandAndAwaitResponse,
			SupportCodeLibrary:          opts.supportCodeLibrary,
		})
		if err != nil {
			return false, err
		}
		testCaseResult := testCaseRunner.Run()
		testRunResult.Update(testCaseResult, opts.runtimeConfig.IsStrict)
		if !isSkipped && !testRunResult.Success && opts.runtimeConfig.IsFailFast {
			isSkipped = true
		}
	}
	return testRunResult.Success, nil
}
