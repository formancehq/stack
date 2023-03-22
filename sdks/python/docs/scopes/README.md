# scopes

### Available Operations

* [add_transient_scope](#add_transient_scope) - Add a transient scope to a scope
* [create_scope](#create_scope) - Create scope
* [delete_scope](#delete_scope) - Delete scope
* [delete_transient_scope](#delete_transient_scope) - Delete a transient scope from a scope
* [list_scopes](#list_scopes) - List scopes
* [read_scope](#read_scope) - Read scope
* [update_scope](#update_scope) - Update scope

## add_transient_scope

Add a transient scope to a scope

### Example Usage

```python
import sdk
from sdk.models import operations

s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)

req = operations.AddTransientScopeRequest(
    scope_id='aspernatur',
    transient_scope_id='architecto',
)

res = s.scopes.add_transient_scope(req)

if res.status_code == 200:
    # handle response
```

## create_scope

Create scope

### Example Usage

```python
import sdk
from sdk.models import shared

s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)

req = shared.CreateScopeRequest(
    label='magnam',
    metadata={
        "excepturi": 'ullam',
    },
)

res = s.scopes.create_scope(req)

if res.create_scope_response is not None:
    # handle response
```

## delete_scope

Delete scope

### Example Usage

```python
import sdk
from sdk.models import operations

s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)

req = operations.DeleteScopeRequest(
    scope_id='provident',
)

res = s.scopes.delete_scope(req)

if res.status_code == 200:
    # handle response
```

## delete_transient_scope

Delete a transient scope from a scope

### Example Usage

```python
import sdk
from sdk.models import operations

s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)

req = operations.DeleteTransientScopeRequest(
    scope_id='quos',
    transient_scope_id='sint',
)

res = s.scopes.delete_transient_scope(req)

if res.status_code == 200:
    # handle response
```

## list_scopes

List Scopes

### Example Usage

```python
import sdk


s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)


res = s.scopes.list_scopes()

if res.list_scopes_response is not None:
    # handle response
```

## read_scope

Read scope

### Example Usage

```python
import sdk
from sdk.models import operations

s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)

req = operations.ReadScopeRequest(
    scope_id='accusantium',
)

res = s.scopes.read_scope(req)

if res.read_scope_response is not None:
    # handle response
```

## update_scope

Update scope

### Example Usage

```python
import sdk
from sdk.models import operations, shared

s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)

req = operations.UpdateScopeRequest(
    update_scope_request=shared.UpdateScopeRequest(
        label='mollitia',
        metadata={
            "mollitia": 'ad',
            "eum": 'dolor',
            "necessitatibus": 'odit',
            "nemo": 'quasi',
        },
    ),
    scope_id='iure',
)

res = s.scopes.update_scope(req)

if res.update_scope_response is not None:
    # handle response
```
