package formatter

import (
	"fmt"

	"github.com/cucumber/cucumber-pickle-runner/src/dto"
	"github.com/cucumber/cucumber-pickle-runner/src/dto/event"
	gherkin "github.com/cucumber/gherkin-go"
)

type TestStepAttachment struct {
	Data string
}

type CompletedTestStep struct {
	Attachments []dto.A
}

type CompletedTestCase struct {

}

type TestCaseData struct {
	GherkinDocument  *gherkin.GherkinDocument
	Pickle           *gherkin.Pickle
	SourceLocation *dto.Location
	Steps []*

	TestCasePrepared *event.TestCasePrepared
	TestCaseResult   *dto.TestResult
}

type TestStepData struct {
	GherkinKeyword string
	PickleStep     *gherkin.PickleStep
	TestStep       *event.TestCasePreparedStep
}

type EventDataCollector struct {
	gherkinDocumentMap  map[string]*gherkin.GherkinDocument // key is uri
	pickleMap           map[string]*gherkin.Pickle          // key is uri:line
	testCasePreparedMap map[string]*event.TestCasePrepared  // key is uri:line
	testCaseResultMap   map[string]*dto.TestResult          // key is uri:line
}

func NewEventDataCollector(eventChannel chan gherkin.CucumberEvent) *EventDataCollector {
	e := &EventDataCollector{
		gherkinDocumentMap:  map[string]*gherkin.GherkinDocument{},
		pickleMap:           map[string]*gherkin.Pickle{},
		testCasePreparedMap: map[string]*event.TestCasePrepared{},
		testCaseResultMap:   map[string]*dto.TestResult{},
	}
	go func() {
		for ev := range eventChannel {
			switch t := ev.(type) {
			case *gherkin.GherkinDocumentEvent:
				e.storeGherkinDocument(t)
			case *gherkin.PickleEvent:
			}
		}
	}()
	return e
}

func (e *EventDataCollector) GetTestCaseData(sourceLocation *dto.Location) *TestCaseData {
	key := getTestCaseKey(sourceLocation)
	return &TestCaseData{
		GherkinDocument:  e.gherkinDocumentMap[sourceLocation.URI],
		Pickle:           e.pickleMap[key],
		TestCasePrepared: e.testCasePreparedMap[key],
		TestCaseResult:   e.testCaseResultMap[key],
	}
}

func (e *EventDataCollector) GetTestStepData(testCaseSourceLocation *dto.Location, stepIndex int) *TestStepData {
	testCaseData := e.GetTestCaseData(testCaseSourceLocation)
	result := &TestStepData{
		TestStep: testCaseData.TestCasePrepared.Steps[stepIndex],
	}
	if result.TestStep.SourceLocation != nil {
		line := result.TestStep.SourceLocation.Line
		result.GherkinKeyword = getStepLineToKeywordMap(testCaseData.GherkinDocument)[line]
		result.PickleStep = getStepLineToPickledStepMap(testCaseData.Pickle)[line]
	}
	return result
}

func (e *EventDataCollector) storeGherkinDocument(ev *gherkin.GherkinDocumentEvent) {
	e.gherkinDocumentMap[ev.URI] = ev.Document
}

func (e *EventDataCollector) storePickle(ev *gherkin.PickleEvent) {
	key := getTestCaseKey(dto.NewLocationForPickle(ev.Pickle, ev.URI))
	e.pickleMap[key] = ev.Pickle
}

func (e *EventDataCollector) storeTestCase(ev *event.TestCasePrepared) {
	key := getTestCaseKey(ev.SourceLocation)
	e.testCasePreparedMap[key] = ev
}

func (e *EventDataCollector) storeTestCaseResult(ev *event.TestCaseFinished) {
	key := getTestCaseKey(ev.SourceLocation)
	e.testCaseResultMap[key] = ev.Result
}

func getTestCaseKey(sourceLocation *dto.Location) string {
	return fmt.Sprintf("%s:%d", sourceLocation.URI, sourceLocation.Line)
}

func getStepLineToKeywordMap(g *gherkin.GherkinDocument) map[int]string {
	result := map[int]string{}
	for _, child := range g.Feature.Children {
		var steps []*gherkin.Step
		switch t := child.(type) {
		case *gherkin.Background:
			steps = t.Steps
		case *gherkin.Scenario:
			steps = t.Steps
		case *gherkin.ScenarioOutline:
			steps = t.Steps
		}
		for _, step := range steps {
			result[step.Location.Line] = step.Keyword
		}
	}
	return result
}

func getStepLineToPickledStepMap(p *gherkin.Pickle) map[int]*gherkin.PickleStep {
	result := map[int]*gherkin.PickleStep{}
	for _, step := range p.Steps {
		line := step.Locations[len(step.Locations)-1].Line
		result[line] = step
	}
	return result
}
