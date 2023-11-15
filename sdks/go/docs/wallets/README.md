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
	"math/big"
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
            Amount: big.NewInt(100),
            Final: formance.Bool(true),
        },
        HoldID: "dolores",
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
	"math/big"
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
            ExpiresAt: types.MustTimeFromString("2022-06-24T19:43:09.208Z"),
            Name: "Cynthia Hayes",
            Priority: big.NewInt(18521),
        },
        ID: "2c73d5fe-9b90-4c28-909b-3fe49a8d9cbf",
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
            "quos": "aliquid",
            "dolorem": "dolorem",
        },
        Name: "Norma Erdman",
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
	"math/big"
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
                Amount: big.NewInt(739551),
                Asset: "voluptate",
            },
            Balance: formance.String("dignissimos"),
            Metadata: map[string]string{
                "amet": "dolorum",
                "numquam": "veritatis",
                "ipsa": "ipsa",
                "iure": "odio",
            },
            Reference: formance.String("quaerat"),
            Sources: []shared.Subject{
                shared.Subject{},
                shared.Subject{},
                shared.Subject{},
                shared.Subject{},
            },
        },
        ID: "bf69280d-1ba7-47a8-9ebf-737ae4203ce5",
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
	"math/big"
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
                Amount: big.NewInt(904425),
                Asset: "suscipit",
            },
            Balances: []string{
                "provident",
                "minima",
                "repellendus",
            },
            Description: formance.String("totam"),
            Destination: &shared.Subject{},
            Metadata: map[string]string{
                "alias": "at",
                "quaerat": "tempora",
                "vel": "quod",
            },
            Pending: formance.Bool(false),
        },
        ID: "e2af7a73-cf3b-4e45-bf87-0b326b5a7342",
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
        BalanceName: "cupiditate",
        ID: "cdb1a842-2bb6-479d-a322-715bf0cbb1e3",
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
        HoldID: "veritatis",
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
            "quos": "tempore",
            "cupiditate": "aperiam",
            "delectus": "dolorem",
        },
        PageSize: formance.Int64(292147),
        WalletID: formance.String("labore"),
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
        PageSize: formance.Int64(240829),
        WalletID: formance.String("dolorum"),
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
        ID: "1108e0ad-cf4b-4921-879f-ce953f73ef7f",
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
        ID: "bc7abd74-dd39-4c0f-9d2c-ff7c70a45626",
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
        ID: "d436813f-16d9-4f5f-8e6c-556146c3e250",
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
            "libero": "aut",
            "aut": "deleniti",
            "impedit": "aliquam",
            "fugit": "accusamus",
        },
        Name: formance.String("Ellen Borer"),
        PageSize: formance.Int64(810424),
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
                "eum": "autem",
            },
        },
        ID: "c8dd6b14-4290-4747-8778-a7bd466d28c1",
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
        HoldID: "ipsa",
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
