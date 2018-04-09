package dto

// CommandType is an enumeration of the available values for the Type field
type CommandType string

const (
	// CommandTypeStart is received to start the process
	CommandTypeStart = CommandType("start")
	// CommandTypeActionComplete is received when the caller has completed an action
	CommandTypeActionComplete = CommandType("action_complete")

	// CommandTypeRunBeforeTestRunHooks is sent to have the caller run before test run hooks
	CommandTypeRunBeforeTestRunHooks = CommandType("run_before_test_run_hooks")
	// CommandTypeInitializeTestCase is sent to have the caller initialize a test case
	CommandTypeInitializeTestCase = CommandType("run_initialize_test_case")
	// CommandTypeRunBeforeTestCaseHook is sent to have the caller run a before test case hook
	CommandTypeRunBeforeTestCaseHook = CommandType("run_before_test_case_hook")
	// CommandTypeRunTestStep is sent to have the caller run a step
	CommandTypeRunTestStep = CommandType("run_test_step")
	// CommandTypeRunAfterTestCaseHook is sent to have the caller run a after test case hook
	CommandTypeRunAfterTestCaseHook = CommandType("run_after_test_case_hook")
	// CommandTypeRunAfterTestRunHooks is sent to have the caller run after test run hooks
	CommandTypeRunAfterTestRunHooks = CommandType("run_after_test_run_hooks")
	// CommandTypeEvent is sent when an event occurs
	CommandTypeEvent = CommandType("event")
	// CommandTypeGenerateSnippet is sent to have the caller generate a step snippet for a pattern
	CommandTypeGenerateSnippet = CommandType("generate_snippet")
	// CommandTypeError is sent when an error occurs
	CommandTypeError = CommandType("error")
)
