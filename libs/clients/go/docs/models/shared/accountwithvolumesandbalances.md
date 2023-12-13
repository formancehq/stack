# AccountWithVolumesAndBalances


## Fields

| Field                                                  | Type                                                   | Required                                               | Description                                            | Example                                                |
| ------------------------------------------------------ | ------------------------------------------------------ | ------------------------------------------------------ | ------------------------------------------------------ | ------------------------------------------------------ |
| `Address`                                              | *string*                                               | :heavy_check_mark:                                     | N/A                                                    | users:001                                              |
| `Balances`                                             | map[string][*big.Int](https://pkg.go.dev/math/big#Int) | :heavy_minus_sign:                                     | N/A                                                    | [object Object]                                        |
| `Metadata`                                             | map[string]*interface{}*                               | :heavy_minus_sign:                                     | N/A                                                    | [object Object]                                        |
| `Type`                                                 | **string*                                              | :heavy_minus_sign:                                     | N/A                                                    | virtual                                                |
| `Volumes`                                              | map[string][Volume](../../models/shared/volume.md)     | :heavy_minus_sign:                                     | N/A                                                    | [object Object]                                        |