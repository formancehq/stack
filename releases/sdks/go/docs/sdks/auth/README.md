# Auth
(*Auth*)

### Available Operations

* [CreateClient](#createclient) - Create client
* [CreateSecret](#createsecret) - Add a secret to a client
* [DeleteClient](#deleteclient) - Delete client
* [DeleteSecret](#deletesecret) - Delete a secret from a client
* [GetServerInfo](#getserverinfo) - Get server info
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
	"context"
	"log"
	formancesdkgo "github.com/formancehq/formance-sdk-go"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
)

func main() {
    s := formancesdkgo.New()

    ctx := context.Background()
    res, err := s.Auth.CreateClient(ctx, &shared.CreateClientRequest{
        Metadata: map[string]interface{}{
            "key": "string",
        },
        Name: "string",
        PostLogoutRedirectUris: []string{
            "string",
        },
        RedirectUris: []string{
            "string",
        },
        Scopes: []string{
            "string",
        },
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.CreateClientResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                | Type                                                                     | Required                                                                 | Description                                                              |
| ------------------------------------------------------------------------ | ------------------------------------------------------------------------ | ------------------------------------------------------------------------ | ------------------------------------------------------------------------ |
| `ctx`                                                                    | [context.Context](https://pkg.go.dev/context#Context)                    | :heavy_check_mark:                                                       | The context to use for the request.                                      |
| `request`                                                                | [shared.CreateClientRequest](../../models/shared/createclientrequest.md) | :heavy_check_mark:                                                       | The request object to use for the request.                               |


### Response

**[*operations.CreateClientResponse](../../models/operations/createclientresponse.md), error**


## CreateSecret

Add a secret to a client

### Example Usage

```go
package main

import(
	"context"
	"log"
	formancesdkgo "github.com/formancehq/formance-sdk-go"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
)

func main() {
    s := formancesdkgo.New()

    ctx := context.Background()
    res, err := s.Auth.CreateSecret(ctx, operations.CreateSecretRequest{
        CreateSecretRequest: &shared.CreateSecretRequest{
            Metadata: map[string]interface{}{
                "key": "string",
            },
            Name: "string",
        },
        ClientID: "string",
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.CreateSecretResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                        | Type                                                                             | Required                                                                         | Description                                                                      |
| -------------------------------------------------------------------------------- | -------------------------------------------------------------------------------- | -------------------------------------------------------------------------------- | -------------------------------------------------------------------------------- |
| `ctx`                                                                            | [context.Context](https://pkg.go.dev/context#Context)                            | :heavy_check_mark:                                                               | The context to use for the request.                                              |
| `request`                                                                        | [operations.CreateSecretRequest](../../models/operations/createsecretrequest.md) | :heavy_check_mark:                                                               | The request object to use for the request.                                       |


### Response

**[*operations.CreateSecretResponse](../../models/operations/createsecretresponse.md), error**


## DeleteClient

Delete client

### Example Usage

```go
package main

import(
	"context"
	"log"
	formancesdkgo "github.com/formancehq/formance-sdk-go"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
)

func main() {
    s := formancesdkgo.New()

    ctx := context.Background()
    res, err := s.Auth.DeleteClient(ctx, operations.DeleteClientRequest{
        ClientID: "string",
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

| Parameter                                                                        | Type                                                                             | Required                                                                         | Description                                                                      |
| -------------------------------------------------------------------------------- | -------------------------------------------------------------------------------- | -------------------------------------------------------------------------------- | -------------------------------------------------------------------------------- |
| `ctx`                                                                            | [context.Context](https://pkg.go.dev/context#Context)                            | :heavy_check_mark:                                                               | The context to use for the request.                                              |
| `request`                                                                        | [operations.DeleteClientRequest](../../models/operations/deleteclientrequest.md) | :heavy_check_mark:                                                               | The request object to use for the request.                                       |


### Response

**[*operations.DeleteClientResponse](../../models/operations/deleteclientresponse.md), error**


## DeleteSecret

Delete a secret from a client

### Example Usage

```go
package main

import(
	"context"
	"log"
	formancesdkgo "github.com/formancehq/formance-sdk-go"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
)

func main() {
    s := formancesdkgo.New()

    ctx := context.Background()
    res, err := s.Auth.DeleteSecret(ctx, operations.DeleteSecretRequest{
        ClientID: "string",
        SecretID: "string",
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

| Parameter                                                                        | Type                                                                             | Required                                                                         | Description                                                                      |
| -------------------------------------------------------------------------------- | -------------------------------------------------------------------------------- | -------------------------------------------------------------------------------- | -------------------------------------------------------------------------------- |
| `ctx`                                                                            | [context.Context](https://pkg.go.dev/context#Context)                            | :heavy_check_mark:                                                               | The context to use for the request.                                              |
| `request`                                                                        | [operations.DeleteSecretRequest](../../models/operations/deletesecretrequest.md) | :heavy_check_mark:                                                               | The request object to use for the request.                                       |


### Response

**[*operations.DeleteSecretResponse](../../models/operations/deletesecretresponse.md), error**


## GetServerInfo

Get server info

### Example Usage

```go
package main

import(
	"context"
	"log"
	formancesdkgo "github.com/formancehq/formance-sdk-go"
)

func main() {
    s := formancesdkgo.New()

    ctx := context.Background()
    res, err := s.Auth.GetServerInfo(ctx)
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

**[*operations.GetServerInfoResponse](../../models/operations/getserverinforesponse.md), error**


## ListClients

List clients

### Example Usage

```go
package main

import(
	"context"
	"log"
	formancesdkgo "github.com/formancehq/formance-sdk-go"
)

func main() {
    s := formancesdkgo.New()

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

**[*operations.ListClientsResponse](../../models/operations/listclientsresponse.md), error**


## ListUsers

List users

### Example Usage

```go
package main

import(
	"context"
	"log"
	formancesdkgo "github.com/formancehq/formance-sdk-go"
)

func main() {
    s := formancesdkgo.New()

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

**[*operations.ListUsersResponse](../../models/operations/listusersresponse.md), error**


## ReadClient

Read client

### Example Usage

```go
package main

import(
	"context"
	"log"
	formancesdkgo "github.com/formancehq/formance-sdk-go"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
)

func main() {
    s := formancesdkgo.New()

    ctx := context.Background()
    res, err := s.Auth.ReadClient(ctx, operations.ReadClientRequest{
        ClientID: "string",
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.ReadClientResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                    | Type                                                                         | Required                                                                     | Description                                                                  |
| ---------------------------------------------------------------------------- | ---------------------------------------------------------------------------- | ---------------------------------------------------------------------------- | ---------------------------------------------------------------------------- |
| `ctx`                                                                        | [context.Context](https://pkg.go.dev/context#Context)                        | :heavy_check_mark:                                                           | The context to use for the request.                                          |
| `request`                                                                    | [operations.ReadClientRequest](../../models/operations/readclientrequest.md) | :heavy_check_mark:                                                           | The request object to use for the request.                                   |


### Response

**[*operations.ReadClientResponse](../../models/operations/readclientresponse.md), error**


## ReadUser

Read user

### Example Usage

```go
package main

import(
	"context"
	"log"
	formancesdkgo "github.com/formancehq/formance-sdk-go"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
)

func main() {
    s := formancesdkgo.New()

    ctx := context.Background()
    res, err := s.Auth.ReadUser(ctx, operations.ReadUserRequest{
        UserID: "string",
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.ReadUserResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                | Type                                                                     | Required                                                                 | Description                                                              |
| ------------------------------------------------------------------------ | ------------------------------------------------------------------------ | ------------------------------------------------------------------------ | ------------------------------------------------------------------------ |
| `ctx`                                                                    | [context.Context](https://pkg.go.dev/context#Context)                    | :heavy_check_mark:                                                       | The context to use for the request.                                      |
| `request`                                                                | [operations.ReadUserRequest](../../models/operations/readuserrequest.md) | :heavy_check_mark:                                                       | The request object to use for the request.                               |


### Response

**[*operations.ReadUserResponse](../../models/operations/readuserresponse.md), error**


## UpdateClient

Update client

### Example Usage

```go
package main

import(
	"context"
	"log"
	formancesdkgo "github.com/formancehq/formance-sdk-go"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
)

func main() {
    s := formancesdkgo.New()

    ctx := context.Background()
    res, err := s.Auth.UpdateClient(ctx, operations.UpdateClientRequest{
        UpdateClientRequest: &shared.UpdateClientRequest{
            Metadata: map[string]interface{}{
                "key": "string",
            },
            Name: "string",
            PostLogoutRedirectUris: []string{
                "string",
            },
            RedirectUris: []string{
                "string",
            },
            Scopes: []string{
                "string",
            },
        },
        ClientID: "string",
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.UpdateClientResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                        | Type                                                                             | Required                                                                         | Description                                                                      |
| -------------------------------------------------------------------------------- | -------------------------------------------------------------------------------- | -------------------------------------------------------------------------------- | -------------------------------------------------------------------------------- |
| `ctx`                                                                            | [context.Context](https://pkg.go.dev/context#Context)                            | :heavy_check_mark:                                                               | The context to use for the request.                                              |
| `request`                                                                        | [operations.UpdateClientRequest](../../models/operations/updateclientrequest.md) | :heavy_check_mark:                                                               | The request object to use for the request.                                       |


### Response

**[*operations.UpdateClientResponse](../../models/operations/updateclientresponse.md), error**

