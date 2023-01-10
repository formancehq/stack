# Formance\PaymentsApi

All URIs are relative to http://localhost, except if the operation defines another base path.

| Method | HTTP request | Description |
| ------------- | ------------- | ------------- |
| [**connectorsStripeTransfer()**](PaymentsApi.md#connectorsStripeTransfer) | **POST** /api/payments/connectors/stripe/transfer | Transfer funds between Stripe accounts |
| [**getAllConnectors()**](PaymentsApi.md#getAllConnectors) | **GET** /api/payments/connectors | Get all installed connectors |
| [**getAllConnectorsConfigs()**](PaymentsApi.md#getAllConnectorsConfigs) | **GET** /api/payments/connectors/configs | Get all available connectors configs |
| [**getConnectorTask()**](PaymentsApi.md#getConnectorTask) | **GET** /api/payments/connectors/{connector}/tasks/{taskId} | Read a specific task of the connector |
| [**getPayment()**](PaymentsApi.md#getPayment) | **GET** /api/payments/payments/{paymentId} | Returns a payment. |
| [**installConnector()**](PaymentsApi.md#installConnector) | **POST** /api/payments/connectors/{connector} | Install connector |
| [**listConnectorTasks()**](PaymentsApi.md#listConnectorTasks) | **GET** /api/payments/connectors/{connector}/tasks | List connector tasks |
| [**listPayments()**](PaymentsApi.md#listPayments) | **GET** /api/payments/payments | Returns a list of payments. |
| [**paymentslistAccounts()**](PaymentsApi.md#paymentslistAccounts) | **GET** /api/payments/accounts | Returns a list of accounts. |
| [**readConnectorConfig()**](PaymentsApi.md#readConnectorConfig) | **GET** /api/payments/connectors/{connector}/config | Read connector config |
| [**resetConnector()**](PaymentsApi.md#resetConnector) | **POST** /api/payments/connectors/{connector}/reset | Reset connector |
| [**uninstallConnector()**](PaymentsApi.md#uninstallConnector) | **DELETE** /api/payments/connectors/{connector} | Uninstall connector |


## `connectorsStripeTransfer()`

```php
connectorsStripeTransfer($stripe_transfer_request)
```

Transfer funds between Stripe accounts

Execute a transfer between two Stripe accounts

### Example

```php
<?php
require_once(__DIR__ . '/vendor/autoload.php');


// Configure OAuth2 access token for authorization: Authorization
$config = Formance\Configuration::getDefaultConfiguration()->setAccessToken('YOUR_ACCESS_TOKEN');


$apiInstance = new Formance\Api\PaymentsApi(
    // If you want use custom http client, pass your client which implements `GuzzleHttp\ClientInterface`.
    // This is optional, `GuzzleHttp\Client` will be used as default.
    new GuzzleHttp\Client(),
    $config
);
$stripe_transfer_request = new \Formance\Model\StripeTransferRequest(); // \Formance\Model\StripeTransferRequest

try {
    $apiInstance->connectorsStripeTransfer($stripe_transfer_request);
} catch (Exception $e) {
    echo 'Exception when calling PaymentsApi->connectorsStripeTransfer: ', $e->getMessage(), PHP_EOL;
}
```

### Parameters

| Name | Type | Description  | Notes |
| ------------- | ------------- | ------------- | ------------- |
| **stripe_transfer_request** | [**\Formance\Model\StripeTransferRequest**](../Model/StripeTransferRequest.md)|  | |

### Return type

void (empty response body)

### Authorization

[Authorization](../../README.md#Authorization)

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../../README.md#endpoints)
[[Back to Model list]](../../README.md#models)
[[Back to README]](../../README.md)

## `getAllConnectors()`

```php
getAllConnectors(): \Formance\Model\ListConnectorsResponse
```

Get all installed connectors

Get all installed connectors

### Example

```php
<?php
require_once(__DIR__ . '/vendor/autoload.php');


// Configure OAuth2 access token for authorization: Authorization
$config = Formance\Configuration::getDefaultConfiguration()->setAccessToken('YOUR_ACCESS_TOKEN');


$apiInstance = new Formance\Api\PaymentsApi(
    // If you want use custom http client, pass your client which implements `GuzzleHttp\ClientInterface`.
    // This is optional, `GuzzleHttp\Client` will be used as default.
    new GuzzleHttp\Client(),
    $config
);

try {
    $result = $apiInstance->getAllConnectors();
    print_r($result);
} catch (Exception $e) {
    echo 'Exception when calling PaymentsApi->getAllConnectors: ', $e->getMessage(), PHP_EOL;
}
```

### Parameters

This endpoint does not need any parameter.

### Return type

[**\Formance\Model\ListConnectorsResponse**](../Model/ListConnectorsResponse.md)

### Authorization

[Authorization](../../README.md#Authorization)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`

[[Back to top]](#) [[Back to API list]](../../README.md#endpoints)
[[Back to Model list]](../../README.md#models)
[[Back to README]](../../README.md)

## `getAllConnectorsConfigs()`

```php
getAllConnectorsConfigs(): \Formance\Model\ListConnectorsConfigsResponse
```

Get all available connectors configs

Get all available connectors configs

### Example

```php
<?php
require_once(__DIR__ . '/vendor/autoload.php');


// Configure OAuth2 access token for authorization: Authorization
$config = Formance\Configuration::getDefaultConfiguration()->setAccessToken('YOUR_ACCESS_TOKEN');


$apiInstance = new Formance\Api\PaymentsApi(
    // If you want use custom http client, pass your client which implements `GuzzleHttp\ClientInterface`.
    // This is optional, `GuzzleHttp\Client` will be used as default.
    new GuzzleHttp\Client(),
    $config
);

try {
    $result = $apiInstance->getAllConnectorsConfigs();
    print_r($result);
} catch (Exception $e) {
    echo 'Exception when calling PaymentsApi->getAllConnectorsConfigs: ', $e->getMessage(), PHP_EOL;
}
```

### Parameters

This endpoint does not need any parameter.

### Return type

[**\Formance\Model\ListConnectorsConfigsResponse**](../Model/ListConnectorsConfigsResponse.md)

### Authorization

[Authorization](../../README.md#Authorization)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`

[[Back to top]](#) [[Back to API list]](../../README.md#endpoints)
[[Back to Model list]](../../README.md#models)
[[Back to README]](../../README.md)

## `getConnectorTask()`

```php
getConnectorTask($connector, $task_id): \Formance\Model\ListConnectorTasks200ResponseInner
```

Read a specific task of the connector

Get a specific task associated to the connector

### Example

```php
<?php
require_once(__DIR__ . '/vendor/autoload.php');


// Configure OAuth2 access token for authorization: Authorization
$config = Formance\Configuration::getDefaultConfiguration()->setAccessToken('YOUR_ACCESS_TOKEN');


$apiInstance = new Formance\Api\PaymentsApi(
    // If you want use custom http client, pass your client which implements `GuzzleHttp\ClientInterface`.
    // This is optional, `GuzzleHttp\Client` will be used as default.
    new GuzzleHttp\Client(),
    $config
);
$connector = new \Formance\Model\Connectors(); // Connectors | The connector code
$task_id = task1; // string | The task id

try {
    $result = $apiInstance->getConnectorTask($connector, $task_id);
    print_r($result);
} catch (Exception $e) {
    echo 'Exception when calling PaymentsApi->getConnectorTask: ', $e->getMessage(), PHP_EOL;
}
```

### Parameters

| Name | Type | Description  | Notes |
| ------------- | ------------- | ------------- | ------------- |
| **connector** | [**Connectors**](../Model/.md)| The connector code | |
| **task_id** | **string**| The task id | |

### Return type

[**\Formance\Model\ListConnectorTasks200ResponseInner**](../Model/ListConnectorTasks200ResponseInner.md)

### Authorization

[Authorization](../../README.md#Authorization)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`

[[Back to top]](#) [[Back to API list]](../../README.md#endpoints)
[[Back to Model list]](../../README.md#models)
[[Back to README]](../../README.md)

## `getPayment()`

```php
getPayment($payment_id): \Formance\Model\Payment
```

Returns a payment.

### Example

```php
<?php
require_once(__DIR__ . '/vendor/autoload.php');


// Configure OAuth2 access token for authorization: Authorization
$config = Formance\Configuration::getDefaultConfiguration()->setAccessToken('YOUR_ACCESS_TOKEN');


$apiInstance = new Formance\Api\PaymentsApi(
    // If you want use custom http client, pass your client which implements `GuzzleHttp\ClientInterface`.
    // This is optional, `GuzzleHttp\Client` will be used as default.
    new GuzzleHttp\Client(),
    $config
);
$payment_id = XXX; // string | The payment id

try {
    $result = $apiInstance->getPayment($payment_id);
    print_r($result);
} catch (Exception $e) {
    echo 'Exception when calling PaymentsApi->getPayment: ', $e->getMessage(), PHP_EOL;
}
```

### Parameters

| Name | Type | Description  | Notes |
| ------------- | ------------- | ------------- | ------------- |
| **payment_id** | **string**| The payment id | |

### Return type

[**\Formance\Model\Payment**](../Model/Payment.md)

### Authorization

[Authorization](../../README.md#Authorization)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`

[[Back to top]](#) [[Back to API list]](../../README.md#endpoints)
[[Back to Model list]](../../README.md#models)
[[Back to README]](../../README.md)

## `installConnector()`

```php
installConnector($connector, $connector_config)
```

Install connector

Install connector

### Example

```php
<?php
require_once(__DIR__ . '/vendor/autoload.php');


// Configure OAuth2 access token for authorization: Authorization
$config = Formance\Configuration::getDefaultConfiguration()->setAccessToken('YOUR_ACCESS_TOKEN');


$apiInstance = new Formance\Api\PaymentsApi(
    // If you want use custom http client, pass your client which implements `GuzzleHttp\ClientInterface`.
    // This is optional, `GuzzleHttp\Client` will be used as default.
    new GuzzleHttp\Client(),
    $config
);
$connector = new \Formance\Model\Connectors(); // Connectors | The connector code
$connector_config = new \Formance\Model\ConnectorConfig(); // \Formance\Model\ConnectorConfig

try {
    $apiInstance->installConnector($connector, $connector_config);
} catch (Exception $e) {
    echo 'Exception when calling PaymentsApi->installConnector: ', $e->getMessage(), PHP_EOL;
}
```

### Parameters

| Name | Type | Description  | Notes |
| ------------- | ------------- | ------------- | ------------- |
| **connector** | [**Connectors**](../Model/.md)| The connector code | |
| **connector_config** | [**\Formance\Model\ConnectorConfig**](../Model/ConnectorConfig.md)|  | |

### Return type

void (empty response body)

### Authorization

[Authorization](../../README.md#Authorization)

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../../README.md#endpoints)
[[Back to Model list]](../../README.md#models)
[[Back to README]](../../README.md)

## `listConnectorTasks()`

```php
listConnectorTasks($connector): \Formance\Model\ListConnectorTasks200ResponseInner[]
```

List connector tasks

List all tasks associated with this connector.

### Example

```php
<?php
require_once(__DIR__ . '/vendor/autoload.php');


// Configure OAuth2 access token for authorization: Authorization
$config = Formance\Configuration::getDefaultConfiguration()->setAccessToken('YOUR_ACCESS_TOKEN');


$apiInstance = new Formance\Api\PaymentsApi(
    // If you want use custom http client, pass your client which implements `GuzzleHttp\ClientInterface`.
    // This is optional, `GuzzleHttp\Client` will be used as default.
    new GuzzleHttp\Client(),
    $config
);
$connector = new \Formance\Model\Connectors(); // Connectors | The connector code

try {
    $result = $apiInstance->listConnectorTasks($connector);
    print_r($result);
} catch (Exception $e) {
    echo 'Exception when calling PaymentsApi->listConnectorTasks: ', $e->getMessage(), PHP_EOL;
}
```

### Parameters

| Name | Type | Description  | Notes |
| ------------- | ------------- | ------------- | ------------- |
| **connector** | [**Connectors**](../Model/.md)| The connector code | |

### Return type

[**\Formance\Model\ListConnectorTasks200ResponseInner[]**](../Model/ListConnectorTasks200ResponseInner.md)

### Authorization

[Authorization](../../README.md#Authorization)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`

[[Back to top]](#) [[Back to API list]](../../README.md#endpoints)
[[Back to Model list]](../../README.md#models)
[[Back to README]](../../README.md)

## `listPayments()`

```php
listPayments($limit, $skip, $sort): \Formance\Model\ListPaymentsResponse
```

Returns a list of payments.

### Example

```php
<?php
require_once(__DIR__ . '/vendor/autoload.php');


// Configure OAuth2 access token for authorization: Authorization
$config = Formance\Configuration::getDefaultConfiguration()->setAccessToken('YOUR_ACCESS_TOKEN');


$apiInstance = new Formance\Api\PaymentsApi(
    // If you want use custom http client, pass your client which implements `GuzzleHttp\ClientInterface`.
    // This is optional, `GuzzleHttp\Client` will be used as default.
    new GuzzleHttp\Client(),
    $config
);
$limit = 10; // int | Limit the number of payments to return, pagination can be achieved in conjunction with 'skip' parameter.
$skip = 100; // int | How many payments to skip, pagination can be achieved in conjunction with 'limit' parameter.
$sort = status; // string[] | Field used to sort payments (Default is by date).

try {
    $result = $apiInstance->listPayments($limit, $skip, $sort);
    print_r($result);
} catch (Exception $e) {
    echo 'Exception when calling PaymentsApi->listPayments: ', $e->getMessage(), PHP_EOL;
}
```

### Parameters

| Name | Type | Description  | Notes |
| ------------- | ------------- | ------------- | ------------- |
| **limit** | **int**| Limit the number of payments to return, pagination can be achieved in conjunction with &#39;skip&#39; parameter. | [optional] |
| **skip** | **int**| How many payments to skip, pagination can be achieved in conjunction with &#39;limit&#39; parameter. | [optional] |
| **sort** | [**string[]**](../Model/string.md)| Field used to sort payments (Default is by date). | [optional] |

### Return type

[**\Formance\Model\ListPaymentsResponse**](../Model/ListPaymentsResponse.md)

### Authorization

[Authorization](../../README.md#Authorization)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`

[[Back to top]](#) [[Back to API list]](../../README.md#endpoints)
[[Back to Model list]](../../README.md#models)
[[Back to README]](../../README.md)

## `paymentslistAccounts()`

```php
paymentslistAccounts($limit, $skip, $sort): \Formance\Model\ListAccountsResponse
```

Returns a list of accounts.

### Example

```php
<?php
require_once(__DIR__ . '/vendor/autoload.php');


// Configure OAuth2 access token for authorization: Authorization
$config = Formance\Configuration::getDefaultConfiguration()->setAccessToken('YOUR_ACCESS_TOKEN');


$apiInstance = new Formance\Api\PaymentsApi(
    // If you want use custom http client, pass your client which implements `GuzzleHttp\ClientInterface`.
    // This is optional, `GuzzleHttp\Client` will be used as default.
    new GuzzleHttp\Client(),
    $config
);
$limit = 10; // int | Limit the number of accounts to return, pagination can be achieved in conjunction with 'skip' parameter.
$skip = 100; // int | How many accounts to skip, pagination can be achieved in conjunction with 'limit' parameter.
$sort = status; // string[] | Field used to sort payments (Default is by date).

try {
    $result = $apiInstance->paymentslistAccounts($limit, $skip, $sort);
    print_r($result);
} catch (Exception $e) {
    echo 'Exception when calling PaymentsApi->paymentslistAccounts: ', $e->getMessage(), PHP_EOL;
}
```

### Parameters

| Name | Type | Description  | Notes |
| ------------- | ------------- | ------------- | ------------- |
| **limit** | **int**| Limit the number of accounts to return, pagination can be achieved in conjunction with &#39;skip&#39; parameter. | [optional] |
| **skip** | **int**| How many accounts to skip, pagination can be achieved in conjunction with &#39;limit&#39; parameter. | [optional] |
| **sort** | [**string[]**](../Model/string.md)| Field used to sort payments (Default is by date). | [optional] |

### Return type

[**\Formance\Model\ListAccountsResponse**](../Model/ListAccountsResponse.md)

### Authorization

[Authorization](../../README.md#Authorization)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`

[[Back to top]](#) [[Back to API list]](../../README.md#endpoints)
[[Back to Model list]](../../README.md#models)
[[Back to README]](../../README.md)

## `readConnectorConfig()`

```php
readConnectorConfig($connector): \Formance\Model\ConnectorConfig
```

Read connector config

Read connector config

### Example

```php
<?php
require_once(__DIR__ . '/vendor/autoload.php');


// Configure OAuth2 access token for authorization: Authorization
$config = Formance\Configuration::getDefaultConfiguration()->setAccessToken('YOUR_ACCESS_TOKEN');


$apiInstance = new Formance\Api\PaymentsApi(
    // If you want use custom http client, pass your client which implements `GuzzleHttp\ClientInterface`.
    // This is optional, `GuzzleHttp\Client` will be used as default.
    new GuzzleHttp\Client(),
    $config
);
$connector = new \Formance\Model\Connectors(); // Connectors | The connector code

try {
    $result = $apiInstance->readConnectorConfig($connector);
    print_r($result);
} catch (Exception $e) {
    echo 'Exception when calling PaymentsApi->readConnectorConfig: ', $e->getMessage(), PHP_EOL;
}
```

### Parameters

| Name | Type | Description  | Notes |
| ------------- | ------------- | ------------- | ------------- |
| **connector** | [**Connectors**](../Model/.md)| The connector code | |

### Return type

[**\Formance\Model\ConnectorConfig**](../Model/ConnectorConfig.md)

### Authorization

[Authorization](../../README.md#Authorization)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`

[[Back to top]](#) [[Back to API list]](../../README.md#endpoints)
[[Back to Model list]](../../README.md#models)
[[Back to README]](../../README.md)

## `resetConnector()`

```php
resetConnector($connector)
```

Reset connector

Reset connector. Will remove the connector and ALL PAYMENTS generated with it.

### Example

```php
<?php
require_once(__DIR__ . '/vendor/autoload.php');


// Configure OAuth2 access token for authorization: Authorization
$config = Formance\Configuration::getDefaultConfiguration()->setAccessToken('YOUR_ACCESS_TOKEN');


$apiInstance = new Formance\Api\PaymentsApi(
    // If you want use custom http client, pass your client which implements `GuzzleHttp\ClientInterface`.
    // This is optional, `GuzzleHttp\Client` will be used as default.
    new GuzzleHttp\Client(),
    $config
);
$connector = new \Formance\Model\Connectors(); // Connectors | The connector code

try {
    $apiInstance->resetConnector($connector);
} catch (Exception $e) {
    echo 'Exception when calling PaymentsApi->resetConnector: ', $e->getMessage(), PHP_EOL;
}
```

### Parameters

| Name | Type | Description  | Notes |
| ------------- | ------------- | ------------- | ------------- |
| **connector** | [**Connectors**](../Model/.md)| The connector code | |

### Return type

void (empty response body)

### Authorization

[Authorization](../../README.md#Authorization)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../../README.md#endpoints)
[[Back to Model list]](../../README.md#models)
[[Back to README]](../../README.md)

## `uninstallConnector()`

```php
uninstallConnector($connector)
```

Uninstall connector

Uninstall  connector

### Example

```php
<?php
require_once(__DIR__ . '/vendor/autoload.php');


// Configure OAuth2 access token for authorization: Authorization
$config = Formance\Configuration::getDefaultConfiguration()->setAccessToken('YOUR_ACCESS_TOKEN');


$apiInstance = new Formance\Api\PaymentsApi(
    // If you want use custom http client, pass your client which implements `GuzzleHttp\ClientInterface`.
    // This is optional, `GuzzleHttp\Client` will be used as default.
    new GuzzleHttp\Client(),
    $config
);
$connector = new \Formance\Model\Connectors(); // Connectors | The connector code

try {
    $apiInstance->uninstallConnector($connector);
} catch (Exception $e) {
    echo 'Exception when calling PaymentsApi->uninstallConnector: ', $e->getMessage(), PHP_EOL;
}
```

### Parameters

| Name | Type | Description  | Notes |
| ------------- | ------------- | ------------- | ------------- |
| **connector** | [**Connectors**](../Model/.md)| The connector code | |

### Return type

void (empty response body)

### Authorization

[Authorization](../../README.md#Authorization)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../../README.md#endpoints)
[[Back to Model list]](../../README.md#models)
[[Back to README]](../../README.md)
