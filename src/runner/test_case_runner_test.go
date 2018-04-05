package runner_test

import (
	"github.com/cucumber/cucumber-pickle-runner/src/dto"
	"github.com/cucumber/cucumber-pickle-runner/src/dto/event"
	"github.com/cucumber/cucumber-pickle-runner/src/runner"
	gherkin "github.com/cucumber/gherkin-go"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("TestCaseRunner", func() {
	Context("with a passing step", func() {
		var allCommandsSent []*dto.Command
		var result *dto.TestResult

		BeforeEach(func() {
			allCommandsSent = []*dto.Command{}
			sendCommand := func(command *dto.Command) {
				allCommandsSent = append(allCommandsSent, command)
			}
			sendCommandAndAwaitResponse := func(command *dto.Command) *dto.Command {
				allCommandsSent = append(allCommandsSent, command)
				if command.Type == dto.CommandTypeRunTestStep {
					return &dto.Command{
						Type: dto.CommandTypeActionComplete,
						HookOrStepResult: &dto.TestResult{
							Status:   dto.StatusPassed,
							Duration: 9,
						},
					}
				}
				return &dto.Command{
					Type: dto.CommandTypeActionComplete,
				}
			}
			supportCodeLibrary, err := runner.NewSupportCodeLibrary(&dto.SupportCodeConfig{
				StepDefinitionConfigs: []*dto.StepDefinitionConfig{
					{
						ID: "step1",
						Pattern: dto.Pattern{
							Source: "I have {int} cukes",
							Type:   dto.PatternTypeCucumberExpression,
						},
						Line: 3,
						URI:  "/path/to/steps",
					},
				},
			})
			Expect(err).NotTo(HaveOccurred())
			testCaseRunner, err := runner.NewTestCaseRunner(&runner.NewTestCaseRunnerOptions{
				ID: "testCase1",
				Pickle: &gherkin.Pickle{
					Locations: []gherkin.Location{{Line: 1}},
					Steps: []*gherkin.PickleStep{
						{
							Locations: []gherkin.Location{{Line: 2}},
							Text:      "I have 100 cukes",
						},
					},
				},
				SendCommand:                 sendCommand,
				SendCommandAndAwaitResponse: sendCommandAndAwaitResponse,
				SupportCodeLibrary:          supportCodeLibrary,
				URI:                         "/path/to/feature",
			})
			Expect(err).NotTo(HaveOccurred())
			result = testCaseRunner.Run()
		})

		It("returns a passing result", func() {
			Expect(result).To(Equal(&dto.TestResult{
				Duration: 9,
				Status:   dto.StatusPassed,
			}))
		})

		It("sends 7 commands", func() {
			Expect(allCommandsSent).To(HaveLen(7))
		})

		It("sends the test case prepared event command", func() {
			Expect(allCommandsSent[0]).To(Equal(&dto.Command{
				Type: dto.CommandTypeEvent,
				Event: &event.TestCasePrepared{
					SourceLocation: &event.Location{
						URI:  "/path/to/feature",
						Line: 1,
					},
					Steps: []*event.TestCasePreparedStep{
						{
							SourceLocation: &event.Location{
								URI:  "/path/to/feature",
								Line: 2,
							},
							ActionLocation: &event.Location{
								URI:  "/path/to/steps",
								Line: 3,
							},
						},
					},
				},
			}))
		})

		It("sends the test case started event command", func() {
			Expect(allCommandsSent[1]).To(Equal(&dto.Command{
				Type: dto.CommandTypeEvent,
				Event: &event.TestCaseStarted{
					SourceLocation: &event.Location{
						URI:  "/path/to/feature",
						Line: 1,
					},
				},
			}))
		})

		It("sends the initialize test case command", func() {
			Expect(allCommandsSent[2]).To(Equal(&dto.Command{
				Type:       dto.CommandTypeInitializeTestCase,
				TestCaseID: "testCase1",
			}))
		})

		It("sends the test step started event commands", func() {
			Expect(allCommandsSent[3]).To(Equal(&dto.Command{
				Type: dto.CommandTypeEvent,
				Event: &event.TestStepStarted{
					Index: 0,
					TestCase: &event.TestCase{
						SourceLocation: &event.Location{
							URI:  "/path/to/feature",
							Line: 1,
						},
					},
				},
			}))
		})

		It("sends the run test step command", func() {
			Expect(allCommandsSent[4]).To(Equal(&dto.Command{
				Type:             dto.CommandTypeRunTestStep,
				TestCaseID:       "testCase1",
				StepDefinitionID: "step1",
				PatternMatches: []*dto.PatternMatch{
					{
						Captures:          []string{"100"},
						ParameterTypeName: "int",
					},
				},
			}))
		})

		It("sends the test step finished event command", func() {
			Expect(allCommandsSent[5]).To(Equal(&dto.Command{
				Type: dto.CommandTypeEvent,
				Event: &event.TestStepFinished{
					Index:  0,
					Result: &dto.TestResult{Duration: 9, Status: dto.StatusPassed},
					TestCase: &event.TestCase{
						SourceLocation: &event.Location{
							URI:  "/path/to/feature",
							Line: 1,
						},
					},
				},
			}))
		})

		It("sends the test case finished event command", func() {
			Expect(allCommandsSent[6]).To(Equal(&dto.Command{
				Type: dto.CommandTypeEvent,
				Event: &event.TestCaseFinished{
					SourceLocation: &event.Location{
						URI:  "/path/to/feature",
						Line: 1,
					},
					Result: &dto.TestResult{
						Duration: 9,
						Status:   dto.StatusPassed,
					},
				},
			}))
		})
	})

	Context("with a failing step", func() {
		var allCommandsSent []*dto.Command
		var result *dto.TestResult

		BeforeEach(func() {
			allCommandsSent = []*dto.Command{}
			sendCommand := func(command *dto.Command) {
				allCommandsSent = append(allCommandsSent, command)
			}
			sendCommandAndAwaitResponse := func(command *dto.Command) *dto.Command {
				allCommandsSent = append(allCommandsSent, command)
				if command.Type == dto.CommandTypeRunTestStep {
					return &dto.Command{
						Type: dto.CommandTypeActionComplete,
						HookOrStepResult: &dto.TestResult{
							Status:   dto.StatusFailed,
							Duration: 8,
							Message:  "error message and stacktrace",
						},
					}
				}
				return &dto.Command{
					Type: dto.CommandTypeActionComplete,
				}
			}
			supportCodeLibrary, err := runner.NewSupportCodeLibrary(&dto.SupportCodeConfig{
				StepDefinitionConfigs: []*dto.StepDefinitionConfig{
					{
						ID: "step1",
						Pattern: dto.Pattern{
							Source: "I have {int} cukes",
							Type:   dto.PatternTypeCucumberExpression,
						},
						Line: 3,
						URI:  "/path/to/steps",
					},
				},
			})
			Expect(err).NotTo(HaveOccurred())
			testCaseRunner, err := runner.NewTestCaseRunner(&runner.NewTestCaseRunnerOptions{
				ID: "testCase1",
				Pickle: &gherkin.Pickle{
					Locations: []gherkin.Location{{Line: 1}},
					Steps: []*gherkin.PickleStep{
						{
							Locations: []gherkin.Location{{Line: 2}},
							Text:      "I have 100 cukes",
						},
					},
				},
				SendCommand:                 sendCommand,
				SendCommandAndAwaitResponse: sendCommandAndAwaitResponse,
				SupportCodeLibrary:          supportCodeLibrary,
				URI:                         "/path/to/feature",
			})
			Expect(err).NotTo(HaveOccurred())
			result = testCaseRunner.Run()
		})

		It("returns a failing result", func() {
			Expect(result).To(Equal(&dto.TestResult{
				Status:   dto.StatusFailed,
				Duration: 8,
				Message:  "error message and stacktrace",
			}))
		})

		It("sends 7 commands", func() {
			Expect(allCommandsSent).To(HaveLen(7))
		})

		It("sends the test step finished event command with status failed", func() {
			Expect(allCommandsSent[5]).To(Equal(&dto.Command{
				Type: dto.CommandTypeEvent,
				Event: &event.TestStepFinished{
					Index: 0,
					Result: &dto.TestResult{
						Status:   dto.StatusFailed,
						Duration: 8,
						Message:  "error message and stacktrace",
					},
					TestCase: &event.TestCase{
						SourceLocation: &event.Location{
							URI:  "/path/to/feature",
							Line: 1,
						},
					},
				},
			}))
		})

		It("sends the test case finished event command with status failed", func() {
			Expect(allCommandsSent[6]).To(Equal(&dto.Command{
				Type: dto.CommandTypeEvent,
				Event: &event.TestCaseFinished{
					SourceLocation: &event.Location{
						URI:  "/path/to/feature",
						Line: 1,
					},
					Result: &dto.TestResult{
						Status:   dto.StatusFailed,
						Duration: 8,
						Message:  "error message and stacktrace",
					},
				},
			}))
		})
	})

	Context("with a ambiguous step", func() {
		var allCommandsSent []*dto.Command
		var result *dto.TestResult
		expectedMessage := "Multiple step definitions match:\n" +
			"  I have {int} cukes     /path/to/steps:3  \n" +
			`  ^I have (\d+) cukes$   /path/to/steps:4  ` + "\n"

		BeforeEach(func() {
			allCommandsSent = []*dto.Command{}
			sendCommand := func(command *dto.Command) {
				allCommandsSent = append(allCommandsSent, command)
			}
			sendCommandAndAwaitResponse := func(command *dto.Command) *dto.Command {
				allCommandsSent = append(allCommandsSent, command)
				return &dto.Command{
					Type: dto.CommandTypeActionComplete,
				}
			}
			supportCodeLibrary, err := runner.NewSupportCodeLibrary(&dto.SupportCodeConfig{
				StepDefinitionConfigs: []*dto.StepDefinitionConfig{
					{
						ID: "step1",
						Pattern: dto.Pattern{
							Source: "I have {int} cukes",
							Type:   dto.PatternTypeCucumberExpression,
						},
						Line: 3,
						URI:  "/path/to/steps",
					},
					{
						ID: "step2",
						Pattern: dto.Pattern{
							Source: `^I have (\d+) cukes$`,
							Type:   dto.PatternTypeRegularExpression,
						},
						Line: 4,
						URI:  "/path/to/steps",
					},
				},
			})
			Expect(err).NotTo(HaveOccurred())
			testCaseRunner, err := runner.NewTestCaseRunner(&runner.NewTestCaseRunnerOptions{
				ID: "testCase1",
				Pickle: &gherkin.Pickle{
					Locations: []gherkin.Location{{Line: 1}},
					Steps: []*gherkin.PickleStep{
						{
							Locations: []gherkin.Location{{Line: 2}},
							Text:      "I have 100 cukes",
						},
					},
				},
				SendCommand:                 sendCommand,
				SendCommandAndAwaitResponse: sendCommandAndAwaitResponse,
				SupportCodeLibrary:          supportCodeLibrary,
				URI:                         "/path/to/feature",
			})
			Expect(err).NotTo(HaveOccurred())
			result = testCaseRunner.Run()
		})

		It("returns a ambiguous result", func() {
			Expect(result).To(Equal(&dto.TestResult{
				Status:  dto.StatusAmbiguous,
				Message: expectedMessage,
			}))
		})

		It("sends 6 commands", func() {
			Expect(allCommandsSent).To(HaveLen(6))
		})

		It("sends the test case prepared event command without an action location", func() {
			Expect(allCommandsSent[0]).To(Equal(&dto.Command{
				Type: dto.CommandTypeEvent,
				Event: &event.TestCasePrepared{
					SourceLocation: &event.Location{
						URI:  "/path/to/feature",
						Line: 1,
					},
					Steps: []*event.TestCasePreparedStep{
						{
							SourceLocation: &event.Location{
								URI:  "/path/to/feature",
								Line: 2,
							},
						},
					},
				},
			}))
		})

		It("sends the test step finished event command", func() {
			Expect(allCommandsSent[4]).To(Equal(&dto.Command{
				Type: dto.CommandTypeEvent,
				Event: &event.TestStepFinished{
					Index: 0,
					Result: &dto.TestResult{
						Status:  dto.StatusAmbiguous,
						Message: expectedMessage,
					},
					TestCase: &event.TestCase{
						SourceLocation: &event.Location{
							URI:  "/path/to/feature",
							Line: 1,
						},
					},
				},
			}))
		})

		It("sends the test case finished event command", func() {
			Expect(allCommandsSent[5]).To(Equal(&dto.Command{
				Type: dto.CommandTypeEvent,
				Event: &event.TestCaseFinished{
					SourceLocation: &event.Location{
						URI:  "/path/to/feature",
						Line: 1,
					},
					Result: &dto.TestResult{
						Status:  dto.StatusAmbiguous,
						Message: expectedMessage,
					},
				},
			}))
		})
	})

	Context("with a undefined step", func() {
		var allCommandsSent []*dto.Command
		var result *dto.TestResult

		BeforeEach(func() {
			allCommandsSent = []*dto.Command{}
			sendCommand := func(command *dto.Command) {
				allCommandsSent = append(allCommandsSent, command)
			}
			sendCommandAndAwaitResponse := func(command *dto.Command) *dto.Command {
				allCommandsSent = append(allCommandsSent, command)
				return &dto.Command{
					Type: dto.CommandTypeActionComplete,
				}
			}
			supportCodeLibrary, err := runner.NewSupportCodeLibrary(&dto.SupportCodeConfig{})
			Expect(err).NotTo(HaveOccurred())
			testCaseRunner, err := runner.NewTestCaseRunner(&runner.NewTestCaseRunnerOptions{
				ID: "testCase1",
				Pickle: &gherkin.Pickle{
					Locations: []gherkin.Location{{Line: 1}},
					Steps: []*gherkin.PickleStep{
						{
							Locations: []gherkin.Location{{Line: 2}},
							Text:      "I have 100 cukes",
						},
					},
				},
				SendCommand:                 sendCommand,
				SendCommandAndAwaitResponse: sendCommandAndAwaitResponse,
				SupportCodeLibrary:          supportCodeLibrary,
				URI:                         "/path/to/feature",
			})
			Expect(err).NotTo(HaveOccurred())
			result = testCaseRunner.Run()
		})

		It("returns a undefined result", func() {
			Expect(result).To(Equal(&dto.TestResult{
				Status: dto.StatusUndefined,
			}))
		})

		It("sends 6 commands", func() {
			Expect(allCommandsSent).To(HaveLen(6))
		})

		It("sends the test case prepared event command without an action location", func() {
			Expect(allCommandsSent[0]).To(Equal(&dto.Command{
				Type: dto.CommandTypeEvent,
				Event: &event.TestCasePrepared{
					SourceLocation: &event.Location{
						URI:  "/path/to/feature",
						Line: 1,
					},
					Steps: []*event.TestCasePreparedStep{
						{
							SourceLocation: &event.Location{
								URI:  "/path/to/feature",
								Line: 2,
							},
						},
					},
				},
			}))
		})

		It("sends the test step finished event command with status undefined", func() {
			Expect(allCommandsSent[4]).To(Equal(&dto.Command{
				Type: dto.CommandTypeEvent,
				Event: &event.TestStepFinished{
					Index: 0,
					Result: &dto.TestResult{
						Status: dto.StatusUndefined,
						// TODO message should be the snippet retrieved from generate snippet command
					},
					TestCase: &event.TestCase{
						SourceLocation: &event.Location{
							URI:  "/path/to/feature",
							Line: 1,
						},
					},
				},
			}))
		})

		It("sends the test case finished event command with status undefined", func() {
			Expect(allCommandsSent[5]).To(Equal(&dto.Command{
				Type: dto.CommandTypeEvent,
				Event: &event.TestCaseFinished{
					SourceLocation: &event.Location{
						URI:  "/path/to/feature",
						Line: 1,
					},
					Result: &dto.TestResult{
						Status: dto.StatusUndefined,
					},
				},
			}))
		})
	})
})
