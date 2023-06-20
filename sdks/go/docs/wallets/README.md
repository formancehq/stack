# Wallets

### Available Operations

* [ConfirmHold](#confirmhold) - Confirm a hold
* [CreateBalance](#createbalance) - Create a balance
* [CreateWallet](#createwallet) - Create a new wallet
* [CreditWallet](#creditwallet) - Credit a wallet
* [DebitWallet](#debitwallet) - Debit a wallet
* [GetBalance](#getbalance) - Get detailed balance
* [GetHold](#gethold) - Get a hold
* [GetHolds](#getholds) - Get all holds for a wallet
* [GetTransactions](#gettransactions)
* [GetWallet](#getwallet) - Get a wallet
* [GetWalletSummary](#getwalletsummary) - Get wallet summary
* [ListBalances](#listbalances) - List balances of a wallet
* [ListWallets](#listwallets) - List all wallets
* [UpdateWallet](#updatewallet) - Update a wallet
* [VoidHold](#voidhold) - Cancel a hold
* [WalletsgetServerInfo](#walletsgetserverinfo) - Get server info

## ConfirmHold

Confirm a hold

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
    res, err := s.Wallets.ConfirmHold(ctx, operations.ConfirmHoldRequest{
        ConfirmHoldRequest: &shared.ConfirmHoldRequest{
            Amount: formance.Int64(100),
            Final: formance.Bool(true),
        },
        HoldID: "omnis",
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.StatusCode == http.StatusOK {
        // handle response
    }
}
```

## CreateBalance

Create a balance

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
    res, err := s.Wallets.CreateBalance(ctx, operations.CreateBalanceRequest{
        CreateBalanceRequest: &shared.CreateBalanceRequest{
            ExpiresAt: types.MustTimeFromString("2022-12-24T23:52:02.245Z"),
            Name: "Megan Rau",
            Priority: formance.Int64(290077),
        },
        ID: "69b6e214-1959-4890-afa5-63e2516fe4c8",
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.CreateBalanceResponse != nil {
        // handle response
    }
}
```

## CreateWallet

Create a new wallet

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
    res, err := s.Wallets.CreateWallet(ctx, shared.CreateWalletRequest{
        Metadata: map[string]string{
            "in": "architecto",
            "architecto": "repudiandae",
            "ullam": "expedita",
        },
        Name: "Kristie Spencer",
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.CreateWalletResponse != nil {
        // handle response
    }
}
```

## CreditWallet

Credit a wallet

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
    res, err := s.Wallets.CreditWallet(ctx, operations.CreditWalletRequest{
        CreditWalletRequest: &shared.CreditWalletRequest{
            Amount: shared.Monetary{
                Amount: 868126,
                Asset: "accusantium",
            },
            Balance: formance.String("consequuntur"),
            Metadata: map[string]string{
                "natus": "magni",
                "sunt": "quo",
                "illum": "pariatur",
            },
            Reference: formance.String("maxime"),
            Sources: []shared.Subject{
                shared.Subject{},
                shared.Subject{},
            },
        },
        ID: "92601fb5-76b0-4d5f-8d30-c5fbb2587053",
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.StatusCode == http.StatusOK {
        // handle response
    }
}
```

## DebitWallet

Debit a wallet

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
    res, err := s.Wallets.DebitWallet(ctx, operations.DebitWalletRequest{
        DebitWalletRequest: &shared.DebitWalletRequest{
            Amount: shared.Monetary{
                Amount: 179490,
                Asset: "perferendis",
            },
            Balances: []string{
                "minus",
            },
            Description: formance.String("quam"),
            Destination: &shared.Subject{},
            Metadata: map[string]string{
                "vero": "nostrum",
            },
            Pending: formance.Bool(false),
        },
        ID: "fe9b90c2-8909-4b3f-a49a-8d9cbf486333",
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.DebitWalletResponse != nil {
        // handle response
    }
}
```

## GetBalance

Get detailed balance

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
    res, err := s.Wallets.GetBalance(ctx, operations.GetBalanceRequest{
        BalanceName: "qui",
        ID: "3f9b77f3-a410-4067-8ebf-69280d1ba77a",
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.GetBalanceResponse != nil {
        // handle response
    }
}
```

## GetHold

Get a hold

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
    res, err := s.Wallets.GetHold(ctx, operations.GetHoldRequest{
        HoldID: "deleniti",
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.GetHoldResponse != nil {
        // handle response
    }
}
```

## GetHolds

Get all holds for a wallet

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
    res, err := s.Wallets.GetHolds(ctx, operations.GetHoldsRequest{
        Cursor: formance.String("aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ=="),
        Metadata: map[string]string{
            "necessitatibus": "distinctio",
            "asperiores": "nihil",
            "ipsum": "voluptate",
        },
        PageSize: formance.Int64(663078),
        WalletID: formance.String("saepe"),
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.GetHoldsResponse != nil {
        // handle response
    }
}
```

## GetTransactions

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
    res, err := s.Wallets.GetTransactions(ctx, operations.GetTransactionsRequest{
        Cursor: formance.String("aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ=="),
        PageSize: formance.Int64(263322),
        WalletID: formance.String("aspernatur"),
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.GetTransactionsResponse != nil {
        // handle response
    }
}
```

## GetWallet

Get a wallet

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
    res, err := s.Wallets.GetWallet(ctx, operations.GetWalletRequest{
        ID: "03ce5e6a-95d8-4a0d-846c-e2af7a73cf3b",
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.GetWalletResponse != nil {
        // handle response
    }
}
```

## GetWalletSummary

Get wallet summary

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
    res, err := s.Wallets.GetWalletSummary(ctx, operations.GetWalletSummaryRequest{
        ID: "e453f870-b326-4b5a-b342-9cdb1a8422bb",
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.GetWalletSummaryResponse != nil {
        // handle response
    }
}
```

## ListBalances

List balances of a wallet

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
    res, err := s.Wallets.ListBalances(ctx, operations.ListBalancesRequest{
        ID: "679d2322-715b-4f0c-bb1e-31b8b90f3443",
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.ListBalancesResponse != nil {
        // handle response
    }
}
```

## ListWallets

List all wallets

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
    res, err := s.Wallets.ListWallets(ctx, operations.ListWalletsRequest{
        Cursor: formance.String("aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ=="),
        Metadata: map[string]string{
            "architecto": "quae",
            "aut": "quas",
            "itaque": "consequatur",
        },
        Name: formance.String("Marcos Schaden"),
        PageSize: formance.Int64(703495),
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.ListWalletsResponse != nil {
        // handle response
    }
}
```

## UpdateWallet

Update a wallet

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
    res, err := s.Wallets.UpdateWallet(ctx, operations.UpdateWalletRequest{
        RequestBody: &operations.UpdateWalletRequestBody{
            Metadata: map[string]string{
                "qui": "quae",
                "laudantium": "odio",
                "occaecati": "voluptatibus",
            },
        },
        ID: "ce953f73-ef7f-4bc7-abd7-4dd39c0f5d2c",
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.StatusCode == http.StatusOK {
        // handle response
    }
}
```

## VoidHold

Cancel a hold

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
    res, err := s.Wallets.VoidHold(ctx, operations.VoidHoldRequest{
        HoldID: "maiores",
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.StatusCode == http.StatusOK {
        // handle response
    }
}
```

## WalletsgetServerInfo

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
    res, err := s.Wallets.WalletsgetServerInfo(ctx)
    if err != nil {
        log.Fatal(err)
    }

    if res.ServerInfo != nil {
        // handle response
    }
}
```
