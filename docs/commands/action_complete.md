# Command Type: Action Complete

This command should be sent by the caller after completing an action.

```
{
  "type": "action_complete",

  // the id that of the command that has been completed
  "responseTo": "",

  // the result if the command was to run a test case hook / step
  "hookOrStepResult": {

    // "passed" or "failed" or "pending" or "skipped"
    "status": "",

    // how long it took the step or hook to run
    "duration": 0,

    // if the status is "failed", this should be the error converted to a string
    // exactly how it should be presented to the user
    "message": ""

  },

  // the snippet if the command was to generate a snippet
  // `{{keyword}}` will be replaced by the proper step keyword in the formatters
  "snippet": ""
}
```
