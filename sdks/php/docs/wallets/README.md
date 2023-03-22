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
    $request->holdId = 'rerum';

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
    $request->createBalanceRequest->expiresAt = DateTime::createFromFormat('Y-m-d\TH:i:sP', '2022-09-14T10:27:07.590Z');
    $request->createBalanceRequest->name = 'Irving Jenkins';
    $request->createBalanceRequest->priority = 881104;
    $request->id = '395efb9b-a88f-43a6-a997-074ba4469b6e';

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
        'architecto' => 'magnam',
    ];
    $request->name = 'Miriam Hermann';

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
    $request->creditWalletRequest->amount->amount = 574325;
    $request->creditWalletRequest->amount->asset = 'accusantium';
    $request->creditWalletRequest->balance = 'mollitia';
    $request->creditWalletRequest->metadata = [
        'mollitia' => 'ad',
        'eum' => 'dolor',
        'necessitatibus' => 'odit',
        'nemo' => 'quasi',
    ];
    $request->creditWalletRequest->reference = 'iure';
    $request->creditWalletRequest->sources = [
        new WalletSubject(),
        new LedgerAccountSubject(),
        new WalletSubject(),
        new WalletSubject(),
    ];
    $request->id = 'b711e5b7-fd2e-4d02-8921-cddc692601fb';

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
    $request->debitWalletRequest->amount->amount = 373291;
    $request->debitWalletRequest->amount->asset = 'voluptate';
    $request->debitWalletRequest->balances = [
        'nam',
        'eaque',
    ];
    $request->debitWalletRequest->description = 'pariatur';
    $request->debitWalletRequest->destination = new LedgerAccountSubject();
    $request->debitWalletRequest->destination->identifier = 'voluptatibus';
    $request->debitWalletRequest->destination->type = 'perferendis';
    $request->debitWalletRequest->metadata = [
        'amet' => 'aut',
        'cumque' => 'corporis',
        'hic' => 'libero',
        'nobis' => 'dolores',
    ];
    $request->debitWalletRequest->pending = false;
    $request->id = '58705320-2c73-4d5f-a9b9-0c28909b3fe4';

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
    $request->balanceName = 'iste';
    $request->id = 'a8d9cbf4-8633-4323-b9b7-7f3a4100674e';

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
    $request->holdID = 'quidem';

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
        'voluptas' => 'natus',
        'eos' => 'atque',
        'sit' => 'fugiat',
        'ab' => 'soluta',
    ];
    $request->pageSize = 679393;
    $request->walletID = 'iusto';

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
    $request->pageSize = 453697;
    $request->walletID = 'dolorum';

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
    $request->id = '89ebf737-ae42-403c-a5e6-a95d8a0d446c';

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
    $request->id = 'e2af7a73-cf3b-4e45-bf87-0b326b5a7342';

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
    $request->id = '9cdb1a84-22bb-4679-9232-2715bf0cbb1e';

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
        'veritatis' => 'nobis',
    ];
    $request->name = 'Dr. Randolph McDermott';
    $request->pageSize = 292147;

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
        'adipisci' => 'dolorum',
        'architecto' => 'quae',
    ];
    $request->id = '08e0adcf-4b92-4187-9fce-953f73ef7fbc';

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
    $request->holdId = 'odio';

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
