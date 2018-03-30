package matchers

import (
	"fmt"

	"github.com/cucumber/cucumber-pickle-runner/src/dto"
	"github.com/onsi/gomega/format"
	"github.com/onsi/gomega/types"
)

// BeACommandWithType is a matcher that validates a *dto.Command Type field
func BeACommandWithType(expectedType dto.CommandType) types.GomegaMatcher {
	return &commandWithTypeMatcher{
		expectedType: expectedType,
	}
}

type commandWithTypeMatcher struct {
	expectedType dto.CommandType
}

func (c *commandWithTypeMatcher) Match(actual interface{}) (success bool, err error) {
	command, ok := actual.(*dto.Command)
	if !ok {
		return false, fmt.Errorf("BeACommandWithType matcher expects a *dto.Command.  Got:\n%s", format.Object(actual, 1))
	}
	return command.Type == c.expectedType, nil
}

func (c *commandWithTypeMatcher) FailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Expected\n%s\nto have type %s", format.Object(actual, 1), c.expectedType)
}

func (c *commandWithTypeMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Expected\n%s\nnot to have type %s", format.Object(actual, 1), c.expectedType)
}
