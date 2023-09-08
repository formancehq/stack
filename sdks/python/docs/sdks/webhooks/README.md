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
from sdk.models import operations, shared

s = sdk.SDK(
    security=shared.Security(
        authorization="",
    ),
)

req = operations.ActivateConfigRequest(
    id='4997257d-dfb6-445b-929c-cbe2ab182818',
)

res = s.webhooks.activate_config(req)

if res.config_response is not None:
    # handle response
```

### Parameters

| Parameter                                                                            | Type                                                                                 | Required                                                                             | Description                                                                          |
| ------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------ |
| `request`                                                                            | [operations.ActivateConfigRequest](../../models/operations/activateconfigrequest.md) | :heavy_check_mark:                                                                   | The request object to use for the request.                                           |


### Response

**[operations.ActivateConfigResponse](../../models/operations/activateconfigresponse.md)**


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
        authorization="",
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

### Parameters

| Parameter                                                                                    | Type                                                                                         | Required                                                                                     | Description                                                                                  |
| -------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------- |
| `request`                                                                                    | [operations.ChangeConfigSecretRequest](../../models/operations/changeconfigsecretrequest.md) | :heavy_check_mark:                                                                           | The request object to use for the request.                                                   |


### Response

**[operations.ChangeConfigSecretResponse](../../models/operations/changeconfigsecretresponse.md)**


## deactivate_config

Deactivate a webhooks config by ID, to stop receiving webhooks to its endpoint.

### Example Usage

```python
import sdk
from sdk.models import operations, shared

s = sdk.SDK(
    security=shared.Security(
        authorization="",
    ),
)

req = operations.DeactivateConfigRequest(
    id='4997257d-dfb6-445b-929c-cbe2ab182818',
)

res = s.webhooks.deactivate_config(req)

if res.config_response is not None:
    # handle response
```

### Parameters

| Parameter                                                                                | Type                                                                                     | Required                                                                                 | Description                                                                              |
| ---------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------- |
| `request`                                                                                | [operations.DeactivateConfigRequest](../../models/operations/deactivateconfigrequest.md) | :heavy_check_mark:                                                                       | The request object to use for the request.                                               |


### Response

**[operations.DeactivateConfigResponse](../../models/operations/deactivateconfigresponse.md)**


## delete_config

Delete a webhooks config by ID.

### Example Usage

```python
import sdk
from sdk.models import operations, shared

s = sdk.SDK(
    security=shared.Security(
        authorization="",
    ),
)

req = operations.DeleteConfigRequest(
    id='4997257d-dfb6-445b-929c-cbe2ab182818',
)

res = s.webhooks.delete_config(req)

if res.status_code == 200:
    # handle response
```

### Parameters

| Parameter                                                                        | Type                                                                             | Required                                                                         | Description                                                                      |
| -------------------------------------------------------------------------------- | -------------------------------------------------------------------------------- | -------------------------------------------------------------------------------- | -------------------------------------------------------------------------------- |
| `request`                                                                        | [operations.DeleteConfigRequest](../../models/operations/deleteconfigrequest.md) | :heavy_check_mark:                                                               | The request object to use for the request.                                       |


### Response

**[operations.DeleteConfigResponse](../../models/operations/deleteconfigresponse.md)**


## get_many_configs

Sorted by updated date descending

### Example Usage

```python
import sdk
from sdk.models import operations, shared

s = sdk.SDK(
    security=shared.Security(
        authorization="",
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

### Parameters

| Parameter                                                                            | Type                                                                                 | Required                                                                             | Description                                                                          |
| ------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------ |
| `request`                                                                            | [operations.GetManyConfigsRequest](../../models/operations/getmanyconfigsrequest.md) | :heavy_check_mark:                                                                   | The request object to use for the request.                                           |


### Response

**[operations.GetManyConfigsResponse](../../models/operations/getmanyconfigsresponse.md)**


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
        authorization="",
    ),
)

req = shared.ConfigUser(
    endpoint='https://example.com',
    event_types=[
        'TYPE1',
        'TYPE1',
        'TYPE1',
        'TYPE1',
    ],
    secret='V0bivxRWveaoz08afqjU6Ko/jwO0Cb+3',
)

res = s.webhooks.insert_config(req)

if res.config_response is not None:
    # handle response
```

### Parameters

| Parameter                                              | Type                                                   | Required                                               | Description                                            |
| ------------------------------------------------------ | ------------------------------------------------------ | ------------------------------------------------------ | ------------------------------------------------------ |
| `request`                                              | [shared.ConfigUser](../../models/shared/configuser.md) | :heavy_check_mark:                                     | The request object to use for the request.             |


### Response

**[operations.InsertConfigResponse](../../models/operations/insertconfigresponse.md)**


## test_config

Test a config by sending a webhook to its endpoint.

### Example Usage

```python
import sdk
from sdk.models import operations, shared

s = sdk.SDK(
    security=shared.Security(
        authorization="",
    ),
)

req = operations.TestConfigRequest(
    id='4997257d-dfb6-445b-929c-cbe2ab182818',
)

res = s.webhooks.test_config(req)

if res.attempt_response is not None:
    # handle response
```

### Parameters

| Parameter                                                                    | Type                                                                         | Required                                                                     | Description                                                                  |
| ---------------------------------------------------------------------------- | ---------------------------------------------------------------------------- | ---------------------------------------------------------------------------- | ---------------------------------------------------------------------------- |
| `request`                                                                    | [operations.TestConfigRequest](../../models/operations/testconfigrequest.md) | :heavy_check_mark:                                                           | The request object to use for the request.                                   |


### Response

**[operations.TestConfigResponse](../../models/operations/testconfigresponse.md)**

