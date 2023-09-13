# mapping

### Available Operations

* [getMapping](#getmapping) - Get the mapping of a ledger
* [updateMapping](#updatemapping) - Update the mapping of a ledger

## getMapping

Get the mapping of a ledger

### Example Usage

```php
<?php

declare(strict_types=1);
require_once 'vendor/autoload.php';

use \formance\stack\SDK;
use \formance\stack\Models\Shared\Security;
use \formance\stack\Models\Operations\GetMappingRequest;

$sdk = SDK::builder()
    ->build();

try {
    $request = new GetMappingRequest();
    $request->ledger = 'ledger001';

    $response = $sdk->mapping->getMapping($request);

    if ($response->mappingResponse !== null) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

### Parameters

| Parameter                                                                                           | Type                                                                                                | Required                                                                                            | Description                                                                                         |
| --------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------- |
| `$request`                                                                                          | [\formance\stack\Models\Operations\GetMappingRequest](../../models/operations/GetMappingRequest.md) | :heavy_check_mark:                                                                                  | The request object to use for the request.                                                          |


### Response

**[?\formance\stack\Models\Operations\GetMappingResponse](../../models/operations/GetMappingResponse.md)**


## updateMapping

Update the mapping of a ledger

### Example Usage

```php
<?php

declare(strict_types=1);
require_once 'vendor/autoload.php';

use \formance\stack\SDK;
use \formance\stack\Models\Shared\Security;
use \formance\stack\Models\Operations\UpdateMappingRequest;
use \formance\stack\Models\Shared\Mapping;
use \formance\stack\Models\Shared\Contract;
use \formance\stack\Models\Shared\ContractExpr;

$sdk = SDK::builder()
    ->build();

try {
    $request = new UpdateMappingRequest();
    $request->mapping = new Mapping();
    $request->mapping->contracts = [
        new Contract(),
        new Contract(),
    ];
    $request->ledger = 'ledger001';

    $response = $sdk->mapping->updateMapping($request);

    if ($response->mappingResponse !== null) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

### Parameters

| Parameter                                                                                                 | Type                                                                                                      | Required                                                                                                  | Description                                                                                               |
| --------------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------- |
| `$request`                                                                                                | [\formance\stack\Models\Operations\UpdateMappingRequest](../../models/operations/UpdateMappingRequest.md) | :heavy_check_mark:                                                                                        | The request object to use for the request.                                                                |


### Response

**[?\formance\stack\Models\Operations\UpdateMappingResponse](../../models/operations/UpdateMappingResponse.md)**

