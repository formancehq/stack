# \ConnectorsApi

All URIs are relative to *http://localhost:8080*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CreateConnector**](ConnectorsApi.md#CreateConnector) | **Post** /connectors | Create connector
[**DeleteConnector**](ConnectorsApi.md#DeleteConnector) | **Delete** /connectors/{connectorID} | Delete connector
[**GetConnectorState**](ConnectorsApi.md#GetConnectorState) | **Get** /connectors/{connectorID} | Get connector state
[**ListConnectors**](ConnectorsApi.md#ListConnectors) | **Get** /connectors | List connectors



## CreateConnector

> CreateConnector201Response CreateConnector(ctx).Body(body).Execute()

Create connector

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
    body := ConnectorConfiguration(987) // ConnectorConfiguration |  (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.ConnectorsApi.CreateConnector(context.Background()).Body(body).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `ConnectorsApi.CreateConnector``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `CreateConnector`: CreateConnector201Response
    fmt.Fprintf(os.Stdout, "Response from `ConnectorsApi.CreateConnector`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiCreateConnectorRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | **ConnectorConfiguration** |  | 

### Return type

[**CreateConnector201Response**](CreateConnector201Response.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteConnector

> DeleteConnector(ctx, connectorID).Execute()

Delete connector

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
    connectorID := "connectorID_example" // string | The connector id

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    r, err := apiClient.ConnectorsApi.DeleteConnector(context.Background(), connectorID).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `ConnectorsApi.DeleteConnector``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**connectorID** | **string** | The connector id | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteConnectorRequest struct via the builder pattern


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


## GetConnectorState

> CreateConnector201Response GetConnectorState(ctx, connectorID).Execute()

Get connector state

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
    connectorID := "connectorID_example" // string | The connector id

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.ConnectorsApi.GetConnectorState(context.Background(), connectorID).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `ConnectorsApi.GetConnectorState``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `GetConnectorState`: CreateConnector201Response
    fmt.Fprintf(os.Stdout, "Response from `ConnectorsApi.GetConnectorState`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**connectorID** | **string** | The connector id | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetConnectorStateRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**CreateConnector201Response**](CreateConnector201Response.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListConnectors

> ListConnectors200Response ListConnectors(ctx).Execute()

List connectors

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
    resp, r, err := apiClient.ConnectorsApi.ListConnectors(context.Background()).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `ConnectorsApi.ListConnectors``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ListConnectors`: ListConnectors200Response
    fmt.Fprintf(os.Stdout, "Response from `ConnectorsApi.ListConnectors`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiListConnectorsRequest struct via the builder pattern


### Return type

[**ListConnectors200Response**](ListConnectors200Response.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

