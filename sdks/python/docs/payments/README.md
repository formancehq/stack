# payments

### Available Operations

* [connectors_transfer](#connectors_transfer) - Transfer funds between Connector accounts
* [create_bank_account](#create_bank_account) - Create a BankAccount in Payments and on the PSP
* [create_transfer_initiation](#create_transfer_initiation) - Create a TransferInitiation
* [delete_transfer_initiation](#delete_transfer_initiation) - Delete a transfer initiation
* [get_account_balances](#get_account_balances) - Get account balances
* [get_bank_account](#get_bank_account) - Get a bank account created by user on Formance
* [get_connector_task](#get_connector_task) - Read a specific task of the connector
* [get_payment](#get_payment) - Get a payment
* [get_transfer_initiation](#get_transfer_initiation) - Get a transfer initiation
* [install_connector](#install_connector) - Install a connector
* [list_all_connectors](#list_all_connectors) - List all installed connectors
* [list_bank_accounts](#list_bank_accounts) - List bank accounts created by user on Formance
* [list_configs_available_connectors](#list_configs_available_connectors) - List the configs of each available connector
* [list_connector_tasks](#list_connector_tasks) - List tasks from a connector
* [list_connectors_transfers](#list_connectors_transfers) - List transfers and their statuses
* [list_payments](#list_payments) - List payments
* [list_transfer_initiations](#list_transfer_initiations) - List Transfer Initiations
* [paymentsget_account](#paymentsget_account) - Get an account
* [paymentsget_server_info](#paymentsget_server_info) - Get server info
* [paymentslist_accounts](#paymentslist_accounts) - List accounts
* [read_connector_config](#read_connector_config) - Read the config of a connector
* [reset_connector](#reset_connector) - Reset a connector
* [retry_transfer_initiation](#retry_transfer_initiation) - Retry a failed transfer initiation
* [udpate_transfer_initiation_status](#udpate_transfer_initiation_status) - Update the status of a transfer initiation
* [uninstall_connector](#uninstall_connector) - Uninstall a connector
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
    country='GB',
    iban='quasi',
    name='My account',
    provider=shared.Connector.MONEYCORP,
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
    created_at=dateutil.parser.isoparse('2022-03-02T21:33:21.372Z'),
    description='enim',
    destination_account_id='consequatur',
    provider=shared.Connector.BANKING_CIRCLE,
    reference='XXX',
    source_account_id='quibusdam',
    type=shared.TransferInitiationRequestType.TRANSFER,
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
    transfer_id='deserunt',
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
    account_id='distinctio',
    asset='quibusdam',
    cursor='aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ==',
    from_=dateutil.parser.isoparse('2022-09-26T08:57:48.803Z'),
    limit=183191,
    page_size=397821,
    sort=[
        'quos',
        'perferendis',
        'magni',
    ],
    to=dateutil.parser.isoparse('2021-11-22T01:26:35.048Z'),
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
    bank_account_id='alias',
)

res = s.payments.get_bank_account(req)

if res.bank_account_response is not None:
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
    connector=shared.Connector.DUMMY_PAY,
    task_id='dolorum',
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
    payment_id='excepturi',
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
    transfer_id='tempora',
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
    request_body=shared.BankingCircleConfig(
        authorization_endpoint='XXX',
        endpoint='XXX',
        password='XXX',
        polling_period='60s',
        user_certificate='XXX',
        user_certificate_key='XXX',
        username='XXX',
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
    page_size=288476,
    sort=[
        'eum',
        'non',
        'eligendi',
        'sint',
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
    connector=shared.Connector.MODULR,
    cursor='aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ==',
    page_size=592042,
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
    connector=shared.Connector.MONEYCORP,
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
    page_size=572252,
    sort=[
        'dolor',
        'debitis',
        'a',
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
    page_size=680056,
    query='in',
    sort=[
        'illum',
        'maiores',
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
    account_id='rerum',
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
    page_size=116202,
    sort=[
        'cumque',
        'facere',
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
    connector=shared.Connector.MODULR,
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
    transfer_id='laborum',
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
        status=shared.UpdateTransferInitiationStatusRequestStatus.VALIDATED,
    ),
    transfer_id='non',
)

res = s.payments.udpate_transfer_initiation_status(req)

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
    connector=shared.Connector.CURRENCY_CLOUD,
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
        key='enim',
    ),
    payment_id='accusamus',
)

res = s.payments.update_metadata(req)

if res.status_code == 200:
    # handle response
```
