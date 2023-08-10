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
    after='architecto',
    cursor='aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ==',
    end_time=dateutil.parser.isoparse('2022-08-01T12:28:44.292Z'),
    ledger='ledger001',
    page_size=635059,
    start_time=dateutil.parser.isoparse('2022-01-02T17:10:32.894Z'),
)

res = s.logs.list_logs(req)

if res.logs_cursor_response is not None:
    # handle response
```
