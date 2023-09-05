# accounts

### Available Operations

* [add_metadata_to_account](#add_metadata_to_account) - Add metadata to an account
* [count_accounts](#count_accounts) - Count the accounts from a ledger
* [get_account](#get_account) - Get account by its address
* [list_accounts](#list_accounts) - List accounts from a ledger

## add_metadata_to_account

Add metadata to an account

### Example Usage

```python
import sdk
from sdk.models import operations, shared

s = sdk.SDK(
    security=shared.Security(
        authorization="",
    ),
)

req = operations.AddMetadataToAccountRequest(
    request_body={
        "provident": 'distinctio',
        "quibusdam": 'unde',
        "nulla": 'corrupti',
    },
    address='users:001',
    ledger='ledger001',
)

res = s.accounts.add_metadata_to_account(req)

if res.status_code == 200:
    # handle response
```

### Parameters

| Parameter                                                                                        | Type                                                                                             | Required                                                                                         | Description                                                                                      |
| ------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------ |
| `request`                                                                                        | [operations.AddMetadataToAccountRequest](../../models/operations/addmetadatatoaccountrequest.md) | :heavy_check_mark:                                                                               | The request object to use for the request.                                                       |


### Response

**[operations.AddMetadataToAccountResponse](../../models/operations/addmetadatatoaccountresponse.md)**


## count_accounts

Count the accounts from a ledger

### Example Usage

```python
import sdk
from sdk.models import operations, shared

s = sdk.SDK(
    security=shared.Security(
        authorization="",
    ),
)

req = operations.CountAccountsRequest(
    address='users:.+',
    ledger='ledger001',
    metadata=operations.CountAccountsMetadata(),
)

res = s.accounts.count_accounts(req)

if res.status_code == 200:
    # handle response
```

### Parameters

| Parameter                                                                          | Type                                                                               | Required                                                                           | Description                                                                        |
| ---------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------- |
| `request`                                                                          | [operations.CountAccountsRequest](../../models/operations/countaccountsrequest.md) | :heavy_check_mark:                                                                 | The request object to use for the request.                                         |


### Response

**[operations.CountAccountsResponse](../../models/operations/countaccountsresponse.md)**


## get_account

Get account by its address

### Example Usage

```python
import sdk
from sdk.models import operations, shared

s = sdk.SDK(
    security=shared.Security(
        authorization="",
    ),
)

req = operations.GetAccountRequest(
    address='users:001',
    ledger='ledger001',
)

res = s.accounts.get_account(req)

if res.account_response is not None:
    # handle response
```

### Parameters

| Parameter                                                                    | Type                                                                         | Required                                                                     | Description                                                                  |
| ---------------------------------------------------------------------------- | ---------------------------------------------------------------------------- | ---------------------------------------------------------------------------- | ---------------------------------------------------------------------------- |
| `request`                                                                    | [operations.GetAccountRequest](../../models/operations/getaccountrequest.md) | :heavy_check_mark:                                                           | The request object to use for the request.                                   |


### Response

**[operations.GetAccountResponse](../../models/operations/getaccountresponse.md)**


## list_accounts

List accounts from a ledger, sorted by address in descending order.

### Example Usage

```python
import sdk
from sdk.models import operations, shared

s = sdk.SDK(
    security=shared.Security(
        authorization="",
    ),
)

req = operations.ListAccountsRequest(
    address='users:.+',
    after='users:003',
    balance=2400,
    balance_operator=operations.ListAccountsBalanceOperator.GTE,
    cursor='aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ==',
    ledger='ledger001',
    metadata=operations.ListAccountsMetadata(),
    page_size=847252,
)

res = s.accounts.list_accounts(req)

if res.accounts_cursor_response is not None:
    # handle response
```

### Parameters

| Parameter                                                                        | Type                                                                             | Required                                                                         | Description                                                                      |
| -------------------------------------------------------------------------------- | -------------------------------------------------------------------------------- | -------------------------------------------------------------------------------- | -------------------------------------------------------------------------------- |
| `request`                                                                        | [operations.ListAccountsRequest](../../models/operations/listaccountsrequest.md) | :heavy_check_mark:                                                               | The request object to use for the request.                                       |


### Response

**[operations.ListAccountsResponse](../../models/operations/listaccountsresponse.md)**

