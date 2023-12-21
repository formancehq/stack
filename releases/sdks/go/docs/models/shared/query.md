# Query


## Fields

| Field                                        | Type                                         | Required                                     | Description                                  | Example                                      |
| -------------------------------------------- | -------------------------------------------- | -------------------------------------------- | -------------------------------------------- | -------------------------------------------- |
| `After`                                      | []*string*                                   | :heavy_minus_sign:                           | N/A                                          | users:002                                    |
| `Cursor`                                     | **string*                                    | :heavy_minus_sign:                           | N/A                                          | YXVsdCBhbmQgYSBtYXhpbXVtIG1heF9yZXN1bHRzLol= |
| `Ledgers`                                    | []*string*                                   | :heavy_minus_sign:                           | N/A                                          | quickstart                                   |
| `PageSize`                                   | **int64*                                     | :heavy_minus_sign:                           | N/A                                          |                                              |
| `Policy`                                     | **string*                                    | :heavy_minus_sign:                           | N/A                                          | OR                                           |
| `Raw`                                        | [*QueryRaw](../../models/shared/queryraw.md) | :heavy_minus_sign:                           | N/A                                          |                                              |
| `Sort`                                       | **string*                                    | :heavy_minus_sign:                           | N/A                                          | id:asc                                       |
| `Target`                                     | **string*                                    | :heavy_minus_sign:                           | N/A                                          |                                              |
| `Terms`                                      | []*string*                                   | :heavy_minus_sign:                           | N/A                                          | destination=central_bank1                    |