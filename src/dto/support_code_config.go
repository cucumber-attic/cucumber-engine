package dto

// SupportCodeConfig is the configuration for the support code
type SupportCodeConfig struct {
	BeforeTestCaseHookDefinitionConfigs []*TestCaseHookDefinitionConfig `json:"beforeTestCaseHookDefinitions"`
	AfterTestCaseHookDefinitionConfigs  []*TestCaseHookDefinitionConfig `json:"afterTestCaseHookDefinitions"`
	StepDefinitionConfigs               []*StepDefinitionConfig         `json:"stepDefinitions"`
	ParameterTypeConfigs                []*ParameterTypeConfig          `json:"parameterTypes"`
}
