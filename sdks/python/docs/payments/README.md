# payments

### Available Operations

* [connectors_transfer](#connectors_transfer) - Transfer funds between Connector accounts
* [create_bank_account](#create_bank_account) - Create a BankAccount in Payments and on the PSP
* [create_transfer_initiation](#create_transfer_initiation) - Create a TransferInitiation
* [delete_transfer_initiation](#delete_transfer_initiation) - Delete a transfer initiation
* [get_account_balances](#get_account_balances) - Get account balances
* [get_bank_account](#get_bank_account) - Get a bank account created by user on Formance
* [~~get_connector_task~~](#get_connector_task) - Read a specific task of the connector :warning: **Deprecated**
* [get_connector_task_v1](#get_connector_task_v1) - Read a specific task of the connector
* [get_payment](#get_payment) - Get a payment
* [get_transfer_initiation](#get_transfer_initiation) - Get a transfer initiation
* [install_connector](#install_connector) - Install a connector
* [list_all_connectors](#list_all_connectors) - List all installed connectors
* [list_bank_accounts](#list_bank_accounts) - List bank accounts created by user on Formance
* [list_configs_available_connectors](#list_configs_available_connectors) - List the configs of each available connector
* [~~list_connector_tasks~~](#list_connector_tasks) - List tasks from a connector :warning: **Deprecated**
* [list_connector_tasks_v1](#list_connector_tasks_v1) - List tasks from a connector
* [list_payments](#list_payments) - List payments
* [list_transfer_initiations](#list_transfer_initiations) - List Transfer Initiations
* [paymentsget_account](#paymentsget_account) - Get an account
* [paymentsget_server_info](#paymentsget_server_info) - Get server info
* [paymentslist_accounts](#paymentslist_accounts) - List accounts
* [~~read_connector_config~~](#read_connector_config) - Read the config of a connector :warning: **Deprecated**
* [read_connector_config_v1](#read_connector_config_v1) - Read the config of a connector
* [~~reset_connector~~](#reset_connector) - Reset a connector :warning: **Deprecated**
* [reset_connector_v1](#reset_connector_v1) - Reset a connector
* [retry_transfer_initiation](#retry_transfer_initiation) - Retry a failed transfer initiation
* [udpate_transfer_initiation_status](#udpate_transfer_initiation_status) - Update the status of a transfer initiation
* [~~uninstall_connector~~](#uninstall_connector) - Uninstall a connector :warning: **Deprecated**
* [uninstall_connector_v1](#uninstall_connector_v1) - Uninstall a connector
* [update_metadata](#update_metadata) - Update metadata

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
    connector=shared.Connector.CURRENCY_CLOUD,
)

res = s.payments.connectors_transfer(req)

if res.transfer_response is not None:
    # handle response
```

## create_bank_account

Create a bank account in Payments and on the PSP.

### Example Usage

```python
import sdk
from sdk.models import shared

s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)

req = shared.BankAccountRequest(
    account_number='voluptates',
    connector_id='quasi',
    country='GB',
    iban='repudiandae',
    name='My account',
    swift_bic_code='sint',
)

res = s.payments.create_bank_account(req)

if res.bank_account_response is not None:
    # handle response
```

## create_transfer_initiation

Create a transfer initiation

### Example Usage

```python
import sdk
import dateutil.parser
from sdk.models import shared

s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)

req = shared.TransferInitiationRequest(
    amount=83112,
    asset='USD',
    connector_id='itaque',
    description='incidunt',
    destination_account_id='enim',
    provider=shared.Connector.STRIPE,
    reference='XXX',
    scheduled_at=dateutil.parser.isoparse('2021-04-26T02:10:00.226Z'),
    source_account_id='explicabo',
    type=shared.TransferInitiationRequestType.PAYOUT,
    validated=False,
)

res = s.payments.create_transfer_initiation(req)

if res.transfer_initiation_response is not None:
    # handle response
```

## delete_transfer_initiation

Delete a transfer initiation by its id.

### Example Usage

```python
import sdk
from sdk.models import operations

s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)

req = operations.DeleteTransferInitiationRequest(
    transfer_id='distinctio',
)

res = s.payments.delete_transfer_initiation(req)

if res.status_code == 200:
    # handle response
```

## get_account_balances

Get account balances

### Example Usage

```python
import sdk
import dateutil.parser
from sdk.models import operations

s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)

req = operations.GetAccountBalancesRequest(
    account_id='quibusdam',
    asset='labore',
    cursor='aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ==',
    from_=dateutil.parser.isoparse('2022-10-26T03:14:36.345Z'),
    limit=397821,
    page_size=586513,
    sort=[
        'perferendis',
        'magni',
        'assumenda',
    ],
    to=dateutil.parser.isoparse('2022-12-30T06:52:02.282Z'),
)

res = s.payments.get_account_balances(req)

if res.balances_cursor is not None:
    # handle response
```

## get_bank_account

Get a bank account created by user on Formance

### Example Usage

```python
import sdk
from sdk.models import operations

s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)

req = operations.GetBankAccountRequest(
    bank_account_id='fugit',
)

res = s.payments.get_bank_account(req)

if res.bank_account_response is not None:
    # handle response
```

## ~~get_connector_task~~

Get a specific task associated to the connector.

> :warning: **DEPRECATED**: this method will be removed in a future release, please migrate away from it as soon as possible.

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
    task_id='excepturi',
)

res = s.payments.get_connector_task(req)

if res.task_response is not None:
    # handle response
```

## get_connector_task_v1

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

req = operations.GetConnectorTaskV1Request(
    connector=shared.Connector.WISE,
    connector_id='facilis',
    task_id='tempore',
)

res = s.payments.get_connector_task_v1(req)

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
    payment_id='labore',
)

res = s.payments.get_payment(req)

if res.payment_response is not None:
    # handle response
```

## get_transfer_initiation

Get a transfer initiation

### Example Usage

```python
import sdk
from sdk.models import operations

s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)

req = operations.GetTransferInitiationRequest(
    transfer_id='delectus',
)

res = s.payments.get_transfer_initiation(req)

if res.transfer_initiation_response is not None:
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
    request_body=shared.ModulrConfig(
        api_key='XXX',
        api_secret='XXX',
        endpoint='XXX',
        name='My Modulr Account',
        polling_period='60s',
    ),
    connector=shared.Connector.DUMMY_PAY,
)

res = s.payments.install_connector(req)

if res.connector_response is not None:
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

## list_bank_accounts

List all bank accounts created by user on Formance.

### Example Usage

```python
import sdk
from sdk.models import operations

s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)

req = operations.ListBankAccountsRequest(
    cursor='aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ==',
    page_size=756107,
    sort=[
        'aliquid',
        'provident',
        'necessitatibus',
    ],
)

res = s.payments.list_bank_accounts(req)

if res.bank_accounts_cursor is not None:
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

## ~~list_connector_tasks~~

List all tasks associated with this connector.

> :warning: **DEPRECATED**: this method will be removed in a future release, please migrate away from it as soon as possible.

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
    connector=shared.Connector.CURRENCY_CLOUD,
    cursor='aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ==',
    page_size=638921,
)

res = s.payments.list_connector_tasks(req)

if res.tasks_cursor is not None:
    # handle response
```

## list_connector_tasks_v1

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

req = operations.ListConnectorTasksV1Request(
    connector=shared.Connector.DUMMY_PAY,
    connector_id='debitis',
    cursor='aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ==',
    page_size=952749,
)

res = s.payments.list_connector_tasks_v1(req)

if res.tasks_cursor is not None:
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
    page_size=680056,
    sort=[
        'in',
        'illum',
    ],
)

res = s.payments.list_payments(req)

if res.payments_cursor is not None:
    # handle response
```

## list_transfer_initiations

List Transfer Initiations

### Example Usage

```python
import sdk
from sdk.models import operations

s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)

req = operations.ListTransferInitiationsRequest(
    cursor='aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ==',
    page_size=978571,
    query='rerum',
    sort=[
        'magnam',
    ],
)

res = s.payments.list_transfer_initiations(req)

if res.transfer_initiations_cursor is not None:
    # handle response
```

## paymentsget_account

Get an account

### Example Usage

```python
import sdk
from sdk.models import operations

s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)

req = operations.PaymentsgetAccountRequest(
    account_id='cumque',
)

res = s.payments.paymentsget_account(req)

if res.payments_account_response is not None:
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
    page_size=813798,
    sort=[
        'aliquid',
        'laborum',
    ],
)

res = s.payments.paymentslist_accounts(req)

if res.accounts_cursor is not None:
    # handle response
```

## ~~read_connector_config~~

Read connector config

> :warning: **DEPRECATED**: this method will be removed in a future release, please migrate away from it as soon as possible.

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
    connector=shared.Connector.MONEYCORP,
)

res = s.payments.read_connector_config(req)

if res.connector_config_response is not None:
    # handle response
```

## read_connector_config_v1

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

req = operations.ReadConnectorConfigV1Request(
    connector=shared.Connector.DUMMY_PAY,
    connector_id='occaecati',
)

res = s.payments.read_connector_config_v1(req)

if res.connector_config_response is not None:
    # handle response
```

## ~~reset_connector~~

Reset a connector by its name.
It will remove the connector and ALL PAYMENTS generated with it.


> :warning: **DEPRECATED**: this method will be removed in a future release, please migrate away from it as soon as possible.

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
    connector=shared.Connector.WISE,
)

res = s.payments.reset_connector(req)

if res.status_code == 200:
    # handle response
```

## reset_connector_v1

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

req = operations.ResetConnectorV1Request(
    connector=shared.Connector.MONEYCORP,
    connector_id='delectus',
)

res = s.payments.reset_connector_v1(req)

if res.status_code == 200:
    # handle response
```

## retry_transfer_initiation

Retry a failed transfer initiation

### Example Usage

```python
import sdk
from sdk.models import operations

s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)

req = operations.RetryTransferInitiationRequest(
    transfer_id='quidem',
)

res = s.payments.retry_transfer_initiation(req)

if res.status_code == 200:
    # handle response
```

## udpate_transfer_initiation_status

Update a transfer initiation status

### Example Usage

```python
import sdk
from sdk.models import operations, shared

s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)

req = operations.UdpateTransferInitiationStatusRequest(
    update_transfer_initiation_status_request=shared.UpdateTransferInitiationStatusRequest(
        status=shared.UpdateTransferInitiationStatusRequestStatus.FAILED,
    ),
    transfer_id='nam',
)

res = s.payments.udpate_transfer_initiation_status(req)

if res.status_code == 200:
    # handle response
```

## ~~uninstall_connector~~

Uninstall a connector by its name.

> :warning: **DEPRECATED**: this method will be removed in a future release, please migrate away from it as soon as possible.

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
    connector=shared.Connector.BANKING_CIRCLE,
)

res = s.payments.uninstall_connector(req)

if res.status_code == 200:
    # handle response
```

## uninstall_connector_v1

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

req = operations.UninstallConnectorV1Request(
    connector=shared.Connector.CURRENCY_CLOUD,
    connector_id='deleniti',
)

res = s.payments.uninstall_connector_v1(req)

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
        key='sapiente',
    ),
    payment_id='amet',
)

res = s.payments.update_metadata(req)

if res.status_code == 200:
    # handle response
```
