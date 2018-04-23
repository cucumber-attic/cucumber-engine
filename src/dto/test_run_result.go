package dto

// TestRunResult is the result of a test run
type TestRunResult struct {
	Duration int  `json:"duration"`
	Success  bool `json:"success"`
}

// NewTestRunResult creates a new test run result
func NewTestRunResult() *TestRunResult {
	return &TestRunResult{
		Success:  true,
		Duration: 0,
	}
}

// Update updates the test run result with a test case result
func (t *TestRunResult) Update(testCaseResult *TestResult, isStrict bool) {
	t.Duration += testCaseResult.Duration
	if shouldCauseFailure(testCaseResult.Status, isStrict) {
		t.Success = false
	}
}

func shouldCauseFailure(status Status, isStrict bool) bool {
	return status == StatusAmbiguous ||
		status == StatusFailed ||
		status == StatusUndefined ||
		(status == StatusPending && isStrict)
}
