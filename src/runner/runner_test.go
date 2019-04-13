package runner_test

import (
	"path"
	"runtime"
	"time"

	"github.com/cucumber/cucumber-engine/src/runner"
	helpers "github.com/cucumber/cucumber-engine/test/helpers"
	. "github.com/cucumber/cucumber-engine/test/matchers"
	messages "github.com/cucumber/cucumber-messages-go/v2"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Runner", func() {
	_, filename, _, _ := runtime.Caller(0)
	rootDir := path.Join(filename, "..", "..", "..")

	Context("with a feature with a single scenario with three steps", func() {
		featurePath := path.Join(rootDir, "test", "fixtures", "a.feature")

		Context("all steps are undefined", func() {
			var allMessagesSent []*messages.Wrapper

			BeforeEach(func() {
				allMessagesSent = runWithConfigAndResponder(
					&messages.SourcesConfig{
						AbsolutePaths: []string{featurePath},
						Filters:       &messages.SourcesFilterConfig{},
						Language:      "en",
						Order:         &messages.SourcesOrder{},
					},
					&messages.RuntimeConfig{
						MaxParallel: 1,
					},
					&messages.SupportCodeConfig{},
					func(commandChan chan *messages.Wrapper, incoming *messages.Wrapper) {
						switch x := incoming.Message.(type) {
						case *messages.Wrapper_CommandRunBeforeTestRunHooks:
							commandChan <- helpers.CreateActionCompleteMessage(x.CommandRunBeforeTestRunHooks.ActionId)
						case *messages.Wrapper_CommandRunAfterTestRunHooks:
							commandChan <- helpers.CreateActionCompleteMessage(x.CommandRunAfterTestRunHooks.ActionId)
						case *messages.Wrapper_CommandInitializeTestCase:
							commandChan <- helpers.CreateActionCompleteMessage(x.CommandInitializeTestCase.ActionId)
						case *messages.Wrapper_CommandGenerateSnippet:
							commandChan <- helpers.CreateActionCompleteMessageWithSnippet(x.CommandGenerateSnippet.ActionId, "snippet")
						}
					},
				)
			})

			It("sends 21 commands", func() {
				Expect(allMessagesSent).To(HaveLen(21))
				Expect(allMessagesSent[0]).To(BeAMessageOfType(&messages.Source{}))
				Expect(allMessagesSent[1]).To(BeAMessageOfType(&messages.GherkinDocument{}))
				Expect(allMessagesSent[2]).To(BeAMessageOfType(&messages.Pickle{}))
				Expect(allMessagesSent[3]).To(BeAMessageOfType(&messages.PickleAccepted{}))
				Expect(allMessagesSent[4]).To(BeAMessageOfType(&messages.TestRunStarted{}))
				Expect(allMessagesSent[5]).To(BeAMessageOfType(&messages.CommandRunBeforeTestRunHooks{}))
				Expect(allMessagesSent[6]).To(BeAMessageOfType(&messages.TestCasePrepared{}))
				Expect(allMessagesSent[7]).To(BeAMessageOfType(&messages.TestCaseStarted{}))
				Expect(allMessagesSent[8]).To(BeAMessageOfType(&messages.CommandInitializeTestCase{}))
				Expect(allMessagesSent[9]).To(BeAMessageOfType(&messages.TestStepStarted{}))
				Expect(allMessagesSent[10]).To(BeAMessageOfType(&messages.CommandGenerateSnippet{}))
				Expect(allMessagesSent[11]).To(BeAMessageOfType(&messages.TestStepFinished{}))
				Expect(allMessagesSent[12]).To(BeAMessageOfType(&messages.TestStepStarted{}))
				Expect(allMessagesSent[13]).To(BeAMessageOfType(&messages.CommandGenerateSnippet{}))
				Expect(allMessagesSent[14]).To(BeAMessageOfType(&messages.TestStepFinished{}))
				Expect(allMessagesSent[15]).To(BeAMessageOfType(&messages.TestStepStarted{}))
				Expect(allMessagesSent[16]).To(BeAMessageOfType(&messages.CommandGenerateSnippet{}))
				Expect(allMessagesSent[17]).To(BeAMessageOfType(&messages.TestStepFinished{}))
				Expect(allMessagesSent[18]).To(BeAMessageOfType(&messages.TestCaseFinished{}))
				Expect(allMessagesSent[19]).To(BeAMessageOfType(&messages.CommandRunAfterTestRunHooks{}))
				Expect(allMessagesSent[20]).To(Equal(&messages.Wrapper{
					Message: &messages.Wrapper_TestRunFinished{
						TestRunFinished: &messages.TestRunFinished{
							Success: false,
						},
					},
				}))
			})
		})
	})

	Context("all pickles gets rejected", func() {
		featurePath := path.Join(rootDir, "test", "fixtures", "a.feature")
		var allMessagesSent []*messages.Wrapper

		BeforeEach(func() {
			allMessagesSent = runWithConfigAndResponder(
				&messages.SourcesConfig{
					AbsolutePaths: []string{featurePath},
					Filters: &messages.SourcesFilterConfig{
						TagExpression: "@tagA",
					},
					Language: "en",
					Order:    &messages.SourcesOrder{},
				},
				&messages.RuntimeConfig{},
				&messages.SupportCodeConfig{},
				func(commandChan chan *messages.Wrapper, incoming *messages.Wrapper) {},
			)
		})

		It("does not run test run hooks", func() {
			Expect(allMessagesSent).To(HaveLen(6))
			Expect(allMessagesSent[0]).To(BeAMessageOfType(&messages.Source{}))
			Expect(allMessagesSent[1]).To(BeAMessageOfType(&messages.GherkinDocument{}))
			Expect(allMessagesSent[2]).To(BeAMessageOfType(&messages.Pickle{}))
			Expect(allMessagesSent[3]).To(BeAMessageOfType(&messages.PickleRejected{}))
			Expect(allMessagesSent[4]).To(BeAMessageOfType(&messages.TestRunStarted{}))
			Expect(allMessagesSent[5]).To(Equal(&messages.Wrapper{
				Message: &messages.Wrapper_TestRunFinished{
					TestRunFinished: &messages.TestRunFinished{
						Success: true,
					},
				},
			}))
		})
	})

	Context("with test case hooks", func() {
		featurePath := path.Join(rootDir, "test", "fixtures", "tags.feature")
		var allMessagesSent []*messages.Wrapper

		BeforeEach(func() {
			allMessagesSent = runWithConfigAndResponder(
				&messages.SourcesConfig{
					AbsolutePaths: []string{featurePath},
					Filters:       &messages.SourcesFilterConfig{},
					Language:      "en",
					Order:         &messages.SourcesOrder{},
				},
				&messages.RuntimeConfig{
					MaxParallel: 1,
				},
				&messages.SupportCodeConfig{
					BeforeTestCaseHookDefinitionConfigs: []*messages.TestCaseHookDefinitionConfig{
						{
							Id: "1",
							Location: &messages.SourceReference{
								Uri:      "hooks.js",
								Location: &messages.Location{Line: 11},
							},
						},
						{
							Id:            "2",
							TagExpression: "@tagA",
							Location: &messages.SourceReference{
								Uri:      "hooks.js",
								Location: &messages.Location{Line: 12},
							},
						},
					},
					AfterTestCaseHookDefinitionConfigs: []*messages.TestCaseHookDefinitionConfig{
						{
							Id:            "3",
							TagExpression: "@tagA",
							Location: &messages.SourceReference{
								Uri:      "hooks.js",
								Location: &messages.Location{Line: 13},
							},
						},
						{
							Id: "1",
							Location: &messages.SourceReference{
								Uri:      "hooks.js",
								Location: &messages.Location{Line: 14},
							},
						},
					},
				},
				func(commandChan chan *messages.Wrapper, incoming *messages.Wrapper) {
					switch x := incoming.Message.(type) {
					case *messages.Wrapper_CommandRunBeforeTestRunHooks:
						commandChan <- helpers.CreateActionCompleteMessage(x.CommandRunBeforeTestRunHooks.ActionId)
					case *messages.Wrapper_CommandRunAfterTestRunHooks:
						commandChan <- helpers.CreateActionCompleteMessage(x.CommandRunAfterTestRunHooks.ActionId)
					case *messages.Wrapper_CommandInitializeTestCase:
						commandChan <- helpers.CreateActionCompleteMessage(x.CommandInitializeTestCase.ActionId)
					case *messages.Wrapper_CommandRunBeforeTestCaseHook:
						commandChan <- helpers.CreateActionCompleteMessageWithTestResult(x.CommandRunBeforeTestCaseHook.ActionId, &messages.TestResult{Status: messages.Status_PASSED})
					case *messages.Wrapper_CommandRunAfterTestCaseHook:
						commandChan <- helpers.CreateActionCompleteMessageWithTestResult(x.CommandRunAfterTestCaseHook.ActionId, &messages.TestResult{Status: messages.Status_PASSED})
					case *messages.Wrapper_CommandGenerateSnippet:
						commandChan <- helpers.CreateActionCompleteMessageWithSnippet(x.CommandGenerateSnippet.ActionId, "snippet")
					}
				},
			)
		})

		It("runs test case hooks only for pickles that match the tag expression", func() {
			testCasePreparedMessages := []*messages.TestCasePrepared{}
			for _, msg := range allMessagesSent {
				if wrapper, ok := msg.Message.(*messages.Wrapper_TestCasePrepared); ok {
					testCasePreparedMessages = append(testCasePreparedMessages, wrapper.TestCasePrepared)
				}
			}
			Expect(testCasePreparedMessages).To(HaveLen(2))
			Expect(testCasePreparedMessages[0]).To(Equal(&messages.TestCasePrepared{
				PickleId: "A1",
				Steps: []*messages.TestCasePreparedStep{
					{
						ActionLocation: &messages.SourceReference{
							Uri:      "hooks.js",
							Location: &messages.Location{Line: 11},
						},
					},
					{
						SourceLocation: &messages.SourceReference{
							Uri:      featurePath,
							Location: &messages.Location{Line: 3, Column: 10},
						},
					},
					{
						ActionLocation: &messages.SourceReference{
							Uri:      "hooks.js",
							Location: &messages.Location{Line: 14},
						},
					},
				},
			}))
			Expect(testCasePreparedMessages[1]).To(Equal(&messages.TestCasePrepared{
				PickleId: "A2",
				Steps: []*messages.TestCasePreparedStep{
					{
						ActionLocation: &messages.SourceReference{
							Uri:      "hooks.js",
							Location: &messages.Location{Line: 11},
						},
					},
					{
						ActionLocation: &messages.SourceReference{
							Uri:      "hooks.js",
							Location: &messages.Location{Line: 12},
						},
					},
					{
						SourceLocation: &messages.SourceReference{
							Uri:      featurePath,
							Location: &messages.Location{Line: 7, Column: 10},
						},
					},
					{
						ActionLocation: &messages.SourceReference{
							Uri:      "hooks.js",
							Location: &messages.Location{Line: 13},
						},
					},
					{
						ActionLocation: &messages.SourceReference{
							Uri:      "hooks.js",
							Location: &messages.Location{Line: 14},
						},
					},
				},
			}))
		})
	})

	Context("running in parallel with three pickles", func() {
		featurePath := path.Join(rootDir, "test", "fixtures", "many.feature")
		var allMessagesSent []*messages.Wrapper

		Context("maxParallel is 2", func() {
			BeforeEach(func() {
				allMessagesSent = runWithConfigAndResponder(
					&messages.SourcesConfig{
						AbsolutePaths: []string{featurePath},
						Filters:       &messages.SourcesFilterConfig{},
						Language:      "en",
						Order:         &messages.SourcesOrder{},
					},
					&messages.RuntimeConfig{
						MaxParallel: 2,
					},
					&messages.SupportCodeConfig{},
					func(commandChan chan *messages.Wrapper, incoming *messages.Wrapper) {
						switch x := incoming.Message.(type) {
						case *messages.Wrapper_CommandRunBeforeTestRunHooks:
							commandChan <- helpers.CreateActionCompleteMessage(x.CommandRunBeforeTestRunHooks.ActionId)
						case *messages.Wrapper_CommandRunAfterTestRunHooks:
							commandChan <- helpers.CreateActionCompleteMessage(x.CommandRunAfterTestRunHooks.ActionId)
						case *messages.Wrapper_CommandInitializeTestCase:
							commandChan <- helpers.CreateActionCompleteMessage(x.CommandInitializeTestCase.ActionId)
						case *messages.Wrapper_CommandGenerateSnippet:
							time.Sleep(100 * time.Millisecond)
							commandChan <- helpers.CreateActionCompleteMessageWithSnippet(x.CommandGenerateSnippet.ActionId, "snippet")
						}
					},
				)
			})

			It("runs exactly two test cases at once", func() {
				maxRunning := determineMaxRunning(allMessagesSent)
				Expect(maxRunning).To(Equal(2))
			})
		})

		Context("maxParallel is 0", func() {
			BeforeEach(func() {
				allMessagesSent = runWithConfigAndResponder(
					&messages.SourcesConfig{
						AbsolutePaths: []string{featurePath},
						Filters:       &messages.SourcesFilterConfig{},
						Language:      "en",
						Order:         &messages.SourcesOrder{},
					},
					&messages.RuntimeConfig{
						MaxParallel: 0,
					},
					&messages.SupportCodeConfig{},
					func(commandChan chan *messages.Wrapper, incoming *messages.Wrapper) {
						switch x := incoming.Message.(type) {
						case *messages.Wrapper_CommandRunBeforeTestRunHooks:
							commandChan <- helpers.CreateActionCompleteMessage(x.CommandRunBeforeTestRunHooks.ActionId)
						case *messages.Wrapper_CommandRunAfterTestRunHooks:
							commandChan <- helpers.CreateActionCompleteMessage(x.CommandRunAfterTestRunHooks.ActionId)
						case *messages.Wrapper_CommandInitializeTestCase:
							commandChan <- helpers.CreateActionCompleteMessage(x.CommandInitializeTestCase.ActionId)
						case *messages.Wrapper_CommandGenerateSnippet:
							time.Sleep(100 * time.Millisecond)
							commandChan <- helpers.CreateActionCompleteMessageWithSnippet(x.CommandGenerateSnippet.ActionId, "snippet")
						}
					},
				)
			})

			It("runs all test cases at once", func() {
				maxRunning := determineMaxRunning(allMessagesSent)
				Expect(maxRunning).To(Equal(5))
			})
		})
	})
})

