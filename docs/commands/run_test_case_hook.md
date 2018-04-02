# Command Type: Run Before / After Test Case Hook

This command is sent by the program asking the caller to run a test case hook. Once complete, the caller should send an [action complete](./action_complete.md) command with the result.

```
{
  "type": "run_before_test_case_hook", // or "run_after_test_case_hook"

  // id of the test case
  "testCaseId": "",

  // id of the test case hook definition to run
  "testCaseHookDefinitionId": "",

}
```
