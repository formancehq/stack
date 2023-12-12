# AddMetadataOnTransactionRequest


## Fields

| Field                    | Type                     | Required                 | Description              | Example                  |
| ------------------------ | ------------------------ | ------------------------ | ------------------------ | ------------------------ |
| `RequestBody`            | map[string]*interface{}* | :heavy_minus_sign:       | metadata                 |                          |
| `Ledger`                 | *string*                 | :heavy_check_mark:       | Name of the ledger.      | ledger001                |
| `Txid`                   | *int64*                  | :heavy_check_mark:       | Transaction ID.          | 1234                     |