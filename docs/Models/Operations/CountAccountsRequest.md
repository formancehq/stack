# CountAccountsRequest


## Fields

| Field                                                                                                                            | Type                                                                                                                             | Required                                                                                                                         | Description                                                                                                                      | Example                                                                                                                          |
| -------------------------------------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------- |
| `address`                                                                                                                        | *?string*                                                                                                                        | :heavy_minus_sign:                                                                                                               | Filter accounts by address pattern (regular expression placed between ^ and $).                                                  | users:.+                                                                                                                         |
| `ledger`                                                                                                                         | *string*                                                                                                                         | :heavy_check_mark:                                                                                                               | Name of the ledger.                                                                                                              | ledger001                                                                                                                        |
| `metadata`                                                                                                                       | array<string, *mixed*>                                                                                                           | :heavy_minus_sign:                                                                                                               | Filter accounts by metadata key value pairs. The filter can be used like this metadata[key]=value1&metadata[a.nested.key]=value2 |                                                                                                                                  |