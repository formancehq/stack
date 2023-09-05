# stats

### Available Operations

* [readStats](#readstats) - Get statistics from a ledger

## readStats

Get statistics from a ledger. (aggregate metrics on accounts and transactions)


### Example Usage

```php
<?php

declare(strict_types=1);
require_once 'vendor/autoload.php';

use \formance\stack\SDK;
use \formance\stack\Models\Shared\Security;
use \formance\stack\Models\Operations\ReadStatsRequest;

$sdk = SDK::builder()
    ->build();

try {
    $request = new ReadStatsRequest();
    $request->ledger = 'ledger001';

    $response = $sdk->stats->readStats($request);

    if ($response->statsResponse !== null) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

### Parameters

| Parameter                                                                                         | Type                                                                                              | Required                                                                                          | Description                                                                                       |
| ------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------- |
| `$request`                                                                                        | [\formance\stack\Models\Operations\ReadStatsRequest](../../models/operations/ReadStatsRequest.md) | :heavy_check_mark:                                                                                | The request object to use for the request.                                                        |


### Response

**[?\formance\stack\Models\Operations\ReadStatsResponse](../../models/operations/ReadStatsResponse.md)**

