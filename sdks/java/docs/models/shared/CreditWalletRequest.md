# CreditWalletRequest


## Fields

| Field                                       | Type                                        | Required                                    | Description                                 |
| ------------------------------------------- | ------------------------------------------- | ------------------------------------------- | ------------------------------------------- |
| `amount`                                    | [Monetary](../../models/shared/Monetary.md) | :heavy_check_mark:                          | N/A                                         |
| `balance`                                   | *String*                                    | :heavy_minus_sign:                          | The balance to credit                       |
| `metadata`                                  | Map<String, *String*>                       | :heavy_check_mark:                          | Metadata associated with the wallet.        |
| `reference`                                 | *String*                                    | :heavy_minus_sign:                          | N/A                                         |
| `sources`                                   | List<*Object*>                              | :heavy_check_mark:                          | N/A                                         |