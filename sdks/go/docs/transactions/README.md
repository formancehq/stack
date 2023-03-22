# Transactions

### Available Operations

* [AddMetadataOnTransaction](#addmetadataontransaction) - Set the metadata of a transaction by its ID
* [CountTransactions](#counttransactions) - Count the transactions from a ledger
* [CreateTransaction](#createtransaction) - Create a new transaction to a ledger
* [GetTransaction](#gettransaction) - Get transaction from a ledger by its ID
* [ListTransactions](#listtransactions) - List transactions from a ledger
* [RevertTransaction](#reverttransaction) - Revert a ledger transaction by its ID

## AddMetadataOnTransaction

Set the metadata of a transaction by its ID

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
    res, err := s.Transactions.AddMetadataOnTransaction(ctx, operations.AddMetadataOnTransactionRequest{
        IdempotencyKey: formance.String("incidunt"),
        RequestBody: map[string]string{
            "consequatur": "est",
            "quibusdam": "explicabo",
        },
        Async: formance.Bool(true),
        DryRun: formance.Bool(true),
        Ledger: "ledger001",
        Txid: 1234,
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.StatusCode == http.StatusOK {
        // handle response
    }
}
```

## CountTransactions

Count the transactions from a ledger

### Example Usage

```go
package main

import(
	"context"
	"log"
	"github.com/formancehq/formance-sdk-go"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/types"
)

func main() {
    s := formance.New(
        formance.WithSecurity(shared.Security{
            Authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
        }),
    )

    ctx := context.Background()
    res, err := s.Transactions.CountTransactions(ctx, operations.CountTransactionsRequest{
        Account: formance.String("users:001"),
        Destination: formance.String("users:001"),
        EndTime: types.MustTimeFromString("2021-07-27T01:56:50.693Z"),
        Ledger: "ledger001",
        Metadata: map[string]string{
            "labore": "modi",
            "qui": "aliquid",
            "cupiditate": "quos",
            "perferendis": "magni",
        },
        Reference: formance.String("ref:001"),
        Source: formance.String("users:001"),
        StartTime: types.MustTimeFromString("2021-11-22T01:26:35.048Z"),
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.StatusCode == http.StatusOK {
        // handle response
    }
}
```

## CreateTransaction

Create a new transaction to a ledger

### Example Usage

```go
package main

import(
	"context"
	"log"
	"github.com/formancehq/formance-sdk-go"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/formancehq/formance-sdk-go/pkg/types"
)

func main() {
    s := formance.New(
        formance.WithSecurity(shared.Security{
            Authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
        }),
    )

    ctx := context.Background()
    res, err := s.Transactions.CreateTransaction(ctx, operations.CreateTransactionRequest{
        IdempotencyKey: formance.String("alias"),
        PostTransaction: shared.PostTransaction{
            Metadata: map[string]string{
                "dolorum": "excepturi",
            },
            Postings: []shared.Posting{
                shared.Posting{
                    Amount: 100,
                    Asset: "COIN",
                    Destination: "users:002",
                    Source: "users:001",
                },
                shared.Posting{
                    Amount: 100,
                    Asset: "COIN",
                    Destination: "users:002",
                    Source: "users:001",
                },
            },
            Reference: formance.String("ref:001"),
            Script: &shared.PostTransactionScript{
                Plain: "vars {
            account $user
            }
            send [COIN 10] (
            	source = @world
            	destination = $user
            )
            ",
                Vars: map[string]interface{}{
                    "tempore": "labore",
                    "delectus": "eum",
                    "non": "eligendi",
                },
            },
            Timestamp: types.MustTimeFromString("2022-03-17T20:21:28.792Z"),
        },
        Async: formance.Bool(true),
        DryRun: formance.Bool(true),
        Ledger: "ledger001",
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.CreateTransactionResponse != nil {
        // handle response
    }
}
```

## GetTransaction

Get transaction from a ledger by its ID

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
    res, err := s.Transactions.GetTransaction(ctx, operations.GetTransactionRequest{
        Ledger: "ledger001",
        Txid: 1234,
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.GetTransactionResponse != nil {
        // handle response
    }
}
```

## ListTransactions

List transactions from a ledger, sorted by txid in descending order.

### Example Usage

```go
package main

import(
	"context"
	"log"
	"github.com/formancehq/formance-sdk-go"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/types"
)

func main() {
    s := formance.New(
        formance.WithSecurity(shared.Security{
            Authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
        }),
    )

    ctx := context.Background()
    res, err := s.Transactions.ListTransactions(ctx, operations.ListTransactionsRequest{
        Account: formance.String("users:001"),
        Cursor: formance.String("aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ=="),
        Destination: formance.String("users:001"),
        EndTime: types.MustTimeFromString("2021-03-17T21:24:26.606Z"),
        Ledger: "ledger001",
        Metadata: map[string]string{
            "officia": "dolor",
            "debitis": "a",
            "dolorum": "in",
        },
        PageSize: formance.Int64(449198),
        Reference: formance.String("ref:001"),
        Source: formance.String("users:001"),
        StartTime: types.MustTimeFromString("2020-01-25T11:09:22.009Z"),
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.TransactionsCursorResponse != nil {
        // handle response
    }
}
```

## RevertTransaction

Revert a ledger transaction by its ID

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
    res, err := s.Transactions.RevertTransaction(ctx, operations.RevertTransactionRequest{
        Ledger: "ledger001",
        Txid: 1234,
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.RevertTransactionResponse != nil {
        // handle response
    }
}
```
