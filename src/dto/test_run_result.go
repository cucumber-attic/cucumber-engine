package dto

import messages "github.com/cucumber/cucumber-messages-go/v3"

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

func shouldCauseFailure(status messages.TestResult_Status, isStrict bool) bool {
	return status == messages.TestResult_AMBIGUOUS ||
		status == messages.TestResult_FAILED ||
		status == messages.TestResult_UNDEFINED ||
		(status == messages.TestResult_PENDING && isStrict)
}
