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
                    "amet": 'deserunt',
                    "nisi": 'vel',
                    "natus": 'omnis',
                    "molestiae": 'perferendis',
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
                timestamp=dateutil.parser.isoparse('2022-04-14T15:11:13.227Z'),
            ),
            shared.TransactionData(
                metadata={
                    "labore": 'labore',
                    "suscipit": 'natus',
                    "nobis": 'eum',
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
                    shared.Posting(
                        amount=100,
                        asset='COIN',
                        destination='users:002',
                        source='users:001',
                    ),
                ],
                reference='ref:001',
                timestamp=dateutil.parser.isoparse('2022-11-24T10:55:00.183Z'),
            ),
            shared.TransactionData(
                metadata={
                    "et": 'excepturi',
                    "ullam": 'provident',
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
                timestamp=dateutil.parser.isoparse('2022-12-07T10:53:17.121Z'),
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
        "reiciendis": 'mollitia',
        "ad": 'eum',
        "dolor": 'necessitatibus',
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
    end_time=dateutil.parser.isoparse('2022-08-19T20:09:28.183Z'),
    ledger='ledger001',
    metadata={
        "iure": 'doloribus',
    },
    reference='ref:001',
    source='users:001',
    start_time=dateutil.parser.isoparse('2022-03-21T22:14:24.691Z'),
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
    idempotency_key='maxime',
    post_transaction=shared.PostTransaction(
        metadata={
            "facilis": 'in',
            "architecto": 'architecto',
            "repudiandae": 'ullam',
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
                "repellat": 'quibusdam',
                "sed": 'saepe',
            },
        ),
        timestamp=dateutil.parser.isoparse('2022-11-20T20:56:20.791Z'),
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
    after='consequuntur',
    cursor='aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ==',
    destination='users:001',
    end_time=dateutil.parser.isoparse('2021-10-08T15:23:46.576Z'),
    ledger='ledger001',
    metadata={
        "sunt": 'quo',
    },
    page_size=848009,
    reference='ref:001',
    source='users:001',
    start_time=dateutil.parser.isoparse('2020-07-30T23:39:27.609Z'),
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
