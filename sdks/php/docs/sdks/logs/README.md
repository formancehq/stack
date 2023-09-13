# logs

### Available Operations

* [listLogs](#listlogs) - List the logs from a ledger

## listLogs

List the logs from a ledger, sorted by ID in descending order.

### Example Usage

```php
<?php

declare(strict_types=1);
require_once 'vendor/autoload.php';

use \formance\stack\SDK;
use \formance\stack\Models\Shared\Security;
use \formance\stack\Models\Operations\ListLogsRequest;

$sdk = SDK::builder()
    ->build();

try {
    $request = new ListLogsRequest();
    $request->after = 'corporis';
    $request->cursor = 'aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ==';
    $request->endTime = DateTime::createFromFormat('Y-m-d\TH:i:sP', '2022-04-01T23:59:21.675Z');
    $request->ledger = 'ledger001';
    $request->pageSize = 315428;
    $request->startTime = DateTime::createFromFormat('Y-m-d\TH:i:sP', '2022-04-10T11:47:13.463Z');

    $response = $sdk->logs->listLogs($request);

    if ($response->logsCursorResponse !== null) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

### Parameters

| Parameter                                                                                       | Type                                                                                            | Required                                                                                        | Description                                                                                     |
| ----------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------- |
| `$request`                                                                                      | [\formance\stack\Models\Operations\ListLogsRequest](../../models/operations/ListLogsRequest.md) | :heavy_check_mark:                                                                              | The request object to use for the request.                                                      |


### Response

**[?\formance\stack\Models\Operations\ListLogsResponse](../../models/operations/ListLogsResponse.md)**

