# CreditWalletRequest


## Fields

| Field                                       | Type                                        | Required                                    | Description                                 |
| ------------------------------------------- | ------------------------------------------- | ------------------------------------------- | ------------------------------------------- |
| `amount`                                    | [Monetary](../../models/shared/monetary.md) | :heavy_check_mark:                          | N/A                                         |
| `balance`                                   | *string*                                    | :heavy_minus_sign:                          | The balance to credit                       |
| `metadata`                                  | Record<string, *string*>                    | :heavy_check_mark:                          | Metadata associated with the wallet.        |
| `reference`                                 | *string*                                    | :heavy_minus_sign:                          | N/A                                         |
| `sources`                                   | *any*[]                                     | :heavy_check_mark:                          | N/A                                         |