# Orchestration
(*Orchestration*)

### Available Operations

* [CancelEvent](#cancelevent) - Cancel a running workflow
* [CreateTrigger](#createtrigger) - Create trigger
* [CreateWorkflow](#createworkflow) - Create workflow
* [DeleteTrigger](#deletetrigger) - Delete trigger
* [DeleteWorkflow](#deleteworkflow) - Delete a flow by id
* [GetInstance](#getinstance) - Get a workflow instance by id
* [GetInstanceHistory](#getinstancehistory) - Get a workflow instance history by id
* [GetInstanceStageHistory](#getinstancestagehistory) - Get a workflow instance stage history
* [GetWorkflow](#getworkflow) - Get a flow by id
* [ListInstances](#listinstances) - List instances of a workflow
* [ListTriggers](#listtriggers) - List triggers
* [ListTriggersOccurrences](#listtriggersoccurrences) - List triggers occurrences
* [ListWorkflows](#listworkflows) - List registered workflows
* [OrchestrationgetServerInfo](#orchestrationgetserverinfo) - Get server info
* [ReadTrigger](#readtrigger) - Read trigger
* [RunWorkflow](#runworkflow) - Run workflow
* [SendEvent](#sendevent) - Send an event to a running workflow
* [TestTrigger](#testtrigger) - Test trigger
* [V2CancelEvent](#v2cancelevent) - Cancel a running workflow
* [V2CreateTrigger](#v2createtrigger) - Create trigger
* [V2CreateWorkflow](#v2createworkflow) - Create workflow
* [V2DeleteTrigger](#v2deletetrigger) - Delete trigger
* [V2DeleteWorkflow](#v2deleteworkflow) - Delete a flow by id
* [V2GetInstance](#v2getinstance) - Get a workflow instance by id
* [V2GetInstanceHistory](#v2getinstancehistory) - Get a workflow instance history by id
* [V2GetInstanceStageHistory](#v2getinstancestagehistory) - Get a workflow instance stage history
* [V2GetServerInfo](#v2getserverinfo) - Get server info
* [V2GetWorkflow](#v2getworkflow) - Get a flow by id
* [V2ListInstances](#v2listinstances) - List instances of a workflow
* [V2ListTriggers](#v2listtriggers) - List triggers
* [V2ListTriggersOccurrences](#v2listtriggersoccurrences) - List triggers occurrences
* [V2ListWorkflows](#v2listworkflows) - List registered workflows
* [V2ReadTrigger](#v2readtrigger) - Read trigger
* [V2RunWorkflow](#v2runworkflow) - Run workflow
* [V2SendEvent](#v2sendevent) - Send an event to a running workflow

## CancelEvent

Cancel a running workflow

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	formancesdkgo "github.com/formancehq/formance-sdk-go"
	"context"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"log"
	"net/http"
)

func main() {
    s := formancesdkgo.New(
        formancesdkgo.WithSecurity(shared.Security{
            Authorization: "Bearer <YOUR_ACCESS_TOKEN_HERE>",
        }),
    )

    ctx := context.Background()
    res, err := s.Orchestration.CancelEvent(ctx, operations.CancelEventRequest{
        InstanceID: "string",
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

| Parameter                                                                          | Type                                                                               | Required                                                                           | Description                                                                        |
| ---------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------- |
| `ctx`                                                                              | [context.Context](https://pkg.go.dev/context#Context)                              | :heavy_check_mark:                                                                 | The context to use for the request.                                                |
| `request`                                                                          | [operations.CancelEventRequest](../../pkg/models/operations/canceleventrequest.md) | :heavy_check_mark:                                                                 | The request object to use for the request.                                         |


### Response

**[*operations.CancelEventResponse](../../pkg/models/operations/canceleventresponse.md), error**
| Error Object       | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4xx-5xx            | */*                |

## CreateTrigger

Create trigger

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	formancesdkgo "github.com/formancehq/formance-sdk-go"
	"context"
	"log"
)

func main() {
    s := formancesdkgo.New(
        formancesdkgo.WithSecurity(shared.Security{
            Authorization: "Bearer <YOUR_ACCESS_TOKEN_HERE>",
        }),
    )

    ctx := context.Background()
    res, err := s.Orchestration.CreateTrigger(ctx, &shared.TriggerData{
        Event: "string",
        Vars: map[string]interface{}{
            "key": "string",
        },
        WorkflowID: "string",
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.CreateTriggerResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                    | Type                                                         | Required                                                     | Description                                                  |
| ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ |
| `ctx`                                                        | [context.Context](https://pkg.go.dev/context#Context)        | :heavy_check_mark:                                           | The context to use for the request.                          |
| `request`                                                    | [shared.TriggerData](../../pkg/models/shared/triggerdata.md) | :heavy_check_mark:                                           | The request object to use for the request.                   |


### Response

**[*operations.CreateTriggerResponse](../../pkg/models/operations/createtriggerresponse.md), error**
| Error Object       | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4xx-5xx            | */*                |

## CreateWorkflow

Create a workflow

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	formancesdkgo "github.com/formancehq/formance-sdk-go"
	"context"
	"log"
)

func main() {
    s := formancesdkgo.New(
        formancesdkgo.WithSecurity(shared.Security{
            Authorization: "Bearer <YOUR_ACCESS_TOKEN_HERE>",
        }),
    )

    ctx := context.Background()
    res, err := s.Orchestration.CreateWorkflow(ctx, &shared.CreateWorkflowRequest{
        Stages: []map[string]interface{}{
            map[string]interface{}{
                "key": "string",
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

| Parameter                                                                        | Type                                                                             | Required                                                                         | Description                                                                      |
| -------------------------------------------------------------------------------- | -------------------------------------------------------------------------------- | -------------------------------------------------------------------------------- | -------------------------------------------------------------------------------- |
| `ctx`                                                                            | [context.Context](https://pkg.go.dev/context#Context)                            | :heavy_check_mark:                                                               | The context to use for the request.                                              |
| `request`                                                                        | [shared.CreateWorkflowRequest](../../pkg/models/shared/createworkflowrequest.md) | :heavy_check_mark:                                                               | The request object to use for the request.                                       |


### Response

**[*operations.CreateWorkflowResponse](../../pkg/models/operations/createworkflowresponse.md), error**
| Error Object       | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4xx-5xx            | */*                |

## DeleteTrigger

Read trigger

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	formancesdkgo "github.com/formancehq/formance-sdk-go"
	"context"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"log"
	"net/http"
)

func main() {
    s := formancesdkgo.New(
        formancesdkgo.WithSecurity(shared.Security{
            Authorization: "Bearer <YOUR_ACCESS_TOKEN_HERE>",
        }),
    )

    ctx := context.Background()
    res, err := s.Orchestration.DeleteTrigger(ctx, operations.DeleteTriggerRequest{
        TriggerID: "string",
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

| Parameter                                                                              | Type                                                                                   | Required                                                                               | Description                                                                            |
| -------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------- |
| `ctx`                                                                                  | [context.Context](https://pkg.go.dev/context#Context)                                  | :heavy_check_mark:                                                                     | The context to use for the request.                                                    |
| `request`                                                                              | [operations.DeleteTriggerRequest](../../pkg/models/operations/deletetriggerrequest.md) | :heavy_check_mark:                                                                     | The request object to use for the request.                                             |


### Response

**[*operations.DeleteTriggerResponse](../../pkg/models/operations/deletetriggerresponse.md), error**
| Error Object       | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4xx-5xx            | */*                |

## DeleteWorkflow

Delete a flow by id

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	formancesdkgo "github.com/formancehq/formance-sdk-go"
	"context"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"log"
	"net/http"
)

func main() {
    s := formancesdkgo.New(
        formancesdkgo.WithSecurity(shared.Security{
            Authorization: "Bearer <YOUR_ACCESS_TOKEN_HERE>",
        }),
    )

    ctx := context.Background()
    res, err := s.Orchestration.DeleteWorkflow(ctx, operations.DeleteWorkflowRequest{
        FlowID: "string",
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

| Parameter                                                                                | Type                                                                                     | Required                                                                                 | Description                                                                              |
| ---------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------- |
| `ctx`                                                                                    | [context.Context](https://pkg.go.dev/context#Context)                                    | :heavy_check_mark:                                                                       | The context to use for the request.                                                      |
| `request`                                                                                | [operations.DeleteWorkflowRequest](../../pkg/models/operations/deleteworkflowrequest.md) | :heavy_check_mark:                                                                       | The request object to use for the request.                                               |


### Response

**[*operations.DeleteWorkflowResponse](../../pkg/models/operations/deleteworkflowresponse.md), error**
| Error Object       | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4xx-5xx            | */*                |

## GetInstance

Get a workflow instance by id

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	formancesdkgo "github.com/formancehq/formance-sdk-go"
	"context"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"log"
)

func main() {
    s := formancesdkgo.New(
        formancesdkgo.WithSecurity(shared.Security{
            Authorization: "Bearer <YOUR_ACCESS_TOKEN_HERE>",
        }),
    )

    ctx := context.Background()
    res, err := s.Orchestration.GetInstance(ctx, operations.GetInstanceRequest{
        InstanceID: "string",
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

| Parameter                                                                          | Type                                                                               | Required                                                                           | Description                                                                        |
| ---------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------- |
| `ctx`                                                                              | [context.Context](https://pkg.go.dev/context#Context)                              | :heavy_check_mark:                                                                 | The context to use for the request.                                                |
| `request`                                                                          | [operations.GetInstanceRequest](../../pkg/models/operations/getinstancerequest.md) | :heavy_check_mark:                                                                 | The request object to use for the request.                                         |


### Response

**[*operations.GetInstanceResponse](../../pkg/models/operations/getinstanceresponse.md), error**
| Error Object       | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4xx-5xx            | */*                |

## GetInstanceHistory

Get a workflow instance history by id

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	formancesdkgo "github.com/formancehq/formance-sdk-go"
	"context"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"log"
)

func main() {
    s := formancesdkgo.New(
        formancesdkgo.WithSecurity(shared.Security{
            Authorization: "Bearer <YOUR_ACCESS_TOKEN_HERE>",
        }),
    )

    ctx := context.Background()
    res, err := s.Orchestration.GetInstanceHistory(ctx, operations.GetInstanceHistoryRequest{
        InstanceID: "string",
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

| Parameter                                                                                        | Type                                                                                             | Required                                                                                         | Description                                                                                      |
| ------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------ |
| `ctx`                                                                                            | [context.Context](https://pkg.go.dev/context#Context)                                            | :heavy_check_mark:                                                                               | The context to use for the request.                                                              |
| `request`                                                                                        | [operations.GetInstanceHistoryRequest](../../pkg/models/operations/getinstancehistoryrequest.md) | :heavy_check_mark:                                                                               | The request object to use for the request.                                                       |


### Response

**[*operations.GetInstanceHistoryResponse](../../pkg/models/operations/getinstancehistoryresponse.md), error**
| Error Object       | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4xx-5xx            | */*                |

## GetInstanceStageHistory

Get a workflow instance stage history

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	formancesdkgo "github.com/formancehq/formance-sdk-go"
	"context"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"log"
)

func main() {
    s := formancesdkgo.New(
        formancesdkgo.WithSecurity(shared.Security{
            Authorization: "Bearer <YOUR_ACCESS_TOKEN_HERE>",
        }),
    )

    ctx := context.Background()
    res, err := s.Orchestration.GetInstanceStageHistory(ctx, operations.GetInstanceStageHistoryRequest{
        InstanceID: "string",
        Number: 600636,
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

| Parameter                                                                                                  | Type                                                                                                       | Required                                                                                                   | Description                                                                                                |
| ---------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------- |
| `ctx`                                                                                                      | [context.Context](https://pkg.go.dev/context#Context)                                                      | :heavy_check_mark:                                                                                         | The context to use for the request.                                                                        |
| `request`                                                                                                  | [operations.GetInstanceStageHistoryRequest](../../pkg/models/operations/getinstancestagehistoryrequest.md) | :heavy_check_mark:                                                                                         | The request object to use for the request.                                                                 |


### Response

**[*operations.GetInstanceStageHistoryResponse](../../pkg/models/operations/getinstancestagehistoryresponse.md), error**
| Error Object       | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4xx-5xx            | */*                |

## GetWorkflow

Get a flow by id

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	formancesdkgo "github.com/formancehq/formance-sdk-go"
	"context"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"log"
)

func main() {
    s := formancesdkgo.New(
        formancesdkgo.WithSecurity(shared.Security{
            Authorization: "Bearer <YOUR_ACCESS_TOKEN_HERE>",
        }),
    )

    ctx := context.Background()
    res, err := s.Orchestration.GetWorkflow(ctx, operations.GetWorkflowRequest{
        FlowID: "string",
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

| Parameter                                                                          | Type                                                                               | Required                                                                           | Description                                                                        |
| ---------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------- |
| `ctx`                                                                              | [context.Context](https://pkg.go.dev/context#Context)                              | :heavy_check_mark:                                                                 | The context to use for the request.                                                |
| `request`                                                                          | [operations.GetWorkflowRequest](../../pkg/models/operations/getworkflowrequest.md) | :heavy_check_mark:                                                                 | The request object to use for the request.                                         |


### Response

**[*operations.GetWorkflowResponse](../../pkg/models/operations/getworkflowresponse.md), error**
| Error Object       | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4xx-5xx            | */*                |

## ListInstances

List instances of a workflow

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	formancesdkgo "github.com/formancehq/formance-sdk-go"
	"context"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"log"
)

func main() {
    s := formancesdkgo.New(
        formancesdkgo.WithSecurity(shared.Security{
            Authorization: "Bearer <YOUR_ACCESS_TOKEN_HERE>",
        }),
    )

    ctx := context.Background()
    res, err := s.Orchestration.ListInstances(ctx, operations.ListInstancesRequest{})
    if err != nil {
        log.Fatal(err)
    }

    if res.ListRunsResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                              | Type                                                                                   | Required                                                                               | Description                                                                            |
| -------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------- |
| `ctx`                                                                                  | [context.Context](https://pkg.go.dev/context#Context)                                  | :heavy_check_mark:                                                                     | The context to use for the request.                                                    |
| `request`                                                                              | [operations.ListInstancesRequest](../../pkg/models/operations/listinstancesrequest.md) | :heavy_check_mark:                                                                     | The request object to use for the request.                                             |


### Response

**[*operations.ListInstancesResponse](../../pkg/models/operations/listinstancesresponse.md), error**
| Error Object       | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4xx-5xx            | */*                |

## ListTriggers

List triggers

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	formancesdkgo "github.com/formancehq/formance-sdk-go"
	"context"
	"log"
)

func main() {
    s := formancesdkgo.New(
        formancesdkgo.WithSecurity(shared.Security{
            Authorization: "Bearer <YOUR_ACCESS_TOKEN_HERE>",
        }),
    )

    ctx := context.Background()
    res, err := s.Orchestration.ListTriggers(ctx)
    if err != nil {
        log.Fatal(err)
    }

    if res.ListTriggersResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                             | Type                                                  | Required                                              | Description                                           |
| ----------------------------------------------------- | ----------------------------------------------------- | ----------------------------------------------------- | ----------------------------------------------------- |
| `ctx`                                                 | [context.Context](https://pkg.go.dev/context#Context) | :heavy_check_mark:                                    | The context to use for the request.                   |


### Response

**[*operations.ListTriggersResponse](../../pkg/models/operations/listtriggersresponse.md), error**
| Error Object       | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4xx-5xx            | */*                |

## ListTriggersOccurrences

List triggers occurrences

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	formancesdkgo "github.com/formancehq/formance-sdk-go"
	"context"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"log"
)

func main() {
    s := formancesdkgo.New(
        formancesdkgo.WithSecurity(shared.Security{
            Authorization: "Bearer <YOUR_ACCESS_TOKEN_HERE>",
        }),
    )

    ctx := context.Background()
    res, err := s.Orchestration.ListTriggersOccurrences(ctx, operations.ListTriggersOccurrencesRequest{
        TriggerID: "string",
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.ListTriggersOccurrencesResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                                                  | Type                                                                                                       | Required                                                                                                   | Description                                                                                                |
| ---------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------- |
| `ctx`                                                                                                      | [context.Context](https://pkg.go.dev/context#Context)                                                      | :heavy_check_mark:                                                                                         | The context to use for the request.                                                                        |
| `request`                                                                                                  | [operations.ListTriggersOccurrencesRequest](../../pkg/models/operations/listtriggersoccurrencesrequest.md) | :heavy_check_mark:                                                                                         | The request object to use for the request.                                                                 |


### Response

**[*operations.ListTriggersOccurrencesResponse](../../pkg/models/operations/listtriggersoccurrencesresponse.md), error**
| Error Object       | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4xx-5xx            | */*                |

## ListWorkflows

List registered workflows

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	formancesdkgo "github.com/formancehq/formance-sdk-go"
	"context"
	"log"
)

func main() {
    s := formancesdkgo.New(
        formancesdkgo.WithSecurity(shared.Security{
            Authorization: "Bearer <YOUR_ACCESS_TOKEN_HERE>",
        }),
    )

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

**[*operations.ListWorkflowsResponse](../../pkg/models/operations/listworkflowsresponse.md), error**
| Error Object       | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4xx-5xx            | */*                |

## OrchestrationgetServerInfo

Get server info

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	formancesdkgo "github.com/formancehq/formance-sdk-go"
	"context"
	"log"
)

func main() {
    s := formancesdkgo.New(
        formancesdkgo.WithSecurity(shared.Security{
            Authorization: "Bearer <YOUR_ACCESS_TOKEN_HERE>",
        }),
    )

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

**[*operations.OrchestrationgetServerInfoResponse](../../pkg/models/operations/orchestrationgetserverinforesponse.md), error**
| Error Object       | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4xx-5xx            | */*                |

## ReadTrigger

Read trigger

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	formancesdkgo "github.com/formancehq/formance-sdk-go"
	"context"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"log"
)

func main() {
    s := formancesdkgo.New(
        formancesdkgo.WithSecurity(shared.Security{
            Authorization: "Bearer <YOUR_ACCESS_TOKEN_HERE>",
        }),
    )

    ctx := context.Background()
    res, err := s.Orchestration.ReadTrigger(ctx, operations.ReadTriggerRequest{
        TriggerID: "string",
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.ReadTriggerResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                          | Type                                                                               | Required                                                                           | Description                                                                        |
| ---------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------- |
| `ctx`                                                                              | [context.Context](https://pkg.go.dev/context#Context)                              | :heavy_check_mark:                                                                 | The context to use for the request.                                                |
| `request`                                                                          | [operations.ReadTriggerRequest](../../pkg/models/operations/readtriggerrequest.md) | :heavy_check_mark:                                                                 | The request object to use for the request.                                         |


### Response

**[*operations.ReadTriggerResponse](../../pkg/models/operations/readtriggerresponse.md), error**
| Error Object       | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4xx-5xx            | */*                |

## RunWorkflow

Run workflow

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	formancesdkgo "github.com/formancehq/formance-sdk-go"
	"context"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"log"
)

func main() {
    s := formancesdkgo.New(
        formancesdkgo.WithSecurity(shared.Security{
            Authorization: "Bearer <YOUR_ACCESS_TOKEN_HERE>",
        }),
    )

    ctx := context.Background()
    res, err := s.Orchestration.RunWorkflow(ctx, operations.RunWorkflowRequest{
        RequestBody: map[string]string{
            "key": "string",
        },
        WorkflowID: "string",
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

| Parameter                                                                          | Type                                                                               | Required                                                                           | Description                                                                        |
| ---------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------- |
| `ctx`                                                                              | [context.Context](https://pkg.go.dev/context#Context)                              | :heavy_check_mark:                                                                 | The context to use for the request.                                                |
| `request`                                                                          | [operations.RunWorkflowRequest](../../pkg/models/operations/runworkflowrequest.md) | :heavy_check_mark:                                                                 | The request object to use for the request.                                         |


### Response

**[*operations.RunWorkflowResponse](../../pkg/models/operations/runworkflowresponse.md), error**
| Error Object       | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4xx-5xx            | */*                |

## SendEvent

Send an event to a running workflow

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	formancesdkgo "github.com/formancehq/formance-sdk-go"
	"context"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"log"
	"net/http"
)

func main() {
    s := formancesdkgo.New(
        formancesdkgo.WithSecurity(shared.Security{
            Authorization: "Bearer <YOUR_ACCESS_TOKEN_HERE>",
        }),
    )

    ctx := context.Background()
    res, err := s.Orchestration.SendEvent(ctx, operations.SendEventRequest{
        RequestBody: &operations.SendEventRequestBody{
            Name: "string",
        },
        InstanceID: "string",
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
| `request`                                                                      | [operations.SendEventRequest](../../pkg/models/operations/sendeventrequest.md) | :heavy_check_mark:                                                             | The request object to use for the request.                                     |


### Response

**[*operations.SendEventResponse](../../pkg/models/operations/sendeventresponse.md), error**
| Error Object       | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4xx-5xx            | */*                |

## TestTrigger

Test trigger

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	formancesdkgo "github.com/formancehq/formance-sdk-go"
	"context"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"log"
)

func main() {
    s := formancesdkgo.New(
        formancesdkgo.WithSecurity(shared.Security{
            Authorization: "Bearer <YOUR_ACCESS_TOKEN_HERE>",
        }),
    )

    ctx := context.Background()
    res, err := s.Orchestration.TestTrigger(ctx, operations.TestTriggerRequest{
        RequestBody: map[string]interface{}{
            "key": "string",
        },
        TriggerID: "string",
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.V2TestTriggerResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                          | Type                                                                               | Required                                                                           | Description                                                                        |
| ---------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------- |
| `ctx`                                                                              | [context.Context](https://pkg.go.dev/context#Context)                              | :heavy_check_mark:                                                                 | The context to use for the request.                                                |
| `request`                                                                          | [operations.TestTriggerRequest](../../pkg/models/operations/testtriggerrequest.md) | :heavy_check_mark:                                                                 | The request object to use for the request.                                         |


### Response

**[*operations.TestTriggerResponse](../../pkg/models/operations/testtriggerresponse.md), error**
| Error Object       | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4xx-5xx            | */*                |

## V2CancelEvent

Cancel a running workflow

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	formancesdkgo "github.com/formancehq/formance-sdk-go"
	"context"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"log"
	"net/http"
)

func main() {
    s := formancesdkgo.New(
        formancesdkgo.WithSecurity(shared.Security{
            Authorization: "Bearer <YOUR_ACCESS_TOKEN_HERE>",
        }),
    )

    ctx := context.Background()
    res, err := s.Orchestration.V2CancelEvent(ctx, operations.V2CancelEventRequest{
        InstanceID: "string",
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

| Parameter                                                                              | Type                                                                                   | Required                                                                               | Description                                                                            |
| -------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------- |
| `ctx`                                                                                  | [context.Context](https://pkg.go.dev/context#Context)                                  | :heavy_check_mark:                                                                     | The context to use for the request.                                                    |
| `request`                                                                              | [operations.V2CancelEventRequest](../../pkg/models/operations/v2canceleventrequest.md) | :heavy_check_mark:                                                                     | The request object to use for the request.                                             |


### Response

**[*operations.V2CancelEventResponse](../../pkg/models/operations/v2canceleventresponse.md), error**
| Error Object       | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4xx-5xx            | */*                |

## V2CreateTrigger

Create trigger

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	formancesdkgo "github.com/formancehq/formance-sdk-go"
	"context"
	"log"
)

func main() {
    s := formancesdkgo.New(
        formancesdkgo.WithSecurity(shared.Security{
            Authorization: "Bearer <YOUR_ACCESS_TOKEN_HERE>",
        }),
    )

    ctx := context.Background()
    res, err := s.Orchestration.V2CreateTrigger(ctx, &shared.V2TriggerData{
        Event: "string",
        Vars: map[string]interface{}{
            "key": "string",
        },
        WorkflowID: "string",
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.V2CreateTriggerResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                        | Type                                                             | Required                                                         | Description                                                      |
| ---------------------------------------------------------------- | ---------------------------------------------------------------- | ---------------------------------------------------------------- | ---------------------------------------------------------------- |
| `ctx`                                                            | [context.Context](https://pkg.go.dev/context#Context)            | :heavy_check_mark:                                               | The context to use for the request.                              |
| `request`                                                        | [shared.V2TriggerData](../../pkg/models/shared/v2triggerdata.md) | :heavy_check_mark:                                               | The request object to use for the request.                       |


### Response

**[*operations.V2CreateTriggerResponse](../../pkg/models/operations/v2createtriggerresponse.md), error**
| Error Object       | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4xx-5xx            | */*                |

## V2CreateWorkflow

Create a workflow

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	formancesdkgo "github.com/formancehq/formance-sdk-go"
	"context"
	"log"
)

func main() {
    s := formancesdkgo.New(
        formancesdkgo.WithSecurity(shared.Security{
            Authorization: "Bearer <YOUR_ACCESS_TOKEN_HERE>",
        }),
    )

    ctx := context.Background()
    res, err := s.Orchestration.V2CreateWorkflow(ctx, &shared.V2CreateWorkflowRequest{
        Stages: []map[string]interface{}{
            map[string]interface{}{
                "key": "string",
            },
        },
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.V2CreateWorkflowResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                            | Type                                                                                 | Required                                                                             | Description                                                                          |
| ------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------ |
| `ctx`                                                                                | [context.Context](https://pkg.go.dev/context#Context)                                | :heavy_check_mark:                                                                   | The context to use for the request.                                                  |
| `request`                                                                            | [shared.V2CreateWorkflowRequest](../../pkg/models/shared/v2createworkflowrequest.md) | :heavy_check_mark:                                                                   | The request object to use for the request.                                           |


### Response

**[*operations.V2CreateWorkflowResponse](../../pkg/models/operations/v2createworkflowresponse.md), error**
| Error Object       | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4xx-5xx            | */*                |

## V2DeleteTrigger

Read trigger

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	formancesdkgo "github.com/formancehq/formance-sdk-go"
	"context"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"log"
	"net/http"
)

func main() {
    s := formancesdkgo.New(
        formancesdkgo.WithSecurity(shared.Security{
            Authorization: "Bearer <YOUR_ACCESS_TOKEN_HERE>",
        }),
    )

    ctx := context.Background()
    res, err := s.Orchestration.V2DeleteTrigger(ctx, operations.V2DeleteTriggerRequest{
        TriggerID: "string",
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

| Parameter                                                                                  | Type                                                                                       | Required                                                                                   | Description                                                                                |
| ------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------ |
| `ctx`                                                                                      | [context.Context](https://pkg.go.dev/context#Context)                                      | :heavy_check_mark:                                                                         | The context to use for the request.                                                        |
| `request`                                                                                  | [operations.V2DeleteTriggerRequest](../../pkg/models/operations/v2deletetriggerrequest.md) | :heavy_check_mark:                                                                         | The request object to use for the request.                                                 |


### Response

**[*operations.V2DeleteTriggerResponse](../../pkg/models/operations/v2deletetriggerresponse.md), error**
| Error Object       | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.V2Error  | 404                | application/json   |
| sdkerrors.SDKError | 4xx-5xx            | */*                |

## V2DeleteWorkflow

Delete a flow by id

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	formancesdkgo "github.com/formancehq/formance-sdk-go"
	"context"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"log"
	"net/http"
)

func main() {
    s := formancesdkgo.New(
        formancesdkgo.WithSecurity(shared.Security{
            Authorization: "Bearer <YOUR_ACCESS_TOKEN_HERE>",
        }),
    )

    ctx := context.Background()
    res, err := s.Orchestration.V2DeleteWorkflow(ctx, operations.V2DeleteWorkflowRequest{
        FlowID: "string",
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

| Parameter                                                                                    | Type                                                                                         | Required                                                                                     | Description                                                                                  |
| -------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------- |
| `ctx`                                                                                        | [context.Context](https://pkg.go.dev/context#Context)                                        | :heavy_check_mark:                                                                           | The context to use for the request.                                                          |
| `request`                                                                                    | [operations.V2DeleteWorkflowRequest](../../pkg/models/operations/v2deleteworkflowrequest.md) | :heavy_check_mark:                                                                           | The request object to use for the request.                                                   |


### Response

**[*operations.V2DeleteWorkflowResponse](../../pkg/models/operations/v2deleteworkflowresponse.md), error**
| Error Object       | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.V2Error  | 400,404            | application/json   |
| sdkerrors.SDKError | 4xx-5xx            | */*                |

## V2GetInstance

Get a workflow instance by id

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	formancesdkgo "github.com/formancehq/formance-sdk-go"
	"context"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"log"
)

func main() {
    s := formancesdkgo.New(
        formancesdkgo.WithSecurity(shared.Security{
            Authorization: "Bearer <YOUR_ACCESS_TOKEN_HERE>",
        }),
    )

    ctx := context.Background()
    res, err := s.Orchestration.V2GetInstance(ctx, operations.V2GetInstanceRequest{
        InstanceID: "string",
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.V2GetWorkflowInstanceResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                              | Type                                                                                   | Required                                                                               | Description                                                                            |
| -------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------- |
| `ctx`                                                                                  | [context.Context](https://pkg.go.dev/context#Context)                                  | :heavy_check_mark:                                                                     | The context to use for the request.                                                    |
| `request`                                                                              | [operations.V2GetInstanceRequest](../../pkg/models/operations/v2getinstancerequest.md) | :heavy_check_mark:                                                                     | The request object to use for the request.                                             |


### Response

**[*operations.V2GetInstanceResponse](../../pkg/models/operations/v2getinstanceresponse.md), error**
| Error Object       | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4xx-5xx            | */*                |

## V2GetInstanceHistory

Get a workflow instance history by id

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	formancesdkgo "github.com/formancehq/formance-sdk-go"
	"context"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"log"
)

func main() {
    s := formancesdkgo.New(
        formancesdkgo.WithSecurity(shared.Security{
            Authorization: "Bearer <YOUR_ACCESS_TOKEN_HERE>",
        }),
    )

    ctx := context.Background()
    res, err := s.Orchestration.V2GetInstanceHistory(ctx, operations.V2GetInstanceHistoryRequest{
        InstanceID: "string",
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.V2GetWorkflowInstanceHistoryResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                                            | Type                                                                                                 | Required                                                                                             | Description                                                                                          |
| ---------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------- |
| `ctx`                                                                                                | [context.Context](https://pkg.go.dev/context#Context)                                                | :heavy_check_mark:                                                                                   | The context to use for the request.                                                                  |
| `request`                                                                                            | [operations.V2GetInstanceHistoryRequest](../../pkg/models/operations/v2getinstancehistoryrequest.md) | :heavy_check_mark:                                                                                   | The request object to use for the request.                                                           |


### Response

**[*operations.V2GetInstanceHistoryResponse](../../pkg/models/operations/v2getinstancehistoryresponse.md), error**
| Error Object       | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4xx-5xx            | */*                |

## V2GetInstanceStageHistory

Get a workflow instance stage history

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	formancesdkgo "github.com/formancehq/formance-sdk-go"
	"context"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"log"
)

func main() {
    s := formancesdkgo.New(
        formancesdkgo.WithSecurity(shared.Security{
            Authorization: "Bearer <YOUR_ACCESS_TOKEN_HERE>",
        }),
    )

    ctx := context.Background()
    res, err := s.Orchestration.V2GetInstanceStageHistory(ctx, operations.V2GetInstanceStageHistoryRequest{
        InstanceID: "string",
        Number: 465454,
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.V2GetWorkflowInstanceHistoryStageResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                                                      | Type                                                                                                           | Required                                                                                                       | Description                                                                                                    |
| -------------------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------- |
| `ctx`                                                                                                          | [context.Context](https://pkg.go.dev/context#Context)                                                          | :heavy_check_mark:                                                                                             | The context to use for the request.                                                                            |
| `request`                                                                                                      | [operations.V2GetInstanceStageHistoryRequest](../../pkg/models/operations/v2getinstancestagehistoryrequest.md) | :heavy_check_mark:                                                                                             | The request object to use for the request.                                                                     |


### Response

**[*operations.V2GetInstanceStageHistoryResponse](../../pkg/models/operations/v2getinstancestagehistoryresponse.md), error**
| Error Object       | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4xx-5xx            | */*                |

## V2GetServerInfo

Get server info

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	formancesdkgo "github.com/formancehq/formance-sdk-go"
	"context"
	"log"
)

func main() {
    s := formancesdkgo.New(
        formancesdkgo.WithSecurity(shared.Security{
            Authorization: "Bearer <YOUR_ACCESS_TOKEN_HERE>",
        }),
    )

    ctx := context.Background()
    res, err := s.Orchestration.V2GetServerInfo(ctx)
    if err != nil {
        log.Fatal(err)
    }

    if res.V2ServerInfo != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                             | Type                                                  | Required                                              | Description                                           |
| ----------------------------------------------------- | ----------------------------------------------------- | ----------------------------------------------------- | ----------------------------------------------------- |
| `ctx`                                                 | [context.Context](https://pkg.go.dev/context#Context) | :heavy_check_mark:                                    | The context to use for the request.                   |


### Response

**[*operations.V2GetServerInfoResponse](../../pkg/models/operations/v2getserverinforesponse.md), error**
| Error Object       | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4xx-5xx            | */*                |

## V2GetWorkflow

Get a flow by id

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	formancesdkgo "github.com/formancehq/formance-sdk-go"
	"context"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"log"
)

func main() {
    s := formancesdkgo.New(
        formancesdkgo.WithSecurity(shared.Security{
            Authorization: "Bearer <YOUR_ACCESS_TOKEN_HERE>",
        }),
    )

    ctx := context.Background()
    res, err := s.Orchestration.V2GetWorkflow(ctx, operations.V2GetWorkflowRequest{
        FlowID: "string",
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.V2GetWorkflowResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                              | Type                                                                                   | Required                                                                               | Description                                                                            |
| -------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------- |
| `ctx`                                                                                  | [context.Context](https://pkg.go.dev/context#Context)                                  | :heavy_check_mark:                                                                     | The context to use for the request.                                                    |
| `request`                                                                              | [operations.V2GetWorkflowRequest](../../pkg/models/operations/v2getworkflowrequest.md) | :heavy_check_mark:                                                                     | The request object to use for the request.                                             |


### Response

**[*operations.V2GetWorkflowResponse](../../pkg/models/operations/v2getworkflowresponse.md), error**
| Error Object       | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4xx-5xx            | */*                |

## V2ListInstances

List instances of a workflow

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	formancesdkgo "github.com/formancehq/formance-sdk-go"
	"context"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"log"
)

func main() {
    s := formancesdkgo.New(
        formancesdkgo.WithSecurity(shared.Security{
            Authorization: "Bearer <YOUR_ACCESS_TOKEN_HERE>",
        }),
    )

    ctx := context.Background()
    res, err := s.Orchestration.V2ListInstances(ctx, operations.V2ListInstancesRequest{})
    if err != nil {
        log.Fatal(err)
    }

    if res.V2ListRunsResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                                  | Type                                                                                       | Required                                                                                   | Description                                                                                |
| ------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------ |
| `ctx`                                                                                      | [context.Context](https://pkg.go.dev/context#Context)                                      | :heavy_check_mark:                                                                         | The context to use for the request.                                                        |
| `request`                                                                                  | [operations.V2ListInstancesRequest](../../pkg/models/operations/v2listinstancesrequest.md) | :heavy_check_mark:                                                                         | The request object to use for the request.                                                 |


### Response

**[*operations.V2ListInstancesResponse](../../pkg/models/operations/v2listinstancesresponse.md), error**
| Error Object       | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4xx-5xx            | */*                |

## V2ListTriggers

List triggers

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	formancesdkgo "github.com/formancehq/formance-sdk-go"
	"context"
	"log"
)

func main() {
    s := formancesdkgo.New(
        formancesdkgo.WithSecurity(shared.Security{
            Authorization: "Bearer <YOUR_ACCESS_TOKEN_HERE>",
        }),
    )

    ctx := context.Background()
    res, err := s.Orchestration.V2ListTriggers(ctx)
    if err != nil {
        log.Fatal(err)
    }

    if res.V2ListTriggersResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                             | Type                                                  | Required                                              | Description                                           |
| ----------------------------------------------------- | ----------------------------------------------------- | ----------------------------------------------------- | ----------------------------------------------------- |
| `ctx`                                                 | [context.Context](https://pkg.go.dev/context#Context) | :heavy_check_mark:                                    | The context to use for the request.                   |


### Response

**[*operations.V2ListTriggersResponse](../../pkg/models/operations/v2listtriggersresponse.md), error**
| Error Object       | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4xx-5xx            | */*                |

## V2ListTriggersOccurrences

List triggers occurrences

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	formancesdkgo "github.com/formancehq/formance-sdk-go"
	"context"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"log"
)

func main() {
    s := formancesdkgo.New(
        formancesdkgo.WithSecurity(shared.Security{
            Authorization: "Bearer <YOUR_ACCESS_TOKEN_HERE>",
        }),
    )

    ctx := context.Background()
    res, err := s.Orchestration.V2ListTriggersOccurrences(ctx, operations.V2ListTriggersOccurrencesRequest{
        TriggerID: "string",
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.V2ListTriggersOccurrencesResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                                                      | Type                                                                                                           | Required                                                                                                       | Description                                                                                                    |
| -------------------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------- |
| `ctx`                                                                                                          | [context.Context](https://pkg.go.dev/context#Context)                                                          | :heavy_check_mark:                                                                                             | The context to use for the request.                                                                            |
| `request`                                                                                                      | [operations.V2ListTriggersOccurrencesRequest](../../pkg/models/operations/v2listtriggersoccurrencesrequest.md) | :heavy_check_mark:                                                                                             | The request object to use for the request.                                                                     |


### Response

**[*operations.V2ListTriggersOccurrencesResponse](../../pkg/models/operations/v2listtriggersoccurrencesresponse.md), error**
| Error Object       | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4xx-5xx            | */*                |

## V2ListWorkflows

List registered workflows

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	formancesdkgo "github.com/formancehq/formance-sdk-go"
	"context"
	"log"
)

func main() {
    s := formancesdkgo.New(
        formancesdkgo.WithSecurity(shared.Security{
            Authorization: "Bearer <YOUR_ACCESS_TOKEN_HERE>",
        }),
    )

    ctx := context.Background()
    res, err := s.Orchestration.V2ListWorkflows(ctx)
    if err != nil {
        log.Fatal(err)
    }

    if res.V2ListWorkflowsResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                             | Type                                                  | Required                                              | Description                                           |
| ----------------------------------------------------- | ----------------------------------------------------- | ----------------------------------------------------- | ----------------------------------------------------- |
| `ctx`                                                 | [context.Context](https://pkg.go.dev/context#Context) | :heavy_check_mark:                                    | The context to use for the request.                   |


### Response

**[*operations.V2ListWorkflowsResponse](../../pkg/models/operations/v2listworkflowsresponse.md), error**
| Error Object       | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4xx-5xx            | */*                |

## V2ReadTrigger

Read trigger

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	formancesdkgo "github.com/formancehq/formance-sdk-go"
	"context"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"log"
)

func main() {
    s := formancesdkgo.New(
        formancesdkgo.WithSecurity(shared.Security{
            Authorization: "Bearer <YOUR_ACCESS_TOKEN_HERE>",
        }),
    )

    ctx := context.Background()
    res, err := s.Orchestration.V2ReadTrigger(ctx, operations.V2ReadTriggerRequest{
        TriggerID: "string",
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.V2ReadTriggerResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                              | Type                                                                                   | Required                                                                               | Description                                                                            |
| -------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------- |
| `ctx`                                                                                  | [context.Context](https://pkg.go.dev/context#Context)                                  | :heavy_check_mark:                                                                     | The context to use for the request.                                                    |
| `request`                                                                              | [operations.V2ReadTriggerRequest](../../pkg/models/operations/v2readtriggerrequest.md) | :heavy_check_mark:                                                                     | The request object to use for the request.                                             |


### Response

**[*operations.V2ReadTriggerResponse](../../pkg/models/operations/v2readtriggerresponse.md), error**
| Error Object       | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4xx-5xx            | */*                |

## V2RunWorkflow

Run workflow

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	formancesdkgo "github.com/formancehq/formance-sdk-go"
	"context"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"log"
)

func main() {
    s := formancesdkgo.New(
        formancesdkgo.WithSecurity(shared.Security{
            Authorization: "Bearer <YOUR_ACCESS_TOKEN_HERE>",
        }),
    )

    ctx := context.Background()
    res, err := s.Orchestration.V2RunWorkflow(ctx, operations.V2RunWorkflowRequest{
        RequestBody: map[string]string{
            "key": "string",
        },
        WorkflowID: "string",
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.V2RunWorkflowResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                              | Type                                                                                   | Required                                                                               | Description                                                                            |
| -------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------- |
| `ctx`                                                                                  | [context.Context](https://pkg.go.dev/context#Context)                                  | :heavy_check_mark:                                                                     | The context to use for the request.                                                    |
| `request`                                                                              | [operations.V2RunWorkflowRequest](../../pkg/models/operations/v2runworkflowrequest.md) | :heavy_check_mark:                                                                     | The request object to use for the request.                                             |


### Response

**[*operations.V2RunWorkflowResponse](../../pkg/models/operations/v2runworkflowresponse.md), error**
| Error Object       | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4xx-5xx            | */*                |

## V2SendEvent

Send an event to a running workflow

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	formancesdkgo "github.com/formancehq/formance-sdk-go"
	"context"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"log"
	"net/http"
)

func main() {
    s := formancesdkgo.New(
        formancesdkgo.WithSecurity(shared.Security{
            Authorization: "Bearer <YOUR_ACCESS_TOKEN_HERE>",
        }),
    )

    ctx := context.Background()
    res, err := s.Orchestration.V2SendEvent(ctx, operations.V2SendEventRequest{
        RequestBody: &operations.V2SendEventRequestBody{
            Name: "string",
        },
        InstanceID: "string",
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

| Parameter                                                                          | Type                                                                               | Required                                                                           | Description                                                                        |
| ---------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------- |
| `ctx`                                                                              | [context.Context](https://pkg.go.dev/context#Context)                              | :heavy_check_mark:                                                                 | The context to use for the request.                                                |
| `request`                                                                          | [operations.V2SendEventRequest](../../pkg/models/operations/v2sendeventrequest.md) | :heavy_check_mark:                                                                 | The request object to use for the request.                                         |


### Response

**[*operations.V2SendEventResponse](../../pkg/models/operations/v2sendeventresponse.md), error**
| Error Object       | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4xx-5xx            | */*                |
