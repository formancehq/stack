# Wallet


## Fields

| Field                                                                  | Type                                                                   | Required                                                               | Description                                                            |
| ---------------------------------------------------------------------- | ---------------------------------------------------------------------- | ---------------------------------------------------------------------- | ---------------------------------------------------------------------- |
| `Balances`                                                             | [*shared.WalletBalances](../../../pkg/models/shared/walletbalances.md) | :heavy_minus_sign:                                                     | N/A                                                                    |
| `CreatedAt`                                                            | [time.Time](https://pkg.go.dev/time#Time)                              | :heavy_check_mark:                                                     | N/A                                                                    |
| `ID`                                                                   | *string*                                                               | :heavy_check_mark:                                                     | The unique ID of the wallet.                                           |
| `Ledger`                                                               | *string*                                                               | :heavy_check_mark:                                                     | N/A                                                                    |
| `Metadata`                                                             | map[string]*string*                                                    | :heavy_check_mark:                                                     | Metadata associated with the wallet.                                   |
| `Name`                                                                 | *string*                                                               | :heavy_check_mark:                                                     | N/A                                                                    |