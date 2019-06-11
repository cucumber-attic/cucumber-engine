package matchers

import (
	"fmt"
	"reflect"

	messages "github.com/cucumber/cucumber-messages-go/v3"
	"github.com/onsi/gomega/format"
	"github.com/onsi/gomega/types"
)

// BeAMessageOfType is a matcher that validates a *message.Envelope's message type
func BeAMessageOfType(expected interface{}) types.GomegaMatcher {
	return &messageOfTypeMatcher{
		expectedType: reflect.TypeOf(expected),
	}
}

type messageOfTypeMatcher struct {
	expectedType reflect.Type
}

func (m *messageOfTypeMatcher) Match(actual interface{}) (success bool, err error) {
	msg, ok := actual.(*messages.Envelope)
	if !ok {
		return false, fmt.Errorf("BeACommandWithType matcher expects a *messages.Envelope.  Got:\n%s", format.Object(actual, 1))
	}

	var actualInner interface{}
	switch x := msg.Message.(type) {
	case *messages.Envelope_Source:
		actualInner = x.Source
	case *messages.Envelope_GherkinDocument:
		actualInner = x.GherkinDocument
	case *messages.Envelope_Pickle:
		actualInner = x.Pickle
	case *messages.Envelope_Attachment:
		actualInner = x.Attachment
	case *messages.Envelope_TestCaseStarted:
		actualInner = x.TestCaseStarted
	case *messages.Envelope_TestStepStarted:
		actualInner = x.TestStepStarted
	case *messages.Envelope_TestStepFinished:
		actualInner = x.TestStepFinished
	case *messages.Envelope_TestCaseFinished:
		actualInner = x.TestCaseFinished
	case *messages.Envelope_PickleAccepted:
		actualInner = x.PickleAccepted
	case *messages.Envelope_PickleRejected:
		actualInner = x.PickleRejected
	case *messages.Envelope_TestCasePrepared:
		actualInner = x.TestCasePrepared
	case *messages.Envelope_TestRunStarted:
		actualInner = x.TestRunStarted
	case *messages.Envelope_TestRunFinished:
		actualInner = x.TestRunFinished
	case *messages.Envelope_CommandStart:
		actualInner = x.CommandStart
	case *messages.Envelope_CommandActionComplete:
		actualInner = x.CommandActionComplete
	case *messages.Envelope_CommandRunBeforeTestRunHooks:
		actualInner = x.CommandRunBeforeTestRunHooks
	case *messages.Envelope_CommandInitializeTestCase:
		actualInner = x.CommandInitializeTestCase
	case *messages.Envelope_CommandRunBeforeTestCaseHook:
		actualInner = x.CommandRunBeforeTestCaseHook
	case *messages.Envelope_CommandRunTestStep:
		actualInner = x.CommandRunTestStep
	case *messages.Envelope_CommandRunAfterTestCaseHook:
		actualInner = x.CommandRunAfterTestCaseHook
	case *messages.Envelope_CommandRunAfterTestRunHooks:
		actualInner = x.CommandRunAfterTestRunHooks
	case *messages.Envelope_CommandGenerateSnippet:
		actualInner = x.CommandGenerateSnippet
	case *messages.Envelope_CommandError:
		actualInner = x.CommandError
	default:
		panic(fmt.Errorf("Unexpected message type: %v", msg))
	}
	actualType := reflect.TypeOf(actualInner)

	return actualType.AssignableTo(m.expectedType), nil
}

func (m *messageOfTypeMatcher) FailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Expected\n%s\nto have type %s", format.Object(actual, 1), m.expectedType)
}

func (m *messageOfTypeMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Expected\n%s\nnot to have type %s", format.Object(actual, 1), m.expectedType)
}
