# AccountWithVolumesAndBalances


## Fields

| Field                            | Type                             | Required                         | Description                      | Example                          |
| -------------------------------- | -------------------------------- | -------------------------------- | -------------------------------- | -------------------------------- |
| `address`                        | *String*                         | :heavy_check_mark:               | N/A                              | users:001                        |
| `balances`                       | Map<String, *Long*>              | :heavy_check_mark:               | N/A                              |                                  |
| `metadata`                       | Map<String, *String*>            | :heavy_check_mark:               | N/A                              |                                  |
| `type`                           | *String*                         | :heavy_minus_sign:               | N/A                              | virtual                          |
| `volumes`                        | Map<String, Map<String, *Long*>> | :heavy_check_mark:               | N/A                              |                                  |