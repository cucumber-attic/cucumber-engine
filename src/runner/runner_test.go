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

		It("all steps are undefined", func() {
			allCommandsSent := runWithConfigAndResponder(
				&dto.FeaturesConfig{
					AbsolutePaths: []string{featurePath},
					Filters:       &dto.FeaturesFilterConfig{},
				},
				&dto.RuntimeConfig{
					BeforeTestCaseHookDefinitions: []*dto.TestCaseHookDefinition{},
					AfterTestCaseHookDefinitions:  []*dto.TestCaseHookDefinition{},
					StepDefinitions:               []*dto.StepDefinition{},
				},
				func(commandChan chan *dto.Command, command *dto.Command) {
					switch command.Type {
					case dto.CommandTypeRunBeforeTestRunHooks, dto.CommandTypeRunAfterTestRunHooks, dto.CommandTypeInitializeTestCase:
						commandChan <- &dto.Command{
							Type:       dto.CommandTypeActionComplete,
							ResponseTo: command.ID,
						}
					}
				},
			)
			Expect(allCommandsSent).To(HaveLen(18))
			Expect(allCommandsSent[0]).To(BeACommandWithEventAssignableToTypeOf(&gherkin.SourceEvent{}))
			Expect(allCommandsSent[1]).To(BeACommandWithEventAssignableToTypeOf(&gherkin.GherkinDocumentEvent{}))
			Expect(allCommandsSent[2]).To(BeACommandWithEventAssignableToTypeOf(&gherkin.PickleEvent{}))
			Expect(allCommandsSent[3]).To(BeACommandWithEventAssignableToTypeOf(&event.PickleAccepted{}))
			Expect(allCommandsSent[4]).To(BeACommandWithEventAssignableToTypeOf(&event.TestRunStarted{}))
			Expect(allCommandsSent[5]).To(BeACommandWithType(dto.CommandTypeRunBeforeTestRunHooks))
			Expect(allCommandsSent[6]).To(BeACommandWithEventAssignableToTypeOf(&event.TestCasePrepared{}))
			Expect(allCommandsSent[7]).To(Equal(&dto.Command{
				Type: dto.CommandTypeEvent,
				Event: &event.TestCaseStarted{
					SourceLocation: event.Location{
						URI:  featurePath,
						Line: 2,
					},
				},
			}))
			Expect(allCommandsSent[8]).To(BeACommandWithType(dto.CommandTypeInitializeTestCase))
			Expect(allCommandsSent[9]).To(Equal(&dto.Command{
				Type: dto.CommandTypeEvent,
				Event: &event.TestStepStarted{
					Index: 0,
					TestCase: event.TestCase{
						SourceLocation: event.Location{
							URI:  featurePath,
							Line: 2,
						},
					},
				},
			}))
			Expect(allCommandsSent[10]).To(Equal(&dto.Command{
				Type: dto.CommandTypeEvent,
				Event: &event.TestStepFinished{
					Index:  0,
					Result: &dto.TestResult{Status: dto.StatusUndefined},
					TestCase: event.TestCase{
						SourceLocation: event.Location{
							URI:  featurePath,
							Line: 2,
						},
					},
				},
			}))
			Expect(allCommandsSent[11]).To(Equal(&dto.Command{
				Type: dto.CommandTypeEvent,
				Event: &event.TestStepStarted{
					Index: 1,
					TestCase: event.TestCase{
						SourceLocation: event.Location{
							URI:  featurePath,
							Line: 2,
						},
					},
				},
			}))
			Expect(allCommandsSent[12]).To(Equal(&dto.Command{
				Type: dto.CommandTypeEvent,
				Event: &event.TestStepFinished{
					Index:  1,
					Result: &dto.TestResult{Status: dto.StatusUndefined},
					TestCase: event.TestCase{
						SourceLocation: event.Location{
							URI:  featurePath,
							Line: 2,
						},
					},
				},
			}))
			Expect(allCommandsSent[13]).To(Equal(&dto.Command{
				Type: dto.CommandTypeEvent,
				Event: &event.TestStepStarted{
					Index: 2,
					TestCase: event.TestCase{
						SourceLocation: event.Location{
							URI:  featurePath,
							Line: 2,
						},
					},
				},
			}))
			Expect(allCommandsSent[14]).To(Equal(&dto.Command{
				Type: dto.CommandTypeEvent,
				Event: &event.TestStepFinished{
					Index:  2,
					Result: &dto.TestResult{Status: dto.StatusUndefined},
					TestCase: event.TestCase{
						SourceLocation: event.Location{
							URI:  featurePath,
							Line: 2,
						},
					},
				},
			}))
			Expect(allCommandsSent[15]).To(Equal(&dto.Command{
				Type: dto.CommandTypeEvent,
				Event: &event.TestCaseFinished{
					Result: &dto.TestResult{Status: dto.StatusUndefined},
					SourceLocation: event.Location{
						URI:  featurePath,
						Line: 2,
					},
				},
			}))
			Expect(allCommandsSent[16]).To(BeACommandWithType(dto.CommandTypeRunAfterTestRunHooks))
			Expect(allCommandsSent[17]).To(Equal(&dto.Command{
				Type:  dto.CommandTypeEvent,
				Event: &event.TestRunFinished{Success: false},
			}))
		})
	})

	Context("all pickles gets rejected", func() {
		featurePath := path.Join(rootDir, "test", "fixtures", "a.feature")

		It("does not run test run hooks", func() {
			allCommandsSent := runWithConfigAndResponder(
				&dto.FeaturesConfig{
					AbsolutePaths: []string{featurePath},
					Filters: &dto.FeaturesFilterConfig{
						TagExpression: "@tagA",
					},
				},
				&dto.RuntimeConfig{
					BeforeTestCaseHookDefinitions: []*dto.TestCaseHookDefinition{},
					AfterTestCaseHookDefinitions:  []*dto.TestCaseHookDefinition{},
					StepDefinitions:               []*dto.StepDefinition{},
				},
				func(commandChan chan *dto.Command, command *dto.Command) {},
			)
			Expect(allCommandsSent).To(HaveLen(6))
			Expect(allCommandsSent[0]).To(BeACommandWithEventAssignableToTypeOf(&gherkin.SourceEvent{}))
			Expect(allCommandsSent[1]).To(BeACommandWithEventAssignableToTypeOf(&gherkin.GherkinDocumentEvent{}))
			Expect(allCommandsSent[2]).To(BeACommandWithEventAssignableToTypeOf(&gherkin.PickleEvent{}))
			Expect(allCommandsSent[3]).To(BeACommandWithEventAssignableToTypeOf(&event.PickleRejected{}))
			Expect(allCommandsSent[4]).To(BeACommandWithEventAssignableToTypeOf(&event.TestRunStarted{}))
			Expect(allCommandsSent[5]).To(Equal(&dto.Command{
				Type:  dto.CommandTypeEvent,
				Event: &event.TestRunFinished{Success: true},
			}))
		})
	})
})

func runWithConfigAndResponder(featuresConfig *dto.FeaturesConfig, runtimeConfig *dto.RuntimeConfig, responder func(chan *dto.Command, *dto.Command)) []*dto.Command {
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
		Type:           dto.CommandTypeStart,
		FeaturesConfig: featuresConfig,
		RuntimeConfig:  runtimeConfig,
	}
	<-done
	return allCommandsSent
}
