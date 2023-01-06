# formance.TransactionsApi

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**addMetadataOnTransaction**](TransactionsApi.md#addMetadataOnTransaction) | **POST** /api/ledger/{ledger}/transactions/{txid}/metadata | Set the metadata of a transaction by its ID.
[**countTransactions**](TransactionsApi.md#countTransactions) | **HEAD** /api/ledger/{ledger}/transactions | Count the transactions from a ledger.
[**createTransaction**](TransactionsApi.md#createTransaction) | **POST** /api/ledger/{ledger}/transactions | Create a new transaction to a ledger.
[**createTransactions**](TransactionsApi.md#createTransactions) | **POST** /api/ledger/{ledger}/transactions/batch | Create a new batch of transactions to a ledger.
[**getTransaction**](TransactionsApi.md#getTransaction) | **GET** /api/ledger/{ledger}/transactions/{txid} | Get transaction from a ledger by its ID.
[**listTransactions**](TransactionsApi.md#listTransactions) | **GET** /api/ledger/{ledger}/transactions | List transactions from a ledger.
[**revertTransaction**](TransactionsApi.md#revertTransaction) | **POST** /api/ledger/{ledger}/transactions/{txid}/revert | Revert a ledger transaction by its ID.


# **addMetadataOnTransaction**
> void addMetadataOnTransaction()


### Example


```typescript
import { TransactionsApi, createConfiguration } from '@formancehq/formance';
import * as fs from 'fs';

const configuration = createConfiguration();
const apiInstance = new TransactionsApi(configuration);

let body:TransactionsApiAddMetadataOnTransactionRequest = {
  // string | Name of the ledger.
  ledger: "ledger001",
  // number | Transaction ID.
  txid: 1234,
  // { [key: string]: any; } | metadata (optional)
  requestBody: {
    "key": null,
  },
};

apiInstance.addMetadataOnTransaction(body).then((data:any) => {
  console.log('API called successfully. Returned data: ' + data);
}).catch((error:any) => console.error(error));
```


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **requestBody** | **{ [key: string]: any; }**| metadata |
 **ledger** | [**string**] | Name of the ledger. | defaults to undefined
 **txid** | [**number**] | Transaction ID. | defaults to undefined


### Return type

**void**

### Authorization

[Authorization](README.md#Authorization)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
**204** | No Content |  -  |
**400** | Bad Request |  -  |
**404** | Not Found |  -  |
**409** | Conflict |  -  |

[[Back to top]](#) [[Back to API list]](README.md#documentation-for-api-endpoints) [[Back to Model list]](README.md#documentation-for-models) [[Back to README]](README.md)

# **countTransactions**
> void countTransactions()


### Example


```typescript
import { TransactionsApi, createConfiguration } from '@formancehq/formance';
import * as fs from 'fs';

const configuration = createConfiguration();
const apiInstance = new TransactionsApi(configuration);

let body:TransactionsApiCountTransactionsRequest = {
  // string | Name of the ledger.
  ledger: "ledger001",
  // string | Filter transactions by reference field. (optional)
  reference: "ref:001",
  // string | Filter transactions with postings involving given account, either as source or destination (regular expression placed between ^ and $). (optional)
  account: "users:001",
  // string | Filter transactions with postings involving given account at source (regular expression placed between ^ and $). (optional)
  source: "users:001",
  // string | Filter transactions with postings involving given account at destination (regular expression placed between ^ and $). (optional)
  destination: "users:001",
  // any | Filter transactions by metadata key value pairs. Nested objects can be used as seen in the example below. (optional)
  metadata: {},
};

apiInstance.countTransactions(body).then((data:any) => {
  console.log('API called successfully. Returned data: ' + data);
}).catch((error:any) => console.error(error));
```


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ledger** | [**string**] | Name of the ledger. | defaults to undefined
 **reference** | [**string**] | Filter transactions by reference field. | (optional) defaults to undefined
 **account** | [**string**] | Filter transactions with postings involving given account, either as source or destination (regular expression placed between ^ and $). | (optional) defaults to undefined
 **source** | [**string**] | Filter transactions with postings involving given account at source (regular expression placed between ^ and $). | (optional) defaults to undefined
 **destination** | [**string**] | Filter transactions with postings involving given account at destination (regular expression placed between ^ and $). | (optional) defaults to undefined
 **metadata** | **any** | Filter transactions by metadata key value pairs. Nested objects can be used as seen in the example below. | (optional) defaults to undefined


### Return type

**void**

### Authorization

[Authorization](README.md#Authorization)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: Not defined


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | OK |  * Count -  <br>  |

[[Back to top]](#) [[Back to API list]](README.md#documentation-for-api-endpoints) [[Back to Model list]](README.md#documentation-for-models) [[Back to README]](README.md)

# **createTransaction**
> TransactionsResponse createTransaction(transactionData)


### Example


```typescript
import { TransactionsApi, createConfiguration } from '@formancehq/formance';
import * as fs from 'fs';

const configuration = createConfiguration();
const apiInstance = new TransactionsApi(configuration);

let body:TransactionsApiCreateTransactionRequest = {
  // string | Name of the ledger.
  ledger: "ledger001",
  // TransactionData
  transactionData: {
    timestamp: new Date('1970-01-01T00:00:00.00Z'),
    postings: [
      {
        amount: 100,
        asset: "COIN",
        destination: "users:002",
        source: "users:001",
      },
    ],
    reference: "ref:001",
    metadata: {
      "key": null,
    },
  },
  // boolean | Set the preview mode. Preview mode doesn't add the logs to the database or publish a message to the message broker. (optional)
  preview: true,
};

apiInstance.createTransaction(body).then((data:any) => {
  console.log('API called successfully. Returned data: ' + data);
}).catch((error:any) => console.error(error));
```


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **transactionData** | **TransactionData**|  |
 **ledger** | [**string**] | Name of the ledger. | defaults to undefined
 **preview** | [**boolean**] | Set the preview mode. Preview mode doesn&#39;t add the logs to the database or publish a message to the message broker. | (optional) defaults to undefined


### Return type

**TransactionsResponse**

### Authorization

[Authorization](README.md#Authorization)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | OK |  -  |
**304** | Not modified (when preview is enabled) |  -  |
**400** | Bad Request |  -  |
**409** | Conflict |  -  |

[[Back to top]](#) [[Back to API list]](README.md#documentation-for-api-endpoints) [[Back to Model list]](README.md#documentation-for-models) [[Back to README]](README.md)

# **createTransactions**
> TransactionsResponse createTransactions(transactions)


### Example


```typescript
import { TransactionsApi, createConfiguration } from '@formancehq/formance';
import * as fs from 'fs';

const configuration = createConfiguration();
const apiInstance = new TransactionsApi(configuration);

let body:TransactionsApiCreateTransactionsRequest = {
  // string | Name of the ledger.
  ledger: "ledger001",
  // Transactions
  transactions: {
    transactions: [
      {
        timestamp: new Date('1970-01-01T00:00:00.00Z'),
        postings: [
          {
            amount: 100,
            asset: "COIN",
            destination: "users:002",
            source: "users:001",
          },
        ],
        reference: "ref:001",
        metadata: {
          "key": null,
        },
      },
    ],
  },
};

apiInstance.createTransactions(body).then((data:any) => {
  console.log('API called successfully. Returned data: ' + data);
}).catch((error:any) => console.error(error));
```


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **transactions** | **Transactions**|  |
 **ledger** | [**string**] | Name of the ledger. | defaults to undefined


### Return type

**TransactionsResponse**

### Authorization

[Authorization](README.md#Authorization)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | OK |  -  |
**400** | Bad Request |  -  |
**409** | Conflict |  -  |

[[Back to top]](#) [[Back to API list]](README.md#documentation-for-api-endpoints) [[Back to Model list]](README.md#documentation-for-models) [[Back to README]](README.md)

# **getTransaction**
> TransactionResponse getTransaction()


### Example


```typescript
import { TransactionsApi, createConfiguration } from '@formancehq/formance';
import * as fs from 'fs';

const configuration = createConfiguration();
const apiInstance = new TransactionsApi(configuration);

let body:TransactionsApiGetTransactionRequest = {
  // string | Name of the ledger.
  ledger: "ledger001",
  // number | Transaction ID.
  txid: 1234,
};

apiInstance.getTransaction(body).then((data:any) => {
  console.log('API called successfully. Returned data: ' + data);
}).catch((error:any) => console.error(error));
```


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ledger** | [**string**] | Name of the ledger. | defaults to undefined
 **txid** | [**number**] | Transaction ID. | defaults to undefined


### Return type

**TransactionResponse**

### Authorization

[Authorization](README.md#Authorization)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | OK |  -  |
**400** | Bad Request |  -  |
**404** | Not Found |  -  |

[[Back to top]](#) [[Back to API list]](README.md#documentation-for-api-endpoints) [[Back to Model list]](README.md#documentation-for-models) [[Back to README]](README.md)

# **listTransactions**
> ListTransactions200Response listTransactions()

List transactions from a ledger, sorted by txid in descending order.

### Example


```typescript
import { TransactionsApi, createConfiguration } from '@formancehq/formance';
import * as fs from 'fs';

const configuration = createConfiguration();
const apiInstance = new TransactionsApi(configuration);

let body:TransactionsApiListTransactionsRequest = {
  // string | Name of the ledger.
  ledger: "ledger001",
  // number | The maximum number of results to return per page (optional)
  pageSize: 100,
  // string | Pagination cursor, will return transactions after given txid (in descending order). (optional)
  after: "1234",
  // string | Find transactions by reference field. (optional)
  reference: "ref:001",
  // string | Find transactions with postings involving given account, either as source or destination. (optional)
  account: "users:001",
  // string | Find transactions with postings involving given account at source. (optional)
  source: "users:001",
  // string | Find transactions with postings involving given account at destination. (optional)
  destination: "users:001",
  // string | Filter transactions that occurred after this timestamp. The format is RFC3339 and is inclusive (for example, 12:00:01 includes the first second of the minute).  (optional)
  startTime: "start_time_example",
  // string | Filter transactions that occurred before this timestamp. The format is RFC3339 and is exclusive (for example, 12:00:01 excludes the first second of the minute).  (optional)
  endTime: "end_time_example",
  // string | Parameter used in pagination requests. Maximum page size is set to 15. Set to the value of next for the next page of results.  Set to the value of previous for the previous page of results. No other parameters can be set when the pagination token is set.  (optional)
  paginationToken: "aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ==",
  // any | Filter transactions by metadata key value pairs. Nested objects can be used as seen in the example below. (optional)
  metadata: {},
};

apiInstance.listTransactions(body).then((data:any) => {
  console.log('API called successfully. Returned data: ' + data);
}).catch((error:any) => console.error(error));
```


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ledger** | [**string**] | Name of the ledger. | defaults to undefined
 **pageSize** | [**number**] | The maximum number of results to return per page | (optional) defaults to 15
 **after** | [**string**] | Pagination cursor, will return transactions after given txid (in descending order). | (optional) defaults to undefined
 **reference** | [**string**] | Find transactions by reference field. | (optional) defaults to undefined
 **account** | [**string**] | Find transactions with postings involving given account, either as source or destination. | (optional) defaults to undefined
 **source** | [**string**] | Find transactions with postings involving given account at source. | (optional) defaults to undefined
 **destination** | [**string**] | Find transactions with postings involving given account at destination. | (optional) defaults to undefined
 **startTime** | [**string**] | Filter transactions that occurred after this timestamp. The format is RFC3339 and is inclusive (for example, 12:00:01 includes the first second of the minute).  | (optional) defaults to undefined
 **endTime** | [**string**] | Filter transactions that occurred before this timestamp. The format is RFC3339 and is exclusive (for example, 12:00:01 excludes the first second of the minute).  | (optional) defaults to undefined
 **paginationToken** | [**string**] | Parameter used in pagination requests. Maximum page size is set to 15. Set to the value of next for the next page of results.  Set to the value of previous for the previous page of results. No other parameters can be set when the pagination token is set.  | (optional) defaults to undefined
 **metadata** | **any** | Filter transactions by metadata key value pairs. Nested objects can be used as seen in the example below. | (optional) defaults to undefined


### Return type

**ListTransactions200Response**

### Authorization

[Authorization](README.md#Authorization)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | OK |  -  |
**400** | Bad Request |  -  |

[[Back to top]](#) [[Back to API list]](README.md#documentation-for-api-endpoints) [[Back to Model list]](README.md#documentation-for-models) [[Back to README]](README.md)

# **revertTransaction**
> TransactionResponse revertTransaction()


### Example


```typescript
import { TransactionsApi, createConfiguration } from '@formancehq/formance';
import * as fs from 'fs';

const configuration = createConfiguration();
const apiInstance = new TransactionsApi(configuration);

let body:TransactionsApiRevertTransactionRequest = {
  // string | Name of the ledger.
  ledger: "ledger001",
  // number | Transaction ID.
  txid: 1234,
};

apiInstance.revertTransaction(body).then((data:any) => {
  console.log('API called successfully. Returned data: ' + data);
}).catch((error:any) => console.error(error));
```


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ledger** | [**string**] | Name of the ledger. | defaults to undefined
 **txid** | [**number**] | Transaction ID. | defaults to undefined


### Return type

**TransactionResponse**

### Authorization

[Authorization](README.md#Authorization)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | OK |  -  |
**400** | Bad Request |  -  |
**404** | Not Found |  -  |
**409** | Conflict |  -  |

[[Back to top]](#) [[Back to API list]](README.md#documentation-for-api-endpoints) [[Back to Model list]](README.md#documentation-for-models) [[Back to README]](README.md)

