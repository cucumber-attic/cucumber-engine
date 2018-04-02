package dto

// RuntimeConfig is the configuration for the run
type RuntimeConfig struct {
	IsFailFast                    bool                      `json:"is_fail_fast"`
	IsDryRun                      bool                      `json:"is_dry_run"`
	IsStrict                      bool                      `json:"is_strict"`
	BeforeTestCaseHookDefinitions []*TestCaseHookDefinition `json:"before_test_case_hook_definitions"`
	AfterTestCaseHookDefinitions  []*TestCaseHookDefinition `json:"after_test_case_hook_definitions"`
	StepDefinitions               []*StepDefinition         `json:"step_definitions"`
	ParameterTypes                []*ParameterType          `json:"parameter_types"`
}
