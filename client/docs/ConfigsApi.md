# \ConfigsApi

All URIs are relative to *https://.o.formance.cloud/webhooks*

Method | HTTP request | Description
------------- | ------------- | -------------
[**ActivateOneConfig**](ConfigsApi.md#ActivateOneConfig) | **Put** /configs/{id}/activate | Activate one config
[**ChangeOneConfigSecret**](ConfigsApi.md#ChangeOneConfigSecret) | **Put** /configs/{id}/secret/change | Change the signing secret of a config
[**DeactivateOneConfig**](ConfigsApi.md#DeactivateOneConfig) | **Put** /configs/{id}/deactivate | Deactivate one config
[**DeleteOneConfig**](ConfigsApi.md#DeleteOneConfig) | **Delete** /configs/{id} | Delete one config
[**GetManyConfigs**](ConfigsApi.md#GetManyConfigs) | **Get** /configs | Get many configs
[**InsertOneConfig**](ConfigsApi.md#InsertOneConfig) | **Post** /configs | Insert a new config 



## ActivateOneConfig

> GetManyConfigs200Response ActivateOneConfig(ctx, id).Execute()

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
    resp, r, err := apiClient.ConfigsApi.ActivateOneConfig(context.Background(), id).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `ConfigsApi.ActivateOneConfig``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ActivateOneConfig`: GetManyConfigs200Response
    fmt.Fprintf(os.Stdout, "Response from `ConfigsApi.ActivateOneConfig`: %v\n", resp)
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

[**GetManyConfigs200Response**](GetManyConfigs200Response.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ChangeOneConfigSecret

> GetManyConfigs200Response ChangeOneConfigSecret(ctx, id).ChangeOneConfigSecretRequest(changeOneConfigSecretRequest).Execute()

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
    changeOneConfigSecretRequest := *openapiclient.NewChangeOneConfigSecretRequest() // ChangeOneConfigSecretRequest |  (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.ConfigsApi.ChangeOneConfigSecret(context.Background(), id).ChangeOneConfigSecretRequest(changeOneConfigSecretRequest).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `ConfigsApi.ChangeOneConfigSecret``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ChangeOneConfigSecret`: GetManyConfigs200Response
    fmt.Fprintf(os.Stdout, "Response from `ConfigsApi.ChangeOneConfigSecret`: %v\n", resp)
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

[**GetManyConfigs200Response**](GetManyConfigs200Response.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeactivateOneConfig

> GetManyConfigs200Response DeactivateOneConfig(ctx, id).Execute()

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
    resp, r, err := apiClient.ConfigsApi.DeactivateOneConfig(context.Background(), id).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `ConfigsApi.DeactivateOneConfig``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `DeactivateOneConfig`: GetManyConfigs200Response
    fmt.Fprintf(os.Stdout, "Response from `ConfigsApi.DeactivateOneConfig`: %v\n", resp)
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

[**GetManyConfigs200Response**](GetManyConfigs200Response.md)

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
    resp, r, err := apiClient.ConfigsApi.DeleteOneConfig(context.Background(), id).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `ConfigsApi.DeleteOneConfig``: %v\n", err)
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
    resp, r, err := apiClient.ConfigsApi.GetManyConfigs(context.Background()).Id(id).Endpoint(endpoint).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `ConfigsApi.GetManyConfigs``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `GetManyConfigs`: GetManyConfigs200Response
    fmt.Fprintf(os.Stdout, "Response from `ConfigsApi.GetManyConfigs`: %v\n", resp)
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

> string InsertOneConfig(ctx).ConfigUser(configUser).Execute()

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
    resp, r, err := apiClient.ConfigsApi.InsertOneConfig(context.Background()).ConfigUser(configUser).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `ConfigsApi.InsertOneConfig``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `InsertOneConfig`: string
    fmt.Fprintf(os.Stdout, "Response from `ConfigsApi.InsertOneConfig`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiInsertOneConfigRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **configUser** | [**ConfigUser**](ConfigUser.md) |  | 

### Return type

**string**

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: text/plain

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

