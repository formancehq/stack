# transactions

### Available Operations

* [create_transactions](#create_transactions) - Create a new batch of transactions to a ledger
* [add_metadata_on_transaction](#add_metadata_on_transaction) - Set the metadata of a transaction by its ID
* [count_transactions](#count_transactions) - Count the transactions from a ledger
* [create_transaction](#create_transaction) - Create a new transaction to a ledger
* [get_transaction](#get_transaction) - Get transaction from a ledger by its ID
* [list_transactions](#list_transactions) - List transactions from a ledger
* [revert_transaction](#revert_transaction) - Revert a ledger transaction by its ID

## create_transactions

Create a new batch of transactions to a ledger

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

req = operations.CreateTransactionsRequest(
    transactions=shared.Transactions(
        transactions=[
            shared.TransactionData(
                metadata={
                    "sint": 'accusantium',
                    "mollitia": 'reiciendis',
                    "mollitia": 'ad',
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
                timestamp=dateutil.parser.isoparse('2022-02-07T18:15:06.372Z'),
            ),
            shared.TransactionData(
                metadata={
                    "nemo": 'quasi',
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
                timestamp=dateutil.parser.isoparse('2020-04-29T08:15:14.819Z'),
            ),
            shared.TransactionData(
                metadata={
                    "maxime": 'deleniti',
                    "facilis": 'in',
                },
                postings=[
                    shared.Posting(
                        amount=100,
                        asset='COIN',
                        destination='users:002',
                        source='users:001',
                    ),
                ],
                reference='ref:001',
                timestamp=dateutil.parser.isoparse('2022-01-30T09:19:56.236Z'),
            ),
        ],
    ),
    ledger='ledger001',
)

res = s.transactions.create_transactions(req)

if res.transactions_response is not None:
    # handle response
```

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
    request_body={
        "expedita": 'nihil',
        "repellat": 'quibusdam',
    },
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
    end_time=dateutil.parser.isoparse('2022-02-04T19:17:08.641Z'),
    ledger='ledger001',
    metadata={
        "accusantium": 'consequuntur',
        "praesentium": 'natus',
        "magni": 'sunt',
        "quo": 'illum',
    },
    reference='ref:001',
    source='users:001',
    start_time=dateutil.parser.isoparse('2020-07-30T23:39:27.609Z'),
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
    idempotency_key='ea',
    post_transaction=shared.PostTransaction(
        metadata={
            "odit": 'ea',
            "accusantium": 'ab',
            "maiores": 'quidem',
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
                "autem": 'nam',
                "eaque": 'pariatur',
            },
        ),
        timestamp=dateutil.parser.isoparse('2022-01-09T22:25:53.570Z'),
    ),
    ledger='ledger001',
    preview=True,
)

res = s.transactions.create_transaction(req)

if res.transactions_response is not None:
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

if res.transaction_response is not None:
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
    after='perferendis',
    cursor='aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ==',
    destination='users:001',
    end_time=dateutil.parser.isoparse('2022-04-23T08:05:27.776Z'),
    ledger='ledger001',
    metadata={
        "cumque": 'corporis',
    },
    page_size=944124,
    reference='ref:001',
    source='users:001',
    start_time=dateutil.parser.isoparse('2021-07-02T12:00:47.197Z'),
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

if res.transaction_response is not None:
    # handle response
```
