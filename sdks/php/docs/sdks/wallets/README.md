# wallets

### Available Operations

* [confirmHold](#confirmhold) - Confirm a hold
* [createBalance](#createbalance) - Create a balance
* [createWallet](#createwallet) - Create a new wallet
* [creditWallet](#creditwallet) - Credit a wallet
* [debitWallet](#debitwallet) - Debit a wallet
* [getBalance](#getbalance) - Get detailed balance
* [getHold](#gethold) - Get a hold
* [getHolds](#getholds) - Get all holds for a wallet
* [getTransactions](#gettransactions)
* [getWallet](#getwallet) - Get a wallet
* [getWalletSummary](#getwalletsummary) - Get wallet summary
* [listBalances](#listbalances) - List balances of a wallet
* [listWallets](#listwallets) - List all wallets
* [updateWallet](#updatewallet) - Update a wallet
* [voidHold](#voidhold) - Cancel a hold
* [walletsgetServerInfo](#walletsgetserverinfo) - Get server info

## confirmHold

Confirm a hold

### Example Usage

```php
<?php

declare(strict_types=1);
require_once 'vendor/autoload.php';

use \formance\stack\SDK;
use \formance\stack\Models\Shared\Security;
use \formance\stack\Models\Operations\ConfirmHoldRequest;
use \formance\stack\Models\Shared\ConfirmHoldRequest;

$sdk = SDK::builder()
    ->build();

try {
    $request = new ConfirmHoldRequest();
    $request->confirmHoldRequest = new ConfirmHoldRequest();
    $request->confirmHoldRequest->amount = 100;
    $request->confirmHoldRequest->final = true;
    $request->holdId = 'non';

    $response = $sdk->wallets->confirmHold($request);

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
| `$request`                                                                                            | [\formance\stack\Models\Operations\ConfirmHoldRequest](../../models/operations/ConfirmHoldRequest.md) | :heavy_check_mark:                                                                                    | The request object to use for the request.                                                            |


### Response

**[?\formance\stack\Models\Operations\ConfirmHoldResponse](../../models/operations/ConfirmHoldResponse.md)**


## createBalance

Create a balance

### Example Usage

```php
<?php

declare(strict_types=1);
require_once 'vendor/autoload.php';

use \formance\stack\SDK;
use \formance\stack\Models\Shared\Security;
use \formance\stack\Models\Operations\CreateBalanceRequest;
use \formance\stack\Models\Shared\CreateBalanceRequest;

$sdk = SDK::builder()
    ->build();

try {
    $request = new CreateBalanceRequest();
    $request->createBalanceRequest = new CreateBalanceRequest();
    $request->createBalanceRequest->expiresAt = DateTime::createFromFormat('Y-m-d\TH:i:sP', '2021-04-10T02:35:06.342Z');
    $request->createBalanceRequest->name = 'Sherri Tremblay';
    $request->createBalanceRequest->priority = 223081;
    $request->id = 'efa77dfb-14cd-466a-a395-efb9ba88f3a6';

    $response = $sdk->wallets->createBalance($request);

    if ($response->createBalanceResponse !== null) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

### Parameters

| Parameter                                                                                                 | Type                                                                                                      | Required                                                                                                  | Description                                                                                               |
| --------------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------- |
| `$request`                                                                                                | [\formance\stack\Models\Operations\CreateBalanceRequest](../../models/operations/CreateBalanceRequest.md) | :heavy_check_mark:                                                                                        | The request object to use for the request.                                                                |


### Response

**[?\formance\stack\Models\Operations\CreateBalanceResponse](../../models/operations/CreateBalanceResponse.md)**


## createWallet

Create a new wallet

### Example Usage

```php
<?php

declare(strict_types=1);
require_once 'vendor/autoload.php';

use \formance\stack\SDK;
use \formance\stack\Models\Shared\Security;
use \formance\stack\Models\Shared\CreateWalletRequest;

$sdk = SDK::builder()
    ->build();

try {
    $request = new CreateWalletRequest();
    $request->metadata = [
        'natus' => 'omnis',
        'molestiae' => 'perferendis',
    ];
    $request->name = 'Megan Rau';

    $response = $sdk->wallets->createWallet($request);

    if ($response->createWalletResponse !== null) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

### Parameters

| Parameter                                                                                       | Type                                                                                            | Required                                                                                        | Description                                                                                     |
| ----------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------- |
| `$request`                                                                                      | [\formance\stack\Models\Shared\CreateWalletRequest](../../models/shared/CreateWalletRequest.md) | :heavy_check_mark:                                                                              | The request object to use for the request.                                                      |


### Response

**[?\formance\stack\Models\Operations\CreateWalletResponse](../../models/operations/CreateWalletResponse.md)**


## creditWallet

Credit a wallet

### Example Usage

```php
<?php

declare(strict_types=1);
require_once 'vendor/autoload.php';

use \formance\stack\SDK;
use \formance\stack\Models\Shared\Security;
use \formance\stack\Models\Operations\CreditWalletRequest;
use \formance\stack\Models\Shared\CreditWalletRequest;
use \formance\stack\Models\Shared\Monetary;

$sdk = SDK::builder()
    ->build();

try {
    $request = new CreditWalletRequest();
    $request->creditWalletRequest = new CreditWalletRequest();
    $request->creditWalletRequest->amount = new Monetary();
    $request->creditWalletRequest->amount->amount = 290077;
    $request->creditWalletRequest->amount->asset = 'suscipit';
    $request->creditWalletRequest->balance = 'natus';
    $request->creditWalletRequest->metadata = [
        'eum' => 'vero',
        'aspernatur' => 'architecto',
        'magnam' => 'et',
    ];
    $request->creditWalletRequest->reference = 'excepturi';
    $request->creditWalletRequest->sources = [
        new WalletSubject(),
        new WalletSubject(),
    ];
    $request->id = '90afa563-e251-46fe-8c8b-711e5b7fd2ed';

    $response = $sdk->wallets->creditWallet($request);

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
| `$request`                                                                                              | [\formance\stack\Models\Operations\CreditWalletRequest](../../models/operations/CreditWalletRequest.md) | :heavy_check_mark:                                                                                      | The request object to use for the request.                                                              |


### Response

**[?\formance\stack\Models\Operations\CreditWalletResponse](../../models/operations/CreditWalletResponse.md)**


## debitWallet

Debit a wallet

### Example Usage

```php
<?php

declare(strict_types=1);
require_once 'vendor/autoload.php';

use \formance\stack\SDK;
use \formance\stack\Models\Shared\Security;
use \formance\stack\Models\Operations\DebitWalletRequest;
use \formance\stack\Models\Shared\DebitWalletRequest;
use \formance\stack\Models\Shared\Monetary;

$sdk = SDK::builder()
    ->build();

try {
    $request = new DebitWalletRequest();
    $request->debitWalletRequest = new DebitWalletRequest();
    $request->debitWalletRequest->amount = new Monetary();
    $request->debitWalletRequest->amount->amount = 37559;
    $request->debitWalletRequest->amount->asset = 'consequuntur';
    $request->debitWalletRequest->balances = [
        'natus',
        'magni',
        'sunt',
    ];
    $request->debitWalletRequest->description = 'quo';
    $request->debitWalletRequest->destination = new WalletSubject();
    $request->debitWalletRequest->destination->balance = 'pariatur';
    $request->debitWalletRequest->destination->identifier = 'maxime';
    $request->debitWalletRequest->destination->type = 'ea';
    $request->debitWalletRequest->metadata = [
        'odit' => 'ea',
        'accusantium' => 'ab',
        'maiores' => 'quidem',
    ];
    $request->debitWalletRequest->pending = false;
    $request->id = '576b0d5f-0d30-4c5f-bb25-87053202c73d';

    $response = $sdk->wallets->debitWallet($request);

    if ($response->debitWalletResponse !== null) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

### Parameters

| Parameter                                                                                             | Type                                                                                                  | Required                                                                                              | Description                                                                                           |
| ----------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------- |
| `$request`                                                                                            | [\formance\stack\Models\Operations\DebitWalletRequest](../../models/operations/DebitWalletRequest.md) | :heavy_check_mark:                                                                                    | The request object to use for the request.                                                            |


### Response

**[?\formance\stack\Models\Operations\DebitWalletResponse](../../models/operations/DebitWalletResponse.md)**


## getBalance

Get detailed balance

### Example Usage

```php
<?php

declare(strict_types=1);
require_once 'vendor/autoload.php';

use \formance\stack\SDK;
use \formance\stack\Models\Shared\Security;
use \formance\stack\Models\Operations\GetBalanceRequest;

$sdk = SDK::builder()
    ->build();

try {
    $request = new GetBalanceRequest();
    $request->balanceName = 'nostrum';
    $request->id = 'fe9b90c2-8909-4b3f-a49a-8d9cbf486333';

    $response = $sdk->wallets->getBalance($request);

    if ($response->getBalanceResponse !== null) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

### Parameters

| Parameter                                                                                           | Type                                                                                                | Required                                                                                            | Description                                                                                         |
| --------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------- |
| `$request`                                                                                          | [\formance\stack\Models\Operations\GetBalanceRequest](../../models/operations/GetBalanceRequest.md) | :heavy_check_mark:                                                                                  | The request object to use for the request.                                                          |


### Response

**[?\formance\stack\Models\Operations\GetBalanceResponse](../../models/operations/GetBalanceResponse.md)**


## getHold

Get a hold

### Example Usage

```php
<?php

declare(strict_types=1);
require_once 'vendor/autoload.php';

use \formance\stack\SDK;
use \formance\stack\Models\Shared\Security;
use \formance\stack\Models\Operations\GetHoldRequest;

$sdk = SDK::builder()
    ->build();

try {
    $request = new GetHoldRequest();
    $request->holdID = 'qui';

    $response = $sdk->wallets->getHold($request);

    if ($response->getHoldResponse !== null) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

### Parameters

| Parameter                                                                                     | Type                                                                                          | Required                                                                                      | Description                                                                                   |
| --------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------- |
| `$request`                                                                                    | [\formance\stack\Models\Operations\GetHoldRequest](../../models/operations/GetHoldRequest.md) | :heavy_check_mark:                                                                            | The request object to use for the request.                                                    |


### Response

**[?\formance\stack\Models\Operations\GetHoldResponse](../../models/operations/GetHoldResponse.md)**


## getHolds

Get all holds for a wallet

### Example Usage

```php
<?php

declare(strict_types=1);
require_once 'vendor/autoload.php';

use \formance\stack\SDK;
use \formance\stack\Models\Shared\Security;
use \formance\stack\Models\Operations\GetHoldsRequest;

$sdk = SDK::builder()
    ->build();

try {
    $request = new GetHoldsRequest();
    $request->cursor = 'aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ==';
    $request->metadata = [
        'hic' => 'excepturi',
    ];
    $request->pageSize = 739551;
    $request->walletID = 'voluptate';

    $response = $sdk->wallets->getHolds($request);

    if ($response->getHoldsResponse !== null) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

### Parameters

| Parameter                                                                                       | Type                                                                                            | Required                                                                                        | Description                                                                                     |
| ----------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------- |
| `$request`                                                                                      | [\formance\stack\Models\Operations\GetHoldsRequest](../../models/operations/GetHoldsRequest.md) | :heavy_check_mark:                                                                              | The request object to use for the request.                                                      |


### Response

**[?\formance\stack\Models\Operations\GetHoldsResponse](../../models/operations/GetHoldsResponse.md)**


## getTransactions

### Example Usage

```php
<?php

declare(strict_types=1);
require_once 'vendor/autoload.php';

use \formance\stack\SDK;
use \formance\stack\Models\Shared\Security;
use \formance\stack\Models\Operations\GetTransactionsRequest;

$sdk = SDK::builder()
    ->build();

try {
    $request = new GetTransactionsRequest();
    $request->cursor = 'aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ==';
    $request->pageSize = 490459;
    $request->walletID = 'reiciendis';

    $response = $sdk->wallets->getTransactions($request);

    if ($response->getTransactionsResponse !== null) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

### Parameters

| Parameter                                                                                                     | Type                                                                                                          | Required                                                                                                      | Description                                                                                                   |
| ------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------- |
| `$request`                                                                                                    | [\formance\stack\Models\Operations\GetTransactionsRequest](../../models/operations/GetTransactionsRequest.md) | :heavy_check_mark:                                                                                            | The request object to use for the request.                                                                    |


### Response

**[?\formance\stack\Models\Operations\GetTransactionsResponse](../../models/operations/GetTransactionsResponse.md)**


## getWallet

Get a wallet

### Example Usage

```php
<?php

declare(strict_types=1);
require_once 'vendor/autoload.php';

use \formance\stack\SDK;
use \formance\stack\Models\Shared\Security;
use \formance\stack\Models\Operations\GetWalletRequest;

$sdk = SDK::builder()
    ->build();

try {
    $request = new GetWalletRequest();
    $request->id = '3a410067-4ebf-4692-80d1-ba77a89ebf73';

    $response = $sdk->wallets->getWallet($request);

    if ($response->getWalletResponse !== null) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

### Parameters

| Parameter                                                                                         | Type                                                                                              | Required                                                                                          | Description                                                                                       |
| ------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------- |
| `$request`                                                                                        | [\formance\stack\Models\Operations\GetWalletRequest](../../models/operations/GetWalletRequest.md) | :heavy_check_mark:                                                                                | The request object to use for the request.                                                        |


### Response

**[?\formance\stack\Models\Operations\GetWalletResponse](../../models/operations/GetWalletResponse.md)**


## getWalletSummary

Get wallet summary

### Example Usage

```php
<?php

declare(strict_types=1);
require_once 'vendor/autoload.php';

use \formance\stack\SDK;
use \formance\stack\Models\Shared\Security;
use \formance\stack\Models\Operations\GetWalletSummaryRequest;

$sdk = SDK::builder()
    ->build();

try {
    $request = new GetWalletSummaryRequest();
    $request->id = '7ae4203c-e5e6-4a95-98a0-d446ce2af7a7';

    $response = $sdk->wallets->getWalletSummary($request);

    if ($response->getWalletSummaryResponse !== null) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

### Parameters

| Parameter                                                                                                       | Type                                                                                                            | Required                                                                                                        | Description                                                                                                     |
| --------------------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------- |
| `$request`                                                                                                      | [\formance\stack\Models\Operations\GetWalletSummaryRequest](../../models/operations/GetWalletSummaryRequest.md) | :heavy_check_mark:                                                                                              | The request object to use for the request.                                                                      |


### Response

**[?\formance\stack\Models\Operations\GetWalletSummaryResponse](../../models/operations/GetWalletSummaryResponse.md)**


## listBalances

List balances of a wallet

### Example Usage

```php
<?php

declare(strict_types=1);
require_once 'vendor/autoload.php';

use \formance\stack\SDK;
use \formance\stack\Models\Shared\Security;
use \formance\stack\Models\Operations\ListBalancesRequest;

$sdk = SDK::builder()
    ->build();

try {
    $request = new ListBalancesRequest();
    $request->id = '3cf3be45-3f87-40b3-a6b5-a73429cdb1a8';

    $response = $sdk->wallets->listBalances($request);

    if ($response->listBalancesResponse !== null) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

### Parameters

| Parameter                                                                                               | Type                                                                                                    | Required                                                                                                | Description                                                                                             |
| ------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------- |
| `$request`                                                                                              | [\formance\stack\Models\Operations\ListBalancesRequest](../../models/operations/ListBalancesRequest.md) | :heavy_check_mark:                                                                                      | The request object to use for the request.                                                              |


### Response

**[?\formance\stack\Models\Operations\ListBalancesResponse](../../models/operations/ListBalancesResponse.md)**


## listWallets

List all wallets

### Example Usage

```php
<?php

declare(strict_types=1);
require_once 'vendor/autoload.php';

use \formance\stack\SDK;
use \formance\stack\Models\Shared\Security;
use \formance\stack\Models\Operations\ListWalletsRequest;

$sdk = SDK::builder()
    ->build();

try {
    $request = new ListWalletsRequest();
    $request->cursor = 'aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ==';
    $request->metadata = [
        'aspernatur' => 'dolores',
        'distinctio' => 'facilis',
    ];
    $request->name = 'Constance Mann';
    $request->pageSize = 204865;

    $response = $sdk->wallets->listWallets($request);

    if ($response->listWalletsResponse !== null) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

### Parameters

| Parameter                                                                                             | Type                                                                                                  | Required                                                                                              | Description                                                                                           |
| ----------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------- |
| `$request`                                                                                            | [\formance\stack\Models\Operations\ListWalletsRequest](../../models/operations/ListWalletsRequest.md) | :heavy_check_mark:                                                                                    | The request object to use for the request.                                                            |


### Response

**[?\formance\stack\Models\Operations\ListWalletsResponse](../../models/operations/ListWalletsResponse.md)**


## updateWallet

Update a wallet

### Example Usage

```php
<?php

declare(strict_types=1);
require_once 'vendor/autoload.php';

use \formance\stack\SDK;
use \formance\stack\Models\Shared\Security;
use \formance\stack\Models\Operations\UpdateWalletRequest;
use \formance\stack\Models\Operations\UpdateWalletRequestBody;

$sdk = SDK::builder()
    ->build();

try {
    $request = new UpdateWalletRequest();
    $request->requestBody = new UpdateWalletRequestBody();
    $request->requestBody->metadata = [
        'magni' => 'odio',
    ];
    $request->id = '15bf0cbb-1e31-4b8b-90f3-443a1108e0ad';

    $response = $sdk->wallets->updateWallet($request);

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
| `$request`                                                                                              | [\formance\stack\Models\Operations\UpdateWalletRequest](../../models/operations/UpdateWalletRequest.md) | :heavy_check_mark:                                                                                      | The request object to use for the request.                                                              |


### Response

**[?\formance\stack\Models\Operations\UpdateWalletResponse](../../models/operations/UpdateWalletResponse.md)**


## voidHold

Cancel a hold

### Example Usage

```php
<?php

declare(strict_types=1);
require_once 'vendor/autoload.php';

use \formance\stack\SDK;
use \formance\stack\Models\Shared\Security;
use \formance\stack\Models\Operations\VoidHoldRequest;

$sdk = SDK::builder()
    ->build();

try {
    $request = new VoidHoldRequest();
    $request->holdId = 'porro';

    $response = $sdk->wallets->voidHold($request);

    if ($response->statusCode === 200) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

### Parameters

| Parameter                                                                                       | Type                                                                                            | Required                                                                                        | Description                                                                                     |
| ----------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------- |
| `$request`                                                                                      | [\formance\stack\Models\Operations\VoidHoldRequest](../../models/operations/VoidHoldRequest.md) | :heavy_check_mark:                                                                              | The request object to use for the request.                                                      |


### Response

**[?\formance\stack\Models\Operations\VoidHoldResponse](../../models/operations/VoidHoldResponse.md)**


## walletsgetServerInfo

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
    $response = $sdk->wallets->walletsgetServerInfo();

    if ($response->serverInfo !== null) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```


### Response

**[?\formance\stack\Models\Operations\WalletsgetServerInfoResponse](../../models/operations/WalletsgetServerInfoResponse.md)**

