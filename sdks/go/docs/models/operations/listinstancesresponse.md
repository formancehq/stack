# ListInstancesResponse


## Fields

| Field                                                               | Type                                                                | Required                                                            | Description                                                         |
| ------------------------------------------------------------------- | ------------------------------------------------------------------- | ------------------------------------------------------------------- | ------------------------------------------------------------------- |
| `ContentType`                                                       | *string*                                                            | :heavy_check_mark:                                                  | N/A                                                                 |
| `Error`                                                             | [*shared.Error](../../models/shared/error.md)                       | :heavy_minus_sign:                                                  | General error                                                       |
| `ListRunsResponse`                                                  | [*shared.ListRunsResponse](../../models/shared/listrunsresponse.md) | :heavy_minus_sign:                                                  | List of workflow instances                                          |
| `StatusCode`                                                        | *int*                                                               | :heavy_check_mark:                                                  | N/A                                                                 |
| `RawResponse`                                                       | [*http.Response](https://pkg.go.dev/net/http#Response)              | :heavy_minus_sign:                                                  | N/A                                                                 |