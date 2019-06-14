package runner_test

import (
	"github.com/cucumber/cucumber-engine/src/runner"
	"github.com/cucumber/cucumber-engine/test/helpers"
	. "github.com/cucumber/cucumber-engine/test/matchers"
	messages "github.com/cucumber/cucumber-messages-go/v3"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("TestCaseRunner", func() {
	Context("with a passing step", func() {
		var allMessagesSent []*messages.Envelope
		var pickle *messages.Pickle
		var result *messages.TestResult

		BeforeEach(func() {
			allMessagesSent = []*messages.Envelope{}
			sendCommand := func(command *messages.Envelope) {
				allMessagesSent = append(allMessagesSent, command)
			}
			sendCommandAndAwaitResponse := func(incoming *messages.Envelope) *messages.Envelope {
				sendCommand(incoming)
				switch x := incoming.Message.(type) {
				case *messages.Envelope_CommandRunTestStep:
					return helpers.CreateActionCompleteMessageWithTestResult(
						x.CommandRunTestStep.ActionId,
						&messages.TestResult{
							DurationNanoseconds: 9,
							Status:              messages.TestResult_PASSED,
						},
					)
				default:
					return helpers.CreateActionCompleteMessage("")
				}
			}
			supportCodeLibrary, err := runner.NewSupportCodeLibrary(&messages.SupportCodeConfig{
				StepDefinitionConfigs: []*messages.StepDefinitionConfig{
					{
						Id: "step1",
						Pattern: &messages.StepDefinitionPattern{
							Source: "I have {int} cukes",
							Type:   messages.StepDefinitionPatternType_CUCUMBER_EXPRESSION,
						},
						Location: &messages.SourceReference{
							Uri:      "/path/to/steps",
							Location: &messages.Location{Line: 3},
						},
					},
				},
			})
			Expect(err).NotTo(HaveOccurred())
			pickle = &messages.Pickle{
				Locations: []*messages.Location{{Line: 1}},
				Steps: []*messages.Pickle_PickleStep{
					{
						Locations: []*messages.Location{{Line: 2}},
						Text:      "I have 100 cukes",
					},
				},
				Uri: "/path/to/feature",
			}
			testCaseRunner, err := runner.NewTestCaseRunner(&runner.NewTestCaseRunnerOptions{
				Pickle:                      pickle,
				SendCommand:                 sendCommand,
				SendCommandAndAwaitResponse: sendCommandAndAwaitResponse,
				SupportCodeLibrary:          supportCodeLibrary,
			})
			Expect(err).NotTo(HaveOccurred())
			result = testCaseRunner.Run()
		})

		It("returns a passing result", func() {
			Expect(result).To(Equal(&messages.TestResult{
				DurationNanoseconds: 9,
				Status:              messages.TestResult_PASSED,
			}))
		})

		It("sends 7 commands", func() {
			Expect(allMessagesSent).To(HaveLen(7))
		})

		It("sends the test case prepared event command", func() {
			Expect(allMessagesSent[0]).To(Equal(&messages.Envelope{
				Message: &messages.Envelope_TestCasePrepared{
					TestCasePrepared: &messages.TestCasePrepared{
						PickleId: "",
						Steps: []*messages.TestCasePreparedStep{
							{
								SourceLocation: &messages.SourceReference{
									Uri:      "/path/to/feature",
									Location: &messages.Location{Line: 2},
								},
								ActionLocation: &messages.SourceReference{
									Uri:      "/path/to/steps",
									Location: &messages.Location{Line: 3},
								},
							},
						},
					},
				},
			}))
		})

		It("sends the test case started event command", func() {
			Expect(allMessagesSent[1]).To(Equal(&messages.Envelope{
				Message: &messages.Envelope_TestCaseStarted{
					TestCaseStarted: &messages.TestCaseStarted{
						PickleId: "",
					},
				},
			}))
		})

		It("sends the initialize test case command", func() {
			Expect(allMessagesSent[2]).To(Equal(&messages.Envelope{
				Message: &messages.Envelope_CommandInitializeTestCase{
					CommandInitializeTestCase: &messages.CommandInitializeTestCase{
						Pickle: pickle,
					},
				},
			}))
		})

		It("sends the test step started event commands", func() {
			Expect(allMessagesSent[3]).To(Equal(&messages.Envelope{
				Message: &messages.Envelope_TestStepStarted{
					TestStepStarted: &messages.TestStepStarted{
						PickleId: "",
						Index:    0,
					},
				},
			}))
		})

		It("sends the run test step command", func() {
			Expect(allMessagesSent[4]).To(Equal(&messages.Envelope{
				Message: &messages.Envelope_CommandRunTestStep{
					CommandRunTestStep: &messages.CommandRunTestStep{
						StepDefinitionId: "step1",
						PatternMatches: []*messages.PatternMatch{
							{
								Captures:          []string{"100"},
								ParameterTypeName: "int",
							},
						},
					},
				},
			}))
		})

		It("sends the test step finished event command", func() {
			Expect(allMessagesSent[5]).To(Equal(&messages.Envelope{
				Message: &messages.Envelope_TestStepFinished{
					TestStepFinished: &messages.TestStepFinished{
						PickleId: "",
						Index:    0,
						TestResult: &messages.TestResult{
							DurationNanoseconds: 9,
							Status:              messages.TestResult_PASSED,
						},
					},
				},
			}))
		})

		It("sends the test case finished event command", func() {
			Expect(allMessagesSent[6]).To(Equal(&messages.Envelope{
				Message: &messages.Envelope_TestCaseFinished{
					TestCaseFinished: &messages.TestCaseFinished{
						PickleId: "",
						TestResult: &messages.TestResult{
							DurationNanoseconds: 9,
							Status:              messages.TestResult_PASSED,
						},
					},
				},
			}))
		})
	})

	Context("with a failing step", func() {
		var allMessagesSent []*messages.Envelope
		var result *messages.TestResult

		BeforeEach(func() {
			allMessagesSent = []*messages.Envelope{}
			sendCommand := func(incoming *messages.Envelope) {
				allMessagesSent = append(allMessagesSent, incoming)
			}
			sendCommandAndAwaitResponse := func(incoming *messages.Envelope) *messages.Envelope {
				sendCommand(incoming)
				switch x := incoming.Message.(type) {
				case *messages.Envelope_CommandRunTestStep:
					return helpers.CreateActionCompleteMessageWithTestResult(
						x.CommandRunTestStep.ActionId,
						&messages.TestResult{
							Status:              messages.TestResult_FAILED,
							DurationNanoseconds: 8,
							Message:             "error message and stacktrace",
						},
					)
				default:
					return helpers.CreateActionCompleteMessage("")
				}
			}
			supportCodeLibrary, err := runner.NewSupportCodeLibrary(&messages.SupportCodeConfig{
				StepDefinitionConfigs: []*messages.StepDefinitionConfig{
					{
						Id: "step1",
						Pattern: &messages.StepDefinitionPattern{
							Source: "I have {int} cukes",
							Type:   messages.StepDefinitionPatternType_CUCUMBER_EXPRESSION,
						},
						Location: &messages.SourceReference{
							Uri:      "/path/to/steps",
							Location: &messages.Location{Line: 3},
						},
					},
				},
			})
			Expect(err).NotTo(HaveOccurred())
			testCaseRunner, err := runner.NewTestCaseRunner(&runner.NewTestCaseRunnerOptions{
				Pickle: &messages.Pickle{
					Locations: []*messages.Location{{Line: 1}},
					Steps: []*messages.Pickle_PickleStep{
						{
							Locations: []*messages.Location{{Line: 2}},
							Text:      "I have 100 cukes",
						},
					},
					Uri: "/path/to/feature",
				},
				SendCommand:                 sendCommand,
				SendCommandAndAwaitResponse: sendCommandAndAwaitResponse,
				SupportCodeLibrary:          supportCodeLibrary,
			})
			Expect(err).NotTo(HaveOccurred())
			result = testCaseRunner.Run()
		})

		It("returns a failing result", func() {
			Expect(result).To(Equal(&messages.TestResult{
				Status:              messages.TestResult_FAILED,
				DurationNanoseconds: 8,
				Message:             "error message and stacktrace",
			}))
		})

		It("sends 7 commands", func() {
			Expect(allMessagesSent).To(HaveLen(7))
			Expect(allMessagesSent[0]).To(BeAMessageOfType(&messages.TestCasePrepared{}))
			Expect(allMessagesSent[1]).To(BeAMessageOfType(&messages.TestCaseStarted{}))
			Expect(allMessagesSent[2]).To(BeAMessageOfType(&messages.CommandInitializeTestCase{}))
			Expect(allMessagesSent[3]).To(BeAMessageOfType(&messages.TestStepStarted{}))
			Expect(allMessagesSent[4]).To(BeAMessageOfType(&messages.CommandRunTestStep{}))
			Expect(allMessagesSent[5]).To(BeAMessageOfType(&messages.TestStepFinished{}))
			Expect(allMessagesSent[6]).To(BeAMessageOfType(&messages.TestCaseFinished{}))
		})

		It("sends the test step finished event command with status failed", func() {
			Expect(allMessagesSent[5]).To(Equal(&messages.Envelope{
				Message: &messages.Envelope_TestStepFinished{
					TestStepFinished: &messages.TestStepFinished{
						PickleId: "",
						Index:    0,
						TestResult: &messages.TestResult{
							Status:              messages.TestResult_FAILED,
							DurationNanoseconds: 8,
							Message:             "error message and stacktrace",
						},
					},
				},
			}))
		})

		It("sends the test case finished event command with status failed", func() {
			Expect(allMessagesSent[6]).To(Equal(&messages.Envelope{
				Message: &messages.Envelope_TestCaseFinished{
					TestCaseFinished: &messages.TestCaseFinished{
						PickleId: "",
						TestResult: &messages.TestResult{
							Status:              messages.TestResult_FAILED,
							DurationNanoseconds: 8,
							Message:             "error message and stacktrace",
						},
					},
				},
			}))
		})
	})

	Context("with a ambiguous step", func() {
		var allMessagesSent []*messages.Envelope
		var result *messages.TestResult
		expectedMessage := "Multiple step definitions match:\n" +
			"  'I have {int} cukes'   - /path/to/steps:3  \n" +
			`  '^I have (\d+) cukes$' - /path/to/steps:4  ` + "\n"

		BeforeEach(func() {
			allMessagesSent = []*messages.Envelope{}
			sendCommand := func(incoming *messages.Envelope) {
				allMessagesSent = append(allMessagesSent, incoming)
			}
			sendCommandAndAwaitResponse := func(incoming *messages.Envelope) *messages.Envelope {
				sendCommand(incoming)
				return helpers.CreateActionCompleteMessage("")
			}
			supportCodeLibrary, err := runner.NewSupportCodeLibrary(&messages.SupportCodeConfig{
				StepDefinitionConfigs: []*messages.StepDefinitionConfig{
					{
						Id: "step1",
						Pattern: &messages.StepDefinitionPattern{
							Source: "I have {int} cukes",
							Type:   messages.StepDefinitionPatternType_CUCUMBER_EXPRESSION,
						},
						Location: &messages.SourceReference{
							Uri:      "/path/to/steps",
							Location: &messages.Location{Line: 3},
						},
					},
					{
						Id: "step1",
						Pattern: &messages.StepDefinitionPattern{
							Source: `^I have (\d+) cukes$`,
							Type:   messages.StepDefinitionPatternType_REGULAR_EXPRESSION,
						},
						Location: &messages.SourceReference{
							Uri:      "/path/to/steps",
							Location: &messages.Location{Line: 4},
						},
					},
				},
			})
			Expect(err).NotTo(HaveOccurred())
			testCaseRunner, err := runner.NewTestCaseRunner(&runner.NewTestCaseRunnerOptions{
				Pickle: &messages.Pickle{
					Locations: []*messages.Location{{Line: 1}},
					Steps: []*messages.Pickle_PickleStep{
						{
							Locations: []*messages.Location{{Line: 2}},
							Text:      "I have 100 cukes",
						},
					},
					Uri: "/path/to/feature",
				},
				SendCommand:                 sendCommand,
				SendCommandAndAwaitResponse: sendCommandAndAwaitResponse,
				SupportCodeLibrary:          supportCodeLibrary,
			})
			Expect(err).NotTo(HaveOccurred())
			result = testCaseRunner.Run()
		})

		It("returns a ambiguous result", func() {
			Expect(result).To(Equal(&messages.TestResult{
				Status:  messages.TestResult_AMBIGUOUS,
				Message: expectedMessage,
			}))
		})

		It("sends 6 commands", func() {
			Expect(allMessagesSent).To(HaveLen(6))
			Expect(allMessagesSent[0]).To(BeAMessageOfType(&messages.TestCasePrepared{}))
			Expect(allMessagesSent[1]).To(BeAMessageOfType(&messages.TestCaseStarted{}))
			Expect(allMessagesSent[2]).To(BeAMessageOfType(&messages.CommandInitializeTestCase{}))
			Expect(allMessagesSent[3]).To(BeAMessageOfType(&messages.TestStepStarted{}))
			Expect(allMessagesSent[4]).To(BeAMessageOfType(&messages.TestStepFinished{}))
			Expect(allMessagesSent[5]).To(BeAMessageOfType(&messages.TestCaseFinished{}))
		})

		It("sends the test case prepared event command without an action location", func() {
			Expect(allMessagesSent[0]).To(Equal(&messages.Envelope{
				Message: &messages.Envelope_TestCasePrepared{
					TestCasePrepared: &messages.TestCasePrepared{
						PickleId: "",
						Steps: []*messages.TestCasePreparedStep{
							{
								SourceLocation: &messages.SourceReference{
									Uri:      "/path/to/feature",
									Location: &messages.Location{Line: 2},
								},
							},
						},
					},
				},
			}))
		})

		It("sends the test step finished event command", func() {
			Expect(allMessagesSent[4]).To(Equal(&messages.Envelope{
				Message: &messages.Envelope_TestStepFinished{
					TestStepFinished: &messages.TestStepFinished{
						PickleId: "",
						Index:    0,
						TestResult: &messages.TestResult{
							Status:  messages.TestResult_AMBIGUOUS,
							Message: expectedMessage,
						},
					},
				},
			}))
		})

		It("sends the test case finished event command", func() {
			Expect(allMessagesSent[5]).To(Equal(&messages.Envelope{
				Message: &messages.Envelope_TestCaseFinished{
					TestCaseFinished: &messages.TestCaseFinished{
						PickleId: "",
						TestResult: &messages.TestResult{
							Status:  messages.TestResult_AMBIGUOUS,
							Message: expectedMessage,
						},
					},
				},
			}))
		})
	})

	Context("with a ambiguous step and base directory", func() {
		var allMessagesSent []*messages.Envelope
		var result *messages.TestResult

		BeforeEach(func() {
			allMessagesSent = []*messages.Envelope{}
			sendCommand := func(incoming *messages.Envelope) {
				allMessagesSent = append(allMessagesSent, incoming)
			}
			sendCommandAndAwaitResponse := func(incoming *messages.Envelope) *messages.Envelope {
				sendCommand(incoming)
				return helpers.CreateActionCompleteMessage("")
			}
			supportCodeLibrary, err := runner.NewSupportCodeLibrary(&messages.SupportCodeConfig{
				StepDefinitionConfigs: []*messages.StepDefinitionConfig{
					{
						Id: "step1",
						Pattern: &messages.StepDefinitionPattern{
							Source: "I have {int} cukes",
							Type:   messages.StepDefinitionPatternType_CUCUMBER_EXPRESSION,
						},
						Location: &messages.SourceReference{
							Uri:      "/path/to/base/path/to/steps",
							Location: &messages.Location{Line: 3},
						},
					},
					{
						Id: "step1",
						Pattern: &messages.StepDefinitionPattern{
							Source: `^I have (\d+) cukes$`,
							Type:   messages.StepDefinitionPatternType_REGULAR_EXPRESSION,
						},
						Location: &messages.SourceReference{
							Uri:      "/path/to/base/path/to/steps",
							Location: &messages.Location{Line: 4},
						},
					},
				},
			})
			Expect(err).NotTo(HaveOccurred())
			testCaseRunner, err := runner.NewTestCaseRunner(&runner.NewTestCaseRunnerOptions{
				BaseDirectory: "/path/to/base",
				Pickle: &messages.Pickle{
					Locations: []*messages.Location{{Line: 1}},
					Steps: []*messages.Pickle_PickleStep{
						{
							Locations: []*messages.Location{{Line: 2}},
							Text:      "I have 100 cukes",
						},
					},
					Uri: "/path/to/base/path/to/feature",
				},
				SendCommand:                 sendCommand,
				SendCommandAndAwaitResponse: sendCommandAndAwaitResponse,
				SupportCodeLibrary:          supportCodeLibrary,
			})
			Expect(err).NotTo(HaveOccurred())
			result = testCaseRunner.Run()
		})

		It("returns a ambiguous result with the paths ", func() {
			Expect(result).To(Equal(&messages.TestResult{
				Status: messages.TestResult_AMBIGUOUS,
				Message: "Multiple step definitions match:\n" +
					"  'I have {int} cukes'   - path/to/steps:3  \n" +
					`  '^I have (\d+) cukes$' - path/to/steps:4  ` + "\n",
			}))
		})

		It("sends 6 commands", func() {
			Expect(allMessagesSent).To(HaveLen(6))
			Expect(allMessagesSent[0]).To(BeAMessageOfType(&messages.TestCasePrepared{}))
			Expect(allMessagesSent[1]).To(BeAMessageOfType(&messages.TestCaseStarted{}))
			Expect(allMessagesSent[2]).To(BeAMessageOfType(&messages.CommandInitializeTestCase{}))
			Expect(allMessagesSent[3]).To(BeAMessageOfType(&messages.TestStepStarted{}))
			Expect(allMessagesSent[4]).To(BeAMessageOfType(&messages.TestStepFinished{}))
			Expect(allMessagesSent[5]).To(BeAMessageOfType(&messages.TestCaseFinished{}))
		})
	})

	Context("with a undefined step", func() {
		var allMessagesSent []*messages.Envelope
		var result *messages.TestResult
		var snippet = "snippet line1\nsnippet line2\nsnippet line3"

		BeforeEach(func() {
			allMessagesSent = []*messages.Envelope{}
			sendCommand := func(incoming *messages.Envelope) {
				allMessagesSent = append(allMessagesSent, incoming)
			}
			sendCommandAndAwaitResponse := func(incoming *messages.Envelope) *messages.Envelope {
				sendCommand(incoming)
				switch x := incoming.Message.(type) {
				case *messages.Envelope_CommandGenerateSnippet:
					return helpers.CreateActionCompleteMessageWithSnippet(
						x.CommandGenerateSnippet.ActionId,
						snippet,
					)
				default:
					return helpers.CreateActionCompleteMessage("")
				}
			}
			supportCodeLibrary, err := runner.NewSupportCodeLibrary(&messages.SupportCodeConfig{})
			Expect(err).NotTo(HaveOccurred())
			testCaseRunner, err := runner.NewTestCaseRunner(&runner.NewTestCaseRunnerOptions{
				Pickle: &messages.Pickle{
					Locations: []*messages.Location{{Line: 1}},
					Steps: []*messages.Pickle_PickleStep{
						{
							Locations: []*messages.Location{{Line: 2}},
							Text:      "I have 100 cukes",
						},
					},
					Uri: "/path/to/feature",
				},
				SendCommand:                 sendCommand,
				SendCommandAndAwaitResponse: sendCommandAndAwaitResponse,
				SupportCodeLibrary:          supportCodeLibrary,
			})
			Expect(err).NotTo(HaveOccurred())
			result = testCaseRunner.Run()
		})

		It("returns a undefined result", func() {
			Expect(result).To(Equal(&messages.TestResult{
				Status:  messages.TestResult_UNDEFINED,
				Message: snippet,
			}))
		})

		It("sends 7 commands", func() {
			Expect(allMessagesSent).To(HaveLen(7))
			Expect(allMessagesSent[0]).To(BeAMessageOfType(&messages.TestCasePrepared{}))
			Expect(allMessagesSent[1]).To(BeAMessageOfType(&messages.TestCaseStarted{}))
			Expect(allMessagesSent[2]).To(BeAMessageOfType(&messages.CommandInitializeTestCase{}))
			Expect(allMessagesSent[3]).To(BeAMessageOfType(&messages.TestStepStarted{}))
			Expect(allMessagesSent[4]).To(BeAMessageOfType(&messages.CommandGenerateSnippet{}))
			Expect(allMessagesSent[5]).To(BeAMessageOfType(&messages.TestStepFinished{}))
			Expect(allMessagesSent[6]).To(BeAMessageOfType(&messages.TestCaseFinished{}))
		})

		It("sends the test case prepared event command without an action location", func() {
			Expect(allMessagesSent[0]).To(Equal(&messages.Envelope{
				Message: &messages.Envelope_TestCasePrepared{
					TestCasePrepared: &messages.TestCasePrepared{
						PickleId: "",
						Steps: []*messages.TestCasePreparedStep{
							{
								SourceLocation: &messages.SourceReference{
									Uri:      "/path/to/feature",
									Location: &messages.Location{Line: 2},
								},
							},
						},
					},
				},
			}))
		})

		It("sends the generate snippet command", func() {
			Expect(allMessagesSent[4]).To(Equal(&messages.Envelope{
				Message: &messages.Envelope_CommandGenerateSnippet{
					CommandGenerateSnippet: &messages.CommandGenerateSnippet{
						GeneratedExpressions: []*messages.GeneratedExpression{
							{
								Text:               "I have {int} cukes",
								ParameterTypeNames: []string{"int"},
							},
						},
					},
				},
			}))
		})

		It("sends the test step finished event command with status undefined", func() {
			Expect(allMessagesSent[5]).To(Equal(&messages.Envelope{
				Message: &messages.Envelope_TestStepFinished{
					TestStepFinished: &messages.TestStepFinished{
						Index: 0,
						TestResult: &messages.TestResult{
							Status:  messages.TestResult_UNDEFINED,
							Message: snippet,
						},
					},
				},
			}))
		})

		It("sends the test case finished event command with status undefined", func() {
			Expect(allMessagesSent[6]).To(Equal(&messages.Envelope{
				Message: &messages.Envelope_TestCaseFinished{
					TestCaseFinished: &messages.TestCaseFinished{
						TestResult: &messages.TestResult{
							Status:  messages.TestResult_UNDEFINED,
							Message: snippet,
						},
					},
				},
			}))
		})
	})

	Context("with a failing and then skipped step", func() {
		var allMessagesSent []*messages.Envelope
		var result *messages.TestResult

		BeforeEach(func() {
			allMessagesSent = []*messages.Envelope{}
			sendCommand := func(incoming *messages.Envelope) {
				allMessagesSent = append(allMessagesSent, incoming)
			}
			sendCommandAndAwaitResponse := func(incoming *messages.Envelope) *messages.Envelope {
				sendCommand(incoming)
				switch x := incoming.Message.(type) {
				case *messages.Envelope_CommandRunTestStep:
					return helpers.CreateActionCompleteMessageWithTestResult(
						x.CommandRunTestStep.ActionId,
						&messages.TestResult{
							Status:              messages.TestResult_FAILED,
							DurationNanoseconds: 8,
							Message:             "error message and stacktrace",
						},
					)
				default:
					return helpers.CreateActionCompleteMessage("")
				}
			}
			supportCodeLibrary, err := runner.NewSupportCodeLibrary(&messages.SupportCodeConfig{
				StepDefinitionConfigs: []*messages.StepDefinitionConfig{
					{
						Id: "step1",
						Pattern: &messages.StepDefinitionPattern{
							Source: "I have {int} cukes",
							Type:   messages.StepDefinitionPatternType_CUCUMBER_EXPRESSION,
						},
						Location: &messages.SourceReference{
							Uri:      "/path/to/steps",
							Location: &messages.Location{Line: 3},
						},
					},
				},
			})
			Expect(err).NotTo(HaveOccurred())
			testCaseRunner, err := runner.NewTestCaseRunner(&runner.NewTestCaseRunnerOptions{
				Pickle: &messages.Pickle{
					Locations: []*messages.Location{{Line: 1}},
					Steps: []*messages.Pickle_PickleStep{
						{
							Locations: []*messages.Location{{Line: 2}},
							Text:      "I have 100 cukes",
						},
						{
							Locations: []*messages.Location{{Line: 3}},
							Text:      "I have 101 cukes",
						},
					},
					Uri: "/path/to/feature",
				},
				SendCommand:                 sendCommand,
				SendCommandAndAwaitResponse: sendCommandAndAwaitResponse,
				SupportCodeLibrary:          supportCodeLibrary,
			})
			Expect(err).NotTo(HaveOccurred())
			result = testCaseRunner.Run()
		})

		It("returns a failing result", func() {
			Expect(result).To(Equal(&messages.TestResult{
				Status:              messages.TestResult_FAILED,
				DurationNanoseconds: 8,
				Message:             "error message and stacktrace",
			}))
		})

		It("sends 9 commands", func() {
			Expect(allMessagesSent).To(HaveLen(9))
			Expect(allMessagesSent[0]).To(BeAMessageOfType(&messages.TestCasePrepared{}))
			Expect(allMessagesSent[1]).To(BeAMessageOfType(&messages.TestCaseStarted{}))
			Expect(allMessagesSent[2]).To(BeAMessageOfType(&messages.CommandInitializeTestCase{}))
			Expect(allMessagesSent[3]).To(BeAMessageOfType(&messages.TestStepStarted{}))
			Expect(allMessagesSent[4]).To(BeAMessageOfType(&messages.CommandRunTestStep{}))
			Expect(allMessagesSent[5]).To(BeAMessageOfType(&messages.TestStepFinished{}))
			Expect(allMessagesSent[6]).To(BeAMessageOfType(&messages.TestStepStarted{}))
			Expect(allMessagesSent[7]).To(BeAMessageOfType(&messages.TestStepFinished{}))
			Expect(allMessagesSent[8]).To(BeAMessageOfType(&messages.TestCaseFinished{}))
		})

		It("sends the test step finished event command with status skipped for the second step", func() {
			Expect(allMessagesSent[7]).To(Equal(&messages.Envelope{
				Message: &messages.Envelope_TestStepFinished{
					TestStepFinished: &messages.TestStepFinished{
						Index: 1,
						TestResult: &messages.TestResult{
							Status: messages.TestResult_SKIPPED,
						},
					},
				},
			}))
		})
	})

	Context("isSkipped is true (fail fast or dry run)", func() {
		var allMessagesSent []*messages.Envelope
		var result *messages.TestResult

		BeforeEach(func() {
			allMessagesSent = []*messages.Envelope{}
			sendCommand := func(incoming *messages.Envelope) {
				allMessagesSent = append(allMessagesSent, incoming)
			}
			sendCommandAndAwaitResponse := func(incoming *messages.Envelope) *messages.Envelope {
				sendCommand(incoming)
				return helpers.CreateActionCompleteMessage("")
			}
			supportCodeLibrary, err := runner.NewSupportCodeLibrary(&messages.SupportCodeConfig{
				BeforeTestCaseHookDefinitionConfigs: []*messages.TestCaseHookDefinitionConfig{
					{
						Id: "beforeHook1",
						Location: &messages.SourceReference{
							Uri:      "/path/to/hooks",
							Location: &messages.Location{Line: 11},
						},
					},
				},
				AfterTestCaseHookDefinitionConfigs: []*messages.TestCaseHookDefinitionConfig{
					{
						Id: "afterHook1",
						Location: &messages.SourceReference{
							Uri:      "/path/to/hooks",
							Location: &messages.Location{Line: 12},
						},
					},
				},
				StepDefinitionConfigs: []*messages.StepDefinitionConfig{
					{
						Id: "step1",
						Pattern: &messages.StepDefinitionPattern{
							Source: "I have {int} cukes",
							Type:   messages.StepDefinitionPatternType_CUCUMBER_EXPRESSION,
						},
						Location: &messages.SourceReference{
							Uri:      "/path/to/steps",
							Location: &messages.Location{Line: 3},
						},
					},
				},
			})
			Expect(err).NotTo(HaveOccurred())
			testCaseRunner, err := runner.NewTestCaseRunner(&runner.NewTestCaseRunnerOptions{
				IsSkipped: true,
				Pickle: &messages.Pickle{
					Locations: []*messages.Location{{Line: 1}},
					Steps: []*messages.Pickle_PickleStep{
						{
							Locations: []*messages.Location{{Line: 2}},
							Text:      "I have 100 cukes",
						},
					},
					Uri: "/path/to/feature",
				},
				SendCommand:                 sendCommand,
				SendCommandAndAwaitResponse: sendCommandAndAwaitResponse,
				SupportCodeLibrary:          supportCodeLibrary,
			})
			Expect(err).NotTo(HaveOccurred())
			result = testCaseRunner.Run()
		})

		It("returns a skipped result", func() {
			Expect(result).To(Equal(&messages.TestResult{
				Status: messages.TestResult_SKIPPED,
			}))
		})

		It("sends 9 commands", func() {
			Expect(allMessagesSent).To(HaveLen(9))
			Expect(allMessagesSent[0]).To(BeAMessageOfType(&messages.TestCasePrepared{}))
			Expect(allMessagesSent[1]).To(BeAMessageOfType(&messages.TestCaseStarted{}))
			Expect(allMessagesSent[2]).To(BeAMessageOfType(&messages.TestStepStarted{}))
			Expect(allMessagesSent[3]).To(BeAMessageOfType(&messages.TestStepFinished{}))
			Expect(allMessagesSent[4]).To(BeAMessageOfType(&messages.TestStepStarted{}))
			Expect(allMessagesSent[5]).To(BeAMessageOfType(&messages.TestStepFinished{}))
			Expect(allMessagesSent[6]).To(BeAMessageOfType(&messages.TestStepStarted{}))
			Expect(allMessagesSent[7]).To(BeAMessageOfType(&messages.TestStepFinished{}))
			Expect(allMessagesSent[8]).To(BeAMessageOfType(&messages.TestCaseFinished{}))
		})

		It("sends the test step finished event command with status skipped for the before hook", func() {
			Expect(allMessagesSent[3]).To(Equal(&messages.Envelope{
				Message: &messages.Envelope_TestStepFinished{
					TestStepFinished: &messages.TestStepFinished{
						Index: 0,
						TestResult: &messages.TestResult{
							Status: messages.TestResult_SKIPPED,
						},
					},
				},
			}))
		})

		It("sends the test step finished event command with status skipped for the step", func() {
			Expect(allMessagesSent[5]).To(Equal(&messages.Envelope{
				Message: &messages.Envelope_TestStepFinished{
					TestStepFinished: &messages.TestStepFinished{
						Index: 1,
						TestResult: &messages.TestResult{
							Status: messages.TestResult_SKIPPED,
						},
					},
				},
			}))
		})

		It("sends the test step finished event command with status skipped for the after hook", func() {
			Expect(allMessagesSent[7]).To(Equal(&messages.Envelope{
				Message: &messages.Envelope_TestStepFinished{
					TestStepFinished: &messages.TestStepFinished{
						Index: 2,
						TestResult: &messages.TestResult{
							Status: messages.TestResult_SKIPPED,
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
