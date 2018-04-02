# Cucumber Pickle Runner

Shared go binary that can be used by all language implementations.

It takes care of loading the features, filtering the pickles, and orchestrating the test run. It defers running the hooks / steps to the caller. Its primary output is events that conform to the event protocol.

## Usage

* Ensure you have downloaded the proper executable for the user's machine.
* Run the executable
  * The program can be interfaced with newline delimited json over stdin / stdout.
  * The program should be sent a [start](./docs/commands/start.md)
  * The program will then send commands for the caller to complete. The caller should send a [response](./docs/commands/action_complete.md) once the action is complete.
    * [Run before test run hooks](./docs/commands/run_test_run_hooks.md)
    * [Initialize test case](./docs/commands/initialize_test_case.md)
    * [Run before test case hooks](./docs/commands/run_test_case_hook.md)
    * [Run test step](./docs/commands/run_test_step.md)
    * [Run after test case hooks](./docs/commands/run_test_case_hook.md)
    * [Run after test run hooks](./docs/commands/run_test_run_hooks.md)
  * The program will also send [event](./docs/commands/event.md) commands

## Development

#### Setup

* `make setup`

#### Run linter / tests

* `make spec`
