# orchestration

### Available Operations

* [cancel_event](#cancel_event) - Cancel a running workflow
* [create_workflow](#create_workflow) - Create workflow
* [delete_workflow](#delete_workflow) - Delete a flow by id
* [get_instance](#get_instance) - Get a workflow instance by id
* [get_instance_history](#get_instance_history) - Get a workflow instance history by id
* [get_instance_stage_history](#get_instance_stage_history) - Get a workflow instance stage history
* [get_workflow](#get_workflow) - Get a flow by id
* [list_instances](#list_instances) - List instances of a workflow
* [list_workflows](#list_workflows) - List registered workflows
* [orchestrationget_server_info](#orchestrationget_server_info) - Get server info
* [run_workflow](#run_workflow) - Run workflow
* [send_event](#send_event) - Send an event to a running workflow

## cancel_event

Cancel a running workflow

### Example Usage

```python
import sdk
from sdk.models import operations

s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)

req = operations.CancelEventRequest(
    instance_id='quae',
)

res = s.orchestration.cancel_event(req)

if res.status_code == 200:
    # handle response
```

## create_workflow

Create a workflow

### Example Usage

```python
import sdk
from sdk.models import shared

s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)

req = shared.CreateWorkflowRequest(
    name='Alison Mann',
    stages=[
        {
            "rem": 'voluptates',
            "quasi": 'repudiandae',
            "sint": 'veritatis',
        },
        {
            "incidunt": 'enim',
            "consequatur": 'est',
            "quibusdam": 'explicabo',
            "deserunt": 'distinctio',
        },
    ],
)

res = s.orchestration.create_workflow(req)

if res.create_workflow_response is not None:
    # handle response
```

## delete_workflow

Delete a flow by id

### Example Usage

```python
import sdk
from sdk.models import operations

s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)

req = operations.DeleteWorkflowRequest(
    flow_id='quibusdam',
)

res = s.orchestration.delete_workflow(req)

if res.status_code == 200:
    # handle response
```

## get_instance

Get a workflow instance by id

### Example Usage

```python
import sdk
from sdk.models import operations

s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)

req = operations.GetInstanceRequest(
    instance_id='labore',
)

res = s.orchestration.get_instance(req)

if res.get_workflow_instance_response is not None:
    # handle response
```

## get_instance_history

Get a workflow instance history by id

### Example Usage

```python
import sdk
from sdk.models import operations

s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)

req = operations.GetInstanceHistoryRequest(
    instance_id='modi',
)

res = s.orchestration.get_instance_history(req)

if res.get_workflow_instance_history_response is not None:
    # handle response
```

## get_instance_stage_history

Get a workflow instance stage history

### Example Usage

```python
import sdk
from sdk.models import operations

s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)

req = operations.GetInstanceStageHistoryRequest(
    instance_id='qui',
    number=397821,
)

res = s.orchestration.get_instance_stage_history(req)

if res.get_workflow_instance_history_stage_response is not None:
    # handle response
```

## get_workflow

Get a flow by id

### Example Usage

```python
import sdk
from sdk.models import operations

s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)

req = operations.GetWorkflowRequest(
    flow_id='cupiditate',
)

res = s.orchestration.get_workflow(req)

if res.get_workflow_response is not None:
    # handle response
```

## list_instances

List instances of a workflow

### Example Usage

```python
import sdk
from sdk.models import operations

s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)

req = operations.ListInstancesRequest(
    running=False,
    workflow_id='quos',
)

res = s.orchestration.list_instances(req)

if res.list_runs_response is not None:
    # handle response
```

## list_workflows

List registered workflows

### Example Usage

```python
import sdk


s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)


res = s.orchestration.list_workflows()

if res.list_workflows_response is not None:
    # handle response
```

## orchestrationget_server_info

Get server info

### Example Usage

```python
import sdk


s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)


res = s.orchestration.orchestrationget_server_info()

if res.server_info is not None:
    # handle response
```

## run_workflow

Run workflow

### Example Usage

```python
import sdk
from sdk.models import operations

s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)

req = operations.RunWorkflowRequest(
    request_body={
        "magni": 'assumenda',
    },
    wait=False,
    workflow_id='ipsam',
)

res = s.orchestration.run_workflow(req)

if res.run_workflow_response is not None:
    # handle response
```

## send_event

Send an event to a running workflow

### Example Usage

```python
import sdk
from sdk.models import operations

s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)

req = operations.SendEventRequest(
    request_body=operations.SendEventRequestBody(
        name='Denise Pagac',
    ),
    instance_id='facilis',
)

res = s.orchestration.send_event(req)

if res.status_code == 200:
    # handle response
```
