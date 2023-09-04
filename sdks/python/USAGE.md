<!-- Start SDK Example Usage -->


```python
import sdk
from sdk.models import shared

s = sdk.SDK(
    security=shared.Security(
        authorization="",
    ),
)


res = s.get_versions()

if res.get_versions_response is not None:
    # handle response
```
<!-- End SDK Example Usage -->