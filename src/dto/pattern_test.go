package dto_test

import (
	"github.com/cucumber/cucumber-engine/src/dto"
	"github.com/cucumber/cucumber-expressions-go"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Pattern", func() {
	Describe("Expression", func() {
		Context("type is cucumber expression", func() {
			It("returns a cucumber expression", func() {
				pattern := dto.Pattern{
					Type:   dto.PatternTypeCucumberExpression,
					Source: "I have {int} cukes",
				}
				parameterTypeRegistry := cucumberexpressions.NewParameterTypeRegistry()
				expr, err := pattern.Expression(parameterTypeRegistry)
				Expect(err).NotTo(HaveOccurred())
				Expect(expr).To(BeAssignableToTypeOf(&cucumberexpressions.CucumberExpression{}))
			})
		})

		Context("type is regular expression", func() {
			It("returns a regular expression", func() {
				pattern := dto.Pattern{
					Type:   dto.PatternTypeRegularExpression,
					Source: `^I have (\d+) cukes$`,
				}
				parameterTypeRegistry := cucumberexpressions.NewParameterTypeRegistry()
				expr, err := pattern.Expression(parameterTypeRegistry)
				Expect(err).NotTo(HaveOccurred())
				Expect(expr).To(BeAssignableToTypeOf(&cucumberexpressions.RegularExpression{}))
			})
		})

		Context("type is invalid", func() {
			It("returns an error", func() {
				pattern := dto.Pattern{
					Type:   dto.PatternType("invalid"),
					Source: "I have {int} cukes",
				}
				parameterTypeRegistry := cucumberexpressions.NewParameterTypeRegistry()
				_, err := pattern.Expression(parameterTypeRegistry)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("Unexpected pattern type: `invalid`. Should be `cucumber_expression` or `regular_expression`"))
			})
		})
	})
})
