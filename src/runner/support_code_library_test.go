package runner_test

import (
	"github.com/cucumber/cucumber-pickle-runner/src/dto"
	"github.com/cucumber/cucumber-pickle-runner/src/runner"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("SupportCodeLibrary", func() {
	var library *runner.SupportCodeLibrary

	It("returns err on invalid after hook tag expression", func() {
		_, err := runner.NewSupportCodeLibrary(&dto.SupportCodeConfig{
			AfterTestCaseHookDefinitionConfigs: []*dto.TestCaseHookDefinitionConfig{
				{TagExpression: "@tagA @tagB"},
			},
		})
		Expect(err).To(HaveOccurred())
	})

	It("returns err on invalid before hook tag expression", func() {
		_, err := runner.NewSupportCodeLibrary(&dto.SupportCodeConfig{
			BeforeTestCaseHookDefinitionConfigs: []*dto.TestCaseHookDefinitionConfig{
				{TagExpression: "@tagA @tagB"},
			},
		})
		Expect(err).To(HaveOccurred())
	})

	It("returns err on invalid step pattern", func() {
		_, err := runner.NewSupportCodeLibrary(&dto.SupportCodeConfig{
			StepDefinitionConfigs: []*dto.StepDefinitionConfig{
				{Pattern: dto.Pattern{Type: "regular_expression", Source: "*"}},
			},
		})
		Expect(err).To(HaveOccurred())
	})

	It("returns err on parameter type issues", func() {
		_, err := runner.NewSupportCodeLibrary(&dto.SupportCodeConfig{
			ParameterTypeConfigs: []*dto.ParameterTypeConfig{
				{Name: "parameterType1"},
				{Name: "parameterType1"},
			},
		})
		Expect(err).To(HaveOccurred())
	})

	Describe("GetMatchingAfterTestCaseHookDefinitions", func() {
		Context("with an empty config", func() {
			BeforeEach(func() {
				var err error
				library, err = runner.NewSupportCodeLibrary(&dto.SupportCodeConfig{})
				Expect(err).NotTo(HaveOccurred())
			})

			It("returns no matching hook definitions", func() {
				Expect(library.GetMatchingAfterTestCaseHookDefinitions([]string{})).To(BeEmpty())
			})
		})

		Context("with a hook", func() {
			BeforeEach(func() {
				var err error
				library, err = runner.NewSupportCodeLibrary(&dto.SupportCodeConfig{
					AfterTestCaseHookDefinitionConfigs: []*dto.TestCaseHookDefinitionConfig{
						{ID: "afterHook1", TagExpression: "@tagA"},
					},
				})
				Expect(err).NotTo(HaveOccurred())
			})

			It("returns the hook only when the tag expression returns true", func() {
				matching := library.GetMatchingAfterTestCaseHookDefinitions([]string{})
				Expect(matching).To(BeEmpty())
				matching = library.GetMatchingAfterTestCaseHookDefinitions([]string{"@tagA"})
				Expect(matching).To(HaveLen(1))
				Expect(matching[0].ID).To(Equal("afterHook1"))
			})
		})
	})

	Describe("GetMatchingBeforeTestCaseHookDefinitions", func() {
		Context("with an empty config", func() {
			BeforeEach(func() {
				var err error
				library, err = runner.NewSupportCodeLibrary(&dto.SupportCodeConfig{})
				Expect(err).NotTo(HaveOccurred())
			})

			It("returns no matching hook definitions", func() {
				Expect(library.GetMatchingBeforeTestCaseHookDefinitions([]string{})).To(BeEmpty())
			})
		})

		Context("with a hook", func() {
			BeforeEach(func() {
				var err error
				library, err = runner.NewSupportCodeLibrary(&dto.SupportCodeConfig{
					BeforeTestCaseHookDefinitionConfigs: []*dto.TestCaseHookDefinitionConfig{
						{ID: "beforeHook1", TagExpression: "@tagB"},
					},
				})
				Expect(err).NotTo(HaveOccurred())
			})

			It("returns the hook only when the tag expression returns true", func() {
				matching := library.GetMatchingBeforeTestCaseHookDefinitions([]string{})
				Expect(matching).To(BeEmpty())
				matching = library.GetMatchingBeforeTestCaseHookDefinitions([]string{"@tagB"})
				Expect(matching).To(HaveLen(1))
				Expect(matching[0].ID).To(Equal("beforeHook1"))
			})
		})
	})

	Describe("GetMatchingStepDefinitions", func() {
		Context("with an empty config", func() {
			BeforeEach(func() {
				var err error
				library, err = runner.NewSupportCodeLibrary(&dto.SupportCodeConfig{})
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
				library, err = runner.NewSupportCodeLibrary(&dto.SupportCodeConfig{
					StepDefinitionConfigs: []*dto.StepDefinitionConfig{
						{ID: "step1", Pattern: dto.Pattern{Type: "cucumber_expression", Source: "I have {int} cukes"}},
					},
				})
				Expect(err).NotTo(HaveOccurred())
			})

			It("returns the step and the pattern matches", func() {
				stepDefinitions, patternMatch, err := library.GetMatchingStepDefinitions("I have 10 cukes")
				Expect(err).NotTo(HaveOccurred())
				Expect(stepDefinitions).To(HaveLen(1))
				Expect(stepDefinitions[0].ID).To(Equal("step1"))
				Expect(patternMatch).To(Equal([]*dto.PatternMatch{
					{Captures: []string{"10"}, ParameterTypeName: "int"},
				}))
			})
		})

		Context("with multiple steps that match", func() {
			BeforeEach(func() {
				var err error
				library, err = runner.NewSupportCodeLibrary(&dto.SupportCodeConfig{
					StepDefinitionConfigs: []*dto.StepDefinitionConfig{
						{ID: "step1", Pattern: dto.Pattern{Type: "cucumber_expression", Source: "I have {int} cukes"}},
						{ID: "step2", Pattern: dto.Pattern{Type: "regular_expression", Source: `^I have (\d+) cukes$`}},
					},
				})
				Expect(err).NotTo(HaveOccurred())
			})

			It("returns the steps and no pattern matches", func() {
				stepDefinitions, patternMatch, err := library.GetMatchingStepDefinitions("I have 10 cukes")
				Expect(err).NotTo(HaveOccurred())
				Expect(stepDefinitions).To(HaveLen(2))
				Expect(stepDefinitions[0].ID).To(Equal("step1"))
				Expect(stepDefinitions[1].ID).To(Equal("step2"))
				Expect(patternMatch).To(BeNil())
			})
		})

		Context("with an error while matching the step", func() {
			BeforeEach(func() {
				var err error
				library, err = runner.NewSupportCodeLibrary(&dto.SupportCodeConfig{
					ParameterTypeConfigs: []*dto.ParameterTypeConfig{
						{Name: "name1", Regexps: []string{`abc`}},
						{Name: "name2", Regexps: []string{`abc`}},
					},
					StepDefinitionConfigs: []*dto.StepDefinitionConfig{
						{ID: "step1", Pattern: dto.Pattern{Type: "regular_expression", Source: `^I have (abc)$`}},
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
