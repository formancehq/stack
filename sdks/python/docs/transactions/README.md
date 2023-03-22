# transactions

### Available Operations

* [add_metadata_on_transaction](#add_metadata_on_transaction) - Set the metadata of a transaction by its ID
* [count_transactions](#count_transactions) - Count the transactions from a ledger
* [create_transaction](#create_transaction) - Create a new transaction to a ledger
* [get_transaction](#get_transaction) - Get transaction from a ledger by its ID
* [list_transactions](#list_transactions) - List transactions from a ledger
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
    idempotency_key='incidunt',
    request_body={
        "consequatur": 'est',
        "quibusdam": 'explicabo',
    },
    async_=True,
    dry_run=True,
    ledger='ledger001',
    txid=1234,
)

res = s.transactions.add_metadata_on_transaction(req)

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
    end_time=dateutil.parser.isoparse('2021-07-27T01:56:50.693Z'),
    ledger='ledger001',
    metadata={
        "labore": 'modi',
        "qui": 'aliquid',
        "cupiditate": 'quos',
        "perferendis": 'magni',
    },
    reference='ref:001',
    source='users:001',
    start_time=dateutil.parser.isoparse('2021-11-22T01:26:35.048Z'),
)

res = s.transactions.count_transactions(req)

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
    idempotency_key='alias',
    post_transaction=shared.PostTransaction(
        metadata={
            "dolorum": 'excepturi',
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
                "tempore": 'labore',
                "delectus": 'eum',
                "non": 'eligendi',
            },
        ),
        timestamp=dateutil.parser.isoparse('2022-03-17T20:21:28.792Z'),
    ),
    async_=True,
    dry_run=True,
    ledger='ledger001',
)

res = s.transactions.create_transaction(req)

if res.create_transaction_response is not None:
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

res = s.transactions.get_transaction(req)

if res.get_transaction_response is not None:
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
    end_time=dateutil.parser.isoparse('2021-03-17T21:24:26.606Z'),
    ledger='ledger001',
    metadata={
        "officia": 'dolor',
        "debitis": 'a',
        "dolorum": 'in',
    },
    page_size=449198,
    reference='ref:001',
    source='users:001',
    start_time=dateutil.parser.isoparse('2020-01-25T11:09:22.009Z'),
)

res = s.transactions.list_transactions(req)

if res.transactions_cursor_response is not None:
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

res = s.transactions.revert_transaction(req)

if res.revert_transaction_response is not None:
    # handle response
```
