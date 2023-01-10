# PaymentsApi

All URIs are relative to *http://localhost*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**connectorsStripeTransfer**](PaymentsApi.md#connectorsStripeTransfer) | **POST** api/payments/connectors/stripe/transfer | Transfer funds between Stripe accounts |
| [**getAllConnectors**](PaymentsApi.md#getAllConnectors) | **GET** api/payments/connectors | Get all installed connectors |
| [**getAllConnectorsConfigs**](PaymentsApi.md#getAllConnectorsConfigs) | **GET** api/payments/connectors/configs | Get all available connectors configs |
| [**getConnectorTask**](PaymentsApi.md#getConnectorTask) | **GET** api/payments/connectors/{connector}/tasks/{taskId} | Read a specific task of the connector |
| [**getPayment**](PaymentsApi.md#getPayment) | **GET** api/payments/payments/{paymentId} | Returns a payment. |
| [**installConnector**](PaymentsApi.md#installConnector) | **POST** api/payments/connectors/{connector} | Install connector |
| [**listConnectorTasks**](PaymentsApi.md#listConnectorTasks) | **GET** api/payments/connectors/{connector}/tasks | List connector tasks |
| [**listPayments**](PaymentsApi.md#listPayments) | **GET** api/payments/payments | Returns a list of payments. |
| [**paymentslistAccounts**](PaymentsApi.md#paymentslistAccounts) | **GET** api/payments/accounts | Returns a list of accounts. |
| [**readConnectorConfig**](PaymentsApi.md#readConnectorConfig) | **GET** api/payments/connectors/{connector}/config | Read connector config |
| [**resetConnector**](PaymentsApi.md#resetConnector) | **POST** api/payments/connectors/{connector}/reset | Reset connector |
| [**uninstallConnector**](PaymentsApi.md#uninstallConnector) | **DELETE** api/payments/connectors/{connector} | Uninstall connector |



## connectorsStripeTransfer

> connectorsStripeTransfer(stripeTransferRequest)

Transfer funds between Stripe accounts

Execute a transfer between two Stripe accounts

### Example

```java
// Import classes:
import com.formance.formance.ApiClient;
import com.formance.formance.ApiException;
import com.formance.formance.Configuration;
import com.formance.formance.auth.*;
import com.formance.formance.models.*;
import com.formance.formance.api.PaymentsApi;

public class Example {
    public static void main(String[] args) {
        ApiClient defaultClient = Configuration.getDefaultApiClient();
        defaultClient.setBasePath("http://localhost");
        
        // Configure OAuth2 access token for authorization: Authorization
        OAuth Authorization = (OAuth) defaultClient.getAuthentication("Authorization");
        Authorization.setAccessToken("YOUR ACCESS TOKEN");

        PaymentsApi apiInstance = new PaymentsApi(defaultClient);
        StripeTransferRequest stripeTransferRequest = new StripeTransferRequest(); // StripeTransferRequest | 
        try {
            apiInstance.connectorsStripeTransfer(stripeTransferRequest);
        } catch (ApiException e) {
            System.err.println("Exception when calling PaymentsApi#connectorsStripeTransfer");
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
| **stripeTransferRequest** | [**StripeTransferRequest**](StripeTransferRequest.md)|  | |

### Return type

null (empty response body)

### Authorization

[Authorization](../README.md#Authorization)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: Not defined


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Transfer has been executed |  -  |


## getAllConnectors

> ListConnectorsResponse getAllConnectors()

Get all installed connectors

Get all installed connectors

### Example

```java
// Import classes:
import com.formance.formance.ApiClient;
import com.formance.formance.ApiException;
import com.formance.formance.Configuration;
import com.formance.formance.auth.*;
import com.formance.formance.models.*;
import com.formance.formance.api.PaymentsApi;

public class Example {
    public static void main(String[] args) {
        ApiClient defaultClient = Configuration.getDefaultApiClient();
        defaultClient.setBasePath("http://localhost");
        
        // Configure OAuth2 access token for authorization: Authorization
        OAuth Authorization = (OAuth) defaultClient.getAuthentication("Authorization");
        Authorization.setAccessToken("YOUR ACCESS TOKEN");

        PaymentsApi apiInstance = new PaymentsApi(defaultClient);
        try {
            ListConnectorsResponse result = apiInstance.getAllConnectors();
            System.out.println(result);
        } catch (ApiException e) {
            System.err.println("Exception when calling PaymentsApi#getAllConnectors");
            System.err.println("Status code: " + e.getCode());
            System.err.println("Reason: " + e.getResponseBody());
            System.err.println("Response headers: " + e.getResponseHeaders());
            e.printStackTrace();
        }
    }
}
```

### Parameters

This endpoint does not need any parameter.

### Return type

[**ListConnectorsResponse**](ListConnectorsResponse.md)

### Authorization

[Authorization](../README.md#Authorization)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | List of installed connectors |  -  |


## getAllConnectorsConfigs

> ListConnectorsConfigsResponse getAllConnectorsConfigs()

Get all available connectors configs

Get all available connectors configs

### Example

```java
// Import classes:
import com.formance.formance.ApiClient;
import com.formance.formance.ApiException;
import com.formance.formance.Configuration;
import com.formance.formance.auth.*;
import com.formance.formance.models.*;
import com.formance.formance.api.PaymentsApi;

public class Example {
    public static void main(String[] args) {
        ApiClient defaultClient = Configuration.getDefaultApiClient();
        defaultClient.setBasePath("http://localhost");
        
        // Configure OAuth2 access token for authorization: Authorization
        OAuth Authorization = (OAuth) defaultClient.getAuthentication("Authorization");
        Authorization.setAccessToken("YOUR ACCESS TOKEN");

        PaymentsApi apiInstance = new PaymentsApi(defaultClient);
        try {
            ListConnectorsConfigsResponse result = apiInstance.getAllConnectorsConfigs();
            System.out.println(result);
        } catch (ApiException e) {
            System.err.println("Exception when calling PaymentsApi#getAllConnectorsConfigs");
            System.err.println("Status code: " + e.getCode());
            System.err.println("Reason: " + e.getResponseBody());
            System.err.println("Response headers: " + e.getResponseHeaders());
            e.printStackTrace();
        }
    }
}
```

### Parameters

This endpoint does not need any parameter.

### Return type

[**ListConnectorsConfigsResponse**](ListConnectorsConfigsResponse.md)

### Authorization

[Authorization](../README.md#Authorization)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | List of available connectors configs |  -  |


## getConnectorTask

> ListConnectorTasks200ResponseInner getConnectorTask(connector, taskId)

Read a specific task of the connector

Get a specific task associated to the connector

### Example

```java
// Import classes:
import com.formance.formance.ApiClient;
import com.formance.formance.ApiException;
import com.formance.formance.Configuration;
import com.formance.formance.auth.*;
import com.formance.formance.models.*;
import com.formance.formance.api.PaymentsApi;

public class Example {
    public static void main(String[] args) {
        ApiClient defaultClient = Configuration.getDefaultApiClient();
        defaultClient.setBasePath("http://localhost");
        
        // Configure OAuth2 access token for authorization: Authorization
        OAuth Authorization = (OAuth) defaultClient.getAuthentication("Authorization");
        Authorization.setAccessToken("YOUR ACCESS TOKEN");

        PaymentsApi apiInstance = new PaymentsApi(defaultClient);
        Connectors connector = Connectors.fromValue("STRIPE"); // Connectors | The connector code
        String taskId = "task1"; // String | The task id
        try {
            ListConnectorTasks200ResponseInner result = apiInstance.getConnectorTask(connector, taskId);
            System.out.println(result);
        } catch (ApiException e) {
            System.err.println("Exception when calling PaymentsApi#getConnectorTask");
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
| **connector** | [**Connectors**](.md)| The connector code | [enum: STRIPE, DUMMY-PAY, SIE, MODULR, CURRENCY-CLOUD, BANKING-CIRCLE] |
| **taskId** | **String**| The task id | |

### Return type

[**ListConnectorTasks200ResponseInner**](ListConnectorTasks200ResponseInner.md)

### Authorization

[Authorization](../README.md#Authorization)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | The specified task |  -  |


## getPayment

> Payment getPayment(paymentId)

Returns a payment.

### Example

```java
// Import classes:
import com.formance.formance.ApiClient;
import com.formance.formance.ApiException;
import com.formance.formance.Configuration;
import com.formance.formance.auth.*;
import com.formance.formance.models.*;
import com.formance.formance.api.PaymentsApi;

public class Example {
    public static void main(String[] args) {
        ApiClient defaultClient = Configuration.getDefaultApiClient();
        defaultClient.setBasePath("http://localhost");
        
        // Configure OAuth2 access token for authorization: Authorization
        OAuth Authorization = (OAuth) defaultClient.getAuthentication("Authorization");
        Authorization.setAccessToken("YOUR ACCESS TOKEN");

        PaymentsApi apiInstance = new PaymentsApi(defaultClient);
        String paymentId = "XXX"; // String | The payment id
        try {
            Payment result = apiInstance.getPayment(paymentId);
            System.out.println(result);
        } catch (ApiException e) {
            System.err.println("Exception when calling PaymentsApi#getPayment");
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
| **paymentId** | **String**| The payment id | |

### Return type

[**Payment**](Payment.md)

### Authorization

[Authorization](../README.md#Authorization)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | A payment |  -  |


## installConnector

> installConnector(connector, connectorConfig)

Install connector

Install connector

### Example

```java
// Import classes:
import com.formance.formance.ApiClient;
import com.formance.formance.ApiException;
import com.formance.formance.Configuration;
import com.formance.formance.auth.*;
import com.formance.formance.models.*;
import com.formance.formance.api.PaymentsApi;

public class Example {
    public static void main(String[] args) {
        ApiClient defaultClient = Configuration.getDefaultApiClient();
        defaultClient.setBasePath("http://localhost");
        
        // Configure OAuth2 access token for authorization: Authorization
        OAuth Authorization = (OAuth) defaultClient.getAuthentication("Authorization");
        Authorization.setAccessToken("YOUR ACCESS TOKEN");

        PaymentsApi apiInstance = new PaymentsApi(defaultClient);
        Connectors connector = Connectors.fromValue("STRIPE"); // Connectors | The connector code
        ConnectorConfig connectorConfig = new ConnectorConfig(); // ConnectorConfig | 
        try {
            apiInstance.installConnector(connector, connectorConfig);
        } catch (ApiException e) {
            System.err.println("Exception when calling PaymentsApi#installConnector");
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
| **connector** | [**Connectors**](.md)| The connector code | [enum: STRIPE, DUMMY-PAY, SIE, MODULR, CURRENCY-CLOUD, BANKING-CIRCLE] |
| **connectorConfig** | [**ConnectorConfig**](ConnectorConfig.md)|  | |

### Return type

null (empty response body)

### Authorization

[Authorization](../README.md#Authorization)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: Not defined


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **204** | Connector has been installed |  -  |


## listConnectorTasks

> List&lt;ListConnectorTasks200ResponseInner&gt; listConnectorTasks(connector)

List connector tasks

List all tasks associated with this connector.

### Example

```java
// Import classes:
import com.formance.formance.ApiClient;
import com.formance.formance.ApiException;
import com.formance.formance.Configuration;
import com.formance.formance.auth.*;
import com.formance.formance.models.*;
import com.formance.formance.api.PaymentsApi;

public class Example {
    public static void main(String[] args) {
        ApiClient defaultClient = Configuration.getDefaultApiClient();
        defaultClient.setBasePath("http://localhost");
        
        // Configure OAuth2 access token for authorization: Authorization
        OAuth Authorization = (OAuth) defaultClient.getAuthentication("Authorization");
        Authorization.setAccessToken("YOUR ACCESS TOKEN");

        PaymentsApi apiInstance = new PaymentsApi(defaultClient);
        Connectors connector = Connectors.fromValue("STRIPE"); // Connectors | The connector code
        try {
            List<ListConnectorTasks200ResponseInner> result = apiInstance.listConnectorTasks(connector);
            System.out.println(result);
        } catch (ApiException e) {
            System.err.println("Exception when calling PaymentsApi#listConnectorTasks");
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
| **connector** | [**Connectors**](.md)| The connector code | [enum: STRIPE, DUMMY-PAY, SIE, MODULR, CURRENCY-CLOUD, BANKING-CIRCLE] |

### Return type

[**List&lt;ListConnectorTasks200ResponseInner&gt;**](ListConnectorTasks200ResponseInner.md)

### Authorization

[Authorization](../README.md#Authorization)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Task list |  -  |


## listPayments

> ListPaymentsResponse listPayments(limit, skip, sort)

Returns a list of payments.

### Example

```java
// Import classes:
import com.formance.formance.ApiClient;
import com.formance.formance.ApiException;
import com.formance.formance.Configuration;
import com.formance.formance.auth.*;
import com.formance.formance.models.*;
import com.formance.formance.api.PaymentsApi;

public class Example {
    public static void main(String[] args) {
        ApiClient defaultClient = Configuration.getDefaultApiClient();
        defaultClient.setBasePath("http://localhost");
        
        // Configure OAuth2 access token for authorization: Authorization
        OAuth Authorization = (OAuth) defaultClient.getAuthentication("Authorization");
        Authorization.setAccessToken("YOUR ACCESS TOKEN");

        PaymentsApi apiInstance = new PaymentsApi(defaultClient);
        Integer limit = 10; // Integer | Limit the number of payments to return, pagination can be achieved in conjunction with 'skip' parameter.
        Integer skip = 100; // Integer | How many payments to skip, pagination can be achieved in conjunction with 'limit' parameter.
        List<String> sort = Arrays.asList(); // List<String> | Field used to sort payments (Default is by date).
        try {
            ListPaymentsResponse result = apiInstance.listPayments(limit, skip, sort);
            System.out.println(result);
        } catch (ApiException e) {
            System.err.println("Exception when calling PaymentsApi#listPayments");
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
| **limit** | **Integer**| Limit the number of payments to return, pagination can be achieved in conjunction with &#39;skip&#39; parameter. | [optional] |
| **skip** | **Integer**| How many payments to skip, pagination can be achieved in conjunction with &#39;limit&#39; parameter. | [optional] |
| **sort** | [**List&lt;String&gt;**](String.md)| Field used to sort payments (Default is by date). | [optional] |

### Return type

[**ListPaymentsResponse**](ListPaymentsResponse.md)

### Authorization

[Authorization](../README.md#Authorization)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | A JSON array of payments |  -  |


## paymentslistAccounts

> ListAccountsResponse paymentslistAccounts(limit, skip, sort)

Returns a list of accounts.

### Example

```java
// Import classes:
import com.formance.formance.ApiClient;
import com.formance.formance.ApiException;
import com.formance.formance.Configuration;
import com.formance.formance.auth.*;
import com.formance.formance.models.*;
import com.formance.formance.api.PaymentsApi;

public class Example {
    public static void main(String[] args) {
        ApiClient defaultClient = Configuration.getDefaultApiClient();
        defaultClient.setBasePath("http://localhost");
        
        // Configure OAuth2 access token for authorization: Authorization
        OAuth Authorization = (OAuth) defaultClient.getAuthentication("Authorization");
        Authorization.setAccessToken("YOUR ACCESS TOKEN");

        PaymentsApi apiInstance = new PaymentsApi(defaultClient);
        Integer limit = 10; // Integer | Limit the number of accounts to return, pagination can be achieved in conjunction with 'skip' parameter.
        Integer skip = 100; // Integer | How many accounts to skip, pagination can be achieved in conjunction with 'limit' parameter.
        List<String> sort = Arrays.asList(); // List<String> | Field used to sort payments (Default is by date).
        try {
            ListAccountsResponse result = apiInstance.paymentslistAccounts(limit, skip, sort);
            System.out.println(result);
        } catch (ApiException e) {
            System.err.println("Exception when calling PaymentsApi#paymentslistAccounts");
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
| **limit** | **Integer**| Limit the number of accounts to return, pagination can be achieved in conjunction with &#39;skip&#39; parameter. | [optional] |
| **skip** | **Integer**| How many accounts to skip, pagination can be achieved in conjunction with &#39;limit&#39; parameter. | [optional] |
| **sort** | [**List&lt;String&gt;**](String.md)| Field used to sort payments (Default is by date). | [optional] |

### Return type

[**ListAccountsResponse**](ListAccountsResponse.md)

### Authorization

[Authorization](../README.md#Authorization)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | A JSON array of accounts |  -  |


## readConnectorConfig

> ConnectorConfig readConnectorConfig(connector)

Read connector config

Read connector config

### Example

```java
// Import classes:
import com.formance.formance.ApiClient;
import com.formance.formance.ApiException;
import com.formance.formance.Configuration;
import com.formance.formance.auth.*;
import com.formance.formance.models.*;
import com.formance.formance.api.PaymentsApi;

public class Example {
    public static void main(String[] args) {
        ApiClient defaultClient = Configuration.getDefaultApiClient();
        defaultClient.setBasePath("http://localhost");
        
        // Configure OAuth2 access token for authorization: Authorization
        OAuth Authorization = (OAuth) defaultClient.getAuthentication("Authorization");
        Authorization.setAccessToken("YOUR ACCESS TOKEN");

        PaymentsApi apiInstance = new PaymentsApi(defaultClient);
        Connectors connector = Connectors.fromValue("STRIPE"); // Connectors | The connector code
        try {
            ConnectorConfig result = apiInstance.readConnectorConfig(connector);
            System.out.println(result);
        } catch (ApiException e) {
            System.err.println("Exception when calling PaymentsApi#readConnectorConfig");
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
| **connector** | [**Connectors**](.md)| The connector code | [enum: STRIPE, DUMMY-PAY, SIE, MODULR, CURRENCY-CLOUD, BANKING-CIRCLE] |

### Return type

[**ConnectorConfig**](ConnectorConfig.md)

### Authorization

[Authorization](../README.md#Authorization)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Connector config |  -  |


## resetConnector

> resetConnector(connector)

Reset connector

Reset connector. Will remove the connector and ALL PAYMENTS generated with it.

### Example

```java
// Import classes:
import com.formance.formance.ApiClient;
import com.formance.formance.ApiException;
import com.formance.formance.Configuration;
import com.formance.formance.auth.*;
import com.formance.formance.models.*;
import com.formance.formance.api.PaymentsApi;

public class Example {
    public static void main(String[] args) {
        ApiClient defaultClient = Configuration.getDefaultApiClient();
        defaultClient.setBasePath("http://localhost");
        
        // Configure OAuth2 access token for authorization: Authorization
        OAuth Authorization = (OAuth) defaultClient.getAuthentication("Authorization");
        Authorization.setAccessToken("YOUR ACCESS TOKEN");

        PaymentsApi apiInstance = new PaymentsApi(defaultClient);
        Connectors connector = Connectors.fromValue("STRIPE"); // Connectors | The connector code
        try {
            apiInstance.resetConnector(connector);
        } catch (ApiException e) {
            System.err.println("Exception when calling PaymentsApi#resetConnector");
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
| **connector** | [**Connectors**](.md)| The connector code | [enum: STRIPE, DUMMY-PAY, SIE, MODULR, CURRENCY-CLOUD, BANKING-CIRCLE] |

### Return type

null (empty response body)

### Authorization

[Authorization](../README.md#Authorization)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: Not defined


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **204** | Connector has been reset |  -  |


## uninstallConnector

> uninstallConnector(connector)

Uninstall connector

Uninstall  connector

### Example

```java
// Import classes:
import com.formance.formance.ApiClient;
import com.formance.formance.ApiException;
import com.formance.formance.Configuration;
import com.formance.formance.auth.*;
import com.formance.formance.models.*;
import com.formance.formance.api.PaymentsApi;

public class Example {
    public static void main(String[] args) {
        ApiClient defaultClient = Configuration.getDefaultApiClient();
        defaultClient.setBasePath("http://localhost");
        
        // Configure OAuth2 access token for authorization: Authorization
        OAuth Authorization = (OAuth) defaultClient.getAuthentication("Authorization");
        Authorization.setAccessToken("YOUR ACCESS TOKEN");

        PaymentsApi apiInstance = new PaymentsApi(defaultClient);
        Connectors connector = Connectors.fromValue("STRIPE"); // Connectors | The connector code
        try {
            apiInstance.uninstallConnector(connector);
        } catch (ApiException e) {
            System.err.println("Exception when calling PaymentsApi#uninstallConnector");
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
| **connector** | [**Connectors**](.md)| The connector code | [enum: STRIPE, DUMMY-PAY, SIE, MODULR, CURRENCY-CLOUD, BANKING-CIRCLE] |

### Return type

null (empty response body)

### Authorization

[Authorization](../README.md#Authorization)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: Not defined


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **204** | Connector has been uninstalled |  -  |

