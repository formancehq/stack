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

$sdk = SDK::builder()
    ->build();

try {
    $request = new Query();
    $request->after = [
        'users:002',
        'users:002',
        'users:002',
    ];
    $request->cursor = 'YXVsdCBhbmQgYSBtYXhpbXVtIG1heF9yZXN1bHRzLol=';
    $request->ledgers = [
        'quickstart',
        'quickstart',
        'quickstart',
    ];
    $request->pageSize = 725255;
    $request->policy = 'OR';
    $request->raw = [
        'blanditiis' => 'deleniti',
        'sapiente' => 'amet',
        'deserunt' => 'nisi',
    ];
    $request->sort = 'txid:asc';
    $request->target = 'vel';
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
