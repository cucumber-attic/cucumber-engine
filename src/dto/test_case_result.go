package dto

// TestResult is the result of a test case
type TestResult struct {
	Duration  int    `json:"duration"`
	Status    string `json:"status"`
	Exception string `json:"exception"`
}
