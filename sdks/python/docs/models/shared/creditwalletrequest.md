# CreditWalletRequest


## Fields

| Field                                       | Type                                        | Required                                    | Description                                 |
| ------------------------------------------- | ------------------------------------------- | ------------------------------------------- | ------------------------------------------- |
| `amount`                                    | [Monetary](../../models/shared/monetary.md) | :heavy_check_mark:                          | N/A                                         |
| `balance`                                   | *Optional[str]*                             | :heavy_minus_sign:                          | The balance to credit                       |
| `metadata`                                  | dict[str, *str*]                            | :heavy_check_mark:                          | Metadata associated with the wallet.        |
| `reference`                                 | *Optional[str]*                             | :heavy_minus_sign:                          | N/A                                         |
| `sources`                                   | list[*Any*]                                 | :heavy_check_mark:                          | N/A                                         |