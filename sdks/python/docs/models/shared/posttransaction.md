# PostTransaction

The request body must contain at least one of the following objects:
  - `postings`: suitable for simple transactions
  - `script`: enabling more complex transactions with Numscript



## Fields

| Field                                                                           | Type                                                                            | Required                                                                        | Description                                                                     | Example                                                                         |
| ------------------------------------------------------------------------------- | ------------------------------------------------------------------------------- | ------------------------------------------------------------------------------- | ------------------------------------------------------------------------------- | ------------------------------------------------------------------------------- |
| `metadata`                                                                      | dict[str, *Any*]                                                                | :heavy_minus_sign:                                                              | N/A                                                                             |                                                                                 |
| `postings`                                                                      | list[[Posting](../../models/shared/posting.md)]                                 | :heavy_minus_sign:                                                              | N/A                                                                             |                                                                                 |
| `reference`                                                                     | *Optional[str]*                                                                 | :heavy_minus_sign:                                                              | N/A                                                                             | ref:001                                                                         |
| `script`                                                                        | [Optional[PostTransactionScript]](../../models/shared/posttransactionscript.md) | :heavy_minus_sign:                                                              | N/A                                                                             |                                                                                 |
| `timestamp`                                                                     | [date](https://docs.python.org/3/library/datetime.html#date-objects)            | :heavy_minus_sign:                                                              | N/A                                                                             |                                                                                 |