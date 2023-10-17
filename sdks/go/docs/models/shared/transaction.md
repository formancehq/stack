# Transaction


## Fields

| Field                                       | Type                                        | Required                                    | Description                                 | Example                                     |
| ------------------------------------------- | ------------------------------------------- | ------------------------------------------- | ------------------------------------------- | ------------------------------------------- |
| `ID`                                        | *int64*                                     | :heavy_check_mark:                          | N/A                                         |                                             |
| `Metadata`                                  | map[string]*string*                         | :heavy_check_mark:                          | N/A                                         |                                             |
| `Postings`                                  | [][Posting](../../models/shared/posting.md) | :heavy_check_mark:                          | N/A                                         |                                             |
| `Reference`                                 | **string*                                   | :heavy_minus_sign:                          | N/A                                         | ref:001                                     |
| `Reverted`                                  | *bool*                                      | :heavy_check_mark:                          | N/A                                         |                                             |
| `Timestamp`                                 | [time.Time](https://pkg.go.dev/time#Time)   | :heavy_check_mark:                          | N/A                                         |                                             |