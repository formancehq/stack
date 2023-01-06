# \PaymentsApi

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**ConnectorsStripeTransfer**](PaymentsApi.md#ConnectorsStripeTransfer) | **Post** /api/payments/connectors/stripe/transfer | Transfer funds between Stripe accounts
[**GetAllConnectors**](PaymentsApi.md#GetAllConnectors) | **Get** /api/payments/connectors | Get all installed connectors
[**GetAllConnectorsConfigs**](PaymentsApi.md#GetAllConnectorsConfigs) | **Get** /api/payments/connectors/configs | Get all available connectors configs
[**GetConnectorTask**](PaymentsApi.md#GetConnectorTask) | **Get** /api/payments/connectors/{connector}/tasks/{taskId} | Read a specific task of the connector
[**GetPayment**](PaymentsApi.md#GetPayment) | **Get** /api/payments/payments/{paymentId} | Returns a payment.
[**InstallConnector**](PaymentsApi.md#InstallConnector) | **Post** /api/payments/connectors/{connector} | Install connector
[**ListConnectorTasks**](PaymentsApi.md#ListConnectorTasks) | **Get** /api/payments/connectors/{connector}/tasks | List connector tasks
[**ListPayments**](PaymentsApi.md#ListPayments) | **Get** /api/payments/payments | Returns a list of payments.
[**ReadConnectorConfig**](PaymentsApi.md#ReadConnectorConfig) | **Get** /api/payments/connectors/{connector}/config | Read connector config
[**ResetConnector**](PaymentsApi.md#ResetConnector) | **Post** /api/payments/connectors/{connector}/reset | Reset connector
[**UninstallConnector**](PaymentsApi.md#UninstallConnector) | **Delete** /api/payments/connectors/{connector} | Uninstall connector



## ConnectorsStripeTransfer

> ConnectorsStripeTransfer(ctx).StripeTransferRequest(stripeTransferRequest).Execute()

Transfer funds between Stripe accounts



### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    client "./openapi"
)

func main() {
    stripeTransferRequest := *client.NewStripeTransferRequest() // StripeTransferRequest | 

    configuration := client.NewConfiguration()
    apiClient := client.NewAPIClient(configuration)
    resp, r, err := apiClient.PaymentsApi.ConnectorsStripeTransfer(context.Background()).StripeTransferRequest(stripeTransferRequest).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `PaymentsApi.ConnectorsStripeTransfer``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiConnectorsStripeTransferRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **stripeTransferRequest** | [**StripeTransferRequest**](StripeTransferRequest.md) |  | 

### Return type

 (empty response body)

### Authorization

[Authorization](../README.md#Authorization)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetAllConnectors

> ListConnectorsResponse GetAllConnectors(ctx).Execute()

Get all installed connectors



### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    client "./openapi"
)

func main() {

    configuration := client.NewConfiguration()
    apiClient := client.NewAPIClient(configuration)
    resp, r, err := apiClient.PaymentsApi.GetAllConnectors(context.Background()).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `PaymentsApi.GetAllConnectors``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `GetAllConnectors`: ListConnectorsResponse
    fmt.Fprintf(os.Stdout, "Response from `PaymentsApi.GetAllConnectors`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiGetAllConnectorsRequest struct via the builder pattern


### Return type

[**ListConnectorsResponse**](ListConnectorsResponse.md)

### Authorization

[Authorization](../README.md#Authorization)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetAllConnectorsConfigs

> ListConnectorsConfigsResponse GetAllConnectorsConfigs(ctx).Execute()

Get all available connectors configs



### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    client "./openapi"
)

func main() {

    configuration := client.NewConfiguration()
    apiClient := client.NewAPIClient(configuration)
    resp, r, err := apiClient.PaymentsApi.GetAllConnectorsConfigs(context.Background()).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `PaymentsApi.GetAllConnectorsConfigs``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `GetAllConnectorsConfigs`: ListConnectorsConfigsResponse
    fmt.Fprintf(os.Stdout, "Response from `PaymentsApi.GetAllConnectorsConfigs`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiGetAllConnectorsConfigsRequest struct via the builder pattern


### Return type

[**ListConnectorsConfigsResponse**](ListConnectorsConfigsResponse.md)

### Authorization

[Authorization](../README.md#Authorization)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetConnectorTask

> ListConnectorTasks200ResponseInner GetConnectorTask(ctx, connector, taskId).Execute()

Read a specific task of the connector



### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    client "./openapi"
)

func main() {
    connector := client.Connectors("STRIPE") // Connectors | The connector code
    taskId := "task1" // string | The task id

    configuration := client.NewConfiguration()
    apiClient := client.NewAPIClient(configuration)
    resp, r, err := apiClient.PaymentsApi.GetConnectorTask(context.Background(), connector, taskId).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `PaymentsApi.GetConnectorTask``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `GetConnectorTask`: ListConnectorTasks200ResponseInner
    fmt.Fprintf(os.Stdout, "Response from `PaymentsApi.GetConnectorTask`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**connector** | [**Connectors**](.md) | The connector code | 
**taskId** | **string** | The task id | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetConnectorTaskRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

[**ListConnectorTasks200ResponseInner**](ListConnectorTasks200ResponseInner.md)

### Authorization

[Authorization](../README.md#Authorization)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetPayment

> Payment GetPayment(ctx, paymentId).Execute()

Returns a payment.

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    client "./openapi"
)

func main() {
    paymentId := "XXX" // string | The payment id

    configuration := client.NewConfiguration()
    apiClient := client.NewAPIClient(configuration)
    resp, r, err := apiClient.PaymentsApi.GetPayment(context.Background(), paymentId).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `PaymentsApi.GetPayment``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `GetPayment`: Payment
    fmt.Fprintf(os.Stdout, "Response from `PaymentsApi.GetPayment`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**paymentId** | **string** | The payment id | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetPaymentRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**Payment**](Payment.md)

### Authorization

[Authorization](../README.md#Authorization)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## InstallConnector

> InstallConnector(ctx, connector).ConnectorConfig(connectorConfig).Execute()

Install connector



### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    client "./openapi"
)

func main() {
    connector := client.Connectors("STRIPE") // Connectors | The connector code
    connectorConfig := client.ConnectorConfig{BankingCircleConfig: client.NewBankingCircleConfig("XXX", "XXX", "XXX", "XXX")} // ConnectorConfig | 

    configuration := client.NewConfiguration()
    apiClient := client.NewAPIClient(configuration)
    resp, r, err := apiClient.PaymentsApi.InstallConnector(context.Background(), connector).ConnectorConfig(connectorConfig).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `PaymentsApi.InstallConnector``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**connector** | [**Connectors**](.md) | The connector code | 

### Other Parameters

Other parameters are passed through a pointer to a apiInstallConnectorRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **connectorConfig** | [**ConnectorConfig**](ConnectorConfig.md) |  | 

### Return type

 (empty response body)

### Authorization

[Authorization](../README.md#Authorization)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListConnectorTasks

> []ListConnectorTasks200ResponseInner ListConnectorTasks(ctx, connector).Execute()

List connector tasks



### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    client "./openapi"
)

func main() {
    connector := client.Connectors("STRIPE") // Connectors | The connector code

    configuration := client.NewConfiguration()
    apiClient := client.NewAPIClient(configuration)
    resp, r, err := apiClient.PaymentsApi.ListConnectorTasks(context.Background(), connector).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `PaymentsApi.ListConnectorTasks``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ListConnectorTasks`: []ListConnectorTasks200ResponseInner
    fmt.Fprintf(os.Stdout, "Response from `PaymentsApi.ListConnectorTasks`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**connector** | [**Connectors**](.md) | The connector code | 

### Other Parameters

Other parameters are passed through a pointer to a apiListConnectorTasksRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**[]ListConnectorTasks200ResponseInner**](ListConnectorTasks200ResponseInner.md)

### Authorization

[Authorization](../README.md#Authorization)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListPayments

> ListPaymentsResponse ListPayments(ctx).Limit(limit).Skip(skip).Sort(sort).Execute()

Returns a list of payments.

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    client "./openapi"
)

func main() {
    limit := int32(10) // int32 | Limit the number of payments to return, pagination can be achieved in conjunction with 'skip' parameter. (optional)
    skip := int32(100) // int32 | How many payments to skip, pagination can be achieved in conjunction with 'limit' parameter. (optional)
    sort := []string{"Inner_example"} // []string | Field used to sort payments (Default is by date). (optional)

    configuration := client.NewConfiguration()
    apiClient := client.NewAPIClient(configuration)
    resp, r, err := apiClient.PaymentsApi.ListPayments(context.Background()).Limit(limit).Skip(skip).Sort(sort).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `PaymentsApi.ListPayments``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ListPayments`: ListPaymentsResponse
    fmt.Fprintf(os.Stdout, "Response from `PaymentsApi.ListPayments`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListPaymentsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **limit** | **int32** | Limit the number of payments to return, pagination can be achieved in conjunction with &#39;skip&#39; parameter. | 
 **skip** | **int32** | How many payments to skip, pagination can be achieved in conjunction with &#39;limit&#39; parameter. | 
 **sort** | **[]string** | Field used to sort payments (Default is by date). | 

### Return type

[**ListPaymentsResponse**](ListPaymentsResponse.md)

### Authorization

[Authorization](../README.md#Authorization)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ReadConnectorConfig

> ConnectorConfig ReadConnectorConfig(ctx, connector).Execute()

Read connector config



### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    client "./openapi"
)

func main() {
    connector := client.Connectors("STRIPE") // Connectors | The connector code

    configuration := client.NewConfiguration()
    apiClient := client.NewAPIClient(configuration)
    resp, r, err := apiClient.PaymentsApi.ReadConnectorConfig(context.Background(), connector).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `PaymentsApi.ReadConnectorConfig``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ReadConnectorConfig`: ConnectorConfig
    fmt.Fprintf(os.Stdout, "Response from `PaymentsApi.ReadConnectorConfig`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**connector** | [**Connectors**](.md) | The connector code | 

### Other Parameters

Other parameters are passed through a pointer to a apiReadConnectorConfigRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**ConnectorConfig**](ConnectorConfig.md)

### Authorization

[Authorization](../README.md#Authorization)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ResetConnector

> ResetConnector(ctx, connector).Execute()

Reset connector



### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    client "./openapi"
)

func main() {
    connector := client.Connectors("STRIPE") // Connectors | The connector code

    configuration := client.NewConfiguration()
    apiClient := client.NewAPIClient(configuration)
    resp, r, err := apiClient.PaymentsApi.ResetConnector(context.Background(), connector).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `PaymentsApi.ResetConnector``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**connector** | [**Connectors**](.md) | The connector code | 

### Other Parameters

Other parameters are passed through a pointer to a apiResetConnectorRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

 (empty response body)

### Authorization

[Authorization](../README.md#Authorization)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UninstallConnector

> UninstallConnector(ctx, connector).Execute()

Uninstall connector



### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    client "./openapi"
)

func main() {
    connector := client.Connectors("STRIPE") // Connectors | The connector code

    configuration := client.NewConfiguration()
    apiClient := client.NewAPIClient(configuration)
    resp, r, err := apiClient.PaymentsApi.UninstallConnector(context.Background(), connector).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `PaymentsApi.UninstallConnector``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**connector** | [**Connectors**](.md) | The connector code | 

### Other Parameters

Other parameters are passed through a pointer to a apiUninstallConnectorRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

 (empty response body)

### Authorization

[Authorization](../README.md#Authorization)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

