package runner

import (
	dto "github.com/cucumber/cucumber-engine/src/dto"
	messages "github.com/cucumber/cucumber-messages-go/v2"
	uuid "github.com/satori/go.uuid"
)

type runNextTestCaseResult struct {
	err            error
	testCaseResult *messages.TestResult
}

type parallelTestCaseRunnerMaster struct {
	baseDirectory               string
	nextPickleIndex             int
	pickles                     []*messages.Pickle
	runtimeConfig               *messages.RuntimeConfig
	sendCommand                 func(*messages.Wrapper)
	sendCommandAndAwaitResponse func(*messages.Wrapper) *messages.Wrapper
	supportCodeLibrary          *SupportCodeLibrary
}

func newParallelTestCaseRunnerMaster(opts *runTestCasesOptions) *parallelTestCaseRunnerMaster {
	return &parallelTestCaseRunnerMaster{
		baseDirectory:               opts.baseDirectory,
		nextPickleIndex:             0,
		pickles:                     opts.pickles,
		runtimeConfig:               opts.runtimeConfig,
		sendCommand:                 opts.sendCommand,
		sendCommandAndAwaitResponse: opts.sendCommandAndAwaitResponse,
		supportCodeLibrary:          opts.supportCodeLibrary,
	}
}

func (p *parallelTestCaseRunnerMaster) run() (bool, error) {
	testRunResult := dto.NewTestRunResult()
	isSkipped := p.runtimeConfig.IsDryRun
	numRunning := 0
	toStart := int(p.runtimeConfig.MaxParallel)
	if toStart == 0 || toStart > len(p.pickles) {
		toStart = len(p.pickles)
	}
	onFinish := make(chan *runNextTestCaseResult, toStart)
	for i := 1; i <= toStart; i++ {
		p.runNextTestCase(isSkipped, onFinish)
		numRunning++
	}
	for numRunning > 0 {
		result := <-onFinish
		if result.err != nil {
			return false, result.err
		}
		testRunResult.Update(result.testCaseResult, p.runtimeConfig.IsStrict)
		if !isSkipped && !testRunResult.Success && p.runtimeConfig.IsFailFast {
			isSkipped = true
		}
		if p.nextPickleIndex == len(p.pickles) {
			numRunning--
		} else {
			p.runNextTestCase(isSkipped, onFinish)
		}
	}
	return testRunResult.Success, nil
}

func (p *parallelTestCaseRunnerMaster) runNextTestCase(isSkipped bool, onFinish chan *runNextTestCaseResult) {
	pickle := p.pickles[p.nextPickleIndex]
	p.nextPickleIndex++
	go func() {
		testCaseRunner, err := NewTestCaseRunner(&NewTestCaseRunnerOptions{
			BaseDirectory:               p.baseDirectory,
			ID:                          uuid.NewV4().String(),
			IsSkipped:                   isSkipped,
			Pickle:                      pickle,
			SendCommand:                 p.sendCommand,
			SendCommandAndAwaitResponse: p.sendCommandAndAwaitResponse,
			SupportCodeLibrary:          p.supportCodeLibrary,
		})
		if err != nil {
			onFinish <- &runNextTestCaseResult{err: err}
		}
		onFinish <- &runNextTestCaseResult{testCaseResult: testCaseRunner.Run()}
	}()
}
