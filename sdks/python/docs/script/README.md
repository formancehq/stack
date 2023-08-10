# script

### Available Operations

* [~~run_script~~](#run_script) - Execute a Numscript :warning: **Deprecated**

## ~~run_script~~

This route is deprecated, and has been merged into `POST /{ledger}/transactions`.


> :warning: **DEPRECATED**: this method will be removed in a future release, please migrate away from it as soon as possible.

### Example Usage

```python
import sdk
from sdk.models import operations, shared

s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)

req = operations.RunScriptRequest(
    script=shared.Script(
        metadata={
            "dolor": 'debitis',
            "a": 'dolorum',
            "in": 'in',
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
        vars={
            "maiores": 'rerum',
            "dicta": 'magnam',
            "cumque": 'facere',
            "ea": 'aliquid',
        },
    ),
    ledger='ledger001',
    preview=True,
)

res = s.script.run_script(req)

if res.script_response is not None:
    # handle response
```
