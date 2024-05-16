# Auth
(*Auth*)

### Available Operations

* [CreateClient](#createclient) - Create client
* [CreateSecret](#createsecret) - Add a secret to a client
* [DeleteClient](#deleteclient) - Delete client
* [DeleteSecret](#deletesecret) - Delete a secret from a client
* [ListClients](#listclients) - List clients
* [ListUsers](#listusers) - List users
* [ReadClient](#readclient) - Read client
* [ReadUser](#readuser) - Read user
* [UpdateClient](#updateclient) - Update client

## CreateClient

Create client

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

    var request *shared.CreateClientRequest = &shared.CreateClientRequest{
        Name: "<value>",
    }
    
    ctx := context.Background()
    res, err := s.Auth.CreateClient(ctx, request)
    if err != nil {
        log.Fatal(err)
    }
    if res.CreateClientResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                    | Type                                                                         | Required                                                                     | Description                                                                  |
| ---------------------------------------------------------------------------- | ---------------------------------------------------------------------------- | ---------------------------------------------------------------------------- | ---------------------------------------------------------------------------- |
| `ctx`                                                                        | [context.Context](https://pkg.go.dev/context#Context)                        | :heavy_check_mark:                                                           | The context to use for the request.                                          |
| `request`                                                                    | [shared.CreateClientRequest](../../pkg/models/shared/createclientrequest.md) | :heavy_check_mark:                                                           | The request object to use for the request.                                   |


### Response

**[*operations.CreateClientResponse](../../pkg/models/operations/createclientresponse.md), error**
| Error Object       | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4xx-5xx            | */*                |

## CreateSecret

Add a secret to a client

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

    request := operations.CreateSecretRequest{
        ClientID: "<value>",
    }
    
    ctx := context.Background()
    res, err := s.Auth.CreateSecret(ctx, request)
    if err != nil {
        log.Fatal(err)
    }
    if res.CreateSecretResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                            | Type                                                                                 | Required                                                                             | Description                                                                          |
| ------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------ |
| `ctx`                                                                                | [context.Context](https://pkg.go.dev/context#Context)                                | :heavy_check_mark:                                                                   | The context to use for the request.                                                  |
| `request`                                                                            | [operations.CreateSecretRequest](../../pkg/models/operations/createsecretrequest.md) | :heavy_check_mark:                                                                   | The request object to use for the request.                                           |


### Response

**[*operations.CreateSecretResponse](../../pkg/models/operations/createsecretresponse.md), error**
| Error Object       | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4xx-5xx            | */*                |

## DeleteClient

Delete client

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

    request := operations.DeleteClientRequest{
        ClientID: "<value>",
    }
    
    ctx := context.Background()
    res, err := s.Auth.DeleteClient(ctx, request)
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
| `request`                                                                            | [operations.DeleteClientRequest](../../pkg/models/operations/deleteclientrequest.md) | :heavy_check_mark:                                                                   | The request object to use for the request.                                           |


### Response

**[*operations.DeleteClientResponse](../../pkg/models/operations/deleteclientresponse.md), error**
| Error Object       | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4xx-5xx            | */*                |

## DeleteSecret

Delete a secret from a client

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

    request := operations.DeleteSecretRequest{
        ClientID: "<value>",
        SecretID: "<value>",
    }
    
    ctx := context.Background()
    res, err := s.Auth.DeleteSecret(ctx, request)
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
| `request`                                                                            | [operations.DeleteSecretRequest](../../pkg/models/operations/deletesecretrequest.md) | :heavy_check_mark:                                                                   | The request object to use for the request.                                           |


### Response

**[*operations.DeleteSecretResponse](../../pkg/models/operations/deletesecretresponse.md), error**
| Error Object       | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4xx-5xx            | */*                |

## ListClients

List clients

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
    res, err := s.Auth.ListClients(ctx)
    if err != nil {
        log.Fatal(err)
    }
    if res.ListClientsResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                             | Type                                                  | Required                                              | Description                                           |
| ----------------------------------------------------- | ----------------------------------------------------- | ----------------------------------------------------- | ----------------------------------------------------- |
| `ctx`                                                 | [context.Context](https://pkg.go.dev/context#Context) | :heavy_check_mark:                                    | The context to use for the request.                   |


### Response

**[*operations.ListClientsResponse](../../pkg/models/operations/listclientsresponse.md), error**
| Error Object       | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4xx-5xx            | */*                |

## ListUsers

List users

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
    res, err := s.Auth.ListUsers(ctx)
    if err != nil {
        log.Fatal(err)
    }
    if res.ListUsersResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                             | Type                                                  | Required                                              | Description                                           |
| ----------------------------------------------------- | ----------------------------------------------------- | ----------------------------------------------------- | ----------------------------------------------------- |
| `ctx`                                                 | [context.Context](https://pkg.go.dev/context#Context) | :heavy_check_mark:                                    | The context to use for the request.                   |


### Response

**[*operations.ListUsersResponse](../../pkg/models/operations/listusersresponse.md), error**
| Error Object       | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4xx-5xx            | */*                |

## ReadClient

Read client

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

    request := operations.ReadClientRequest{
        ClientID: "<value>",
    }
    
    ctx := context.Background()
    res, err := s.Auth.ReadClient(ctx, request)
    if err != nil {
        log.Fatal(err)
    }
    if res.ReadClientResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                        | Type                                                                             | Required                                                                         | Description                                                                      |
| -------------------------------------------------------------------------------- | -------------------------------------------------------------------------------- | -------------------------------------------------------------------------------- | -------------------------------------------------------------------------------- |
| `ctx`                                                                            | [context.Context](https://pkg.go.dev/context#Context)                            | :heavy_check_mark:                                                               | The context to use for the request.                                              |
| `request`                                                                        | [operations.ReadClientRequest](../../pkg/models/operations/readclientrequest.md) | :heavy_check_mark:                                                               | The request object to use for the request.                                       |


### Response

**[*operations.ReadClientResponse](../../pkg/models/operations/readclientresponse.md), error**
| Error Object       | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4xx-5xx            | */*                |

## ReadUser

Read user

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

    request := operations.ReadUserRequest{
        UserID: "<value>",
    }
    
    ctx := context.Background()
    res, err := s.Auth.ReadUser(ctx, request)
    if err != nil {
        log.Fatal(err)
    }
    if res.ReadUserResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                    | Type                                                                         | Required                                                                     | Description                                                                  |
| ---------------------------------------------------------------------------- | ---------------------------------------------------------------------------- | ---------------------------------------------------------------------------- | ---------------------------------------------------------------------------- |
| `ctx`                                                                        | [context.Context](https://pkg.go.dev/context#Context)                        | :heavy_check_mark:                                                           | The context to use for the request.                                          |
| `request`                                                                    | [operations.ReadUserRequest](../../pkg/models/operations/readuserrequest.md) | :heavy_check_mark:                                                           | The request object to use for the request.                                   |


### Response

**[*operations.ReadUserResponse](../../pkg/models/operations/readuserresponse.md), error**
| Error Object       | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4xx-5xx            | */*                |

## UpdateClient

Update client

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

    request := operations.UpdateClientRequest{
        ClientID: "<value>",
    }
    
    ctx := context.Background()
    res, err := s.Auth.UpdateClient(ctx, request)
    if err != nil {
        log.Fatal(err)
    }
    if res.UpdateClientResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                            | Type                                                                                 | Required                                                                             | Description                                                                          |
| ------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------ |
| `ctx`                                                                                | [context.Context](https://pkg.go.dev/context#Context)                                | :heavy_check_mark:                                                                   | The context to use for the request.                                                  |
| `request`                                                                            | [operations.UpdateClientRequest](../../pkg/models/operations/updateclientrequest.md) | :heavy_check_mark:                                                                   | The request object to use for the request.                                           |


### Response

**[*operations.UpdateClientResponse](../../pkg/models/operations/updateclientresponse.md), error**
| Error Object       | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4xx-5xx            | */*                |
