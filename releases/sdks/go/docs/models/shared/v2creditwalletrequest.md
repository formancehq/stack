# V2CreditWalletRequest


## Fields

| Field                                           | Type                                            | Required                                        | Description                                     |
| ----------------------------------------------- | ----------------------------------------------- | ----------------------------------------------- | ----------------------------------------------- |
| `Amount`                                        | [V2Monetary](../../models/shared/v2monetary.md) | :heavy_check_mark:                              | N/A                                             |
| `Balance`                                       | **string*                                       | :heavy_minus_sign:                              | The balance to credit                           |
| `Metadata`                                      | map[string]*string*                             | :heavy_check_mark:                              | Metadata associated with the wallet.            |
| `Reference`                                     | **string*                                       | :heavy_minus_sign:                              | N/A                                             |
| `Sources`                                       | [][V2Subject](../../models/shared/v2subject.md) | :heavy_check_mark:                              | N/A                                             |