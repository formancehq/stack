# script

### Available Operations

* [~~run_script~~](#run_script) - Execute a Numscript :warning: **Deprecated**

## ~~run_script~~

This route is deprecated, and has been merged into `POST /{ledger}/transactions`.


> :warning: **DEPRECATED**: This will be removed in a future release, please migrate away from it as soon as possible.

### Example Usage

```python
import sdk
from sdk.models import operations, shared

s = sdk.SDK(
    security=shared.Security(
        authorization="",
    ),
)

req = operations.RunScriptRequest(
    script=shared.Script(
        metadata={
            "rem": 'voluptates',
            "quasi": 'repudiandae',
            "sint": 'veritatis',
        },
        plain='vars {
    account $user
    }
    send [COIN 10] (
    	source = @world
    	destination = $user
    )
    ',
        reference='order_1234',
        vars=shared.ScriptVars(),
    ),
    ledger='ledger001',
    preview=True,
)

res = s.script.run_script(req)

if res.script_response is not None:
    # handle response
```

### Parameters

| Parameter                                                                  | Type                                                                       | Required                                                                   | Description                                                                |
| -------------------------------------------------------------------------- | -------------------------------------------------------------------------- | -------------------------------------------------------------------------- | -------------------------------------------------------------------------- |
| `request`                                                                  | [operations.RunScriptRequest](../../models/operations/runscriptrequest.md) | :heavy_check_mark:                                                         | The request object to use for the request.                                 |


### Response

**[operations.RunScriptResponse](../../models/operations/runscriptresponse.md)**

