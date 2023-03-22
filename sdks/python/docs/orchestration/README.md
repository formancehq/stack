# orchestration

### Available Operations

* [cancel_event](#cancel_event) - Cancel a running workflow
* [create_workflow](#create_workflow) - Create workflow
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
    instance_id='culpa',
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
    name='Darrin Brakus',
    stages=[
        {
            "repellat": 'mollitia',
        },
        {
            "numquam": 'commodi',
            "quam": 'molestiae',
            "velit": 'error',
        },
        {
            "quis": 'vitae',
        },
    ],
)

res = s.orchestration.create_workflow(req)

if res.create_workflow_response is not None:
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
    instance_id='laborum',
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
    instance_id='animi',
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
    instance_id='enim',
    number=138183,
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
    flow_id='quo',
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
    workflow_id='sequi',
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
        "ipsam": 'id',
        "possimus": 'aut',
        "quasi": 'error',
        "temporibus": 'laborum',
    },
    wait=False,
    workflow_id='quasi',
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
        name='Jan Thiel',
    ),
    instance_id='voluptatibus',
)

res = s.orchestration.send_event(req)

if res.status_code == 200:
    # handle response
```
