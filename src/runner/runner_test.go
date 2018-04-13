package runner_test

import (
	"path"
	"runtime"

	"github.com/cucumber/cucumber-pickle-runner/src/dto"
	"github.com/cucumber/cucumber-pickle-runner/src/dto/event"
	"github.com/cucumber/cucumber-pickle-runner/src/runner"
	. "github.com/cucumber/cucumber-pickle-runner/test/matchers"
	gherkin "github.com/cucumber/gherkin-go"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Runner", func() {
	_, filename, _, _ := runtime.Caller(0)
	rootDir := path.Join(filename, "..", "..", "..")

	Context("with a feature with a single scenario with three steps", func() {
		featurePath := path.Join(rootDir, "test", "fixtures", "a.feature")

		Context("all steps are undefined", func() {
			var allCommandsSent []*dto.Command

			BeforeEach(func() {
				allCommandsSent = runWithConfigAndResponder(
					&dto.FeaturesConfig{
						AbsolutePaths: []string{featurePath},
						Filters:       &dto.FeaturesFilterConfig{},
						Language:      "en",
						Order:         &dto.FeaturesOrder{},
					},
					&dto.RuntimeConfig{},
					&dto.SupportCodeConfig{},
					func(commandChan chan *dto.Command, command *dto.Command) {
						switch command.Type {
						case dto.CommandTypeRunBeforeTestRunHooks, dto.CommandTypeRunAfterTestRunHooks, dto.CommandTypeInitializeTestCase:
							commandChan <- &dto.Command{
								Type:       dto.CommandTypeActionComplete,
								ResponseTo: command.ID,
							}
						case dto.CommandTypeGenerateSnippet:
							commandChan <- &dto.Command{
								Type:       dto.CommandTypeActionComplete,
								ResponseTo: command.ID,
								Snippet:    "snippet",
							}
						}
					},
				)
			})

			It("sends 21 commands", func() {
				Expect(allCommandsSent).To(HaveLen(21))
				Expect(allCommandsSent[0]).To(BeACommandWithEventAssignableToTypeOf(&gherkin.SourceEvent{}))
				Expect(allCommandsSent[1]).To(BeACommandWithEventAssignableToTypeOf(&gherkin.GherkinDocumentEvent{}))
				Expect(allCommandsSent[2]).To(BeACommandWithEventAssignableToTypeOf(&gherkin.PickleEvent{}))
				Expect(allCommandsSent[3]).To(BeACommandWithEventAssignableToTypeOf(&event.PickleAccepted{}))
				Expect(allCommandsSent[4]).To(BeACommandWithEventAssignableToTypeOf(&event.TestRunStarted{}))
				Expect(allCommandsSent[5]).To(BeACommandWithType(dto.CommandTypeRunBeforeTestRunHooks))
				Expect(allCommandsSent[6]).To(BeACommandWithEventAssignableToTypeOf(&event.TestCasePrepared{}))
				Expect(allCommandsSent[7]).To(BeACommandWithEventAssignableToTypeOf(&event.TestCaseStarted{}))
				Expect(allCommandsSent[8]).To(BeACommandWithType(dto.CommandTypeInitializeTestCase))
				Expect(allCommandsSent[9]).To(BeACommandWithEventAssignableToTypeOf(&event.TestStepStarted{}))
				Expect(allCommandsSent[10]).To(BeACommandWithType(dto.CommandTypeGenerateSnippet))
				Expect(allCommandsSent[11]).To(BeACommandWithEventAssignableToTypeOf(&event.TestStepFinished{}))
				Expect(allCommandsSent[12]).To(BeACommandWithEventAssignableToTypeOf(&event.TestStepStarted{}))
				Expect(allCommandsSent[13]).To(BeACommandWithType(dto.CommandTypeGenerateSnippet))
				Expect(allCommandsSent[14]).To(BeACommandWithEventAssignableToTypeOf(&event.TestStepFinished{}))
				Expect(allCommandsSent[15]).To(BeACommandWithEventAssignableToTypeOf(&event.TestStepStarted{}))
				Expect(allCommandsSent[16]).To(BeACommandWithType(dto.CommandTypeGenerateSnippet))
				Expect(allCommandsSent[17]).To(BeACommandWithEventAssignableToTypeOf(&event.TestStepFinished{}))
				Expect(allCommandsSent[18]).To(BeACommandWithEventAssignableToTypeOf(&event.TestCaseFinished{}))
				Expect(allCommandsSent[19]).To(BeACommandWithType(dto.CommandTypeRunAfterTestRunHooks))
				Expect(allCommandsSent[20]).To(Equal(&dto.Command{
					Type:  dto.CommandTypeEvent,
					Event: &event.TestRunFinished{Result: &dto.TestRunResult{Success: false}},
				}))
			})
		})
	})

	Context("all pickles gets rejected", func() {
		featurePath := path.Join(rootDir, "test", "fixtures", "a.feature")
		var allCommandsSent []*dto.Command

		BeforeEach(func() {
			allCommandsSent = runWithConfigAndResponder(
				&dto.FeaturesConfig{
					AbsolutePaths: []string{featurePath},
					Filters: &dto.FeaturesFilterConfig{
						TagExpression: "@tagA",
					},
					Language: "en",
					Order:    &dto.FeaturesOrder{},
				},
				&dto.RuntimeConfig{},
				&dto.SupportCodeConfig{},
				func(commandChan chan *dto.Command, command *dto.Command) {},
			)
		})

		It("does not run test run hooks", func() {
			Expect(allCommandsSent).To(HaveLen(6))
			Expect(allCommandsSent[0]).To(BeACommandWithEventAssignableToTypeOf(&gherkin.SourceEvent{}))
			Expect(allCommandsSent[1]).To(BeACommandWithEventAssignableToTypeOf(&gherkin.GherkinDocumentEvent{}))
			Expect(allCommandsSent[2]).To(BeACommandWithEventAssignableToTypeOf(&gherkin.PickleEvent{}))
			Expect(allCommandsSent[3]).To(BeACommandWithEventAssignableToTypeOf(&event.PickleRejected{}))
			Expect(allCommandsSent[4]).To(BeACommandWithEventAssignableToTypeOf(&event.TestRunStarted{}))
			Expect(allCommandsSent[5]).To(Equal(&dto.Command{
				Type:  dto.CommandTypeEvent,
				Event: &event.TestRunFinished{Result: &dto.TestRunResult{Success: true}},
			}))
		})
	})

	Context("with test case hooks", func() {
		featurePath := path.Join(rootDir, "test", "fixtures", "tags.feature")
		var allCommandsSent []*dto.Command

		BeforeEach(func() {
			allCommandsSent = runWithConfigAndResponder(
				&dto.FeaturesConfig{
					AbsolutePaths: []string{featurePath},
					Filters:       &dto.FeaturesFilterConfig{},
					Language:      "en",
					Order:         &dto.FeaturesOrder{},
				},
				&dto.RuntimeConfig{},
				&dto.SupportCodeConfig{
					BeforeTestCaseHookDefinitionConfigs: []*dto.TestCaseHookDefinitionConfig{
						{ID: "1", URI: "hooks.js", Line: 11},
						{ID: "2", TagExpression: "@tagA", URI: "hooks.js", Line: 12},
					},
					AfterTestCaseHookDefinitionConfigs: []*dto.TestCaseHookDefinitionConfig{
						{ID: "3", TagExpression: "@tagA", URI: "hooks.js", Line: 13},
						{ID: "4", URI: "hooks.js", Line: 14},
					},
				},
				func(commandChan chan *dto.Command, command *dto.Command) {
					switch command.Type {
					case dto.CommandTypeRunBeforeTestRunHooks, dto.CommandTypeRunAfterTestRunHooks, dto.CommandTypeInitializeTestCase:
						commandChan <- &dto.Command{
							Type:       dto.CommandTypeActionComplete,
							ResponseTo: command.ID,
						}
					case dto.CommandTypeRunBeforeTestCaseHook, dto.CommandTypeRunAfterTestCaseHook:
						commandChan <- &dto.Command{
							Type:       dto.CommandTypeActionComplete,
							ResponseTo: command.ID,
							HookOrStepResult: &dto.TestResult{
								Status: dto.StatusPassed,
							},
						}
					case dto.CommandTypeGenerateSnippet:
						commandChan <- &dto.Command{
							Type:       dto.CommandTypeActionComplete,
							ResponseTo: command.ID,
							Snippet:    "snippet",
						}
					}
				},
			)
		})

		It("runs test case hooks only for pickles that match the tag expression", func() {
			testCasePreparedEvents := []*event.TestCasePrepared{}
			for _, command := range allCommandsSent {
				if command.Type == dto.CommandTypeEvent {
					if testCasePreparedEvent, ok := command.Event.(*event.TestCasePrepared); ok {
						testCasePreparedEvents = append(testCasePreparedEvents, testCasePreparedEvent)
					}
				}
			}
			Expect(testCasePreparedEvents).To(HaveLen(2))
			Expect(testCasePreparedEvents[0]).To(Equal(&event.TestCasePrepared{
				SourceLocation: &dto.Location{
					Line: 2,
					URI:  featurePath,
				},
				Steps: []*event.TestCasePreparedStep{
					{ActionLocation: &dto.Location{Line: 11, URI: "hooks.js"}},
					{SourceLocation: &dto.Location{Line: 3, URI: featurePath}},
					{ActionLocation: &dto.Location{Line: 14, URI: "hooks.js"}},
				},
			}))
			Expect(testCasePreparedEvents[1]).To(Equal(&event.TestCasePrepared{
				SourceLocation: &dto.Location{
					Line: 6,
					URI:  featurePath,
				},
				Steps: []*event.TestCasePreparedStep{
					{ActionLocation: &dto.Location{Line: 11, URI: "hooks.js"}},
					{ActionLocation: &dto.Location{Line: 12, URI: "hooks.js"}},
					{SourceLocation: &dto.Location{Line: 7, URI: featurePath}},
					{ActionLocation: &dto.Location{Line: 13, URI: "hooks.js"}},
					{ActionLocation: &dto.Location{Line: 14, URI: "hooks.js"}},
				},
			}))
		})
	})
})

func runWithConfigAndResponder(featuresConfig *dto.FeaturesConfig, runtimeConfig *dto.RuntimeConfig, supportCodeConfig *dto.SupportCodeConfig, responder func(chan *dto.Command, *dto.Command)) []*dto.Command {
	allCommandsSent := []*dto.Command{}
	r := runner.NewRunner()
	incoming, outgoing := r.GetCommandChannels()
	done := make(chan bool)
	go func() {
		for command := range outgoing {
			allCommandsSent = append(allCommandsSent, command)
			responder(incoming, command)
		}
		done <- true
	}()
	incoming <- &dto.Command{
		Type:              dto.CommandTypeStart,
		FeaturesConfig:    featuresConfig,
		RuntimeConfig:     runtimeConfig,
		SupportCodeConfig: supportCodeConfig,
	}
	<-done
	return allCommandsSent
}
