# PostTransaction

The request body must contain at least one of the following objects:
  - `postings`: suitable for simple transactions
  - `script`: enabling more complex transactions with Numscript



## Fields

| Field                                                                  | Type                                                                   | Required                                                               | Description                                                            | Example                                                                |
| ---------------------------------------------------------------------- | ---------------------------------------------------------------------- | ---------------------------------------------------------------------- | ---------------------------------------------------------------------- | ---------------------------------------------------------------------- |
| `Metadata`                                                             | map[string]*interface{}*                                               | :heavy_minus_sign:                                                     | N/A                                                                    |                                                                        |
| `Postings`                                                             | [][Posting](../../models/shared/posting.md)                            | :heavy_minus_sign:                                                     | N/A                                                                    |                                                                        |
| `Reference`                                                            | **string*                                                              | :heavy_minus_sign:                                                     | N/A                                                                    | ref:001                                                                |
| `Script`                                                               | [*PostTransactionScript](../../models/shared/posttransactionscript.md) | :heavy_minus_sign:                                                     | N/A                                                                    |                                                                        |
| `Timestamp`                                                            | [*time.Time](https://pkg.go.dev/time#Time)                             | :heavy_minus_sign:                                                     | N/A                                                                    |                                                                        |