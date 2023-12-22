# Wallets
(*Wallets*)

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
	formancesdkgo "github.com/formancehq/formance-sdk-go"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"math/big"
)

func main() {
    s := formancesdkgo.New()

    ctx := context.Background()
    res, err := s.Wallets.ConfirmHold(ctx, operations.ConfirmHoldRequest{
        ConfirmHoldRequest: &shared.ConfirmHoldRequest{
            Amount: big.NewInt(100),
            Final: formancesdkgo.Bool(true),
        },
        HoldID: "string",
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

| Parameter                                                                      | Type                                                                           | Required                                                                       | Description                                                                    |
| ------------------------------------------------------------------------------ | ------------------------------------------------------------------------------ | ------------------------------------------------------------------------------ | ------------------------------------------------------------------------------ |
| `ctx`                                                                          | [context.Context](https://pkg.go.dev/context#Context)                          | :heavy_check_mark:                                                             | The context to use for the request.                                            |
| `request`                                                                      | [operations.ConfirmHoldRequest](../../models/operations/confirmholdrequest.md) | :heavy_check_mark:                                                             | The request object to use for the request.                                     |


### Response

**[*operations.ConfirmHoldResponse](../../models/operations/confirmholdresponse.md), error**


## CreateBalance

Create a balance

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
    res, err := s.Wallets.CreateBalance(ctx, operations.CreateBalanceRequest{
        CreateBalanceRequest: &shared.CreateBalanceRequest{
            Name: "string",
        },
        ID: "<ID>",
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.CreateBalanceResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                          | Type                                                                               | Required                                                                           | Description                                                                        |
| ---------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------- |
| `ctx`                                                                              | [context.Context](https://pkg.go.dev/context#Context)                              | :heavy_check_mark:                                                                 | The context to use for the request.                                                |
| `request`                                                                          | [operations.CreateBalanceRequest](../../models/operations/createbalancerequest.md) | :heavy_check_mark:                                                                 | The request object to use for the request.                                         |


### Response

**[*operations.CreateBalanceResponse](../../models/operations/createbalanceresponse.md), error**


## CreateWallet

Create a new wallet

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
    res, err := s.Wallets.CreateWallet(ctx, &shared.CreateWalletRequest{
        Metadata: map[string]string{
            "key": "string",
        },
        Name: "string",
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.CreateWalletResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                | Type                                                                     | Required                                                                 | Description                                                              |
| ------------------------------------------------------------------------ | ------------------------------------------------------------------------ | ------------------------------------------------------------------------ | ------------------------------------------------------------------------ |
| `ctx`                                                                    | [context.Context](https://pkg.go.dev/context#Context)                    | :heavy_check_mark:                                                       | The context to use for the request.                                      |
| `request`                                                                | [shared.CreateWalletRequest](../../models/shared/createwalletrequest.md) | :heavy_check_mark:                                                       | The request object to use for the request.                               |


### Response

**[*operations.CreateWalletResponse](../../models/operations/createwalletresponse.md), error**


## CreditWallet

Credit a wallet

### Example Usage

```go
package main

import(
	"context"
	"log"
	formancesdkgo "github.com/formancehq/formance-sdk-go"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"math/big"
)

func main() {
    s := formancesdkgo.New()

    ctx := context.Background()
    res, err := s.Wallets.CreditWallet(ctx, operations.CreditWalletRequest{
        CreditWalletRequest: &shared.CreditWalletRequest{
            Amount: shared.Monetary{
                Amount: big.NewInt(201874),
                Asset: "string",
            },
            Metadata: map[string]string{
                "key": "string",
            },
            Sources: []shared.Subject{
                shared.CreateSubjectLedgerAccountSubject(
                    shared.LedgerAccountSubject{
                        Identifier: "string",
                        Type: "string",
                    },
                ),
            },
        },
        ID: "<ID>",
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
| `request`                                                                        | [operations.CreditWalletRequest](../../models/operations/creditwalletrequest.md) | :heavy_check_mark:                                                               | The request object to use for the request.                                       |


### Response

**[*operations.CreditWalletResponse](../../models/operations/creditwalletresponse.md), error**


## DebitWallet

Debit a wallet

### Example Usage

```go
package main

import(
	"context"
	"log"
	formancesdkgo "github.com/formancehq/formance-sdk-go"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"math/big"
)

func main() {
    s := formancesdkgo.New()

    ctx := context.Background()
    res, err := s.Wallets.DebitWallet(ctx, operations.DebitWalletRequest{
        DebitWalletRequest: &shared.DebitWalletRequest{
            Amount: shared.Monetary{
                Amount: big.NewInt(245256),
                Asset: "string",
            },
            Balances: []string{
                "string",
            },
            Destination: shared.CreateSubjectLedgerAccountSubject(
                    shared.LedgerAccountSubject{
                        Identifier: "string",
                        Type: "string",
                    },
            ),
            Metadata: map[string]string{
                "key": "string",
            },
        },
        ID: "<ID>",
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.DebitWalletResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                      | Type                                                                           | Required                                                                       | Description                                                                    |
| ------------------------------------------------------------------------------ | ------------------------------------------------------------------------------ | ------------------------------------------------------------------------------ | ------------------------------------------------------------------------------ |
| `ctx`                                                                          | [context.Context](https://pkg.go.dev/context#Context)                          | :heavy_check_mark:                                                             | The context to use for the request.                                            |
| `request`                                                                      | [operations.DebitWalletRequest](../../models/operations/debitwalletrequest.md) | :heavy_check_mark:                                                             | The request object to use for the request.                                     |


### Response

**[*operations.DebitWalletResponse](../../models/operations/debitwalletresponse.md), error**


## GetBalance

Get detailed balance

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
    res, err := s.Wallets.GetBalance(ctx, operations.GetBalanceRequest{
        BalanceName: "string",
        ID: "<ID>",
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.GetBalanceResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                    | Type                                                                         | Required                                                                     | Description                                                                  |
| ---------------------------------------------------------------------------- | ---------------------------------------------------------------------------- | ---------------------------------------------------------------------------- | ---------------------------------------------------------------------------- |
| `ctx`                                                                        | [context.Context](https://pkg.go.dev/context#Context)                        | :heavy_check_mark:                                                           | The context to use for the request.                                          |
| `request`                                                                    | [operations.GetBalanceRequest](../../models/operations/getbalancerequest.md) | :heavy_check_mark:                                                           | The request object to use for the request.                                   |


### Response

**[*operations.GetBalanceResponse](../../models/operations/getbalanceresponse.md), error**


## GetHold

Get a hold

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
    res, err := s.Wallets.GetHold(ctx, operations.GetHoldRequest{
        HoldID: "string",
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.GetHoldResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                              | Type                                                                   | Required                                                               | Description                                                            |
| ---------------------------------------------------------------------- | ---------------------------------------------------------------------- | ---------------------------------------------------------------------- | ---------------------------------------------------------------------- |
| `ctx`                                                                  | [context.Context](https://pkg.go.dev/context#Context)                  | :heavy_check_mark:                                                     | The context to use for the request.                                    |
| `request`                                                              | [operations.GetHoldRequest](../../models/operations/getholdrequest.md) | :heavy_check_mark:                                                     | The request object to use for the request.                             |


### Response

**[*operations.GetHoldResponse](../../models/operations/getholdresponse.md), error**


## GetHolds

Get all holds for a wallet

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
    res, err := s.Wallets.GetHolds(ctx, operations.GetHoldsRequest{
        Cursor: formancesdkgo.String("aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ=="),
        Metadata: map[string]string{
            "key": "string",
        },
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.GetHoldsResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                | Type                                                                     | Required                                                                 | Description                                                              |
| ------------------------------------------------------------------------ | ------------------------------------------------------------------------ | ------------------------------------------------------------------------ | ------------------------------------------------------------------------ |
| `ctx`                                                                    | [context.Context](https://pkg.go.dev/context#Context)                    | :heavy_check_mark:                                                       | The context to use for the request.                                      |
| `request`                                                                | [operations.GetHoldsRequest](../../models/operations/getholdsrequest.md) | :heavy_check_mark:                                                       | The request object to use for the request.                               |


### Response

**[*operations.GetHoldsResponse](../../models/operations/getholdsresponse.md), error**


## GetTransactions

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
    res, err := s.Wallets.GetTransactions(ctx, operations.GetTransactionsRequest{
        Cursor: formancesdkgo.String("aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ=="),
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.GetTransactionsResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                              | Type                                                                                   | Required                                                                               | Description                                                                            |
| -------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------- |
| `ctx`                                                                                  | [context.Context](https://pkg.go.dev/context#Context)                                  | :heavy_check_mark:                                                                     | The context to use for the request.                                                    |
| `request`                                                                              | [operations.GetTransactionsRequest](../../models/operations/gettransactionsrequest.md) | :heavy_check_mark:                                                                     | The request object to use for the request.                                             |


### Response

**[*operations.GetTransactionsResponse](../../models/operations/gettransactionsresponse.md), error**


## GetWallet

Get a wallet

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
    res, err := s.Wallets.GetWallet(ctx, operations.GetWalletRequest{
        ID: "<ID>",
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.GetWalletResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                  | Type                                                                       | Required                                                                   | Description                                                                |
| -------------------------------------------------------------------------- | -------------------------------------------------------------------------- | -------------------------------------------------------------------------- | -------------------------------------------------------------------------- |
| `ctx`                                                                      | [context.Context](https://pkg.go.dev/context#Context)                      | :heavy_check_mark:                                                         | The context to use for the request.                                        |
| `request`                                                                  | [operations.GetWalletRequest](../../models/operations/getwalletrequest.md) | :heavy_check_mark:                                                         | The request object to use for the request.                                 |


### Response

**[*operations.GetWalletResponse](../../models/operations/getwalletresponse.md), error**


## GetWalletSummary

Get wallet summary

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
    res, err := s.Wallets.GetWalletSummary(ctx, operations.GetWalletSummaryRequest{
        ID: "<ID>",
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.GetWalletSummaryResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                                | Type                                                                                     | Required                                                                                 | Description                                                                              |
| ---------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------- |
| `ctx`                                                                                    | [context.Context](https://pkg.go.dev/context#Context)                                    | :heavy_check_mark:                                                                       | The context to use for the request.                                                      |
| `request`                                                                                | [operations.GetWalletSummaryRequest](../../models/operations/getwalletsummaryrequest.md) | :heavy_check_mark:                                                                       | The request object to use for the request.                                               |


### Response

**[*operations.GetWalletSummaryResponse](../../models/operations/getwalletsummaryresponse.md), error**


## ListBalances

List balances of a wallet

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
    res, err := s.Wallets.ListBalances(ctx, operations.ListBalancesRequest{
        ID: "<ID>",
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.ListBalancesResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                        | Type                                                                             | Required                                                                         | Description                                                                      |
| -------------------------------------------------------------------------------- | -------------------------------------------------------------------------------- | -------------------------------------------------------------------------------- | -------------------------------------------------------------------------------- |
| `ctx`                                                                            | [context.Context](https://pkg.go.dev/context#Context)                            | :heavy_check_mark:                                                               | The context to use for the request.                                              |
| `request`                                                                        | [operations.ListBalancesRequest](../../models/operations/listbalancesrequest.md) | :heavy_check_mark:                                                               | The request object to use for the request.                                       |


### Response

**[*operations.ListBalancesResponse](../../models/operations/listbalancesresponse.md), error**


## ListWallets

List all wallets

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
    res, err := s.Wallets.ListWallets(ctx, operations.ListWalletsRequest{
        Cursor: formancesdkgo.String("aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ=="),
        Metadata: map[string]string{
            "key": "string",
        },
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.ListWalletsResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                      | Type                                                                           | Required                                                                       | Description                                                                    |
| ------------------------------------------------------------------------------ | ------------------------------------------------------------------------------ | ------------------------------------------------------------------------------ | ------------------------------------------------------------------------------ |
| `ctx`                                                                          | [context.Context](https://pkg.go.dev/context#Context)                          | :heavy_check_mark:                                                             | The context to use for the request.                                            |
| `request`                                                                      | [operations.ListWalletsRequest](../../models/operations/listwalletsrequest.md) | :heavy_check_mark:                                                             | The request object to use for the request.                                     |


### Response

**[*operations.ListWalletsResponse](../../models/operations/listwalletsresponse.md), error**


## UpdateWallet

Update a wallet

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
    res, err := s.Wallets.UpdateWallet(ctx, operations.UpdateWalletRequest{
        RequestBody: &operations.UpdateWalletRequestBody{
            Metadata: map[string]string{
                "key": "string",
            },
        },
        ID: "<ID>",
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
| `request`                                                                        | [operations.UpdateWalletRequest](../../models/operations/updatewalletrequest.md) | :heavy_check_mark:                                                               | The request object to use for the request.                                       |


### Response

**[*operations.UpdateWalletResponse](../../models/operations/updatewalletresponse.md), error**


## VoidHold

Cancel a hold

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
    res, err := s.Wallets.VoidHold(ctx, operations.VoidHoldRequest{
        HoldID: "string",
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

| Parameter                                                                | Type                                                                     | Required                                                                 | Description                                                              |
| ------------------------------------------------------------------------ | ------------------------------------------------------------------------ | ------------------------------------------------------------------------ | ------------------------------------------------------------------------ |
| `ctx`                                                                    | [context.Context](https://pkg.go.dev/context#Context)                    | :heavy_check_mark:                                                       | The context to use for the request.                                      |
| `request`                                                                | [operations.VoidHoldRequest](../../models/operations/voidholdrequest.md) | :heavy_check_mark:                                                       | The request object to use for the request.                               |


### Response

**[*operations.VoidHoldResponse](../../models/operations/voidholdresponse.md), error**


## WalletsgetServerInfo

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
    res, err := s.Wallets.WalletsgetServerInfo(ctx)
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

**[*operations.WalletsgetServerInfoResponse](../../models/operations/walletsgetserverinforesponse.md), error**

