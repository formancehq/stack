# search

### Available Operations

* [search](#search) - Search
* [searchget_server_info](#searchget_server_info) - Get server info

## search

ElasticSearch query engine

### Example Usage

```python
import sdk
from sdk.models import shared

s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)

req = shared.Query(
    after=[
        'users:002',
        'users:002',
        'users:002',
    ],
    cursor='YXVsdCBhbmQgYSBtYXhpbXVtIG1heF9yZXN1bHRzLol=',
    ledgers=[
        'quickstart',
        'quickstart',
        'quickstart',
    ],
    page_size=725255,
    policy='OR',
    raw={
        "blanditiis": 'deleniti',
        "sapiente": 'amet',
        "deserunt": 'nisi',
    },
    sort='txid:asc',
    target='vel',
    terms=[
        'destination=central_bank1',
        'destination=central_bank1',
        'destination=central_bank1',
    ],
)

res = s.search.search(req)

if res.response is not None:
    # handle response
```

## searchget_server_info

Get server info

### Example Usage

```python
import sdk


s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)


res = s.search.searchget_server_info()

if res.server_info is not None:
    # handle response
```
