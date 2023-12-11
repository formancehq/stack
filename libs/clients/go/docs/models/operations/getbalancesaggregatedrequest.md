# GetBalancesAggregatedRequest


## Fields

| Field                                                                     | Type                                                                      | Required                                                                  | Description                                                               | Example                                                                   |
| ------------------------------------------------------------------------- | ------------------------------------------------------------------------- | ------------------------------------------------------------------------- | ------------------------------------------------------------------------- | ------------------------------------------------------------------------- |
| `Address`                                                                 | **string*                                                                 | :heavy_minus_sign:                                                        | Filter balances involving given account, either as source or destination. | users:001                                                                 |
| `Ledger`                                                                  | *string*                                                                  | :heavy_check_mark:                                                        | Name of the ledger.                                                       | ledger001                                                                 |