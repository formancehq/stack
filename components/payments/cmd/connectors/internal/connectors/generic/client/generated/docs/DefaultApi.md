# \DefaultApi

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**GetAccountBalances**](DefaultApi.md#GetAccountBalances) | **Get** /accounts/{accountId}/balances | Get account balance
[**GetAccounts**](DefaultApi.md#GetAccounts) | **Get** /accounts | Get all accounts
[**GetBeneficiaries**](DefaultApi.md#GetBeneficiaries) | **Get** /beneficiaries | Get all beneficiaries
[**GetTransactions**](DefaultApi.md#GetTransactions) | **Get** /transactions | Get all transactions



## GetAccountBalances

> Balances GetAccountBalances(ctx, accountId).Execute()

Get account balance

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "github.com/formancehq/payments/genericclient"
)

func main() {
    accountId := "accountId_example" // string | 

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.DefaultApi.GetAccountBalances(context.Background(), accountId).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultApi.GetAccountBalances``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `GetAccountBalances`: Balances
    fmt.Fprintf(os.Stdout, "Response from `DefaultApi.GetAccountBalances`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**accountId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetAccountBalancesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**Balances**](Balances.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetAccounts

> []Account GetAccounts(ctx).PageSize(pageSize).Page(page).Sort(sort).CreatedAtFrom(createdAtFrom).Execute()

Get all accounts

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    "time"
    openapiclient "github.com/formancehq/payments/genericclient"
)

func main() {
    pageSize := int64(100) // int64 | Number of items per page (optional) (default to 100)
    page := int64(1) // int64 | Page number (optional) (default to 1)
    sort := "createdAt:asc" // string | Sort order (optional)
    createdAtFrom := time.Now() // time.Time | Filter by created at date (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.DefaultApi.GetAccounts(context.Background()).PageSize(pageSize).Page(page).Sort(sort).CreatedAtFrom(createdAtFrom).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultApi.GetAccounts``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `GetAccounts`: []Account
    fmt.Fprintf(os.Stdout, "Response from `DefaultApi.GetAccounts`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiGetAccountsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **pageSize** | **int64** | Number of items per page | [default to 100]
 **page** | **int64** | Page number | [default to 1]
 **sort** | **string** | Sort order | 
 **createdAtFrom** | **time.Time** | Filter by created at date | 

### Return type

[**[]Account**](Account.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetBeneficiaries

> []Beneficiary GetBeneficiaries(ctx).PageSize(pageSize).Page(page).Sort(sort).CreatedAtFrom(createdAtFrom).Execute()

Get all beneficiaries

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    "time"
    openapiclient "github.com/formancehq/payments/genericclient"
)

func main() {
    pageSize := int64(100) // int64 | Number of items per page (optional) (default to 100)
    page := int64(1) // int64 | Page number (optional) (default to 1)
    sort := "createdAt:asc" // string | Sort order (optional)
    createdAtFrom := time.Now() // time.Time | Filter by created at date (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.DefaultApi.GetBeneficiaries(context.Background()).PageSize(pageSize).Page(page).Sort(sort).CreatedAtFrom(createdAtFrom).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultApi.GetBeneficiaries``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `GetBeneficiaries`: []Beneficiary
    fmt.Fprintf(os.Stdout, "Response from `DefaultApi.GetBeneficiaries`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiGetBeneficiariesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **pageSize** | **int64** | Number of items per page | [default to 100]
 **page** | **int64** | Page number | [default to 1]
 **sort** | **string** | Sort order | 
 **createdAtFrom** | **time.Time** | Filter by created at date | 

### Return type

[**[]Beneficiary**](Beneficiary.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetTransactions

> []Transaction GetTransactions(ctx).PageSize(pageSize).Page(page).Sort(sort).UpdatedAtFrom(updatedAtFrom).Execute()

Get all transactions

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    "time"
    openapiclient "github.com/formancehq/payments/genericclient"
)

func main() {
    pageSize := int64(100) // int64 | Number of items per page (optional) (default to 100)
    page := int64(1) // int64 | Page number (optional) (default to 1)
    sort := "createdAt:asc" // string | Sort order (optional)
    updatedAtFrom := time.Now() // time.Time | Filter by updated at date (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.DefaultApi.GetTransactions(context.Background()).PageSize(pageSize).Page(page).Sort(sort).UpdatedAtFrom(updatedAtFrom).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultApi.GetTransactions``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `GetTransactions`: []Transaction
    fmt.Fprintf(os.Stdout, "Response from `DefaultApi.GetTransactions`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiGetTransactionsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **pageSize** | **int64** | Number of items per page | [default to 100]
 **page** | **int64** | Page number | [default to 1]
 **sort** | **string** | Sort order | 
 **updatedAtFrom** | **time.Time** | Filter by updated at date | 

### Return type

[**[]Transaction**](Transaction.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

