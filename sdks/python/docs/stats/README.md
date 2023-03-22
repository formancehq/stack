# stats

### Available Operations

* [read_stats](#read_stats) - Get statistics from a ledger

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

res = s.stats.read_stats(req)

if res.stats_response is not None:
    # handle response
```
