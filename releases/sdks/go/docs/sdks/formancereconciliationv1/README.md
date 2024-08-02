# FormanceReconciliationV1
(*Reconciliation.V1*)

### Available Operations

* [CreatePolicy](#createpolicy) - Create a policy
* [DeletePolicy](#deletepolicy) - Delete a policy
* [GetPolicy](#getpolicy) - Get a policy
* [GetReconciliation](#getreconciliation) - Get a reconciliation
* [ListPolicies](#listpolicies) - List policies
* [ListReconciliations](#listreconciliations) - List reconciliations
* [Reconcile](#reconcile) - Reconcile using a policy
* [ReconciliationgetServerInfo](#reconciliationgetserverinfo) - Get server info

## CreatePolicy

Create a policy

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"os"
	"github.com/formancehq/formance-sdk-go/v2"
	"context"
	"log"
)

func main() {
    s := v2.New(
        v2.WithSecurity(shared.Security{
            Authorization: os.Getenv("AUTHORIZATION"),
        }),
    )
    request := shared.PolicyRequest{
        LedgerName: "default",
        LedgerQuery: map[string]any{
            "key": "<value>",
        },
        Name: "XXX",
        PaymentsPoolID: "XXX",
    }
    ctx := context.Background()
    res, err := s.Reconciliation.V1.CreatePolicy(ctx, request)
    if err != nil {
        log.Fatal(err)
    }
    if res.PolicyResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                        | Type                                                             | Required                                                         | Description                                                      |
| ---------------------------------------------------------------- | ---------------------------------------------------------------- | ---------------------------------------------------------------- | ---------------------------------------------------------------- |
| `ctx`                                                            | [context.Context](https://pkg.go.dev/context#Context)            | :heavy_check_mark:                                               | The context to use for the request.                              |
| `request`                                                        | [shared.PolicyRequest](../../pkg/models/shared/policyrequest.md) | :heavy_check_mark:                                               | The request object to use for the request.                       |
| `opts`                                                           | [][operations.Option](../../pkg/models/operations/option.md)     | :heavy_minus_sign:                                               | The options for this request.                                    |


### Response

**[*operations.CreatePolicyResponse](../../pkg/models/operations/createpolicyresponse.md), error**
| Error Object                          | Status Code                           | Content Type                          |
| ------------------------------------- | ------------------------------------- | ------------------------------------- |
| sdkerrors.ReconciliationErrorResponse | default                               | application/json                      |
| sdkerrors.SDKError                    | 4xx-5xx                               | */*                                   |

## DeletePolicy

Delete a policy by its id.

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"os"
	"github.com/formancehq/formance-sdk-go/v2"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"context"
	"log"
)

func main() {
    s := v2.New(
        v2.WithSecurity(shared.Security{
            Authorization: os.Getenv("AUTHORIZATION"),
        }),
    )
    request := operations.DeletePolicyRequest{
        PolicyID: "XXX",
    }
    ctx := context.Background()
    res, err := s.Reconciliation.V1.DeletePolicy(ctx, request)
    if err != nil {
        log.Fatal(err)
    }
    if res != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                            | Type                                                                                 | Required                                                                             | Description                                                                          |
| ------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------ |
| `ctx`                                                                                | [context.Context](https://pkg.go.dev/context#Context)                                | :heavy_check_mark:                                                                   | The context to use for the request.                                                  |
| `request`                                                                            | [operations.DeletePolicyRequest](../../pkg/models/operations/deletepolicyrequest.md) | :heavy_check_mark:                                                                   | The request object to use for the request.                                           |
| `opts`                                                                               | [][operations.Option](../../pkg/models/operations/option.md)                         | :heavy_minus_sign:                                                                   | The options for this request.                                                        |


### Response

**[*operations.DeletePolicyResponse](../../pkg/models/operations/deletepolicyresponse.md), error**
| Error Object                          | Status Code                           | Content Type                          |
| ------------------------------------- | ------------------------------------- | ------------------------------------- |
| sdkerrors.ReconciliationErrorResponse | default                               | application/json                      |
| sdkerrors.SDKError                    | 4xx-5xx                               | */*                                   |

## GetPolicy

Get a policy

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"os"
	"github.com/formancehq/formance-sdk-go/v2"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"context"
	"log"
)

func main() {
    s := v2.New(
        v2.WithSecurity(shared.Security{
            Authorization: os.Getenv("AUTHORIZATION"),
        }),
    )
    request := operations.GetPolicyRequest{
        PolicyID: "XXX",
    }
    ctx := context.Background()
    res, err := s.Reconciliation.V1.GetPolicy(ctx, request)
    if err != nil {
        log.Fatal(err)
    }
    if res.PolicyResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                      | Type                                                                           | Required                                                                       | Description                                                                    |
| ------------------------------------------------------------------------------ | ------------------------------------------------------------------------------ | ------------------------------------------------------------------------------ | ------------------------------------------------------------------------------ |
| `ctx`                                                                          | [context.Context](https://pkg.go.dev/context#Context)                          | :heavy_check_mark:                                                             | The context to use for the request.                                            |
| `request`                                                                      | [operations.GetPolicyRequest](../../pkg/models/operations/getpolicyrequest.md) | :heavy_check_mark:                                                             | The request object to use for the request.                                     |
| `opts`                                                                         | [][operations.Option](../../pkg/models/operations/option.md)                   | :heavy_minus_sign:                                                             | The options for this request.                                                  |


### Response

**[*operations.GetPolicyResponse](../../pkg/models/operations/getpolicyresponse.md), error**
| Error Object                          | Status Code                           | Content Type                          |
| ------------------------------------- | ------------------------------------- | ------------------------------------- |
| sdkerrors.ReconciliationErrorResponse | default                               | application/json                      |
| sdkerrors.SDKError                    | 4xx-5xx                               | */*                                   |

## GetReconciliation

Get a reconciliation

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"os"
	"github.com/formancehq/formance-sdk-go/v2"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"context"
	"log"
)

func main() {
    s := v2.New(
        v2.WithSecurity(shared.Security{
            Authorization: os.Getenv("AUTHORIZATION"),
        }),
    )
    request := operations.GetReconciliationRequest{
        ReconciliationID: "XXX",
    }
    ctx := context.Background()
    res, err := s.Reconciliation.V1.GetReconciliation(ctx, request)
    if err != nil {
        log.Fatal(err)
    }
    if res.ReconciliationResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                                      | Type                                                                                           | Required                                                                                       | Description                                                                                    |
| ---------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------- |
| `ctx`                                                                                          | [context.Context](https://pkg.go.dev/context#Context)                                          | :heavy_check_mark:                                                                             | The context to use for the request.                                                            |
| `request`                                                                                      | [operations.GetReconciliationRequest](../../pkg/models/operations/getreconciliationrequest.md) | :heavy_check_mark:                                                                             | The request object to use for the request.                                                     |
| `opts`                                                                                         | [][operations.Option](../../pkg/models/operations/option.md)                                   | :heavy_minus_sign:                                                                             | The options for this request.                                                                  |


### Response

**[*operations.GetReconciliationResponse](../../pkg/models/operations/getreconciliationresponse.md), error**
| Error Object                          | Status Code                           | Content Type                          |
| ------------------------------------- | ------------------------------------- | ------------------------------------- |
| sdkerrors.ReconciliationErrorResponse | default                               | application/json                      |
| sdkerrors.SDKError                    | 4xx-5xx                               | */*                                   |

## ListPolicies

List policies

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"os"
	"github.com/formancehq/formance-sdk-go/v2"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"context"
	"log"
)

func main() {
    s := v2.New(
        v2.WithSecurity(shared.Security{
            Authorization: os.Getenv("AUTHORIZATION"),
        }),
    )
    request := operations.ListPoliciesRequest{
        Cursor: v2.String("aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ=="),
        PageSize: v2.Int64(100),
    }
    ctx := context.Background()
    res, err := s.Reconciliation.V1.ListPolicies(ctx, request)
    if err != nil {
        log.Fatal(err)
    }
    if res.PoliciesCursorResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                            | Type                                                                                 | Required                                                                             | Description                                                                          |
| ------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------ |
| `ctx`                                                                                | [context.Context](https://pkg.go.dev/context#Context)                                | :heavy_check_mark:                                                                   | The context to use for the request.                                                  |
| `request`                                                                            | [operations.ListPoliciesRequest](../../pkg/models/operations/listpoliciesrequest.md) | :heavy_check_mark:                                                                   | The request object to use for the request.                                           |
| `opts`                                                                               | [][operations.Option](../../pkg/models/operations/option.md)                         | :heavy_minus_sign:                                                                   | The options for this request.                                                        |


### Response

**[*operations.ListPoliciesResponse](../../pkg/models/operations/listpoliciesresponse.md), error**
| Error Object                          | Status Code                           | Content Type                          |
| ------------------------------------- | ------------------------------------- | ------------------------------------- |
| sdkerrors.ReconciliationErrorResponse | default                               | application/json                      |
| sdkerrors.SDKError                    | 4xx-5xx                               | */*                                   |

## ListReconciliations

List reconciliations

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"os"
	"github.com/formancehq/formance-sdk-go/v2"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"context"
	"log"
)

func main() {
    s := v2.New(
        v2.WithSecurity(shared.Security{
            Authorization: os.Getenv("AUTHORIZATION"),
        }),
    )
    request := operations.ListReconciliationsRequest{
        Cursor: v2.String("aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ=="),
        PageSize: v2.Int64(100),
    }
    ctx := context.Background()
    res, err := s.Reconciliation.V1.ListReconciliations(ctx, request)
    if err != nil {
        log.Fatal(err)
    }
    if res.ReconciliationsCursorResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                                          | Type                                                                                               | Required                                                                                           | Description                                                                                        |
| -------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------- |
| `ctx`                                                                                              | [context.Context](https://pkg.go.dev/context#Context)                                              | :heavy_check_mark:                                                                                 | The context to use for the request.                                                                |
| `request`                                                                                          | [operations.ListReconciliationsRequest](../../pkg/models/operations/listreconciliationsrequest.md) | :heavy_check_mark:                                                                                 | The request object to use for the request.                                                         |
| `opts`                                                                                             | [][operations.Option](../../pkg/models/operations/option.md)                                       | :heavy_minus_sign:                                                                                 | The options for this request.                                                                      |


### Response

**[*operations.ListReconciliationsResponse](../../pkg/models/operations/listreconciliationsresponse.md), error**
| Error Object                          | Status Code                           | Content Type                          |
| ------------------------------------- | ------------------------------------- | ------------------------------------- |
| sdkerrors.ReconciliationErrorResponse | default                               | application/json                      |
| sdkerrors.SDKError                    | 4xx-5xx                               | */*                                   |

## Reconcile

Reconcile using a policy

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"os"
	"github.com/formancehq/formance-sdk-go/v2"
	"github.com/formancehq/formance-sdk-go/v2/pkg/types"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"context"
	"log"
)

func main() {
    s := v2.New(
        v2.WithSecurity(shared.Security{
            Authorization: os.Getenv("AUTHORIZATION"),
        }),
    )
    request := operations.ReconcileRequest{
        ReconciliationRequest: shared.ReconciliationRequest{
            ReconciledAtLedger: types.MustTimeFromString("2021-01-01T00:00:00.000Z"),
            ReconciledAtPayments: types.MustTimeFromString("2021-01-01T00:00:00.000Z"),
        },
        PolicyID: "XXX",
    }
    ctx := context.Background()
    res, err := s.Reconciliation.V1.Reconcile(ctx, request)
    if err != nil {
        log.Fatal(err)
    }
    if res.ReconciliationResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                      | Type                                                                           | Required                                                                       | Description                                                                    |
| ------------------------------------------------------------------------------ | ------------------------------------------------------------------------------ | ------------------------------------------------------------------------------ | ------------------------------------------------------------------------------ |
| `ctx`                                                                          | [context.Context](https://pkg.go.dev/context#Context)                          | :heavy_check_mark:                                                             | The context to use for the request.                                            |
| `request`                                                                      | [operations.ReconcileRequest](../../pkg/models/operations/reconcilerequest.md) | :heavy_check_mark:                                                             | The request object to use for the request.                                     |
| `opts`                                                                         | [][operations.Option](../../pkg/models/operations/option.md)                   | :heavy_minus_sign:                                                             | The options for this request.                                                  |


### Response

**[*operations.ReconcileResponse](../../pkg/models/operations/reconcileresponse.md), error**
| Error Object                          | Status Code                           | Content Type                          |
| ------------------------------------- | ------------------------------------- | ------------------------------------- |
| sdkerrors.ReconciliationErrorResponse | default                               | application/json                      |
| sdkerrors.SDKError                    | 4xx-5xx                               | */*                                   |

## ReconciliationgetServerInfo

Get server info

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"os"
	"github.com/formancehq/formance-sdk-go/v2"
	"context"
	"log"
)

func main() {
    s := v2.New(
        v2.WithSecurity(shared.Security{
            Authorization: os.Getenv("AUTHORIZATION"),
        }),
    )

    ctx := context.Background()
    res, err := s.Reconciliation.V1.ReconciliationgetServerInfo(ctx)
    if err != nil {
        log.Fatal(err)
    }
    if res.ServerInfo != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                    | Type                                                         | Required                                                     | Description                                                  |
| ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ |
| `ctx`                                                        | [context.Context](https://pkg.go.dev/context#Context)        | :heavy_check_mark:                                           | The context to use for the request.                          |
| `opts`                                                       | [][operations.Option](../../pkg/models/operations/option.md) | :heavy_minus_sign:                                           | The options for this request.                                |


### Response

**[*operations.ReconciliationgetServerInfoResponse](../../pkg/models/operations/reconciliationgetserverinforesponse.md), error**
| Error Object                          | Status Code                           | Content Type                          |
| ------------------------------------- | ------------------------------------- | ------------------------------------- |
| sdkerrors.ReconciliationErrorResponse | default                               | application/json                      |
| sdkerrors.SDKError                    | 4xx-5xx                               | */*                                   |
