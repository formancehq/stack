# \PaymentsApi

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**ConnectorsStripeTransfer**](PaymentsApi.md#ConnectorsStripeTransfer) | **Post** /connectors/stripe/transfer | Transfer funds between Stripe accounts
[**GetAllConnectors**](PaymentsApi.md#GetAllConnectors) | **Get** /connectors | Get all installed connectors
[**GetAllConnectorsConfigs**](PaymentsApi.md#GetAllConnectorsConfigs) | **Get** /connectors/configs | Get all available connectors configs
[**GetConnectorTask**](PaymentsApi.md#GetConnectorTask) | **Get** /connectors/{connector}/tasks/{taskId} | Read a specific task of the connector
[**GetPayment**](PaymentsApi.md#GetPayment) | **Get** /payments/{paymentId} | Returns a payment.
[**InstallConnector**](PaymentsApi.md#InstallConnector) | **Post** /connectors/{connector} | Install connector
[**ListConnectorTasks**](PaymentsApi.md#ListConnectorTasks) | **Get** /connectors/{connector}/tasks | List connector tasks
[**ListPayments**](PaymentsApi.md#ListPayments) | **Get** /payments | Returns a list of payments.
[**ReadConnectorConfig**](PaymentsApi.md#ReadConnectorConfig) | **Get** /connectors/{connector}/config | Read connector config
[**ResetConnector**](PaymentsApi.md#ResetConnector) | **Post** /connectors/{connector}/reset | Reset connector
[**UninstallConnector**](PaymentsApi.md#UninstallConnector) | **Delete** /connectors/{connector} | Uninstall connector



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
    openapiclient "./openapi"
)

func main() {
    stripeTransferRequest := *openapiclient.NewStripeTransferRequest() // StripeTransferRequest |

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
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

No authorization required

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
    openapiclient "./openapi"
)

func main() {

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
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

No authorization required

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
    openapiclient "./openapi"
)

func main() {

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
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

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetConnectorTask

> interface{} GetConnectorTask(ctx, connector, taskId).Execute()

Read a specific task of the connector



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
    connector := *openapiclient.NewConnectors() // Connectors | The connector code
    taskId := TODO // interface{} | The task id

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.PaymentsApi.GetConnectorTask(context.Background(), connector, taskId).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `PaymentsApi.GetConnectorTask``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `GetConnectorTask`: interface{}
    fmt.Fprintf(os.Stdout, "Response from `PaymentsApi.GetConnectorTask`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**connector** | [**Connectors**](.md) | The connector code |
**taskId** | [**interface{}**](.md) | The task id |

### Other Parameters

Other parameters are passed through a pointer to a apiGetConnectorTaskRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

**interface{}**

### Authorization

No authorization required

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
    openapiclient "./openapi"
)

func main() {
    paymentId := TODO // interface{} | The payment id

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
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
**paymentId** | [**interface{}**](.md) | The payment id |

### Other Parameters

Other parameters are passed through a pointer to a apiGetPaymentRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**Payment**](Payment.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## InstallConnector

> InstallConnector(ctx, connector).Body(body).Execute()

Install connector



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
    connector := *openapiclient.NewConnectors() // Connectors | The connector code
    body := interface{}(987) // interface{} |

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.PaymentsApi.InstallConnector(context.Background(), connector).Body(body).Execute()
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

 **body** | **interface{}** |  |

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListConnectorTasks

> interface{} ListConnectorTasks(ctx, connector).Execute()

List connector tasks



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
    connector := *openapiclient.NewConnectors() // Connectors | The connector code

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.PaymentsApi.ListConnectorTasks(context.Background(), connector).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `PaymentsApi.ListConnectorTasks``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ListConnectorTasks`: interface{}
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

**interface{}**

### Authorization

No authorization required

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
    openapiclient "./openapi"
)

func main() {
    limit := TODO // interface{} | Limit the number of payments to return, pagination can be achieved in conjunction with 'skip' parameter. (optional)
    skip := TODO // interface{} | How many payments to skip, pagination can be achieved in conjunction with 'limit' parameter. (optional)
    sort := TODO // interface{} | Field used to sort payments (Default is by date). (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
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
 **limit** | [**interface{}**](interface{}.md) | Limit the number of payments to return, pagination can be achieved in conjunction with &#39;skip&#39; parameter. |
 **skip** | [**interface{}**](interface{}.md) | How many payments to skip, pagination can be achieved in conjunction with &#39;limit&#39; parameter. |
 **sort** | [**interface{}**](interface{}.md) | Field used to sort payments (Default is by date). |

### Return type

[**ListPaymentsResponse**](ListPaymentsResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ReadConnectorConfig

> interface{} ReadConnectorConfig(ctx, connector).Execute()

Read connector config



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
    connector := *openapiclient.NewConnectors() // Connectors | The connector code

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.PaymentsApi.ReadConnectorConfig(context.Background(), connector).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `PaymentsApi.ReadConnectorConfig``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ReadConnectorConfig`: interface{}
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

**interface{}**

### Authorization

No authorization required

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
    openapiclient "./openapi"
)

func main() {
    connector := *openapiclient.NewConnectors() // Connectors | The connector code

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
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

No authorization required

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
    openapiclient "./openapi"
)

func main() {
    connector := *openapiclient.NewConnectors() // Connectors | The connector code

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
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

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)
