# GetBalancesAggregatedRequest


## Fields

| Field                                                                     | Type                                                                      | Required                                                                  | Description                                                               | Example                                                                   |
| ------------------------------------------------------------------------- | ------------------------------------------------------------------------- | ------------------------------------------------------------------------- | ------------------------------------------------------------------------- | ------------------------------------------------------------------------- |
| `address`                                                                 | *string*                                                                  | :heavy_minus_sign:                                                        | Filter balances involving given account, either as source or destination. | users:001                                                                 |
| `ledger`                                                                  | *string*                                                                  | :heavy_check_mark:                                                        | Name of the ledger.                                                       | ledger001                                                                 |