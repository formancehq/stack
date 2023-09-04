# AccountWithVolumesAndBalances


## Fields

| Field                       | Type                        | Required                    | Description                 | Example                     |
| --------------------------- | --------------------------- | --------------------------- | --------------------------- | --------------------------- |
| `address`                   | *str*                       | :heavy_check_mark:          | N/A                         | users:001                   |
| `balances`                  | dict[str, *int*]            | :heavy_check_mark:          | N/A                         |                             |
| `metadata`                  | dict[str, *str*]            | :heavy_check_mark:          | N/A                         |                             |
| `type`                      | *Optional[str]*             | :heavy_minus_sign:          | N/A                         | virtual                     |
| `volumes`                   | dict[str, dict[str, *int*]] | :heavy_check_mark:          | N/A                         |                             |