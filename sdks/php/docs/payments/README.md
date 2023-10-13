# payments

### Available Operations

* [connectorsTransfer](#connectorstransfer) - Transfer funds between Connector accounts
* [createBankAccount](#createbankaccount) - Create a BankAccount in Payments and on the PSP
* [createTransferInitiation](#createtransferinitiation) - Create a TransferInitiation
* [deleteTransferInitiation](#deletetransferinitiation) - Delete a transfer initiation
* [getAccountBalances](#getaccountbalances) - Get account balances
* [getBankAccount](#getbankaccount) - Get a bank account created by user on Formance
* [getConnectorTask](#getconnectortask) - Read a specific task of the connector
* [getPayment](#getpayment) - Get a payment
* [getTransferInitiation](#gettransferinitiation) - Get a transfer initiation
* [installConnector](#installconnector) - Install a connector
* [listAllConnectors](#listallconnectors) - List all installed connectors
* [listBankAccounts](#listbankaccounts) - List bank accounts created by user on Formance
* [listConfigsAvailableConnectors](#listconfigsavailableconnectors) - List the configs of each available connector
* [listConnectorTasks](#listconnectortasks) - List tasks from a connector
* [listConnectorsTransfers](#listconnectorstransfers) - List transfers and their statuses
* [listPayments](#listpayments) - List payments
* [listTransferInitiations](#listtransferinitiations) - List Transfer Initiations
* [paymentsgetAccount](#paymentsgetaccount) - Get an account
* [paymentsgetServerInfo](#paymentsgetserverinfo) - Get server info
* [paymentslistAccounts](#paymentslistaccounts) - List accounts
* [readConnectorConfig](#readconnectorconfig) - Read the config of a connector
* [resetConnector](#resetconnector) - Reset a connector
* [retryTransferInitiation](#retrytransferinitiation) - Retry a failed transfer initiation
* [udpateTransferInitiationStatus](#udpatetransferinitiationstatus) - Update the status of a transfer initiation
* [uninstallConnector](#uninstallconnector) - Uninstall a connector
* [updateMetadata](#updatemetadata) - Update metadata

## connectorsTransfer

Execute a transfer between two accounts.

### Example Usage

```php
<?php

declare(strict_types=1);
require_once 'vendor/autoload.php';

use \formance\stack\SDK;
use \formance\stack\Models\Shared\Security;
use \formance\stack\Models\Operations\ConnectorsTransferRequest;
use \formance\stack\Models\Shared\TransferRequest;
use \formance\stack\Models\Shared\Connector;

$sdk = SDK::builder()
    ->build();

try {
    $request = new ConnectorsTransferRequest();
    $request->transferRequest = new TransferRequest();
    $request->transferRequest->amount = 100;
    $request->transferRequest->asset = 'USD';
    $request->transferRequest->destination = 'acct_1Gqj58KZcSIg2N2q';
    $request->transferRequest->source = 'acct_1Gqj58KZcSIg2N2q';
    $request->connector = Connector::STRIPE;

    $response = $sdk->payments->connectorsTransfer($request);

    if ($response->transferResponse !== null) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

## createBankAccount

Create a bank account in Payments and on the PSP.

### Example Usage

```php
<?php

declare(strict_types=1);
require_once 'vendor/autoload.php';

use \formance\stack\SDK;
use \formance\stack\Models\Shared\Security;
use \formance\stack\Models\Shared\BankAccountRequest;
use \formance\stack\Models\Shared\Connector;

$sdk = SDK::builder()
    ->build();

try {
    $request = new BankAccountRequest();
    $request->accountNumber = 'corporis';
    $request->country = 'GB';
    $request->iban = 'dolore';
    $request->name = 'My account';
    $request->provider = Connector::MODULR;
    $request->swiftBicCode = 'dicta';

    $response = $sdk->payments->createBankAccount($request);

    if ($response->bankAccountResponse !== null) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

## createTransferInitiation

Create a transfer initiation

### Example Usage

```php
<?php

declare(strict_types=1);
require_once 'vendor/autoload.php';

use \formance\stack\SDK;
use \formance\stack\Models\Shared\Security;
use \formance\stack\Models\Shared\TransferInitiationRequest;
use \formance\stack\Models\Shared\Connector;
use \formance\stack\Models\Shared\TransferInitiationRequestType;

$sdk = SDK::builder()
    ->build();

try {
    $request = new TransferInitiationRequest();
    $request->amount = 688661;
    $request->asset = 'USD';
    $request->createdAt = DateTime::createFromFormat('Y-m-d\TH:i:sP', '2022-02-13T15:01:52.114Z');
    $request->description = 'commodi';
    $request->destinationAccountID = 'repudiandae';
    $request->provider = Connector::STRIPE;
    $request->reference = 'XXX';
    $request->sourceAccountID = 'ipsum';
    $request->type = TransferInitiationRequestType::PAYOUT;
    $request->validated = false;

    $response = $sdk->payments->createTransferInitiation($request);

    if ($response->transferInitiationResponse !== null) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

## deleteTransferInitiation

Delete a transfer initiation by its id.

### Example Usage

```php
<?php

declare(strict_types=1);
require_once 'vendor/autoload.php';

use \formance\stack\SDK;
use \formance\stack\Models\Shared\Security;
use \formance\stack\Models\Operations\DeleteTransferInitiationRequest;

$sdk = SDK::builder()
    ->build();

try {
    $request = new DeleteTransferInitiationRequest();
    $request->transferId = 'molestias';

    $response = $sdk->payments->deleteTransferInitiation($request);

    if ($response->statusCode === 200) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

## getAccountBalances

Get account balances

### Example Usage

```php
<?php

declare(strict_types=1);
require_once 'vendor/autoload.php';

use \formance\stack\SDK;
use \formance\stack\Models\Shared\Security;
use \formance\stack\Models\Operations\GetAccountBalancesRequest;

$sdk = SDK::builder()
    ->build();

try {
    $request = new GetAccountBalancesRequest();
    $request->accountId = 'excepturi';
    $request->asset = 'pariatur';
    $request->cursor = 'aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ==';
    $request->from = DateTime::createFromFormat('Y-m-d\TH:i:sP', '2022-06-29T05:25:54.356Z');
    $request->limit = 523248;
    $request->pageSize = 916723;
    $request->sort = [
        'repudiandae',
    ];
    $request->to = DateTime::createFromFormat('Y-m-d\TH:i:sP', '2022-11-01T07:52:08.326Z');

    $response = $sdk->payments->getAccountBalances($request);

    if ($response->balancesCursor !== null) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

## getBankAccount

Get a bank account created by user on Formance

### Example Usage

```php
<?php

declare(strict_types=1);
require_once 'vendor/autoload.php';

use \formance\stack\SDK;
use \formance\stack\Models\Shared\Security;
use \formance\stack\Models\Operations\GetBankAccountRequest;

$sdk = SDK::builder()
    ->build();

try {
    $request = new GetBankAccountRequest();
    $request->bankAccountId = 'itaque';

    $response = $sdk->payments->getBankAccount($request);

    if ($response->bankAccountResponse !== null) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

## getConnectorTask

Get a specific task associated to the connector.

### Example Usage

```php
<?php

declare(strict_types=1);
require_once 'vendor/autoload.php';

use \formance\stack\SDK;
use \formance\stack\Models\Shared\Security;
use \formance\stack\Models\Operations\GetConnectorTaskRequest;
use \formance\stack\Models\Shared\Connector;

$sdk = SDK::builder()
    ->build();

try {
    $request = new GetConnectorTaskRequest();
    $request->connector = Connector::WISE;
    $request->taskId = 'enim';

    $response = $sdk->payments->getConnectorTask($request);

    if ($response->taskResponse !== null) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

## getPayment

Get a payment

### Example Usage

```php
<?php

declare(strict_types=1);
require_once 'vendor/autoload.php';

use \formance\stack\SDK;
use \formance\stack\Models\Shared\Security;
use \formance\stack\Models\Operations\GetPaymentRequest;

$sdk = SDK::builder()
    ->build();

try {
    $request = new GetPaymentRequest();
    $request->paymentId = 'consequatur';

    $response = $sdk->payments->getPayment($request);

    if ($response->paymentResponse !== null) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

## getTransferInitiation

Get a transfer initiation

### Example Usage

```php
<?php

declare(strict_types=1);
require_once 'vendor/autoload.php';

use \formance\stack\SDK;
use \formance\stack\Models\Shared\Security;
use \formance\stack\Models\Operations\GetTransferInitiationRequest;

$sdk = SDK::builder()
    ->build();

try {
    $request = new GetTransferInitiationRequest();
    $request->transferId = 'est';

    $response = $sdk->payments->getTransferInitiation($request);

    if ($response->transferInitiationResponse !== null) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

## installConnector

Install a connector by its name and config.

### Example Usage

```php
<?php

declare(strict_types=1);
require_once 'vendor/autoload.php';

use \formance\stack\SDK;
use \formance\stack\Models\Shared\Security;
use \formance\stack\Models\Operations\InstallConnectorRequest;
use \formance\stack\Models\Shared\Connector;

$sdk = SDK::builder()
    ->build();

try {
    $request = new InstallConnectorRequest();
    $request->requestBody = new MangoPayConfig();
    $request->requestBody->apiKey = 'XXX';
    $request->requestBody->clientID = 'XXX';
    $request->requestBody->endpoint = 'XXX';
    $request->requestBody->pollingPeriod = '60s';
    $request->connector = Connector::DUMMY_PAY;

    $response = $sdk->payments->installConnector($request);

    if ($response->statusCode === 200) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

## listAllConnectors

List all installed connectors.

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
    $response = $sdk->payments->listAllConnectors();

    if ($response->connectorsResponse !== null) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

## listBankAccounts

List all bank accounts created by user on Formance.

### Example Usage

```php
<?php

declare(strict_types=1);
require_once 'vendor/autoload.php';

use \formance\stack\SDK;
use \formance\stack\Models\Shared\Security;
use \formance\stack\Models\Operations\ListBankAccountsRequest;

$sdk = SDK::builder()
    ->build();

try {
    $request = new ListBankAccountsRequest();
    $request->cursor = 'aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ==';
    $request->pageSize = 647174;
    $request->sort = [
        'quibusdam',
        'labore',
        'modi',
    ];

    $response = $sdk->payments->listBankAccounts($request);

    if ($response->bankAccountsCursor !== null) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

## listConfigsAvailableConnectors

List the configs of each available connector.

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
    $response = $sdk->payments->listConfigsAvailableConnectors();

    if ($response->connectorsConfigsResponse !== null) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

## listConnectorTasks

List all tasks associated with this connector.

### Example Usage

```php
<?php

declare(strict_types=1);
require_once 'vendor/autoload.php';

use \formance\stack\SDK;
use \formance\stack\Models\Shared\Security;
use \formance\stack\Models\Operations\ListConnectorTasksRequest;
use \formance\stack\Models\Shared\Connector;

$sdk = SDK::builder()
    ->build();

try {
    $request = new ListConnectorTasksRequest();
    $request->connector = Connector::DUMMY_PAY;
    $request->cursor = 'aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ==';
    $request->pageSize = 397821;

    $response = $sdk->payments->listConnectorTasks($request);

    if ($response->tasksCursor !== null) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

## listConnectorsTransfers

List transfers

### Example Usage

```php
<?php

declare(strict_types=1);
require_once 'vendor/autoload.php';

use \formance\stack\SDK;
use \formance\stack\Models\Shared\Security;
use \formance\stack\Models\Operations\ListConnectorsTransfersRequest;
use \formance\stack\Models\Shared\Connector;

$sdk = SDK::builder()
    ->build();

try {
    $request = new ListConnectorsTransfersRequest();
    $request->connector = Connector::CURRENCY_CLOUD;

    $response = $sdk->payments->listConnectorsTransfers($request);

    if ($response->transfersResponse !== null) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

## listPayments

List payments

### Example Usage

```php
<?php

declare(strict_types=1);
require_once 'vendor/autoload.php';

use \formance\stack\SDK;
use \formance\stack\Models\Shared\Security;
use \formance\stack\Models\Operations\ListPaymentsRequest;

$sdk = SDK::builder()
    ->build();

try {
    $request = new ListPaymentsRequest();
    $request->cursor = 'aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ==';
    $request->pageSize = 552822;
    $request->sort = [
        'magni',
    ];

    $response = $sdk->payments->listPayments($request);

    if ($response->paymentsCursor !== null) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

## listTransferInitiations

List Transfer Initiations

### Example Usage

```php
<?php

declare(strict_types=1);
require_once 'vendor/autoload.php';

use \formance\stack\SDK;
use \formance\stack\Models\Shared\Security;
use \formance\stack\Models\Operations\ListTransferInitiationsRequest;

$sdk = SDK::builder()
    ->build();

try {
    $request = new ListTransferInitiationsRequest();
    $request->cursor = 'aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ==';
    $request->pageSize = 828940;
    $request->query = 'ipsam';
    $request->sort = [
        'fugit',
    ];

    $response = $sdk->payments->listTransferInitiations($request);

    if ($response->transferInitiationsCursor !== null) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

## paymentsgetAccount

Get an account

### Example Usage

```php
<?php

declare(strict_types=1);
require_once 'vendor/autoload.php';

use \formance\stack\SDK;
use \formance\stack\Models\Shared\Security;
use \formance\stack\Models\Operations\PaymentsgetAccountRequest;

$sdk = SDK::builder()
    ->build();

try {
    $request = new PaymentsgetAccountRequest();
    $request->accountId = 'dolorum';

    $response = $sdk->payments->paymentsgetAccount($request);

    if ($response->paymentsAccountResponse !== null) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

## paymentsgetServerInfo

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
    $response = $sdk->payments->paymentsgetServerInfo();

    if ($response->serverInfo !== null) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

## paymentslistAccounts

List accounts

### Example Usage

```php
<?php

declare(strict_types=1);
require_once 'vendor/autoload.php';

use \formance\stack\SDK;
use \formance\stack\Models\Shared\Security;
use \formance\stack\Models\Operations\PaymentslistAccountsRequest;

$sdk = SDK::builder()
    ->build();

try {
    $request = new PaymentslistAccountsRequest();
    $request->cursor = 'aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ==';
    $request->pageSize = 569618;
    $request->sort = [
        'facilis',
        'tempore',
    ];

    $response = $sdk->payments->paymentslistAccounts($request);

    if ($response->accountsCursor !== null) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

## readConnectorConfig

Read connector config

### Example Usage

```php
<?php

declare(strict_types=1);
require_once 'vendor/autoload.php';

use \formance\stack\SDK;
use \formance\stack\Models\Shared\Security;
use \formance\stack\Models\Operations\ReadConnectorConfigRequest;
use \formance\stack\Models\Shared\Connector;

$sdk = SDK::builder()
    ->build();

try {
    $request = new ReadConnectorConfigRequest();
    $request->connector = Connector::WISE;

    $response = $sdk->payments->readConnectorConfig($request);

    if ($response->connectorConfigResponse !== null) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

## resetConnector

Reset a connector by its name.
It will remove the connector and ALL PAYMENTS generated with it.


### Example Usage

```php
<?php

declare(strict_types=1);
require_once 'vendor/autoload.php';

use \formance\stack\SDK;
use \formance\stack\Models\Shared\Security;
use \formance\stack\Models\Operations\ResetConnectorRequest;
use \formance\stack\Models\Shared\Connector;

$sdk = SDK::builder()
    ->build();

try {
    $request = new ResetConnectorRequest();
    $request->connector = Connector::MONEYCORP;

    $response = $sdk->payments->resetConnector($request);

    if ($response->statusCode === 200) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

## retryTransferInitiation

Retry a failed transfer initiation

### Example Usage

```php
<?php

declare(strict_types=1);
require_once 'vendor/autoload.php';

use \formance\stack\SDK;
use \formance\stack\Models\Shared\Security;
use \formance\stack\Models\Operations\RetryTransferInitiationRequest;

$sdk = SDK::builder()
    ->build();

try {
    $request = new RetryTransferInitiationRequest();
    $request->transferId = 'eum';

    $response = $sdk->payments->retryTransferInitiation($request);

    if ($response->statusCode === 200) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

## udpateTransferInitiationStatus

Update a transfer initiation status

### Example Usage

```php
<?php

declare(strict_types=1);
require_once 'vendor/autoload.php';

use \formance\stack\SDK;
use \formance\stack\Models\Shared\Security;
use \formance\stack\Models\Operations\UdpateTransferInitiationStatusRequest;
use \formance\stack\Models\Shared\UpdateTransferInitiationStatusRequest;
use \formance\stack\Models\Shared\UpdateTransferInitiationStatusRequestStatus;

$sdk = SDK::builder()
    ->build();

try {
    $request = new UdpateTransferInitiationStatusRequest();
    $request->updateTransferInitiationStatusRequest = new UpdateTransferInitiationStatusRequest();
    $request->updateTransferInitiationStatusRequest->status = UpdateTransferInitiationStatusRequestStatus::PROCESSING;
    $request->transferId = 'eligendi';

    $response = $sdk->payments->udpateTransferInitiationStatus($request);

    if ($response->statusCode === 200) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

## uninstallConnector

Uninstall a connector by its name.

### Example Usage

```php
<?php

declare(strict_types=1);
require_once 'vendor/autoload.php';

use \formance\stack\SDK;
use \formance\stack\Models\Shared\Security;
use \formance\stack\Models\Operations\UninstallConnectorRequest;
use \formance\stack\Models\Shared\Connector;

$sdk = SDK::builder()
    ->build();

try {
    $request = new UninstallConnectorRequest();
    $request->connector = Connector::CURRENCY_CLOUD;

    $response = $sdk->payments->uninstallConnector($request);

    if ($response->statusCode === 200) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

## updateMetadata

Update metadata

### Example Usage

```php
<?php

declare(strict_types=1);
require_once 'vendor/autoload.php';

use \formance\stack\SDK;
use \formance\stack\Models\Shared\Security;
use \formance\stack\Models\Operations\UpdateMetadataRequest;
use \formance\stack\Models\Shared\PaymentMetadata;

$sdk = SDK::builder()
    ->build();

try {
    $request = new UpdateMetadataRequest();
    $request->paymentMetadata = new PaymentMetadata();
    $request->paymentMetadata->key = 'aliquid';
    $request->paymentId = 'provident';

    $response = $sdk->payments->updateMetadata($request);

    if ($response->statusCode === 200) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```
