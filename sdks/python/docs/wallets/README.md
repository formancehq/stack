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
    hold_id='blanditiis',
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
        expires_at=dateutil.parser.isoparse('2021-02-02T01:24:52.629Z'),
        name='Sandy Huels',
        priority=606393,
    ),
    id='7074ba44-69b6-4e21-8195-9890afa563e2',
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
        "quasi": 'iure',
        "doloribus": 'debitis',
    },
    name='Jasmine Lind',
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
            amount=100226,
            asset='architecto',
        ),
        balance='repudiandae',
        metadata={
            "expedita": 'nihil',
            "repellat": 'quibusdam',
        },
        reference='sed',
        sources=[
            shared.WalletSubject(
                balance='accusantium',
                identifier='consequuntur',
                type='praesentium',
            ),
            shared.WalletSubject(
                balance='magni',
                identifier='sunt',
                type='quo',
            ),
            shared.WalletSubject(
                balance='pariatur',
                identifier='maxime',
                type='ea',
            ),
            shared.WalletSubject(
                balance='odit',
                identifier='ea',
                type='accusantium',
            ),
        ],
    ),
    id='1fb576b0-d5f0-4d30-85fb-b2587053202c',
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
            amount=463451,
            asset='dolor',
        ),
        balances=[
            'nostrum',
            'hic',
            'recusandae',
            'omnis',
        ],
        description='facilis',
        destination=shared.WalletSubject(
            balance='voluptatem',
            identifier='porro',
            type='consequuntur',
        ),
        metadata={
            "error": 'eaque',
            "occaecati": 'rerum',
            "adipisci": 'asperiores',
        },
        pending=False,
    ),
    id='e49a8d9c-bf48-4633-b23f-9b77f3a41006',
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
    balance_name='odio',
    id='4ebf6928-0d1b-4a77-a89e-bf737ae4203c',
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
    hold_id='accusamus',
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
        "saepe": 'suscipit',
        "deserunt": 'provident',
    },
    page_size=324683,
    wallet_id='repellendus',
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
    page_size=519711,
    wallet_id='similique',
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
    id='0d446ce2-af7a-473c-b3be-453f870b326b',
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
    id='5a73429c-db1a-4842-abb6-79d2322715bf',
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
    id='0cbb1e31-b8b9-40f3-843a-1108e0adcf4b',
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
        "qui": 'quae',
        "laudantium": 'odio',
        "occaecati": 'voluptatibus',
    },
    name='Ignacio Moen',
    page_size=961571,
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
            "consectetur": 'vero',
            "tenetur": 'dignissimos',
        },
    ),
    id='fbc7abd7-4dd3-49c0-b5d2-cff7c70a4562',
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
    hold_id='vel',
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
