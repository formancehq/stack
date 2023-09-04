# PostTransaction

The request body must contain at least one of the following objects:
  - `postings`: suitable for simple transactions
  - `script`: enabling more complex transactions with Numscript



## Fields

| Field                                                                  | Type                                                                   | Required                                                               | Description                                                            | Example                                                                |
| ---------------------------------------------------------------------- | ---------------------------------------------------------------------- | ---------------------------------------------------------------------- | ---------------------------------------------------------------------- | ---------------------------------------------------------------------- |
| `metadata`                                                             | array<string, *string*>                                                | :heavy_check_mark:                                                     | N/A                                                                    |                                                                        |
| `postings`                                                             | array<[Posting](../../models/shared/Posting.md)>                       | :heavy_minus_sign:                                                     | N/A                                                                    |                                                                        |
| `reference`                                                            | *?string*                                                              | :heavy_minus_sign:                                                     | N/A                                                                    | ref:001                                                                |
| `script`                                                               | [?PostTransactionScript](../../models/shared/PostTransactionScript.md) | :heavy_minus_sign:                                                     | N/A                                                                    |                                                                        |
| `timestamp`                                                            | [\DateTime](https://www.php.net/manual/en/class.datetime.php)          | :heavy_minus_sign:                                                     | N/A                                                                    |                                                                        |