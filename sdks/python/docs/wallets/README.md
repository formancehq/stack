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
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)

req = operations.ConfirmHoldRequest(
    confirm_hold_request=shared.ConfirmHoldRequest(
        amount=100,
        final=True,
    ),
    hold_id='labore',
)

res = s.wallets.confirm_hold(req)

if res.status_code == 200:
    # handle response
```

## create_balance

Create a balance

### Example Usage

```python
import sdk
import dateutil.parser
from sdk.models import operations, shared

s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)

req = operations.CreateBalanceRequest(
    create_balance_request=shared.CreateBalanceRequest(
        expires_at=dateutil.parser.isoparse('2022-08-14T00:52:14.624Z'),
        name='Robin Keebler',
        priority=102863,
    ),
    id='41959890-afa5-463e-a516-fe4c8b711e5b',
)

res = s.wallets.create_balance(req)

if res.create_balance_response is not None:
    # handle response
```

## create_wallet

Create a new wallet

### Example Usage

```python
import sdk
from sdk.models import shared

s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)

req = shared.CreateWalletRequest(
    metadata={
        "repellat": 'quibusdam',
        "sed": 'saepe',
    },
    name='Edward Crooks',
)

res = s.wallets.create_wallet(req)

if res.create_wallet_response is not None:
    # handle response
```

## credit_wallet

Credit a wallet

### Example Usage

```python
import sdk
from sdk.models import operations, shared

s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)

req = operations.CreditWalletRequest(
    credit_wallet_request=shared.CreditWalletRequest(
        amount=shared.Monetary(
            amount=166847,
            asset='sunt',
        ),
        balance='quo',
        metadata={
            "pariatur": 'maxime',
            "ea": 'excepturi',
            "odit": 'ea',
            "accusantium": 'ab',
        },
        reference='maiores',
        sources=[
            shared.LedgerAccountSubject(
                identifier='voluptate',
                type='autem',
            ),
            shared.WalletSubject(
                balance='eaque',
                identifier='pariatur',
                type='nemo',
            ),
            shared.WalletSubject(
                balance='perferendis',
                identifier='fugiat',
                type='amet',
            ),
        ],
    ),
    id='0c5fbb25-8705-4320-ac73-d5fe9b90c289',
)

res = s.wallets.credit_wallet(req)

if res.status_code == 200:
    # handle response
```

## debit_wallet

Debit a wallet

### Example Usage

```python
import sdk
from sdk.models import operations, shared

s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)

req = operations.DebitWalletRequest(
    debit_wallet_request=shared.DebitWalletRequest(
        amount=shared.Monetary(
            amount=50370,
            asset='occaecati',
        ),
        balances=[
            'adipisci',
            'asperiores',
            'earum',
        ],
        description='modi',
        destination=shared.WalletSubject(
            balance='dolorum',
            identifier='deleniti',
            type='pariatur',
        ),
        metadata={
            "nobis": 'libero',
            "delectus": 'quaerat',
            "quos": 'aliquid',
        },
        pending=False,
    ),
    id='33323f9b-77f3-4a41-8067-4ebf69280d1b',
)

res = s.wallets.debit_wallet(req)

if res.debit_wallet_response is not None:
    # handle response
```

## get_balance

Get detailed balance

### Example Usage

```python
import sdk
from sdk.models import operations

s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)

req = operations.GetBalanceRequest(
    balance_name='dolorum',
    id='77a89ebf-737a-4e42-83ce-5e6a95d8a0d4',
)

res = s.wallets.get_balance(req)

if res.get_balance_response is not None:
    # handle response
```

## get_hold

Get a hold

### Example Usage

```python
import sdk
from sdk.models import operations

s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)

req = operations.GetHoldRequest(
    hold_id='tempora',
)

res = s.wallets.get_hold(req)

if res.get_hold_response is not None:
    # handle response
```

## get_holds

Get all holds for a wallet

### Example Usage

```python
import sdk
from sdk.models import operations

s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)

req = operations.GetHoldsRequest(
    cursor='aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ==',
    metadata={
        "quod": 'officiis',
        "qui": 'dolorum',
    },
    page_size=952792,
    wallet_id='esse',
)

res = s.wallets.get_holds(req)

if res.get_holds_response is not None:
    # handle response
```

## get_transactions

### Example Usage

```python
import sdk
from sdk.models import operations

s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)

req = operations.GetTransactionsRequest(
    cursor='aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ==',
    page_size=687488,
    wallet_id='iusto',
)

res = s.wallets.get_transactions(req)

if res.get_transactions_response is not None:
    # handle response
```

## get_wallet

Get a wallet

### Example Usage

```python
import sdk
from sdk.models import operations

s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)

req = operations.GetWalletRequest(
    id='3cf3be45-3f87-40b3-a6b5-a73429cdb1a8',
)

res = s.wallets.get_wallet(req)

if res.get_wallet_response is not None:
    # handle response
```

## get_wallet_summary

Get wallet summary

### Example Usage

```python
import sdk
from sdk.models import operations

s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)

req = operations.GetWalletSummaryRequest(
    id='422bb679-d232-4271-9bf0-cbb1e31b8b90',
)

res = s.wallets.get_wallet_summary(req)

if res.get_wallet_summary_response is not None:
    # handle response
```

## list_balances

List balances of a wallet

### Example Usage

```python
import sdk
from sdk.models import operations

s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)

req = operations.ListBalancesRequest(
    id='f3443a11-08e0-4adc-b4b9-21879fce953f',
)

res = s.wallets.list_balances(req)

if res.list_balances_response is not None:
    # handle response
```

## list_wallets

List all wallets

### Example Usage

```python
import sdk
from sdk.models import operations

s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)

req = operations.ListWalletsRequest(
    cursor='aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ==',
    metadata={
        "consectetur": 'vero',
        "tenetur": 'dignissimos',
    },
    name='Kelvin Schmidt',
    page_size=708548,
)

res = s.wallets.list_wallets(req)

if res.list_wallets_response is not None:
    # handle response
```

## update_wallet

Update a wallet

### Example Usage

```python
import sdk
from sdk.models import operations

s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)

req = operations.UpdateWalletRequest(
    request_body=operations.UpdateWalletRequestBody(
        metadata={
            "ducimus": 'dolore',
            "quibusdam": 'illum',
            "sequi": 'natus',
            "impedit": 'aut',
        },
    ),
    id='f5d2cff7-c70a-4456-a6d4-36813f16d9f5',
)

res = s.wallets.update_wallet(req)

if res.status_code == 200:
    # handle response
```

## void_hold

Cancel a hold

### Example Usage

```python
import sdk
from sdk.models import operations

s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)

req = operations.VoidHoldRequest(
    hold_id='sapiente',
)

res = s.wallets.void_hold(req)

if res.status_code == 200:
    # handle response
```

## walletsget_server_info

Get server info

### Example Usage

```python
import sdk


s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)


res = s.wallets.walletsget_server_info()

if res.server_info is not None:
    # handle response
```
