# search

### Available Operations

* [search](#search) - Search
* [searchgetServerInfo](#searchgetserverinfo) - Get server info

## search

ElasticSearch query engine

### Example Usage

```php
<?php

declare(strict_types=1);
require_once 'vendor/autoload.php';

use \formance\stack\SDK;
use \formance\stack\Models\Shared\Security;
use \formance\stack\Models\Shared\Query;
use \formance\stack\Models\Shared\QueryRaw;

$sdk = SDK::builder()
    ->build();

try {
    $request = new Query();
    $request->after = [
        'users:002',
        'users:002',
        'users:002',
        'users:002',
    ];
    $request->cursor = 'YXVsdCBhbmQgYSBtYXhpbXVtIG1heF9yZXN1bHRzLol=';
    $request->ledgers = [
        'quickstart',
        'quickstart',
    ];
    $request->pageSize = 318569;
    $request->policy = 'OR';
    $request->raw = new QueryRaw();
    $request->sort = 'txid:asc';
    $request->target = 'consequatur';
    $request->terms = [
        'destination=central_bank1',
        'destination=central_bank1',
        'destination=central_bank1',
    ];

    $response = $sdk->search->search($request);

    if ($response->response !== null) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

### Parameters

| Parameter                                                           | Type                                                                | Required                                                            | Description                                                         |
| ------------------------------------------------------------------- | ------------------------------------------------------------------- | ------------------------------------------------------------------- | ------------------------------------------------------------------- |
| `$request`                                                          | [\formance\stack\Models\Shared\Query](../../models/shared/Query.md) | :heavy_check_mark:                                                  | The request object to use for the request.                          |


### Response

**[?\formance\stack\Models\Operations\SearchResponse](../../models/operations/SearchResponse.md)**


## searchgetServerInfo

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
    $response = $sdk->search->searchgetServerInfo();

    if ($response->serverInfo !== null) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```


### Response

**[?\formance\stack\Models\Operations\SearchgetServerInfoResponse](../../models/operations/SearchgetServerInfoResponse.md)**

