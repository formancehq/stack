# ActivateConfigResponse


## Fields

| Field                                                           | Type                                                            | Required                                                        | Description                                                     |
| --------------------------------------------------------------- | --------------------------------------------------------------- | --------------------------------------------------------------- | --------------------------------------------------------------- |
| `ConfigResponse`                                                | [*shared.ConfigResponse](../../models/shared/configresponse.md) | :heavy_minus_sign:                                              | Config successfully activated.                                  |
| `ContentType`                                                   | *string*                                                        | :heavy_check_mark:                                              | N/A                                                             |
| `ErrorResponse`                                                 | [*shared.ErrorResponse](../../models/shared/errorresponse.md)   | :heavy_minus_sign:                                              | Error                                                           |
| `StatusCode`                                                    | *int*                                                           | :heavy_check_mark:                                              | N/A                                                             |
| `RawResponse`                                                   | [*http.Response](https://pkg.go.dev/net/http#Response)          | :heavy_minus_sign:                                              | N/A                                                             |