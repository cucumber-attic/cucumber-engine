package matchers

import (
	"fmt"
	"reflect"

	"github.com/cucumber/cucumber-engine/src/dto"
	"github.com/onsi/gomega/format"
	"github.com/onsi/gomega/types"
)

// BeACommandWithEventAssignableToTypeOf is a matcher that validates a *dto.Command Type is
// event and that the event is AssignableTo to the type of the given object
func BeACommandWithEventAssignableToTypeOf(expected interface{}) types.GomegaMatcher {
	return &commandWithEventAssignableToTypeOf{
		expected: expected,
	}
}

type commandWithEventAssignableToTypeOf struct {
	expected interface{}
}

func (c *commandWithEventAssignableToTypeOf) Match(actual interface{}) (success bool, err error) {
	command, ok := actual.(*dto.Command)
	if !ok {
		return false, fmt.Errorf("BeACommandWithEventOfType matcher expects a *dto.Command.  Got:\n%s", format.Object(actual, 1))
	}
	if command.Type != dto.CommandTypeEvent {
		return false, fmt.Errorf("BeACommandWithEventOfType matcher expects a *dto.Command with type dto.CommandTypeEvent.  Got:\n%s", format.Object(actual, 1))
	}
	actualType := reflect.TypeOf(command.Event)
	expectedType := reflect.TypeOf(c.expected)
	return actualType.AssignableTo(expectedType), nil
}

func (c *commandWithEventAssignableToTypeOf) FailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Expected\n%s\nto have event with type %s. Got: %s ", format.Object(actual, 1), reflect.TypeOf(c.expected), reflect.TypeOf(actual.(*dto.Command).Event))
}

func (c *commandWithEventAssignableToTypeOf) NegatedFailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Expected\n%s\nnot to have event with type %s. Got: %s ", format.Object(actual, 1), reflect.TypeOf(c.expected), reflect.TypeOf(actual.(*dto.Command).Event))
}
