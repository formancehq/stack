# payments

### Available Operations

* [connectors_stripe_transfer](#connectors_stripe_transfer) - Transfer funds between Stripe accounts
* [connectors_transfer](#connectors_transfer) - Transfer funds between Connector accounts
* [get_account_balances](#get_account_balances) - Get account balances
* [get_connector_task](#get_connector_task) - Read a specific task of the connector
* [get_payment](#get_payment) - Get a payment
* [install_connector](#install_connector) - Install a connector
* [list_all_connectors](#list_all_connectors) - List all installed connectors
* [list_configs_available_connectors](#list_configs_available_connectors) - List the configs of each available connector
* [list_connector_tasks](#list_connector_tasks) - List tasks from a connector
* [list_connectors_transfers](#list_connectors_transfers) - List transfers and their statuses
* [list_payments](#list_payments) - List payments
* [paymentsget_account](#paymentsget_account) - Get an account
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
        authorization="",
    ),
)

req = shared.StripeTransferRequest(
    amount=100,
    asset='USD',
    destination='acct_1Gqj58KZcSIg2N2q',
    metadata=shared.StripeTransferRequestMetadata(),
)

res = s.payments.connectors_stripe_transfer(req)

if res.stripe_transfer_response is not None:
    # handle response
```

### Parameters

| Parameter                                                                    | Type                                                                         | Required                                                                     | Description                                                                  |
| ---------------------------------------------------------------------------- | ---------------------------------------------------------------------------- | ---------------------------------------------------------------------------- | ---------------------------------------------------------------------------- |
| `request`                                                                    | [shared.StripeTransferRequest](../../models/shared/stripetransferrequest.md) | :heavy_check_mark:                                                           | The request object to use for the request.                                   |


### Response

**[operations.ConnectorsStripeTransferResponse](../../models/operations/connectorsstripetransferresponse.md)**


## connectors_transfer

Execute a transfer between two accounts.

### Example Usage

```python
import sdk
from sdk.models import operations, shared

s = sdk.SDK(
    security=shared.Security(
        authorization="",
    ),
)

req = operations.ConnectorsTransferRequest(
    transfer_request=shared.TransferRequest(
        amount=100,
        asset='USD',
        destination='acct_1Gqj58KZcSIg2N2q',
        source='acct_1Gqj58KZcSIg2N2q',
    ),
    connector=shared.Connector.STRIPE,
)

res = s.payments.connectors_transfer(req)

if res.transfer_response is not None:
    # handle response
```

### Parameters

| Parameter                                                                                    | Type                                                                                         | Required                                                                                     | Description                                                                                  |
| -------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------- |
| `request`                                                                                    | [operations.ConnectorsTransferRequest](../../models/operations/connectorstransferrequest.md) | :heavy_check_mark:                                                                           | The request object to use for the request.                                                   |


### Response

**[operations.ConnectorsTransferResponse](../../models/operations/connectorstransferresponse.md)**


## get_account_balances

Get account balances

### Example Usage

```python
import sdk
import dateutil.parser
from sdk.models import operations, shared

s = sdk.SDK(
    security=shared.Security(
        authorization="",
    ),
)

req = operations.GetAccountBalancesRequest(
    account_id='quasi',
    asset='error',
    cursor='aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ==',
    from_=dateutil.parser.isoparse('2020-12-24T08:13:29.299Z'),
    limit=96098,
    page_size=971945,
    sort=[
        'vero',
        'nihil',
        'praesentium',
        'voluptatibus',
    ],
    to=dateutil.parser.isoparse('2022-05-25T05:33:11.349Z'),
)

res = s.payments.get_account_balances(req)

if res.balances_cursor is not None:
    # handle response
```

### Parameters

| Parameter                                                                                    | Type                                                                                         | Required                                                                                     | Description                                                                                  |
| -------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------- |
| `request`                                                                                    | [operations.GetAccountBalancesRequest](../../models/operations/getaccountbalancesrequest.md) | :heavy_check_mark:                                                                           | The request object to use for the request.                                                   |


### Response

**[operations.GetAccountBalancesResponse](../../models/operations/getaccountbalancesresponse.md)**


## get_connector_task

Get a specific task associated to the connector.

### Example Usage

```python
import sdk
from sdk.models import operations, shared

s = sdk.SDK(
    security=shared.Security(
        authorization="",
    ),
)

req = operations.GetConnectorTaskRequest(
    connector=shared.Connector.MODULR,
    task_id='cum',
)

res = s.payments.get_connector_task(req)

if res.task_response is not None:
    # handle response
```

### Parameters

| Parameter                                                                                | Type                                                                                     | Required                                                                                 | Description                                                                              |
| ---------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------- |
| `request`                                                                                | [operations.GetConnectorTaskRequest](../../models/operations/getconnectortaskrequest.md) | :heavy_check_mark:                                                                       | The request object to use for the request.                                               |


### Response

**[operations.GetConnectorTaskResponse](../../models/operations/getconnectortaskresponse.md)**


## get_payment

Get a payment

### Example Usage

```python
import sdk
from sdk.models import operations, shared

s = sdk.SDK(
    security=shared.Security(
        authorization="",
    ),
)

req = operations.GetPaymentRequest(
    payment_id='perferendis',
)

res = s.payments.get_payment(req)

if res.payment_response is not None:
    # handle response
```

### Parameters

| Parameter                                                                    | Type                                                                         | Required                                                                     | Description                                                                  |
| ---------------------------------------------------------------------------- | ---------------------------------------------------------------------------- | ---------------------------------------------------------------------------- | ---------------------------------------------------------------------------- |
| `request`                                                                    | [operations.GetPaymentRequest](../../models/operations/getpaymentrequest.md) | :heavy_check_mark:                                                           | The request object to use for the request.                                   |


### Response

**[operations.GetPaymentResponse](../../models/operations/getpaymentresponse.md)**


## install_connector

Install a connector by its name and config.

### Example Usage

```python
import sdk
from sdk.models import operations, shared

s = sdk.SDK(
    security=shared.Security(
        authorization="",
    ),
)

req = operations.InstallConnectorRequest(
    request_body=shared.StripeConfig(
        api_key='XXX',
        page_size=50,
        polling_period='60s',
    ),
    connector=shared.Connector.MODULR,
)

res = s.payments.install_connector(req)

if res.status_code == 200:
    # handle response
```

### Parameters

| Parameter                                                                                | Type                                                                                     | Required                                                                                 | Description                                                                              |
| ---------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------- |
| `request`                                                                                | [operations.InstallConnectorRequest](../../models/operations/installconnectorrequest.md) | :heavy_check_mark:                                                                       | The request object to use for the request.                                               |


### Response

**[operations.InstallConnectorResponse](../../models/operations/installconnectorresponse.md)**


## list_all_connectors

List all installed connectors.

### Example Usage

```python
import sdk
from sdk.models import shared

s = sdk.SDK(
    security=shared.Security(
        authorization="",
    ),
)


res = s.payments.list_all_connectors()

if res.connectors_response is not None:
    # handle response
```


### Response

**[operations.ListAllConnectorsResponse](../../models/operations/listallconnectorsresponse.md)**


## list_configs_available_connectors

List the configs of each available connector.

### Example Usage

```python
import sdk
from sdk.models import shared

s = sdk.SDK(
    security=shared.Security(
        authorization="",
    ),
)


res = s.payments.list_configs_available_connectors()

if res.connectors_configs_response is not None:
    # handle response
```


### Response

**[operations.ListConfigsAvailableConnectorsResponse](../../models/operations/listconfigsavailableconnectorsresponse.md)**


## list_connector_tasks

List all tasks associated with this connector.

### Example Usage

```python
import sdk
from sdk.models import operations, shared

s = sdk.SDK(
    security=shared.Security(
        authorization="",
    ),
)

req = operations.ListConnectorTasksRequest(
    connector=shared.Connector.WISE,
    cursor='aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ==',
    page_size=979587,
)

res = s.payments.list_connector_tasks(req)

if res.tasks_cursor is not None:
    # handle response
```

### Parameters

| Parameter                                                                                    | Type                                                                                         | Required                                                                                     | Description                                                                                  |
| -------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------- |
| `request`                                                                                    | [operations.ListConnectorTasksRequest](../../models/operations/listconnectortasksrequest.md) | :heavy_check_mark:                                                                           | The request object to use for the request.                                                   |


### Response

**[operations.ListConnectorTasksResponse](../../models/operations/listconnectortasksresponse.md)**


## list_connectors_transfers

List transfers

### Example Usage

```python
import sdk
from sdk.models import operations, shared

s = sdk.SDK(
    security=shared.Security(
        authorization="",
    ),
)

req = operations.ListConnectorsTransfersRequest(
    connector=shared.Connector.STRIPE,
)

res = s.payments.list_connectors_transfers(req)

if res.transfers_response is not None:
    # handle response
```

### Parameters

| Parameter                                                                                              | Type                                                                                                   | Required                                                                                               | Description                                                                                            |
| ------------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------------ |
| `request`                                                                                              | [operations.ListConnectorsTransfersRequest](../../models/operations/listconnectorstransfersrequest.md) | :heavy_check_mark:                                                                                     | The request object to use for the request.                                                             |


### Response

**[operations.ListConnectorsTransfersResponse](../../models/operations/listconnectorstransfersresponse.md)**


## list_payments

List payments

### Example Usage

```python
import sdk
from sdk.models import operations, shared

s = sdk.SDK(
    security=shared.Security(
        authorization="",
    ),
)

req = operations.ListPaymentsRequest(
    cursor='aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ==',
    page_size=359444,
    sort=[
        'iusto',
        'dicta',
    ],
)

res = s.payments.list_payments(req)

if res.payments_cursor is not None:
    # handle response
```

### Parameters

| Parameter                                                                        | Type                                                                             | Required                                                                         | Description                                                                      |
| -------------------------------------------------------------------------------- | -------------------------------------------------------------------------------- | -------------------------------------------------------------------------------- | -------------------------------------------------------------------------------- |
| `request`                                                                        | [operations.ListPaymentsRequest](../../models/operations/listpaymentsrequest.md) | :heavy_check_mark:                                                               | The request object to use for the request.                                       |


### Response

**[operations.ListPaymentsResponse](../../models/operations/listpaymentsresponse.md)**


## paymentsget_account

Get an account

### Example Usage

```python
import sdk
from sdk.models import operations, shared

s = sdk.SDK(
    security=shared.Security(
        authorization="",
    ),
)

req = operations.PaymentsgetAccountRequest(
    account_id='harum',
)

res = s.payments.paymentsget_account(req)

if res.payments_account_response is not None:
    # handle response
```

### Parameters

| Parameter                                                                                    | Type                                                                                         | Required                                                                                     | Description                                                                                  |
| -------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------- |
| `request`                                                                                    | [operations.PaymentsgetAccountRequest](../../models/operations/paymentsgetaccountrequest.md) | :heavy_check_mark:                                                                           | The request object to use for the request.                                                   |


### Response

**[operations.PaymentsgetAccountResponse](../../models/operations/paymentsgetaccountresponse.md)**


## paymentsget_server_info

Get server info

### Example Usage

```python
import sdk
from sdk.models import shared

s = sdk.SDK(
    security=shared.Security(
        authorization="",
    ),
)


res = s.payments.paymentsget_server_info()

if res.server_info is not None:
    # handle response
```


### Response

**[operations.PaymentsgetServerInfoResponse](../../models/operations/paymentsgetserverinforesponse.md)**


## paymentslist_accounts

List accounts

### Example Usage

```python
import sdk
from sdk.models import operations, shared

s = sdk.SDK(
    security=shared.Security(
        authorization="",
    ),
)

req = operations.PaymentslistAccountsRequest(
    cursor='aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ==',
    page_size=317983,
    sort=[
        'commodi',
        'repudiandae',
        'quae',
        'ipsum',
    ],
)

res = s.payments.paymentslist_accounts(req)

if res.accounts_cursor is not None:
    # handle response
```

### Parameters

| Parameter                                                                                        | Type                                                                                             | Required                                                                                         | Description                                                                                      |
| ------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------ |
| `request`                                                                                        | [operations.PaymentslistAccountsRequest](../../models/operations/paymentslistaccountsrequest.md) | :heavy_check_mark:                                                                               | The request object to use for the request.                                                       |


### Response

**[operations.PaymentslistAccountsResponse](../../models/operations/paymentslistaccountsresponse.md)**


## read_connector_config

Read connector config

### Example Usage

```python
import sdk
from sdk.models import operations, shared

s = sdk.SDK(
    security=shared.Security(
        authorization="",
    ),
)

req = operations.ReadConnectorConfigRequest(
    connector=shared.Connector.BANKING_CIRCLE,
)

res = s.payments.read_connector_config(req)

if res.connector_config_response is not None:
    # handle response
```

### Parameters

| Parameter                                                                                      | Type                                                                                           | Required                                                                                       | Description                                                                                    |
| ---------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------- |
| `request`                                                                                      | [operations.ReadConnectorConfigRequest](../../models/operations/readconnectorconfigrequest.md) | :heavy_check_mark:                                                                             | The request object to use for the request.                                                     |


### Response

**[operations.ReadConnectorConfigResponse](../../models/operations/readconnectorconfigresponse.md)**


## reset_connector

Reset a connector by its name.
It will remove the connector and ALL PAYMENTS generated with it.


### Example Usage

```python
import sdk
from sdk.models import operations, shared

s = sdk.SDK(
    security=shared.Security(
        authorization="",
    ),
)

req = operations.ResetConnectorRequest(
    connector=shared.Connector.CURRENCY_CLOUD,
)

res = s.payments.reset_connector(req)

if res.status_code == 200:
    # handle response
```

### Parameters

| Parameter                                                                            | Type                                                                                 | Required                                                                             | Description                                                                          |
| ------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------ |
| `request`                                                                            | [operations.ResetConnectorRequest](../../models/operations/resetconnectorrequest.md) | :heavy_check_mark:                                                                   | The request object to use for the request.                                           |


### Response

**[operations.ResetConnectorResponse](../../models/operations/resetconnectorresponse.md)**


## uninstall_connector

Uninstall a connector by its name.

### Example Usage

```python
import sdk
from sdk.models import operations, shared

s = sdk.SDK(
    security=shared.Security(
        authorization="",
    ),
)

req = operations.UninstallConnectorRequest(
    connector=shared.Connector.CURRENCY_CLOUD,
)

res = s.payments.uninstall_connector(req)

if res.status_code == 200:
    # handle response
```

### Parameters

| Parameter                                                                                    | Type                                                                                         | Required                                                                                     | Description                                                                                  |
| -------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------- |
| `request`                                                                                    | [operations.UninstallConnectorRequest](../../models/operations/uninstallconnectorrequest.md) | :heavy_check_mark:                                                                           | The request object to use for the request.                                                   |


### Response

**[operations.UninstallConnectorResponse](../../models/operations/uninstallconnectorresponse.md)**


## update_metadata

Update metadata

### Example Usage

```python
import sdk
from sdk.models import operations, shared

s = sdk.SDK(
    security=shared.Security(
        authorization="",
    ),
)

req = operations.UpdateMetadataRequest(
    payment_metadata=shared.PaymentMetadata(
        key='pariatur',
    ),
    payment_id='modi',
)

res = s.payments.update_metadata(req)

if res.status_code == 200:
    # handle response
```

### Parameters

| Parameter                                                                            | Type                                                                                 | Required                                                                             | Description                                                                          |
| ------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------ |
| `request`                                                                            | [operations.UpdateMetadataRequest](../../models/operations/updatemetadatarequest.md) | :heavy_check_mark:                                                                   | The request object to use for the request.                                           |


### Response

**[operations.UpdateMetadataResponse](../../models/operations/updatemetadataresponse.md)**

