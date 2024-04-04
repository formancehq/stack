# V2GetVolumesWithBalancesRequest


## Fields

| Field                                        | Type                                         | Required                                     | Description                                  | Example                                      |
| -------------------------------------------- | -------------------------------------------- | -------------------------------------------- | -------------------------------------------- | -------------------------------------------- |
| `RequestBody`                                | map[string]*interface{}*                     | :heavy_minus_sign:                           | N/A                                          |                                              |
| `Ledger`                                     | *string*                                     | :heavy_check_mark:                           | Name of the ledger.                          | ledger001                                    |
| `Oot`                                        | [*time.Time](https://pkg.go.dev/time#Time)   | :heavy_minus_sign:                           | N/A                                          |                                              |
| `Pit`                                        | [*time.Time](https://pkg.go.dev/time#Time)   | :heavy_minus_sign:                           | N/A                                          |                                              |
| `UseInsertionDate`                           | **bool*                                      | :heavy_minus_sign:                           | Use insertion date instead of effective date |                                              |