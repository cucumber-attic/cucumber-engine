package helpers

import (
	messages "github.com/cucumber/cucumber-messages-go/v3"
)

// CreateActionCompleteMessage returns a CommandActionComplete message for the given actionID
func CreateActionCompleteMessage(actionID string) *messages.Envelope {
	return &messages.Envelope{
		Message: &messages.Envelope_CommandActionComplete{
			CommandActionComplete: &messages.CommandActionComplete{
				CompletedId: actionID,
			},
		},
	}
}

// CreateActionCompleteMessageWithSnippet returns a CommandActionComplete message for the given actionID and snippet
func CreateActionCompleteMessageWithSnippet(actionID string, snippet string) *messages.Envelope {
	return &messages.Envelope{
		Message: &messages.Envelope_CommandActionComplete{
			CommandActionComplete: &messages.CommandActionComplete{
				CompletedId: actionID,
				Result: &messages.CommandActionComplete_Snippet{
					Snippet: snippet,
				},
			},
		},
	}
}

// CreateActionCompleteMessageWithTestResult returns a CommandActionComplete message for the given actionID and test result
func CreateActionCompleteMessageWithTestResult(actionID string, testResult *messages.TestResult) *messages.Envelope {
	return &messages.Envelope{
		Message: &messages.Envelope_CommandActionComplete{
			CommandActionComplete: &messages.CommandActionComplete{
				CompletedId: actionID,
				Result: &messages.CommandActionComplete_TestResult{
					TestResult: testResult,
				},
			},
		},
	}
}
