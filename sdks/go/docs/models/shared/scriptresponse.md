# ScriptResponse

On success, it will return a 200 status code, and the resulting transaction under the `transaction` field.

On failure, it will also return a 200 status code, and the following fields:
  - `details`: contains a URL. When there is an error parsing Numscript, the result can be difficult to read—the provided URL will render the error in an easy-to-read format.
  - `errorCode` and `error_code` (deprecated): contains the string code of the error
  - `errorMessage` and `error_message` (deprecated): contains a human-readable indication of what went wrong, for example that an account had insufficient funds, or that there was an error in the provided Numscript.



## Fields

| Field                                                                                        | Type                                                                                         | Required                                                                                     | Description                                                                                  | Example                                                                                      |
| -------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------- |
| `Details`                                                                                    | **string*                                                                                    | :heavy_minus_sign:                                                                           | N/A                                                                                          | https://play.numscript.org/?payload=eyJlcnJvciI6ImFjY291bnQgaGFkIGluc3VmZmljaWVudCBmdW5kcyJ9 |
| `ErrorCode`                                                                                  | [*ErrorsEnum](../../models/shared/errorsenum.md)                                             | :heavy_minus_sign:                                                                           | N/A                                                                                          | INSUFFICIENT_FUND                                                                            |
| `ErrorMessage`                                                                               | **string*                                                                                    | :heavy_minus_sign:                                                                           | N/A                                                                                          | account had insufficient funds                                                               |
| `Transaction`                                                                                | [*Transaction](../../models/shared/transaction.md)                                           | :heavy_minus_sign:                                                                           | N/A                                                                                          |                                                                                              |