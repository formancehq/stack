# ScriptResponse


## Fields

| Field                                                                                        | Type                                                                                         | Required                                                                                     | Description                                                                                  | Example                                                                                      |
| -------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------- |
| `details`                                                                                    | *?string*                                                                                    | :heavy_minus_sign:                                                                           | N/A                                                                                          | https://play.numscript.org/?payload=eyJlcnJvciI6ImFjY291bnQgaGFkIGluc3VmZmljaWVudCBmdW5kcyJ9 |
| `errorCode`                                                                                  | [?\formance\stack\Models\Shared\ErrorsEnum](../../Models/Shared/ErrorsEnum.md)               | :heavy_minus_sign:                                                                           | N/A                                                                                          | INSUFFICIENT_FUND                                                                            |
| `errorMessage`                                                                               | *?string*                                                                                    | :heavy_minus_sign:                                                                           | N/A                                                                                          | account had insufficient funds                                                               |
| `transaction`                                                                                | [?\formance\stack\Models\Shared\Transaction](../../Models/Shared/Transaction.md)             | :heavy_minus_sign:                                                                           | N/A                                                                                          |                                                                                              |