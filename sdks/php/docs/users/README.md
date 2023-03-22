# users

### Available Operations

* [listUsers](#listusers) - List users
* [readUser](#readuser) - Read user

## listUsers

List users

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
    $response = $sdk->users->listUsers();

    if ($response->listUsersResponse !== null) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

## readUser

Read user

### Example Usage

```php
<?php

declare(strict_types=1);
require_once 'vendor/autoload.php';

use \OpenAPI\OpenAPI\SDK;
use \OpenAPI\OpenAPI\Models\Shared\Security;
use \OpenAPI\OpenAPI\Models\Operations\ReadUserRequest;

$sdk = SDK::builder()
    ->build();

try {
    $request = new ReadUserRequest();
    $request->userId = 'dolores';

    $response = $sdk->users->readUser($request);

    if ($response->readUserResponse !== null) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```
