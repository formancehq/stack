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
    $request->holdId = 'blanditiis';

    $response = $sdk->wallets->confirmHold($request);

    if ($response->statusCode === 200) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

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
    $request->createBalanceRequest->expiresAt = DateTime::createFromFormat('Y-m-d\TH:i:sP', '2021-02-02T01:24:52.629Z');
    $request->createBalanceRequest->name = 'Sandy Huels';
    $request->createBalanceRequest->priority = 606393;
    $request->id = '7074ba44-69b6-4e21-8195-9890afa563e2';

    $response = $sdk->wallets->createBalance($request);

    if ($response->createBalanceResponse !== null) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

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
        'quasi' => 'iure',
        'doloribus' => 'debitis',
    ];
    $request->name = 'Jasmine Lind';

    $response = $sdk->wallets->createWallet($request);

    if ($response->createWalletResponse !== null) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

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
    $request->creditWalletRequest->amount->amount = 100226;
    $request->creditWalletRequest->amount->asset = 'architecto';
    $request->creditWalletRequest->balance = 'repudiandae';
    $request->creditWalletRequest->metadata = [
        'expedita' => 'nihil',
        'repellat' => 'quibusdam',
    ];
    $request->creditWalletRequest->reference = 'sed';
    $request->creditWalletRequest->sources = [
        new WalletSubject(),
        new LedgerAccountSubject(),
        new LedgerAccountSubject(),
        new WalletSubject(),
    ];
    $request->id = '921cddc6-9260-41fb-976b-0d5f0d30c5fb';

    $response = $sdk->wallets->creditWallet($request);

    if ($response->statusCode === 200) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

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
    $request->debitWalletRequest->amount->amount = 749999;
    $request->debitWalletRequest->amount->asset = 'dolores';
    $request->debitWalletRequest->balances = [
        'totam',
        'dignissimos',
    ];
    $request->debitWalletRequest->description = 'eaque';
    $request->debitWalletRequest->destination = new LedgerAccountSubject();
    $request->debitWalletRequest->destination->identifier = 'nesciunt';
    $request->debitWalletRequest->destination->type = 'eos';
    $request->debitWalletRequest->metadata = [
        'dolores' => 'minus',
    ];
    $request->debitWalletRequest->pending = false;
    $request->id = '73d5fe9b-90c2-4890-9b3f-e49a8d9cbf48';

    $response = $sdk->wallets->debitWallet($request);

    if ($response->debitWalletResponse !== null) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

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
    $request->balanceName = 'aliquid';
    $request->id = '33323f9b-77f3-4a41-8067-4ebf69280d1b';

    $response = $sdk->wallets->getBalance($request);

    if ($response->getBalanceResponse !== null) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

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
    $request->holdID = 'dolorum';

    $response = $sdk->wallets->getHold($request);

    if ($response->getHoldResponse !== null) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

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
        'voluptate' => 'dolorum',
        'deleniti' => 'omnis',
    ];
    $request->pageSize = 896672;
    $request->walletID = 'distinctio';

    $response = $sdk->wallets->getHolds($request);

    if ($response->getHoldsResponse !== null) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

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
    $request->pageSize = 990339;
    $request->walletID = 'nihil';

    $response = $sdk->wallets->getTransactions($request);

    if ($response->getTransactionsResponse !== null) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

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
    $request->id = '37ae4203-ce5e-46a9-9d8a-0d446ce2af7a';

    $response = $sdk->wallets->getWallet($request);

    if ($response->getWalletResponse !== null) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

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
    $request->id = '73cf3be4-53f8-470b-b26b-5a73429cdb1a';

    $response = $sdk->wallets->getWalletSummary($request);

    if ($response->getWalletSummaryResponse !== null) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

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
    $request->id = '8422bb67-9d23-4227-95bf-0cbb1e31b8b9';

    $response = $sdk->wallets->listBalances($request);

    if ($response->listBalancesResponse !== null) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

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
        'delectus' => 'dolorem',
    ];
    $request->name = 'Clara Fisher Jr.';
    $request->pageSize = 16429;

    $response = $sdk->wallets->listWallets($request);

    if ($response->listWalletsResponse !== null) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

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
        'itaque' => 'consequatur',
        'est' => 'repellendus',
        'porro' => 'doloribus',
    ];
    $request->id = '4b921879-fce9-453f-b3ef-7fbc7abd74dd';

    $response = $sdk->wallets->updateWallet($request);

    if ($response->statusCode === 200) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

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
    $request->holdId = 'sequi';

    $response = $sdk->wallets->voidHold($request);

    if ($response->statusCode === 200) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

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
