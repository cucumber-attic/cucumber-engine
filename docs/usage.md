# Usage

* Define your user interface for executing cucumber (often a CLI) and how a user defines steps/hooks and other configurations.
* Download the binary [here](https://github.com/cucumber/cucumber-engine/releases)
  * The binary can be downloaded during your package install process or on the fly as needed
  * Determine the user os and system architecture and map it to correct binary. If the os/system architecture is unexpected/unsupported, ask the user to open a issue on the language specific repo with their os/system architecture. Then open an issue on this repo if needed.
    * Note on windows: running the 386 binary on an AMD64 causes backgrounds to not be parsed correctly (found during cucumber-js integration)
* Start a subprocess that runs the binary.
  * The program can be interfaced with newline delimited json over `stdin` / `stdout`.
    * `stderr` of the program should be redirected to `stderr` of the caller
  * The program should be sent a [start](./commands/start.md) command immediately
  * The program will then send commands for the caller to complete. The caller should send a [response](./commands/action_complete.md) once the action is complete.
    * [Run before test run hooks](./commands/run_test_run_hooks.md)
    * [Initialize test case](./commands/initialize_test_case.md)
    * [Run before test case hooks](./commands/run_test_case_hook.md)
    * [Generate snippet](./commands/generate_snippet.md)
    * [Run test step](./commands/run_test_step.md)
    * [Run after test case hooks](./commands/run_test_case_hook.md)
    * [Run after test run hooks](./commands/run_test_run_hooks.md)
  * The program will also send [event](./commands/event.md) commands.
    * Use the `test-run-finished` event to see the result of the test run. Once this event is received, close the stdin stream of the program which will cause the program to exit.
  * The program may send an [error](./commands/error.md) commands
