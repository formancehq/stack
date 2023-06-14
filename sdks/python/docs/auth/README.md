# auth

### Available Operations

* [add_scope_to_client](#add_scope_to_client) - Add scope to client
* [add_transient_scope](#add_transient_scope) - Add a transient scope to a scope
* [create_client](#create_client) - Create client
* [create_scope](#create_scope) - Create scope
* [create_secret](#create_secret) - Add a secret to a client
* [delete_client](#delete_client) - Delete client
* [delete_scope](#delete_scope) - Delete scope
* [delete_scope_from_client](#delete_scope_from_client) - Delete scope from client
* [delete_secret](#delete_secret) - Delete a secret from a client
* [delete_transient_scope](#delete_transient_scope) - Delete a transient scope from a scope
* [get_server_info](#get_server_info) - Get server info
* [list_clients](#list_clients) - List clients
* [list_scopes](#list_scopes) - List scopes
* [list_users](#list_users) - List users
* [read_client](#read_client) - Read client
* [read_scope](#read_scope) - Read scope
* [read_user](#read_user) - Read user
* [update_client](#update_client) - Update client
* [update_scope](#update_scope) - Update scope

## add_scope_to_client

Add scope to client

### Example Usage

```python
import sdk
from sdk.models import operations

s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)

req = operations.AddScopeToClientRequest(
    client_id='corrupti',
    scope_id='provident',
)

res = s.auth.add_scope_to_client(req)

if res.status_code == 200:
    # handle response
```

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
    scope_id='distinctio',
    transient_scope_id='quibusdam',
)

res = s.auth.add_transient_scope(req)

if res.status_code == 200:
    # handle response
```

## create_client

Create client

### Example Usage

```python
import sdk
from sdk.models import shared

s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)

req = shared.CreateClientRequest(
    description='unde',
    metadata={
        "corrupti": 'illum',
        "vel": 'error',
        "deserunt": 'suscipit',
        "iure": 'magnam',
    },
    name='Larry Windler',
    post_logout_redirect_uris=[
        'minus',
        'placeat',
    ],
    public=False,
    redirect_uris=[
        'iusto',
        'excepturi',
        'nisi',
    ],
    trusted=False,
)

res = s.auth.create_client(req)

if res.create_client_response is not None:
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
    label='recusandae',
    metadata={
        "ab": 'quis',
        "veritatis": 'deserunt',
        "perferendis": 'ipsam',
        "repellendus": 'sapiente',
    },
)

res = s.auth.create_scope(req)

if res.create_scope_response is not None:
    # handle response
```

## create_secret

Add a secret to a client

### Example Usage

```python
import sdk
from sdk.models import operations, shared

s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)

req = operations.CreateSecretRequest(
    create_secret_request=shared.CreateSecretRequest(
        metadata={
            "odit": 'at',
            "at": 'maiores',
            "molestiae": 'quod',
            "quod": 'esse',
        },
        name='Miss Lowell Parisian',
    ),
    client_id='occaecati',
)

res = s.auth.create_secret(req)

if res.create_secret_response is not None:
    # handle response
```

## delete_client

Delete client

### Example Usage

```python
import sdk
from sdk.models import operations

s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)

req = operations.DeleteClientRequest(
    client_id='fugit',
)

res = s.auth.delete_client(req)

if res.status_code == 200:
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
    scope_id='deleniti',
)

res = s.auth.delete_scope(req)

if res.status_code == 200:
    # handle response
```

## delete_scope_from_client

Delete scope from client

### Example Usage

```python
import sdk
from sdk.models import operations

s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)

req = operations.DeleteScopeFromClientRequest(
    client_id='hic',
    scope_id='optio',
)

res = s.auth.delete_scope_from_client(req)

if res.status_code == 200:
    # handle response
```

## delete_secret

Delete a secret from a client

### Example Usage

```python
import sdk
from sdk.models import operations

s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)

req = operations.DeleteSecretRequest(
    client_id='totam',
    secret_id='beatae',
)

res = s.auth.delete_secret(req)

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
    scope_id='commodi',
    transient_scope_id='molestiae',
)

res = s.auth.delete_transient_scope(req)

if res.status_code == 200:
    # handle response
```

## get_server_info

Get server info

### Example Usage

```python
import sdk


s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)


res = s.auth.get_server_info()

if res.server_info is not None:
    # handle response
```

## list_clients

List clients

### Example Usage

```python
import sdk


s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)


res = s.auth.list_clients()

if res.list_clients_response is not None:
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


res = s.auth.list_scopes()

if res.list_scopes_response is not None:
    # handle response
```

## list_users

List users

### Example Usage

```python
import sdk


s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)


res = s.auth.list_users()

if res.list_users_response is not None:
    # handle response
```

## read_client

Read client

### Example Usage

```python
import sdk
from sdk.models import operations

s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)

req = operations.ReadClientRequest(
    client_id='modi',
)

res = s.auth.read_client(req)

if res.read_client_response is not None:
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
    scope_id='qui',
)

res = s.auth.read_scope(req)

if res.read_scope_response is not None:
    # handle response
```

## read_user

Read user

### Example Usage

```python
import sdk
from sdk.models import operations

s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)

req = operations.ReadUserRequest(
    user_id='impedit',
)

res = s.auth.read_user(req)

if res.read_user_response is not None:
    # handle response
```

## update_client

Update client

### Example Usage

```python
import sdk
from sdk.models import operations, shared

s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)

req = operations.UpdateClientRequest(
    update_client_request=shared.UpdateClientRequest(
        description='cum',
        metadata={
            "ipsum": 'excepturi',
            "aspernatur": 'perferendis',
        },
        name='Faye Cormier',
        post_logout_redirect_uris=[
            'laboriosam',
            'hic',
            'saepe',
        ],
        public=False,
        redirect_uris=[
            'in',
            'corporis',
            'iste',
        ],
        trusted=False,
    ),
    client_id='iure',
)

res = s.auth.update_client(req)

if res.update_client_response is not None:
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
        label='saepe',
        metadata={
            "architecto": 'ipsa',
            "reiciendis": 'est',
            "mollitia": 'laborum',
        },
    ),
    scope_id='dolores',
)

res = s.auth.update_scope(req)

if res.update_scope_response is not None:
    # handle response
```
