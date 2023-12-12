# PolicyRequest


## Fields

| Field                                                 | Type                                                  | Required                                              | Description                                           | Example                                               |
| ----------------------------------------------------- | ----------------------------------------------------- | ----------------------------------------------------- | ----------------------------------------------------- | ----------------------------------------------------- |
| `LedgerName`                                          | *string*                                              | :heavy_check_mark:                                    | N/A                                                   | default                                               |
| `LedgerQuery`                                         | *string*                                              | :heavy_check_mark:                                    | N/A                                                   | {"$match": {"metadata[reconciliation]": "pool:main"}} |
| `Name`                                                | *string*                                              | :heavy_check_mark:                                    | N/A                                                   | XXX                                                   |
| `PaymentsPoolID`                                      | *string*                                              | :heavy_check_mark:                                    | N/A                                                   | XXX                                                   |