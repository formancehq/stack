# OrchestrationTransaction


## Fields

| Field                                       | Type                                        | Required                                    | Description                                 | Example                                     |
| ------------------------------------------- | ------------------------------------------- | ------------------------------------------- | ------------------------------------------- | ------------------------------------------- |
| `ID`                                        | [*big.Int](https://pkg.go.dev/math/big#Int) | :heavy_check_mark:                          | N/A                                         |                                             |
| `Metadata`                                  | map[string]*string*                         | :heavy_check_mark:                          | N/A                                         | [object Object]                             |
| `Postings`                                  | [][Posting](../../models/shared/posting.md) | :heavy_check_mark:                          | N/A                                         |                                             |
| `Reference`                                 | **string*                                   | :heavy_minus_sign:                          | N/A                                         | ref:001                                     |
| `Reverted`                                  | *bool*                                      | :heavy_check_mark:                          | N/A                                         |                                             |
| `Timestamp`                                 | [time.Time](https://pkg.go.dev/time#Time)   | :heavy_check_mark:                          | N/A                                         |                                             |