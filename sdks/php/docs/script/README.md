# script

### Available Operations

* [~~runScript~~](#runscript) - Execute a Numscript :warning: **Deprecated**

## ~~runScript~~

This route is deprecated, and has been merged into `POST /{ledger}/transactions`.


> :warning: **DEPRECATED**: this method will be removed in a future release, please migrate away from it as soon as possible.

### Example Usage

```php
<?php

declare(strict_types=1);
require_once 'vendor/autoload.php';

use \formance\stack\SDK;
use \formance\stack\Models\Shared\Security;
use \formance\stack\Models\Operations\RunScriptRequest;
use \formance\stack\Models\Shared\Script;

$sdk = SDK::builder()
    ->build();

try {
    $request = new RunScriptRequest();
    $request->script = new Script();
    $request->script->metadata = [
        'sint' => 'officia',
        'dolor' => 'debitis',
        'a' => 'dolorum',
        'in' => 'in',
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
    $request->script->vars = [
        'maiores' => 'rerum',
        'dicta' => 'magnam',
        'cumque' => 'facere',
        'ea' => 'aliquid',
    ];
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
