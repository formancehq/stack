# Ledger

### Available Operations

* [AddMetadataOnTransaction](#addmetadataontransaction) - Set the metadata of a transaction by its ID
* [AddMetadataToAccount](#addmetadatatoaccount) - Add metadata to an account
* [CountAccounts](#countaccounts) - Count the accounts from a ledger
* [CountTransactions](#counttransactions) - Count the transactions from a ledger
* [CreateTransaction](#createtransaction) - Create a new transaction to a ledger
* [GetAccount](#getaccount) - Get account by its address
* [GetBalances](#getbalances) - Get the balances from a ledger's account
* [GetBalancesAggregated](#getbalancesaggregated) - Get the aggregated balances from selected accounts
* [GetInfo](#getinfo) - Show server information
* [GetLedgerInfo](#getledgerinfo) - Get information about a ledger
* [GetTransaction](#gettransaction) - Get transaction from a ledger by its ID
* [ListAccounts](#listaccounts) - List accounts from a ledger
* [ListLogs](#listlogs) - List the logs from a ledger
* [ListTransactions](#listtransactions) - List transactions from a ledger
* [ReadStats](#readstats) - Get statistics from a ledger
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
    res, err := s.Ledger.AddMetadataOnTransaction(ctx, operations.AddMetadataOnTransactionRequest{
        IdempotencyKey: formance.String("dolorem"),
        RequestBody: map[string]string{
            "explicabo": "nobis",
            "enim": "omnis",
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
    res, err := s.Ledger.AddMetadataToAccount(ctx, operations.AddMetadataToAccountRequest{
        IdempotencyKey: formance.String("nemo"),
        RequestBody: map[string]string{
            "excepturi": "accusantium",
            "iure": "culpa",
        },
        Address: "users:001",
        Async: formance.Bool(true),
        DryRun: formance.Bool(true),
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
    res, err := s.Ledger.CountAccounts(ctx, operations.CountAccountsRequest{
        Address: formance.String("users:.+"),
        Ledger: "ledger001",
        Metadata: map[string]interface{}{
            "sapiente": "architecto",
            "mollitia": "dolorem",
            "culpa": "consequuntur",
            "repellat": "mollitia",
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
    res, err := s.Ledger.CountTransactions(ctx, operations.CountTransactionsRequest{
        Account: formance.String("users:001"),
        Destination: formance.String("users:001"),
        EndTime: types.MustTimeFromString("2022-06-30T02:19:51.375Z"),
        Ledger: "ledger001",
        Metadata: map[string]string{
            "quam": "molestiae",
            "velit": "error",
        },
        Reference: formance.String("ref:001"),
        Source: formance.String("users:001"),
        StartTime: types.MustTimeFromString("2022-08-30T15:03:11.112Z"),
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
    res, err := s.Ledger.CreateTransaction(ctx, operations.CreateTransactionRequest{
        IdempotencyKey: formance.String("vitae"),
        PostTransaction: shared.PostTransaction{
            Metadata: map[string]string{
                "animi": "enim",
                "odit": "quo",
                "sequi": "tenetur",
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
                    "possimus": "aut",
                    "quasi": "error",
                    "temporibus": "laborum",
                },
            },
            Timestamp: types.MustTimeFromString("2022-01-11T05:45:42.485Z"),
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
    res, err := s.Ledger.GetAccount(ctx, operations.GetAccountRequest{
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

## GetBalances

Get the balances from a ledger's account

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
    res, err := s.Ledger.GetBalances(ctx, operations.GetBalancesRequest{
        Address: formance.String("users:001"),
        Cursor: formance.String("aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ=="),
        Ledger: "ledger001",
        PageSize: formance.Int64(976460),
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.BalancesCursorResponse != nil {
        // handle response
    }
}
```

## GetBalancesAggregated

Get the aggregated balances from selected accounts

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
    res, err := s.Ledger.GetBalancesAggregated(ctx, operations.GetBalancesAggregatedRequest{
        Address: formance.String("users:001"),
        Ledger: "ledger001",
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.AggregateBalancesResponse != nil {
        // handle response
    }
}
```

## GetInfo

Show server information

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
    res, err := s.Ledger.GetInfo(ctx)
    if err != nil {
        log.Fatal(err)
    }

    if res.ConfigInfoResponse != nil {
        // handle response
    }
}
```

## GetLedgerInfo

Get information about a ledger

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
    res, err := s.Ledger.GetLedgerInfo(ctx, operations.GetLedgerInfoRequest{
        Ledger: "ledger001",
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.LedgerInfoResponse != nil {
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
    res, err := s.Ledger.GetTransaction(ctx, operations.GetTransactionRequest{
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
    res, err := s.Ledger.ListAccounts(ctx, operations.ListAccountsRequest{
        Address: formance.String("users:.+"),
        Balance: formance.Int64(2400),
        BalanceOperator: operations.ListAccountsBalanceOperatorGte.ToPointer(),
        Cursor: formance.String("aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ=="),
        Ledger: "ledger001",
        Metadata: map[string]string{
            "nihil": "praesentium",
            "voluptatibus": "ipsa",
            "omnis": "voluptate",
            "cum": "perferendis",
        },
        PageSize: formance.Int64(39187),
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.AccountsCursorResponse != nil {
        // handle response
    }
}
```

## ListLogs

List the logs from a ledger, sorted by ID in descending order.

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
    res, err := s.Ledger.ListLogs(ctx, operations.ListLogsRequest{
        Cursor: formance.String("aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ=="),
        EndTime: types.MustTimeFromString("2022-09-19T18:36:39.009Z"),
        Ledger: "ledger001",
        PageSize: formance.Int64(979587),
        StartTime: types.MustTimeFromString("2022-08-22T19:15:58.586Z"),
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.LogsCursorResponse != nil {
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
    res, err := s.Ledger.ListTransactions(ctx, operations.ListTransactionsRequest{
        Account: formance.String("users:001"),
        Cursor: formance.String("aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ=="),
        Destination: formance.String("users:001"),
        EndTime: types.MustTimeFromString("2022-07-09T11:22:20.922Z"),
        Ledger: "ledger001",
        Metadata: map[string]string{
            "harum": "enim",
        },
        PageSize: formance.Int64(880476),
        Reference: formance.String("ref:001"),
        Source: formance.String("users:001"),
        StartTime: types.MustTimeFromString("2022-01-30T20:15:26.045Z"),
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.TransactionsCursorResponse != nil {
        // handle response
    }
}
```

## ReadStats

Get statistics from a ledger. (aggregate metrics on accounts and transactions)


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
    res, err := s.Ledger.ReadStats(ctx, operations.ReadStatsRequest{
        Ledger: "ledger001",
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.StatsResponse != nil {
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
    res, err := s.Ledger.RevertTransaction(ctx, operations.RevertTransactionRequest{
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
