# ExpandedDebitHold


## Fields

| Field                                             | Type                                              | Required                                          | Description                                       | Example                                           |
| ------------------------------------------------- | ------------------------------------------------- | ------------------------------------------------- | ------------------------------------------------- | ------------------------------------------------- |
| `description`                                     | *String*                                          | :heavy_check_mark:                                | N/A                                               |                                                   |
| `destination`                                     | *Object*                                          | :heavy_minus_sign:                                | N/A                                               |                                                   |
| `id`                                              | *String*                                          | :heavy_check_mark:                                | The unique ID of the hold.                        |                                                   |
| `metadata`                                        | Map<String, *String*>                             | :heavy_check_mark:                                | Metadata associated with the hold.                |                                                   |
| `originalAmount`                                  | *Long*                                            | :heavy_check_mark:                                | Original amount on hold                           | 100                                               |
| `remaining`                                       | *Long*                                            | :heavy_check_mark:                                | Remaining amount on hold                          | 10                                                |
| `walletID`                                        | *String*                                          | :heavy_check_mark:                                | The ID of the wallet the hold is associated with. |                                                   |