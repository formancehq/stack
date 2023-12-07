# orchestration

### Available Operations

* [cancelEvent](#cancelevent) - Cancel a running workflow
* [createTrigger](#createtrigger) - Create trigger
* [createWorkflow](#createworkflow) - Create workflow
* [deleteTrigger](#deletetrigger) - Delete trigger
* [deleteWorkflow](#deleteworkflow) - Delete a flow by id
* [getInstance](#getinstance) - Get a workflow instance by id
* [getInstanceHistory](#getinstancehistory) - Get a workflow instance history by id
* [getInstanceStageHistory](#getinstancestagehistory) - Get a workflow instance stage history
* [getWorkflow](#getworkflow) - Get a flow by id
* [listInstances](#listinstances) - List instances of a workflow
* [listTriggers](#listtriggers) - List triggers
* [listTriggersOccurrences](#listtriggersoccurrences) - List triggers occurrences
* [listWorkflows](#listworkflows) - List registered workflows
* [orchestrationgetServerInfo](#orchestrationgetserverinfo) - Get server info
* [readTrigger](#readtrigger) - Read trigger
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
    $request->instanceID = 'occaecati';

    $response = $sdk->orchestration->cancelEvent($request);

    if ($response->statusCode === 200) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

## createTrigger

Create trigger

### Example Usage

```php
<?php

declare(strict_types=1);
require_once 'vendor/autoload.php';

use \formance\stack\SDK;
use \formance\stack\Models\Shared\Security;
use \formance\stack\Models\Shared\TriggerData;

$sdk = SDK::builder()
    ->build();

try {
    $request = new TriggerData();
    $request->event = 'numquam';
    $request->filter = 'commodi';
    $request->vars = [
        'molestiae' => 'velit',
        'error' => 'quia',
    ];
    $request->workflowID = 'quis';

    $response = $sdk->orchestration->createTrigger($request);

    if ($response->createTriggerResponse !== null) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

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
    $request->name = 'Kayla O'Kon';
    $request->stages = [
        [
            'tenetur' => 'ipsam',
        ],
        [
            'possimus' => 'aut',
            'quasi' => 'error',
            'temporibus' => 'laborum',
        ],
        [
            'reiciendis' => 'voluptatibus',
        ],
        [
            'nihil' => 'praesentium',
            'voluptatibus' => 'ipsa',
            'omnis' => 'voluptate',
            'cum' => 'perferendis',
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

## deleteTrigger

Read trigger

### Example Usage

```php
<?php

declare(strict_types=1);
require_once 'vendor/autoload.php';

use \formance\stack\SDK;
use \formance\stack\Models\Shared\Security;
use \formance\stack\Models\Operations\DeleteTriggerRequest;

$sdk = SDK::builder()
    ->build();

try {
    $request = new DeleteTriggerRequest();
    $request->triggerID = 'doloremque';

    $response = $sdk->orchestration->deleteTrigger($request);

    if ($response->statusCode === 200) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

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
    $request->flowId = 'reprehenderit';

    $response = $sdk->orchestration->deleteWorkflow($request);

    if ($response->statusCode === 200) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

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
    $request->instanceID = 'ut';

    $response = $sdk->orchestration->getInstance($request);

    if ($response->getWorkflowInstanceResponse !== null) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

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
    $request->instanceID = 'maiores';

    $response = $sdk->orchestration->getInstanceHistory($request);

    if ($response->getWorkflowInstanceHistoryResponse !== null) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

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
    $request->instanceID = 'dicta';
    $request->number = 359444;

    $response = $sdk->orchestration->getInstanceStageHistory($request);

    if ($response->getWorkflowInstanceHistoryStageResponse !== null) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

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
    $request->flowId = 'dolore';

    $response = $sdk->orchestration->getWorkflow($request);

    if ($response->getWorkflowResponse !== null) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

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
    $request->workflowID = 'iusto';

    $response = $sdk->orchestration->listInstances($request);

    if ($response->listRunsResponse !== null) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

## listTriggers

List triggers

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
    $response = $sdk->orchestration->listTriggers();

    if ($response->listTriggersResponse !== null) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

## listTriggersOccurrences

List triggers occurrences

### Example Usage

```php
<?php

declare(strict_types=1);
require_once 'vendor/autoload.php';

use \formance\stack\SDK;
use \formance\stack\Models\Shared\Security;
use \formance\stack\Models\Operations\ListTriggersOccurrencesRequest;

$sdk = SDK::builder()
    ->build();

try {
    $request = new ListTriggersOccurrencesRequest();
    $request->triggerID = 'dicta';

    $response = $sdk->orchestration->listTriggersOccurrences($request);

    if ($response->listTriggersOccurrencesResponse !== null) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

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

## readTrigger

Read trigger

### Example Usage

```php
<?php

declare(strict_types=1);
require_once 'vendor/autoload.php';

use \formance\stack\SDK;
use \formance\stack\Models\Shared\Security;
use \formance\stack\Models\Operations\ReadTriggerRequest;

$sdk = SDK::builder()
    ->build();

try {
    $request = new ReadTriggerRequest();
    $request->triggerID = 'harum';

    $response = $sdk->orchestration->readTrigger($request);

    if ($response->readTriggerResponse !== null) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

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
        'accusamus' => 'commodi',
        'repudiandae' => 'quae',
    ];
    $request->wait = false;
    $request->workflowID = 'ipsum';

    $response = $sdk->orchestration->runWorkflow($request);

    if ($response->runWorkflowResponse !== null) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

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
    $request->requestBody->name = 'Virgil Mante';
    $request->instanceID = 'praesentium';

    $response = $sdk->orchestration->sendEvent($request);

    if ($response->statusCode === 200) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```
