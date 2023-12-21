# PaymentslistAccountsResponse


## Fields

| Field                                                           | Type                                                            | Required                                                        | Description                                                     |
| --------------------------------------------------------------- | --------------------------------------------------------------- | --------------------------------------------------------------- | --------------------------------------------------------------- |
| `AccountsCursor`                                                | [*shared.AccountsCursor](../../models/shared/accountscursor.md) | :heavy_minus_sign:                                              | OK                                                              |
| `ContentType`                                                   | *string*                                                        | :heavy_check_mark:                                              | HTTP response content type for this operation                   |
| `StatusCode`                                                    | *int*                                                           | :heavy_check_mark:                                              | HTTP response status code for this operation                    |
| `RawResponse`                                                   | [*http.Response](https://pkg.go.dev/net/http#Response)          | :heavy_minus_sign:                                              | Raw HTTP response; suitable for custom response parsing         |