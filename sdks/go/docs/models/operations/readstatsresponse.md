# ReadStatsResponse


## Fields

| Field                                                         | Type                                                          | Required                                                      | Description                                                   |
| ------------------------------------------------------------- | ------------------------------------------------------------- | ------------------------------------------------------------- | ------------------------------------------------------------- |
| `ContentType`                                                 | *string*                                                      | :heavy_check_mark:                                            | N/A                                                           |
| `ErrorResponse`                                               | [*shared.ErrorResponse](../../models/shared/errorresponse.md) | :heavy_minus_sign:                                            | Error                                                         |
| `StatsResponse`                                               | [*shared.StatsResponse](../../models/shared/statsresponse.md) | :heavy_minus_sign:                                            | OK                                                            |
| `StatusCode`                                                  | *int*                                                         | :heavy_check_mark:                                            | N/A                                                           |
| `RawResponse`                                                 | [*http.Response](https://pkg.go.dev/net/http#Response)        | :heavy_minus_sign:                                            | N/A                                                           |