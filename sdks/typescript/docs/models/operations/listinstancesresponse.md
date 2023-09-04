# ListInstancesResponse


## Fields

| Field                                                              | Type                                                               | Required                                                           | Description                                                        |
| ------------------------------------------------------------------ | ------------------------------------------------------------------ | ------------------------------------------------------------------ | ------------------------------------------------------------------ |
| `contentType`                                                      | *string*                                                           | :heavy_check_mark:                                                 | N/A                                                                |
| `error`                                                            | [shared.ErrorT](../../models/shared/errort.md)                     | :heavy_minus_sign:                                                 | General error                                                      |
| `listRunsResponse`                                                 | [shared.ListRunsResponse](../../models/shared/listrunsresponse.md) | :heavy_minus_sign:                                                 | List of workflow instances                                         |
| `statusCode`                                                       | *number*                                                           | :heavy_check_mark:                                                 | N/A                                                                |
| `rawResponse`                                                      | [AxiosResponse>](https://axios-http.com/docs/res_schema)           | :heavy_minus_sign:                                                 | N/A                                                                |