# transactions

### Available Operations

* [addMetadataOnTransaction](#addmetadataontransaction) - Set the metadata of a transaction by its ID
* [countTransactions](#counttransactions) - Count the transactions from a ledger
* [createTransaction](#createtransaction) - Create a new transaction to a ledger
* [getTransaction](#gettransaction) - Get transaction from a ledger by its ID
* [listTransactions](#listtransactions) - List transactions from a ledger
* [revertTransaction](#reverttransaction) - Revert a ledger transaction by its ID

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
    $request->idempotencyKey = 'incidunt';
    $request->requestBody = [
        'consequatur' => 'est',
        'quibusdam' => 'explicabo',
    ];
    $request->async = true;
    $request->dryRun = true;
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
    $request->endTime = DateTime::createFromFormat('Y-m-d\TH:i:sP', '2021-07-27T01:56:50.693Z');
    $request->ledger = 'ledger001';
    $request->metadata = [
        'labore' => 'modi',
        'qui' => 'aliquid',
        'cupiditate' => 'quos',
        'perferendis' => 'magni',
    ];
    $request->reference = 'ref:001';
    $request->source = 'users:001';
    $request->startTime = DateTime::createFromFormat('Y-m-d\TH:i:sP', '2021-11-22T01:26:35.048Z');

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
    $request->idempotencyKey = 'alias';
    $request->postTransaction = new PostTransaction();
    $request->postTransaction->metadata = [
        'dolorum' => 'excepturi',
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
        'tempore' => 'labore',
        'delectus' => 'eum',
        'non' => 'eligendi',
    ];
    $request->postTransaction->timestamp = DateTime::createFromFormat('Y-m-d\TH:i:sP', '2022-03-17T20:21:28.792Z');
    $request->async = true;
    $request->dryRun = true;
    $request->ledger = 'ledger001';

    $response = $sdk->transactions->createTransaction($request);

    if ($response->createTransactionResponse !== null) {
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

    if ($response->getTransactionResponse !== null) {
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
    $request->cursor = 'aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ==';
    $request->destination = 'users:001';
    $request->endTime = DateTime::createFromFormat('Y-m-d\TH:i:sP', '2021-03-17T21:24:26.606Z');
    $request->ledger = 'ledger001';
    $request->metadata = [
        'officia' => 'dolor',
        'debitis' => 'a',
        'dolorum' => 'in',
    ];
    $request->pageSize = 449198;
    $request->reference = 'ref:001';
    $request->source = 'users:001';
    $request->startTime = DateTime::createFromFormat('Y-m-d\TH:i:sP', '2020-01-25T11:09:22.009Z');

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

    if ($response->revertTransactionResponse !== null) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```
