# Auth

### Available Operations

* [AddScopeToClient](#addscopetoclient) - Add scope to client
* [AddTransientScope](#addtransientscope) - Add a transient scope to a scope
* [CreateClient](#createclient) - Create client
* [CreateScope](#createscope) - Create scope
* [CreateSecret](#createsecret) - Add a secret to a client
* [DeleteClient](#deleteclient) - Delete client
* [DeleteScope](#deletescope) - Delete scope
* [DeleteScopeFromClient](#deletescopefromclient) - Delete scope from client
* [DeleteSecret](#deletesecret) - Delete a secret from a client
* [DeleteTransientScope](#deletetransientscope) - Delete a transient scope from a scope
* [GetServerInfo](#getserverinfo) - Get server info
* [ListClients](#listclients) - List clients
* [ListScopes](#listscopes) - List scopes
* [ListUsers](#listusers) - List users
* [ReadClient](#readclient) - Read client
* [ReadScope](#readscope) - Read scope
* [ReadUser](#readuser) - Read user
* [UpdateClient](#updateclient) - Update client
* [UpdateScope](#updatescope) - Update scope

## AddScopeToClient

Add scope to client

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
    s := formance.New(
        formance.WithSecurity(shared.Security{
            Authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
        }),
    )

    ctx := context.Background()
    res, err := s.Auth.AddScopeToClient(ctx, operations.AddScopeToClientRequest{
        ClientID: "corrupti",
        ScopeID: "provident",
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.StatusCode == http.StatusOK {
        // handle response
    }
}
```

## AddTransientScope

Add a transient scope to a scope

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
    s := formance.New(
        formance.WithSecurity(shared.Security{
            Authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
        }),
    )

    ctx := context.Background()
    res, err := s.Auth.AddTransientScope(ctx, operations.AddTransientScopeRequest{
        ScopeID: "distinctio",
        TransientScopeID: "quibusdam",
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.StatusCode == http.StatusOK {
        // handle response
    }
}
```

## CreateClient

Create client

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
    s := formance.New(
        formance.WithSecurity(shared.Security{
            Authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
        }),
    )

    ctx := context.Background()
    res, err := s.Auth.CreateClient(ctx, shared.CreateClientRequest{
        Description: formance.String("unde"),
        Metadata: map[string]interface{}{
            "corrupti": "illum",
            "vel": "error",
            "deserunt": "suscipit",
            "iure": "magnam",
        },
        Name: "Larry Windler",
        PostLogoutRedirectUris: []string{
            "minus",
            "placeat",
        },
        Public: formance.Bool(false),
        RedirectUris: []string{
            "iusto",
            "excepturi",
            "nisi",
        },
        Trusted: formance.Bool(false),
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.CreateClientResponse != nil {
        // handle response
    }
}
```

## CreateScope

Create scope

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
    s := formance.New(
        formance.WithSecurity(shared.Security{
            Authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
        }),
    )

    ctx := context.Background()
    res, err := s.Auth.CreateScope(ctx, shared.CreateScopeRequest{
        Label: "recusandae",
        Metadata: map[string]interface{}{
            "ab": "quis",
            "veritatis": "deserunt",
            "perferendis": "ipsam",
            "repellendus": "sapiente",
        },
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.CreateScopeResponse != nil {
        // handle response
    }
}
```

## CreateSecret

Add a secret to a client

### Example Usage

```go
package main

import(
	"context"
	"log"
	"github.com/formancehq/formance-sdk-go"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
)

func main() {
    s := formance.New(
        formance.WithSecurity(shared.Security{
            Authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
        }),
    )

    ctx := context.Background()
    res, err := s.Auth.CreateSecret(ctx, operations.CreateSecretRequest{
        CreateSecretRequest: &shared.CreateSecretRequest{
            Metadata: map[string]interface{}{
                "odit": "at",
                "at": "maiores",
                "molestiae": "quod",
                "quod": "esse",
            },
            Name: "Miss Lowell Parisian",
        },
        ClientID: "occaecati",
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.CreateSecretResponse != nil {
        // handle response
    }
}
```

## DeleteClient

Delete client

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
    s := formance.New(
        formance.WithSecurity(shared.Security{
            Authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
        }),
    )

    ctx := context.Background()
    res, err := s.Auth.DeleteClient(ctx, operations.DeleteClientRequest{
        ClientID: "fugit",
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.StatusCode == http.StatusOK {
        // handle response
    }
}
```

## DeleteScope

Delete scope

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
    s := formance.New(
        formance.WithSecurity(shared.Security{
            Authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
        }),
    )

    ctx := context.Background()
    res, err := s.Auth.DeleteScope(ctx, operations.DeleteScopeRequest{
        ScopeID: "deleniti",
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.StatusCode == http.StatusOK {
        // handle response
    }
}
```

## DeleteScopeFromClient

Delete scope from client

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
    s := formance.New(
        formance.WithSecurity(shared.Security{
            Authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
        }),
    )

    ctx := context.Background()
    res, err := s.Auth.DeleteScopeFromClient(ctx, operations.DeleteScopeFromClientRequest{
        ClientID: "hic",
        ScopeID: "optio",
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.StatusCode == http.StatusOK {
        // handle response
    }
}
```

## DeleteSecret

Delete a secret from a client

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
    s := formance.New(
        formance.WithSecurity(shared.Security{
            Authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
        }),
    )

    ctx := context.Background()
    res, err := s.Auth.DeleteSecret(ctx, operations.DeleteSecretRequest{
        ClientID: "totam",
        SecretID: "beatae",
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.StatusCode == http.StatusOK {
        // handle response
    }
}
```

## DeleteTransientScope

Delete a transient scope from a scope

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
    s := formance.New(
        formance.WithSecurity(shared.Security{
            Authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
        }),
    )

    ctx := context.Background()
    res, err := s.Auth.DeleteTransientScope(ctx, operations.DeleteTransientScopeRequest{
        ScopeID: "commodi",
        TransientScopeID: "molestiae",
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.StatusCode == http.StatusOK {
        // handle response
    }
}
```

## GetServerInfo

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
    s := formance.New(
        formance.WithSecurity(shared.Security{
            Authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
        }),
    )

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

## ListClients

List clients

### Example Usage

```go
package main

import(
	"context"
	"log"
	"github.com/formancehq/formance-sdk-go"
)

func main() {
    s := formance.New(
        formance.WithSecurity(shared.Security{
            Authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
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

## ListScopes

List Scopes

### Example Usage

```go
package main

import(
	"context"
	"log"
	"github.com/formancehq/formance-sdk-go"
)

func main() {
    s := formance.New(
        formance.WithSecurity(shared.Security{
            Authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
        }),
    )

    ctx := context.Background()
    res, err := s.Auth.ListScopes(ctx)
    if err != nil {
        log.Fatal(err)
    }

    if res.ListScopesResponse != nil {
        // handle response
    }
}
```

## ListUsers

List users

### Example Usage

```go
package main

import(
	"context"
	"log"
	"github.com/formancehq/formance-sdk-go"
)

func main() {
    s := formance.New(
        formance.WithSecurity(shared.Security{
            Authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
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

## ReadClient

Read client

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
    s := formance.New(
        formance.WithSecurity(shared.Security{
            Authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
        }),
    )

    ctx := context.Background()
    res, err := s.Auth.ReadClient(ctx, operations.ReadClientRequest{
        ClientID: "modi",
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.ReadClientResponse != nil {
        // handle response
    }
}
```

## ReadScope

Read scope

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
    s := formance.New(
        formance.WithSecurity(shared.Security{
            Authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
        }),
    )

    ctx := context.Background()
    res, err := s.Auth.ReadScope(ctx, operations.ReadScopeRequest{
        ScopeID: "qui",
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.ReadScopeResponse != nil {
        // handle response
    }
}
```

## ReadUser

Read user

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
    s := formance.New(
        formance.WithSecurity(shared.Security{
            Authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
        }),
    )

    ctx := context.Background()
    res, err := s.Auth.ReadUser(ctx, operations.ReadUserRequest{
        UserID: "impedit",
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.ReadUserResponse != nil {
        // handle response
    }
}
```

## UpdateClient

Update client

### Example Usage

```go
package main

import(
	"context"
	"log"
	"github.com/formancehq/formance-sdk-go"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
)

func main() {
    s := formance.New(
        formance.WithSecurity(shared.Security{
            Authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
        }),
    )

    ctx := context.Background()
    res, err := s.Auth.UpdateClient(ctx, operations.UpdateClientRequest{
        UpdateClientRequest: &shared.UpdateClientRequest{
            Description: formance.String("cum"),
            Metadata: map[string]interface{}{
                "ipsum": "excepturi",
                "aspernatur": "perferendis",
            },
            Name: "Faye Cormier",
            PostLogoutRedirectUris: []string{
                "laboriosam",
                "hic",
                "saepe",
            },
            Public: formance.Bool(false),
            RedirectUris: []string{
                "in",
                "corporis",
                "iste",
            },
            Trusted: formance.Bool(false),
        },
        ClientID: "iure",
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.UpdateClientResponse != nil {
        // handle response
    }
}
```

## UpdateScope

Update scope

### Example Usage

```go
package main

import(
	"context"
	"log"
	"github.com/formancehq/formance-sdk-go"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
)

func main() {
    s := formance.New(
        formance.WithSecurity(shared.Security{
            Authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
        }),
    )

    ctx := context.Background()
    res, err := s.Auth.UpdateScope(ctx, operations.UpdateScopeRequest{
        UpdateScopeRequest: &shared.UpdateScopeRequest{
            Label: "saepe",
            Metadata: map[string]interface{}{
                "architecto": "ipsa",
                "reiciendis": "est",
                "mollitia": "laborum",
            },
        },
        ScopeID: "dolores",
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.UpdateScopeResponse != nil {
        // handle response
    }
}
```
