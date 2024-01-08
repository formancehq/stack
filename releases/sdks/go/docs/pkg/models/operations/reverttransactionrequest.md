# RevertTransactionRequest


## Fields

| Field                                       | Type                                        | Required                                    | Description                                 | Example                                     |
| ------------------------------------------- | ------------------------------------------- | ------------------------------------------- | ------------------------------------------- | ------------------------------------------- |
| `DisableChecks`                             | **bool*                                     | :heavy_minus_sign:                          | Allow to disable balances checks            |                                             |
| `Ledger`                                    | *string*                                    | :heavy_check_mark:                          | Name of the ledger.                         | ledger001                                   |
| `Txid`                                      | [*big.Int](https://pkg.go.dev/math/big#Int) | :heavy_check_mark:                          | Transaction ID.                             | 1234                                        |