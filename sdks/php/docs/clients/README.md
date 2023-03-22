# clients

### Available Operations

* [addScopeToClient](#addscopetoclient) - Add scope to client
* [createClient](#createclient) - Create client
* [createSecret](#createsecret) - Add a secret to a client
* [deleteClient](#deleteclient) - Delete client
* [deleteScopeFromClient](#deletescopefromclient) - Delete scope from client
* [deleteSecret](#deletesecret) - Delete a secret from a client
* [listClients](#listclients) - List clients
* [readClient](#readclient) - Read client
* [updateClient](#updateclient) - Update client

## addScopeToClient

Add scope to client

### Example Usage

```php
<?php

declare(strict_types=1);
require_once 'vendor/autoload.php';

use \OpenAPI\OpenAPI\SDK;
use \OpenAPI\OpenAPI\Models\Shared\Security;
use \OpenAPI\OpenAPI\Models\Operations\AddScopeToClientRequest;

$sdk = SDK::builder()
    ->build();

try {
    $request = new AddScopeToClientRequest();
    $request->clientId = 'velit';
    $request->scopeId = 'error';

    $response = $sdk->clients->addScopeToClient($request);

    if ($response->statusCode === 200) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

## createClient

Create client

### Example Usage

```php
<?php

declare(strict_types=1);
require_once 'vendor/autoload.php';

use \OpenAPI\OpenAPI\SDK;
use \OpenAPI\OpenAPI\Models\Shared\Security;
use \OpenAPI\OpenAPI\Models\Shared\CreateClientRequest;

$sdk = SDK::builder()
    ->build();

try {
    $request = new CreateClientRequest();
    $request->description = 'quia';
    $request->metadata = [
        'vitae' => 'laborum',
        'animi' => 'enim',
    ];
    $request->name = 'Angelica Dietrich';
    $request->postLogoutRedirectUris = [
        'possimus',
        'aut',
        'quasi',
    ];
    $request->public = false;
    $request->redirectUris = [
        'temporibus',
        'laborum',
        'quasi',
    ];
    $request->trusted = false;

    $response = $sdk->clients->createClient($request);

    if ($response->createClientResponse !== null) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

## createSecret

Add a secret to a client

### Example Usage

```php
<?php

declare(strict_types=1);
require_once 'vendor/autoload.php';

use \OpenAPI\OpenAPI\SDK;
use \OpenAPI\OpenAPI\Models\Shared\Security;
use \OpenAPI\OpenAPI\Models\Operations\CreateSecretRequest;
use \OpenAPI\OpenAPI\Models\Shared\CreateSecretRequest;

$sdk = SDK::builder()
    ->build();

try {
    $request = new CreateSecretRequest();
    $request->createSecretRequest = new CreateSecretRequest();
    $request->createSecretRequest->metadata = [
        'voluptatibus' => 'vero',
        'nihil' => 'praesentium',
        'voluptatibus' => 'ipsa',
        'omnis' => 'voluptate',
    ];
    $request->createSecretRequest->name = 'Thomas Batz';
    $request->clientId = 'maiores';

    $response = $sdk->clients->createSecret($request);

    if ($response->createSecretResponse !== null) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

## deleteClient

Delete client

### Example Usage

```php
<?php

declare(strict_types=1);
require_once 'vendor/autoload.php';

use \OpenAPI\OpenAPI\SDK;
use \OpenAPI\OpenAPI\Models\Shared\Security;
use \OpenAPI\OpenAPI\Models\Operations\DeleteClientRequest;

$sdk = SDK::builder()
    ->build();

try {
    $request = new DeleteClientRequest();
    $request->clientId = 'dicta';

    $response = $sdk->clients->deleteClient($request);

    if ($response->statusCode === 200) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

## deleteScopeFromClient

Delete scope from client

### Example Usage

```php
<?php

declare(strict_types=1);
require_once 'vendor/autoload.php';

use \OpenAPI\OpenAPI\SDK;
use \OpenAPI\OpenAPI\Models\Shared\Security;
use \OpenAPI\OpenAPI\Models\Operations\DeleteScopeFromClientRequest;

$sdk = SDK::builder()
    ->build();

try {
    $request = new DeleteScopeFromClientRequest();
    $request->clientId = 'corporis';
    $request->scopeId = 'dolore';

    $response = $sdk->clients->deleteScopeFromClient($request);

    if ($response->statusCode === 200) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

## deleteSecret

Delete a secret from a client

### Example Usage

```php
<?php

declare(strict_types=1);
require_once 'vendor/autoload.php';

use \OpenAPI\OpenAPI\SDK;
use \OpenAPI\OpenAPI\Models\Shared\Security;
use \OpenAPI\OpenAPI\Models\Operations\DeleteSecretRequest;

$sdk = SDK::builder()
    ->build();

try {
    $request = new DeleteSecretRequest();
    $request->clientId = 'iusto';
    $request->secretId = 'dicta';

    $response = $sdk->clients->deleteSecret($request);

    if ($response->statusCode === 200) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

## listClients

List clients

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
    $response = $sdk->clients->listClients();

    if ($response->listClientsResponse !== null) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

## readClient

Read client

### Example Usage

```php
<?php

declare(strict_types=1);
require_once 'vendor/autoload.php';

use \OpenAPI\OpenAPI\SDK;
use \OpenAPI\OpenAPI\Models\Shared\Security;
use \OpenAPI\OpenAPI\Models\Operations\ReadClientRequest;

$sdk = SDK::builder()
    ->build();

try {
    $request = new ReadClientRequest();
    $request->clientId = 'harum';

    $response = $sdk->clients->readClient($request);

    if ($response->readClientResponse !== null) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

## updateClient

Update client

### Example Usage

```php
<?php

declare(strict_types=1);
require_once 'vendor/autoload.php';

use \OpenAPI\OpenAPI\SDK;
use \OpenAPI\OpenAPI\Models\Shared\Security;
use \OpenAPI\OpenAPI\Models\Operations\UpdateClientRequest;
use \OpenAPI\OpenAPI\Models\Shared\UpdateClientRequest;

$sdk = SDK::builder()
    ->build();

try {
    $request = new UpdateClientRequest();
    $request->updateClientRequest = new UpdateClientRequest();
    $request->updateClientRequest->description = 'enim';
    $request->updateClientRequest->metadata = [
        'commodi' => 'repudiandae',
        'quae' => 'ipsum',
        'quidem' => 'molestias',
        'excepturi' => 'pariatur',
    ];
    $request->updateClientRequest->name = 'Irma Ledner DVM';
    $request->updateClientRequest->postLogoutRedirectUris = [
        'veritatis',
        'itaque',
        'incidunt',
    ];
    $request->updateClientRequest->public = false;
    $request->updateClientRequest->redirectUris = [
        'consequatur',
        'est',
    ];
    $request->updateClientRequest->trusted = false;
    $request->clientId = 'quibusdam';

    $response = $sdk->clients->updateClient($request);

    if ($response->updateClientResponse !== null) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```
