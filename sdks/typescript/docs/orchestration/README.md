# orchestration

### Available Operations

* [cancelEvent](#cancelevent) - Cancel a running workflow
* [createTrigger](#createtrigger) - Create trigger
* [createWorkflow](#createworkflow) - Create workflow
* [deleteTrigger](#deletetrigger) - Delete trigger
* [deleteWorkflow](#deleteworkflow) - Delete a flow by id
* [getInstance](#getinstance) - Get a workflow instance by id
* [getInstanceHistory](#getinstancehistory) - Get a workflow instance history by id
* [getInstanceStageHistory](#getinstancestagehistory) - Get a workflow instance stage history
* [getWorkflow](#getworkflow) - Get a flow by id
* [listInstances](#listinstances) - List instances of a workflow
* [listTriggers](#listtriggers) - List triggers
* [listTriggersOccurrences](#listtriggersoccurrences) - List triggers occurrences
* [listWorkflows](#listworkflows) - List registered workflows
* [orchestrationgetServerInfo](#orchestrationgetserverinfo) - Get server info
* [readTrigger](#readtrigger) - Read trigger
* [runWorkflow](#runworkflow) - Run workflow
* [sendEvent](#sendevent) - Send an event to a running workflow

## cancelEvent

Cancel a running workflow

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { CancelEventResponse } from "@formance/formance-sdk/dist/sdk/models/operations";
import { ErrorErrorCode } from "@formance/formance-sdk/dist/sdk/models/shared";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.orchestration.cancelEvent({
  instanceID: "ipsam",
}).then((res: CancelEventResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## createTrigger

Create trigger

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { CreateTriggerResponse } from "@formance/formance-sdk/dist/sdk/models/operations";
import { ErrorErrorCode } from "@formance/formance-sdk/dist/sdk/models/shared";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.orchestration.createTrigger({
  event: "id",
  filter: "possimus",
  vars: {
    "quasi": "error",
  },
  workflowID: "temporibus",
}).then((res: CreateTriggerResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## createWorkflow

Create a workflow

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { CreateWorkflowResponse } from "@formance/formance-sdk/dist/sdk/models/operations";
import { ErrorErrorCode } from "@formance/formance-sdk/dist/sdk/models/shared";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.orchestration.createWorkflow({
  name: "Ryan Witting",
  stages: [
    {
      "voluptatibus": "ipsa",
      "omnis": "voluptate",
      "cum": "perferendis",
    },
    {
      "reprehenderit": "ut",
    },
  ],
}).then((res: CreateWorkflowResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## deleteTrigger

Read trigger

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { DeleteTriggerResponse } from "@formance/formance-sdk/dist/sdk/models/operations";
import { ErrorErrorCode } from "@formance/formance-sdk/dist/sdk/models/shared";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.orchestration.deleteTrigger({
  triggerID: "maiores",
}).then((res: DeleteTriggerResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## deleteWorkflow

Delete a flow by id

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { DeleteWorkflowResponse } from "@formance/formance-sdk/dist/sdk/models/operations";
import { ErrorErrorCode } from "@formance/formance-sdk/dist/sdk/models/shared";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.orchestration.deleteWorkflow({
  flowId: "dicta",
}).then((res: DeleteWorkflowResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## getInstance

Get a workflow instance by id

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { GetInstanceResponse } from "@formance/formance-sdk/dist/sdk/models/operations";
import { ErrorErrorCode } from "@formance/formance-sdk/dist/sdk/models/shared";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.orchestration.getInstance({
  instanceID: "corporis",
}).then((res: GetInstanceResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## getInstanceHistory

Get a workflow instance history by id

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { GetInstanceHistoryResponse } from "@formance/formance-sdk/dist/sdk/models/operations";
import { ErrorErrorCode } from "@formance/formance-sdk/dist/sdk/models/shared";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.orchestration.getInstanceHistory({
  instanceID: "dolore",
}).then((res: GetInstanceHistoryResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## getInstanceStageHistory

Get a workflow instance stage history

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { GetInstanceStageHistoryResponse } from "@formance/formance-sdk/dist/sdk/models/operations";
import { Connector, ErrorErrorCode, PaymentScheme, PaymentStatus, PaymentType } from "@formance/formance-sdk/dist/sdk/models/shared";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.orchestration.getInstanceStageHistory({
  instanceID: "iusto",
  number: 118727,
}).then((res: GetInstanceStageHistoryResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## getWorkflow

Get a flow by id

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { GetWorkflowResponse } from "@formance/formance-sdk/dist/sdk/models/operations";
import { ErrorErrorCode } from "@formance/formance-sdk/dist/sdk/models/shared";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.orchestration.getWorkflow({
  flowId: "harum",
}).then((res: GetWorkflowResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## listInstances

List instances of a workflow

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { ListInstancesResponse } from "@formance/formance-sdk/dist/sdk/models/operations";
import { ErrorErrorCode } from "@formance/formance-sdk/dist/sdk/models/shared";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.orchestration.listInstances({
  running: false,
  workflowID: "enim",
}).then((res: ListInstancesResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## listTriggers

List triggers

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { ListTriggersResponse } from "@formance/formance-sdk/dist/sdk/models/operations";
import { ErrorErrorCode } from "@formance/formance-sdk/dist/sdk/models/shared";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.orchestration.listTriggers().then((res: ListTriggersResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## listTriggersOccurrences

List triggers occurrences

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { ListTriggersOccurrencesResponse } from "@formance/formance-sdk/dist/sdk/models/operations";
import { ErrorErrorCode } from "@formance/formance-sdk/dist/sdk/models/shared";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.orchestration.listTriggersOccurrences({
  triggerID: "accusamus",
}).then((res: ListTriggersOccurrencesResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## listWorkflows

List registered workflows

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { ListWorkflowsResponse } from "@formance/formance-sdk/dist/sdk/models/operations";
import { ErrorErrorCode } from "@formance/formance-sdk/dist/sdk/models/shared";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.orchestration.listWorkflows().then((res: ListWorkflowsResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## orchestrationgetServerInfo

Get server info

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { OrchestrationgetServerInfoResponse } from "@formance/formance-sdk/dist/sdk/models/operations";
import { ErrorErrorCode } from "@formance/formance-sdk/dist/sdk/models/shared";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.orchestration.orchestrationgetServerInfo().then((res: OrchestrationgetServerInfoResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## readTrigger

Read trigger

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { ReadTriggerResponse } from "@formance/formance-sdk/dist/sdk/models/operations";
import { ErrorErrorCode } from "@formance/formance-sdk/dist/sdk/models/shared";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.orchestration.readTrigger({
  triggerID: "commodi",
}).then((res: ReadTriggerResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## runWorkflow

Run workflow

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { RunWorkflowResponse } from "@formance/formance-sdk/dist/sdk/models/operations";
import { ErrorErrorCode } from "@formance/formance-sdk/dist/sdk/models/shared";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.orchestration.runWorkflow({
  requestBody: {
    "quae": "ipsum",
    "quidem": "molestias",
    "excepturi": "pariatur",
    "modi": "praesentium",
  },
  wait: false,
  workflowID: "rem",
}).then((res: RunWorkflowResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## sendEvent

Send an event to a running workflow

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { SendEventResponse } from "@formance/formance-sdk/dist/sdk/models/operations";
import { ErrorErrorCode } from "@formance/formance-sdk/dist/sdk/models/shared";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.orchestration.sendEvent({
  requestBody: {
    name: "Carl Waelchi DVM",
  },
  instanceID: "incidunt",
}).then((res: SendEventResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```
