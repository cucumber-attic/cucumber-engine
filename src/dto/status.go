package dto

// Status is an enumeration of the available values for the Status field in the TestResult struct
type Status string

var (
	// StatusAmbiguous is the status for a step with multiple definitions
	StatusAmbiguous = Status("ambiguous")
	// StatusFailed is the status for a hook/step that failed
	StatusFailed = Status("failed")
	// StatusPassed is the status for a hook/step that passod
	StatusPassed = Status("passed")
	// StatusPending is the status for a step with an incomplete definition
	StatusPending = Status("pending")
	// StatusSkipped is the status for a hook/step that is skipped deliberately
	// to cause the scenarin to be skipped or there was a previous error
	StatusSkipped = Status("skipped")
	// StatusUndefined is the status for a step without a definition
	StatusUndefined = Status("undefined")
)
