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
        authorization="",
    ),
)

req = shared.Query(
    after=[
        'users:002',
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
    page_size=588465,
    policy='OR',
    raw=shared.QueryRaw(),
    sort='id:asc',
    target='nam',
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

### Parameters

| Parameter                                    | Type                                         | Required                                     | Description                                  |
| -------------------------------------------- | -------------------------------------------- | -------------------------------------------- | -------------------------------------------- |
| `request`                                    | [shared.Query](../../models/shared/query.md) | :heavy_check_mark:                           | The request object to use for the request.   |


### Response

**[operations.SearchResponse](../../models/operations/searchresponse.md)**


## searchget_server_info

Get server info

### Example Usage

```python
import sdk
from sdk.models import shared

s = sdk.SDK(
    security=shared.Security(
        authorization="",
    ),
)


res = s.search.searchget_server_info()

if res.server_info is not None:
    # handle response
```


### Response

**[operations.SearchgetServerInfoResponse](../../models/operations/searchgetserverinforesponse.md)**

