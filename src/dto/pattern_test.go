package dto_test

import (
	"github.com/cucumber/cucumber-engine/src/dto"
	cucumberexpressions "github.com/cucumber/cucumber-expressions-go"
	messages "github.com/cucumber/cucumber-messages-go/v3"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Pattern", func() {
	Describe("Expression", func() {
		Context("type is cucumber expression", func() {
			It("returns a cucumber expression", func() {
				pattern := &messages.StepDefinitionPattern{
					Type:   messages.StepDefinitionPatternType_CUCUMBER_EXPRESSION,
					Source: "I have {int} cukes",
				}
				parameterTypeRegistry := cucumberexpressions.NewParameterTypeRegistry()
				expr, err := dto.GetExpression(pattern, parameterTypeRegistry)
				Expect(err).NotTo(HaveOccurred())
				Expect(expr).To(BeAssignableToTypeOf(&cucumberexpressions.CucumberExpression{}))
			})
		})

		Context("type is regular expression", func() {
			It("returns a regular expression", func() {
				pattern := &messages.StepDefinitionPattern{
					Type:   messages.StepDefinitionPatternType_REGULAR_EXPRESSION,
					Source: `^I have (\d+) cukes$`,
				}
				parameterTypeRegistry := cucumberexpressions.NewParameterTypeRegistry()
				expr, err := dto.GetExpression(pattern, parameterTypeRegistry)
				Expect(err).NotTo(HaveOccurred())
				Expect(expr).To(BeAssignableToTypeOf(&cucumberexpressions.RegularExpression{}))
			})
		})

		Context("type is invalid", func() {
			It("returns an error", func() {
				pattern := &messages.StepDefinitionPattern{
					Type:   3,
					Source: "I have {int} cukes",
				}
				parameterTypeRegistry := cucumberexpressions.NewParameterTypeRegistry()
				_, err := dto.GetExpression(pattern, parameterTypeRegistry)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("Unexpected pattern type: `3`. Should be `CUCUMBER_EXPRESSION` or `REGULAR_EXPRESSION`"))
			})
		})
	})
})
