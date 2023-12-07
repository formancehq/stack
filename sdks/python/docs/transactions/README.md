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
            shared.TransactionData(
                metadata={
                    "expedita": 'nihil',
                    "repellat": 'quibusdam',
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
                timestamp=dateutil.parser.isoparse('2020-05-25T09:38:49.528Z'),
            ),
            shared.TransactionData(
                metadata={
                    "consequuntur": 'praesentium',
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
                    shared.Posting(
                        amount=100,
                        asset='COIN',
                        destination='users:002',
                        source='users:001',
                    ),
                ],
                reference='ref:001',
                timestamp=dateutil.parser.isoparse('2022-11-16T19:20:12.159Z'),
            ),
            shared.TransactionData(
                metadata={
                    "illum": 'pariatur',
                    "maxime": 'ea',
                    "excepturi": 'odit',
                    "ea": 'accusantium',
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
                timestamp=dateutil.parser.isoparse('2020-11-28T07:34:18.392Z'),
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
        "voluptate": 'autem',
        "nam": 'eaque',
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
    end_time=dateutil.parser.isoparse('2021-11-26T18:45:44.366Z'),
    ledger='ledger001',
    metadata={
        "perferendis": 'fugiat',
        "amet": 'aut',
        "cumque": 'corporis',
        "hic": 'libero',
    },
    reference='ref:001',
    source='users:001',
    start_time=dateutil.parser.isoparse('2022-08-28T17:02:52.151Z'),
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
    idempotency_key='quis',
    post_transaction=shared.PostTransaction(
        metadata={
            "dignissimos": 'eaque',
            "quis": 'nesciunt',
            "eos": 'perferendis',
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
                "quam": 'dolor',
                "vero": 'nostrum',
                "hic": 'recusandae',
                "omnis": 'facilis',
            },
        ),
        timestamp=dateutil.parser.isoparse('2022-12-08T18:10:54.422Z'),
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
    after='porro',
    cursor='aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ==',
    destination='users:001',
    end_time=dateutil.parser.isoparse('2022-07-02T11:46:10.299Z'),
    ledger='ledger001',
    metadata={
        "eaque": 'occaecati',
        "rerum": 'adipisci',
        "asperiores": 'earum',
    },
    page_size=267262,
    reference='ref:001',
    source='users:001',
    start_time=dateutil.parser.isoparse('2021-08-23T06:19:56.211Z'),
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
