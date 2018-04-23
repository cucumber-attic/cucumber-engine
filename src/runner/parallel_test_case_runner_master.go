package runner

import (
	"github.com/cucumber/cucumber-pickle-runner/src/dto"
	gherkin "github.com/cucumber/gherkin-go"
	uuid "github.com/satori/go.uuid"
)

type runNextTestCaseResult struct {
	err            error
	testCaseResult *dto.TestResult
}

type parallelTestCaseRunnerMaster struct {
	baseDirectory               string
	nextPickleEventIndex        int
	maxParallel                 int
	pickleEvents                []*gherkin.PickleEvent
	runtimeConfig               *dto.RuntimeConfig
	sendCommand                 func(*dto.Command)
	sendCommandAndAwaitResponse func(*dto.Command) *dto.Command
	supportCodeLibrary          *SupportCodeLibrary
	testRunResult               *dto.TestRunResult
}

func newParallelTestCaseRunnerMaster(opts *runTestCasesOptions) *parallelTestCaseRunnerMaster {
	return &parallelTestCaseRunnerMaster{
		baseDirectory:               opts.baseDirectory,
		nextPickleEventIndex:        0,
		pickleEvents:                opts.pickleEvents,
		runtimeConfig:               opts.runtimeConfig,
		sendCommand:                 opts.sendCommand,
		sendCommandAndAwaitResponse: opts.sendCommandAndAwaitResponse,
		supportCodeLibrary:          opts.supportCodeLibrary,
	}
}

func (p *parallelTestCaseRunnerMaster) run() (*dto.TestRunResult, error) {
	testRunResult := dto.NewTestRunResult()
	isSkipped := p.runtimeConfig.IsDryRun
	numRunning := 0
	toStart := p.runtimeConfig.MaxParallel
	if toStart == -1 || toStart > len(p.pickleEvents) {
		toStart = len(p.pickleEvents)
	}
	onFinish := make(chan *runNextTestCaseResult, toStart)
	for i := 1; i <= toStart; i++ {
		p.runNextTestCase(isSkipped, onFinish)
		numRunning++
	}
	for numRunning > 0 {
		result := <-onFinish
		if result.err != nil {
			return nil, result.err
		}
		testRunResult.Update(result.testCaseResult, p.runtimeConfig.IsStrict)
		if !isSkipped && !testRunResult.Success && p.runtimeConfig.IsFailFast {
			isSkipped = true
		}
		if p.nextPickleEventIndex == len(p.pickleEvents) {
			numRunning--
		} else {
			p.runNextTestCase(isSkipped, onFinish)
		}
	}
	return testRunResult, nil
}

func (p *parallelTestCaseRunnerMaster) runNextTestCase(isSkipped bool, onFinish chan *runNextTestCaseResult) {
	pickleEvent := p.pickleEvents[p.nextPickleEventIndex]
	p.nextPickleEventIndex++
	go func() {
		testCaseRunner, err := NewTestCaseRunner(&NewTestCaseRunnerOptions{
			BaseDirectory:               p.baseDirectory,
			ID:                          uuid.NewV4().String(),
			IsSkipped:                   isSkipped,
			Pickle:                      pickleEvent.Pickle,
			SendCommand:                 p.sendCommand,
			SendCommandAndAwaitResponse: p.sendCommandAndAwaitResponse,
			SupportCodeLibrary:          p.supportCodeLibrary,
			URI:                         pickleEvent.URI,
		})
		if err != nil {
			onFinish <- &runNextTestCaseResult{err: err}
		}
		onFinish <- &runNextTestCaseResult{testCaseResult: testCaseRunner.Run()}
	}()
}
