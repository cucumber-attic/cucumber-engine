package event_test

import (
	"github.com/cucumber/cucumber-pickle-runner/src/dto"
	"github.com/cucumber/cucumber-pickle-runner/src/dto/event"
	gherkin "github.com/cucumber/gherkin-go"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("NewTestCasePrepared", func() {
	Context("with no steps/hooks", func() {
		It("the event has no steps", func() {
			testCasePrepared := event.NewTestCasePrepared(event.NewTestCasePreparedOptions{
				Pickle: &gherkin.Pickle{
					Locations: []gherkin.Location{{Line: 1}},
				},
				URI: "a.feature",
			})
			Expect(testCasePrepared.SourceLocation).To(Equal(event.Location{URI: "a.feature", Line: 1}))
			Expect(testCasePrepared.Steps).To(BeEmpty())
		})
	})

	Context("with a step with no definitions", func() {
		It("the event does not have an actionLocation for the step", func() {
			testCasePrepared := event.NewTestCasePrepared(event.NewTestCasePreparedOptions{
				Pickle: &gherkin.Pickle{
					Locations: []gherkin.Location{{Line: 1}},
					Steps: []*gherkin.PickleStep{
						{Locations: []gherkin.Location{{Line: 2}}},
					},
				},
				StepIndexToStepDefinitions: [][]*dto.StepDefinition{
					{},
				},
				URI: "a.feature",
			})
			Expect(testCasePrepared.SourceLocation).To(Equal(event.Location{URI: "a.feature", Line: 1}))
			Expect(testCasePrepared.Steps).To(Equal([]event.TestCasePreparedStep{
				{SourceLocation: event.Location{URI: "a.feature", Line: 2}},
			}))
		})
	})

	Context("with step with one definition", func() {
		It("the event has an actionLocation for the step", func() {
			testCasePrepared := event.NewTestCasePrepared(event.NewTestCasePreparedOptions{
				Pickle: &gherkin.Pickle{
					Locations: []gherkin.Location{{Line: 1}},
					Steps: []*gherkin.PickleStep{
						{Locations: []gherkin.Location{{Line: 2}}},
					},
				},
				StepIndexToStepDefinitions: [][]*dto.StepDefinition{
					{
						{URI: "steps.js", Line: 3},
					},
				},
				URI: "a.feature",
			})
			Expect(testCasePrepared.SourceLocation).To(Equal(event.Location{URI: "a.feature", Line: 1}))
			Expect(testCasePrepared.Steps).To(Equal([]event.TestCasePreparedStep{
				{
					SourceLocation: event.Location{URI: "a.feature", Line: 2},
					ActionLocation: event.Location{URI: "steps.js", Line: 3},
				},
			}))
		})
	})

	Context("with step with multiple definitions", func() {
		It("the event does not have an actionLocation for the step", func() {
			testCasePrepared := event.NewTestCasePrepared(event.NewTestCasePreparedOptions{
				Pickle: &gherkin.Pickle{
					Locations: []gherkin.Location{{Line: 1}},
					Steps: []*gherkin.PickleStep{
						{Locations: []gherkin.Location{{Line: 2}}},
					},
				},
				StepIndexToStepDefinitions: [][]*dto.StepDefinition{
					{
						{URI: "steps.js", Line: 3},
						{URI: "steps.js", Line: 4},
					},
				},
				URI: "a.feature",
			})
			Expect(testCasePrepared.SourceLocation).To(Equal(event.Location{URI: "a.feature", Line: 1}))
			Expect(testCasePrepared.Steps).To(Equal([]event.TestCasePreparedStep{
				{
					SourceLocation: event.Location{URI: "a.feature", Line: 2},
				},
			}))
		})
	})

	Context("with step and hooks", func() {
		It("the event has both listed as steps", func() {
			testCasePrepared := event.NewTestCasePrepared(event.NewTestCasePreparedOptions{
				Pickle: &gherkin.Pickle{
					Locations: []gherkin.Location{{Line: 1}},
					Steps: []*gherkin.PickleStep{
						{Locations: []gherkin.Location{{Line: 2}}},
					},
				},
				BeforeTestCaseHookDefinitions: []*dto.TestCaseHookDefinition{
					{URI: "steps.js", Line: 10},
				},
				AfterTestCaseHookDefinitions: []*dto.TestCaseHookDefinition{
					{URI: "steps.js", Line: 11},
				},
				StepIndexToStepDefinitions: [][]*dto.StepDefinition{
					{{URI: "steps.js", Line: 3}},
				},
				URI: "a.feature",
			})
			Expect(testCasePrepared.SourceLocation).To(Equal(event.Location{URI: "a.feature", Line: 1}))
			Expect(testCasePrepared.Steps).To(Equal([]event.TestCasePreparedStep{
				{
					ActionLocation: event.Location{URI: "steps.js", Line: 10},
				},
				{
					ActionLocation: event.Location{URI: "steps.js", Line: 3},
					SourceLocation: event.Location{URI: "a.feature", Line: 2},
				},
				{
					ActionLocation: event.Location{URI: "steps.js", Line: 11},
				},
			}))
		})
	})
})
