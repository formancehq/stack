# TransactionData


## Fields

| Field                                                     | Type                                                      | Required                                                  | Description                                               | Example                                                   |
| --------------------------------------------------------- | --------------------------------------------------------- | --------------------------------------------------------- | --------------------------------------------------------- | --------------------------------------------------------- |
| `Metadata`                                                | map[string]*any*                                          | :heavy_minus_sign:                                        | N/A                                                       |                                                           |
| `Postings`                                                | [][shared.Posting](../../../pkg/models/shared/posting.md) | :heavy_check_mark:                                        | N/A                                                       |                                                           |
| `Reference`                                               | **string*                                                 | :heavy_minus_sign:                                        | N/A                                                       | ref:001                                                   |
| `Timestamp`                                               | [*time.Time](https://pkg.go.dev/time#Time)                | :heavy_minus_sign:                                        | N/A                                                       |                                                           |