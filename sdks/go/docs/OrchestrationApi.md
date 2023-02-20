# \OrchestrationApi

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CancelEvent**](OrchestrationApi.md#CancelEvent) | **Put** /api/orchestration/instances/{instanceID}/abort | Cancel a running workflow
[**CreateWorkflow**](OrchestrationApi.md#CreateWorkflow) | **Post** /api/orchestration/workflows | Create workflow
[**GetInstance**](OrchestrationApi.md#GetInstance) | **Get** /api/orchestration/instances/{instanceID} | Get a workflow instance by id
[**GetInstanceHistory**](OrchestrationApi.md#GetInstanceHistory) | **Get** /api/orchestration/instances/{instanceID}/history | Get a workflow instance history by id
[**GetInstanceStageHistory**](OrchestrationApi.md#GetInstanceStageHistory) | **Get** /api/orchestration/instances/{instanceID}/stages/{number}/history | Get a workflow instance stage history
[**GetWorkflow**](OrchestrationApi.md#GetWorkflow) | **Get** /api/orchestration/workflows/{flowId} | Get a flow by id
[**ListInstances**](OrchestrationApi.md#ListInstances) | **Get** /api/orchestration/instances | List instances of a workflow
[**ListWorkflows**](OrchestrationApi.md#ListWorkflows) | **Get** /api/orchestration/workflows | List registered workflows
[**OrchestrationgetServerInfo**](OrchestrationApi.md#OrchestrationgetServerInfo) | **Get** /api/orchestration/_info | Get server info
[**RunWorkflow**](OrchestrationApi.md#RunWorkflow) | **Post** /api/orchestration/workflows/{workflowID}/instances | Run workflow
[**SendEvent**](OrchestrationApi.md#SendEvent) | **Post** /api/orchestration/instances/{instanceID}/events | Send an event to a running workflow



## CancelEvent

> CancelEvent(ctx, instanceID).Execute()

Cancel a running workflow



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
    instanceID := "xxx" // string | The instance id

    configuration := client.NewConfiguration()
    apiClient := client.NewAPIClient(configuration)
    r, err := apiClient.OrchestrationApi.CancelEvent(context.Background(), instanceID).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `OrchestrationApi.CancelEvent``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**instanceID** | **string** | The instance id | 

### Other Parameters

Other parameters are passed through a pointer to a apiCancelEventRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

 (empty response body)

### Authorization

[Authorization](../README.md#Authorization)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## CreateWorkflow

> CreateWorkflowResponse CreateWorkflow(ctx).Body(body).Execute()

Create workflow



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
    body := WorkflowConfig(987) // WorkflowConfig |  (optional)

    configuration := client.NewConfiguration()
    apiClient := client.NewAPIClient(configuration)
    resp, r, err := apiClient.OrchestrationApi.CreateWorkflow(context.Background()).Body(body).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `OrchestrationApi.CreateWorkflow``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `CreateWorkflow`: CreateWorkflowResponse
    fmt.Fprintf(os.Stdout, "Response from `OrchestrationApi.CreateWorkflow`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiCreateWorkflowRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | **WorkflowConfig** |  | 

### Return type

[**CreateWorkflowResponse**](CreateWorkflowResponse.md)

### Authorization

[Authorization](../README.md#Authorization)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetInstance

> GetWorkflowInstanceResponse GetInstance(ctx, instanceID).Execute()

Get a workflow instance by id



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
    instanceID := "xxx" // string | The instance id

    configuration := client.NewConfiguration()
    apiClient := client.NewAPIClient(configuration)
    resp, r, err := apiClient.OrchestrationApi.GetInstance(context.Background(), instanceID).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `OrchestrationApi.GetInstance``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `GetInstance`: GetWorkflowInstanceResponse
    fmt.Fprintf(os.Stdout, "Response from `OrchestrationApi.GetInstance`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**instanceID** | **string** | The instance id | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetInstanceRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**GetWorkflowInstanceResponse**](GetWorkflowInstanceResponse.md)

### Authorization

[Authorization](../README.md#Authorization)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetInstanceHistory

> GetWorkflowInstanceHistoryResponse GetInstanceHistory(ctx, instanceID).Execute()

Get a workflow instance history by id



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
    instanceID := "xxx" // string | The instance id

    configuration := client.NewConfiguration()
    apiClient := client.NewAPIClient(configuration)
    resp, r, err := apiClient.OrchestrationApi.GetInstanceHistory(context.Background(), instanceID).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `OrchestrationApi.GetInstanceHistory``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `GetInstanceHistory`: GetWorkflowInstanceHistoryResponse
    fmt.Fprintf(os.Stdout, "Response from `OrchestrationApi.GetInstanceHistory`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**instanceID** | **string** | The instance id | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetInstanceHistoryRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**GetWorkflowInstanceHistoryResponse**](GetWorkflowInstanceHistoryResponse.md)

### Authorization

[Authorization](../README.md#Authorization)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetInstanceStageHistory

> GetWorkflowInstanceHistoryStageResponse GetInstanceStageHistory(ctx, instanceID, number).Execute()

Get a workflow instance stage history



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
    instanceID := "xxx" // string | The instance id
    number := int32(0) // int32 | The stage number

    configuration := client.NewConfiguration()
    apiClient := client.NewAPIClient(configuration)
    resp, r, err := apiClient.OrchestrationApi.GetInstanceStageHistory(context.Background(), instanceID, number).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `OrchestrationApi.GetInstanceStageHistory``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `GetInstanceStageHistory`: GetWorkflowInstanceHistoryStageResponse
    fmt.Fprintf(os.Stdout, "Response from `OrchestrationApi.GetInstanceStageHistory`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**instanceID** | **string** | The instance id | 
**number** | **int32** | The stage number | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetInstanceStageHistoryRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

[**GetWorkflowInstanceHistoryStageResponse**](GetWorkflowInstanceHistoryStageResponse.md)

### Authorization

[Authorization](../README.md#Authorization)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetWorkflow

> GetWorkflowResponse GetWorkflow(ctx, flowId).Execute()

Get a flow by id



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
    flowId := "xxx" // string | The flow id

    configuration := client.NewConfiguration()
    apiClient := client.NewAPIClient(configuration)
    resp, r, err := apiClient.OrchestrationApi.GetWorkflow(context.Background(), flowId).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `OrchestrationApi.GetWorkflow``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `GetWorkflow`: GetWorkflowResponse
    fmt.Fprintf(os.Stdout, "Response from `OrchestrationApi.GetWorkflow`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**flowId** | **string** | The flow id | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetWorkflowRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**GetWorkflowResponse**](GetWorkflowResponse.md)

### Authorization

[Authorization](../README.md#Authorization)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListInstances

> ListRunsResponse ListInstances(ctx).WorkflowID(workflowID).Running(running).Execute()

List instances of a workflow



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
    workflowID := "xxx" // string | A workflow id (optional)
    running := false // bool | Filter running instances (optional)

    configuration := client.NewConfiguration()
    apiClient := client.NewAPIClient(configuration)
    resp, r, err := apiClient.OrchestrationApi.ListInstances(context.Background()).WorkflowID(workflowID).Running(running).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `OrchestrationApi.ListInstances``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ListInstances`: ListRunsResponse
    fmt.Fprintf(os.Stdout, "Response from `OrchestrationApi.ListInstances`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListInstancesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **workflowID** | **string** | A workflow id | 
 **running** | **bool** | Filter running instances | 

### Return type

[**ListRunsResponse**](ListRunsResponse.md)

### Authorization

[Authorization](../README.md#Authorization)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListWorkflows

> ListWorkflowsResponse ListWorkflows(ctx).Execute()

List registered workflows



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
    resp, r, err := apiClient.OrchestrationApi.ListWorkflows(context.Background()).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `OrchestrationApi.ListWorkflows``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ListWorkflows`: ListWorkflowsResponse
    fmt.Fprintf(os.Stdout, "Response from `OrchestrationApi.ListWorkflows`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiListWorkflowsRequest struct via the builder pattern


### Return type

[**ListWorkflowsResponse**](ListWorkflowsResponse.md)

### Authorization

[Authorization](../README.md#Authorization)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## OrchestrationgetServerInfo

> ServerInfo OrchestrationgetServerInfo(ctx).Execute()

Get server info

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
    resp, r, err := apiClient.OrchestrationApi.OrchestrationgetServerInfo(context.Background()).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `OrchestrationApi.OrchestrationgetServerInfo``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `OrchestrationgetServerInfo`: ServerInfo
    fmt.Fprintf(os.Stdout, "Response from `OrchestrationApi.OrchestrationgetServerInfo`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiOrchestrationgetServerInfoRequest struct via the builder pattern


### Return type

[**ServerInfo**](ServerInfo.md)

### Authorization

[Authorization](../README.md#Authorization)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## RunWorkflow

> RunWorkflowResponse RunWorkflow(ctx, workflowID).Wait(wait).RequestBody(requestBody).Execute()

Run workflow



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
    workflowID := "xxx" // string | The flow id
    wait := true // bool | Wait end of the workflow before return (optional)
    requestBody := map[string]string{"key": "Inner_example"} // map[string]string |  (optional)

    configuration := client.NewConfiguration()
    apiClient := client.NewAPIClient(configuration)
    resp, r, err := apiClient.OrchestrationApi.RunWorkflow(context.Background(), workflowID).Wait(wait).RequestBody(requestBody).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `OrchestrationApi.RunWorkflow``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `RunWorkflow`: RunWorkflowResponse
    fmt.Fprintf(os.Stdout, "Response from `OrchestrationApi.RunWorkflow`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**workflowID** | **string** | The flow id | 

### Other Parameters

Other parameters are passed through a pointer to a apiRunWorkflowRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **wait** | **bool** | Wait end of the workflow before return | 
 **requestBody** | **map[string]string** |  | 

### Return type

[**RunWorkflowResponse**](RunWorkflowResponse.md)

### Authorization

[Authorization](../README.md#Authorization)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## SendEvent

> SendEvent(ctx, instanceID).SendEventRequest(sendEventRequest).Execute()

Send an event to a running workflow



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
    instanceID := "xxx" // string | The instance id
    sendEventRequest := *client.NewSendEventRequest("Name_example") // SendEventRequest |  (optional)

    configuration := client.NewConfiguration()
    apiClient := client.NewAPIClient(configuration)
    r, err := apiClient.OrchestrationApi.SendEvent(context.Background(), instanceID).SendEventRequest(sendEventRequest).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `OrchestrationApi.SendEvent``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**instanceID** | **string** | The instance id | 

### Other Parameters

Other parameters are passed through a pointer to a apiSendEventRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **sendEventRequest** | [**SendEventRequest**](SendEventRequest.md) |  | 

### Return type

 (empty response body)

### Authorization

[Authorization](../README.md#Authorization)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

