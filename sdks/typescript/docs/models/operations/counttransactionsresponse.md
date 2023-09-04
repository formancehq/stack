# CountTransactionsResponse


## Fields

| Field                                                        | Type                                                         | Required                                                     | Description                                                  |
| ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ |
| `contentType`                                                | *string*                                                     | :heavy_check_mark:                                           | N/A                                                          |
| `errorResponse`                                              | [shared.ErrorResponse](../../models/shared/errorresponse.md) | :heavy_minus_sign:                                           | Error                                                        |
| `headers`                                                    | Record<string, *string*[]>                                   | :heavy_minus_sign:                                           | N/A                                                          |
| `statusCode`                                                 | *number*                                                     | :heavy_check_mark:                                           | N/A                                                          |
| `rawResponse`                                                | [AxiosResponse>](https://axios-http.com/docs/res_schema)     | :heavy_minus_sign:                                           | N/A                                                          |