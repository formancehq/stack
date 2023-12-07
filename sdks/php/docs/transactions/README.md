# transactions

### Available Operations

* [createTransactions](#createtransactions) - Create a new batch of transactions to a ledger
* [addMetadataOnTransaction](#addmetadataontransaction) - Set the metadata of a transaction by its ID
* [countTransactions](#counttransactions) - Count the transactions from a ledger
* [createTransaction](#createtransaction) - Create a new transaction to a ledger
* [getTransaction](#gettransaction) - Get transaction from a ledger by its ID
* [listTransactions](#listtransactions) - List transactions from a ledger
* [revertTransaction](#reverttransaction) - Revert a ledger transaction by its ID

## createTransactions

Create a new batch of transactions to a ledger

### Example Usage

```php
<?php

declare(strict_types=1);
require_once 'vendor/autoload.php';

use \formance\stack\SDK;
use \formance\stack\Models\Shared\Security;
use \formance\stack\Models\Operations\CreateTransactionsRequest;
use \formance\stack\Models\Shared\Transactions;
use \formance\stack\Models\Shared\TransactionData;
use \formance\stack\Models\Shared\Posting;

$sdk = SDK::builder()
    ->build();

try {
    $request = new CreateTransactionsRequest();
    $request->transactions = new Transactions();
    $request->transactions->transactions = [
        new TransactionData(),
        new TransactionData(),
        new TransactionData(),
    ];
    $request->ledger = 'ledger001';

    $response = $sdk->transactions->createTransactions($request);

    if ($response->transactionsResponse !== null) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

## addMetadataOnTransaction

Set the metadata of a transaction by its ID

### Example Usage

```php
<?php

declare(strict_types=1);
require_once 'vendor/autoload.php';

use \formance\stack\SDK;
use \formance\stack\Models\Shared\Security;
use \formance\stack\Models\Operations\AddMetadataOnTransactionRequest;

$sdk = SDK::builder()
    ->build();

try {
    $request = new AddMetadataOnTransactionRequest();
    $request->requestBody = [
        'sint' => 'accusantium',
        'mollitia' => 'reiciendis',
        'mollitia' => 'ad',
    ];
    $request->ledger = 'ledger001';
    $request->txid = 1234;

    $response = $sdk->transactions->addMetadataOnTransaction($request);

    if ($response->statusCode === 200) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

## countTransactions

Count the transactions from a ledger

### Example Usage

```php
<?php

declare(strict_types=1);
require_once 'vendor/autoload.php';

use \formance\stack\SDK;
use \formance\stack\Models\Shared\Security;
use \formance\stack\Models\Operations\CountTransactionsRequest;

$sdk = SDK::builder()
    ->build();

try {
    $request = new CountTransactionsRequest();
    $request->account = 'users:001';
    $request->destination = 'users:001';
    $request->endTime = DateTime::createFromFormat('Y-m-d\TH:i:sP', '2022-10-12T05:44:19.260Z');
    $request->ledger = 'ledger001';
    $request->metadata = [
        'odit' => 'nemo',
        'quasi' => 'iure',
        'doloribus' => 'debitis',
        'eius' => 'maxime',
    ];
    $request->reference = 'ref:001';
    $request->source = 'users:001';
    $request->startTime = DateTime::createFromFormat('Y-m-d\TH:i:sP', '2021-08-05T03:52:18.835Z');

    $response = $sdk->transactions->countTransactions($request);

    if ($response->statusCode === 200) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

## createTransaction

Create a new transaction to a ledger

### Example Usage

```php
<?php

declare(strict_types=1);
require_once 'vendor/autoload.php';

use \formance\stack\SDK;
use \formance\stack\Models\Shared\Security;
use \formance\stack\Models\Operations\CreateTransactionRequest;
use \formance\stack\Models\Shared\PostTransaction;
use \formance\stack\Models\Shared\Posting;
use \formance\stack\Models\Shared\PostTransactionScript;

$sdk = SDK::builder()
    ->build();

try {
    $request = new CreateTransactionRequest();
    $request->idempotencyKey = 'in';
    $request->postTransaction = new PostTransaction();
    $request->postTransaction->metadata = [
        'architecto' => 'repudiandae',
    ];
    $request->postTransaction->postings = [
        new Posting(),
        new Posting(),
    ];
    $request->postTransaction->reference = 'ref:001';
    $request->postTransaction->script = new PostTransactionScript();
    $request->postTransaction->script->plain = 'vars {
    account $user
    }
    send [COIN 10] (
    	source = @world
    	destination = $user
    )
    ';
    $request->postTransaction->script->vars = [
        'nihil' => 'repellat',
        'quibusdam' => 'sed',
        'saepe' => 'pariatur',
    ];
    $request->postTransaction->timestamp = DateTime::createFromFormat('Y-m-d\TH:i:sP', '2022-11-02T16:33:42.080Z');
    $request->ledger = 'ledger001';
    $request->preview = true;

    $response = $sdk->transactions->createTransaction($request);

    if ($response->transactionsResponse !== null) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

## getTransaction

Get transaction from a ledger by its ID

### Example Usage

```php
<?php

declare(strict_types=1);
require_once 'vendor/autoload.php';

use \formance\stack\SDK;
use \formance\stack\Models\Shared\Security;
use \formance\stack\Models\Operations\GetTransactionRequest;

$sdk = SDK::builder()
    ->build();

try {
    $request = new GetTransactionRequest();
    $request->ledger = 'ledger001';
    $request->txid = 1234;

    $response = $sdk->transactions->getTransaction($request);

    if ($response->transactionResponse !== null) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

## listTransactions

List transactions from a ledger, sorted by txid in descending order.

### Example Usage

```php
<?php

declare(strict_types=1);
require_once 'vendor/autoload.php';

use \formance\stack\SDK;
use \formance\stack\Models\Shared\Security;
use \formance\stack\Models\Operations\ListTransactionsRequest;

$sdk = SDK::builder()
    ->build();

try {
    $request = new ListTransactionsRequest();
    $request->account = 'users:001';
    $request->after = 'praesentium';
    $request->cursor = 'aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ==';
    $request->destination = 'users:001';
    $request->endTime = DateTime::createFromFormat('Y-m-d\TH:i:sP', '2022-09-01T04:49:52.515Z');
    $request->ledger = 'ledger001';
    $request->metadata = [
        'quo' => 'illum',
    ];
    $request->pageSize = 864934;
    $request->reference = 'ref:001';
    $request->source = 'users:001';
    $request->startTime = DateTime::createFromFormat('Y-m-d\TH:i:sP', '2021-10-07T12:29:38.109Z');

    $response = $sdk->transactions->listTransactions($request);

    if ($response->transactionsCursorResponse !== null) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

## revertTransaction

Revert a ledger transaction by its ID

### Example Usage

```php
<?php

declare(strict_types=1);
require_once 'vendor/autoload.php';

use \formance\stack\SDK;
use \formance\stack\Models\Shared\Security;
use \formance\stack\Models\Operations\RevertTransactionRequest;

$sdk = SDK::builder()
    ->build();

try {
    $request = new RevertTransactionRequest();
    $request->ledger = 'ledger001';
    $request->txid = 1234;

    $response = $sdk->transactions->revertTransaction($request);

    if ($response->transactionResponse !== null) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```
