# scopes

### Available Operations

* [addTransientScope](#addtransientscope) - Add a transient scope to a scope
* [createScope](#createscope) - Create scope
* [deleteScope](#deletescope) - Delete scope
* [deleteTransientScope](#deletetransientscope) - Delete a transient scope from a scope
* [listScopes](#listscopes) - List scopes
* [readScope](#readscope) - Read scope
* [updateScope](#updatescope) - Update scope

## addTransientScope

Add a transient scope to a scope

### Example Usage

```php
<?php

declare(strict_types=1);
require_once 'vendor/autoload.php';

use \OpenAPI\OpenAPI\SDK;
use \OpenAPI\OpenAPI\Models\Shared\Security;
use \OpenAPI\OpenAPI\Models\Operations\AddTransientScopeRequest;

$sdk = SDK::builder()
    ->build();

try {
    $request = new AddTransientScopeRequest();
    $request->scopeId = 'aspernatur';
    $request->transientScopeId = 'architecto';

    $response = $sdk->scopes->addTransientScope($request);

    if ($response->statusCode === 200) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

## createScope

Create scope

### Example Usage

```php
<?php

declare(strict_types=1);
require_once 'vendor/autoload.php';

use \OpenAPI\OpenAPI\SDK;
use \OpenAPI\OpenAPI\Models\Shared\Security;
use \OpenAPI\OpenAPI\Models\Shared\CreateScopeRequest;

$sdk = SDK::builder()
    ->build();

try {
    $request = new CreateScopeRequest();
    $request->label = 'magnam';
    $request->metadata = [
        'excepturi' => 'ullam',
    ];

    $response = $sdk->scopes->createScope($request);

    if ($response->createScopeResponse !== null) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

## deleteScope

Delete scope

### Example Usage

```php
<?php

declare(strict_types=1);
require_once 'vendor/autoload.php';

use \OpenAPI\OpenAPI\SDK;
use \OpenAPI\OpenAPI\Models\Shared\Security;
use \OpenAPI\OpenAPI\Models\Operations\DeleteScopeRequest;

$sdk = SDK::builder()
    ->build();

try {
    $request = new DeleteScopeRequest();
    $request->scopeId = 'provident';

    $response = $sdk->scopes->deleteScope($request);

    if ($response->statusCode === 200) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

## deleteTransientScope

Delete a transient scope from a scope

### Example Usage

```php
<?php

declare(strict_types=1);
require_once 'vendor/autoload.php';

use \OpenAPI\OpenAPI\SDK;
use \OpenAPI\OpenAPI\Models\Shared\Security;
use \OpenAPI\OpenAPI\Models\Operations\DeleteTransientScopeRequest;

$sdk = SDK::builder()
    ->build();

try {
    $request = new DeleteTransientScopeRequest();
    $request->scopeId = 'quos';
    $request->transientScopeId = 'sint';

    $response = $sdk->scopes->deleteTransientScope($request);

    if ($response->statusCode === 200) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

## listScopes

List Scopes

### Example Usage

```php
<?php

declare(strict_types=1);
require_once 'vendor/autoload.php';

use \OpenAPI\OpenAPI\SDK;
use \OpenAPI\OpenAPI\Models\Shared\Security;

$sdk = SDK::builder()
    ->build();

try {
    $response = $sdk->scopes->listScopes();

    if ($response->listScopesResponse !== null) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

## readScope

Read scope

### Example Usage

```php
<?php

declare(strict_types=1);
require_once 'vendor/autoload.php';

use \OpenAPI\OpenAPI\SDK;
use \OpenAPI\OpenAPI\Models\Shared\Security;
use \OpenAPI\OpenAPI\Models\Operations\ReadScopeRequest;

$sdk = SDK::builder()
    ->build();

try {
    $request = new ReadScopeRequest();
    $request->scopeId = 'accusantium';

    $response = $sdk->scopes->readScope($request);

    if ($response->readScopeResponse !== null) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

## updateScope

Update scope

### Example Usage

```php
<?php

declare(strict_types=1);
require_once 'vendor/autoload.php';

use \OpenAPI\OpenAPI\SDK;
use \OpenAPI\OpenAPI\Models\Shared\Security;
use \OpenAPI\OpenAPI\Models\Operations\UpdateScopeRequest;
use \OpenAPI\OpenAPI\Models\Shared\UpdateScopeRequest;

$sdk = SDK::builder()
    ->build();

try {
    $request = new UpdateScopeRequest();
    $request->updateScopeRequest = new UpdateScopeRequest();
    $request->updateScopeRequest->label = 'mollitia';
    $request->updateScopeRequest->metadata = [
        'mollitia' => 'ad',
        'eum' => 'dolor',
        'necessitatibus' => 'odit',
        'nemo' => 'quasi',
    ];
    $request->scopeId = 'iure';

    $response = $sdk->scopes->updateScope($request);

    if ($response->updateScopeResponse !== null) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```
