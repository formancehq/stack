# logs

### Available Operations

* [list_logs](#list_logs) - List the logs from a ledger

## list_logs

List the logs from a ledger, sorted by ID in descending order.

### Example Usage

```python
import sdk
import dateutil.parser
from sdk.models import operations, shared

s = sdk.SDK(
    security=shared.Security(
        authorization="",
    ),
)

req = operations.ListLogsRequest(
    after='corporis',
    cursor='aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ==',
    end_time=dateutil.parser.isoparse('2022-04-01T23:59:21.675Z'),
    ledger='ledger001',
    page_size=315428,
    start_time=dateutil.parser.isoparse('2022-04-10T11:47:13.463Z'),
)

res = s.logs.list_logs(req)

if res.logs_cursor_response is not None:
    # handle response
```

### Parameters

| Parameter                                                                | Type                                                                     | Required                                                                 | Description                                                              |
| ------------------------------------------------------------------------ | ------------------------------------------------------------------------ | ------------------------------------------------------------------------ | ------------------------------------------------------------------------ |
| `request`                                                                | [operations.ListLogsRequest](../../models/operations/listlogsrequest.md) | :heavy_check_mark:                                                       | The request object to use for the request.                               |


### Response

**[operations.ListLogsResponse](../../models/operations/listlogsresponse.md)**

