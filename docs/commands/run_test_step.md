# Command Type: Run Test Step

This command is sent by the program asking the caller to run a test step. Once complete, the caller should send an [action complete](./action_complete.md) command with the result.

```
{
  "type": "run_test_step",

  // id of the test case
  "testCaseId": "",

  // id of the step to run
  "testStepDefinitionId": "",

  // captures groups
  "patternCaptureGroups": [
    {
      // the string value from the matching step
      "value": "",

      // the parameter type name this capture should be transformed with
      //   empty string if none
      "parameterTypeName": "",

    },
    // ...
  ],

  // gherkin step arguments (data table / doc string)
  arguments: [],
}
```
