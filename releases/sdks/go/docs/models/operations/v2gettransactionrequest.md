# V2GetTransactionRequest


## Fields

| Field                                       | Type                                        | Required                                    | Description                                 | Example                                     |
| ------------------------------------------- | ------------------------------------------- | ------------------------------------------- | ------------------------------------------- | ------------------------------------------- |
| `Expand`                                    | **string*                                   | :heavy_minus_sign:                          | N/A                                         |                                             |
| `ID`                                        | [*big.Int](https://pkg.go.dev/math/big#Int) | :heavy_check_mark:                          | Transaction ID.                             | 1234                                        |
| `Ledger`                                    | *string*                                    | :heavy_check_mark:                          | Name of the ledger.                         | ledger001                                   |
| `Pit`                                       | [*time.Time](https://pkg.go.dev/time#Time)  | :heavy_minus_sign:                          | N/A                                         |                                             |