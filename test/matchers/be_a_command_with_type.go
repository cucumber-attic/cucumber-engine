package matchers

import (
	"fmt"
	"reflect"

	messages "github.com/cucumber/cucumber-messages-go/v2"
	"github.com/onsi/gomega/format"
	"github.com/onsi/gomega/types"
)

// BeAMessageOfType is a matcher that validates a *message.Wrapper's message type
func BeAMessageOfType(expected interface{}) types.GomegaMatcher {
	return &messageOfTypeMatcher{
		expectedType: reflect.TypeOf(expected),
	}
}

type messageOfTypeMatcher struct {
	expectedType reflect.Type
}

func (m *messageOfTypeMatcher) Match(actual interface{}) (success bool, err error) {
	msg, ok := actual.(*messages.Wrapper)
	if !ok {
		return false, fmt.Errorf("BeACommandWithType matcher expects a *messages.Wrapper.  Got:\n%s", format.Object(actual, 1))
	}

	var actualInner interface{}
	switch x := msg.Message.(type) {
	case *messages.Wrapper_Source:
		actualInner = x.Source
	case *messages.Wrapper_GherkinDocument:
		actualInner = x.GherkinDocument
	case *messages.Wrapper_Pickle:
		actualInner = x.Pickle
	case *messages.Wrapper_Attachment:
		actualInner = x.Attachment
	case *messages.Wrapper_TestCaseStarted:
		actualInner = x.TestCaseStarted
	case *messages.Wrapper_TestStepStarted:
		actualInner = x.TestStepStarted
	case *messages.Wrapper_TestStepFinished:
		actualInner = x.TestStepFinished
	case *messages.Wrapper_TestCaseFinished:
		actualInner = x.TestCaseFinished
	case *messages.Wrapper_PickleAccepted:
		actualInner = x.PickleAccepted
	case *messages.Wrapper_PickleRejected:
		actualInner = x.PickleRejected
	case *messages.Wrapper_TestCasePrepared:
		actualInner = x.TestCasePrepared
	case *messages.Wrapper_TestRunStarted:
		actualInner = x.TestRunStarted
	case *messages.Wrapper_TestRunFinished:
		actualInner = x.TestRunFinished
	case *messages.Wrapper_CommandStart:
		actualInner = x.CommandStart
	case *messages.Wrapper_CommandActionComplete:
		actualInner = x.CommandActionComplete
	case *messages.Wrapper_CommandRunBeforeTestRunHooks:
		actualInner = x.CommandRunBeforeTestRunHooks
	case *messages.Wrapper_CommandInitializeTestCase:
		actualInner = x.CommandInitializeTestCase
	case *messages.Wrapper_CommandRunBeforeTestCaseHook:
		actualInner = x.CommandRunBeforeTestCaseHook
	case *messages.Wrapper_CommandRunTestStep:
		actualInner = x.CommandRunTestStep
	case *messages.Wrapper_CommandRunAfterTestCaseHook:
		actualInner = x.CommandRunAfterTestCaseHook
	case *messages.Wrapper_CommandRunAfterTestRunHooks:
		actualInner = x.CommandRunAfterTestRunHooks
	case *messages.Wrapper_CommandGenerateSnippet:
		actualInner = x.CommandGenerateSnippet
	case *messages.Wrapper_CommandError:
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
