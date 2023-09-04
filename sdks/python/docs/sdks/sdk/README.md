# SDK

## Overview

Formance Stack API: Open, modular foundation for unique payments flows

# Introduction
This API is documented in **OpenAPI format**.

# Authentication
Formance Stack offers one forms of authentication:
  - OAuth2
OAuth2 - an open protocol to allow secure authorization in a simple
and standard method from web, mobile and desktop applications.
<SecurityDefinitions />


### Available Operations

* [get_versions](#get_versions) - Show stack version information

## get_versions

Show stack version information

### Example Usage

```python
import sdk
from sdk.models import shared

s = sdk.SDK(
    security=shared.Security(
        authorization="",
    ),
)


res = s.sdk.get_versions()

if res.get_versions_response is not None:
    # handle response
```


### Response

**[operations.GetVersionsResponse](../../models/operations/getversionsresponse.md)**

