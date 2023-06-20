# ledger

### Available Operations

* [add_metadata_on_transaction](#add_metadata_on_transaction) - Set the metadata of a transaction by its ID
* [add_metadata_to_account](#add_metadata_to_account) - Add metadata to an account
* [count_accounts](#count_accounts) - Count the accounts from a ledger
* [count_transactions](#count_transactions) - Count the transactions from a ledger
* [create_transaction](#create_transaction) - Create a new transaction to a ledger
* [get_account](#get_account) - Get account by its address
* [get_balances](#get_balances) - Get the balances from a ledger's account
* [get_balances_aggregated](#get_balances_aggregated) - Get the aggregated balances from selected accounts
* [get_info](#get_info) - Show server information
* [get_ledger_info](#get_ledger_info) - Get information about a ledger
* [get_transaction](#get_transaction) - Get transaction from a ledger by its ID
* [list_accounts](#list_accounts) - List accounts from a ledger
* [list_logs](#list_logs) - List the logs from a ledger
* [list_transactions](#list_transactions) - List transactions from a ledger
* [read_stats](#read_stats) - Get statistics from a ledger
* [revert_transaction](#revert_transaction) - Revert a ledger transaction by its ID

## add_metadata_on_transaction

Set the metadata of a transaction by its ID

### Example Usage

```python
import sdk
from sdk.models import operations

s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)

req = operations.AddMetadataOnTransactionRequest(
    idempotency_key='dolorem',
    request_body={
        "explicabo": 'nobis',
        "enim": 'omnis',
    },
    async_=True,
    dry_run=True,
    ledger='ledger001',
    txid=1234,
)

res = s.ledger.add_metadata_on_transaction(req)

if res.status_code == 200:
    # handle response
```

## add_metadata_to_account

Add metadata to an account

### Example Usage

```python
import sdk
from sdk.models import operations

s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)

req = operations.AddMetadataToAccountRequest(
    idempotency_key='nemo',
    request_body={
        "excepturi": 'accusantium',
        "iure": 'culpa',
    },
    address='users:001',
    async_=True,
    dry_run=True,
    ledger='ledger001',
)

res = s.ledger.add_metadata_to_account(req)

if res.status_code == 200:
    # handle response
```

## count_accounts

Count the accounts from a ledger

### Example Usage

```python
import sdk
from sdk.models import operations

s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)

req = operations.CountAccountsRequest(
    address='users:.+',
    ledger='ledger001',
    metadata={
        "sapiente": 'architecto',
        "mollitia": 'dolorem',
        "culpa": 'consequuntur',
        "repellat": 'mollitia',
    },
)

res = s.ledger.count_accounts(req)

if res.status_code == 200:
    # handle response
```

## count_transactions

Count the transactions from a ledger

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

req = operations.CountTransactionsRequest(
    account='users:001',
    destination='users:001',
    end_time=dateutil.parser.isoparse('2022-06-30T02:19:51.375Z'),
    ledger='ledger001',
    metadata={
        "quam": 'molestiae',
        "velit": 'error',
    },
    reference='ref:001',
    source='users:001',
    start_time=dateutil.parser.isoparse('2022-08-30T15:03:11.112Z'),
)

res = s.ledger.count_transactions(req)

if res.status_code == 200:
    # handle response
```

## create_transaction

Create a new transaction to a ledger

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

req = operations.CreateTransactionRequest(
    idempotency_key='vitae',
    post_transaction=shared.PostTransaction(
        metadata={
            "animi": 'enim',
            "odit": 'quo',
            "sequi": 'tenetur',
        },
        postings=[
            shared.Posting(
                amount=100,
                asset='COIN',
                destination='users:002',
                source='users:001',
            ),
            shared.Posting(
                amount=100,
                asset='COIN',
                destination='users:002',
                source='users:001',
            ),
        ],
        reference='ref:001',
        script=shared.PostTransactionScript(
            plain='vars {
        account $user
        }
        send [COIN 10] (
        	source = @world
        	destination = $user
        )
        ',
            vars={
                "possimus": 'aut',
                "quasi": 'error',
                "temporibus": 'laborum',
            },
        ),
        timestamp=dateutil.parser.isoparse('2022-01-11T05:45:42.485Z'),
    ),
    async_=True,
    dry_run=True,
    ledger='ledger001',
)

res = s.ledger.create_transaction(req)

if res.create_transaction_response is not None:
    # handle response
```

## get_account

Get account by its address

### Example Usage

```python
import sdk
from sdk.models import operations

s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)

req = operations.GetAccountRequest(
    address='users:001',
    ledger='ledger001',
)

res = s.ledger.get_account(req)

if res.account_response is not None:
    # handle response
```

## get_balances

Get the balances from a ledger's account

### Example Usage

```python
import sdk
from sdk.models import operations

s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)

req = operations.GetBalancesRequest(
    address='users:001',
    cursor='aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ==',
    ledger='ledger001',
    page_size=976460,
)

res = s.ledger.get_balances(req)

if res.balances_cursor_response is not None:
    # handle response
```

## get_balances_aggregated

Get the aggregated balances from selected accounts

### Example Usage

```python
import sdk
from sdk.models import operations

s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)

req = operations.GetBalancesAggregatedRequest(
    address='users:001',
    ledger='ledger001',
)

res = s.ledger.get_balances_aggregated(req)

if res.aggregate_balances_response is not None:
    # handle response
```

## get_info

Show server information

### Example Usage

```python
import sdk


s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)


res = s.ledger.get_info()

if res.config_info_response is not None:
    # handle response
```

## get_ledger_info

Get information about a ledger

### Example Usage

```python
import sdk
from sdk.models import operations

s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)

req = operations.GetLedgerInfoRequest(
    ledger='ledger001',
)

res = s.ledger.get_ledger_info(req)

if res.ledger_info_response is not None:
    # handle response
```

## get_transaction

Get transaction from a ledger by its ID

### Example Usage

```python
import sdk
from sdk.models import operations

s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)

req = operations.GetTransactionRequest(
    ledger='ledger001',
    txid=1234,
)

res = s.ledger.get_transaction(req)

if res.get_transaction_response is not None:
    # handle response
```

## list_accounts

List accounts from a ledger, sorted by address in descending order.

### Example Usage

```python
import sdk
from sdk.models import operations

s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)

req = operations.ListAccountsRequest(
    address='users:.+',
    balance=2400,
    balance_operator=operations.ListAccountsBalanceOperator.GTE,
    cursor='aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ==',
    ledger='ledger001',
    metadata={
        "nihil": 'praesentium',
        "voluptatibus": 'ipsa',
        "omnis": 'voluptate',
        "cum": 'perferendis',
    },
    page_size=39187,
)

res = s.ledger.list_accounts(req)

if res.accounts_cursor_response is not None:
    # handle response
```

## list_logs

List the logs from a ledger, sorted by ID in descending order.

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

req = operations.ListLogsRequest(
    cursor='aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ==',
    end_time=dateutil.parser.isoparse('2022-09-19T18:36:39.009Z'),
    ledger='ledger001',
    page_size=979587,
    start_time=dateutil.parser.isoparse('2022-08-22T19:15:58.586Z'),
)

res = s.ledger.list_logs(req)

if res.logs_cursor_response is not None:
    # handle response
```

## list_transactions

List transactions from a ledger, sorted by txid in descending order.

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

req = operations.ListTransactionsRequest(
    account='users:001',
    cursor='aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ==',
    destination='users:001',
    end_time=dateutil.parser.isoparse('2022-07-09T11:22:20.922Z'),
    ledger='ledger001',
    metadata={
        "harum": 'enim',
    },
    page_size=880476,
    reference='ref:001',
    source='users:001',
    start_time=dateutil.parser.isoparse('2022-01-30T20:15:26.045Z'),
)

res = s.ledger.list_transactions(req)

if res.transactions_cursor_response is not None:
    # handle response
```

## read_stats

Get statistics from a ledger. (aggregate metrics on accounts and transactions)


### Example Usage

```python
import sdk
from sdk.models import operations

s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)

req = operations.ReadStatsRequest(
    ledger='ledger001',
)

res = s.ledger.read_stats(req)

if res.stats_response is not None:
    # handle response
```

## revert_transaction

Revert a ledger transaction by its ID

### Example Usage

```python
import sdk
from sdk.models import operations

s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)

req = operations.RevertTransactionRequest(
    ledger='ledger001',
    txid=1234,
)

res = s.ledger.revert_transaction(req)

if res.revert_transaction_response is not None:
    # handle response
```
