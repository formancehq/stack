# AddMetadataOnTransactionRequest


## Fields

| Field                                       | Type                                        | Required                                    | Description                                 | Example                                     |
| ------------------------------------------- | ------------------------------------------- | ------------------------------------------- | ------------------------------------------- | ------------------------------------------- |
| `RequestBody`                               | map[string]*any*                            | :heavy_minus_sign:                          | metadata                                    |                                             |
| `Ledger`                                    | *string*                                    | :heavy_check_mark:                          | Name of the ledger.                         | ledger001                                   |
| `Txid`                                      | [*big.Int](https://pkg.go.dev/math/big#Int) | :heavy_check_mark:                          | Transaction ID.                             | 1234                                        |