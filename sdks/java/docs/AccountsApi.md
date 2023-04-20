# AccountsApi

All URIs are relative to *http://localhost*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**addMetadataToAccount**](AccountsApi.md#addMetadataToAccount) | **POST** api/ledger/{ledger}/accounts/{address}/metadata | Add metadata to an account |
| [**countAccounts**](AccountsApi.md#countAccounts) | **HEAD** api/ledger/{ledger}/accounts | Count the accounts from a ledger |
| [**getAccount**](AccountsApi.md#getAccount) | **GET** api/ledger/{ledger}/accounts/{address} | Get account by its address |
| [**listAccounts**](AccountsApi.md#listAccounts) | **GET** api/ledger/{ledger}/accounts | List accounts from a ledger |



## addMetadataToAccount

> addMetadataToAccount(ledger, address, requestBody, dryRun, async, idempotencyKey)

Add metadata to an account

### Example

```java
// Import classes:
import com.formance.formance.ApiClient;
import com.formance.formance.ApiException;
import com.formance.formance.Configuration;
import com.formance.formance.auth.*;
import com.formance.formance.models.*;
import com.formance.formance.api.AccountsApi;

public class Example {
    public static void main(String[] args) {
        ApiClient defaultClient = Configuration.getDefaultApiClient();
        defaultClient.setBasePath("http://localhost");
        
        // Configure OAuth2 access token for authorization: Authorization
        OAuth Authorization = (OAuth) defaultClient.getAuthentication("Authorization");
        Authorization.setAccessToken("YOUR ACCESS TOKEN");

        AccountsApi apiInstance = new AccountsApi(defaultClient);
        String ledger = "ledger001"; // String | Name of the ledger.
        String address = "users:001"; // String | Exact address of the account. It must match the following regular expressions pattern: ``` ^\\w+(:\\w+)*$ ``` 
        Map<String, String> requestBody = new HashMap(); // Map<String, String> | metadata
        Boolean dryRun = true; // Boolean | Set the dry run mode. Dry run mode doesn't add the logs to the database or publish a message to the message broker.
        Boolean async = true; // Boolean | Set async mode.
        String idempotencyKey = "idempotencyKey_example"; // String | Use an idempotency key
        try {
            apiInstance.addMetadataToAccount(ledger, address, requestBody, dryRun, async, idempotencyKey);
        } catch (ApiException e) {
            System.err.println("Exception when calling AccountsApi#addMetadataToAccount");
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
| **address** | **String**| Exact address of the account. It must match the following regular expressions pattern: &#x60;&#x60;&#x60; ^\\w+(:\\w+)*$ &#x60;&#x60;&#x60;  | |
| **requestBody** | [**Map&lt;String, String&gt;**](String.md)| metadata | |
| **dryRun** | **Boolean**| Set the dry run mode. Dry run mode doesn&#39;t add the logs to the database or publish a message to the message broker. | [optional] |
| **async** | **Boolean**| Set async mode. | [optional] |
| **idempotencyKey** | **String**| Use an idempotency key | [optional] |

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


## countAccounts

> countAccounts(ledger, address, metadata)

Count the accounts from a ledger

### Example

```java
// Import classes:
import com.formance.formance.ApiClient;
import com.formance.formance.ApiException;
import com.formance.formance.Configuration;
import com.formance.formance.auth.*;
import com.formance.formance.models.*;
import com.formance.formance.api.AccountsApi;

public class Example {
    public static void main(String[] args) {
        ApiClient defaultClient = Configuration.getDefaultApiClient();
        defaultClient.setBasePath("http://localhost");
        
        // Configure OAuth2 access token for authorization: Authorization
        OAuth Authorization = (OAuth) defaultClient.getAuthentication("Authorization");
        Authorization.setAccessToken("YOUR ACCESS TOKEN");

        AccountsApi apiInstance = new AccountsApi(defaultClient);
        String ledger = "ledger001"; // String | Name of the ledger.
        String address = "users:.+"; // String | Filter accounts by address pattern (regular expression placed between ^ and $).
        Object metadata = new HashMap(); // Object | Filter accounts by metadata key value pairs. The filter can be used like this metadata[key]=value1&metadata[a.nested.key]=value2
        try {
            apiInstance.countAccounts(ledger, address, metadata);
        } catch (ApiException e) {
            System.err.println("Exception when calling AccountsApi#countAccounts");
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
| **address** | **String**| Filter accounts by address pattern (regular expression placed between ^ and $). | [optional] |
| **metadata** | [**Object**](.md)| Filter accounts by metadata key value pairs. The filter can be used like this metadata[key]&#x3D;value1&amp;metadata[a.nested.key]&#x3D;value2 | [optional] |

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
| **204** | OK |  * Count -  <br>  |
| **0** | Error |  -  |


## getAccount

> AccountResponse getAccount(ledger, address)

Get account by its address

### Example

```java
// Import classes:
import com.formance.formance.ApiClient;
import com.formance.formance.ApiException;
import com.formance.formance.Configuration;
import com.formance.formance.auth.*;
import com.formance.formance.models.*;
import com.formance.formance.api.AccountsApi;

public class Example {
    public static void main(String[] args) {
        ApiClient defaultClient = Configuration.getDefaultApiClient();
        defaultClient.setBasePath("http://localhost");
        
        // Configure OAuth2 access token for authorization: Authorization
        OAuth Authorization = (OAuth) defaultClient.getAuthentication("Authorization");
        Authorization.setAccessToken("YOUR ACCESS TOKEN");

        AccountsApi apiInstance = new AccountsApi(defaultClient);
        String ledger = "ledger001"; // String | Name of the ledger.
        String address = "users:001"; // String | Exact address of the account. It must match the following regular expressions pattern: ``` ^\\w+(:\\w+)*$ ``` 
        try {
            AccountResponse result = apiInstance.getAccount(ledger, address);
            System.out.println(result);
        } catch (ApiException e) {
            System.err.println("Exception when calling AccountsApi#getAccount");
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
| **address** | **String**| Exact address of the account. It must match the following regular expressions pattern: &#x60;&#x60;&#x60; ^\\w+(:\\w+)*$ &#x60;&#x60;&#x60;  | |

### Return type

[**AccountResponse**](AccountResponse.md)

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


## listAccounts

> AccountsCursorResponse listAccounts(ledger, pageSize, address, metadata, balance, balanceOperator, cursor)

List accounts from a ledger

List accounts from a ledger, sorted by address in descending order.

### Example

```java
// Import classes:
import com.formance.formance.ApiClient;
import com.formance.formance.ApiException;
import com.formance.formance.Configuration;
import com.formance.formance.auth.*;
import com.formance.formance.models.*;
import com.formance.formance.api.AccountsApi;

public class Example {
    public static void main(String[] args) {
        ApiClient defaultClient = Configuration.getDefaultApiClient();
        defaultClient.setBasePath("http://localhost");
        
        // Configure OAuth2 access token for authorization: Authorization
        OAuth Authorization = (OAuth) defaultClient.getAuthentication("Authorization");
        Authorization.setAccessToken("YOUR ACCESS TOKEN");

        AccountsApi apiInstance = new AccountsApi(defaultClient);
        String ledger = "ledger001"; // String | Name of the ledger.
        Long pageSize = 100L; // Long | The maximum number of results to return per page. 
        String address = "users:.+"; // String | Filter accounts by address pattern (regular expression placed between ^ and $).
        Map<String, String> metadata = new HashMap(); // Map<String, String> | Filter accounts by metadata key value pairs. Nested objects can be used as seen in the example below.
        Long balance = 2400L; // Long | Filter accounts by their balance (default operator is gte)
        String balanceOperator = "gte"; // String | Operator used for the filtering of balances can be greater than/equal, less than/equal, greater than, less than, equal or not. 
        String cursor = "aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ=="; // String | Parameter used in pagination requests. Maximum page size is set to 15. Set to the value of next for the next page of results. Set to the value of previous for the previous page of results. No other parameters can be set when this parameter is set. 
        try {
            AccountsCursorResponse result = apiInstance.listAccounts(ledger, pageSize, address, metadata, balance, balanceOperator, cursor);
            System.out.println(result);
        } catch (ApiException e) {
            System.err.println("Exception when calling AccountsApi#listAccounts");
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
| **pageSize** | **Long**| The maximum number of results to return per page.  | [optional] |
| **address** | **String**| Filter accounts by address pattern (regular expression placed between ^ and $). | [optional] |
| **metadata** | [**Map&lt;String, String&gt;**](String.md)| Filter accounts by metadata key value pairs. Nested objects can be used as seen in the example below. | [optional] |
| **balance** | **Long**| Filter accounts by their balance (default operator is gte) | [optional] |
| **balanceOperator** | **String**| Operator used for the filtering of balances can be greater than/equal, less than/equal, greater than, less than, equal or not.  | [optional] [enum: gte, lte, gt, lt, e, ne] |
| **cursor** | **String**| Parameter used in pagination requests. Maximum page size is set to 15. Set to the value of next for the next page of results. Set to the value of previous for the previous page of results. No other parameters can be set when this parameter is set.  | [optional] |

### Return type

[**AccountsCursorResponse**](AccountsCursorResponse.md)

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

