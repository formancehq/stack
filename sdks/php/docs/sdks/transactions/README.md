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

### Parameters

| Parameter                                                                                                           | Type                                                                                                                | Required                                                                                                            | Description                                                                                                         |
| ------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------- |
| `$request`                                                                                                          | [\formance\stack\Models\Operations\CreateTransactionsRequest](../../models/operations/CreateTransactionsRequest.md) | :heavy_check_mark:                                                                                                  | The request object to use for the request.                                                                          |


### Response

**[?\formance\stack\Models\Operations\CreateTransactionsResponse](../../models/operations/CreateTransactionsResponse.md)**


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
        'deserunt' => 'distinctio',
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

### Parameters

| Parameter                                                                                                                       | Type                                                                                                                            | Required                                                                                                                        | Description                                                                                                                     |
| ------------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------- |
| `$request`                                                                                                                      | [\formance\stack\Models\Operations\AddMetadataOnTransactionRequest](../../models/operations/AddMetadataOnTransactionRequest.md) | :heavy_check_mark:                                                                                                              | The request object to use for the request.                                                                                      |


### Response

**[?\formance\stack\Models\Operations\AddMetadataOnTransactionResponse](../../models/operations/AddMetadataOnTransactionResponse.md)**


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
use \formance\stack\Models\Operations\CountTransactionsMetadata;

$sdk = SDK::builder()
    ->build();

try {
    $request = new CountTransactionsRequest();
    $request->account = 'users:001';
    $request->destination = 'users:001';
    $request->endTime = DateTime::createFromFormat('Y-m-d\TH:i:sP', '2022-02-18T02:24:27.835Z');
    $request->ledger = 'ledger001';
    $request->metadata = new CountTransactionsMetadata();
    $request->reference = 'ref:001';
    $request->source = 'users:001';
    $request->startTime = DateTime::createFromFormat('Y-m-d\TH:i:sP', '2022-10-26T03:14:36.345Z');

    $response = $sdk->transactions->countTransactions($request);

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
| `$request`                                                                                                        | [\formance\stack\Models\Operations\CountTransactionsRequest](../../models/operations/CountTransactionsRequest.md) | :heavy_check_mark:                                                                                                | The request object to use for the request.                                                                        |


### Response

**[?\formance\stack\Models\Operations\CountTransactionsResponse](../../models/operations/CountTransactionsResponse.md)**


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
use \formance\stack\Models\Shared\PostTransactionScriptVars;

$sdk = SDK::builder()
    ->build();

try {
    $request = new CreateTransactionRequest();
    $request->idempotencyKey = 'aliquid';
    $request->postTransaction = new PostTransaction();
    $request->postTransaction->metadata = [
        'quos' => 'perferendis',
        'magni' => 'assumenda',
        'ipsam' => 'alias',
    ];
    $request->postTransaction->postings = [
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
    $request->postTransaction->script->vars = new PostTransactionScriptVars();
    $request->postTransaction->timestamp = DateTime::createFromFormat('Y-m-d\TH:i:sP', '2021-11-11T04:17:07.569Z');
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

### Parameters

| Parameter                                                                                                         | Type                                                                                                              | Required                                                                                                          | Description                                                                                                       |
| ----------------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------- |
| `$request`                                                                                                        | [\formance\stack\Models\Operations\CreateTransactionRequest](../../models/operations/CreateTransactionRequest.md) | :heavy_check_mark:                                                                                                | The request object to use for the request.                                                                        |


### Response

**[?\formance\stack\Models\Operations\CreateTransactionResponse](../../models/operations/CreateTransactionResponse.md)**


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

### Parameters

| Parameter                                                                                                   | Type                                                                                                        | Required                                                                                                    | Description                                                                                                 |
| ----------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------- |
| `$request`                                                                                                  | [\formance\stack\Models\Operations\GetTransactionRequest](../../models/operations/GetTransactionRequest.md) | :heavy_check_mark:                                                                                          | The request object to use for the request.                                                                  |


### Response

**[?\formance\stack\Models\Operations\GetTransactionResponse](../../models/operations/GetTransactionResponse.md)**


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
use \formance\stack\Models\Operations\ListTransactionsMetadata;

$sdk = SDK::builder()
    ->build();

try {
    $request = new ListTransactionsRequest();
    $request->account = 'users:001';
    $request->after = 'tempora';
    $request->cursor = 'aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ==';
    $request->destination = 'users:001';
    $request->endTime = DateTime::createFromFormat('Y-m-d\TH:i:sP', '2021-07-13T07:24:02.478Z');
    $request->ledger = 'ledger001';
    $request->metadata = new ListTransactionsMetadata();
    $request->pageSize = 288476;
    $request->reference = 'ref:001';
    $request->source = 'users:001';
    $request->startTime = DateTime::createFromFormat('Y-m-d\TH:i:sP', '2021-09-13T13:11:22.288Z');

    $response = $sdk->transactions->listTransactions($request);

    if ($response->transactionsCursorResponse !== null) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

### Parameters

| Parameter                                                                                                       | Type                                                                                                            | Required                                                                                                        | Description                                                                                                     |
| --------------------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------- |
| `$request`                                                                                                      | [\formance\stack\Models\Operations\ListTransactionsRequest](../../models/operations/ListTransactionsRequest.md) | :heavy_check_mark:                                                                                              | The request object to use for the request.                                                                      |


### Response

**[?\formance\stack\Models\Operations\ListTransactionsResponse](../../models/operations/ListTransactionsResponse.md)**


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

### Parameters

| Parameter                                                                                                         | Type                                                                                                              | Required                                                                                                          | Description                                                                                                       |
| ----------------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------- |
| `$request`                                                                                                        | [\formance\stack\Models\Operations\RevertTransactionRequest](../../models/operations/RevertTransactionRequest.md) | :heavy_check_mark:                                                                                                | The request object to use for the request.                                                                        |


### Response

**[?\formance\stack\Models\Operations\RevertTransactionResponse](../../models/operations/RevertTransactionResponse.md)**

