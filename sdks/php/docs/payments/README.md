# payments

### Available Operations

* [connectorsStripeTransfer](#connectorsstripetransfer) - Transfer funds between Stripe accounts
* [connectorsTransfer](#connectorstransfer) - Transfer funds between Connector accounts
* [getAccountBalances](#getaccountbalances) - Get account balances
* [getConnectorTask](#getconnectortask) - Read a specific task of the connector
* [getPayment](#getpayment) - Get a payment
* [installConnector](#installconnector) - Install a connector
* [listAllConnectors](#listallconnectors) - List all installed connectors
* [listConfigsAvailableConnectors](#listconfigsavailableconnectors) - List the configs of each available connector
* [listConnectorTasks](#listconnectortasks) - List tasks from a connector
* [listConnectorsTransfers](#listconnectorstransfers) - List transfers and their statuses
* [listPayments](#listpayments) - List payments
* [paymentsgetAccount](#paymentsgetaccount) - Get an account
* [paymentsgetServerInfo](#paymentsgetserverinfo) - Get server info
* [paymentslistAccounts](#paymentslistaccounts) - List accounts
* [readConnectorConfig](#readconnectorconfig) - Read the config of a connector
* [resetConnector](#resetconnector) - Reset a connector
* [uninstallConnector](#uninstallconnector) - Uninstall a connector
* [updateMetadata](#updatemetadata) - Update metadata

## connectorsStripeTransfer

Execute a transfer between two Stripe accounts.

### Example Usage

```php
<?php

declare(strict_types=1);
require_once 'vendor/autoload.php';

use \formance\stack\SDK;
use \formance\stack\Models\Shared\Security;
use \formance\stack\Models\Shared\StripeTransferRequest;

$sdk = SDK::builder()
    ->build();

try {
    $request = new StripeTransferRequest();
    $request->amount = 100;
    $request->asset = 'USD';
    $request->destination = 'acct_1Gqj58KZcSIg2N2q';
    $request->metadata = [
        'corporis' => 'dolore',
    ];

    $response = $sdk->payments->connectorsStripeTransfer($request);

    if ($response->stripeTransferResponse !== null) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

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
    $request->connector = Connector::MODULR;

    $response = $sdk->payments->connectorsTransfer($request);

    if ($response->transferResponse !== null) {
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
    $request->accountId = 'dicta';
    $request->asset = 'harum';
    $request->cursor = 'aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ==';
    $request->from = DateTime::createFromFormat('Y-m-d\TH:i:sP', '2022-02-13T15:01:52.114Z');
    $request->limit = 414263;
    $request->pageSize = 918236;
    $request->sort = [
        'ipsum',
    ];
    $request->to = DateTime::createFromFormat('Y-m-d\TH:i:sP', '2021-11-14T09:53:27.431Z');

    $response = $sdk->payments->getAccountBalances($request);

    if ($response->balancesCursor !== null) {
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
    $request->connector = Connector::CURRENCY_CLOUD;
    $request->taskId = 'pariatur';

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
    $request->paymentId = 'modi';

    $response = $sdk->payments->getPayment($request);

    if ($response->paymentResponse !== null) {
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
    $request->requestBody = new CurrencyCloudConfig();
    $request->requestBody->apiKey = 'XXX';
    $request->requestBody->endpoint = 'XXX';
    $request->requestBody->loginID = 'XXX';
    $request->requestBody->pollingPeriod = '60s';
    $request->connector = Connector::CURRENCY_CLOUD;

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
    $request->connector = Connector::MONEYCORP;
    $request->cursor = 'aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ==';
    $request->pageSize = 93940;

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
    $request->connector = Connector::MONEYCORP;

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
    $request->pageSize = 575947;
    $request->sort = [
        'itaque',
    ];

    $response = $sdk->payments->listPayments($request);

    if ($response->paymentsCursor !== null) {
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
    $request->accountId = 'incidunt';

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
    $request->pageSize = 318569;
    $request->sort = [
        'est',
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
    $request->connector = Connector::MANGOPAY;

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
    $request->connector = Connector::DUMMY_PAY;

    $response = $sdk->payments->resetConnector($request);

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
    $request->connector = Connector::BANKING_CIRCLE;

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
    $request->paymentMetadata->key = 'distinctio';
    $request->paymentId = 'quibusdam';

    $response = $sdk->payments->updateMetadata($request);

    if ($response->statusCode === 200) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```
