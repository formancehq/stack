# orchestration

### Available Operations

* [cancel_event](#cancel_event) - Cancel a running workflow
* [create_trigger](#create_trigger) - Create trigger
* [create_workflow](#create_workflow) - Create workflow
* [delete_trigger](#delete_trigger) - Delete trigger
* [delete_workflow](#delete_workflow) - Delete a flow by id
* [get_instance](#get_instance) - Get a workflow instance by id
* [get_instance_history](#get_instance_history) - Get a workflow instance history by id
* [get_instance_stage_history](#get_instance_stage_history) - Get a workflow instance stage history
* [get_workflow](#get_workflow) - Get a flow by id
* [list_instances](#list_instances) - List instances of a workflow
* [list_triggers](#list_triggers) - List triggers
* [list_triggers_occurrences](#list_triggers_occurrences) - List triggers occurrences
* [list_workflows](#list_workflows) - List registered workflows
* [orchestrationget_server_info](#orchestrationget_server_info) - Get server info
* [read_trigger](#read_trigger) - Read trigger
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
    instance_id='ipsam',
)

res = s.orchestration.cancel_event(req)

if res.status_code == 200:
    # handle response
```

## create_trigger

Create trigger

### Example Usage

```python
import sdk
from sdk.models import shared

s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)

req = shared.TriggerData(
    event='id',
    filter='possimus',
    vars={
        "quasi": 'error',
    },
    workflow_id='temporibus',
)

res = s.orchestration.create_trigger(req)

if res.create_trigger_response is not None:
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
    name='Ryan Witting',
    stages=[
        {
            "voluptatibus": 'ipsa',
            "omnis": 'voluptate',
            "cum": 'perferendis',
        },
        {
            "reprehenderit": 'ut',
        },
    ],
)

res = s.orchestration.create_workflow(req)

if res.create_workflow_response is not None:
    # handle response
```

## delete_trigger

Read trigger

### Example Usage

```python
import sdk
from sdk.models import operations

s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)

req = operations.DeleteTriggerRequest(
    trigger_id='maiores',
)

res = s.orchestration.delete_trigger(req)

if res.status_code == 200:
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
    flow_id='dicta',
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
    instance_id='corporis',
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
    instance_id='dolore',
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
    instance_id='iusto',
    number=118727,
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
    flow_id='harum',
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
    workflow_id='enim',
)

res = s.orchestration.list_instances(req)

if res.list_runs_response is not None:
    # handle response
```

## list_triggers

List triggers

### Example Usage

```python
import sdk


s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)


res = s.orchestration.list_triggers()

if res.list_triggers_response is not None:
    # handle response
```

## list_triggers_occurrences

List triggers occurrences

### Example Usage

```python
import sdk
from sdk.models import operations

s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)

req = operations.ListTriggersOccurrencesRequest(
    trigger_id='accusamus',
)

res = s.orchestration.list_triggers_occurrences(req)

if res.list_triggers_occurrences_response is not None:
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

## read_trigger

Read trigger

### Example Usage

```python
import sdk
from sdk.models import operations

s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)

req = operations.ReadTriggerRequest(
    trigger_id='commodi',
)

res = s.orchestration.read_trigger(req)

if res.read_trigger_response is not None:
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
        "quae": 'ipsum',
        "quidem": 'molestias',
        "excepturi": 'pariatur',
        "modi": 'praesentium',
    },
    wait=False,
    workflow_id='rem',
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
        name='Carl Waelchi DVM',
    ),
    instance_id='incidunt',
)

res = s.orchestration.send_event(req)

if res.status_code == 200:
    # handle response
```
