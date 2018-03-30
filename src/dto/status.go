package dto

var (
	// StatusAmbiguous is the status for a step with multiple definitions
	StatusAmbiguous = "ambiguous"
	// StatusFailed is the status for a hook/step that failed
	StatusFailed = "failed"
	// StatusPassed is the status for a hook/step that passod
	StatusPassed = "passed"
	// StatusPending is the status for a step with an incomplete definition
	StatusPending = "pending"
	// StatusSkipped is the status for a hook/step that is skipped deliberately
	// to cause the scenarin to be skipped or there was a previous error
	StatusSkipped = "skipped"
	// StatusUndefined is the status for a step without a definition
	StatusUndefined = "undefined"
)
