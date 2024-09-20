# OrchestrationCreditWalletRequest


## Fields

| Field                                                     | Type                                                      | Required                                                  | Description                                               |
| --------------------------------------------------------- | --------------------------------------------------------- | --------------------------------------------------------- | --------------------------------------------------------- |
| `Amount`                                                  | [shared.Monetary](../../../pkg/models/shared/monetary.md) | :heavy_check_mark:                                        | N/A                                                       |
| `Balance`                                                 | **string*                                                 | :heavy_minus_sign:                                        | The balance to credit                                     |
| `Metadata`                                                | map[string]*string*                                       | :heavy_check_mark:                                        | Metadata associated with the wallet.                      |
| `Reference`                                               | **string*                                                 | :heavy_minus_sign:                                        | N/A                                                       |
| `Sources`                                                 | [][shared.Subject](../../../pkg/models/shared/subject.md) | :heavy_check_mark:                                        | N/A                                                       |
| `Timestamp`                                               | [*time.Time](https://pkg.go.dev/time#Time)                | :heavy_minus_sign:                                        | N/A                                                       |