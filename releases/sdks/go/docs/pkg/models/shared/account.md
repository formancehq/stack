# Account


## Fields

| Field                                         | Type                                          | Required                                      | Description                                   | Example                                       |
| --------------------------------------------- | --------------------------------------------- | --------------------------------------------- | --------------------------------------------- | --------------------------------------------- |
| `Address`                                     | *string*                                      | :heavy_check_mark:                            | N/A                                           | users:001                                     |
| `Metadata`                                    | map[string]*interface{}*                      | :heavy_minus_sign:                            | N/A                                           | {"admin":true,"a":{"nested":{"key":"value"}}} |
| `Type`                                        | **string*                                     | :heavy_minus_sign:                            | N/A                                           | virtual                                       |