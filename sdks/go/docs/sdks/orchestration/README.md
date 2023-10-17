# Orchestration

### Available Operations

* [CancelEvent](#cancelevent) - Cancel a running workflow
* [CreateWorkflow](#createworkflow) - Create workflow
* [DeleteWorkflow](#deleteworkflow) - Delete a flow by id
* [GetInstance](#getinstance) - Get a workflow instance by id
* [GetInstanceHistory](#getinstancehistory) - Get a workflow instance history by id
* [GetInstanceStageHistory](#getinstancestagehistory) - Get a workflow instance stage history
* [GetWorkflow](#getworkflow) - Get a flow by id
* [ListInstances](#listinstances) - List instances of a workflow
* [ListWorkflows](#listworkflows) - List registered workflows
* [OrchestrationgetServerInfo](#orchestrationgetserverinfo) - Get server info
* [RunWorkflow](#runworkflow) - Run workflow
* [SendEvent](#sendevent) - Send an event to a running workflow

## CancelEvent

Cancel a running workflow

### Example Usage

```go
package main

import(
	"context"
	"log"
	"github.com/formancehq/formance-sdk-go"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
)

func main() {
    s := formance.New()

    ctx := context.Background()
    res, err := s.Orchestration.CancelEvent(ctx, operations.CancelEventRequest{
        InstanceID: "est",
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.StatusCode == http.StatusOK {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                      | Type                                                                           | Required                                                                       | Description                                                                    |
| ------------------------------------------------------------------------------ | ------------------------------------------------------------------------------ | ------------------------------------------------------------------------------ | ------------------------------------------------------------------------------ |
| `ctx`                                                                          | [context.Context](https://pkg.go.dev/context#Context)                          | :heavy_check_mark:                                                             | The context to use for the request.                                            |
| `request`                                                                      | [operations.CancelEventRequest](../../models/operations/canceleventrequest.md) | :heavy_check_mark:                                                             | The request object to use for the request.                                     |


### Response

**[*operations.CancelEventResponse](../../models/operations/canceleventresponse.md), error**


## CreateWorkflow

Create a workflow

### Example Usage

```go
package main

import(
	"context"
	"log"
	"github.com/formancehq/formance-sdk-go"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
)

func main() {
    s := formance.New()

    ctx := context.Background()
    res, err := s.Orchestration.CreateWorkflow(ctx, shared.CreateWorkflowRequest{
        Name: formance.String("Benjamin O'Connell"),
        Stages: []map[string]interface{}{
            map[string]interface{}{
                "qui": "aliquid",
                "cupiditate": "quos",
            },
            map[string]interface{}{
                "magni": "assumenda",
            },
        },
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.CreateWorkflowResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                    | Type                                                                         | Required                                                                     | Description                                                                  |
| ---------------------------------------------------------------------------- | ---------------------------------------------------------------------------- | ---------------------------------------------------------------------------- | ---------------------------------------------------------------------------- |
| `ctx`                                                                        | [context.Context](https://pkg.go.dev/context#Context)                        | :heavy_check_mark:                                                           | The context to use for the request.                                          |
| `request`                                                                    | [shared.CreateWorkflowRequest](../../models/shared/createworkflowrequest.md) | :heavy_check_mark:                                                           | The request object to use for the request.                                   |


### Response

**[*operations.CreateWorkflowResponse](../../models/operations/createworkflowresponse.md), error**


## DeleteWorkflow

Delete a flow by id

### Example Usage

```go
package main

import(
	"context"
	"log"
	"github.com/formancehq/formance-sdk-go"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
)

func main() {
    s := formance.New()

    ctx := context.Background()
    res, err := s.Orchestration.DeleteWorkflow(ctx, operations.DeleteWorkflowRequest{
        FlowID: "ipsam",
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.StatusCode == http.StatusOK {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                            | Type                                                                                 | Required                                                                             | Description                                                                          |
| ------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------ |
| `ctx`                                                                                | [context.Context](https://pkg.go.dev/context#Context)                                | :heavy_check_mark:                                                                   | The context to use for the request.                                                  |
| `request`                                                                            | [operations.DeleteWorkflowRequest](../../models/operations/deleteworkflowrequest.md) | :heavy_check_mark:                                                                   | The request object to use for the request.                                           |


### Response

**[*operations.DeleteWorkflowResponse](../../models/operations/deleteworkflowresponse.md), error**


## GetInstance

Get a workflow instance by id

### Example Usage

```go
package main

import(
	"context"
	"log"
	"github.com/formancehq/formance-sdk-go"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
)

func main() {
    s := formance.New()

    ctx := context.Background()
    res, err := s.Orchestration.GetInstance(ctx, operations.GetInstanceRequest{
        InstanceID: "alias",
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.GetWorkflowInstanceResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                      | Type                                                                           | Required                                                                       | Description                                                                    |
| ------------------------------------------------------------------------------ | ------------------------------------------------------------------------------ | ------------------------------------------------------------------------------ | ------------------------------------------------------------------------------ |
| `ctx`                                                                          | [context.Context](https://pkg.go.dev/context#Context)                          | :heavy_check_mark:                                                             | The context to use for the request.                                            |
| `request`                                                                      | [operations.GetInstanceRequest](../../models/operations/getinstancerequest.md) | :heavy_check_mark:                                                             | The request object to use for the request.                                     |


### Response

**[*operations.GetInstanceResponse](../../models/operations/getinstanceresponse.md), error**


## GetInstanceHistory

Get a workflow instance history by id

### Example Usage

```go
package main

import(
	"context"
	"log"
	"github.com/formancehq/formance-sdk-go"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
)

func main() {
    s := formance.New()

    ctx := context.Background()
    res, err := s.Orchestration.GetInstanceHistory(ctx, operations.GetInstanceHistoryRequest{
        InstanceID: "fugit",
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.GetWorkflowInstanceHistoryResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                                    | Type                                                                                         | Required                                                                                     | Description                                                                                  |
| -------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------- |
| `ctx`                                                                                        | [context.Context](https://pkg.go.dev/context#Context)                                        | :heavy_check_mark:                                                                           | The context to use for the request.                                                          |
| `request`                                                                                    | [operations.GetInstanceHistoryRequest](../../models/operations/getinstancehistoryrequest.md) | :heavy_check_mark:                                                                           | The request object to use for the request.                                                   |


### Response

**[*operations.GetInstanceHistoryResponse](../../models/operations/getinstancehistoryresponse.md), error**


## GetInstanceStageHistory

Get a workflow instance stage history

### Example Usage

```go
package main

import(
	"context"
	"log"
	"github.com/formancehq/formance-sdk-go"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
)

func main() {
    s := formance.New()

    ctx := context.Background()
    res, err := s.Orchestration.GetInstanceStageHistory(ctx, operations.GetInstanceStageHistoryRequest{
        InstanceID: "dolorum",
        Number: 569618,
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.GetWorkflowInstanceHistoryStageResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                                              | Type                                                                                                   | Required                                                                                               | Description                                                                                            |
| ------------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------------ |
| `ctx`                                                                                                  | [context.Context](https://pkg.go.dev/context#Context)                                                  | :heavy_check_mark:                                                                                     | The context to use for the request.                                                                    |
| `request`                                                                                              | [operations.GetInstanceStageHistoryRequest](../../models/operations/getinstancestagehistoryrequest.md) | :heavy_check_mark:                                                                                     | The request object to use for the request.                                                             |


### Response

**[*operations.GetInstanceStageHistoryResponse](../../models/operations/getinstancestagehistoryresponse.md), error**


## GetWorkflow

Get a flow by id

### Example Usage

```go
package main

import(
	"context"
	"log"
	"github.com/formancehq/formance-sdk-go"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
)

func main() {
    s := formance.New()

    ctx := context.Background()
    res, err := s.Orchestration.GetWorkflow(ctx, operations.GetWorkflowRequest{
        FlowID: "tempora",
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.GetWorkflowResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                      | Type                                                                           | Required                                                                       | Description                                                                    |
| ------------------------------------------------------------------------------ | ------------------------------------------------------------------------------ | ------------------------------------------------------------------------------ | ------------------------------------------------------------------------------ |
| `ctx`                                                                          | [context.Context](https://pkg.go.dev/context#Context)                          | :heavy_check_mark:                                                             | The context to use for the request.                                            |
| `request`                                                                      | [operations.GetWorkflowRequest](../../models/operations/getworkflowrequest.md) | :heavy_check_mark:                                                             | The request object to use for the request.                                     |


### Response

**[*operations.GetWorkflowResponse](../../models/operations/getworkflowresponse.md), error**


## ListInstances

List instances of a workflow

### Example Usage

```go
package main

import(
	"context"
	"log"
	"github.com/formancehq/formance-sdk-go"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
)

func main() {
    s := formance.New()

    ctx := context.Background()
    res, err := s.Orchestration.ListInstances(ctx, operations.ListInstancesRequest{
        Running: formance.Bool(false),
        WorkflowID: formance.String("facilis"),
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.ListRunsResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                          | Type                                                                               | Required                                                                           | Description                                                                        |
| ---------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------- |
| `ctx`                                                                              | [context.Context](https://pkg.go.dev/context#Context)                              | :heavy_check_mark:                                                                 | The context to use for the request.                                                |
| `request`                                                                          | [operations.ListInstancesRequest](../../models/operations/listinstancesrequest.md) | :heavy_check_mark:                                                                 | The request object to use for the request.                                         |


### Response

**[*operations.ListInstancesResponse](../../models/operations/listinstancesresponse.md), error**


## ListWorkflows

List registered workflows

### Example Usage

```go
package main

import(
	"context"
	"log"
	"github.com/formancehq/formance-sdk-go"
)

func main() {
    s := formance.New()

    ctx := context.Background()
    res, err := s.Orchestration.ListWorkflows(ctx)
    if err != nil {
        log.Fatal(err)
    }

    if res.ListWorkflowsResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                             | Type                                                  | Required                                              | Description                                           |
| ----------------------------------------------------- | ----------------------------------------------------- | ----------------------------------------------------- | ----------------------------------------------------- |
| `ctx`                                                 | [context.Context](https://pkg.go.dev/context#Context) | :heavy_check_mark:                                    | The context to use for the request.                   |


### Response

**[*operations.ListWorkflowsResponse](../../models/operations/listworkflowsresponse.md), error**


## OrchestrationgetServerInfo

Get server info

### Example Usage

```go
package main

import(
	"context"
	"log"
	"github.com/formancehq/formance-sdk-go"
)

func main() {
    s := formance.New()

    ctx := context.Background()
    res, err := s.Orchestration.OrchestrationgetServerInfo(ctx)
    if err != nil {
        log.Fatal(err)
    }

    if res.ServerInfo != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                             | Type                                                  | Required                                              | Description                                           |
| ----------------------------------------------------- | ----------------------------------------------------- | ----------------------------------------------------- | ----------------------------------------------------- |
| `ctx`                                                 | [context.Context](https://pkg.go.dev/context#Context) | :heavy_check_mark:                                    | The context to use for the request.                   |


### Response

**[*operations.OrchestrationgetServerInfoResponse](../../models/operations/orchestrationgetserverinforesponse.md), error**


## RunWorkflow

Run workflow

### Example Usage

```go
package main

import(
	"context"
	"log"
	"github.com/formancehq/formance-sdk-go"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
)

func main() {
    s := formance.New()

    ctx := context.Background()
    res, err := s.Orchestration.RunWorkflow(ctx, operations.RunWorkflowRequest{
        RequestBody: map[string]string{
            "labore": "delectus",
            "eum": "non",
            "eligendi": "sint",
        },
        Wait: formance.Bool(false),
        WorkflowID: "aliquid",
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.RunWorkflowResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                      | Type                                                                           | Required                                                                       | Description                                                                    |
| ------------------------------------------------------------------------------ | ------------------------------------------------------------------------------ | ------------------------------------------------------------------------------ | ------------------------------------------------------------------------------ |
| `ctx`                                                                          | [context.Context](https://pkg.go.dev/context#Context)                          | :heavy_check_mark:                                                             | The context to use for the request.                                            |
| `request`                                                                      | [operations.RunWorkflowRequest](../../models/operations/runworkflowrequest.md) | :heavy_check_mark:                                                             | The request object to use for the request.                                     |


### Response

**[*operations.RunWorkflowResponse](../../models/operations/runworkflowresponse.md), error**


## SendEvent

Send an event to a running workflow

### Example Usage

```go
package main

import(
	"context"
	"log"
	"github.com/formancehq/formance-sdk-go"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
)

func main() {
    s := formance.New()

    ctx := context.Background()
    res, err := s.Orchestration.SendEvent(ctx, operations.SendEventRequest{
        RequestBody: &operations.SendEventRequestBody{
            Name: "Terence Marquardt",
        },
        InstanceID: "debitis",
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.StatusCode == http.StatusOK {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                  | Type                                                                       | Required                                                                   | Description                                                                |
| -------------------------------------------------------------------------- | -------------------------------------------------------------------------- | -------------------------------------------------------------------------- | -------------------------------------------------------------------------- |
| `ctx`                                                                      | [context.Context](https://pkg.go.dev/context#Context)                      | :heavy_check_mark:                                                         | The context to use for the request.                                        |
| `request`                                                                  | [operations.SendEventRequest](../../models/operations/sendeventrequest.md) | :heavy_check_mark:                                                         | The request object to use for the request.                                 |


### Response

**[*operations.SendEventResponse](../../models/operations/sendeventresponse.md), error**

