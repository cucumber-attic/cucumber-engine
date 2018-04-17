# Command Type: Run Before / After Test Case Hook

This command is sent by the program asking the caller to run a test case hook. Once complete, the caller should send an [action complete](./action_complete.md) command with the result.

```
{
  "type": "run_before_test_case_hook", // or "run_after_test_case_hook"

  // id of the test case
  "testCaseId": "",

  // id of the test case hook definition to run
  "testCaseHookDefinitionId": "",

  // if running an after hook, this is the result of the test case thus far
  // optionally passed to the hook
  result: {

    // "passed" or "failed" or "pending" or "skipped" or "undefined"
    status: "",

    message: "",

    // how long it took the test case to run
    duration: 0,
  },

}
```
