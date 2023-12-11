# V2Hold


## Fields

| Field                                             | Type                                              | Required                                          | Description                                       |
| ------------------------------------------------- | ------------------------------------------------- | ------------------------------------------------- | ------------------------------------------------- |
| `Description`                                     | *string*                                          | :heavy_check_mark:                                | N/A                                               |
| `Destination`                                     | [*V2Subject](../../models/shared/v2subject.md)    | :heavy_minus_sign:                                | N/A                                               |
| `ID`                                              | *string*                                          | :heavy_check_mark:                                | The unique ID of the hold.                        |
| `Metadata`                                        | map[string]*string*                               | :heavy_check_mark:                                | Metadata associated with the hold.                |
| `WalletID`                                        | *string*                                          | :heavy_check_mark:                                | The ID of the wallet the hold is associated with. |