# Command Type: Start

This command should be sent from the calling program immediately

```
{
  "type": "start",

  // the base directory of the features / support code
  // paths displayed will be relative to this
  "baseDirectory": "",

  "featuresConfig": {

    // array of paths to features that need to be loaded
    "absolutePaths": [],

    // filters to select specific scenarios to run
    "filters": {

      // array of strings, which will become regular expressions, that a scenario name must match
      // if empty, all features match
      // if multiple names provided, a scenario needs to match only one
      "names": [],

      // tag expression for what scenarios should run on
      "tagExpression": "",

      // map from feature path to array of line numbers for what scenarios to run
      // if a feature path is not present, it will run all scenarios in that feature
      "lines": {
        "/path/to/feature": [1],
        //...
      },
    },

    // the default language of feature files
    language: "",

    // what order to run scenarios in
    order: {
      // "defined" or "random"
      type: "",
      // if type is random, the seed is required
      seed: 0,
    },
  },

  "runtimeConfig": {

    // if true, after the first failure, the remaining scenarios are skipped
    "isFailFast": false,

    // if true, do not run any steps
    "isDryRun": false,

    // if true, pending steps cause the test run to fail
    "isStrict": false,

    // if -1, runs all test cases in parallel,
    // if 0, 1 or undefined, runs the test cases sequentially
    // otherwise specifies the number the can be run sequentially
    maxParallel: 1,

  },

  "supportCodeConfig": {

    // hooks to run before each test case
    "beforeTestCaseHookDefinitions": [
      {
        // a unique id for the before hook
        "id": "",

        // tag expression for what scenarios this hook should run on
        "tagExpression": "",

        // uri (absolute path) / line for where the hook was defined
        "uri": "",
        "line": ""
      }
      // ...
    ],

    // hooks to run after each test case (same format as beforeTestCaseHookDefinitions)
    "afterTestCaseHookDefinitions": [],

    "stepDefinitions": [
      {

        // a unique id for the step
        "id": "",

        "pattern": {
          // text or regexp as string
          // regexp should follow https://golang.org/pkg/regexp/syntax/
          "source": "",
          // "regular_expression" or "cucumber_expression"
          "type": ""
        },

        // uri (absolute path) / line for where the step was defined
        "uri": "",
        "line": ""
      }
      // ...
    ],

    "parameterTypes": [
      {
        // a unique name for the parameter type
        "name": "",

        // array of regexp sources
        "regexps": [],

        // whether or not this is the preferred parameter type for regular expressions
        "preferForRegexpMatch": false,

        // whether or not to use this when suggesting snippets
        "useForSnippets": false,
      }
      // ...
    ]
  }
}
```
