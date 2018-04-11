package dto

// TestRunResult is the result of a test run
type TestRunResult struct {
	Duration int  `json:"duration"`
	Success  bool `json:"success"`
}
