<!-- Start SDK Example Usage -->
```python
import sdk


s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)


res = s.get_versions()

if res.get_versions_response is not None:
    # handle response
```
<!-- End SDK Example Usage -->