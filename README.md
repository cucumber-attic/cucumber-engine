# cucumber-pickle-runner

Shared go binary that can be used by all language implementations.

It takes care of loading the features, filtering the pickles, and orchestrating the test run. It defers running the hooks / steps to the caller. Its primary output is events that conform to the event protocol.

# Setup

```
glide install
go get github.com/onsi/ginkgo/ginkgo
```

# Usage

* Ensure you have downloaded the proper executable for the user's machine.
* Run the executable
  * The program can be interfaced with newline delimited json over stdin / stdout.
  * The program should be sent a [start](./docs/commands/start.md)
  * The program will then send commands for the caller to complete. The caller should send a [response](./docs/commands/action_complete.md) once the action is complete.
    * Run before test run hooks
    * Initialize test case
    * Run before test case hooks
    * Run step
    * Run after test case hooks
    * Run after test run hooks
  * The program will also send [event](./docs/commands/event.md) commands
