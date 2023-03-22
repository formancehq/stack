# accounts

### Available Operations

* [addMetadataToAccount](#addmetadatatoaccount) - Add metadata to an account
* [countAccounts](#countaccounts) - Count the accounts from a ledger
* [getAccount](#getaccount) - Get account by its address
* [listAccounts](#listaccounts) - List accounts from a ledger

## addMetadataToAccount

Add metadata to an account

### Example Usage

```php
<?php

declare(strict_types=1);
require_once 'vendor/autoload.php';

use \formance\stack\SDK;
use \formance\stack\Models\Shared\Security;
use \formance\stack\Models\Operations\AddMetadataToAccountRequest;

$sdk = SDK::builder()
    ->build();

try {
    $request = new AddMetadataToAccountRequest();
    $request->idempotencyKey = 'corrupti';
    $request->requestBody = [
        'distinctio' => 'quibusdam',
        'unde' => 'nulla',
        'corrupti' => 'illum',
    ];
    $request->address = 'users:001';
    $request->async = true;
    $request->dryRun = true;
    $request->ledger = 'ledger001';

    $response = $sdk->accounts->addMetadataToAccount($request);

    if ($response->statusCode === 200) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

## countAccounts

Count the accounts from a ledger

### Example Usage

```php
<?php

declare(strict_types=1);
require_once 'vendor/autoload.php';

use \formance\stack\SDK;
use \formance\stack\Models\Shared\Security;
use \formance\stack\Models\Operations\CountAccountsRequest;

$sdk = SDK::builder()
    ->build();

try {
    $request = new CountAccountsRequest();
    $request->address = 'users:.+';
    $request->ledger = 'ledger001';
    $request->metadata = [
        'error' => 'deserunt',
        'suscipit' => 'iure',
    ];

    $response = $sdk->accounts->countAccounts($request);

    if ($response->statusCode === 200) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

## getAccount

Get account by its address

### Example Usage

```php
<?php

declare(strict_types=1);
require_once 'vendor/autoload.php';

use \formance\stack\SDK;
use \formance\stack\Models\Shared\Security;
use \formance\stack\Models\Operations\GetAccountRequest;

$sdk = SDK::builder()
    ->build();

try {
    $request = new GetAccountRequest();
    $request->address = 'users:001';
    $request->ledger = 'ledger001';

    $response = $sdk->accounts->getAccount($request);

    if ($response->accountResponse !== null) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

## listAccounts

List accounts from a ledger, sorted by address in descending order.

### Example Usage

```php
<?php

declare(strict_types=1);
require_once 'vendor/autoload.php';

use \formance\stack\SDK;
use \formance\stack\Models\Shared\Security;
use \formance\stack\Models\Operations\ListAccountsRequest;
use \formance\stack\Models\Operations\ListAccountsBalanceOperator;

$sdk = SDK::builder()
    ->build();

try {
    $request = new ListAccountsRequest();
    $request->address = 'users:.+';
    $request->balance = 2400;
    $request->balanceOperator = ListAccountsBalanceOperator::GTE;
    $request->cursor = 'aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ==';
    $request->ledger = 'ledger001';
    $request->metadata = [
        'debitis' => 'ipsa',
        'delectus' => 'tempora',
    ];
    $request->pageSize = 383441;

    $response = $sdk->accounts->listAccounts($request);

    if ($response->accountsCursorResponse !== null) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```
