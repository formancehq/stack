# Accounts

### Available Operations

* [AddMetadataToAccount](#addmetadatatoaccount) - Add metadata to an account
* [CountAccounts](#countaccounts) - Count the accounts from a ledger
* [GetAccount](#getaccount) - Get account by its address
* [ListAccounts](#listaccounts) - List accounts from a ledger

## AddMetadataToAccount

Add metadata to an account

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
    res, err := s.Accounts.AddMetadataToAccount(ctx, operations.AddMetadataToAccountRequest{
        RequestBody: map[string]interface{}{
            "provident": "distinctio",
            "quibusdam": "unde",
            "nulla": "corrupti",
        },
        Address: "users:001",
        Ledger: "ledger001",
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.StatusCode == http.StatusOK {
        // handle response
    }
}
```

## CountAccounts

Count the accounts from a ledger

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
    res, err := s.Accounts.CountAccounts(ctx, operations.CountAccountsRequest{
        Address: formance.String("users:.+"),
        Ledger: "ledger001",
        Metadata: map[string]interface{}{
            "vel": "error",
            "deserunt": "suscipit",
            "iure": "magnam",
            "debitis": "ipsa",
        },
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.StatusCode == http.StatusOK {
        // handle response
    }
}
```

## GetAccount

Get account by its address

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
    res, err := s.Accounts.GetAccount(ctx, operations.GetAccountRequest{
        Address: "users:001",
        Ledger: "ledger001",
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.AccountResponse != nil {
        // handle response
    }
}
```

## ListAccounts

List accounts from a ledger, sorted by address in descending order.

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
    res, err := s.Accounts.ListAccounts(ctx, operations.ListAccountsRequest{
        Address: formance.String("users:.+"),
        After: formance.String("users:003"),
        Balance: formance.Int64(2400),
        BalanceOperator: operations.ListAccountsBalanceOperatorGte.ToPointer(),
        Cursor: formance.String("aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ=="),
        Ledger: "ledger001",
        Metadata: map[string]interface{}{
            "tempora": "suscipit",
            "molestiae": "minus",
            "placeat": "voluptatum",
            "iusto": "excepturi",
        },
        PageSize: formance.Int64(392785),
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.AccountsCursorResponse != nil {
        // handle response
    }
}
```
