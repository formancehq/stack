# \TransactionsApi

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**AddMetadataOnTransaction**](TransactionsApi.md#AddMetadataOnTransaction) | **Post** /api/ledger/{ledger}/transactions/{txid}/metadata | Set the metadata of a transaction by its ID
[**CountTransactions**](TransactionsApi.md#CountTransactions) | **Head** /api/ledger/{ledger}/transactions | Count the transactions from a ledger
[**CreateTransaction**](TransactionsApi.md#CreateTransaction) | **Post** /api/ledger/{ledger}/transactions | Create a new transaction to a ledger
[**GetTransaction**](TransactionsApi.md#GetTransaction) | **Get** /api/ledger/{ledger}/transactions/{txid} | Get transaction from a ledger by its ID
[**ListTransactions**](TransactionsApi.md#ListTransactions) | **Get** /api/ledger/{ledger}/transactions | List transactions from a ledger
[**RevertTransaction**](TransactionsApi.md#RevertTransaction) | **Post** /api/ledger/{ledger}/transactions/{txid}/revert | Revert a ledger transaction by its ID



## AddMetadataOnTransaction

> AddMetadataOnTransaction(ctx, ledger, txid).RequestBody(requestBody).Execute()

Set the metadata of a transaction by its ID

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
    ledger := "ledger001" // string | Name of the ledger.
    txid := int64(1234) // int64 | Transaction ID.
    requestBody := map[string]string{"key": "Inner_example"} // map[string]string | metadata (optional)

    configuration := client.NewConfiguration()
    apiClient := client.NewAPIClient(configuration)
    r, err := apiClient.TransactionsApi.AddMetadataOnTransaction(context.Background(), ledger, txid).RequestBody(requestBody).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `TransactionsApi.AddMetadataOnTransaction``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**ledger** | **string** | Name of the ledger. | 
**txid** | **int64** | Transaction ID. | 

### Other Parameters

Other parameters are passed through a pointer to a apiAddMetadataOnTransactionRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **requestBody** | **map[string]string** | metadata | 

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


## CountTransactions

> CountTransactions(ctx, ledger).Reference(reference).Account(account).Source(source).Destination(destination).StartTime(startTime).EndTime(endTime).Metadata(metadata).Execute()

Count the transactions from a ledger

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    "time"
    client "github.com/formancehq/formance-sdk-go"
)

func main() {
    ledger := "ledger001" // string | Name of the ledger.
    reference := "ref:001" // string | Filter transactions by reference field. (optional)
    account := "users:001" // string | Filter transactions with postings involving given account, either as source or destination (regular expression placed between ^ and $). (optional)
    source := "users:001" // string | Filter transactions with postings involving given account at source (regular expression placed between ^ and $). (optional)
    destination := "users:001" // string | Filter transactions with postings involving given account at destination (regular expression placed between ^ and $). (optional)
    startTime := time.Now() // time.Time | Filter transactions that occurred after this timestamp. The format is RFC3339 and is inclusive (for example, \"2023-01-02T15:04:01Z\" includes the first second of 4th minute).  (optional)
    endTime := time.Now() // time.Time | Filter transactions that occurred before this timestamp. The format is RFC3339 and is exclusive (for example, \"2023-01-02T15:04:01Z\" excludes the first second of 4th minute).  (optional)
    metadata := map[string]interface{}{"key": map[string]interface{}(123)} // map[string]interface{} | Filter transactions by metadata key value pairs. Nested objects can be used as seen in the example below. (optional)

    configuration := client.NewConfiguration()
    apiClient := client.NewAPIClient(configuration)
    r, err := apiClient.TransactionsApi.CountTransactions(context.Background(), ledger).Reference(reference).Account(account).Source(source).Destination(destination).StartTime(startTime).EndTime(endTime).Metadata(metadata).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `TransactionsApi.CountTransactions``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**ledger** | **string** | Name of the ledger. | 

### Other Parameters

Other parameters are passed through a pointer to a apiCountTransactionsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **reference** | **string** | Filter transactions by reference field. | 
 **account** | **string** | Filter transactions with postings involving given account, either as source or destination (regular expression placed between ^ and $). | 
 **source** | **string** | Filter transactions with postings involving given account at source (regular expression placed between ^ and $). | 
 **destination** | **string** | Filter transactions with postings involving given account at destination (regular expression placed between ^ and $). | 
 **startTime** | **time.Time** | Filter transactions that occurred after this timestamp. The format is RFC3339 and is inclusive (for example, \&quot;2023-01-02T15:04:01Z\&quot; includes the first second of 4th minute).  | 
 **endTime** | **time.Time** | Filter transactions that occurred before this timestamp. The format is RFC3339 and is exclusive (for example, \&quot;2023-01-02T15:04:01Z\&quot; excludes the first second of 4th minute).  | 
 **metadata** | [**map[string]interface{}**](map[string]interface{}.md) | Filter transactions by metadata key value pairs. Nested objects can be used as seen in the example below. | 

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


## CreateTransaction

> TransactionResponse CreateTransaction(ctx, ledger).PostTransaction(postTransaction).Preview(preview).Execute()

Create a new transaction to a ledger

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
    ledger := "ledger001" // string | Name of the ledger.
    postTransaction := *client.NewPostTransaction(map[string]string{"key": "Inner_example"}) // PostTransaction | The request body must contain at least one of the following objects:   - `postings`: suitable for simple transactions   - `script`: enabling more complex transactions with Numscript 
    preview := true // bool | Set the preview mode. Preview mode doesn't add the logs to the database or publish a message to the message broker. (optional)

    configuration := client.NewConfiguration()
    apiClient := client.NewAPIClient(configuration)
    resp, r, err := apiClient.TransactionsApi.CreateTransaction(context.Background(), ledger).PostTransaction(postTransaction).Preview(preview).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `TransactionsApi.CreateTransaction``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `CreateTransaction`: TransactionResponse
    fmt.Fprintf(os.Stdout, "Response from `TransactionsApi.CreateTransaction`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**ledger** | **string** | Name of the ledger. | 

### Other Parameters

Other parameters are passed through a pointer to a apiCreateTransactionRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **postTransaction** | [**PostTransaction**](PostTransaction.md) | The request body must contain at least one of the following objects:   - &#x60;postings&#x60;: suitable for simple transactions   - &#x60;script&#x60;: enabling more complex transactions with Numscript  | 
 **preview** | **bool** | Set the preview mode. Preview mode doesn&#39;t add the logs to the database or publish a message to the message broker. | 

### Return type

[**TransactionResponse**](TransactionResponse.md)

### Authorization

[Authorization](../README.md#Authorization)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetTransaction

> TransactionResponse GetTransaction(ctx, ledger, txid).Execute()

Get transaction from a ledger by its ID

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
    ledger := "ledger001" // string | Name of the ledger.
    txid := int64(1234) // int64 | Transaction ID.

    configuration := client.NewConfiguration()
    apiClient := client.NewAPIClient(configuration)
    resp, r, err := apiClient.TransactionsApi.GetTransaction(context.Background(), ledger, txid).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `TransactionsApi.GetTransaction``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `GetTransaction`: TransactionResponse
    fmt.Fprintf(os.Stdout, "Response from `TransactionsApi.GetTransaction`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**ledger** | **string** | Name of the ledger. | 
**txid** | **int64** | Transaction ID. | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetTransactionRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

[**TransactionResponse**](TransactionResponse.md)

### Authorization

[Authorization](../README.md#Authorization)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListTransactions

> TransactionsCursorResponse ListTransactions(ctx, ledger).PageSize(pageSize).After(after).Reference(reference).Account(account).Source(source).Destination(destination).StartTime(startTime).EndTime(endTime).Cursor(cursor).Metadata(metadata).Execute()

List transactions from a ledger



### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    "time"
    client "github.com/formancehq/formance-sdk-go"
)

func main() {
    ledger := "ledger001" // string | Name of the ledger.
    pageSize := int64(100) // int64 | The maximum number of results to return per page.  (optional) (default to 15)
    after := "1234" // string | Pagination cursor, will return transactions after given txid (in descending order). (optional)
    reference := "ref:001" // string | Find transactions by reference field. (optional)
    account := "users:001" // string | Filter transactions with postings involving given account, either as source or destination (regular expression placed between ^ and $). (optional)
    source := "users:001" // string | Filter transactions with postings involving given account at source (regular expression placed between ^ and $). (optional)
    destination := "users:001" // string | Filter transactions with postings involving given account at destination (regular expression placed between ^ and $). (optional)
    startTime := time.Now() // time.Time | Filter transactions that occurred after this timestamp. The format is RFC3339 and is inclusive (for example, \"2023-01-02T15:04:01Z\" includes the first second of 4th minute).  (optional)
    endTime := time.Now() // time.Time | Filter transactions that occurred before this timestamp. The format is RFC3339 and is exclusive (for example, \"2023-01-02T15:04:01Z\" excludes the first second of 4th minute).  (optional)
    cursor := "aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ==" // string | Parameter used in pagination requests. Maximum page size is set to 15. Set to the value of next for the next page of results. Set to the value of previous for the previous page of results. No other parameters can be set when this parameter is set.  (optional)
    metadata := map[string]string{"key": map[string]string{"key": "Inner_example"}} // map[string]string | Filter transactions by metadata key value pairs. (optional)

    configuration := client.NewConfiguration()
    apiClient := client.NewAPIClient(configuration)
    resp, r, err := apiClient.TransactionsApi.ListTransactions(context.Background(), ledger).PageSize(pageSize).After(after).Reference(reference).Account(account).Source(source).Destination(destination).StartTime(startTime).EndTime(endTime).Cursor(cursor).Metadata(metadata).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `TransactionsApi.ListTransactions``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ListTransactions`: TransactionsCursorResponse
    fmt.Fprintf(os.Stdout, "Response from `TransactionsApi.ListTransactions`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**ledger** | **string** | Name of the ledger. | 

### Other Parameters

Other parameters are passed through a pointer to a apiListTransactionsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **pageSize** | **int64** | The maximum number of results to return per page.  | [default to 15]
 **after** | **string** | Pagination cursor, will return transactions after given txid (in descending order). | 
 **reference** | **string** | Find transactions by reference field. | 
 **account** | **string** | Filter transactions with postings involving given account, either as source or destination (regular expression placed between ^ and $). | 
 **source** | **string** | Filter transactions with postings involving given account at source (regular expression placed between ^ and $). | 
 **destination** | **string** | Filter transactions with postings involving given account at destination (regular expression placed between ^ and $). | 
 **startTime** | **time.Time** | Filter transactions that occurred after this timestamp. The format is RFC3339 and is inclusive (for example, \&quot;2023-01-02T15:04:01Z\&quot; includes the first second of 4th minute).  | 
 **endTime** | **time.Time** | Filter transactions that occurred before this timestamp. The format is RFC3339 and is exclusive (for example, \&quot;2023-01-02T15:04:01Z\&quot; excludes the first second of 4th minute).  | 
 **cursor** | **string** | Parameter used in pagination requests. Maximum page size is set to 15. Set to the value of next for the next page of results. Set to the value of previous for the previous page of results. No other parameters can be set when this parameter is set.  | 
 **metadata** | **map[string]map[string]string** | Filter transactions by metadata key value pairs. | 

### Return type

[**TransactionsCursorResponse**](TransactionsCursorResponse.md)

### Authorization

[Authorization](../README.md#Authorization)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## RevertTransaction

> TransactionResponse RevertTransaction(ctx, ledger, txid).Execute()

Revert a ledger transaction by its ID

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
    ledger := "ledger001" // string | Name of the ledger.
    txid := int64(1234) // int64 | Transaction ID.

    configuration := client.NewConfiguration()
    apiClient := client.NewAPIClient(configuration)
    resp, r, err := apiClient.TransactionsApi.RevertTransaction(context.Background(), ledger, txid).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `TransactionsApi.RevertTransaction``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `RevertTransaction`: TransactionResponse
    fmt.Fprintf(os.Stdout, "Response from `TransactionsApi.RevertTransaction`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**ledger** | **string** | Name of the ledger. | 
**txid** | **int64** | Transaction ID. | 

### Other Parameters

Other parameters are passed through a pointer to a apiRevertTransactionRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

[**TransactionResponse**](TransactionResponse.md)

### Authorization

[Authorization](../README.md#Authorization)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

