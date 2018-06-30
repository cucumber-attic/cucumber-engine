package runner_test

import (
	"github.com/cucumber/cucumber-engine/src/dto"
	"github.com/cucumber/cucumber-engine/src/dto/event"
	"github.com/cucumber/cucumber-engine/src/runner"
	. "github.com/cucumber/cucumber-engine/test/matchers"
	gherkin "github.com/cucumber/gherkin-go"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("TestCaseRunner", func() {
	Context("with a passing step", func() {
		var allCommandsSent []*dto.Command
		var pickle *gherkin.Pickle
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
						Result: &dto.TestResult{
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
			pickle = &gherkin.Pickle{
				Locations: []gherkin.Location{{Line: 1}},
				Steps: []*gherkin.PickleStep{
					{
						Locations: []gherkin.Location{{Line: 2}},
						Text:      "I have 100 cukes",
					},
				},
			}
			testCaseRunner, err := runner.NewTestCaseRunner(&runner.NewTestCaseRunnerOptions{
				ID:                          "testCase1",
				Pickle:                      pickle,
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
					SourceLocation: &dto.Location{
						URI:  "/path/to/feature",
						Line: 1,
					},
					Steps: []*event.TestCasePreparedStep{
						{
							SourceLocation: &dto.Location{
								URI:  "/path/to/feature",
								Line: 2,
							},
							ActionLocation: &dto.Location{
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
					SourceLocation: &dto.Location{
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
				TestCase: &dto.TestCase{
					SourceLocation: &dto.Location{
						URI:  "/path/to/feature",
						Line: 1,
					},
				},
				Pickle: pickle,
			}))
		})

		It("sends the test step started event commands", func() {
			Expect(allCommandsSent[3]).To(Equal(&dto.Command{
				Type: dto.CommandTypeEvent,
				Event: &event.TestStepStarted{
					Index: 0,
					TestCase: &dto.TestCase{
						SourceLocation: &dto.Location{
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
					TestCase: &dto.TestCase{
						SourceLocation: &dto.Location{
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
					SourceLocation: &dto.Location{
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
						Result: &dto.TestResult{
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
			Expect(allCommandsSent[0]).To(BeACommandWithEventAssignableToTypeOf(&event.TestCasePrepared{}))
			Expect(allCommandsSent[1]).To(BeACommandWithEventAssignableToTypeOf(&event.TestCaseStarted{}))
			Expect(allCommandsSent[2]).To(BeACommandWithType(dto.CommandTypeInitializeTestCase))
			Expect(allCommandsSent[3]).To(BeACommandWithEventAssignableToTypeOf(&event.TestStepStarted{}))
			Expect(allCommandsSent[4]).To(BeACommandWithType(dto.CommandTypeRunTestStep))
			Expect(allCommandsSent[5]).To(BeACommandWithEventAssignableToTypeOf(&event.TestStepFinished{}))
			Expect(allCommandsSent[6]).To(BeACommandWithEventAssignableToTypeOf(&event.TestCaseFinished{}))
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
					TestCase: &dto.TestCase{
						SourceLocation: &dto.Location{
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
					SourceLocation: &dto.Location{
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
			"  'I have {int} cukes'   - /path/to/steps:3  \n" +
			`  '^I have (\d+) cukes$' - /path/to/steps:4  ` + "\n"

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
			Expect(allCommandsSent[0]).To(BeACommandWithEventAssignableToTypeOf(&event.TestCasePrepared{}))
			Expect(allCommandsSent[1]).To(BeACommandWithEventAssignableToTypeOf(&event.TestCaseStarted{}))
			Expect(allCommandsSent[2]).To(BeACommandWithType(dto.CommandTypeInitializeTestCase))
			Expect(allCommandsSent[3]).To(BeACommandWithEventAssignableToTypeOf(&event.TestStepStarted{}))
			Expect(allCommandsSent[4]).To(BeACommandWithEventAssignableToTypeOf(&event.TestStepFinished{}))
			Expect(allCommandsSent[5]).To(BeACommandWithEventAssignableToTypeOf(&event.TestCaseFinished{}))
		})

		It("sends the test case prepared event command without an action location", func() {
			Expect(allCommandsSent[0]).To(Equal(&dto.Command{
				Type: dto.CommandTypeEvent,
				Event: &event.TestCasePrepared{
					SourceLocation: &dto.Location{
						URI:  "/path/to/feature",
						Line: 1,
					},
					Steps: []*event.TestCasePreparedStep{
						{
							SourceLocation: &dto.Location{
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
					TestCase: &dto.TestCase{
						SourceLocation: &dto.Location{
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
					SourceLocation: &dto.Location{
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

	Context("with a ambiguous step and base directory", func() {
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
			supportCodeLibrary, err := runner.NewSupportCodeLibrary(&dto.SupportCodeConfig{
				StepDefinitionConfigs: []*dto.StepDefinitionConfig{
					{
						ID: "step1",
						Pattern: dto.Pattern{
							Source: "I have {int} cukes",
							Type:   dto.PatternTypeCucumberExpression,
						},
						Line: 3,
						URI:  "/path/to/base/path/to/steps",
					},
					{
						ID: "step2",
						Pattern: dto.Pattern{
							Source: `^I have (\d+) cukes$`,
							Type:   dto.PatternTypeRegularExpression,
						},
						Line: 4,
						URI:  "/path/to/base/path/to/steps",
					},
				},
			})
			Expect(err).NotTo(HaveOccurred())
			testCaseRunner, err := runner.NewTestCaseRunner(&runner.NewTestCaseRunnerOptions{
				BaseDirectory: "/path/to/base",
				ID:            "testCase1",
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
				URI:                         "/path/to/base/path/to/feature",
			})
			Expect(err).NotTo(HaveOccurred())
			result = testCaseRunner.Run()
		})

		It("returns a ambiguous result with the paths ", func() {
			Expect(result).To(Equal(&dto.TestResult{
				Status: dto.StatusAmbiguous,
				Message: "Multiple step definitions match:\n" +
					"  'I have {int} cukes'   - path/to/steps:3  \n" +
					`  '^I have (\d+) cukes$' - path/to/steps:4  ` + "\n",
			}))
		})

		It("sends 6 commands", func() {
			Expect(allCommandsSent).To(HaveLen(6))
			Expect(allCommandsSent[0]).To(BeACommandWithEventAssignableToTypeOf(&event.TestCasePrepared{}))
			Expect(allCommandsSent[1]).To(BeACommandWithEventAssignableToTypeOf(&event.TestCaseStarted{}))
			Expect(allCommandsSent[2]).To(BeACommandWithType(dto.CommandTypeInitializeTestCase))
			Expect(allCommandsSent[3]).To(BeACommandWithEventAssignableToTypeOf(&event.TestStepStarted{}))
			Expect(allCommandsSent[4]).To(BeACommandWithEventAssignableToTypeOf(&event.TestStepFinished{}))
			Expect(allCommandsSent[5]).To(BeACommandWithEventAssignableToTypeOf(&event.TestCaseFinished{}))
		})
	})

	Context("with a undefined step", func() {
		var allCommandsSent []*dto.Command
		var result *dto.TestResult
		var snippet = "snippet line1\nsnippet line2\nsnippet line3"

		BeforeEach(func() {
			allCommandsSent = []*dto.Command{}
			sendCommand := func(command *dto.Command) {
				allCommandsSent = append(allCommandsSent, command)
			}
			sendCommandAndAwaitResponse := func(command *dto.Command) *dto.Command {
				allCommandsSent = append(allCommandsSent, command)
				if command.Type == dto.CommandTypeGenerateSnippet {
					return &dto.Command{
						Type:    dto.CommandTypeActionComplete,
						Snippet: snippet,
					}
				}
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
				Status:  dto.StatusUndefined,
				Message: snippet,
			}))
		})

		It("sends 7 commands", func() {
			Expect(allCommandsSent).To(HaveLen(7))
			Expect(allCommandsSent[0]).To(BeACommandWithEventAssignableToTypeOf(&event.TestCasePrepared{}))
			Expect(allCommandsSent[1]).To(BeACommandWithEventAssignableToTypeOf(&event.TestCaseStarted{}))
			Expect(allCommandsSent[2]).To(BeACommandWithType(dto.CommandTypeInitializeTestCase))
			Expect(allCommandsSent[3]).To(BeACommandWithEventAssignableToTypeOf(&event.TestStepStarted{}))
			Expect(allCommandsSent[4]).To(BeACommandWithType(dto.CommandTypeGenerateSnippet))
			Expect(allCommandsSent[5]).To(BeACommandWithEventAssignableToTypeOf(&event.TestStepFinished{}))
			Expect(allCommandsSent[6]).To(BeACommandWithEventAssignableToTypeOf(&event.TestCaseFinished{}))
		})

		It("sends the test case prepared event command without an action location", func() {
			Expect(allCommandsSent[0]).To(Equal(&dto.Command{
				Type: dto.CommandTypeEvent,
				Event: &event.TestCasePrepared{
					SourceLocation: &dto.Location{
						URI:  "/path/to/feature",
						Line: 1,
					},
					Steps: []*event.TestCasePreparedStep{
						{
							SourceLocation: &dto.Location{
								URI:  "/path/to/feature",
								Line: 2,
							},
						},
					},
				},
			}))
		})

		It("sends the generate snippet command", func() {
			Expect(allCommandsSent[4]).To(Equal(&dto.Command{
				Type: dto.CommandTypeGenerateSnippet,
				GeneratedExpressions: []*dto.GeneratedExpression{
					{
						Text:               "I have {int} cukes",
						ParameterTypeNames: []string{"int"},
					},
				},
			}))
		})

		It("sends the test step finished event command with status undefined", func() {
			Expect(allCommandsSent[5]).To(Equal(&dto.Command{
				Type: dto.CommandTypeEvent,
				Event: &event.TestStepFinished{
					Index: 0,
					Result: &dto.TestResult{
						Status:  dto.StatusUndefined,
						Message: snippet,
					},
					TestCase: &dto.TestCase{
						SourceLocation: &dto.Location{
							URI:  "/path/to/feature",
							Line: 1,
						},
					},
				},
			}))
		})

		It("sends the test case finished event command with status undefined", func() {
			Expect(allCommandsSent[6]).To(Equal(&dto.Command{
				Type: dto.CommandTypeEvent,
				Event: &event.TestCaseFinished{
					SourceLocation: &dto.Location{
						URI:  "/path/to/feature",
						Line: 1,
					},
					Result: &dto.TestResult{
						Status:  dto.StatusUndefined,
						Message: snippet,
					},
				},
			}))
		})
	})

	Context("with a failing and then skipped step", func() {
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
						Result: &dto.TestResult{
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
						{
							Locations: []gherkin.Location{{Line: 3}},
							Text:      "I have 101 cukes",
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

		It("sends 9 commands", func() {
			Expect(allCommandsSent).To(HaveLen(9))
			Expect(allCommandsSent[0]).To(BeACommandWithEventAssignableToTypeOf(&event.TestCasePrepared{}))
			Expect(allCommandsSent[1]).To(BeACommandWithEventAssignableToTypeOf(&event.TestCaseStarted{}))
			Expect(allCommandsSent[2]).To(BeACommandWithType(dto.CommandTypeInitializeTestCase))
			Expect(allCommandsSent[3]).To(BeACommandWithEventAssignableToTypeOf(&event.TestStepStarted{}))
			Expect(allCommandsSent[4]).To(BeACommandWithType(dto.CommandTypeRunTestStep))
			Expect(allCommandsSent[5]).To(BeACommandWithEventAssignableToTypeOf(&event.TestStepFinished{}))
			Expect(allCommandsSent[6]).To(BeACommandWithEventAssignableToTypeOf(&event.TestStepStarted{}))
			Expect(allCommandsSent[7]).To(BeACommandWithEventAssignableToTypeOf(&event.TestStepFinished{}))
			Expect(allCommandsSent[8]).To(BeACommandWithEventAssignableToTypeOf(&event.TestCaseFinished{}))
		})

		It("sends the test step finished event command with status skipped for the second step", func() {
			Expect(allCommandsSent[7]).To(Equal(&dto.Command{
				Type: dto.CommandTypeEvent,
				Event: &event.TestStepFinished{
					Index:  1,
					Result: &dto.TestResult{Status: dto.StatusSkipped},
					TestCase: &dto.TestCase{
						SourceLocation: &dto.Location{
							URI:  "/path/to/feature",
							Line: 1,
						},
					},
				},
			}))
		})
	})

	Context("isSkipped is true (fail fast or dry run)", func() {
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
			supportCodeLibrary, err := runner.NewSupportCodeLibrary(&dto.SupportCodeConfig{
				BeforeTestCaseHookDefinitionConfigs: []*dto.TestCaseHookDefinitionConfig{
					{ID: "beforeHook1", URI: "/path/to/hooks", Line: 11},
				},
				AfterTestCaseHookDefinitionConfigs: []*dto.TestCaseHookDefinitionConfig{
					{ID: "afterHook1", URI: "/path/to/hooks", Line: 12},
				},
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
				ID:        "testCase1",
				IsSkipped: true,
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

		It("returns a skipped result", func() {
			Expect(result).To(Equal(&dto.TestResult{
				Status: dto.StatusSkipped,
			}))
		})

		It("sends 9 commands", func() {
			Expect(allCommandsSent).To(HaveLen(9))
			Expect(allCommandsSent[0]).To(BeACommandWithEventAssignableToTypeOf(&event.TestCasePrepared{}))
			Expect(allCommandsSent[1]).To(BeACommandWithEventAssignableToTypeOf(&event.TestCaseStarted{}))
			Expect(allCommandsSent[2]).To(BeACommandWithEventAssignableToTypeOf(&event.TestStepStarted{}))
			Expect(allCommandsSent[3]).To(BeACommandWithEventAssignableToTypeOf(&event.TestStepFinished{}))
			Expect(allCommandsSent[4]).To(BeACommandWithEventAssignableToTypeOf(&event.TestStepStarted{}))
			Expect(allCommandsSent[5]).To(BeACommandWithEventAssignableToTypeOf(&event.TestStepFinished{}))
			Expect(allCommandsSent[6]).To(BeACommandWithEventAssignableToTypeOf(&event.TestStepStarted{}))
			Expect(allCommandsSent[7]).To(BeACommandWithEventAssignableToTypeOf(&event.TestStepFinished{}))
			Expect(allCommandsSent[8]).To(BeACommandWithEventAssignableToTypeOf(&event.TestCaseFinished{}))
		})

		It("sends the test step finished event command with status skipped for the before hook", func() {
			Expect(allCommandsSent[3]).To(Equal(&dto.Command{
				Type: dto.CommandTypeEvent,
				Event: &event.TestStepFinished{
					Index:  0,
					Result: &dto.TestResult{Status: dto.StatusSkipped},
					TestCase: &dto.TestCase{
						SourceLocation: &dto.Location{
							URI:  "/path/to/feature",
							Line: 1,
						},
					},
				},
			}))
		})

		It("sends the test step finished event command with status skipped for the step", func() {
			Expect(allCommandsSent[5]).To(Equal(&dto.Command{
				Type: dto.CommandTypeEvent,
				Event: &event.TestStepFinished{
					Index:  1,
					Result: &dto.TestResult{Status: dto.StatusSkipped},
					TestCase: &dto.TestCase{
						SourceLocation: &dto.Location{
							URI:  "/path/to/feature",
							Line: 1,
						},
					},
				},
			}))
		})

		It("sends the test step finished event command with status skipped for the after hook", func() {
			Expect(allCommandsSent[7]).To(Equal(&dto.Command{
				Type: dto.CommandTypeEvent,
				Event: &event.TestStepFinished{
					Index:  2,
					Result: &dto.TestResult{Status: dto.StatusSkipped},
					TestCase: &dto.TestCase{
						SourceLocation: &dto.Location{
							URI:  "/path/to/feature",
							Line: 1,
						},
					},
				},
			}))
		})
	})

	Context("with a passing step and before hook", func() {})

	Context("with a passing step and after hook", func() {})

	Context("with a failing before hook", func() {
		// skips the steps
	})

	Context("with a failing before hook and a passing after hook", func() {
		// it runs the after hook
	})

	Context("with a failing step and a passing after hook", func() {
		// it runs the after hook
	})

	Context("with a multiple before hooks and an after hook", func() {
		// it runs the after hook but not later before hooks
	})
})