func determineMaxRunning(allMessagesSent []*messages.Wrapper) int {
	maxRunning := 0
	currentRunning := 0
	for _, msg := range allMessagesSent {
		if _, ok := msg.Message.(*messages.Wrapper_CommandInitializeTestCase); ok {
			currentRunning++
			if currentRunning > maxRunning {
				maxRunning = currentRunning
			}
		}
		if _, ok := msg.Message.(*messages.Wrapper_TestCaseFinished); ok {
			currentRunning--
		}
	}
	return maxRunning
}

func runWithConfigAndResponder(sourcesConfig *messages.SourcesConfig, runtimeConfig *messages.RuntimeConfig, supportCodeConfig *messages.SupportCodeConfig, responder func(chan *messages.Wrapper, *messages.Wrapper)) []*messages.Wrapper {
	allMessagesSent := []*messages.Wrapper{}
	r := runner.NewRunner()
	incoming, outgoing := r.GetCommandChannels()
	done := make(chan bool)
	go func() {
		for msg := range outgoing {
			allMessagesSent = append(allMessagesSent, msg)
			responder(incoming, msg)
		}
		done <- true
	}()
	incoming <- &messages.Wrapper{
		Message: &messages.Wrapper_CommandStart{
			CommandStart: &messages.CommandStart{
				SourcesConfig:     sourcesConfig,
				RuntimeConfig:     runtimeConfig,
				SupportCodeConfig: supportCodeConfig,
			},
		},
	}
	<-done
	return allMessagesSent
}
