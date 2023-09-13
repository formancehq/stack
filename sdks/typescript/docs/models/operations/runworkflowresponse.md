# RunWorkflowResponse


## Fields

| Field                                                                    | Type                                                                     | Required                                                                 | Description                                                              |
| ------------------------------------------------------------------------ | ------------------------------------------------------------------------ | ------------------------------------------------------------------------ | ------------------------------------------------------------------------ |
| `contentType`                                                            | *string*                                                                 | :heavy_check_mark:                                                       | N/A                                                                      |
| `error`                                                                  | [shared.ErrorT](../../models/shared/errort.md)                           | :heavy_minus_sign:                                                       | General error                                                            |
| `runWorkflowResponse`                                                    | [shared.RunWorkflowResponse](../../models/shared/runworkflowresponse.md) | :heavy_minus_sign:                                                       | The workflow instance                                                    |
| `statusCode`                                                             | *number*                                                                 | :heavy_check_mark:                                                       | N/A                                                                      |
| `rawResponse`                                                            | [AxiosResponse>](https://axios-http.com/docs/res_schema)                 | :heavy_minus_sign:                                                       | N/A                                                                      |