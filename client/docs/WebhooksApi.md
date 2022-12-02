# \WebhooksApi

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**ActivateOneConfig**](WebhooksApi.md#ActivateOneConfig) | **Put** /configs/{id}/activate | Activate one config
[**ChangeOneConfigSecret**](WebhooksApi.md#ChangeOneConfigSecret) | **Put** /configs/{id}/secret/change | Change the signing secret of a config
[**DeactivateOneConfig**](WebhooksApi.md#DeactivateOneConfig) | **Put** /configs/{id}/deactivate | Deactivate one config
[**DeleteOneConfig**](WebhooksApi.md#DeleteOneConfig) | **Delete** /configs/{id} | Delete one config
[**GetManyConfigs**](WebhooksApi.md#GetManyConfigs) | **Get** /configs | Get many configs
[**InsertOneConfig**](WebhooksApi.md#InsertOneConfig) | **Post** /configs | Insert a new config 
[**TestOneConfig**](WebhooksApi.md#TestOneConfig) | **Get** /configs/{id}/test | Test one config



## ActivateOneConfig

> ConfigResponse ActivateOneConfig(ctx, id).Execute()

Activate one config

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {
    id := "4997257d-dfb6-445b-929c-cbe2ab182818" // string | Config ID

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.WebhooksApi.ActivateOneConfig(context.Background(), id).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `WebhooksApi.ActivateOneConfig``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ActivateOneConfig`: ConfigResponse
    fmt.Fprintf(os.Stdout, "Response from `WebhooksApi.ActivateOneConfig`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | Config ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiActivateOneConfigRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**ConfigResponse**](ConfigResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ChangeOneConfigSecret

> ConfigResponse ChangeOneConfigSecret(ctx, id).ChangeOneConfigSecretRequest(changeOneConfigSecretRequest).Execute()

Change the signing secret of a config



### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {
    id := "4997257d-dfb6-445b-929c-cbe2ab182818" // string | Config ID
    changeOneConfigSecretRequest := *openapiclient.NewChangeOneConfigSecretRequest("V0bivxRWveaoz08afqjU6Ko/jwO0Cb+3") // ChangeOneConfigSecretRequest |  (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.WebhooksApi.ChangeOneConfigSecret(context.Background(), id).ChangeOneConfigSecretRequest(changeOneConfigSecretRequest).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `WebhooksApi.ChangeOneConfigSecret``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ChangeOneConfigSecret`: ConfigResponse
    fmt.Fprintf(os.Stdout, "Response from `WebhooksApi.ChangeOneConfigSecret`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | Config ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiChangeOneConfigSecretRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **changeOneConfigSecretRequest** | [**ChangeOneConfigSecretRequest**](ChangeOneConfigSecretRequest.md) |  | 

### Return type

[**ConfigResponse**](ConfigResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeactivateOneConfig

> ConfigResponse DeactivateOneConfig(ctx, id).Execute()

Deactivate one config

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {
    id := "4997257d-dfb6-445b-929c-cbe2ab182818" // string | Config ID

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.WebhooksApi.DeactivateOneConfig(context.Background(), id).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `WebhooksApi.DeactivateOneConfig``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `DeactivateOneConfig`: ConfigResponse
    fmt.Fprintf(os.Stdout, "Response from `WebhooksApi.DeactivateOneConfig`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | Config ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeactivateOneConfigRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**ConfigResponse**](ConfigResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteOneConfig

> DeleteOneConfig(ctx, id).Execute()

Delete one config

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {
    id := "4997257d-dfb6-445b-929c-cbe2ab182818" // string | Config ID

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.WebhooksApi.DeleteOneConfig(context.Background(), id).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `WebhooksApi.DeleteOneConfig``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | Config ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteOneConfigRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetManyConfigs

> GetManyConfigs200Response GetManyConfigs(ctx).Id(id).Endpoint(endpoint).Execute()

Get many configs



### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {
    id := "4997257d-dfb6-445b-929c-cbe2ab182818" // string | Optional filter by Config ID (optional)
    endpoint := "https://example.com" // string | Optional filter by endpoint URL (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.WebhooksApi.GetManyConfigs(context.Background()).Id(id).Endpoint(endpoint).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `WebhooksApi.GetManyConfigs``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `GetManyConfigs`: GetManyConfigs200Response
    fmt.Fprintf(os.Stdout, "Response from `WebhooksApi.GetManyConfigs`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiGetManyConfigsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **string** | Optional filter by Config ID | 
 **endpoint** | **string** | Optional filter by endpoint URL | 

### Return type

[**GetManyConfigs200Response**](GetManyConfigs200Response.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## InsertOneConfig

> ConfigResponse InsertOneConfig(ctx).ConfigUser(configUser).Execute()

Insert a new config 



### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {
    configUser := *openapiclient.NewConfigUser() // ConfigUser | 

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.WebhooksApi.InsertOneConfig(context.Background()).ConfigUser(configUser).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `WebhooksApi.InsertOneConfig``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `InsertOneConfig`: ConfigResponse
    fmt.Fprintf(os.Stdout, "Response from `WebhooksApi.InsertOneConfig`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiInsertOneConfigRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **configUser** | [**ConfigUser**](ConfigUser.md) |  | 

### Return type

[**ConfigResponse**](ConfigResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json, text/plain

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## TestOneConfig

> AttemptResponse TestOneConfig(ctx, id).Execute()

Test one config



### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {
    id := "4997257d-dfb6-445b-929c-cbe2ab182818" // string | Config ID

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.WebhooksApi.TestOneConfig(context.Background(), id).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `WebhooksApi.TestOneConfig``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `TestOneConfig`: AttemptResponse
    fmt.Fprintf(os.Stdout, "Response from `WebhooksApi.TestOneConfig`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | Config ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiTestOneConfigRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**AttemptResponse**](AttemptResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

