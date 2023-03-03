# \ServerApi

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**GetInfo**](ServerApi.md#GetInfo) | **Get** /api/ledger/_info | Show server information



## GetInfo

> ConfigInfo GetInfo(ctx).Execute()

Show server information

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    client "github.com/formancehq/formance-sdk-go"
)

func main() {

    configuration := client.NewConfiguration()
    apiClient := client.NewAPIClient(configuration)
    resp, r, err := apiClient.ServerApi.GetInfo(context.Background()).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `ServerApi.GetInfo``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `GetInfo`: ConfigInfo
    fmt.Fprintf(os.Stdout, "Response from `ServerApi.GetInfo`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiGetInfoRequest struct via the builder pattern


### Return type

[**ConfigInfo**](ConfigInfo.md)

### Authorization

[Authorization](../README.md#Authorization)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

