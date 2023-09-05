# script

### Available Operations

* [~~runScript~~](#runscript) - Execute a Numscript :warning: **Deprecated**

## ~~runScript~~

This route is deprecated, and has been merged into `POST /{ledger}/transactions`.


> :warning: **DEPRECATED**: This will be removed in a future release, please migrate away from it as soon as possible.

### Example Usage

```php
<?php

declare(strict_types=1);
require_once 'vendor/autoload.php';

use \formance\stack\SDK;
use \formance\stack\Models\Shared\Security;
use \formance\stack\Models\Operations\RunScriptRequest;
use \formance\stack\Models\Shared\Script;
use \formance\stack\Models\Shared\ScriptVars;

$sdk = SDK::builder()
    ->build();

try {
    $request = new RunScriptRequest();
    $request->script = new Script();
    $request->script->metadata = [
        'rem' => 'voluptates',
        'quasi' => 'repudiandae',
        'sint' => 'veritatis',
    ];
    $request->script->plain = 'vars {
    account $user
    }
    send [COIN 10] (
    	source = @world
    	destination = $user
    )
    ';
    $request->script->reference = 'order_1234';
    $request->script->vars = new ScriptVars();
    $request->ledger = 'ledger001';
    $request->preview = true;

    $response = $sdk->script->runScript($request);

    if ($response->scriptResponse !== null) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

### Parameters

| Parameter                                                                                         | Type                                                                                              | Required                                                                                          | Description                                                                                       |
| ------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------- |
| `$request`                                                                                        | [\formance\stack\Models\Operations\RunScriptRequest](../../models/operations/RunScriptRequest.md) | :heavy_check_mark:                                                                                | The request object to use for the request.                                                        |


### Response

**[?\formance\stack\Models\Operations\RunScriptResponse](../../models/operations/RunScriptResponse.md)**

