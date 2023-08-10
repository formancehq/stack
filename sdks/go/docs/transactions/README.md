# Transactions

### Available Operations

* [CreateTransactions](#createtransactions) - Create a new batch of transactions to a ledger
* [AddMetadataOnTransaction](#addmetadataontransaction) - Set the metadata of a transaction by its ID
* [CountTransactions](#counttransactions) - Count the transactions from a ledger
* [CreateTransaction](#createtransaction) - Create a new transaction to a ledger
* [GetTransaction](#gettransaction) - Get transaction from a ledger by its ID
* [ListTransactions](#listtransactions) - List transactions from a ledger
* [RevertTransaction](#reverttransaction) - Revert a ledger transaction by its ID

## CreateTransactions

Create a new batch of transactions to a ledger

### Example Usage

```go
package main

import(
	"context"
	"log"
	"github.com/formancehq/formance-sdk-go"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"math/big"
	"github.com/formancehq/formance-sdk-go/pkg/types"
)

func main() {
    s := formance.New(
        formance.WithSecurity(shared.Security{
            Authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
        }),
    )

    ctx := context.Background()
    res, err := s.Transactions.CreateTransactions(ctx, operations.CreateTransactionsRequest{
        Transactions: shared.Transactions{
            Transactions: []shared.TransactionData{
                shared.TransactionData{
                    Metadata: map[string]interface{}{
                        "amet": "deserunt",
                        "nisi": "vel",
                        "natus": "omnis",
                        "molestiae": "perferendis",
                    },
                    Postings: []shared.Posting{
                        shared.Posting{
                            Amount: big.NewInt(100),
                            Asset: "COIN",
                            Destination: "users:002",
                            Source: "users:001",
                        },
                        shared.Posting{
                            Amount: big.NewInt(100),
                            Asset: "COIN",
                            Destination: "users:002",
                            Source: "users:001",
                        },
                    },
                    Reference: formance.String("ref:001"),
                    Timestamp: types.MustTimeFromString("2022-04-14T15:11:13.227Z"),
                },
                shared.TransactionData{
                    Metadata: map[string]interface{}{
                        "labore": "labore",
                        "suscipit": "natus",
                        "nobis": "eum",
                    },
                    Postings: []shared.Posting{
                        shared.Posting{
                            Amount: big.NewInt(100),
                            Asset: "COIN",
                            Destination: "users:002",
                            Source: "users:001",
                        },
                        shared.Posting{
                            Amount: big.NewInt(100),
                            Asset: "COIN",
                            Destination: "users:002",
                            Source: "users:001",
                        },
                        shared.Posting{
                            Amount: big.NewInt(100),
                            Asset: "COIN",
                            Destination: "users:002",
                            Source: "users:001",
                        },
                        shared.Posting{
                            Amount: big.NewInt(100),
                            Asset: "COIN",
                            Destination: "users:002",
                            Source: "users:001",
                        },
                    },
                    Reference: formance.String("ref:001"),
                    Timestamp: types.MustTimeFromString("2022-11-24T10:55:00.183Z"),
                },
                shared.TransactionData{
                    Metadata: map[string]interface{}{
                        "et": "excepturi",
                        "ullam": "provident",
                    },
                    Postings: []shared.Posting{
                        shared.Posting{
                            Amount: big.NewInt(100),
                            Asset: "COIN",
                            Destination: "users:002",
                            Source: "users:001",
                        },
                        shared.Posting{
                            Amount: big.NewInt(100),
                            Asset: "COIN",
                            Destination: "users:002",
                            Source: "users:001",
                        },
                        shared.Posting{
                            Amount: big.NewInt(100),
                            Asset: "COIN",
                            Destination: "users:002",
                            Source: "users:001",
                        },
                    },
                    Reference: formance.String("ref:001"),
                    Timestamp: types.MustTimeFromString("2022-12-07T10:53:17.121Z"),
                },
            },
        },
        Ledger: "ledger001",
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.TransactionsResponse != nil {
        // handle response
    }
}
```

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
        RequestBody: map[string]interface{}{
            "reiciendis": "mollitia",
            "ad": "eum",
            "dolor": "necessitatibus",
        },
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
        EndTime: types.MustTimeFromString("2022-08-19T20:09:28.183Z"),
        Ledger: "ledger001",
        Metadata: map[string]interface{}{
            "iure": "doloribus",
        },
        Reference: formance.String("ref:001"),
        Source: formance.String("users:001"),
        StartTime: types.MustTimeFromString("2022-03-21T22:14:24.691Z"),
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
	"math/big"
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
        IdempotencyKey: formance.String("maxime"),
        PostTransaction: shared.PostTransaction{
            Metadata: map[string]interface{}{
                "facilis": "in",
                "architecto": "architecto",
                "repudiandae": "ullam",
            },
            Postings: []shared.Posting{
                shared.Posting{
                    Amount: big.NewInt(100),
                    Asset: "COIN",
                    Destination: "users:002",
                    Source: "users:001",
                },
                shared.Posting{
                    Amount: big.NewInt(100),
                    Asset: "COIN",
                    Destination: "users:002",
                    Source: "users:001",
                },
                shared.Posting{
                    Amount: big.NewInt(100),
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
                    "repellat": "quibusdam",
                    "sed": "saepe",
                },
            },
            Timestamp: types.MustTimeFromString("2022-11-20T20:56:20.791Z"),
        },
        Ledger: "ledger001",
        Preview: formance.Bool(true),
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.TransactionsResponse != nil {
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

    if res.TransactionResponse != nil {
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
        After: formance.String("consequuntur"),
        Cursor: formance.String("aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ=="),
        Destination: formance.String("users:001"),
        EndTime: types.MustTimeFromString("2021-10-08T15:23:46.576Z"),
        Ledger: "ledger001",
        Metadata: map[string]interface{}{
            "sunt": "quo",
        },
        PageSize: formance.Int64(848009),
        Reference: formance.String("ref:001"),
        Source: formance.String("users:001"),
        StartTime: types.MustTimeFromString("2020-07-30T23:39:27.609Z"),
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

    if res.TransactionResponse != nil {
        // handle response
    }
}
```
