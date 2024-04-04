# V2RevertTransactionRequest


## Fields

| Field                                                   | Type                                                    | Required                                                | Description                                             | Example                                                 |
| ------------------------------------------------------- | ------------------------------------------------------- | ------------------------------------------------------- | ------------------------------------------------------- | ------------------------------------------------------- |
| `AtEffectiveDate`                                       | **bool*                                                 | :heavy_minus_sign:                                      | Revert transaction at effective date of the original tx |                                                         |
| `Force`                                                 | **bool*                                                 | :heavy_minus_sign:                                      | Force revert                                            |                                                         |
| `ID`                                                    | [*big.Int](https://pkg.go.dev/math/big#Int)             | :heavy_check_mark:                                      | Transaction ID.                                         | 1234                                                    |
| `Ledger`                                                | *string*                                                | :heavy_check_mark:                                      | Name of the ledger.                                     | ledger001                                               |