# logs

### Available Operations

* [list_logs](#list_logs) - List the logs from a ledger

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
    end_time=dateutil.parser.isoparse('2022-09-04T08:35:09.957Z'),
    ledger='ledger001',
    page_size=570197,
    start_time=dateutil.parser.isoparse('2022-07-24T21:51:02.112Z'),
)

res = s.logs.list_logs(req)

if res.logs_cursor_response is not None:
    # handle response
```
