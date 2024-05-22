# Webhooks
(*Webhooks*)

### Available Operations

* [AbortWaitingAttempt](#abortwaitingattempt) - Abort one waiting attempt
* [ActivateConfig](#activateconfig) - Activate one config
* [ActivateHook](#activatehook) - Activate one Hook
* [ChangeConfigSecret](#changeconfigsecret) - Change the signing secret of a config
* [DeactivateConfig](#deactivateconfig) - Deactivate one config
* [DeactivateHook](#deactivatehook) - Deactivate one Hook
* [DeleteConfig](#deleteconfig) - Delete one config
* [DeleteHook](#deletehook) - Delete one Hook
* [GetAbortedAttempts](#getabortedattempts) - Get aborted Attempts
* [GetHook](#gethook) - Get one Hook by its ID
* [GetManyConfigs](#getmanyconfigs) - Get many configs
* [GetManyHooks](#getmanyhooks) - Get Many hooks
* [GetWaitingAttempts](#getwaitingattempts) - Get Waiting Attempts
* [InsertConfig](#insertconfig) - Insert a new config
* [InsertHook](#inserthook) - Insert new Hook
* [RetryWaitingAttempt](#retrywaitingattempt) - Retry one waiting Attempt
* [RetryWaitingAttempts](#retrywaitingattempts) - Retry all the waiting attempts
* [TestConfig](#testconfig) - Test one config
* [TestHook](#testhook) - Test one Hook
* [UpdateEndpointHook](#updateendpointhook) - Change the endpoint of one Hook
* [UpdateRetryHook](#updateretryhook) - Change the retry attribute of one Hook
* [UpdateSecretHook](#updatesecrethook) - Change the secret of one Hook

## AbortWaitingAttempt

Abort one waiting attempt

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/formancehq/formance-sdk-go/v2"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"context"
	"log"
)

func main() {
    s := v2.New(
        v2.WithSecurity(shared.Security{
            Authorization: "<YOUR_AUTHORIZATION_HERE>",
        }),
    )

    request := operations.AbortWaitingAttemptRequest{
        AttemptID: "4997257d-dfb6-445b-929c-cbe2ab182818",
    }
    
    ctx := context.Background()
    res, err := s.Webhooks.AbortWaitingAttempt(ctx, request)
    if err != nil {
        log.Fatal(err)
    }
    if res.V2AttemptResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                                          | Type                                                                                               | Required                                                                                           | Description                                                                                        |
| -------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------- |
| `ctx`                                                                                              | [context.Context](https://pkg.go.dev/context#Context)                                              | :heavy_check_mark:                                                                                 | The context to use for the request.                                                                |
| `request`                                                                                          | [operations.AbortWaitingAttemptRequest](../../pkg/models/operations/abortwaitingattemptrequest.md) | :heavy_check_mark:                                                                                 | The request object to use for the request.                                                         |


### Response

**[*operations.AbortWaitingAttemptResponse](../../pkg/models/operations/abortwaitingattemptresponse.md), error**
| Error Object                    | Status Code                     | Content Type                    |
| ------------------------------- | ------------------------------- | ------------------------------- |
| sdkerrors.WebhooksErrorResponse | default                         | application/json                |
| sdkerrors.SDKError              | 4xx-5xx                         | */*                             |

## ActivateConfig

Activate a webhooks config by ID, to start receiving webhooks to its endpoint.

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/formancehq/formance-sdk-go/v2"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"context"
	"log"
)

func main() {
    s := v2.New(
        v2.WithSecurity(shared.Security{
            Authorization: "<YOUR_AUTHORIZATION_HERE>",
        }),
    )

    request := operations.ActivateConfigRequest{
        ID: "4997257d-dfb6-445b-929c-cbe2ab182818",
    }
    
    ctx := context.Background()
    res, err := s.Webhooks.ActivateConfig(ctx, request)
    if err != nil {
        log.Fatal(err)
    }
    if res.ConfigResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                                | Type                                                                                     | Required                                                                                 | Description                                                                              |
| ---------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------- |
| `ctx`                                                                                    | [context.Context](https://pkg.go.dev/context#Context)                                    | :heavy_check_mark:                                                                       | The context to use for the request.                                                      |
| `request`                                                                                | [operations.ActivateConfigRequest](../../pkg/models/operations/activateconfigrequest.md) | :heavy_check_mark:                                                                       | The request object to use for the request.                                               |


### Response

**[*operations.ActivateConfigResponse](../../pkg/models/operations/activateconfigresponse.md), error**
| Error Object                    | Status Code                     | Content Type                    |
| ------------------------------- | ------------------------------- | ------------------------------- |
| sdkerrors.WebhooksErrorResponse | default                         | application/json                |
| sdkerrors.SDKError              | 4xx-5xx                         | */*                             |

## ActivateHook

Activate one hook

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/formancehq/formance-sdk-go/v2"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"context"
	"log"
)

func main() {
    s := v2.New(
        v2.WithSecurity(shared.Security{
            Authorization: "<YOUR_AUTHORIZATION_HERE>",
        }),
    )

    request := operations.ActivateHookRequest{
        HookID: "4997257d-dfb6-445b-929c-cbe2ab182818",
    }
    
    ctx := context.Background()
    res, err := s.Webhooks.ActivateHook(ctx, request)
    if err != nil {
        log.Fatal(err)
    }
    if res.V2HookResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                            | Type                                                                                 | Required                                                                             | Description                                                                          |
| ------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------ |
| `ctx`                                                                                | [context.Context](https://pkg.go.dev/context#Context)                                | :heavy_check_mark:                                                                   | The context to use for the request.                                                  |
| `request`                                                                            | [operations.ActivateHookRequest](../../pkg/models/operations/activatehookrequest.md) | :heavy_check_mark:                                                                   | The request object to use for the request.                                           |


### Response

**[*operations.ActivateHookResponse](../../pkg/models/operations/activatehookresponse.md), error**
| Error Object                    | Status Code                     | Content Type                    |
| ------------------------------- | ------------------------------- | ------------------------------- |
| sdkerrors.WebhooksErrorResponse | default                         | application/json                |
| sdkerrors.SDKError              | 4xx-5xx                         | */*                             |

## ChangeConfigSecret

Change the signing secret of the endpoint of a webhooks config.

If not passed or empty, a secret is automatically generated.
The format is a random string of bytes of size 24, base64 encoded. (larger size after encoding)


### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/formancehq/formance-sdk-go/v2"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"context"
	"log"
)

func main() {
    s := v2.New(
        v2.WithSecurity(shared.Security{
            Authorization: "<YOUR_AUTHORIZATION_HERE>",
        }),
    )

    request := operations.ChangeConfigSecretRequest{
        ConfigChangeSecret: &shared.ConfigChangeSecret{
            Secret: "V0bivxRWveaoz08afqjU6Ko/jwO0Cb+3",
        },
        ID: "4997257d-dfb6-445b-929c-cbe2ab182818",
    }
    
    ctx := context.Background()
    res, err := s.Webhooks.ChangeConfigSecret(ctx, request)
    if err != nil {
        log.Fatal(err)
    }
    if res.ConfigResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                                        | Type                                                                                             | Required                                                                                         | Description                                                                                      |
| ------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------ |
| `ctx`                                                                                            | [context.Context](https://pkg.go.dev/context#Context)                                            | :heavy_check_mark:                                                                               | The context to use for the request.                                                              |
| `request`                                                                                        | [operations.ChangeConfigSecretRequest](../../pkg/models/operations/changeconfigsecretrequest.md) | :heavy_check_mark:                                                                               | The request object to use for the request.                                                       |


### Response

**[*operations.ChangeConfigSecretResponse](../../pkg/models/operations/changeconfigsecretresponse.md), error**
| Error Object                    | Status Code                     | Content Type                    |
| ------------------------------- | ------------------------------- | ------------------------------- |
| sdkerrors.WebhooksErrorResponse | default                         | application/json                |
| sdkerrors.SDKError              | 4xx-5xx                         | */*                             |

## DeactivateConfig

Deactivate a webhooks config by ID, to stop receiving webhooks to its endpoint.

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/formancehq/formance-sdk-go/v2"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"context"
	"log"
)

func main() {
    s := v2.New(
        v2.WithSecurity(shared.Security{
            Authorization: "<YOUR_AUTHORIZATION_HERE>",
        }),
    )

    request := operations.DeactivateConfigRequest{
        ID: "4997257d-dfb6-445b-929c-cbe2ab182818",
    }
    
    ctx := context.Background()
    res, err := s.Webhooks.DeactivateConfig(ctx, request)
    if err != nil {
        log.Fatal(err)
    }
    if res.ConfigResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                                    | Type                                                                                         | Required                                                                                     | Description                                                                                  |
| -------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------- |
| `ctx`                                                                                        | [context.Context](https://pkg.go.dev/context#Context)                                        | :heavy_check_mark:                                                                           | The context to use for the request.                                                          |
| `request`                                                                                    | [operations.DeactivateConfigRequest](../../pkg/models/operations/deactivateconfigrequest.md) | :heavy_check_mark:                                                                           | The request object to use for the request.                                                   |


### Response

**[*operations.DeactivateConfigResponse](../../pkg/models/operations/deactivateconfigresponse.md), error**
| Error Object                    | Status Code                     | Content Type                    |
| ------------------------------- | ------------------------------- | ------------------------------- |
| sdkerrors.WebhooksErrorResponse | default                         | application/json                |
| sdkerrors.SDKError              | 4xx-5xx                         | */*                             |

## DeactivateHook

Deactivate one hook

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/formancehq/formance-sdk-go/v2"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"context"
	"log"
)

func main() {
    s := v2.New(
        v2.WithSecurity(shared.Security{
            Authorization: "<YOUR_AUTHORIZATION_HERE>",
        }),
    )

    request := operations.DeactivateHookRequest{
        HookID: "4997257d-dfb6-445b-929c-cbe2ab182818",
    }
    
    ctx := context.Background()
    res, err := s.Webhooks.DeactivateHook(ctx, request)
    if err != nil {
        log.Fatal(err)
    }
    if res.V2HookResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                                | Type                                                                                     | Required                                                                                 | Description                                                                              |
| ---------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------- |
| `ctx`                                                                                    | [context.Context](https://pkg.go.dev/context#Context)                                    | :heavy_check_mark:                                                                       | The context to use for the request.                                                      |
| `request`                                                                                | [operations.DeactivateHookRequest](../../pkg/models/operations/deactivatehookrequest.md) | :heavy_check_mark:                                                                       | The request object to use for the request.                                               |


### Response

**[*operations.DeactivateHookResponse](../../pkg/models/operations/deactivatehookresponse.md), error**
| Error Object                    | Status Code                     | Content Type                    |
| ------------------------------- | ------------------------------- | ------------------------------- |
| sdkerrors.WebhooksErrorResponse | default                         | application/json                |
| sdkerrors.SDKError              | 4xx-5xx                         | */*                             |

## DeleteConfig

Delete a webhooks config by ID.

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/formancehq/formance-sdk-go/v2"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"context"
	"log"
)

func main() {
    s := v2.New(
        v2.WithSecurity(shared.Security{
            Authorization: "<YOUR_AUTHORIZATION_HERE>",
        }),
    )

    request := operations.DeleteConfigRequest{
        ID: "4997257d-dfb6-445b-929c-cbe2ab182818",
    }
    
    ctx := context.Background()
    res, err := s.Webhooks.DeleteConfig(ctx, request)
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
| `request`                                                                            | [operations.DeleteConfigRequest](../../pkg/models/operations/deleteconfigrequest.md) | :heavy_check_mark:                                                                   | The request object to use for the request.                                           |


### Response

**[*operations.DeleteConfigResponse](../../pkg/models/operations/deleteconfigresponse.md), error**
| Error Object                    | Status Code                     | Content Type                    |
| ------------------------------- | ------------------------------- | ------------------------------- |
| sdkerrors.WebhooksErrorResponse | default                         | application/json                |
| sdkerrors.SDKError              | 4xx-5xx                         | */*                             |

## DeleteHook

Set the status of one Hook as "DELETED"

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/formancehq/formance-sdk-go/v2"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"context"
	"log"
)

func main() {
    s := v2.New(
        v2.WithSecurity(shared.Security{
            Authorization: "<YOUR_AUTHORIZATION_HERE>",
        }),
    )

    request := operations.DeleteHookRequest{
        HookID: "4997257d-dfb6-445b-929c-cbe2ab182818",
    }
    
    ctx := context.Background()
    res, err := s.Webhooks.DeleteHook(ctx, request)
    if err != nil {
        log.Fatal(err)
    }
    if res.V2HookResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                        | Type                                                                             | Required                                                                         | Description                                                                      |
| -------------------------------------------------------------------------------- | -------------------------------------------------------------------------------- | -------------------------------------------------------------------------------- | -------------------------------------------------------------------------------- |
| `ctx`                                                                            | [context.Context](https://pkg.go.dev/context#Context)                            | :heavy_check_mark:                                                               | The context to use for the request.                                              |
| `request`                                                                        | [operations.DeleteHookRequest](../../pkg/models/operations/deletehookrequest.md) | :heavy_check_mark:                                                               | The request object to use for the request.                                       |


### Response

**[*operations.DeleteHookResponse](../../pkg/models/operations/deletehookresponse.md), error**
| Error Object                    | Status Code                     | Content Type                    |
| ------------------------------- | ------------------------------- | ------------------------------- |
| sdkerrors.WebhooksErrorResponse | default                         | application/json                |
| sdkerrors.SDKError              | 4xx-5xx                         | */*                             |

## GetAbortedAttempts

Get Aborted Attempts

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/formancehq/formance-sdk-go/v2"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"context"
	"log"
)

func main() {
    s := v2.New(
        v2.WithSecurity(shared.Security{
            Authorization: "<YOUR_AUTHORIZATION_HERE>",
        }),
    )

    request := operations.GetAbortedAttemptsRequest{}
    
    ctx := context.Background()
    res, err := s.Webhooks.GetAbortedAttempts(ctx, request)
    if err != nil {
        log.Fatal(err)
    }
    if res.V2AttemptCursorResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                                        | Type                                                                                             | Required                                                                                         | Description                                                                                      |
| ------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------ |
| `ctx`                                                                                            | [context.Context](https://pkg.go.dev/context#Context)                                            | :heavy_check_mark:                                                                               | The context to use for the request.                                                              |
| `request`                                                                                        | [operations.GetAbortedAttemptsRequest](../../pkg/models/operations/getabortedattemptsrequest.md) | :heavy_check_mark:                                                                               | The request object to use for the request.                                                       |


### Response

**[*operations.GetAbortedAttemptsResponse](../../pkg/models/operations/getabortedattemptsresponse.md), error**
| Error Object                    | Status Code                     | Content Type                    |
| ------------------------------- | ------------------------------- | ------------------------------- |
| sdkerrors.WebhooksErrorResponse | default                         | application/json                |
| sdkerrors.SDKError              | 4xx-5xx                         | */*                             |

## GetHook

Get one Hook by its ID

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/formancehq/formance-sdk-go/v2"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"context"
	"log"
)

func main() {
    s := v2.New(
        v2.WithSecurity(shared.Security{
            Authorization: "<YOUR_AUTHORIZATION_HERE>",
        }),
    )

    request := operations.GetHookRequest{
        HookID: "4997257d-dfb6-445b-929c-cbe2ab182818",
    }
    
    ctx := context.Background()
    res, err := s.Webhooks.GetHook(ctx, request)
    if err != nil {
        log.Fatal(err)
    }
    if res.V2HookResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                  | Type                                                                       | Required                                                                   | Description                                                                |
| -------------------------------------------------------------------------- | -------------------------------------------------------------------------- | -------------------------------------------------------------------------- | -------------------------------------------------------------------------- |
| `ctx`                                                                      | [context.Context](https://pkg.go.dev/context#Context)                      | :heavy_check_mark:                                                         | The context to use for the request.                                        |
| `request`                                                                  | [operations.GetHookRequest](../../pkg/models/operations/gethookrequest.md) | :heavy_check_mark:                                                         | The request object to use for the request.                                 |


### Response

**[*operations.GetHookResponse](../../pkg/models/operations/gethookresponse.md), error**
| Error Object                    | Status Code                     | Content Type                    |
| ------------------------------- | ------------------------------- | ------------------------------- |
| sdkerrors.WebhooksErrorResponse | default                         | application/json                |
| sdkerrors.SDKError              | 4xx-5xx                         | */*                             |

## GetManyConfigs

Sorted by updated date descending

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/formancehq/formance-sdk-go/v2"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"context"
	"log"
)

func main() {
    s := v2.New(
        v2.WithSecurity(shared.Security{
            Authorization: "<YOUR_AUTHORIZATION_HERE>",
        }),
    )

    request := operations.GetManyConfigsRequest{
        Endpoint: v2.String("https://example.com"),
        ID: v2.String("4997257d-dfb6-445b-929c-cbe2ab182818"),
    }
    
    ctx := context.Background()
    res, err := s.Webhooks.GetManyConfigs(ctx, request)
    if err != nil {
        log.Fatal(err)
    }
    if res.ConfigsResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                                | Type                                                                                     | Required                                                                                 | Description                                                                              |
| ---------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------- |
| `ctx`                                                                                    | [context.Context](https://pkg.go.dev/context#Context)                                    | :heavy_check_mark:                                                                       | The context to use for the request.                                                      |
| `request`                                                                                | [operations.GetManyConfigsRequest](../../pkg/models/operations/getmanyconfigsrequest.md) | :heavy_check_mark:                                                                       | The request object to use for the request.                                               |


### Response

**[*operations.GetManyConfigsResponse](../../pkg/models/operations/getmanyconfigsresponse.md), error**
| Error Object                    | Status Code                     | Content Type                    |
| ------------------------------- | ------------------------------- | ------------------------------- |
| sdkerrors.WebhooksErrorResponse | default                         | application/json                |
| sdkerrors.SDKError              | 4xx-5xx                         | */*                             |

## GetManyHooks

List of Available hooks

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/formancehq/formance-sdk-go/v2"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"context"
	"log"
)

func main() {
    s := v2.New(
        v2.WithSecurity(shared.Security{
            Authorization: "<YOUR_AUTHORIZATION_HERE>",
        }),
    )

    request := operations.GetManyHooksRequest{
        Endpoint: v2.String("https://example.com"),
    }
    
    ctx := context.Background()
    res, err := s.Webhooks.GetManyHooks(ctx, request)
    if err != nil {
        log.Fatal(err)
    }
    if res.V2HookCursorResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                            | Type                                                                                 | Required                                                                             | Description                                                                          |
| ------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------ |
| `ctx`                                                                                | [context.Context](https://pkg.go.dev/context#Context)                                | :heavy_check_mark:                                                                   | The context to use for the request.                                                  |
| `request`                                                                            | [operations.GetManyHooksRequest](../../pkg/models/operations/getmanyhooksrequest.md) | :heavy_check_mark:                                                                   | The request object to use for the request.                                           |


### Response

**[*operations.GetManyHooksResponse](../../pkg/models/operations/getmanyhooksresponse.md), error**
| Error Object                    | Status Code                     | Content Type                    |
| ------------------------------- | ------------------------------- | ------------------------------- |
| sdkerrors.WebhooksErrorResponse | default                         | application/json                |
| sdkerrors.SDKError              | 4xx-5xx                         | */*                             |

## GetWaitingAttempts

Get waiting attempts

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/formancehq/formance-sdk-go/v2"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"context"
	"log"
)

func main() {
    s := v2.New(
        v2.WithSecurity(shared.Security{
            Authorization: "<YOUR_AUTHORIZATION_HERE>",
        }),
    )

    request := operations.GetWaitingAttemptsRequest{}
    
    ctx := context.Background()
    res, err := s.Webhooks.GetWaitingAttempts(ctx, request)
    if err != nil {
        log.Fatal(err)
    }
    if res.V2AttemptCursorResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                                        | Type                                                                                             | Required                                                                                         | Description                                                                                      |
| ------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------ |
| `ctx`                                                                                            | [context.Context](https://pkg.go.dev/context#Context)                                            | :heavy_check_mark:                                                                               | The context to use for the request.                                                              |
| `request`                                                                                        | [operations.GetWaitingAttemptsRequest](../../pkg/models/operations/getwaitingattemptsrequest.md) | :heavy_check_mark:                                                                               | The request object to use for the request.                                                       |


### Response

**[*operations.GetWaitingAttemptsResponse](../../pkg/models/operations/getwaitingattemptsresponse.md), error**
| Error Object                    | Status Code                     | Content Type                    |
| ------------------------------- | ------------------------------- | ------------------------------- |
| sdkerrors.WebhooksErrorResponse | default                         | application/json                |
| sdkerrors.SDKError              | 4xx-5xx                         | */*                             |

## InsertConfig

Insert a new webhooks config.

The endpoint should be a valid https URL and be unique.

The secret is the endpoint's verification secret.
If not passed or empty, a secret is automatically generated.
The format is a random string of bytes of size 24, base64 encoded. (larger size after encoding)

All eventTypes are converted to lower-case when inserted.


### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/formancehq/formance-sdk-go/v2"
	"context"
	"log"
)

func main() {
    s := v2.New(
        v2.WithSecurity(shared.Security{
            Authorization: "<YOUR_AUTHORIZATION_HERE>",
        }),
    )

    request := shared.ConfigUser{
        Endpoint: "https://example.com",
        EventTypes: []string{
            "TYPE1",
            "TYPE2",
        },
        Name: v2.String("customer_payment"),
        Secret: v2.String("V0bivxRWveaoz08afqjU6Ko/jwO0Cb+3"),
    }
    
    ctx := context.Background()
    res, err := s.Webhooks.InsertConfig(ctx, request)
    if err != nil {
        log.Fatal(err)
    }
    if res.ConfigResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                  | Type                                                       | Required                                                   | Description                                                |
| ---------------------------------------------------------- | ---------------------------------------------------------- | ---------------------------------------------------------- | ---------------------------------------------------------- |
| `ctx`                                                      | [context.Context](https://pkg.go.dev/context#Context)      | :heavy_check_mark:                                         | The context to use for the request.                        |
| `request`                                                  | [shared.ConfigUser](../../pkg/models/shared/configuser.md) | :heavy_check_mark:                                         | The request object to use for the request.                 |


### Response

**[*operations.InsertConfigResponse](../../pkg/models/operations/insertconfigresponse.md), error**
| Error Object                    | Status Code                     | Content Type                    |
| ------------------------------- | ------------------------------- | ------------------------------- |
| sdkerrors.WebhooksErrorResponse | default                         | application/json                |
| sdkerrors.SDKError              | 4xx-5xx                         | */*                             |

## InsertHook

Insert new Hook

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/formancehq/formance-sdk-go/v2"
	"context"
	"log"
)

func main() {
    s := v2.New(
        v2.WithSecurity(shared.Security{
            Authorization: "<YOUR_AUTHORIZATION_HERE>",
        }),
    )

    request := shared.V2HookBodyParams{
        Endpoint: "<value>",
        Events: []string{
            "<value>",
        },
    }
    
    ctx := context.Background()
    res, err := s.Webhooks.InsertHook(ctx, request)
    if err != nil {
        log.Fatal(err)
    }
    if res.V2HookResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                              | Type                                                                   | Required                                                               | Description                                                            |
| ---------------------------------------------------------------------- | ---------------------------------------------------------------------- | ---------------------------------------------------------------------- | ---------------------------------------------------------------------- |
| `ctx`                                                                  | [context.Context](https://pkg.go.dev/context#Context)                  | :heavy_check_mark:                                                     | The context to use for the request.                                    |
| `request`                                                              | [shared.V2HookBodyParams](../../pkg/models/shared/v2hookbodyparams.md) | :heavy_check_mark:                                                     | The request object to use for the request.                             |


### Response

**[*operations.InsertHookResponse](../../pkg/models/operations/inserthookresponse.md), error**
| Error Object                    | Status Code                     | Content Type                    |
| ------------------------------- | ------------------------------- | ------------------------------- |
| sdkerrors.WebhooksErrorResponse | default                         | application/json                |
| sdkerrors.SDKError              | 4xx-5xx                         | */*                             |

## RetryWaitingAttempt

Flush one waiting attempt

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/formancehq/formance-sdk-go/v2"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"context"
	"log"
)

func main() {
    s := v2.New(
        v2.WithSecurity(shared.Security{
            Authorization: "<YOUR_AUTHORIZATION_HERE>",
        }),
    )

    request := operations.RetryWaitingAttemptRequest{
        AttemptID: "4997257d-dfb6-445b-929c-cbe2ab182818",
    }
    
    ctx := context.Background()
    res, err := s.Webhooks.RetryWaitingAttempt(ctx, request)
    if err != nil {
        log.Fatal(err)
    }
    if res != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                                          | Type                                                                                               | Required                                                                                           | Description                                                                                        |
| -------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------- |
| `ctx`                                                                                              | [context.Context](https://pkg.go.dev/context#Context)                                              | :heavy_check_mark:                                                                                 | The context to use for the request.                                                                |
| `request`                                                                                          | [operations.RetryWaitingAttemptRequest](../../pkg/models/operations/retrywaitingattemptrequest.md) | :heavy_check_mark:                                                                                 | The request object to use for the request.                                                         |


### Response

**[*operations.RetryWaitingAttemptResponse](../../pkg/models/operations/retrywaitingattemptresponse.md), error**
| Error Object                    | Status Code                     | Content Type                    |
| ------------------------------- | ------------------------------- | ------------------------------- |
| sdkerrors.WebhooksErrorResponse | default                         | application/json                |
| sdkerrors.SDKError              | 4xx-5xx                         | */*                             |

## RetryWaitingAttempts

Flush all waiting attempts

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/formancehq/formance-sdk-go/v2"
	"context"
	"log"
)

func main() {
    s := v2.New(
        v2.WithSecurity(shared.Security{
            Authorization: "<YOUR_AUTHORIZATION_HERE>",
        }),
    )


    
    ctx := context.Background()
    res, err := s.Webhooks.RetryWaitingAttempts(ctx)
    if err != nil {
        log.Fatal(err)
    }
    if res != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                             | Type                                                  | Required                                              | Description                                           |
| ----------------------------------------------------- | ----------------------------------------------------- | ----------------------------------------------------- | ----------------------------------------------------- |
| `ctx`                                                 | [context.Context](https://pkg.go.dev/context#Context) | :heavy_check_mark:                                    | The context to use for the request.                   |


### Response

**[*operations.RetryWaitingAttemptsResponse](../../pkg/models/operations/retrywaitingattemptsresponse.md), error**
| Error Object                    | Status Code                     | Content Type                    |
| ------------------------------- | ------------------------------- | ------------------------------- |
| sdkerrors.WebhooksErrorResponse | default                         | application/json                |
| sdkerrors.SDKError              | 4xx-5xx                         | */*                             |

## TestConfig

Test a config by sending a webhook to its endpoint.

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/formancehq/formance-sdk-go/v2"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"context"
	"log"
)

func main() {
    s := v2.New(
        v2.WithSecurity(shared.Security{
            Authorization: "<YOUR_AUTHORIZATION_HERE>",
        }),
    )

    request := operations.TestConfigRequest{
        ID: "4997257d-dfb6-445b-929c-cbe2ab182818",
    }
    
    ctx := context.Background()
    res, err := s.Webhooks.TestConfig(ctx, request)
    if err != nil {
        log.Fatal(err)
    }
    if res.AttemptResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                        | Type                                                                             | Required                                                                         | Description                                                                      |
| -------------------------------------------------------------------------------- | -------------------------------------------------------------------------------- | -------------------------------------------------------------------------------- | -------------------------------------------------------------------------------- |
| `ctx`                                                                            | [context.Context](https://pkg.go.dev/context#Context)                            | :heavy_check_mark:                                                               | The context to use for the request.                                              |
| `request`                                                                        | [operations.TestConfigRequest](../../pkg/models/operations/testconfigrequest.md) | :heavy_check_mark:                                                               | The request object to use for the request.                                       |


### Response

**[*operations.TestConfigResponse](../../pkg/models/operations/testconfigresponse.md), error**
| Error Object                    | Status Code                     | Content Type                    |
| ------------------------------- | ------------------------------- | ------------------------------- |
| sdkerrors.WebhooksErrorResponse | default                         | application/json                |
| sdkerrors.SDKError              | 4xx-5xx                         | */*                             |

## TestHook

Test one hook by its id

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/formancehq/formance-sdk-go/v2"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"context"
	"log"
)

func main() {
    s := v2.New(
        v2.WithSecurity(shared.Security{
            Authorization: "<YOUR_AUTHORIZATION_HERE>",
        }),
    )

    request := operations.TestHookRequest{
        RequestBody: operations.TestHookRequestBody{},
        HookID: "4997257d-dfb6-445b-929c-cbe2ab182818",
    }
    
    ctx := context.Background()
    res, err := s.Webhooks.TestHook(ctx, request)
    if err != nil {
        log.Fatal(err)
    }
    if res.V2AttemptResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                    | Type                                                                         | Required                                                                     | Description                                                                  |
| ---------------------------------------------------------------------------- | ---------------------------------------------------------------------------- | ---------------------------------------------------------------------------- | ---------------------------------------------------------------------------- |
| `ctx`                                                                        | [context.Context](https://pkg.go.dev/context#Context)                        | :heavy_check_mark:                                                           | The context to use for the request.                                          |
| `request`                                                                    | [operations.TestHookRequest](../../pkg/models/operations/testhookrequest.md) | :heavy_check_mark:                                                           | The request object to use for the request.                                   |


### Response

**[*operations.TestHookResponse](../../pkg/models/operations/testhookresponse.md), error**
| Error Object                    | Status Code                     | Content Type                    |
| ------------------------------- | ------------------------------- | ------------------------------- |
| sdkerrors.WebhooksErrorResponse | default                         | application/json                |
| sdkerrors.SDKError              | 4xx-5xx                         | */*                             |

## UpdateEndpointHook

Change the endpoint of one hook

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/formancehq/formance-sdk-go/v2"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"context"
	"log"
)

func main() {
    s := v2.New(
        v2.WithSecurity(shared.Security{
            Authorization: "<YOUR_AUTHORIZATION_HERE>",
        }),
    )

    request := operations.UpdateEndpointHookRequest{
        RequestBody: operations.UpdateEndpointHookRequestBody{},
        HookID: "4997257d-dfb6-445b-929c-cbe2ab182818",
    }
    
    ctx := context.Background()
    res, err := s.Webhooks.UpdateEndpointHook(ctx, request)
    if err != nil {
        log.Fatal(err)
    }
    if res.V2HookResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                                        | Type                                                                                             | Required                                                                                         | Description                                                                                      |
| ------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------ |
| `ctx`                                                                                            | [context.Context](https://pkg.go.dev/context#Context)                                            | :heavy_check_mark:                                                                               | The context to use for the request.                                                              |
| `request`                                                                                        | [operations.UpdateEndpointHookRequest](../../pkg/models/operations/updateendpointhookrequest.md) | :heavy_check_mark:                                                                               | The request object to use for the request.                                                       |


### Response

**[*operations.UpdateEndpointHookResponse](../../pkg/models/operations/updateendpointhookresponse.md), error**
| Error Object                    | Status Code                     | Content Type                    |
| ------------------------------- | ------------------------------- | ------------------------------- |
| sdkerrors.WebhooksErrorResponse | default                         | application/json                |
| sdkerrors.SDKError              | 4xx-5xx                         | */*                             |

## UpdateRetryHook

Change the retry attribute

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/formancehq/formance-sdk-go/v2"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"context"
	"log"
)

func main() {
    s := v2.New(
        v2.WithSecurity(shared.Security{
            Authorization: "<YOUR_AUTHORIZATION_HERE>",
        }),
    )

    request := operations.UpdateRetryHookRequest{
        RequestBody: operations.UpdateRetryHookRequestBody{},
        HookID: "4997257d-dfb6-445b-929c-cbe2ab182818",
    }
    
    ctx := context.Background()
    res, err := s.Webhooks.UpdateRetryHook(ctx, request)
    if err != nil {
        log.Fatal(err)
    }
    if res.V2HookResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                                  | Type                                                                                       | Required                                                                                   | Description                                                                                |
| ------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------ |
| `ctx`                                                                                      | [context.Context](https://pkg.go.dev/context#Context)                                      | :heavy_check_mark:                                                                         | The context to use for the request.                                                        |
| `request`                                                                                  | [operations.UpdateRetryHookRequest](../../pkg/models/operations/updateretryhookrequest.md) | :heavy_check_mark:                                                                         | The request object to use for the request.                                                 |


### Response

**[*operations.UpdateRetryHookResponse](../../pkg/models/operations/updateretryhookresponse.md), error**
| Error Object                    | Status Code                     | Content Type                    |
| ------------------------------- | ------------------------------- | ------------------------------- |
| sdkerrors.WebhooksErrorResponse | default                         | application/json                |
| sdkerrors.SDKError              | 4xx-5xx                         | */*                             |

## UpdateSecretHook

Change the secret of one Hook

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/formancehq/formance-sdk-go/v2"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"context"
	"log"
)

func main() {
    s := v2.New(
        v2.WithSecurity(shared.Security{
            Authorization: "<YOUR_AUTHORIZATION_HERE>",
        }),
    )

    request := operations.UpdateSecretHookRequest{
        RequestBody: operations.UpdateSecretHookRequestBody{},
        HookID: "4997257d-dfb6-445b-929c-cbe2ab182818",
    }
    
    ctx := context.Background()
    res, err := s.Webhooks.UpdateSecretHook(ctx, request)
    if err != nil {
        log.Fatal(err)
    }
    if res.V2HookResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                                    | Type                                                                                         | Required                                                                                     | Description                                                                                  |
| -------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------- |
| `ctx`                                                                                        | [context.Context](https://pkg.go.dev/context#Context)                                        | :heavy_check_mark:                                                                           | The context to use for the request.                                                          |
| `request`                                                                                    | [operations.UpdateSecretHookRequest](../../pkg/models/operations/updatesecrethookrequest.md) | :heavy_check_mark:                                                                           | The request object to use for the request.                                                   |


### Response

**[*operations.UpdateSecretHookResponse](../../pkg/models/operations/updatesecrethookresponse.md), error**
| Error Object                    | Status Code                     | Content Type                    |
| ------------------------------- | ------------------------------- | ------------------------------- |
| sdkerrors.WebhooksErrorResponse | default                         | application/json                |
| sdkerrors.SDKError              | 4xx-5xx                         | */*                             |
