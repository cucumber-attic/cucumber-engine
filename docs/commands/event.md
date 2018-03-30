# Command Type: Event

This command is sent by the program whenever an event occurs. The events are defined [here](https://docs.cucumber.io/event-protocol/) and include these [proposed updates](https://github.com/cucumber/cucumber/pull/172). The caller should pass events to formatters and use the `test-run-finished` event to see the result of the test run.

```
{
  "type": "event",

  // an event object
  "event": {}
}
```
