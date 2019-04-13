package runner_test

import (
	"github.com/cucumber/cucumber-engine/src/runner"
	messages "github.com/cucumber/cucumber-messages-go/v2"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("SupportCodeLibrary", func() {
	var library *runner.SupportCodeLibrary

	It("returns err on invalid after hook tag expression", func() {
		_, err := runner.NewSupportCodeLibrary(&messages.SupportCodeConfig{
			AfterTestCaseHookDefinitionConfigs: []*messages.TestCaseHookDefinitionConfig{
				{TagExpression: "@tagA @tagB"},
			},
		})
		Expect(err).To(HaveOccurred())
	})

	It("returns err on invalid before hook tag expression", func() {
		_, err := runner.NewSupportCodeLibrary(&messages.SupportCodeConfig{
			BeforeTestCaseHookDefinitionConfigs: []*messages.TestCaseHookDefinitionConfig{
				{TagExpression: "@tagA @tagB"},
			},
		})
		Expect(err).To(HaveOccurred())
	})

	It("returns err on invalid step pattern", func() {
		_, err := runner.NewSupportCodeLibrary(&messages.SupportCodeConfig{
			StepDefinitionConfigs: []*messages.StepDefinitionConfig{
				{
					Pattern: &messages.StepDefinitionPattern{
						Type:   messages.StepDefinitionPatternType_REGULAR_EXPRESSION,
						Source: "*",
					},
				},
			},
		})
		Expect(err).To(HaveOccurred())
	})

	It("returns err on parameter type issues", func() {
		_, err := runner.NewSupportCodeLibrary(&messages.SupportCodeConfig{
			ParameterTypeConfigs: []*messages.ParameterTypeConfig{
				{Name: "parameterType1"},
				{Name: "parameterType1"},
			},
		})
		Expect(err).To(HaveOccurred())
	})

	Describe("GenerateExpressions", func() {
		It("returns the generated expressions", func() {
			var err error
			library, err = runner.NewSupportCodeLibrary(&messages.SupportCodeConfig{})
			Expect(err).NotTo(HaveOccurred())
			Expect(library.GenerateExpressions(`I have 100 cukes`)).To(Equal([]*messages.GeneratedExpression{
				{
					Text:               "I have {int} cukes",
					ParameterTypeNames: []string{"int"},
				},
			}))
		})
	})

	Describe("GetMatchingAfterTestCaseHookDefinitions", func() {
		Context("with an empty config", func() {
			BeforeEach(func() {
				var err error
				library, err = runner.NewSupportCodeLibrary(&messages.SupportCodeConfig{})
				Expect(err).NotTo(HaveOccurred())
			})

			It("returns no matching hook definitions", func() {
				Expect(library.GetMatchingAfterTestCaseHookDefinitions([]string{})).To(BeEmpty())
			})
		})

		Context("with a hook", func() {
			BeforeEach(func() {
				var err error
				library, err = runner.NewSupportCodeLibrary(&messages.SupportCodeConfig{
					AfterTestCaseHookDefinitionConfigs: []*messages.TestCaseHookDefinitionConfig{
						{Id: "afterHook1", TagExpression: "@tagA"},
					},
				})
				Expect(err).NotTo(HaveOccurred())
			})

			It("returns the hook only when the tag expression returns true", func() {
				matching := library.GetMatchingAfterTestCaseHookDefinitions([]string{})
				Expect(matching).To(BeEmpty())
				matching = library.GetMatchingAfterTestCaseHookDefinitions([]string{"@tagA"})
				Expect(matching).To(HaveLen(1))
				Expect(matching[0].Config.Id).To(Equal("afterHook1"))
			})
		})
	})

	Describe("GetMatchingBeforeTestCaseHookDefinitions", func() {
		Context("with an empty config", func() {
			BeforeEach(func() {
				var err error
				library, err = runner.NewSupportCodeLibrary(&messages.SupportCodeConfig{})
				Expect(err).NotTo(HaveOccurred())
			})

			It("returns no matching hook definitions", func() {
				Expect(library.GetMatchingBeforeTestCaseHookDefinitions([]string{})).To(BeEmpty())
			})
		})

		Context("with a hook", func() {
			BeforeEach(func() {
				var err error
				library, err = runner.NewSupportCodeLibrary(&messages.SupportCodeConfig{
					BeforeTestCaseHookDefinitionConfigs: []*messages.TestCaseHookDefinitionConfig{
						{Id: "beforeHook1", TagExpression: "@tagB"},
					},
				})
				Expect(err).NotTo(HaveOccurred())
			})

			It("returns the hook only when the tag expression returns true", func() {
				matching := library.GetMatchingBeforeTestCaseHookDefinitions([]string{})
				Expect(matching).To(BeEmpty())
				matching = library.GetMatchingBeforeTestCaseHookDefinitions([]string{"@tagB"})
				Expect(matching).To(HaveLen(1))
				Expect(matching[0].Config.Id).To(Equal("beforeHook1"))
			})
		})
	})

	Describe("GetMatchingStepDefinitions", func() {
		Context("with an empty config", func() {
			BeforeEach(func() {
				var err error
				library, err = runner.NewSupportCodeLibrary(&messages.SupportCodeConfig{})
				Expect(err).NotTo(HaveOccurred())
			})

			It("returns no matching step definitions", func() {
				stepDefinitions, patternMatch, err := library.GetMatchingStepDefinitions("a step")
				Expect(err).NotTo(HaveOccurred())
				Expect(stepDefinitions).To(BeEmpty())
				Expect(patternMatch).To(BeEmpty())
			})
		})

		Context("with a step", func() {
			BeforeEach(func() {
				var err error
				library, err = runner.NewSupportCodeLibrary(&messages.SupportCodeConfig{
					StepDefinitionConfigs: []*messages.StepDefinitionConfig{
						{
							Id: "step1",
							Pattern: &messages.StepDefinitionPattern{
								Type:   messages.StepDefinitionPatternType_CUCUMBER_EXPRESSION,
								Source: "I have {int} cukes",
							},
						},
					},
				})
				Expect(err).NotTo(HaveOccurred())
			})

			It("returns the step and the pattern matches", func() {
				stepDefinitions, patternMatch, err := library.GetMatchingStepDefinitions("I have 10 cukes")
				Expect(err).NotTo(HaveOccurred())
				Expect(stepDefinitions).To(HaveLen(1))
				Expect(stepDefinitions[0].Config.Id).To(Equal("step1"))
				Expect(patternMatch).To(Equal([]*messages.PatternMatch{
					{Captures: []string{"10"}, ParameterTypeName: "int"},
				}))
			})
		})

		Context("with multiple steps that match", func() {
			BeforeEach(func() {
				var err error
				library, err = runner.NewSupportCodeLibrary(&messages.SupportCodeConfig{
					StepDefinitionConfigs: []*messages.StepDefinitionConfig{
						{
							Id: "step1",
							Pattern: &messages.StepDefinitionPattern{
								Type:   messages.StepDefinitionPatternType_CUCUMBER_EXPRESSION,
								Source: "I have {int} cukes",
							},
						},
						{
							Id: "step2",
							Pattern: &messages.StepDefinitionPattern{
								Type:   messages.StepDefinitionPatternType_REGULAR_EXPRESSION,
								Source: `^I have (\d+) cukes$`,
							},
						},
					},
				})
				Expect(err).NotTo(HaveOccurred())
			})

			It("returns the steps and no pattern matches", func() {
				stepDefinitions, patternMatch, err := library.GetMatchingStepDefinitions("I have 10 cukes")
				Expect(err).NotTo(HaveOccurred())
				Expect(stepDefinitions).To(HaveLen(2))
				Expect(stepDefinitions[0].Config.Id).To(Equal("step1"))
				Expect(stepDefinitions[1].Config.Id).To(Equal("step2"))
				Expect(patternMatch).To(BeNil())
			})
		})

		Context("with an error while matching the step", func() {
			BeforeEach(func() {
				var err error
				library, err = runner.NewSupportCodeLibrary(&messages.SupportCodeConfig{
					ParameterTypeConfigs: []*messages.ParameterTypeConfig{
						{Name: "name1", RegularExpressions: []string{`abc`}},
						{Name: "name2", RegularExpressions: []string{`abc`}},
					},
					StepDefinitionConfigs: []*messages.StepDefinitionConfig{
						{
							Id: "step1",
							Pattern: &messages.StepDefinitionPattern{
								Type:   messages.StepDefinitionPatternType_REGULAR_EXPRESSION,
								Source: `^I have (abc)$`,
							},
						},
					},
				})
				Expect(err).NotTo(HaveOccurred())
			})

			It("returns the error", func() {
				_, _, err := library.GetMatchingStepDefinitions("I have abc")
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
