# Command Type: Generate Snippet

This command is sent by the program asking the caller to generate a step definition snippet. Once complete, the caller should send an [action complete](./action_complete.md) with the snippet. The snippet can include `{{keywordType}}` which will be replaced by the proper step keyword (Given / When / Then) in the formatters.

```
{
  "type": "generate_snippet",

  "generatedExpressions": [
    {
      // cucumber expression as text
      "text": "",
      // array of parameter type names
      "parameterTypeNames": []
    }
  ],

  // gherkin step arguments (data table / doc string)
  pickleArguments: [],

}
```
