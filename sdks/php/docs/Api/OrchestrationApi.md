# Formance\OrchestrationApi

All URIs are relative to http://localhost, except if the operation defines another base path.

| Method | HTTP request | Description |
| ------------- | ------------- | ------------- |
| [**createWorkflow()**](OrchestrationApi.md#createWorkflow) | **POST** /api/orchestration/workflows | Create workflow |
| [**getInstance()**](OrchestrationApi.md#getInstance) | **GET** /api/orchestration/instances/{instanceID} | Get a workflow instance by id |
| [**getInstanceHistory()**](OrchestrationApi.md#getInstanceHistory) | **GET** /api/orchestration/instances/{instanceID}/history | Get a workflow instance history by id |
| [**getInstanceStageHistory()**](OrchestrationApi.md#getInstanceStageHistory) | **GET** /api/orchestration/instances/{instanceID}/stages/{number}/history | Get a workflow instance stage history |
| [**getWorkflow()**](OrchestrationApi.md#getWorkflow) | **GET** /api/orchestration/workflows/{flowId} | Get a flow by id |
| [**listInstances()**](OrchestrationApi.md#listInstances) | **GET** /api/orchestration/instances | List instances of a workflow |
| [**listWorkflows()**](OrchestrationApi.md#listWorkflows) | **GET** /api/orchestration/workflows | List registered workflows |
| [**orchestrationgetServerInfo()**](OrchestrationApi.md#orchestrationgetServerInfo) | **GET** /api/orchestration/_info | Get server info |
| [**runWorkflow()**](OrchestrationApi.md#runWorkflow) | **POST** /api/orchestration/workflows/{workflowID}/instances | Run workflow |


## `createWorkflow()`

```php
createWorkflow($body): \Formance\Model\CreateWorkflowResponse
```

Create workflow

Create a workflow

### Example

```php
<?php
require_once(__DIR__ . '/vendor/autoload.php');


// Configure OAuth2 access token for authorization: Authorization
$config = Formance\Configuration::getDefaultConfiguration()->setAccessToken('YOUR_ACCESS_TOKEN');


$apiInstance = new Formance\Api\OrchestrationApi(
    // If you want use custom http client, pass your client which implements `GuzzleHttp\ClientInterface`.
    // This is optional, `GuzzleHttp\Client` will be used as default.
    new GuzzleHttp\Client(),
    $config
);
$body = new \Formance\Model\WorkflowConfig(); // \Formance\Model\WorkflowConfig

try {
    $result = $apiInstance->createWorkflow($body);
    print_r($result);
} catch (Exception $e) {
    echo 'Exception when calling OrchestrationApi->createWorkflow: ', $e->getMessage(), PHP_EOL;
}
```

### Parameters

| Name | Type | Description  | Notes |
| ------------- | ------------- | ------------- | ------------- |
| **body** | **\Formance\Model\WorkflowConfig**|  | [optional] |

### Return type

[**\Formance\Model\CreateWorkflowResponse**](../Model/CreateWorkflowResponse.md)

### Authorization

[Authorization](../../README.md#Authorization)

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`

[[Back to top]](#) [[Back to API list]](../../README.md#endpoints)
[[Back to Model list]](../../README.md#models)
[[Back to README]](../../README.md)

## `getInstance()`

```php
getInstance($instance_id): \Formance\Model\GetWorkflowInstanceResponse
```

Get a workflow instance by id

Get a workflow instance by id

### Example

```php
<?php
require_once(__DIR__ . '/vendor/autoload.php');


// Configure OAuth2 access token for authorization: Authorization
$config = Formance\Configuration::getDefaultConfiguration()->setAccessToken('YOUR_ACCESS_TOKEN');


$apiInstance = new Formance\Api\OrchestrationApi(
    // If you want use custom http client, pass your client which implements `GuzzleHttp\ClientInterface`.
    // This is optional, `GuzzleHttp\Client` will be used as default.
    new GuzzleHttp\Client(),
    $config
);
$instance_id = xxx; // string | The instance id

try {
    $result = $apiInstance->getInstance($instance_id);
    print_r($result);
} catch (Exception $e) {
    echo 'Exception when calling OrchestrationApi->getInstance: ', $e->getMessage(), PHP_EOL;
}
```

### Parameters

| Name | Type | Description  | Notes |
| ------------- | ------------- | ------------- | ------------- |
| **instance_id** | **string**| The instance id | |

### Return type

[**\Formance\Model\GetWorkflowInstanceResponse**](../Model/GetWorkflowInstanceResponse.md)

### Authorization

[Authorization](../../README.md#Authorization)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`

[[Back to top]](#) [[Back to API list]](../../README.md#endpoints)
[[Back to Model list]](../../README.md#models)
[[Back to README]](../../README.md)

## `getInstanceHistory()`

```php
getInstanceHistory($instance_id): \Formance\Model\GetWorkflowInstanceHistoryResponse
```

Get a workflow instance history by id

Get a workflow instance history by id

### Example

```php
<?php
require_once(__DIR__ . '/vendor/autoload.php');


// Configure OAuth2 access token for authorization: Authorization
$config = Formance\Configuration::getDefaultConfiguration()->setAccessToken('YOUR_ACCESS_TOKEN');


$apiInstance = new Formance\Api\OrchestrationApi(
    // If you want use custom http client, pass your client which implements `GuzzleHttp\ClientInterface`.
    // This is optional, `GuzzleHttp\Client` will be used as default.
    new GuzzleHttp\Client(),
    $config
);
$instance_id = xxx; // string | The instance id

try {
    $result = $apiInstance->getInstanceHistory($instance_id);
    print_r($result);
} catch (Exception $e) {
    echo 'Exception when calling OrchestrationApi->getInstanceHistory: ', $e->getMessage(), PHP_EOL;
}
```

### Parameters

| Name | Type | Description  | Notes |
| ------------- | ------------- | ------------- | ------------- |
| **instance_id** | **string**| The instance id | |

### Return type

[**\Formance\Model\GetWorkflowInstanceHistoryResponse**](../Model/GetWorkflowInstanceHistoryResponse.md)

### Authorization

[Authorization](../../README.md#Authorization)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`

[[Back to top]](#) [[Back to API list]](../../README.md#endpoints)
[[Back to Model list]](../../README.md#models)
[[Back to README]](../../README.md)

## `getInstanceStageHistory()`

```php
getInstanceStageHistory($instance_id, $number): \Formance\Model\GetWorkflowInstanceHistoryStageResponse
```

Get a workflow instance stage history

Get a workflow instance stage history

### Example

```php
<?php
require_once(__DIR__ . '/vendor/autoload.php');


// Configure OAuth2 access token for authorization: Authorization
$config = Formance\Configuration::getDefaultConfiguration()->setAccessToken('YOUR_ACCESS_TOKEN');


$apiInstance = new Formance\Api\OrchestrationApi(
    // If you want use custom http client, pass your client which implements `GuzzleHttp\ClientInterface`.
    // This is optional, `GuzzleHttp\Client` will be used as default.
    new GuzzleHttp\Client(),
    $config
);
$instance_id = xxx; // string | The instance id
$number = 0; // int | The stage number

try {
    $result = $apiInstance->getInstanceStageHistory($instance_id, $number);
    print_r($result);
} catch (Exception $e) {
    echo 'Exception when calling OrchestrationApi->getInstanceStageHistory: ', $e->getMessage(), PHP_EOL;
}
```

### Parameters

| Name | Type | Description  | Notes |
| ------------- | ------------- | ------------- | ------------- |
| **instance_id** | **string**| The instance id | |
| **number** | **int**| The stage number | |

### Return type

[**\Formance\Model\GetWorkflowInstanceHistoryStageResponse**](../Model/GetWorkflowInstanceHistoryStageResponse.md)

### Authorization

[Authorization](../../README.md#Authorization)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`

[[Back to top]](#) [[Back to API list]](../../README.md#endpoints)
[[Back to Model list]](../../README.md#models)
[[Back to README]](../../README.md)

## `getWorkflow()`

```php
getWorkflow($flow_id): \Formance\Model\GetWorkflowResponse
```

Get a flow by id

Get a flow by id

### Example

```php
<?php
require_once(__DIR__ . '/vendor/autoload.php');


// Configure OAuth2 access token for authorization: Authorization
$config = Formance\Configuration::getDefaultConfiguration()->setAccessToken('YOUR_ACCESS_TOKEN');


$apiInstance = new Formance\Api\OrchestrationApi(
    // If you want use custom http client, pass your client which implements `GuzzleHttp\ClientInterface`.
    // This is optional, `GuzzleHttp\Client` will be used as default.
    new GuzzleHttp\Client(),
    $config
);
$flow_id = xxx; // string | The flow id

try {
    $result = $apiInstance->getWorkflow($flow_id);
    print_r($result);
} catch (Exception $e) {
    echo 'Exception when calling OrchestrationApi->getWorkflow: ', $e->getMessage(), PHP_EOL;
}
```

### Parameters

| Name | Type | Description  | Notes |
| ------------- | ------------- | ------------- | ------------- |
| **flow_id** | **string**| The flow id | |

### Return type

[**\Formance\Model\GetWorkflowResponse**](../Model/GetWorkflowResponse.md)

### Authorization

[Authorization](../../README.md#Authorization)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`

[[Back to top]](#) [[Back to API list]](../../README.md#endpoints)
[[Back to Model list]](../../README.md#models)
[[Back to README]](../../README.md)

## `listInstances()`

```php
listInstances($workflow_id): \Formance\Model\ListRunsResponse
```

List instances of a workflow

List instances of a workflow

### Example

```php
<?php
require_once(__DIR__ . '/vendor/autoload.php');


// Configure OAuth2 access token for authorization: Authorization
$config = Formance\Configuration::getDefaultConfiguration()->setAccessToken('YOUR_ACCESS_TOKEN');


$apiInstance = new Formance\Api\OrchestrationApi(
    // If you want use custom http client, pass your client which implements `GuzzleHttp\ClientInterface`.
    // This is optional, `GuzzleHttp\Client` will be used as default.
    new GuzzleHttp\Client(),
    $config
);
$workflow_id = xxx; // string | A workflow id

try {
    $result = $apiInstance->listInstances($workflow_id);
    print_r($result);
} catch (Exception $e) {
    echo 'Exception when calling OrchestrationApi->listInstances: ', $e->getMessage(), PHP_EOL;
}
```

### Parameters

| Name | Type | Description  | Notes |
| ------------- | ------------- | ------------- | ------------- |
| **workflow_id** | **string**| A workflow id | |

### Return type

[**\Formance\Model\ListRunsResponse**](../Model/ListRunsResponse.md)

### Authorization

[Authorization](../../README.md#Authorization)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`

[[Back to top]](#) [[Back to API list]](../../README.md#endpoints)
[[Back to Model list]](../../README.md#models)
[[Back to README]](../../README.md)

## `listWorkflows()`

```php
listWorkflows(): \Formance\Model\ListWorkflowsResponse
```

List registered workflows

List registered workflows

### Example

```php
<?php
require_once(__DIR__ . '/vendor/autoload.php');


// Configure OAuth2 access token for authorization: Authorization
$config = Formance\Configuration::getDefaultConfiguration()->setAccessToken('YOUR_ACCESS_TOKEN');


$apiInstance = new Formance\Api\OrchestrationApi(
    // If you want use custom http client, pass your client which implements `GuzzleHttp\ClientInterface`.
    // This is optional, `GuzzleHttp\Client` will be used as default.
    new GuzzleHttp\Client(),
    $config
);

try {
    $result = $apiInstance->listWorkflows();
    print_r($result);
} catch (Exception $e) {
    echo 'Exception when calling OrchestrationApi->listWorkflows: ', $e->getMessage(), PHP_EOL;
}
```

### Parameters

This endpoint does not need any parameter.

### Return type

[**\Formance\Model\ListWorkflowsResponse**](../Model/ListWorkflowsResponse.md)

### Authorization

[Authorization](../../README.md#Authorization)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`

[[Back to top]](#) [[Back to API list]](../../README.md#endpoints)
[[Back to Model list]](../../README.md#models)
[[Back to README]](../../README.md)

## `orchestrationgetServerInfo()`

```php
orchestrationgetServerInfo(): \Formance\Model\ServerInfo
```

Get server info

### Example

```php
<?php
require_once(__DIR__ . '/vendor/autoload.php');


// Configure OAuth2 access token for authorization: Authorization
$config = Formance\Configuration::getDefaultConfiguration()->setAccessToken('YOUR_ACCESS_TOKEN');


$apiInstance = new Formance\Api\OrchestrationApi(
    // If you want use custom http client, pass your client which implements `GuzzleHttp\ClientInterface`.
    // This is optional, `GuzzleHttp\Client` will be used as default.
    new GuzzleHttp\Client(),
    $config
);

try {
    $result = $apiInstance->orchestrationgetServerInfo();
    print_r($result);
} catch (Exception $e) {
    echo 'Exception when calling OrchestrationApi->orchestrationgetServerInfo: ', $e->getMessage(), PHP_EOL;
}
```

### Parameters

This endpoint does not need any parameter.

### Return type

[**\Formance\Model\ServerInfo**](../Model/ServerInfo.md)

### Authorization

[Authorization](../../README.md#Authorization)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`

[[Back to top]](#) [[Back to API list]](../../README.md#endpoints)
[[Back to Model list]](../../README.md#models)
[[Back to README]](../../README.md)

## `runWorkflow()`

```php
runWorkflow($workflow_id, $wait, $request_body): \Formance\Model\RunWorkflowResponse
```

Run workflow

Run workflow

### Example

```php
<?php
require_once(__DIR__ . '/vendor/autoload.php');


// Configure OAuth2 access token for authorization: Authorization
$config = Formance\Configuration::getDefaultConfiguration()->setAccessToken('YOUR_ACCESS_TOKEN');


$apiInstance = new Formance\Api\OrchestrationApi(
    // If you want use custom http client, pass your client which implements `GuzzleHttp\ClientInterface`.
    // This is optional, `GuzzleHttp\Client` will be used as default.
    new GuzzleHttp\Client(),
    $config
);
$workflow_id = xxx; // string | The flow id
$wait = True; // bool | Wait end of the workflow before return
$request_body = array('key' => 'request_body_example'); // array<string,string>

try {
    $result = $apiInstance->runWorkflow($workflow_id, $wait, $request_body);
    print_r($result);
} catch (Exception $e) {
    echo 'Exception when calling OrchestrationApi->runWorkflow: ', $e->getMessage(), PHP_EOL;
}
```

### Parameters

| Name | Type | Description  | Notes |
| ------------- | ------------- | ------------- | ------------- |
| **workflow_id** | **string**| The flow id | |
| **wait** | **bool**| Wait end of the workflow before return | [optional] |
| **request_body** | [**array<string,string>**](../Model/string.md)|  | [optional] |

### Return type

[**\Formance\Model\RunWorkflowResponse**](../Model/RunWorkflowResponse.md)

### Authorization

[Authorization](../../README.md#Authorization)

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`

[[Back to top]](#) [[Back to API list]](../../README.md#endpoints)
[[Back to Model list]](../../README.md#models)
[[Back to README]](../../README.md)
