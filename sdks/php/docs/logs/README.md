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
    $request->cursor = 'aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ==';
    $request->endTime = DateTime::createFromFormat('Y-m-d\TH:i:sP', '2022-09-04T08:35:09.957Z');
    $request->ledger = 'ledger001';
    $request->pageSize = 570197;
    $request->startTime = DateTime::createFromFormat('Y-m-d\TH:i:sP', '2022-07-24T21:51:02.112Z');

    $response = $sdk->logs->listLogs($request);

    if ($response->logsCursorResponse !== null) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```
