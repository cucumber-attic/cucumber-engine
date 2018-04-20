# Command Type: Run Before / After Test Run Hooks

This command is sent by the program asking the caller to run the test run hooks. Once complete, the caller should send an [action complete](./action_complete.md). If there are user errors, the caller should kill the program and exit with an appropriate error message.

```
{
  "type": "run_before_test_run_hooks", // or "run_after_test_run_hooks"
}
```
