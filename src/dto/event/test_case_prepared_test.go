package event_test

import (
	"github.com/cucumber/cucumber-engine/src/dto"
	"github.com/cucumber/cucumber-engine/src/dto/event"
	messages "github.com/cucumber/cucumber-messages-go/v2"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("NewTestCasePrepared", func() {
	Context("with no steps/hooks", func() {
		It("the event has no steps", func() {
			testCasePrepared := event.NewTestCasePrepared(event.NewTestCasePreparedOptions{
				Pickle: &messages.Pickle{Uri: "a.feature"},
			})
			Expect(testCasePrepared.Steps).To(BeEmpty())
		})
	})

	Context("with a step with no definitions", func() {
		It("the event does not have an actionLocation for the step", func() {
			testCasePrepared := event.NewTestCasePrepared(event.NewTestCasePreparedOptions{
				Pickle: &messages.Pickle{
					Locations: []*messages.Location{{Line: 1}},
					Steps: []*messages.PickleStep{
						{Locations: []*messages.Location{{Line: 2}}},
					},
					Uri: "a.feature",
				},
				StepIndexToStepDefinitions: [][]*dto.StepDefinition{
					{},
				},
			})
			Expect(testCasePrepared.Steps).To(Equal([]*messages.TestCasePreparedStep{
				{
					SourceLocation: &messages.SourceReference{
						Uri:      "a.feature",
						Location: &messages.Location{Line: 2},
					},
				},
			}))
		})
	})

	Context("with step with one definition", func() {
		It("the event has an actionLocation for the step", func() {
			testCasePrepared := event.NewTestCasePrepared(event.NewTestCasePreparedOptions{
				Pickle: &messages.Pickle{
					Locations: []*messages.Location{{Line: 1}},
					Steps: []*messages.PickleStep{
						{Locations: []*messages.Location{{Line: 2}}},
					},
					Uri: "a.feature",
				},
				StepIndexToStepDefinitions: [][]*dto.StepDefinition{
					{
						{
							Config: &messages.StepDefinitionConfig{
								Location: &messages.SourceReference{
									Uri:      "steps.js",
									Location: &messages.Location{Line: 3},
								},
							},
						},
					},
				},
			})
			Expect(testCasePrepared.Steps).To(Equal([]*messages.TestCasePreparedStep{
				{
					SourceLocation: &messages.SourceReference{
						Uri:      "a.feature",
						Location: &messages.Location{Line: 2},
					},
					ActionLocation: &messages.SourceReference{
						Uri:      "steps.js",
						Location: &messages.Location{Line: 3},
					},
				},
			}))
		})
	})

	Context("with step with multiple definitions", func() {
		It("the event does not have an actionLocation for the step", func() {
			testCasePrepared := event.NewTestCasePrepared(event.NewTestCasePreparedOptions{
				Pickle: &messages.Pickle{
					Locations: []*messages.Location{{Line: 1}},
					Steps: []*messages.PickleStep{
						{Locations: []*messages.Location{{Line: 2}}},
					},
					Uri: "a.feature",
				},
				StepIndexToStepDefinitions: [][]*dto.StepDefinition{
					{
						{
							Config: &messages.StepDefinitionConfig{
								Location: &messages.SourceReference{
									Uri:      "steps.js",
									Location: &messages.Location{Line: 3},
								},
							},
						},
						{
							Config: &messages.StepDefinitionConfig{
								Location: &messages.SourceReference{
									Uri:      "steps.js",
									Location: &messages.Location{Line: 4},
								},
							},
						},
					},
				},
			})
			Expect(testCasePrepared.Steps).To(Equal([]*messages.TestCasePreparedStep{
				{
					SourceLocation: &messages.SourceReference{
						Uri:      "a.feature",
						Location: &messages.Location{Line: 2},
					},
				},
			}))
		})
	})

	Context("with step and hooks", func() {
		It("the event has both listed as steps", func() {
			testCasePrepared := event.NewTestCasePrepared(event.NewTestCasePreparedOptions{
				Pickle: &messages.Pickle{
					Locations: []*messages.Location{{Line: 1}},
					Steps: []*messages.PickleStep{
						{Locations: []*messages.Location{{Line: 2}}},
					},
					Uri: "a.feature",
				},
				BeforeTestCaseHookDefinitions: []*dto.TestCaseHookDefinition{
					{
						Config: &messages.TestCaseHookDefinitionConfig{
							Location: &messages.SourceReference{
								Uri:      "steps.js",
								Location: &messages.Location{Line: 10},
							},
						},
					},
				},
				AfterTestCaseHookDefinitions: []*dto.TestCaseHookDefinition{
					{
						Config: &messages.TestCaseHookDefinitionConfig{
							Location: &messages.SourceReference{
								Uri:      "steps.js",
								Location: &messages.Location{Line: 11},
							},
						},
					},
				},
				StepIndexToStepDefinitions: [][]*dto.StepDefinition{
					{
						{
							Config: &messages.StepDefinitionConfig{
								Location: &messages.SourceReference{
									Uri:      "steps.js",
									Location: &messages.Location{Line: 3},
								},
							},
						},
					},
				},
			})
			Expect(testCasePrepared.Steps).To(Equal([]*messages.TestCasePreparedStep{
				{
					ActionLocation: &messages.SourceReference{
						Uri:      "steps.js",
						Location: &messages.Location{Line: 10},
					},
				},
				{
					ActionLocation: &messages.SourceReference{
						Uri:      "steps.js",
						Location: &messages.Location{Line: 3},
					},
					SourceLocation: &messages.SourceReference{
						Uri:      "a.feature",
						Location: &messages.Location{Line: 2},
					},
				},
				{
					ActionLocation: &messages.SourceReference{
						Uri:      "steps.js",
						Location: &messages.Location{Line: 11},
					},
				},
			}))
		})
	})
})
