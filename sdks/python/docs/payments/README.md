# payments

### Available Operations

* [connectors_stripe_transfer](#connectors_stripe_transfer) - Transfer funds between Stripe accounts
* [connectors_transfer](#connectors_transfer) - Transfer funds between Connector accounts
* [get_connector_task](#get_connector_task) - Read a specific task of the connector
* [get_payment](#get_payment) - Get a payment
* [install_connector](#install_connector) - Install a connector
* [list_all_connectors](#list_all_connectors) - List all installed connectors
* [list_configs_available_connectors](#list_configs_available_connectors) - List the configs of each available connector
* [list_connector_tasks](#list_connector_tasks) - List tasks from a connector
* [list_connectors_transfers](#list_connectors_transfers) - List transfers and their statuses
* [list_payments](#list_payments) - List payments
* [paymentsget_server_info](#paymentsget_server_info) - Get server info
* [paymentslist_accounts](#paymentslist_accounts) - List accounts
* [read_connector_config](#read_connector_config) - Read the config of a connector
* [reset_connector](#reset_connector) - Reset a connector
* [uninstall_connector](#uninstall_connector) - Uninstall a connector
* [update_metadata](#update_metadata) - Update metadata

## connectors_stripe_transfer

Execute a transfer between two Stripe accounts.

### Example Usage

```python
import sdk
from sdk.models import shared

s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)

req = shared.StripeTransferRequest(
    amount=100,
    asset='USD',
    destination='acct_1Gqj58KZcSIg2N2q',
    metadata={
        "non": 'eligendi',
        "sint": 'aliquid',
    },
)

res = s.payments.connectors_stripe_transfer(req)

if res.stripe_transfer_response is not None:
    # handle response
```

## connectors_transfer

Execute a transfer between two accounts.

### Example Usage

```python
import sdk
from sdk.models import operations, shared

s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)

req = operations.ConnectorsTransferRequest(
    transfer_request=shared.TransferRequest(
        amount=100,
        asset='USD',
        destination='acct_1Gqj58KZcSIg2N2q',
        source='acct_1Gqj58KZcSIg2N2q',
    ),
    connector=shared.Connector.MODULR,
)

res = s.payments.connectors_transfer(req)

if res.transfer_response is not None:
    # handle response
```

## get_connector_task

Get a specific task associated to the connector.

### Example Usage

```python
import sdk
from sdk.models import operations, shared

s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)

req = operations.GetConnectorTaskRequest(
    connector=shared.Connector.BANKING_CIRCLE,
    task_id='sint',
)

res = s.payments.get_connector_task(req)

if res.task_response is not None:
    # handle response
```

## get_payment

Get a payment

### Example Usage

```python
import sdk
from sdk.models import operations

s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)

req = operations.GetPaymentRequest(
    payment_id='officia',
)

res = s.payments.get_payment(req)

if res.payment_response is not None:
    # handle response
```

## install_connector

Install a connector by its name and config.

### Example Usage

```python
import sdk
from sdk.models import operations, shared

s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)

req = operations.InstallConnectorRequest(
    request_body=shared.DummyPayConfig(
        directory='/tmp/dummypay',
        file_generation_period='60s',
        file_polling_period='60s',
    ),
    connector=shared.Connector.BANKING_CIRCLE,
)

res = s.payments.install_connector(req)

if res.status_code == 200:
    # handle response
```

## list_all_connectors

List all installed connectors.

### Example Usage

```python
import sdk


s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)


res = s.payments.list_all_connectors()

if res.connectors_response is not None:
    # handle response
```

## list_configs_available_connectors

List the configs of each available connector.

### Example Usage

```python
import sdk


s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)


res = s.payments.list_configs_available_connectors()

if res.connectors_configs_response is not None:
    # handle response
```

## list_connector_tasks

List all tasks associated with this connector.

### Example Usage

```python
import sdk
from sdk.models import operations, shared

s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)

req = operations.ListConnectorTasksRequest(
    connector=shared.Connector.BANKING_CIRCLE,
    cursor='aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ==',
    page_size=680056,
)

res = s.payments.list_connector_tasks(req)

if res.tasks_cursor is not None:
    # handle response
```

## list_connectors_transfers

List transfers

### Example Usage

```python
import sdk
from sdk.models import operations, shared

s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)

req = operations.ListConnectorsTransfersRequest(
    connector=shared.Connector.WISE,
)

res = s.payments.list_connectors_transfers(req)

if res.transfers_response is not None:
    # handle response
```

## list_payments

List payments

### Example Usage

```python
import sdk
from sdk.models import operations

s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)

req = operations.ListPaymentsRequest(
    cursor='aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ==',
    page_size=449198,
    sort=[
        'maiores',
        'rerum',
        'dicta',
        'magnam',
    ],
)

res = s.payments.list_payments(req)

if res.payments_cursor is not None:
    # handle response
```

## paymentsget_server_info

Get server info

### Example Usage

```python
import sdk


s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)


res = s.payments.paymentsget_server_info()

if res.server_info is not None:
    # handle response
```

## paymentslist_accounts

List accounts

### Example Usage

```python
import sdk
from sdk.models import operations

s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)

req = operations.PaymentslistAccountsRequest(
    cursor='aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ==',
    page_size=767024,
    sort=[
        'ea',
        'aliquid',
        'laborum',
        'accusamus',
    ],
)

res = s.payments.paymentslist_accounts(req)

if res.accounts_cursor is not None:
    # handle response
```

## read_connector_config

Read connector config

### Example Usage

```python
import sdk
from sdk.models import operations, shared

s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)

req = operations.ReadConnectorConfigRequest(
    connector=shared.Connector.DUMMY_PAY,
)

res = s.payments.read_connector_config(req)

if res.connector_config_response is not None:
    # handle response
```

## reset_connector

Reset a connector by its name.
It will remove the connector and ALL PAYMENTS generated with it.


### Example Usage

```python
import sdk
from sdk.models import operations, shared

s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)

req = operations.ResetConnectorRequest(
    connector=shared.Connector.MODULR,
)

res = s.payments.reset_connector(req)

if res.status_code == 200:
    # handle response
```

## uninstall_connector

Uninstall a connector by its name.

### Example Usage

```python
import sdk
from sdk.models import operations, shared

s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)

req = operations.UninstallConnectorRequest(
    connector=shared.Connector.DUMMY_PAY,
)

res = s.payments.uninstall_connector(req)

if res.status_code == 200:
    # handle response
```

## update_metadata

Update metadata

### Example Usage

```python
import sdk
from sdk.models import operations, shared

s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)

req = operations.UpdateMetadataRequest(
    payment_metadata=shared.PaymentMetadata(
        key='accusamus',
    ),
    payment_id='delectus',
)

res = s.payments.update_metadata(req)

if res.status_code == 200:
    # handle response
```
