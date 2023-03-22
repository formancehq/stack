# clients

### Available Operations

* [add_scope_to_client](#add_scope_to_client) - Add scope to client
* [create_client](#create_client) - Create client
* [create_secret](#create_secret) - Add a secret to a client
* [delete_client](#delete_client) - Delete client
* [delete_scope_from_client](#delete_scope_from_client) - Delete scope from client
* [delete_secret](#delete_secret) - Delete a secret from a client
* [list_clients](#list_clients) - List clients
* [read_client](#read_client) - Read client
* [update_client](#update_client) - Update client

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
    client_id='velit',
    scope_id='error',
)

res = s.clients.add_scope_to_client(req)

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
    description='quia',
    metadata={
        "vitae": 'laborum',
        "animi": 'enim',
    },
    name='Angelica Dietrich',
    post_logout_redirect_uris=[
        'possimus',
        'aut',
        'quasi',
    ],
    public=False,
    redirect_uris=[
        'temporibus',
        'laborum',
        'quasi',
    ],
    trusted=False,
)

res = s.clients.create_client(req)

if res.create_client_response is not None:
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
            "voluptatibus": 'vero',
            "nihil": 'praesentium',
            "voluptatibus": 'ipsa',
            "omnis": 'voluptate',
        },
        name='Thomas Batz',
    ),
    client_id='maiores',
)

res = s.clients.create_secret(req)

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
    client_id='dicta',
)

res = s.clients.delete_client(req)

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
    client_id='corporis',
    scope_id='dolore',
)

res = s.clients.delete_scope_from_client(req)

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
    client_id='iusto',
    secret_id='dicta',
)

res = s.clients.delete_secret(req)

if res.status_code == 200:
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


res = s.clients.list_clients()

if res.list_clients_response is not None:
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
    client_id='harum',
)

res = s.clients.read_client(req)

if res.read_client_response is not None:
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
        description='enim',
        metadata={
            "commodi": 'repudiandae',
            "quae": 'ipsum',
            "quidem": 'molestias',
            "excepturi": 'pariatur',
        },
        name='Irma Ledner DVM',
        post_logout_redirect_uris=[
            'veritatis',
            'itaque',
            'incidunt',
        ],
        public=False,
        redirect_uris=[
            'consequatur',
            'est',
        ],
        trusted=False,
    ),
    client_id='quibusdam',
)

res = s.clients.update_client(req)

if res.update_client_response is not None:
    # handle response
```
