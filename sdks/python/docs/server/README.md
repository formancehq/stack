# server

### Available Operations

* [get_info](#get_info) - Show server information

## get_info

Show server information

### Example Usage

```python
import sdk


s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)


res = s.server.get_info()

if res.config_info_response is not None:
    # handle response
```
