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
    $request->holdId = 'omnis';

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
    $request->createBalanceRequest->expiresAt = DateTime::createFromFormat('Y-m-d\TH:i:sP', '2022-12-24T23:52:02.245Z');
    $request->createBalanceRequest->name = 'Megan Rau';
    $request->createBalanceRequest->priority = 290077;
    $request->id = '69b6e214-1959-4890-afa5-63e2516fe4c8';

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
        'in' => 'architecto',
        'architecto' => 'repudiandae',
        'ullam' => 'expedita',
    ];
    $request->name = 'Kristie Spencer';

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
    $request->creditWalletRequest->amount->amount = 868126;
    $request->creditWalletRequest->amount->asset = 'accusantium';
    $request->creditWalletRequest->balance = 'consequuntur';
    $request->creditWalletRequest->metadata = [
        'natus' => 'magni',
        'sunt' => 'quo',
        'illum' => 'pariatur',
    ];
    $request->creditWalletRequest->reference = 'maxime';
    $request->creditWalletRequest->sources = [
        new WalletSubject(),
        new LedgerAccountSubject(),
    ];
    $request->id = '601fb576-b0d5-4f0d-b0c5-fbb258705320';

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
    $request->debitWalletRequest->amount->amount = 170986;
    $request->debitWalletRequest->amount->asset = 'minus';
    $request->debitWalletRequest->balances = [
        'dolor',
        'vero',
    ];
    $request->debitWalletRequest->description = 'nostrum';
    $request->debitWalletRequest->destination = new WalletSubject();
    $request->debitWalletRequest->destination->balance = 'recusandae';
    $request->debitWalletRequest->destination->identifier = 'omnis';
    $request->debitWalletRequest->destination->type = 'facilis';
    $request->debitWalletRequest->metadata = [
        'voluptatem' => 'porro',
        'consequuntur' => 'blanditiis',
        'error' => 'eaque',
    ];
    $request->debitWalletRequest->pending = false;
    $request->id = '9b3fe49a-8d9c-4bf4-8633-323f9b77f3a4';

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
    $request->balanceName = 'veritatis';
    $request->id = '00674ebf-6928-40d1-ba77-a89ebf737ae4';

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
    $request->holdID = 'aspernatur';

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
        'amet' => 'optio',
    ];
    $request->pageSize = 881586;
    $request->walletID = 'ad';

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
    $request->pageSize = 904425;
    $request->walletID = 'suscipit';

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
    $request->id = 'a95d8a0d-446c-4e2a-b7a7-3cf3be453f87';

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
    $request->id = '0b326b5a-7342-49cd-b1a8-422bb679d232';

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
    $request->id = '2715bf0c-bb1e-431b-8b90-f3443a1108e0';

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
        'repellendus' => 'porro',
        'doloribus' => 'ut',
        'facilis' => 'cupiditate',
    ];
    $request->name = 'Amy Langworth';
    $request->pageSize = 977496;

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
        'vero' => 'omnis',
        'quis' => 'ipsum',
        'delectus' => 'voluptate',
        'consectetur' => 'vero',
    ];
    $request->id = 'f7fbc7ab-d74d-4d39-80f5-d2cff7c70a45';

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
    $request->holdId = 'ea';

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
