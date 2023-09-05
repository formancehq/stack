# stats

### Available Operations

* [read_stats](#read_stats) - Get statistics from a ledger

## read_stats

Get statistics from a ledger. (aggregate metrics on accounts and transactions)


### Example Usage

```python
import sdk
from sdk.models import operations, shared

s = sdk.SDK(
    security=shared.Security(
        authorization="",
    ),
)

req = operations.ReadStatsRequest(
    ledger='ledger001',
)

res = s.stats.read_stats(req)

if res.stats_response is not None:
    # handle response
```

### Parameters

| Parameter                                                                  | Type                                                                       | Required                                                                   | Description                                                                |
| -------------------------------------------------------------------------- | -------------------------------------------------------------------------- | -------------------------------------------------------------------------- | -------------------------------------------------------------------------- |
| `request`                                                                  | [operations.ReadStatsRequest](../../models/operations/readstatsrequest.md) | :heavy_check_mark:                                                         | The request object to use for the request.                                 |


### Response

**[operations.ReadStatsResponse](../../models/operations/readstatsresponse.md)**

