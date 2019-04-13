package dto

import messages "github.com/cucumber/cucumber-messages-go/v2"

// TestRunResult is the result of a test run
type TestRunResult struct {
	Success bool `json:"success"`
}

// NewTestRunResult creates a new test run result
func NewTestRunResult() *TestRunResult {
	return &TestRunResult{
		Success: true,
	}
}

// Update updates the test run result with a test case result
func (t *TestRunResult) Update(testCaseResult *messages.TestResult, isStrict bool) {
	if shouldCauseFailure(testCaseResult.Status, isStrict) {
		t.Success = false
	}
}

func shouldCauseFailure(status messages.Status, isStrict bool) bool {
	return status == messages.Status_AMBIGUOUS ||
		status == messages.Status_FAILED ||
		status == messages.Status_UNDEFINED ||
		(status == messages.Status_PENDING && isStrict)
}
