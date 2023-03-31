# TransactionsApi

All URIs are relative to *http://localhost*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**addMetadataOnTransaction**](TransactionsApi.md#addMetadataOnTransaction) | **POST** api/ledger/{ledger}/transactions/{txid}/metadata | Set the metadata of a transaction by its ID |
| [**countTransactions**](TransactionsApi.md#countTransactions) | **HEAD** api/ledger/{ledger}/transactions | Count the transactions from a ledger |
| [**createTransaction**](TransactionsApi.md#createTransaction) | **POST** api/ledger/{ledger}/transactions | Create a new transaction to a ledger |
| [**getTransaction**](TransactionsApi.md#getTransaction) | **GET** api/ledger/{ledger}/transactions/{txid} | Get transaction from a ledger by its ID |
| [**listTransactions**](TransactionsApi.md#listTransactions) | **GET** api/ledger/{ledger}/transactions | List transactions from a ledger |
| [**revertTransaction**](TransactionsApi.md#revertTransaction) | **POST** api/ledger/{ledger}/transactions/{txid}/revert | Revert a ledger transaction by its ID |



## addMetadataOnTransaction

> addMetadataOnTransaction(ledger, txid, requestBody)

Set the metadata of a transaction by its ID

### Example

```java
// Import classes:
import com.formance.formance.ApiClient;
import com.formance.formance.ApiException;
import com.formance.formance.Configuration;
import com.formance.formance.auth.*;
import com.formance.formance.models.*;
import com.formance.formance.api.TransactionsApi;

public class Example {
    public static void main(String[] args) {
        ApiClient defaultClient = Configuration.getDefaultApiClient();
        defaultClient.setBasePath("http://localhost");
        
        // Configure OAuth2 access token for authorization: Authorization
        OAuth Authorization = (OAuth) defaultClient.getAuthentication("Authorization");
        Authorization.setAccessToken("YOUR ACCESS TOKEN");

        TransactionsApi apiInstance = new TransactionsApi(defaultClient);
        String ledger = "ledger001"; // String | Name of the ledger.
        Long txid = 1234L; // Long | Transaction ID.
        Map<String, String> requestBody = new HashMap(); // Map<String, String> | metadata
        try {
            apiInstance.addMetadataOnTransaction(ledger, txid, requestBody);
        } catch (ApiException e) {
            System.err.println("Exception when calling TransactionsApi#addMetadataOnTransaction");
            System.err.println("Status code: " + e.getCode());
            System.err.println("Reason: " + e.getResponseBody());
            System.err.println("Response headers: " + e.getResponseHeaders());
            e.printStackTrace();
        }
    }
}
```

### Parameters


| Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **ledger** | **String**| Name of the ledger. | |
| **txid** | **Long**| Transaction ID. | |
| **requestBody** | [**Map&lt;String, String&gt;**](String.md)| metadata | [optional] |

### Return type

null (empty response body)

### Authorization

[Authorization](../README.md#Authorization)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **204** | No Content |  -  |
| **0** | Error |  -  |


## countTransactions

> countTransactions(ledger, reference, account, source, destination, startTime, endTime, metadata)

Count the transactions from a ledger

### Example

```java
// Import classes:
import com.formance.formance.ApiClient;
import com.formance.formance.ApiException;
import com.formance.formance.Configuration;
import com.formance.formance.auth.*;
import com.formance.formance.models.*;
import com.formance.formance.api.TransactionsApi;

public class Example {
    public static void main(String[] args) {
        ApiClient defaultClient = Configuration.getDefaultApiClient();
        defaultClient.setBasePath("http://localhost");
        
        // Configure OAuth2 access token for authorization: Authorization
        OAuth Authorization = (OAuth) defaultClient.getAuthentication("Authorization");
        Authorization.setAccessToken("YOUR ACCESS TOKEN");

        TransactionsApi apiInstance = new TransactionsApi(defaultClient);
        String ledger = "ledger001"; // String | Name of the ledger.
        String reference = "ref:001"; // String | Filter transactions by reference field.
        String account = "users:001"; // String | Filter transactions with postings involving given account, either as source or destination (regular expression placed between ^ and $).
        String source = "users:001"; // String | Filter transactions with postings involving given account at source (regular expression placed between ^ and $).
        String destination = "users:001"; // String | Filter transactions with postings involving given account at destination (regular expression placed between ^ and $).
        OffsetDateTime startTime = OffsetDateTime.now(); // OffsetDateTime | Filter transactions that occurred after this timestamp. The format is RFC3339 and is inclusive (for example, \"2023-01-02T15:04:01Z\" includes the first second of 4th minute). 
        OffsetDateTime endTime = OffsetDateTime.now(); // OffsetDateTime | Filter transactions that occurred before this timestamp. The format is RFC3339 and is exclusive (for example, \"2023-01-02T15:04:01Z\" excludes the first second of 4th minute). 
        Object metadata = new HashMap(); // Object | Filter transactions by metadata key value pairs. Nested objects can be used as seen in the example below.
        try {
            apiInstance.countTransactions(ledger, reference, account, source, destination, startTime, endTime, metadata);
        } catch (ApiException e) {
            System.err.println("Exception when calling TransactionsApi#countTransactions");
            System.err.println("Status code: " + e.getCode());
            System.err.println("Reason: " + e.getResponseBody());
            System.err.println("Response headers: " + e.getResponseHeaders());
            e.printStackTrace();
        }
    }
}
```

### Parameters


| Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **ledger** | **String**| Name of the ledger. | |
| **reference** | **String**| Filter transactions by reference field. | [optional] |
| **account** | **String**| Filter transactions with postings involving given account, either as source or destination (regular expression placed between ^ and $). | [optional] |
| **source** | **String**| Filter transactions with postings involving given account at source (regular expression placed between ^ and $). | [optional] |
| **destination** | **String**| Filter transactions with postings involving given account at destination (regular expression placed between ^ and $). | [optional] |
| **startTime** | **OffsetDateTime**| Filter transactions that occurred after this timestamp. The format is RFC3339 and is inclusive (for example, \&quot;2023-01-02T15:04:01Z\&quot; includes the first second of 4th minute).  | [optional] |
| **endTime** | **OffsetDateTime**| Filter transactions that occurred before this timestamp. The format is RFC3339 and is exclusive (for example, \&quot;2023-01-02T15:04:01Z\&quot; excludes the first second of 4th minute).  | [optional] |
| **metadata** | [**Object**](.md)| Filter transactions by metadata key value pairs. Nested objects can be used as seen in the example below. | [optional] |

### Return type

null (empty response body)

### Authorization

[Authorization](../README.md#Authorization)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | OK |  * Count -  <br>  |
| **0** | Error |  -  |


## createTransaction

> TransactionResponse createTransaction(ledger, postTransaction, preview)

Create a new transaction to a ledger

### Example

```java
// Import classes:
import com.formance.formance.ApiClient;
import com.formance.formance.ApiException;
import com.formance.formance.Configuration;
import com.formance.formance.auth.*;
import com.formance.formance.models.*;
import com.formance.formance.api.TransactionsApi;

public class Example {
    public static void main(String[] args) {
        ApiClient defaultClient = Configuration.getDefaultApiClient();
        defaultClient.setBasePath("http://localhost");
        
        // Configure OAuth2 access token for authorization: Authorization
        OAuth Authorization = (OAuth) defaultClient.getAuthentication("Authorization");
        Authorization.setAccessToken("YOUR ACCESS TOKEN");

        TransactionsApi apiInstance = new TransactionsApi(defaultClient);
        String ledger = "ledger001"; // String | Name of the ledger.
        PostTransaction postTransaction = new PostTransaction(); // PostTransaction | The request body must contain at least one of the following objects:   - `postings`: suitable for simple transactions   - `script`: enabling more complex transactions with Numscript 
        Boolean preview = true; // Boolean | Set the preview mode. Preview mode doesn't add the logs to the database or publish a message to the message broker.
        try {
            TransactionResponse result = apiInstance.createTransaction(ledger, postTransaction, preview);
            System.out.println(result);
        } catch (ApiException e) {
            System.err.println("Exception when calling TransactionsApi#createTransaction");
            System.err.println("Status code: " + e.getCode());
            System.err.println("Reason: " + e.getResponseBody());
            System.err.println("Response headers: " + e.getResponseHeaders());
            e.printStackTrace();
        }
    }
}
```

### Parameters


| Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **ledger** | **String**| Name of the ledger. | |
| **postTransaction** | [**PostTransaction**](PostTransaction.md)| The request body must contain at least one of the following objects:   - &#x60;postings&#x60;: suitable for simple transactions   - &#x60;script&#x60;: enabling more complex transactions with Numscript  | |
| **preview** | **Boolean**| Set the preview mode. Preview mode doesn&#39;t add the logs to the database or publish a message to the message broker. | [optional] |

### Return type

[**TransactionResponse**](TransactionResponse.md)

### Authorization

[Authorization](../README.md#Authorization)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | OK |  -  |
| **0** | Error |  -  |


## getTransaction

> TransactionResponse getTransaction(ledger, txid)

Get transaction from a ledger by its ID

### Example

```java
// Import classes:
import com.formance.formance.ApiClient;
import com.formance.formance.ApiException;
import com.formance.formance.Configuration;
import com.formance.formance.auth.*;
import com.formance.formance.models.*;
import com.formance.formance.api.TransactionsApi;

public class Example {
    public static void main(String[] args) {
        ApiClient defaultClient = Configuration.getDefaultApiClient();
        defaultClient.setBasePath("http://localhost");
        
        // Configure OAuth2 access token for authorization: Authorization
        OAuth Authorization = (OAuth) defaultClient.getAuthentication("Authorization");
        Authorization.setAccessToken("YOUR ACCESS TOKEN");

        TransactionsApi apiInstance = new TransactionsApi(defaultClient);
        String ledger = "ledger001"; // String | Name of the ledger.
        Long txid = 1234L; // Long | Transaction ID.
        try {
            TransactionResponse result = apiInstance.getTransaction(ledger, txid);
            System.out.println(result);
        } catch (ApiException e) {
            System.err.println("Exception when calling TransactionsApi#getTransaction");
            System.err.println("Status code: " + e.getCode());
            System.err.println("Reason: " + e.getResponseBody());
            System.err.println("Response headers: " + e.getResponseHeaders());
            e.printStackTrace();
        }
    }
}
```

### Parameters


| Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **ledger** | **String**| Name of the ledger. | |
| **txid** | **Long**| Transaction ID. | |

### Return type

[**TransactionResponse**](TransactionResponse.md)

### Authorization

[Authorization](../README.md#Authorization)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | OK |  -  |
| **0** | Error |  -  |


## listTransactions

> TransactionsCursorResponse listTransactions(ledger, pageSize, after, reference, account, source, destination, startTime, endTime, cursor, metadata)

List transactions from a ledger

List transactions from a ledger, sorted by txid in descending order.

### Example

```java
// Import classes:
import com.formance.formance.ApiClient;
import com.formance.formance.ApiException;
import com.formance.formance.Configuration;
import com.formance.formance.auth.*;
import com.formance.formance.models.*;
import com.formance.formance.api.TransactionsApi;

public class Example {
    public static void main(String[] args) {
        ApiClient defaultClient = Configuration.getDefaultApiClient();
        defaultClient.setBasePath("http://localhost");
        
        // Configure OAuth2 access token for authorization: Authorization
        OAuth Authorization = (OAuth) defaultClient.getAuthentication("Authorization");
        Authorization.setAccessToken("YOUR ACCESS TOKEN");

        TransactionsApi apiInstance = new TransactionsApi(defaultClient);
        String ledger = "ledger001"; // String | Name of the ledger.
        Long pageSize = 15L; // Long | The maximum number of results to return per page. 
        String after = "1234"; // String | Pagination cursor, will return transactions after given txid (in descending order).
        String reference = "ref:001"; // String | Find transactions by reference field.
        String account = "users:001"; // String | Filter transactions with postings involving given account, either as source or destination (regular expression placed between ^ and $).
        String source = "users:001"; // String | Filter transactions with postings involving given account at source (regular expression placed between ^ and $).
        String destination = "users:001"; // String | Filter transactions with postings involving given account at destination (regular expression placed between ^ and $).
        OffsetDateTime startTime = OffsetDateTime.now(); // OffsetDateTime | Filter transactions that occurred after this timestamp. The format is RFC3339 and is inclusive (for example, \"2023-01-02T15:04:01Z\" includes the first second of 4th minute). 
        OffsetDateTime endTime = OffsetDateTime.now(); // OffsetDateTime | Filter transactions that occurred before this timestamp. The format is RFC3339 and is exclusive (for example, \"2023-01-02T15:04:01Z\" excludes the first second of 4th minute). 
        String cursor = "aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ=="; // String | Parameter used in pagination requests. Maximum page size is set to 15. Set to the value of next for the next page of results. Set to the value of previous for the previous page of results. No other parameters can be set when this parameter is set. 
        Map<String, String> metadata = new HashMap(); // Map<String, String> | Filter transactions by metadata key value pairs.
        try {
            TransactionsCursorResponse result = apiInstance.listTransactions(ledger, pageSize, after, reference, account, source, destination, startTime, endTime, cursor, metadata);
            System.out.println(result);
        } catch (ApiException e) {
            System.err.println("Exception when calling TransactionsApi#listTransactions");
            System.err.println("Status code: " + e.getCode());
            System.err.println("Reason: " + e.getResponseBody());
            System.err.println("Response headers: " + e.getResponseHeaders());
            e.printStackTrace();
        }
    }
}
```

### Parameters


| Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **ledger** | **String**| Name of the ledger. | |
| **pageSize** | **Long**| The maximum number of results to return per page.  | [optional] [default to 15] |
| **after** | **String**| Pagination cursor, will return transactions after given txid (in descending order). | [optional] |
| **reference** | **String**| Find transactions by reference field. | [optional] |
| **account** | **String**| Filter transactions with postings involving given account, either as source or destination (regular expression placed between ^ and $). | [optional] |
| **source** | **String**| Filter transactions with postings involving given account at source (regular expression placed between ^ and $). | [optional] |
| **destination** | **String**| Filter transactions with postings involving given account at destination (regular expression placed between ^ and $). | [optional] |
| **startTime** | **OffsetDateTime**| Filter transactions that occurred after this timestamp. The format is RFC3339 and is inclusive (for example, \&quot;2023-01-02T15:04:01Z\&quot; includes the first second of 4th minute).  | [optional] |
| **endTime** | **OffsetDateTime**| Filter transactions that occurred before this timestamp. The format is RFC3339 and is exclusive (for example, \&quot;2023-01-02T15:04:01Z\&quot; excludes the first second of 4th minute).  | [optional] |
| **cursor** | **String**| Parameter used in pagination requests. Maximum page size is set to 15. Set to the value of next for the next page of results. Set to the value of previous for the previous page of results. No other parameters can be set when this parameter is set.  | [optional] |
| **metadata** | [**Map&lt;String, String&gt;**](String.md)| Filter transactions by metadata key value pairs. | [optional] |

### Return type

[**TransactionsCursorResponse**](TransactionsCursorResponse.md)

### Authorization

[Authorization](../README.md#Authorization)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | OK |  -  |
| **0** | Error |  -  |


## revertTransaction

> TransactionResponse revertTransaction(ledger, txid)

Revert a ledger transaction by its ID

### Example

```java
// Import classes:
import com.formance.formance.ApiClient;
import com.formance.formance.ApiException;
import com.formance.formance.Configuration;
import com.formance.formance.auth.*;
import com.formance.formance.models.*;
import com.formance.formance.api.TransactionsApi;

public class Example {
    public static void main(String[] args) {
        ApiClient defaultClient = Configuration.getDefaultApiClient();
        defaultClient.setBasePath("http://localhost");
        
        // Configure OAuth2 access token for authorization: Authorization
        OAuth Authorization = (OAuth) defaultClient.getAuthentication("Authorization");
        Authorization.setAccessToken("YOUR ACCESS TOKEN");

        TransactionsApi apiInstance = new TransactionsApi(defaultClient);
        String ledger = "ledger001"; // String | Name of the ledger.
        Long txid = 1234L; // Long | Transaction ID.
        try {
            TransactionResponse result = apiInstance.revertTransaction(ledger, txid);
            System.out.println(result);
        } catch (ApiException e) {
            System.err.println("Exception when calling TransactionsApi#revertTransaction");
            System.err.println("Status code: " + e.getCode());
            System.err.println("Reason: " + e.getResponseBody());
            System.err.println("Response headers: " + e.getResponseHeaders());
            e.printStackTrace();
        }
    }
}
```

### Parameters


| Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **ledger** | **String**| Name of the ledger. | |
| **txid** | **Long**| Transaction ID. | |

### Return type

[**TransactionResponse**](TransactionResponse.md)

### Authorization

[Authorization](../README.md#Authorization)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | OK |  -  |
| **0** | Error |  -  |

