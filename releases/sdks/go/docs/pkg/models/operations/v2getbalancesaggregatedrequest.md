# V2GetBalancesAggregatedRequest


## Fields

| Field                                        | Type                                         | Required                                     | Description                                  | Example                                      |
| -------------------------------------------- | -------------------------------------------- | -------------------------------------------- | -------------------------------------------- | -------------------------------------------- |
| `RequestBody`                                | map[string]*interface{}*                     | :heavy_minus_sign:                           | N/A                                          |                                              |
| `Ledger`                                     | *string*                                     | :heavy_check_mark:                           | Name of the ledger.                          | ledger001                                    |
| `Pit`                                        | [*time.Time](https://pkg.go.dev/time#Time)   | :heavy_minus_sign:                           | N/A                                          |                                              |
| `UseInsertionDate`                           | **bool*                                      | :heavy_minus_sign:                           | Use insertion date instead of effective date |                                              |