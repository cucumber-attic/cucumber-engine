package runner_test

import (
	"fmt"
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
			allCommandsSent := []*dto.Command{}
			r := runner.NewRunner()
			incoming, outgoing := r.GetCommandChannels()
			done := make(chan bool)
			go func() {
				for command := range outgoing {
					fmt.Printf("test %+v\n", command)
					if command.Type == dto.CommandTypeEvent {
						fmt.Printf("test %+v\n", command.Event)
					}
					allCommandsSent = append(allCommandsSent, command)
					switch command.Type {
					case dto.CommandTypeRunBeforeTestRunHooks, dto.CommandTypeRunAfterTestRunHooks, dto.CommandTypeInitializeTestCase:
						incoming <- &dto.Command{
							Type:       dto.CommandTypeActionComplete,
							ResponseTo: command.ID,
						}
					}
				}
				done <- true
			}()
			incoming <- &dto.Command{
				Type: dto.CommandTypeStart,
				FeaturesConfig: &dto.FeaturesConfig{
					AbsolutePaths: []string{featurePath},
				},
				RuntimeConfig: &dto.RuntimeConfig{
					BeforeTestCaseHookDefinitions: []*dto.TestCaseHookDefinition{},
					AfterTestCaseHookDefinitions:  []*dto.TestCaseHookDefinition{},
					StepDefinitions:               []*dto.StepDefinition{},
				},
			}
			<-done
			Expect(allCommandsSent).To(HaveLen(17))
			Expect(allCommandsSent[0]).To(BeACommandWithEventAssignableToTypeOf(&gherkin.SourceEvent{}))
			Expect(allCommandsSent[1]).To(BeACommandWithEventAssignableToTypeOf(&gherkin.GherkinDocumentEvent{}))
			Expect(allCommandsSent[2]).To(BeACommandWithEventAssignableToTypeOf(&gherkin.PickleEvent{}))
			Expect(allCommandsSent[3]).To(BeACommandWithEventAssignableToTypeOf(&event.TestRunStarted{}))
			Expect(allCommandsSent[4]).To(BeACommandWithType(dto.CommandTypeRunBeforeTestRunHooks))
			Expect(allCommandsSent[5]).To(BeACommandWithEventAssignableToTypeOf(&event.TestCasePrepared{}))
			Expect(allCommandsSent[6]).To(Equal(&dto.Command{
				Type: dto.CommandTypeEvent,
				Event: &event.TestCaseStarted{
					SourceLocation: event.Location{
						URI:  featurePath,
						Line: 2,
					},
				},
			}))
			Expect(allCommandsSent[7]).To(BeACommandWithType(dto.CommandTypeInitializeTestCase))
			Expect(allCommandsSent[8]).To(Equal(&dto.Command{
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
			Expect(allCommandsSent[9]).To(Equal(&dto.Command{
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
			Expect(allCommandsSent[10]).To(Equal(&dto.Command{
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
			Expect(allCommandsSent[11]).To(Equal(&dto.Command{
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
			Expect(allCommandsSent[12]).To(Equal(&dto.Command{
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
			Expect(allCommandsSent[13]).To(Equal(&dto.Command{
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
			Expect(allCommandsSent[14]).To(Equal(&dto.Command{
				Type: dto.CommandTypeEvent,
				Event: &event.TestCaseFinished{
					Result: &dto.TestResult{Status: dto.StatusUndefined},
					SourceLocation: event.Location{
						URI:  featurePath,
						Line: 2,
					},
				},
			}))
			Expect(allCommandsSent[15]).To(BeACommandWithType(dto.CommandTypeRunAfterTestRunHooks))
			Expect(allCommandsSent[16]).To(Equal(&dto.Command{
				Type:  dto.CommandTypeEvent,
				Event: &event.TestRunFinished{Success: false},
			}))
		})
	})
})
