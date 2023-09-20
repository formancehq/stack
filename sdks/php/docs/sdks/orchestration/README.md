# orchestration

### Available Operations

* [cancelEvent](#cancelevent) - Cancel a running workflow
* [createWorkflow](#createworkflow) - Create workflow
* [deleteWorkflow](#deleteworkflow) - Delete a flow by id
* [getInstance](#getinstance) - Get a workflow instance by id
* [getInstanceHistory](#getinstancehistory) - Get a workflow instance history by id
* [getInstanceStageHistory](#getinstancestagehistory) - Get a workflow instance stage history
* [getWorkflow](#getworkflow) - Get a flow by id
* [listInstances](#listinstances) - List instances of a workflow
* [listWorkflows](#listworkflows) - List registered workflows
* [orchestrationgetServerInfo](#orchestrationgetserverinfo) - Get server info
* [runWorkflow](#runworkflow) - Run workflow
* [sendEvent](#sendevent) - Send an event to a running workflow

## cancelEvent

Cancel a running workflow

### Example Usage

```php
<?php

declare(strict_types=1);
require_once 'vendor/autoload.php';

use \formance\stack\SDK;
use \formance\stack\Models\Shared\Security;
use \formance\stack\Models\Operations\CancelEventRequest;

$sdk = SDK::builder()
    ->build();

try {
    $request = new CancelEventRequest();
    $request->instanceID = 'modi';

    $response = $sdk->orchestration->cancelEvent($request);

    if ($response->statusCode === 200) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

### Parameters

| Parameter                                                                                             | Type                                                                                                  | Required                                                                                              | Description                                                                                           |
| ----------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------- |
| `$request`                                                                                            | [\formance\stack\Models\Operations\CancelEventRequest](../../models/operations/CancelEventRequest.md) | :heavy_check_mark:                                                                                    | The request object to use for the request.                                                            |


### Response

**[?\formance\stack\Models\Operations\CancelEventResponse](../../models/operations/CancelEventResponse.md)**


## createWorkflow

Create a workflow

### Example Usage

```php
<?php

declare(strict_types=1);
require_once 'vendor/autoload.php';

use \formance\stack\SDK;
use \formance\stack\Models\Shared\Security;
use \formance\stack\Models\Shared\CreateWorkflowRequest;

$sdk = SDK::builder()
    ->build();

try {
    $request = new CreateWorkflowRequest();
    $request->name = 'Dr. Jordan Von';
    $request->stages = [
        [
            'incidunt' => 'enim',
            'consequatur' => 'est',
            'quibusdam' => 'explicabo',
            'deserunt' => 'distinctio',
        ],
    ];

    $response = $sdk->orchestration->createWorkflow($request);

    if ($response->createWorkflowResponse !== null) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

### Parameters

| Parameter                                                                                           | Type                                                                                                | Required                                                                                            | Description                                                                                         |
| --------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------- |
| `$request`                                                                                          | [\formance\stack\Models\Shared\CreateWorkflowRequest](../../models/shared/CreateWorkflowRequest.md) | :heavy_check_mark:                                                                                  | The request object to use for the request.                                                          |


### Response

**[?\formance\stack\Models\Operations\CreateWorkflowResponse](../../models/operations/CreateWorkflowResponse.md)**


## deleteWorkflow

Delete a flow by id

### Example Usage

```php
<?php

declare(strict_types=1);
require_once 'vendor/autoload.php';

use \formance\stack\SDK;
use \formance\stack\Models\Shared\Security;
use \formance\stack\Models\Operations\DeleteWorkflowRequest;

$sdk = SDK::builder()
    ->build();

try {
    $request = new DeleteWorkflowRequest();
    $request->flowId = 'quibusdam';

    $response = $sdk->orchestration->deleteWorkflow($request);

    if ($response->statusCode === 200) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

### Parameters

| Parameter                                                                                                   | Type                                                                                                        | Required                                                                                                    | Description                                                                                                 |
| ----------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------- |
| `$request`                                                                                                  | [\formance\stack\Models\Operations\DeleteWorkflowRequest](../../models/operations/DeleteWorkflowRequest.md) | :heavy_check_mark:                                                                                          | The request object to use for the request.                                                                  |


### Response

**[?\formance\stack\Models\Operations\DeleteWorkflowResponse](../../models/operations/DeleteWorkflowResponse.md)**


## getInstance

Get a workflow instance by id

### Example Usage

```php
<?php

declare(strict_types=1);
require_once 'vendor/autoload.php';

use \formance\stack\SDK;
use \formance\stack\Models\Shared\Security;
use \formance\stack\Models\Operations\GetInstanceRequest;

$sdk = SDK::builder()
    ->build();

try {
    $request = new GetInstanceRequest();
    $request->instanceID = 'labore';

    $response = $sdk->orchestration->getInstance($request);

    if ($response->getWorkflowInstanceResponse !== null) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

### Parameters

| Parameter                                                                                             | Type                                                                                                  | Required                                                                                              | Description                                                                                           |
| ----------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------- |
| `$request`                                                                                            | [\formance\stack\Models\Operations\GetInstanceRequest](../../models/operations/GetInstanceRequest.md) | :heavy_check_mark:                                                                                    | The request object to use for the request.                                                            |


### Response

**[?\formance\stack\Models\Operations\GetInstanceResponse](../../models/operations/GetInstanceResponse.md)**


## getInstanceHistory

Get a workflow instance history by id

### Example Usage

```php
<?php

declare(strict_types=1);
require_once 'vendor/autoload.php';

use \formance\stack\SDK;
use \formance\stack\Models\Shared\Security;
use \formance\stack\Models\Operations\GetInstanceHistoryRequest;

$sdk = SDK::builder()
    ->build();

try {
    $request = new GetInstanceHistoryRequest();
    $request->instanceID = 'modi';

    $response = $sdk->orchestration->getInstanceHistory($request);

    if ($response->getWorkflowInstanceHistoryResponse !== null) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

### Parameters

| Parameter                                                                                                           | Type                                                                                                                | Required                                                                                                            | Description                                                                                                         |
| ------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------- |
| `$request`                                                                                                          | [\formance\stack\Models\Operations\GetInstanceHistoryRequest](../../models/operations/GetInstanceHistoryRequest.md) | :heavy_check_mark:                                                                                                  | The request object to use for the request.                                                                          |


### Response

**[?\formance\stack\Models\Operations\GetInstanceHistoryResponse](../../models/operations/GetInstanceHistoryResponse.md)**


## getInstanceStageHistory

Get a workflow instance stage history

### Example Usage

```php
<?php

declare(strict_types=1);
require_once 'vendor/autoload.php';

use \formance\stack\SDK;
use \formance\stack\Models\Shared\Security;
use \formance\stack\Models\Operations\GetInstanceStageHistoryRequest;

$sdk = SDK::builder()
    ->build();

try {
    $request = new GetInstanceStageHistoryRequest();
    $request->instanceID = 'qui';
    $request->number = 397821;

    $response = $sdk->orchestration->getInstanceStageHistory($request);

    if ($response->getWorkflowInstanceHistoryStageResponse !== null) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

### Parameters

| Parameter                                                                                                                     | Type                                                                                                                          | Required                                                                                                                      | Description                                                                                                                   |
| ----------------------------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------- |
| `$request`                                                                                                                    | [\formance\stack\Models\Operations\GetInstanceStageHistoryRequest](../../models/operations/GetInstanceStageHistoryRequest.md) | :heavy_check_mark:                                                                                                            | The request object to use for the request.                                                                                    |


### Response

**[?\formance\stack\Models\Operations\GetInstanceStageHistoryResponse](../../models/operations/GetInstanceStageHistoryResponse.md)**


## getWorkflow

Get a flow by id

### Example Usage

```php
<?php

declare(strict_types=1);
require_once 'vendor/autoload.php';

use \formance\stack\SDK;
use \formance\stack\Models\Shared\Security;
use \formance\stack\Models\Operations\GetWorkflowRequest;

$sdk = SDK::builder()
    ->build();

try {
    $request = new GetWorkflowRequest();
    $request->flowId = 'cupiditate';

    $response = $sdk->orchestration->getWorkflow($request);

    if ($response->getWorkflowResponse !== null) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

### Parameters

| Parameter                                                                                             | Type                                                                                                  | Required                                                                                              | Description                                                                                           |
| ----------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------- |
| `$request`                                                                                            | [\formance\stack\Models\Operations\GetWorkflowRequest](../../models/operations/GetWorkflowRequest.md) | :heavy_check_mark:                                                                                    | The request object to use for the request.                                                            |


### Response

**[?\formance\stack\Models\Operations\GetWorkflowResponse](../../models/operations/GetWorkflowResponse.md)**


## listInstances

List instances of a workflow

### Example Usage

```php
<?php

declare(strict_types=1);
require_once 'vendor/autoload.php';

use \formance\stack\SDK;
use \formance\stack\Models\Shared\Security;
use \formance\stack\Models\Operations\ListInstancesRequest;

$sdk = SDK::builder()
    ->build();

try {
    $request = new ListInstancesRequest();
    $request->running = false;
    $request->workflowID = 'quos';

    $response = $sdk->orchestration->listInstances($request);

    if ($response->listRunsResponse !== null) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

### Parameters

| Parameter                                                                                                 | Type                                                                                                      | Required                                                                                                  | Description                                                                                               |
| --------------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------- |
| `$request`                                                                                                | [\formance\stack\Models\Operations\ListInstancesRequest](../../models/operations/ListInstancesRequest.md) | :heavy_check_mark:                                                                                        | The request object to use for the request.                                                                |


### Response

**[?\formance\stack\Models\Operations\ListInstancesResponse](../../models/operations/ListInstancesResponse.md)**


## listWorkflows

List registered workflows

### Example Usage

```php
<?php

declare(strict_types=1);
require_once 'vendor/autoload.php';

use \formance\stack\SDK;
use \formance\stack\Models\Shared\Security;

$sdk = SDK::builder()
    ->build();

try {
    $response = $sdk->orchestration->listWorkflows();

    if ($response->listWorkflowsResponse !== null) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```


### Response

**[?\formance\stack\Models\Operations\ListWorkflowsResponse](../../models/operations/ListWorkflowsResponse.md)**


## orchestrationgetServerInfo

Get server info

### Example Usage

```php
<?php

declare(strict_types=1);
require_once 'vendor/autoload.php';

use \formance\stack\SDK;
use \formance\stack\Models\Shared\Security;

$sdk = SDK::builder()
    ->build();

try {
    $response = $sdk->orchestration->orchestrationgetServerInfo();

    if ($response->serverInfo !== null) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```


### Response

**[?\formance\stack\Models\Operations\OrchestrationgetServerInfoResponse](../../models/operations/OrchestrationgetServerInfoResponse.md)**


## runWorkflow

Run workflow

### Example Usage

```php
<?php

declare(strict_types=1);
require_once 'vendor/autoload.php';

use \formance\stack\SDK;
use \formance\stack\Models\Shared\Security;
use \formance\stack\Models\Operations\RunWorkflowRequest;

$sdk = SDK::builder()
    ->build();

try {
    $request = new RunWorkflowRequest();
    $request->requestBody = [
        'magni' => 'assumenda',
    ];
    $request->wait = false;
    $request->workflowID = 'ipsam';

    $response = $sdk->orchestration->runWorkflow($request);

    if ($response->runWorkflowResponse !== null) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

### Parameters

| Parameter                                                                                             | Type                                                                                                  | Required                                                                                              | Description                                                                                           |
| ----------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------- |
| `$request`                                                                                            | [\formance\stack\Models\Operations\RunWorkflowRequest](../../models/operations/RunWorkflowRequest.md) | :heavy_check_mark:                                                                                    | The request object to use for the request.                                                            |


### Response

**[?\formance\stack\Models\Operations\RunWorkflowResponse](../../models/operations/RunWorkflowResponse.md)**


## sendEvent

Send an event to a running workflow

### Example Usage

```php
<?php

declare(strict_types=1);
require_once 'vendor/autoload.php';

use \formance\stack\SDK;
use \formance\stack\Models\Shared\Security;
use \formance\stack\Models\Operations\SendEventRequest;
use \formance\stack\Models\Operations\SendEventRequestBody;

$sdk = SDK::builder()
    ->build();

try {
    $request = new SendEventRequest();
    $request->requestBody = new SendEventRequestBody();
    $request->requestBody->name = 'Denise Pagac';
    $request->instanceID = 'facilis';

    $response = $sdk->orchestration->sendEvent($request);

    if ($response->statusCode === 200) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

### Parameters

| Parameter                                                                                         | Type                                                                                              | Required                                                                                          | Description                                                                                       |
| ------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------- |
| `$request`                                                                                        | [\formance\stack\Models\Operations\SendEventRequest](../../models/operations/SendEventRequest.md) | :heavy_check_mark:                                                                                | The request object to use for the request.                                                        |


### Response

**[?\formance\stack\Models\Operations\SendEventResponse](../../models/operations/SendEventResponse.md)**

