# ExpandedDebitHold


## Fields

| Field                                             | Type                                              | Required                                          | Description                                       | Example                                           |
| ------------------------------------------------- | ------------------------------------------------- | ------------------------------------------------- | ------------------------------------------------- | ------------------------------------------------- |
| `description`                                     | *str*                                             | :heavy_check_mark:                                | N/A                                               |                                                   |
| `destination`                                     | *Optional[Any]*                                   | :heavy_minus_sign:                                | N/A                                               |                                                   |
| `id`                                              | *str*                                             | :heavy_check_mark:                                | The unique ID of the hold.                        |                                                   |
| `metadata`                                        | dict[str, *str*]                                  | :heavy_check_mark:                                | Metadata associated with the hold.                |                                                   |
| `original_amount`                                 | *int*                                             | :heavy_check_mark:                                | Original amount on hold                           | 100                                               |
| `remaining`                                       | *int*                                             | :heavy_check_mark:                                | Remaining amount on hold                          | 10                                                |
| `wallet_id`                                       | *str*                                             | :heavy_check_mark:                                | The ID of the wallet the hold is associated with. |                                                   |