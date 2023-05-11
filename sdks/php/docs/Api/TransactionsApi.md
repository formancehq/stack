# Formance\TransactionsApi

All URIs are relative to http://localhost, except if the operation defines another base path.

| Method | HTTP request | Description |
| ------------- | ------------- | ------------- |
| [**addMetadataOnTransaction()**](TransactionsApi.md#addMetadataOnTransaction) | **POST** /api/ledger/{ledger}/transactions/{txid}/metadata | Set the metadata of a transaction by its ID |
| [**countTransactions()**](TransactionsApi.md#countTransactions) | **HEAD** /api/ledger/{ledger}/transactions | Count the transactions from a ledger |
| [**createTransaction()**](TransactionsApi.md#createTransaction) | **POST** /api/ledger/{ledger}/transactions | Create a new transaction to a ledger |
| [**getTransaction()**](TransactionsApi.md#getTransaction) | **GET** /api/ledger/{ledger}/transactions/{txid} | Get transaction from a ledger by its ID |
| [**listTransactions()**](TransactionsApi.md#listTransactions) | **GET** /api/ledger/{ledger}/transactions | List transactions from a ledger |
| [**revertTransaction()**](TransactionsApi.md#revertTransaction) | **POST** /api/ledger/{ledger}/transactions/{txid}/revert | Revert a ledger transaction by its ID |


## `addMetadataOnTransaction()`

```php
addMetadataOnTransaction($ledger, $txid, $request_body)
```

Set the metadata of a transaction by its ID

### Example

```php
<?php
require_once(__DIR__ . '/vendor/autoload.php');


// Configure OAuth2 access token for authorization: Authorization
$config = Formance\Configuration::getDefaultConfiguration()->setAccessToken('YOUR_ACCESS_TOKEN');


$apiInstance = new Formance\Api\TransactionsApi(
    // If you want use custom http client, pass your client which implements `GuzzleHttp\ClientInterface`.
    // This is optional, `GuzzleHttp\Client` will be used as default.
    new GuzzleHttp\Client(),
    $config
);
$ledger = ledger001; // string | Name of the ledger.
$txid = 1234; // int | Transaction ID.
$request_body = array('key' => 'request_body_example'); // array<string,string> | metadata

try {
    $apiInstance->addMetadataOnTransaction($ledger, $txid, $request_body);
} catch (Exception $e) {
    echo 'Exception when calling TransactionsApi->addMetadataOnTransaction: ', $e->getMessage(), PHP_EOL;
}
```

### Parameters

| Name | Type | Description  | Notes |
| ------------- | ------------- | ------------- | ------------- |
| **ledger** | **string**| Name of the ledger. | |
| **txid** | **int**| Transaction ID. | |
| **request_body** | [**array<string,string>**](../Model/string.md)| metadata | [optional] |

### Return type

void (empty response body)

### Authorization

[Authorization](../../README.md#Authorization)

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`

[[Back to top]](#) [[Back to API list]](../../README.md#endpoints)
[[Back to Model list]](../../README.md#models)
[[Back to README]](../../README.md)

## `countTransactions()`

```php
countTransactions($ledger, $reference, $account, $source, $destination, $start_time, $end_time, $metadata)
```

Count the transactions from a ledger

### Example

```php
<?php
require_once(__DIR__ . '/vendor/autoload.php');


// Configure OAuth2 access token for authorization: Authorization
$config = Formance\Configuration::getDefaultConfiguration()->setAccessToken('YOUR_ACCESS_TOKEN');


$apiInstance = new Formance\Api\TransactionsApi(
    // If you want use custom http client, pass your client which implements `GuzzleHttp\ClientInterface`.
    // This is optional, `GuzzleHttp\Client` will be used as default.
    new GuzzleHttp\Client(),
    $config
);
$ledger = ledger001; // string | Name of the ledger.
$reference = ref:001; // string | Filter transactions by reference field.
$account = users:001; // string | Filter transactions with postings involving given account, either as source or destination (regular expression placed between ^ and $).
$source = users:001; // string | Filter transactions with postings involving given account at source (regular expression placed between ^ and $).
$destination = users:001; // string | Filter transactions with postings involving given account at destination (regular expression placed between ^ and $).
$start_time = new \DateTime("2013-10-20T19:20:30+01:00"); // \DateTime | Filter transactions that occurred after this timestamp. The format is RFC3339 and is inclusive (for example, \"2023-01-02T15:04:01Z\" includes the first second of 4th minute).
$end_time = new \DateTime("2013-10-20T19:20:30+01:00"); // \DateTime | Filter transactions that occurred before this timestamp. The format is RFC3339 and is exclusive (for example, \"2023-01-02T15:04:01Z\" excludes the first second of 4th minute).
$metadata = {"admin":"true"}; // array<string,string> | Filter transactions by metadata key value pairs. Nested objects can be used as seen in the example below.

try {
    $apiInstance->countTransactions($ledger, $reference, $account, $source, $destination, $start_time, $end_time, $metadata);
} catch (Exception $e) {
    echo 'Exception when calling TransactionsApi->countTransactions: ', $e->getMessage(), PHP_EOL;
}
```

### Parameters

| Name | Type | Description  | Notes |
| ------------- | ------------- | ------------- | ------------- |
| **ledger** | **string**| Name of the ledger. | |
| **reference** | **string**| Filter transactions by reference field. | [optional] |
| **account** | **string**| Filter transactions with postings involving given account, either as source or destination (regular expression placed between ^ and $). | [optional] |
| **source** | **string**| Filter transactions with postings involving given account at source (regular expression placed between ^ and $). | [optional] |
| **destination** | **string**| Filter transactions with postings involving given account at destination (regular expression placed between ^ and $). | [optional] |
| **start_time** | **\DateTime**| Filter transactions that occurred after this timestamp. The format is RFC3339 and is inclusive (for example, \&quot;2023-01-02T15:04:01Z\&quot; includes the first second of 4th minute). | [optional] |
| **end_time** | **\DateTime**| Filter transactions that occurred before this timestamp. The format is RFC3339 and is exclusive (for example, \&quot;2023-01-02T15:04:01Z\&quot; excludes the first second of 4th minute). | [optional] |
| **metadata** | [**array<string,string>**](../Model/string.md)| Filter transactions by metadata key value pairs. Nested objects can be used as seen in the example below. | [optional] |

### Return type

void (empty response body)

### Authorization

[Authorization](../../README.md#Authorization)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`

[[Back to top]](#) [[Back to API list]](../../README.md#endpoints)
[[Back to Model list]](../../README.md#models)
[[Back to README]](../../README.md)

## `createTransaction()`

```php
createTransaction($ledger, $post_transaction, $preview): \Formance\Model\TransactionResponse
```

Create a new transaction to a ledger

### Example

```php
<?php
require_once(__DIR__ . '/vendor/autoload.php');


// Configure OAuth2 access token for authorization: Authorization
$config = Formance\Configuration::getDefaultConfiguration()->setAccessToken('YOUR_ACCESS_TOKEN');


$apiInstance = new Formance\Api\TransactionsApi(
    // If you want use custom http client, pass your client which implements `GuzzleHttp\ClientInterface`.
    // This is optional, `GuzzleHttp\Client` will be used as default.
    new GuzzleHttp\Client(),
    $config
);
$ledger = ledger001; // string | Name of the ledger.
$post_transaction = new \Formance\Model\PostTransaction(); // \Formance\Model\PostTransaction | The request body must contain at least one of the following objects:   - `postings`: suitable for simple transactions   - `script`: enabling more complex transactions with Numscript
$preview = true; // bool | Set the preview mode. Preview mode doesn't add the logs to the database or publish a message to the message broker.

try {
    $result = $apiInstance->createTransaction($ledger, $post_transaction, $preview);
    print_r($result);
} catch (Exception $e) {
    echo 'Exception when calling TransactionsApi->createTransaction: ', $e->getMessage(), PHP_EOL;
}
```

### Parameters

| Name | Type | Description  | Notes |
| ------------- | ------------- | ------------- | ------------- |
| **ledger** | **string**| Name of the ledger. | |
| **post_transaction** | [**\Formance\Model\PostTransaction**](../Model/PostTransaction.md)| The request body must contain at least one of the following objects:   - &#x60;postings&#x60;: suitable for simple transactions   - &#x60;script&#x60;: enabling more complex transactions with Numscript | |
| **preview** | **bool**| Set the preview mode. Preview mode doesn&#39;t add the logs to the database or publish a message to the message broker. | [optional] |

### Return type

[**\Formance\Model\TransactionResponse**](../Model/TransactionResponse.md)

### Authorization

[Authorization](../../README.md#Authorization)

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`

[[Back to top]](#) [[Back to API list]](../../README.md#endpoints)
[[Back to Model list]](../../README.md#models)
[[Back to README]](../../README.md)

## `getTransaction()`

```php
getTransaction($ledger, $txid): \Formance\Model\TransactionResponse
```

Get transaction from a ledger by its ID

### Example

```php
<?php
require_once(__DIR__ . '/vendor/autoload.php');


// Configure OAuth2 access token for authorization: Authorization
$config = Formance\Configuration::getDefaultConfiguration()->setAccessToken('YOUR_ACCESS_TOKEN');


$apiInstance = new Formance\Api\TransactionsApi(
    // If you want use custom http client, pass your client which implements `GuzzleHttp\ClientInterface`.
    // This is optional, `GuzzleHttp\Client` will be used as default.
    new GuzzleHttp\Client(),
    $config
);
$ledger = ledger001; // string | Name of the ledger.
$txid = 1234; // int | Transaction ID.

try {
    $result = $apiInstance->getTransaction($ledger, $txid);
    print_r($result);
} catch (Exception $e) {
    echo 'Exception when calling TransactionsApi->getTransaction: ', $e->getMessage(), PHP_EOL;
}
```

### Parameters

| Name | Type | Description  | Notes |
| ------------- | ------------- | ------------- | ------------- |
| **ledger** | **string**| Name of the ledger. | |
| **txid** | **int**| Transaction ID. | |

### Return type

[**\Formance\Model\TransactionResponse**](../Model/TransactionResponse.md)

### Authorization

[Authorization](../../README.md#Authorization)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`

[[Back to top]](#) [[Back to API list]](../../README.md#endpoints)
[[Back to Model list]](../../README.md#models)
[[Back to README]](../../README.md)

## `listTransactions()`

```php
listTransactions($ledger, $page_size, $reference, $account, $source, $destination, $start_time, $end_time, $cursor, $metadata): \Formance\Model\TransactionsCursorResponse
```

List transactions from a ledger

List transactions from a ledger, sorted by txid in descending order.

### Example

```php
<?php
require_once(__DIR__ . '/vendor/autoload.php');


// Configure OAuth2 access token for authorization: Authorization
$config = Formance\Configuration::getDefaultConfiguration()->setAccessToken('YOUR_ACCESS_TOKEN');


$apiInstance = new Formance\Api\TransactionsApi(
    // If you want use custom http client, pass your client which implements `GuzzleHttp\ClientInterface`.
    // This is optional, `GuzzleHttp\Client` will be used as default.
    new GuzzleHttp\Client(),
    $config
);
$ledger = ledger001; // string | Name of the ledger.
$page_size = 100; // int | The maximum number of results to return per page.
$reference = ref:001; // string | Find transactions by reference field.
$account = users:001; // string | Filter transactions with postings involving given account, either as source or destination (regular expression placed between ^ and $).
$source = users:001; // string | Filter transactions with postings involving given account at source (regular expression placed between ^ and $).
$destination = users:001; // string | Filter transactions with postings involving given account at destination (regular expression placed between ^ and $).
$start_time = new \DateTime("2013-10-20T19:20:30+01:00"); // \DateTime | Filter transactions that occurred after this timestamp. The format is RFC3339 and is inclusive (for example, \"2023-01-02T15:04:01Z\" includes the first second of 4th minute).
$end_time = new \DateTime("2013-10-20T19:20:30+01:00"); // \DateTime | Filter transactions that occurred before this timestamp. The format is RFC3339 and is exclusive (for example, \"2023-01-02T15:04:01Z\" excludes the first second of 4th minute).
$cursor = aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ==; // string | Parameter used in pagination requests. Maximum page size is set to 15. Set to the value of next for the next page of results. Set to the value of previous for the previous page of results. No other parameters can be set when this parameter is set.
$metadata = {"admin":"true"}; // array<string,string> | Filter transactions by metadata key value pairs.

try {
    $result = $apiInstance->listTransactions($ledger, $page_size, $reference, $account, $source, $destination, $start_time, $end_time, $cursor, $metadata);
    print_r($result);
} catch (Exception $e) {
    echo 'Exception when calling TransactionsApi->listTransactions: ', $e->getMessage(), PHP_EOL;
}
```

### Parameters

| Name | Type | Description  | Notes |
| ------------- | ------------- | ------------- | ------------- |
| **ledger** | **string**| Name of the ledger. | |
| **page_size** | **int**| The maximum number of results to return per page. | [optional] |
| **reference** | **string**| Find transactions by reference field. | [optional] |
| **account** | **string**| Filter transactions with postings involving given account, either as source or destination (regular expression placed between ^ and $). | [optional] |
| **source** | **string**| Filter transactions with postings involving given account at source (regular expression placed between ^ and $). | [optional] |
| **destination** | **string**| Filter transactions with postings involving given account at destination (regular expression placed between ^ and $). | [optional] |
| **start_time** | **\DateTime**| Filter transactions that occurred after this timestamp. The format is RFC3339 and is inclusive (for example, \&quot;2023-01-02T15:04:01Z\&quot; includes the first second of 4th minute). | [optional] |
| **end_time** | **\DateTime**| Filter transactions that occurred before this timestamp. The format is RFC3339 and is exclusive (for example, \&quot;2023-01-02T15:04:01Z\&quot; excludes the first second of 4th minute). | [optional] |
| **cursor** | **string**| Parameter used in pagination requests. Maximum page size is set to 15. Set to the value of next for the next page of results. Set to the value of previous for the previous page of results. No other parameters can be set when this parameter is set. | [optional] |
| **metadata** | [**array<string,string>**](../Model/string.md)| Filter transactions by metadata key value pairs. | [optional] |

### Return type

[**\Formance\Model\TransactionsCursorResponse**](../Model/TransactionsCursorResponse.md)

### Authorization

[Authorization](../../README.md#Authorization)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`

[[Back to top]](#) [[Back to API list]](../../README.md#endpoints)
[[Back to Model list]](../../README.md#models)
[[Back to README]](../../README.md)

## `revertTransaction()`

```php
revertTransaction($ledger, $txid): \Formance\Model\TransactionResponse
```

Revert a ledger transaction by its ID

### Example

```php
<?php
require_once(__DIR__ . '/vendor/autoload.php');


// Configure OAuth2 access token for authorization: Authorization
$config = Formance\Configuration::getDefaultConfiguration()->setAccessToken('YOUR_ACCESS_TOKEN');


$apiInstance = new Formance\Api\TransactionsApi(
    // If you want use custom http client, pass your client which implements `GuzzleHttp\ClientInterface`.
    // This is optional, `GuzzleHttp\Client` will be used as default.
    new GuzzleHttp\Client(),
    $config
);
$ledger = ledger001; // string | Name of the ledger.
$txid = 1234; // int | Transaction ID.

try {
    $result = $apiInstance->revertTransaction($ledger, $txid);
    print_r($result);
} catch (Exception $e) {
    echo 'Exception when calling TransactionsApi->revertTransaction: ', $e->getMessage(), PHP_EOL;
}
```

### Parameters

| Name | Type | Description  | Notes |
| ------------- | ------------- | ------------- | ------------- |
| **ledger** | **string**| Name of the ledger. | |
| **txid** | **int**| Transaction ID. | |

### Return type

[**\Formance\Model\TransactionResponse**](../Model/TransactionResponse.md)

### Authorization

[Authorization](../../README.md#Authorization)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`

[[Back to top]](#) [[Back to API list]](../../README.md#endpoints)
[[Back to Model list]](../../README.md#models)
[[Back to README]](../../README.md)
