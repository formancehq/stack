# orchestration

### Available Operations

* [cancelEvent](#cancelevent) - Cancel a running workflow
* [createWorkflow](#createworkflow) - Create workflow
* [getInstance](#getinstance) - Get a workflow instance by id
* [getInstanceHistory](#getinstancehistory) - Get a workflow instance history by id
* [getInstanceStageHistory](#getinstancestagehistory) - Get a workflow instance stage history
* [getWorkflow](#getworkflow) - Get a flow by id
* [listInstances](#listinstances) - List instances of a workflow
* [listWorkflows](#listworkflows) - List registered workflows
* [orchestrationgetServerInfo](#orchestrationgetserverinfo) - Get server info
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
  instanceID: "quae",
}).then((res: CancelEventResponse) => {
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
  name: "Alison Mann",
  stages: [
    {
      "rem": "voluptates",
      "quasi": "repudiandae",
      "sint": "veritatis",
    },
    {
      "incidunt": "enim",
      "consequatur": "est",
      "quibusdam": "explicabo",
      "deserunt": "distinctio",
    },
  ],
}).then((res: CreateWorkflowResponse) => {
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
  instanceID: "quibusdam",
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
  instanceID: "labore",
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
  instanceID: "modi",
  number: 183191,
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
  flowId: "aliquid",
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
  workflowID: "cupiditate",
}).then((res: ListInstancesResponse) => {
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
    "perferendis": "magni",
    "assumenda": "ipsam",
    "alias": "fugit",
  },
  wait: false,
  workflowID: "dolorum",
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
    name: "Eddie Prosacco",
  },
  instanceID: "delectus",
}).then((res: SendEventResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```
