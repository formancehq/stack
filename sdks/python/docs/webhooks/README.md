# webhooks

### Available Operations

* [activate_config](#activate_config) - Activate one config
* [change_config_secret](#change_config_secret) - Change the signing secret of a config
* [deactivate_config](#deactivate_config) - Deactivate one config
* [delete_config](#delete_config) - Delete one config
* [get_many_configs](#get_many_configs) - Get many configs
* [insert_config](#insert_config) - Insert a new config
* [test_config](#test_config) - Test one config

## activate_config

Activate a webhooks config by ID, to start receiving webhooks to its endpoint.

### Example Usage

```python
import sdk
from sdk.models import operations

s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)

req = operations.ActivateConfigRequest(
    id='4997257d-dfb6-445b-929c-cbe2ab182818',
)

res = s.webhooks.activate_config(req)

if res.config_response is not None:
    # handle response
```

## change_config_secret

Change the signing secret of the endpoint of a webhooks config.

If not passed or empty, a secret is automatically generated.
The format is a random string of bytes of size 24, base64 encoded. (larger size after encoding)


### Example Usage

```python
import sdk
from sdk.models import operations, shared

s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)

req = operations.ChangeConfigSecretRequest(
    config_change_secret=shared.ConfigChangeSecret(
        secret='V0bivxRWveaoz08afqjU6Ko/jwO0Cb+3',
    ),
    id='4997257d-dfb6-445b-929c-cbe2ab182818',
)

res = s.webhooks.change_config_secret(req)

if res.config_response is not None:
    # handle response
```

## deactivate_config

Deactivate a webhooks config by ID, to stop receiving webhooks to its endpoint.

### Example Usage

```python
import sdk
from sdk.models import operations

s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)

req = operations.DeactivateConfigRequest(
    id='4997257d-dfb6-445b-929c-cbe2ab182818',
)

res = s.webhooks.deactivate_config(req)

if res.config_response is not None:
    # handle response
```

## delete_config

Delete a webhooks config by ID.

### Example Usage

```python
import sdk
from sdk.models import operations

s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)

req = operations.DeleteConfigRequest(
    id='4997257d-dfb6-445b-929c-cbe2ab182818',
)

res = s.webhooks.delete_config(req)

if res.status_code == 200:
    # handle response
```

## get_many_configs

Sorted by updated date descending

### Example Usage

```python
import sdk
from sdk.models import operations

s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)

req = operations.GetManyConfigsRequest(
    endpoint='https://example.com',
    id='4997257d-dfb6-445b-929c-cbe2ab182818',
)

res = s.webhooks.get_many_configs(req)

if res.configs_response is not None:
    # handle response
```

## insert_config

Insert a new webhooks config.

The endpoint should be a valid https URL and be unique.

The secret is the endpoint's verification secret.
If not passed or empty, a secret is automatically generated.
The format is a random string of bytes of size 24, base64 encoded. (larger size after encoding)

All eventTypes are converted to lower-case when inserted.


### Example Usage

```python
import sdk
from sdk.models import shared

s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)

req = shared.ConfigUser(
    endpoint='https://example.com',
    event_types=[
        'TYPE1',
    ],
    secret='V0bivxRWveaoz08afqjU6Ko/jwO0Cb+3',
)

res = s.webhooks.insert_config(req)

if res.config_response is not None:
    # handle response
```

## test_config

Test a config by sending a webhook to its endpoint.

### Example Usage

```python
import sdk
from sdk.models import operations

s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)

req = operations.TestConfigRequest(
    id='4997257d-dfb6-445b-929c-cbe2ab182818',
)

res = s.webhooks.test_config(req)

if res.attempt_response is not None:
    # handle response
```
