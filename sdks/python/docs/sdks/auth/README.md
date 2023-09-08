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
from sdk.models import operations, shared

s = sdk.SDK(
    security=shared.Security(
        authorization="",
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

### Parameters

| Parameter                                                                                | Type                                                                                     | Required                                                                                 | Description                                                                              |
| ---------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------- |
| `request`                                                                                | [operations.AddScopeToClientRequest](../../models/operations/addscopetoclientrequest.md) | :heavy_check_mark:                                                                       | The request object to use for the request.                                               |


### Response

**[operations.AddScopeToClientResponse](../../models/operations/addscopetoclientresponse.md)**


## add_transient_scope

Add a transient scope to a scope

### Example Usage

```python
import sdk
from sdk.models import operations, shared

s = sdk.SDK(
    security=shared.Security(
        authorization="",
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

### Parameters

| Parameter                                                                                  | Type                                                                                       | Required                                                                                   | Description                                                                                |
| ------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------ |
| `request`                                                                                  | [operations.AddTransientScopeRequest](../../models/operations/addtransientscoperequest.md) | :heavy_check_mark:                                                                         | The request object to use for the request.                                                 |


### Response

**[operations.AddTransientScopeResponse](../../models/operations/addtransientscoperesponse.md)**


## create_client

Create client

### Example Usage

```python
import sdk
from sdk.models import shared

s = sdk.SDK(
    security=shared.Security(
        authorization="",
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

### Parameters

| Parameter                                                                | Type                                                                     | Required                                                                 | Description                                                              |
| ------------------------------------------------------------------------ | ------------------------------------------------------------------------ | ------------------------------------------------------------------------ | ------------------------------------------------------------------------ |
| `request`                                                                | [shared.CreateClientRequest](../../models/shared/createclientrequest.md) | :heavy_check_mark:                                                       | The request object to use for the request.                               |


### Response

**[operations.CreateClientResponse](../../models/operations/createclientresponse.md)**


## create_scope

Create scope

### Example Usage

```python
import sdk
from sdk.models import shared

s = sdk.SDK(
    security=shared.Security(
        authorization="",
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

### Parameters

| Parameter                                                              | Type                                                                   | Required                                                               | Description                                                            |
| ---------------------------------------------------------------------- | ---------------------------------------------------------------------- | ---------------------------------------------------------------------- | ---------------------------------------------------------------------- |
| `request`                                                              | [shared.CreateScopeRequest](../../models/shared/createscoperequest.md) | :heavy_check_mark:                                                     | The request object to use for the request.                             |


### Response

**[operations.CreateScopeResponse](../../models/operations/createscoperesponse.md)**


## create_secret

Add a secret to a client

### Example Usage

```python
import sdk
from sdk.models import operations, shared

s = sdk.SDK(
    security=shared.Security(
        authorization="",
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

### Parameters

| Parameter                                                                        | Type                                                                             | Required                                                                         | Description                                                                      |
| -------------------------------------------------------------------------------- | -------------------------------------------------------------------------------- | -------------------------------------------------------------------------------- | -------------------------------------------------------------------------------- |
| `request`                                                                        | [operations.CreateSecretRequest](../../models/operations/createsecretrequest.md) | :heavy_check_mark:                                                               | The request object to use for the request.                                       |


### Response

**[operations.CreateSecretResponse](../../models/operations/createsecretresponse.md)**


## delete_client

Delete client

### Example Usage

```python
import sdk
from sdk.models import operations, shared

s = sdk.SDK(
    security=shared.Security(
        authorization="",
    ),
)

req = operations.DeleteClientRequest(
    client_id='fugit',
)

res = s.auth.delete_client(req)

if res.status_code == 200:
    # handle response
```

### Parameters

| Parameter                                                                        | Type                                                                             | Required                                                                         | Description                                                                      |
| -------------------------------------------------------------------------------- | -------------------------------------------------------------------------------- | -------------------------------------------------------------------------------- | -------------------------------------------------------------------------------- |
| `request`                                                                        | [operations.DeleteClientRequest](../../models/operations/deleteclientrequest.md) | :heavy_check_mark:                                                               | The request object to use for the request.                                       |


### Response

**[operations.DeleteClientResponse](../../models/operations/deleteclientresponse.md)**


## delete_scope

Delete scope

### Example Usage

```python
import sdk
from sdk.models import operations, shared

s = sdk.SDK(
    security=shared.Security(
        authorization="",
    ),
)

req = operations.DeleteScopeRequest(
    scope_id='deleniti',
)

res = s.auth.delete_scope(req)

if res.status_code == 200:
    # handle response
```

### Parameters

| Parameter                                                                      | Type                                                                           | Required                                                                       | Description                                                                    |
| ------------------------------------------------------------------------------ | ------------------------------------------------------------------------------ | ------------------------------------------------------------------------------ | ------------------------------------------------------------------------------ |
| `request`                                                                      | [operations.DeleteScopeRequest](../../models/operations/deletescoperequest.md) | :heavy_check_mark:                                                             | The request object to use for the request.                                     |


### Response

**[operations.DeleteScopeResponse](../../models/operations/deletescoperesponse.md)**


## delete_scope_from_client

Delete scope from client

### Example Usage

```python
import sdk
from sdk.models import operations, shared

s = sdk.SDK(
    security=shared.Security(
        authorization="",
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

### Parameters

| Parameter                                                                                          | Type                                                                                               | Required                                                                                           | Description                                                                                        |
| -------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------- |
| `request`                                                                                          | [operations.DeleteScopeFromClientRequest](../../models/operations/deletescopefromclientrequest.md) | :heavy_check_mark:                                                                                 | The request object to use for the request.                                                         |


### Response

**[operations.DeleteScopeFromClientResponse](../../models/operations/deletescopefromclientresponse.md)**


## delete_secret

Delete a secret from a client

### Example Usage

```python
import sdk
from sdk.models import operations, shared

s = sdk.SDK(
    security=shared.Security(
        authorization="",
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

### Parameters

| Parameter                                                                        | Type                                                                             | Required                                                                         | Description                                                                      |
| -------------------------------------------------------------------------------- | -------------------------------------------------------------------------------- | -------------------------------------------------------------------------------- | -------------------------------------------------------------------------------- |
| `request`                                                                        | [operations.DeleteSecretRequest](../../models/operations/deletesecretrequest.md) | :heavy_check_mark:                                                               | The request object to use for the request.                                       |


### Response

**[operations.DeleteSecretResponse](../../models/operations/deletesecretresponse.md)**


## delete_transient_scope

Delete a transient scope from a scope

### Example Usage

```python
import sdk
from sdk.models import operations, shared

s = sdk.SDK(
    security=shared.Security(
        authorization="",
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

### Parameters

| Parameter                                                                                        | Type                                                                                             | Required                                                                                         | Description                                                                                      |
| ------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------ |
| `request`                                                                                        | [operations.DeleteTransientScopeRequest](../../models/operations/deletetransientscoperequest.md) | :heavy_check_mark:                                                                               | The request object to use for the request.                                                       |


### Response

**[operations.DeleteTransientScopeResponse](../../models/operations/deletetransientscoperesponse.md)**


## get_server_info

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


res = s.auth.get_server_info()

if res.server_info is not None:
    # handle response
```


### Response

**[operations.GetServerInfoResponse](../../models/operations/getserverinforesponse.md)**


## list_clients

List clients

### Example Usage

```python
import sdk
from sdk.models import shared

s = sdk.SDK(
    security=shared.Security(
        authorization="",
    ),
)


res = s.auth.list_clients()

if res.list_clients_response is not None:
    # handle response
```


### Response

**[operations.ListClientsResponse](../../models/operations/listclientsresponse.md)**


## list_scopes

List Scopes

### Example Usage

```python
import sdk
from sdk.models import shared

s = sdk.SDK(
    security=shared.Security(
        authorization="",
    ),
)


res = s.auth.list_scopes()

if res.list_scopes_response is not None:
    # handle response
```


### Response

**[operations.ListScopesResponse](../../models/operations/listscopesresponse.md)**


## list_users

List users

### Example Usage

```python
import sdk
from sdk.models import shared

s = sdk.SDK(
    security=shared.Security(
        authorization="",
    ),
)


res = s.auth.list_users()

if res.list_users_response is not None:
    # handle response
```


### Response

**[operations.ListUsersResponse](../../models/operations/listusersresponse.md)**


## read_client

Read client

### Example Usage

```python
import sdk
from sdk.models import operations, shared

s = sdk.SDK(
    security=shared.Security(
        authorization="",
    ),
)

req = operations.ReadClientRequest(
    client_id='modi',
)

res = s.auth.read_client(req)

if res.read_client_response is not None:
    # handle response
```

### Parameters

| Parameter                                                                    | Type                                                                         | Required                                                                     | Description                                                                  |
| ---------------------------------------------------------------------------- | ---------------------------------------------------------------------------- | ---------------------------------------------------------------------------- | ---------------------------------------------------------------------------- |
| `request`                                                                    | [operations.ReadClientRequest](../../models/operations/readclientrequest.md) | :heavy_check_mark:                                                           | The request object to use for the request.                                   |


### Response

**[operations.ReadClientResponse](../../models/operations/readclientresponse.md)**


## read_scope

Read scope

### Example Usage

```python
import sdk
from sdk.models import operations, shared

s = sdk.SDK(
    security=shared.Security(
        authorization="",
    ),
)

req = operations.ReadScopeRequest(
    scope_id='qui',
)

res = s.auth.read_scope(req)

if res.read_scope_response is not None:
    # handle response
```

### Parameters

| Parameter                                                                  | Type                                                                       | Required                                                                   | Description                                                                |
| -------------------------------------------------------------------------- | -------------------------------------------------------------------------- | -------------------------------------------------------------------------- | -------------------------------------------------------------------------- |
| `request`                                                                  | [operations.ReadScopeRequest](../../models/operations/readscoperequest.md) | :heavy_check_mark:                                                         | The request object to use for the request.                                 |


### Response

**[operations.ReadScopeResponse](../../models/operations/readscoperesponse.md)**


## read_user

Read user

### Example Usage

```python
import sdk
from sdk.models import operations, shared

s = sdk.SDK(
    security=shared.Security(
        authorization="",
    ),
)

req = operations.ReadUserRequest(
    user_id='impedit',
)

res = s.auth.read_user(req)

if res.read_user_response is not None:
    # handle response
```

### Parameters

| Parameter                                                                | Type                                                                     | Required                                                                 | Description                                                              |
| ------------------------------------------------------------------------ | ------------------------------------------------------------------------ | ------------------------------------------------------------------------ | ------------------------------------------------------------------------ |
| `request`                                                                | [operations.ReadUserRequest](../../models/operations/readuserrequest.md) | :heavy_check_mark:                                                       | The request object to use for the request.                               |


### Response

**[operations.ReadUserResponse](../../models/operations/readuserresponse.md)**


## update_client

Update client

### Example Usage

```python
import sdk
from sdk.models import operations, shared

s = sdk.SDK(
    security=shared.Security(
        authorization="",
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

### Parameters

| Parameter                                                                        | Type                                                                             | Required                                                                         | Description                                                                      |
| -------------------------------------------------------------------------------- | -------------------------------------------------------------------------------- | -------------------------------------------------------------------------------- | -------------------------------------------------------------------------------- |
| `request`                                                                        | [operations.UpdateClientRequest](../../models/operations/updateclientrequest.md) | :heavy_check_mark:                                                               | The request object to use for the request.                                       |


### Response

**[operations.UpdateClientResponse](../../models/operations/updateclientresponse.md)**


## update_scope

Update scope

### Example Usage

```python
import sdk
from sdk.models import operations, shared

s = sdk.SDK(
    security=shared.Security(
        authorization="",
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

### Parameters

| Parameter                                                                      | Type                                                                           | Required                                                                       | Description                                                                    |
| ------------------------------------------------------------------------------ | ------------------------------------------------------------------------------ | ------------------------------------------------------------------------------ | ------------------------------------------------------------------------------ |
| `request`                                                                      | [operations.UpdateScopeRequest](../../models/operations/updatescoperequest.md) | :heavy_check_mark:                                                             | The request object to use for the request.                                     |


### Response

**[operations.UpdateScopeResponse](../../models/operations/updatescoperesponse.md)**

