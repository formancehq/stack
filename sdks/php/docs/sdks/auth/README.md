# auth

### Available Operations

* [addScopeToClient](#addscopetoclient) - Add scope to client
* [addTransientScope](#addtransientscope) - Add a transient scope to a scope
* [createClient](#createclient) - Create client
* [createScope](#createscope) - Create scope
* [createSecret](#createsecret) - Add a secret to a client
* [deleteClient](#deleteclient) - Delete client
* [deleteScope](#deletescope) - Delete scope
* [deleteScopeFromClient](#deletescopefromclient) - Delete scope from client
* [deleteSecret](#deletesecret) - Delete a secret from a client
* [deleteTransientScope](#deletetransientscope) - Delete a transient scope from a scope
* [getServerInfo](#getserverinfo) - Get server info
* [listClients](#listclients) - List clients
* [listScopes](#listscopes) - List scopes
* [listUsers](#listusers) - List users
* [readClient](#readclient) - Read client
* [readScope](#readscope) - Read scope
* [readUser](#readuser) - Read user
* [updateClient](#updateclient) - Update client
* [updateScope](#updatescope) - Update scope

## addScopeToClient

Add scope to client

### Example Usage

```php
<?php

declare(strict_types=1);
require_once 'vendor/autoload.php';

use \formance\stack\SDK;
use \formance\stack\Models\Shared\Security;
use \formance\stack\Models\Operations\AddScopeToClientRequest;

$sdk = SDK::builder()
    ->build();

try {
    $request = new AddScopeToClientRequest();
    $request->clientId = 'vel';
    $request->scopeId = 'error';

    $response = $sdk->auth->addScopeToClient($request);

    if ($response->statusCode === 200) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

### Parameters

| Parameter                                                                                                       | Type                                                                                                            | Required                                                                                                        | Description                                                                                                     |
| --------------------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------- |
| `$request`                                                                                                      | [\formance\stack\Models\Operations\AddScopeToClientRequest](../../models/operations/AddScopeToClientRequest.md) | :heavy_check_mark:                                                                                              | The request object to use for the request.                                                                      |


### Response

**[?\formance\stack\Models\Operations\AddScopeToClientResponse](../../models/operations/AddScopeToClientResponse.md)**


## addTransientScope

Add a transient scope to a scope

### Example Usage

```php
<?php

declare(strict_types=1);
require_once 'vendor/autoload.php';

use \formance\stack\SDK;
use \formance\stack\Models\Shared\Security;
use \formance\stack\Models\Operations\AddTransientScopeRequest;

$sdk = SDK::builder()
    ->build();

try {
    $request = new AddTransientScopeRequest();
    $request->scopeId = 'deserunt';
    $request->transientScopeId = 'suscipit';

    $response = $sdk->auth->addTransientScope($request);

    if ($response->statusCode === 200) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

### Parameters

| Parameter                                                                                                         | Type                                                                                                              | Required                                                                                                          | Description                                                                                                       |
| ----------------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------- |
| `$request`                                                                                                        | [\formance\stack\Models\Operations\AddTransientScopeRequest](../../models/operations/AddTransientScopeRequest.md) | :heavy_check_mark:                                                                                                | The request object to use for the request.                                                                        |


### Response

**[?\formance\stack\Models\Operations\AddTransientScopeResponse](../../models/operations/AddTransientScopeResponse.md)**


## createClient

Create client

### Example Usage

```php
<?php

declare(strict_types=1);
require_once 'vendor/autoload.php';

use \formance\stack\SDK;
use \formance\stack\Models\Shared\Security;
use \formance\stack\Models\Shared\CreateClientRequest;

$sdk = SDK::builder()
    ->build();

try {
    $request = new CreateClientRequest();
    $request->description = 'iure';
    $request->metadata = [
        'debitis' => 'ipsa',
        'delectus' => 'tempora',
    ];
    $request->name = 'Minnie Schiller';
    $request->postLogoutRedirectUris = [
        'excepturi',
        'nisi',
    ];
    $request->public = false;
    $request->redirectUris = [
        'temporibus',
        'ab',
        'quis',
        'veritatis',
    ];
    $request->trusted = false;

    $response = $sdk->auth->createClient($request);

    if ($response->createClientResponse !== null) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

### Parameters

| Parameter                                                                                       | Type                                                                                            | Required                                                                                        | Description                                                                                     |
| ----------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------- |
| `$request`                                                                                      | [\formance\stack\Models\Shared\CreateClientRequest](../../models/shared/CreateClientRequest.md) | :heavy_check_mark:                                                                              | The request object to use for the request.                                                      |


### Response

**[?\formance\stack\Models\Operations\CreateClientResponse](../../models/operations/CreateClientResponse.md)**


## createScope

Create scope

### Example Usage

```php
<?php

declare(strict_types=1);
require_once 'vendor/autoload.php';

use \formance\stack\SDK;
use \formance\stack\Models\Shared\Security;
use \formance\stack\Models\Shared\CreateScopeRequest;

$sdk = SDK::builder()
    ->build();

try {
    $request = new CreateScopeRequest();
    $request->label = 'deserunt';
    $request->metadata = [
        'ipsam' => 'repellendus',
    ];

    $response = $sdk->auth->createScope($request);

    if ($response->createScopeResponse !== null) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

### Parameters

| Parameter                                                                                     | Type                                                                                          | Required                                                                                      | Description                                                                                   |
| --------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------- |
| `$request`                                                                                    | [\formance\stack\Models\Shared\CreateScopeRequest](../../models/shared/CreateScopeRequest.md) | :heavy_check_mark:                                                                            | The request object to use for the request.                                                    |


### Response

**[?\formance\stack\Models\Operations\CreateScopeResponse](../../models/operations/CreateScopeResponse.md)**


## createSecret

Add a secret to a client

### Example Usage

```php
<?php

declare(strict_types=1);
require_once 'vendor/autoload.php';

use \formance\stack\SDK;
use \formance\stack\Models\Shared\Security;
use \formance\stack\Models\Operations\CreateSecretRequest;
use \formance\stack\Models\Shared\CreateSecretRequest;

$sdk = SDK::builder()
    ->build();

try {
    $request = new CreateSecretRequest();
    $request->createSecretRequest = new CreateSecretRequest();
    $request->createSecretRequest->metadata = [
        'quo' => 'odit',
        'at' => 'at',
        'maiores' => 'molestiae',
        'quod' => 'quod',
    ];
    $request->createSecretRequest->name = 'Deanna Sauer MD';
    $request->clientId = 'officia';

    $response = $sdk->auth->createSecret($request);

    if ($response->createSecretResponse !== null) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

### Parameters

| Parameter                                                                                               | Type                                                                                                    | Required                                                                                                | Description                                                                                             |
| ------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------- |
| `$request`                                                                                              | [\formance\stack\Models\Operations\CreateSecretRequest](../../models/operations/CreateSecretRequest.md) | :heavy_check_mark:                                                                                      | The request object to use for the request.                                                              |


### Response

**[?\formance\stack\Models\Operations\CreateSecretResponse](../../models/operations/CreateSecretResponse.md)**


## deleteClient

Delete client

### Example Usage

```php
<?php

declare(strict_types=1);
require_once 'vendor/autoload.php';

use \formance\stack\SDK;
use \formance\stack\Models\Shared\Security;
use \formance\stack\Models\Operations\DeleteClientRequest;

$sdk = SDK::builder()
    ->build();

try {
    $request = new DeleteClientRequest();
    $request->clientId = 'occaecati';

    $response = $sdk->auth->deleteClient($request);

    if ($response->statusCode === 200) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

### Parameters

| Parameter                                                                                               | Type                                                                                                    | Required                                                                                                | Description                                                                                             |
| ------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------- |
| `$request`                                                                                              | [\formance\stack\Models\Operations\DeleteClientRequest](../../models/operations/DeleteClientRequest.md) | :heavy_check_mark:                                                                                      | The request object to use for the request.                                                              |


### Response

**[?\formance\stack\Models\Operations\DeleteClientResponse](../../models/operations/DeleteClientResponse.md)**


## deleteScope

Delete scope

### Example Usage

```php
<?php

declare(strict_types=1);
require_once 'vendor/autoload.php';

use \formance\stack\SDK;
use \formance\stack\Models\Shared\Security;
use \formance\stack\Models\Operations\DeleteScopeRequest;

$sdk = SDK::builder()
    ->build();

try {
    $request = new DeleteScopeRequest();
    $request->scopeId = 'fugit';

    $response = $sdk->auth->deleteScope($request);

    if ($response->statusCode === 200) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

### Parameters

| Parameter                                                                                             | Type                                                                                                  | Required                                                                                              | Description                                                                                           |
| ----------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------- |
| `$request`                                                                                            | [\formance\stack\Models\Operations\DeleteScopeRequest](../../models/operations/DeleteScopeRequest.md) | :heavy_check_mark:                                                                                    | The request object to use for the request.                                                            |


### Response

**[?\formance\stack\Models\Operations\DeleteScopeResponse](../../models/operations/DeleteScopeResponse.md)**


## deleteScopeFromClient

Delete scope from client

### Example Usage

```php
<?php

declare(strict_types=1);
require_once 'vendor/autoload.php';

use \formance\stack\SDK;
use \formance\stack\Models\Shared\Security;
use \formance\stack\Models\Operations\DeleteScopeFromClientRequest;

$sdk = SDK::builder()
    ->build();

try {
    $request = new DeleteScopeFromClientRequest();
    $request->clientId = 'deleniti';
    $request->scopeId = 'hic';

    $response = $sdk->auth->deleteScopeFromClient($request);

    if ($response->statusCode === 200) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

### Parameters

| Parameter                                                                                                                 | Type                                                                                                                      | Required                                                                                                                  | Description                                                                                                               |
| ------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------- |
| `$request`                                                                                                                | [\formance\stack\Models\Operations\DeleteScopeFromClientRequest](../../models/operations/DeleteScopeFromClientRequest.md) | :heavy_check_mark:                                                                                                        | The request object to use for the request.                                                                                |


### Response

**[?\formance\stack\Models\Operations\DeleteScopeFromClientResponse](../../models/operations/DeleteScopeFromClientResponse.md)**


## deleteSecret

Delete a secret from a client

### Example Usage

```php
<?php

declare(strict_types=1);
require_once 'vendor/autoload.php';

use \formance\stack\SDK;
use \formance\stack\Models\Shared\Security;
use \formance\stack\Models\Operations\DeleteSecretRequest;

$sdk = SDK::builder()
    ->build();

try {
    $request = new DeleteSecretRequest();
    $request->clientId = 'optio';
    $request->secretId = 'totam';

    $response = $sdk->auth->deleteSecret($request);

    if ($response->statusCode === 200) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

### Parameters

| Parameter                                                                                               | Type                                                                                                    | Required                                                                                                | Description                                                                                             |
| ------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------- |
| `$request`                                                                                              | [\formance\stack\Models\Operations\DeleteSecretRequest](../../models/operations/DeleteSecretRequest.md) | :heavy_check_mark:                                                                                      | The request object to use for the request.                                                              |


### Response

**[?\formance\stack\Models\Operations\DeleteSecretResponse](../../models/operations/DeleteSecretResponse.md)**


## deleteTransientScope

Delete a transient scope from a scope

### Example Usage

```php
<?php

declare(strict_types=1);
require_once 'vendor/autoload.php';

use \formance\stack\SDK;
use \formance\stack\Models\Shared\Security;
use \formance\stack\Models\Operations\DeleteTransientScopeRequest;

$sdk = SDK::builder()
    ->build();

try {
    $request = new DeleteTransientScopeRequest();
    $request->scopeId = 'beatae';
    $request->transientScopeId = 'commodi';

    $response = $sdk->auth->deleteTransientScope($request);

    if ($response->statusCode === 200) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

### Parameters

| Parameter                                                                                                               | Type                                                                                                                    | Required                                                                                                                | Description                                                                                                             |
| ----------------------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------- |
| `$request`                                                                                                              | [\formance\stack\Models\Operations\DeleteTransientScopeRequest](../../models/operations/DeleteTransientScopeRequest.md) | :heavy_check_mark:                                                                                                      | The request object to use for the request.                                                                              |


### Response

**[?\formance\stack\Models\Operations\DeleteTransientScopeResponse](../../models/operations/DeleteTransientScopeResponse.md)**


## getServerInfo

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
    $response = $sdk->auth->getServerInfo();

    if ($response->serverInfo !== null) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```


### Response

**[?\formance\stack\Models\Operations\GetServerInfoResponse](../../models/operations/GetServerInfoResponse.md)**


## listClients

List clients

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
    $response = $sdk->auth->listClients();

    if ($response->listClientsResponse !== null) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```


### Response

**[?\formance\stack\Models\Operations\ListClientsResponse](../../models/operations/ListClientsResponse.md)**


## listScopes

List Scopes

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
    $response = $sdk->auth->listScopes();

    if ($response->listScopesResponse !== null) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```


### Response

**[?\formance\stack\Models\Operations\ListScopesResponse](../../models/operations/ListScopesResponse.md)**


## listUsers

List users

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
    $response = $sdk->auth->listUsers();

    if ($response->listUsersResponse !== null) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```


### Response

**[?\formance\stack\Models\Operations\ListUsersResponse](../../models/operations/ListUsersResponse.md)**


## readClient

Read client

### Example Usage

```php
<?php

declare(strict_types=1);
require_once 'vendor/autoload.php';

use \formance\stack\SDK;
use \formance\stack\Models\Shared\Security;
use \formance\stack\Models\Operations\ReadClientRequest;

$sdk = SDK::builder()
    ->build();

try {
    $request = new ReadClientRequest();
    $request->clientId = 'molestiae';

    $response = $sdk->auth->readClient($request);

    if ($response->readClientResponse !== null) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

### Parameters

| Parameter                                                                                           | Type                                                                                                | Required                                                                                            | Description                                                                                         |
| --------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------- |
| `$request`                                                                                          | [\formance\stack\Models\Operations\ReadClientRequest](../../models/operations/ReadClientRequest.md) | :heavy_check_mark:                                                                                  | The request object to use for the request.                                                          |


### Response

**[?\formance\stack\Models\Operations\ReadClientResponse](../../models/operations/ReadClientResponse.md)**


## readScope

Read scope

### Example Usage

```php
<?php

declare(strict_types=1);
require_once 'vendor/autoload.php';

use \formance\stack\SDK;
use \formance\stack\Models\Shared\Security;
use \formance\stack\Models\Operations\ReadScopeRequest;

$sdk = SDK::builder()
    ->build();

try {
    $request = new ReadScopeRequest();
    $request->scopeId = 'modi';

    $response = $sdk->auth->readScope($request);

    if ($response->readScopeResponse !== null) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

### Parameters

| Parameter                                                                                         | Type                                                                                              | Required                                                                                          | Description                                                                                       |
| ------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------- |
| `$request`                                                                                        | [\formance\stack\Models\Operations\ReadScopeRequest](../../models/operations/ReadScopeRequest.md) | :heavy_check_mark:                                                                                | The request object to use for the request.                                                        |


### Response

**[?\formance\stack\Models\Operations\ReadScopeResponse](../../models/operations/ReadScopeResponse.md)**


## readUser

Read user

### Example Usage

```php
<?php

declare(strict_types=1);
require_once 'vendor/autoload.php';

use \formance\stack\SDK;
use \formance\stack\Models\Shared\Security;
use \formance\stack\Models\Operations\ReadUserRequest;

$sdk = SDK::builder()
    ->build();

try {
    $request = new ReadUserRequest();
    $request->userId = 'qui';

    $response = $sdk->auth->readUser($request);

    if ($response->readUserResponse !== null) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

### Parameters

| Parameter                                                                                       | Type                                                                                            | Required                                                                                        | Description                                                                                     |
| ----------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------- |
| `$request`                                                                                      | [\formance\stack\Models\Operations\ReadUserRequest](../../models/operations/ReadUserRequest.md) | :heavy_check_mark:                                                                              | The request object to use for the request.                                                      |


### Response

**[?\formance\stack\Models\Operations\ReadUserResponse](../../models/operations/ReadUserResponse.md)**


## updateClient

Update client

### Example Usage

```php
<?php

declare(strict_types=1);
require_once 'vendor/autoload.php';

use \formance\stack\SDK;
use \formance\stack\Models\Shared\Security;
use \formance\stack\Models\Operations\UpdateClientRequest;
use \formance\stack\Models\Shared\UpdateClientRequest;

$sdk = SDK::builder()
    ->build();

try {
    $request = new UpdateClientRequest();
    $request->updateClientRequest = new UpdateClientRequest();
    $request->updateClientRequest->description = 'impedit';
    $request->updateClientRequest->metadata = [
        'esse' => 'ipsum',
        'excepturi' => 'aspernatur',
        'perferendis' => 'ad',
    ];
    $request->updateClientRequest->name = 'Louis Moore';
    $request->updateClientRequest->postLogoutRedirectUris = [
        'hic',
        'saepe',
    ];
    $request->updateClientRequest->public = false;
    $request->updateClientRequest->redirectUris = [
        'in',
        'corporis',
        'iste',
    ];
    $request->updateClientRequest->trusted = false;
    $request->clientId = 'iure';

    $response = $sdk->auth->updateClient($request);

    if ($response->updateClientResponse !== null) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

### Parameters

| Parameter                                                                                               | Type                                                                                                    | Required                                                                                                | Description                                                                                             |
| ------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------- |
| `$request`                                                                                              | [\formance\stack\Models\Operations\UpdateClientRequest](../../models/operations/UpdateClientRequest.md) | :heavy_check_mark:                                                                                      | The request object to use for the request.                                                              |


### Response

**[?\formance\stack\Models\Operations\UpdateClientResponse](../../models/operations/UpdateClientResponse.md)**


## updateScope

Update scope

### Example Usage

```php
<?php

declare(strict_types=1);
require_once 'vendor/autoload.php';

use \formance\stack\SDK;
use \formance\stack\Models\Shared\Security;
use \formance\stack\Models\Operations\UpdateScopeRequest;
use \formance\stack\Models\Shared\UpdateScopeRequest;

$sdk = SDK::builder()
    ->build();

try {
    $request = new UpdateScopeRequest();
    $request->updateScopeRequest = new UpdateScopeRequest();
    $request->updateScopeRequest->label = 'saepe';
    $request->updateScopeRequest->metadata = [
        'architecto' => 'ipsa',
        'reiciendis' => 'est',
        'mollitia' => 'laborum',
    ];
    $request->scopeId = 'dolores';

    $response = $sdk->auth->updateScope($request);

    if ($response->updateScopeResponse !== null) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

### Parameters

| Parameter                                                                                             | Type                                                                                                  | Required                                                                                              | Description                                                                                           |
| ----------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------- |
| `$request`                                                                                            | [\formance\stack\Models\Operations\UpdateScopeRequest](../../models/operations/UpdateScopeRequest.md) | :heavy_check_mark:                                                                                    | The request object to use for the request.                                                            |


### Response

**[?\formance\stack\Models\Operations\UpdateScopeResponse](../../models/operations/UpdateScopeResponse.md)**

