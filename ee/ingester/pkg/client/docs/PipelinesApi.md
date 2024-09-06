# \PipelinesApi

All URIs are relative to *http://localhost:8080*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CreatePipeline**](PipelinesApi.md#CreatePipeline) | **Post** /pipelines | Create pipeline
[**DeletePipeline**](PipelinesApi.md#DeletePipeline) | **Delete** /pipelines/{pipelineID} | Delete pipeline
[**GetPipelineState**](PipelinesApi.md#GetPipelineState) | **Get** /pipelines/{pipelineID} | Get pipeline state
[**ListPipelines**](PipelinesApi.md#ListPipelines) | **Get** /pipelines | List pipelines
[**PausePipeline**](PipelinesApi.md#PausePipeline) | **Post** /pipelines/{pipelineID}/pause | Pause pipeline
[**ResetPipeline**](PipelinesApi.md#ResetPipeline) | **Post** /pipelines/{pipelineID}/reset | Reset pipeline
[**ResumePipeline**](PipelinesApi.md#ResumePipeline) | **Post** /pipelines/{pipelineID}/resume | Resume pipeline
[**StartPipeline**](PipelinesApi.md#StartPipeline) | **Post** /pipelines/{pipelineID}/start | Start pipeline
[**StopPipeline**](PipelinesApi.md#StopPipeline) | **Post** /pipelines/{pipelineID}/stop | Stop pipeline



## CreatePipeline

> CreatePipeline201Response CreatePipeline(ctx).CreatePipelineRequest(createPipelineRequest).Execute()

Create pipeline

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "github.com/formancehq/ingester/ingesterclient"
)

func main() {
    createPipelineRequest := *openapiclient.NewCreatePipelineRequest("Module_example", "ConnectorID_example") // CreatePipelineRequest |  (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.PipelinesApi.CreatePipeline(context.Background()).CreatePipelineRequest(createPipelineRequest).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `PipelinesApi.CreatePipeline``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `CreatePipeline`: CreatePipeline201Response
    fmt.Fprintf(os.Stdout, "Response from `PipelinesApi.CreatePipeline`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiCreatePipelineRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **createPipelineRequest** | [**CreatePipelineRequest**](CreatePipelineRequest.md) |  | 

### Return type

[**CreatePipeline201Response**](CreatePipeline201Response.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeletePipeline

> DeletePipeline(ctx, pipelineID).Execute()

Delete pipeline

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "github.com/formancehq/ingester/ingesterclient"
)

func main() {
    pipelineID := "pipelineID_example" // string | The pipeline id

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    r, err := apiClient.PipelinesApi.DeletePipeline(context.Background(), pipelineID).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `PipelinesApi.DeletePipeline``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**pipelineID** | **string** | The pipeline id | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeletePipelineRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetPipelineState

> CreatePipeline201Response GetPipelineState(ctx, pipelineID).Execute()

Get pipeline state

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "github.com/formancehq/ingester/ingesterclient"
)

func main() {
    pipelineID := "pipelineID_example" // string | The pipeline id

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.PipelinesApi.GetPipelineState(context.Background(), pipelineID).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `PipelinesApi.GetPipelineState``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `GetPipelineState`: CreatePipeline201Response
    fmt.Fprintf(os.Stdout, "Response from `PipelinesApi.GetPipelineState`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**pipelineID** | **string** | The pipeline id | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetPipelineStateRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**CreatePipeline201Response**](CreatePipeline201Response.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListPipelines

> ListPipelines200Response ListPipelines(ctx).Execute()

List pipelines

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "github.com/formancehq/ingester/ingesterclient"
)

func main() {

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.PipelinesApi.ListPipelines(context.Background()).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `PipelinesApi.ListPipelines``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ListPipelines`: ListPipelines200Response
    fmt.Fprintf(os.Stdout, "Response from `PipelinesApi.ListPipelines`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiListPipelinesRequest struct via the builder pattern


### Return type

[**ListPipelines200Response**](ListPipelines200Response.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## PausePipeline

> PausePipeline(ctx, pipelineID).Execute()

Pause pipeline

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "github.com/formancehq/ingester/ingesterclient"
)

func main() {
    pipelineID := "pipelineID_example" // string | The pipeline id

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    r, err := apiClient.PipelinesApi.PausePipeline(context.Background(), pipelineID).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `PipelinesApi.PausePipeline``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**pipelineID** | **string** | The pipeline id | 

### Other Parameters

Other parameters are passed through a pointer to a apiPausePipelineRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ResetPipeline

> ResetPipeline(ctx, pipelineID).Execute()

Reset pipeline

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "github.com/formancehq/ingester/ingesterclient"
)

func main() {
    pipelineID := "pipelineID_example" // string | The pipeline id

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    r, err := apiClient.PipelinesApi.ResetPipeline(context.Background(), pipelineID).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `PipelinesApi.ResetPipeline``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**pipelineID** | **string** | The pipeline id | 

### Other Parameters

Other parameters are passed through a pointer to a apiResetPipelineRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ResumePipeline

> ResumePipeline(ctx, pipelineID).Execute()

Resume pipeline

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "github.com/formancehq/ingester/ingesterclient"
)

func main() {
    pipelineID := "pipelineID_example" // string | The pipeline id

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    r, err := apiClient.PipelinesApi.ResumePipeline(context.Background(), pipelineID).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `PipelinesApi.ResumePipeline``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**pipelineID** | **string** | The pipeline id | 

### Other Parameters

Other parameters are passed through a pointer to a apiResumePipelineRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## StartPipeline

> StartPipeline(ctx, pipelineID).Execute()

Start pipeline

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "github.com/formancehq/ingester/ingesterclient"
)

func main() {
    pipelineID := "pipelineID_example" // string | The pipeline id

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    r, err := apiClient.PipelinesApi.StartPipeline(context.Background(), pipelineID).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `PipelinesApi.StartPipeline``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**pipelineID** | **string** | The pipeline id | 

### Other Parameters

Other parameters are passed through a pointer to a apiStartPipelineRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## StopPipeline

> StopPipeline(ctx, pipelineID).Execute()

Stop pipeline

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "github.com/formancehq/ingester/ingesterclient"
)

func main() {
    pipelineID := "pipelineID_example" // string | The pipeline id

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    r, err := apiClient.PipelinesApi.StopPipeline(context.Background(), pipelineID).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `PipelinesApi.StopPipeline``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**pipelineID** | **string** | The pipeline id | 

### Other Parameters

Other parameters are passed through a pointer to a apiStopPipelineRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

