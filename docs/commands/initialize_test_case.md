# Command Type: Initialize Test Case

This command is sent by the program asking the caller to initialize a test case. Once complete, the caller should send an [action complete](./action_complete.md). A test case id is included that will be sent with commands to run hooks and steps for this specific test case.

```
{
  "type": "initialize_test_case",

  // id of the test case
  "testCaseId": "",

}
```
