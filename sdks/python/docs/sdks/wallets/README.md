# wallets

### Available Operations

* [confirm_hold](#confirm_hold) - Confirm a hold
* [create_balance](#create_balance) - Create a balance
* [create_wallet](#create_wallet) - Create a new wallet
* [credit_wallet](#credit_wallet) - Credit a wallet
* [debit_wallet](#debit_wallet) - Debit a wallet
* [get_balance](#get_balance) - Get detailed balance
* [get_hold](#get_hold) - Get a hold
* [get_holds](#get_holds) - Get all holds for a wallet
* [get_transactions](#get_transactions)
* [get_wallet](#get_wallet) - Get a wallet
* [get_wallet_summary](#get_wallet_summary) - Get wallet summary
* [list_balances](#list_balances) - List balances of a wallet
* [list_wallets](#list_wallets) - List all wallets
* [update_wallet](#update_wallet) - Update a wallet
* [void_hold](#void_hold) - Cancel a hold
* [walletsget_server_info](#walletsget_server_info) - Get server info

## confirm_hold

Confirm a hold

### Example Usage

```python
import sdk
from sdk.models import operations, shared

s = sdk.SDK(
    security=shared.Security(
        authorization="",
    ),
)

req = operations.ConfirmHoldRequest(
    confirm_hold_request=shared.ConfirmHoldRequest(
        amount=100,
        final=True,
    ),
    hold_id='nam',
)

res = s.wallets.confirm_hold(req)

if res.status_code == 200:
    # handle response
```

### Parameters

| Parameter                                                                      | Type                                                                           | Required                                                                       | Description                                                                    |
| ------------------------------------------------------------------------------ | ------------------------------------------------------------------------------ | ------------------------------------------------------------------------------ | ------------------------------------------------------------------------------ |
| `request`                                                                      | [operations.ConfirmHoldRequest](../../models/operations/confirmholdrequest.md) | :heavy_check_mark:                                                             | The request object to use for the request.                                     |


### Response

**[operations.ConfirmHoldResponse](../../models/operations/confirmholdresponse.md)**


## create_balance

Create a balance

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

req = operations.CreateBalanceRequest(
    create_balance_request=shared.CreateBalanceRequest(
        expires_at=dateutil.parser.isoparse('2021-12-31T00:47:48.012Z'),
        name='Timmy Feeney',
        priority=423855,
    ),
    id='997074ba-4469-4b6e-a141-959890afa563',
)

res = s.wallets.create_balance(req)

if res.create_balance_response is not None:
    # handle response
```

### Parameters

| Parameter                                                                          | Type                                                                               | Required                                                                           | Description                                                                        |
| ---------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------- |
| `request`                                                                          | [operations.CreateBalanceRequest](../../models/operations/createbalancerequest.md) | :heavy_check_mark:                                                                 | The request object to use for the request.                                         |


### Response

**[operations.CreateBalanceResponse](../../models/operations/createbalanceresponse.md)**


## create_wallet

Create a new wallet

### Example Usage

```python
import sdk
from sdk.models import shared

s = sdk.SDK(
    security=shared.Security(
        authorization="",
    ),
)

req = shared.CreateWalletRequest(
    metadata={
        "odit": 'nemo',
        "quasi": 'iure',
        "doloribus": 'debitis',
        "eius": 'maxime',
    },
    name='Mr. Andres King',
)

res = s.wallets.create_wallet(req)

if res.create_wallet_response is not None:
    # handle response
```

### Parameters

| Parameter                                                                | Type                                                                     | Required                                                                 | Description                                                              |
| ------------------------------------------------------------------------ | ------------------------------------------------------------------------ | ------------------------------------------------------------------------ | ------------------------------------------------------------------------ |
| `request`                                                                | [shared.CreateWalletRequest](../../models/shared/createwalletrequest.md) | :heavy_check_mark:                                                       | The request object to use for the request.                               |


### Response

**[operations.CreateWalletResponse](../../models/operations/createwalletresponse.md)**


## credit_wallet

Credit a wallet

### Example Usage

```python
import sdk
from sdk.models import operations, shared

s = sdk.SDK(
    security=shared.Security(
        authorization="",
    ),
)

req = operations.CreditWalletRequest(
    credit_wallet_request=shared.CreditWalletRequest(
        amount=shared.Monetary(
            amount=352312,
            asset='expedita',
        ),
        balance='nihil',
        metadata={
            "quibusdam": 'sed',
            "saepe": 'pariatur',
            "accusantium": 'consequuntur',
            "praesentium": 'natus',
        },
        reference='magni',
        sources=[
            shared.WalletSubject(
                balance='illum',
                identifier='pariatur',
                type='maxime',
            ),
        ],
    ),
    id='692601fb-576b-40d5-b0d3-0c5fbb258705',
)

res = s.wallets.credit_wallet(req)

if res.status_code == 200:
    # handle response
```

### Parameters

| Parameter                                                                        | Type                                                                             | Required                                                                         | Description                                                                      |
| -------------------------------------------------------------------------------- | -------------------------------------------------------------------------------- | -------------------------------------------------------------------------------- | -------------------------------------------------------------------------------- |
| `request`                                                                        | [operations.CreditWalletRequest](../../models/operations/creditwalletrequest.md) | :heavy_check_mark:                                                               | The request object to use for the request.                                       |


### Response

**[operations.CreditWalletResponse](../../models/operations/creditwalletresponse.md)**


## debit_wallet

Debit a wallet

### Example Usage

```python
import sdk
from sdk.models import operations, shared

s = sdk.SDK(
    security=shared.Security(
        authorization="",
    ),
)

req = operations.DebitWalletRequest(
    debit_wallet_request=shared.DebitWalletRequest(
        amount=shared.Monetary(
            amount=199996,
            asset='eos',
        ),
        balances=[
            'dolores',
        ],
        description='minus',
        destination=shared.LedgerAccountSubject(
            identifier='dolor',
            type='vero',
        ),
        metadata={
            "hic": 'recusandae',
            "omnis": 'facilis',
        },
        pending=False,
    ),
    id='90c28909-b3fe-449a-8d9c-bf48633323f9',
)

res = s.wallets.debit_wallet(req)

if res.debit_wallet_response is not None:
    # handle response
```

### Parameters

| Parameter                                                                      | Type                                                                           | Required                                                                       | Description                                                                    |
| ------------------------------------------------------------------------------ | ------------------------------------------------------------------------------ | ------------------------------------------------------------------------------ | ------------------------------------------------------------------------------ |
| `request`                                                                      | [operations.DebitWalletRequest](../../models/operations/debitwalletrequest.md) | :heavy_check_mark:                                                             | The request object to use for the request.                                     |


### Response

**[operations.DebitWalletResponse](../../models/operations/debitwalletresponse.md)**


## get_balance

Get detailed balance

### Example Usage

```python
import sdk
from sdk.models import operations, shared

s = sdk.SDK(
    security=shared.Security(
        authorization="",
    ),
)

req = operations.GetBalanceRequest(
    balance_name='cum',
    id='77f3a410-0674-4ebf-a928-0d1ba77a89eb',
)

res = s.wallets.get_balance(req)

if res.get_balance_response is not None:
    # handle response
```

### Parameters

| Parameter                                                                    | Type                                                                         | Required                                                                     | Description                                                                  |
| ---------------------------------------------------------------------------- | ---------------------------------------------------------------------------- | ---------------------------------------------------------------------------- | ---------------------------------------------------------------------------- |
| `request`                                                                    | [operations.GetBalanceRequest](../../models/operations/getbalancerequest.md) | :heavy_check_mark:                                                           | The request object to use for the request.                                   |


### Response

**[operations.GetBalanceResponse](../../models/operations/getbalanceresponse.md)**


## get_hold

Get a hold

### Example Usage

```python
import sdk
from sdk.models import operations, shared

s = sdk.SDK(
    security=shared.Security(
        authorization="",
    ),
)

req = operations.GetHoldRequest(
    hold_id='asperiores',
)

res = s.wallets.get_hold(req)

if res.get_hold_response is not None:
    # handle response
```

### Parameters

| Parameter                                                              | Type                                                                   | Required                                                               | Description                                                            |
| ---------------------------------------------------------------------- | ---------------------------------------------------------------------- | ---------------------------------------------------------------------- | ---------------------------------------------------------------------- |
| `request`                                                              | [operations.GetHoldRequest](../../models/operations/getholdrequest.md) | :heavy_check_mark:                                                     | The request object to use for the request.                             |


### Response

**[operations.GetHoldResponse](../../models/operations/getholdresponse.md)**


## get_holds

Get all holds for a wallet

### Example Usage

```python
import sdk
from sdk.models import operations, shared

s = sdk.SDK(
    security=shared.Security(
        authorization="",
    ),
)

req = operations.GetHoldsRequest(
    cursor='aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ==',
    metadata={
        "ipsum": 'voluptate',
        "id": 'saepe',
    },
    page_size=263322,
    wallet_id='aspernatur',
)

res = s.wallets.get_holds(req)

if res.get_holds_response is not None:
    # handle response
```

### Parameters

| Parameter                                                                | Type                                                                     | Required                                                                 | Description                                                              |
| ------------------------------------------------------------------------ | ------------------------------------------------------------------------ | ------------------------------------------------------------------------ | ------------------------------------------------------------------------ |
| `request`                                                                | [operations.GetHoldsRequest](../../models/operations/getholdsrequest.md) | :heavy_check_mark:                                                       | The request object to use for the request.                               |


### Response

**[operations.GetHoldsResponse](../../models/operations/getholdsresponse.md)**


## get_transactions

### Example Usage

```python
import sdk
from sdk.models import operations, shared

s = sdk.SDK(
    security=shared.Security(
        authorization="",
    ),
)

req = operations.GetTransactionsRequest(
    cursor='aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ==',
    page_size=20651,
    wallet_id='amet',
)

res = s.wallets.get_transactions(req)

if res.get_transactions_response is not None:
    # handle response
```

### Parameters

| Parameter                                                                              | Type                                                                                   | Required                                                                               | Description                                                                            |
| -------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------- |
| `request`                                                                              | [operations.GetTransactionsRequest](../../models/operations/gettransactionsrequest.md) | :heavy_check_mark:                                                                     | The request object to use for the request.                                             |


### Response

**[operations.GetTransactionsResponse](../../models/operations/gettransactionsresponse.md)**


## get_wallet

Get a wallet

### Example Usage

```python
import sdk
from sdk.models import operations, shared

s = sdk.SDK(
    security=shared.Security(
        authorization="",
    ),
)

req = operations.GetWalletRequest(
    id='ce5e6a95-d8a0-4d44-ace2-af7a73cf3be4',
)

res = s.wallets.get_wallet(req)

if res.get_wallet_response is not None:
    # handle response
```

### Parameters

| Parameter                                                                  | Type                                                                       | Required                                                                   | Description                                                                |
| -------------------------------------------------------------------------- | -------------------------------------------------------------------------- | -------------------------------------------------------------------------- | -------------------------------------------------------------------------- |
| `request`                                                                  | [operations.GetWalletRequest](../../models/operations/getwalletrequest.md) | :heavy_check_mark:                                                         | The request object to use for the request.                                 |


### Response

**[operations.GetWalletResponse](../../models/operations/getwalletresponse.md)**


## get_wallet_summary

Get wallet summary

### Example Usage

```python
import sdk
from sdk.models import operations, shared

s = sdk.SDK(
    security=shared.Security(
        authorization="",
    ),
)

req = operations.GetWalletSummaryRequest(
    id='53f870b3-26b5-4a73-829c-db1a8422bb67',
)

res = s.wallets.get_wallet_summary(req)

if res.get_wallet_summary_response is not None:
    # handle response
```

### Parameters

| Parameter                                                                                | Type                                                                                     | Required                                                                                 | Description                                                                              |
| ---------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------- |
| `request`                                                                                | [operations.GetWalletSummaryRequest](../../models/operations/getwalletsummaryrequest.md) | :heavy_check_mark:                                                                       | The request object to use for the request.                                               |


### Response

**[operations.GetWalletSummaryResponse](../../models/operations/getwalletsummaryresponse.md)**


## list_balances

List balances of a wallet

### Example Usage

```python
import sdk
from sdk.models import operations, shared

s = sdk.SDK(
    security=shared.Security(
        authorization="",
    ),
)

req = operations.ListBalancesRequest(
    id='9d232271-5bf0-4cbb-9e31-b8b90f3443a1',
)

res = s.wallets.list_balances(req)

if res.list_balances_response is not None:
    # handle response
```

### Parameters

| Parameter                                                                        | Type                                                                             | Required                                                                         | Description                                                                      |
| -------------------------------------------------------------------------------- | -------------------------------------------------------------------------------- | -------------------------------------------------------------------------------- | -------------------------------------------------------------------------------- |
| `request`                                                                        | [operations.ListBalancesRequest](../../models/operations/listbalancesrequest.md) | :heavy_check_mark:                                                               | The request object to use for the request.                                       |


### Response

**[operations.ListBalancesResponse](../../models/operations/listbalancesresponse.md)**


## list_wallets

List all wallets

### Example Usage

```python
import sdk
from sdk.models import operations, shared

s = sdk.SDK(
    security=shared.Security(
        authorization="",
    ),
)

req = operations.ListWalletsRequest(
    cursor='aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ==',
    metadata={
        "aut": 'quas',
    },
    name='William Ortiz',
    page_size=984330,
)

res = s.wallets.list_wallets(req)

if res.list_wallets_response is not None:
    # handle response
```

### Parameters

| Parameter                                                                      | Type                                                                           | Required                                                                       | Description                                                                    |
| ------------------------------------------------------------------------------ | ------------------------------------------------------------------------------ | ------------------------------------------------------------------------------ | ------------------------------------------------------------------------------ |
| `request`                                                                      | [operations.ListWalletsRequest](../../models/operations/listwalletsrequest.md) | :heavy_check_mark:                                                             | The request object to use for the request.                                     |


### Response

**[operations.ListWalletsResponse](../../models/operations/listwalletsresponse.md)**


## update_wallet

Update a wallet

### Example Usage

```python
import sdk
from sdk.models import operations, shared

s = sdk.SDK(
    security=shared.Security(
        authorization="",
    ),
)

req = operations.UpdateWalletRequest(
    request_body=operations.UpdateWalletRequestBody(
        metadata={
            "facilis": 'cupiditate',
            "qui": 'quae',
        },
    ),
    id='879fce95-3f73-4ef7-bbc7-abd74dd39c0f',
)

res = s.wallets.update_wallet(req)

if res.status_code == 200:
    # handle response
```

### Parameters

| Parameter                                                                        | Type                                                                             | Required                                                                         | Description                                                                      |
| -------------------------------------------------------------------------------- | -------------------------------------------------------------------------------- | -------------------------------------------------------------------------------- | -------------------------------------------------------------------------------- |
| `request`                                                                        | [operations.UpdateWalletRequest](../../models/operations/updatewalletrequest.md) | :heavy_check_mark:                                                               | The request object to use for the request.                                       |


### Response

**[operations.UpdateWalletResponse](../../models/operations/updatewalletresponse.md)**


## void_hold

Cancel a hold

### Example Usage

```python
import sdk
from sdk.models import operations, shared

s = sdk.SDK(
    security=shared.Security(
        authorization="",
    ),
)

req = operations.VoidHoldRequest(
    hold_id='exercitationem',
)

res = s.wallets.void_hold(req)

if res.status_code == 200:
    # handle response
```

### Parameters

| Parameter                                                                | Type                                                                     | Required                                                                 | Description                                                              |
| ------------------------------------------------------------------------ | ------------------------------------------------------------------------ | ------------------------------------------------------------------------ | ------------------------------------------------------------------------ |
| `request`                                                                | [operations.VoidHoldRequest](../../models/operations/voidholdrequest.md) | :heavy_check_mark:                                                       | The request object to use for the request.                               |


### Response

**[operations.VoidHoldResponse](../../models/operations/voidholdresponse.md)**


## walletsget_server_info

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


res = s.wallets.walletsget_server_info()

if res.server_info is not None:
    # handle response
```


### Response

**[operations.WalletsgetServerInfoResponse](../../models/operations/walletsgetserverinforesponse.md)**

