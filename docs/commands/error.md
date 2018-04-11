# Command Type: Error

This command is sent by the program whenever an error occurs. The caller should print the error and exit. Possible errors are gherkin parse errors, invalid tag expressions, invalid cucumber expressions, and others.

```
{
  "type": "error",

  // an error message
  "error": ""
}
```
