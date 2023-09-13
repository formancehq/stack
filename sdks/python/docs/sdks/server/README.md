# server

### Available Operations

* [get_info](#get_info) - Show server information

## get_info

Show server information

### Example Usage

```python
import sdk
from sdk.models import shared

s = sdk.SDK(
    security=shared.Security(
        authorization="",
    ),
)


res = s.server.get_info()

if res.config_info_response is not None:
    # handle response
```


### Response

**[operations.GetInfoResponse](../../models/operations/getinforesponse.md)**

