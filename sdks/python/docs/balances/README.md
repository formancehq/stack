# balances

### Available Operations

* [get_balances](#get_balances) - Get the balances from a ledger's account
* [get_balances_aggregated](#get_balances_aggregated) - Get the aggregated balances from selected accounts

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
    after='users:003',
    cursor='aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ==',
    ledger='ledger001',
    page_size=958950,
)

res = s.balances.get_balances(req)

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

res = s.balances.get_balances_aggregated(req)

if res.aggregate_balances_response is not None:
    # handle response
```
