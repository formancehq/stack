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
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/formancehq/formance-sdk-go/v2"
	"math/big"
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

    request := operations.ConfirmHoldRequest{
        ConfirmHoldRequest: &shared.ConfirmHoldRequest{
            Amount: big.NewInt(100),
            Final: v2.Bool(true),
        },
        HoldID: "<value>",
    }
    
    ctx := context.Background()
    res, err := s.Wallets.ConfirmHold(ctx, request)
    if err != nil {
        log.Fatal(err)
    }
    if res != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                          | Type                                                                               | Required                                                                           | Description                                                                        |
| ---------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------- |
| `ctx`                                                                              | [context.Context](https://pkg.go.dev/context#Context)                              | :heavy_check_mark:                                                                 | The context to use for the request.                                                |
| `request`                                                                          | [operations.ConfirmHoldRequest](../../pkg/models/operations/confirmholdrequest.md) | :heavy_check_mark:                                                                 | The request object to use for the request.                                         |


### Response

**[*operations.ConfirmHoldResponse](../../pkg/models/operations/confirmholdresponse.md), error**
| Error Object                   | Status Code                    | Content Type                   |
| ------------------------------ | ------------------------------ | ------------------------------ |
| sdkerrors.WalletsErrorResponse | default                        | application/json               |
| sdkerrors.SDKError             | 4xx-5xx                        | */*                            |

## CreateBalance

Create a balance

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

    request := operations.CreateBalanceRequest{
        ID: "<id>",
    }
    
    ctx := context.Background()
    res, err := s.Wallets.CreateBalance(ctx, request)
    if err != nil {
        log.Fatal(err)
    }
    if res.CreateBalanceResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                              | Type                                                                                   | Required                                                                               | Description                                                                            |
| -------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------- |
| `ctx`                                                                                  | [context.Context](https://pkg.go.dev/context#Context)                                  | :heavy_check_mark:                                                                     | The context to use for the request.                                                    |
| `request`                                                                              | [operations.CreateBalanceRequest](../../pkg/models/operations/createbalancerequest.md) | :heavy_check_mark:                                                                     | The request object to use for the request.                                             |


### Response

**[*operations.CreateBalanceResponse](../../pkg/models/operations/createbalanceresponse.md), error**
| Error Object                   | Status Code                    | Content Type                   |
| ------------------------------ | ------------------------------ | ------------------------------ |
| sdkerrors.WalletsErrorResponse | default                        | application/json               |
| sdkerrors.SDKError             | 4xx-5xx                        | */*                            |

## CreateWallet

Create a new wallet

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

    var request *shared.CreateWalletRequest = &shared.CreateWalletRequest{
        Metadata: map[string]string{
            "key": "<value>",
        },
        Name: "<value>",
    }
    
    ctx := context.Background()
    res, err := s.Wallets.CreateWallet(ctx, request)
    if err != nil {
        log.Fatal(err)
    }
    if res.CreateWalletResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                    | Type                                                                         | Required                                                                     | Description                                                                  |
| ---------------------------------------------------------------------------- | ---------------------------------------------------------------------------- | ---------------------------------------------------------------------------- | ---------------------------------------------------------------------------- |
| `ctx`                                                                        | [context.Context](https://pkg.go.dev/context#Context)                        | :heavy_check_mark:                                                           | The context to use for the request.                                          |
| `request`                                                                    | [shared.CreateWalletRequest](../../pkg/models/shared/createwalletrequest.md) | :heavy_check_mark:                                                           | The request object to use for the request.                                   |


### Response

**[*operations.CreateWalletResponse](../../pkg/models/operations/createwalletresponse.md), error**
| Error Object                   | Status Code                    | Content Type                   |
| ------------------------------ | ------------------------------ | ------------------------------ |
| sdkerrors.WalletsErrorResponse | default                        | application/json               |
| sdkerrors.SDKError             | 4xx-5xx                        | */*                            |

## CreditWallet

Credit a wallet

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/formancehq/formance-sdk-go/v2"
	"math/big"
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

    request := operations.CreditWalletRequest{
        CreditWalletRequest: &shared.CreditWalletRequest{
            Amount: shared.Monetary{
                Amount: big.NewInt(100),
                Asset: "USD/2",
            },
            Metadata: map[string]string{
                "key": "",
            },
            Sources: []shared.Subject{
                shared.CreateSubjectLedgerAccountSubject(
                    shared.LedgerAccountSubject{
                        Identifier: "<value>",
                        Type: "<value>",
                    },
                ),
            },
        },
        ID: "<id>",
    }
    
    ctx := context.Background()
    res, err := s.Wallets.CreditWallet(ctx, request)
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
| `request`                                                                            | [operations.CreditWalletRequest](../../pkg/models/operations/creditwalletrequest.md) | :heavy_check_mark:                                                                   | The request object to use for the request.                                           |


### Response

**[*operations.CreditWalletResponse](../../pkg/models/operations/creditwalletresponse.md), error**
| Error Object                   | Status Code                    | Content Type                   |
| ------------------------------ | ------------------------------ | ------------------------------ |
| sdkerrors.WalletsErrorResponse | default                        | application/json               |
| sdkerrors.SDKError             | 4xx-5xx                        | */*                            |

## DebitWallet

Debit a wallet

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/formancehq/formance-sdk-go/v2"
	"math/big"
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

    request := operations.DebitWalletRequest{
        DebitWalletRequest: &shared.DebitWalletRequest{
            Amount: shared.Monetary{
                Amount: big.NewInt(100),
                Asset: "USD/2",
            },
            Metadata: map[string]string{
                "key": "",
            },
            Pending: v2.Bool(true),
        },
        ID: "<id>",
    }
    
    ctx := context.Background()
    res, err := s.Wallets.DebitWallet(ctx, request)
    if err != nil {
        log.Fatal(err)
    }
    if res.DebitWalletResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                          | Type                                                                               | Required                                                                           | Description                                                                        |
| ---------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------- |
| `ctx`                                                                              | [context.Context](https://pkg.go.dev/context#Context)                              | :heavy_check_mark:                                                                 | The context to use for the request.                                                |
| `request`                                                                          | [operations.DebitWalletRequest](../../pkg/models/operations/debitwalletrequest.md) | :heavy_check_mark:                                                                 | The request object to use for the request.                                         |


### Response

**[*operations.DebitWalletResponse](../../pkg/models/operations/debitwalletresponse.md), error**
| Error Object                   | Status Code                    | Content Type                   |
| ------------------------------ | ------------------------------ | ------------------------------ |
| sdkerrors.WalletsErrorResponse | default                        | application/json               |
| sdkerrors.SDKError             | 4xx-5xx                        | */*                            |

## GetBalance

Get detailed balance

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

    request := operations.GetBalanceRequest{
        BalanceName: "<value>",
        ID: "<id>",
    }
    
    ctx := context.Background()
    res, err := s.Wallets.GetBalance(ctx, request)
    if err != nil {
        log.Fatal(err)
    }
    if res.GetBalanceResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                        | Type                                                                             | Required                                                                         | Description                                                                      |
| -------------------------------------------------------------------------------- | -------------------------------------------------------------------------------- | -------------------------------------------------------------------------------- | -------------------------------------------------------------------------------- |
| `ctx`                                                                            | [context.Context](https://pkg.go.dev/context#Context)                            | :heavy_check_mark:                                                               | The context to use for the request.                                              |
| `request`                                                                        | [operations.GetBalanceRequest](../../pkg/models/operations/getbalancerequest.md) | :heavy_check_mark:                                                               | The request object to use for the request.                                       |


### Response

**[*operations.GetBalanceResponse](../../pkg/models/operations/getbalanceresponse.md), error**
| Error Object                   | Status Code                    | Content Type                   |
| ------------------------------ | ------------------------------ | ------------------------------ |
| sdkerrors.WalletsErrorResponse | default                        | application/json               |
| sdkerrors.SDKError             | 4xx-5xx                        | */*                            |

## GetHold

Get a hold

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

    request := operations.GetHoldRequest{
        HoldID: "<value>",
    }
    
    ctx := context.Background()
    res, err := s.Wallets.GetHold(ctx, request)
    if err != nil {
        log.Fatal(err)
    }
    if res.GetHoldResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                  | Type                                                                       | Required                                                                   | Description                                                                |
| -------------------------------------------------------------------------- | -------------------------------------------------------------------------- | -------------------------------------------------------------------------- | -------------------------------------------------------------------------- |
| `ctx`                                                                      | [context.Context](https://pkg.go.dev/context#Context)                      | :heavy_check_mark:                                                         | The context to use for the request.                                        |
| `request`                                                                  | [operations.GetHoldRequest](../../pkg/models/operations/getholdrequest.md) | :heavy_check_mark:                                                         | The request object to use for the request.                                 |


### Response

**[*operations.GetHoldResponse](../../pkg/models/operations/getholdresponse.md), error**
| Error Object                   | Status Code                    | Content Type                   |
| ------------------------------ | ------------------------------ | ------------------------------ |
| sdkerrors.WalletsErrorResponse | default                        | application/json               |
| sdkerrors.SDKError             | 4xx-5xx                        | */*                            |

## GetHolds

Get all holds for a wallet

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

    request := operations.GetHoldsRequest{
        Cursor: v2.String("aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ=="),
        Metadata: map[string]string{
            "admin": "true",
        },
        PageSize: v2.Int64(100),
        WalletID: v2.String("wallet1"),
    }
    
    ctx := context.Background()
    res, err := s.Wallets.GetHolds(ctx, request)
    if err != nil {
        log.Fatal(err)
    }
    if res.GetHoldsResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                    | Type                                                                         | Required                                                                     | Description                                                                  |
| ---------------------------------------------------------------------------- | ---------------------------------------------------------------------------- | ---------------------------------------------------------------------------- | ---------------------------------------------------------------------------- |
| `ctx`                                                                        | [context.Context](https://pkg.go.dev/context#Context)                        | :heavy_check_mark:                                                           | The context to use for the request.                                          |
| `request`                                                                    | [operations.GetHoldsRequest](../../pkg/models/operations/getholdsrequest.md) | :heavy_check_mark:                                                           | The request object to use for the request.                                   |


### Response

**[*operations.GetHoldsResponse](../../pkg/models/operations/getholdsresponse.md), error**
| Error Object                   | Status Code                    | Content Type                   |
| ------------------------------ | ------------------------------ | ------------------------------ |
| sdkerrors.WalletsErrorResponse | default                        | application/json               |
| sdkerrors.SDKError             | 4xx-5xx                        | */*                            |

## GetTransactions

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

    request := operations.GetTransactionsRequest{
        Cursor: v2.String("aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ=="),
        PageSize: v2.Int64(100),
        WalletID: v2.String("wallet1"),
    }
    
    ctx := context.Background()
    res, err := s.Wallets.GetTransactions(ctx, request)
    if err != nil {
        log.Fatal(err)
    }
    if res.GetTransactionsResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                                  | Type                                                                                       | Required                                                                                   | Description                                                                                |
| ------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------ |
| `ctx`                                                                                      | [context.Context](https://pkg.go.dev/context#Context)                                      | :heavy_check_mark:                                                                         | The context to use for the request.                                                        |
| `request`                                                                                  | [operations.GetTransactionsRequest](../../pkg/models/operations/gettransactionsrequest.md) | :heavy_check_mark:                                                                         | The request object to use for the request.                                                 |


### Response

**[*operations.GetTransactionsResponse](../../pkg/models/operations/gettransactionsresponse.md), error**
| Error Object                   | Status Code                    | Content Type                   |
| ------------------------------ | ------------------------------ | ------------------------------ |
| sdkerrors.WalletsErrorResponse | default                        | application/json               |
| sdkerrors.SDKError             | 4xx-5xx                        | */*                            |

## GetWallet

Get a wallet

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

    request := operations.GetWalletRequest{
        ID: "<id>",
    }
    
    ctx := context.Background()
    res, err := s.Wallets.GetWallet(ctx, request)
    if err != nil {
        log.Fatal(err)
    }
    if res.GetWalletResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                      | Type                                                                           | Required                                                                       | Description                                                                    |
| ------------------------------------------------------------------------------ | ------------------------------------------------------------------------------ | ------------------------------------------------------------------------------ | ------------------------------------------------------------------------------ |
| `ctx`                                                                          | [context.Context](https://pkg.go.dev/context#Context)                          | :heavy_check_mark:                                                             | The context to use for the request.                                            |
| `request`                                                                      | [operations.GetWalletRequest](../../pkg/models/operations/getwalletrequest.md) | :heavy_check_mark:                                                             | The request object to use for the request.                                     |


### Response

**[*operations.GetWalletResponse](../../pkg/models/operations/getwalletresponse.md), error**
| Error Object                   | Status Code                    | Content Type                   |
| ------------------------------ | ------------------------------ | ------------------------------ |
| sdkerrors.WalletsErrorResponse | default                        | application/json               |
| sdkerrors.SDKError             | 4xx-5xx                        | */*                            |

## GetWalletSummary

Get wallet summary

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

    request := operations.GetWalletSummaryRequest{
        ID: "<id>",
    }
    
    ctx := context.Background()
    res, err := s.Wallets.GetWalletSummary(ctx, request)
    if err != nil {
        log.Fatal(err)
    }
    if res.GetWalletSummaryResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                                    | Type                                                                                         | Required                                                                                     | Description                                                                                  |
| -------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------- |
| `ctx`                                                                                        | [context.Context](https://pkg.go.dev/context#Context)                                        | :heavy_check_mark:                                                                           | The context to use for the request.                                                          |
| `request`                                                                                    | [operations.GetWalletSummaryRequest](../../pkg/models/operations/getwalletsummaryrequest.md) | :heavy_check_mark:                                                                           | The request object to use for the request.                                                   |


### Response

**[*operations.GetWalletSummaryResponse](../../pkg/models/operations/getwalletsummaryresponse.md), error**
| Error Object                   | Status Code                    | Content Type                   |
| ------------------------------ | ------------------------------ | ------------------------------ |
| sdkerrors.WalletsErrorResponse | default                        | application/json               |
| sdkerrors.SDKError             | 4xx-5xx                        | */*                            |

## ListBalances

List balances of a wallet

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

    request := operations.ListBalancesRequest{
        ID: "<id>",
    }
    
    ctx := context.Background()
    res, err := s.Wallets.ListBalances(ctx, request)
    if err != nil {
        log.Fatal(err)
    }
    if res.ListBalancesResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                            | Type                                                                                 | Required                                                                             | Description                                                                          |
| ------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------ |
| `ctx`                                                                                | [context.Context](https://pkg.go.dev/context#Context)                                | :heavy_check_mark:                                                                   | The context to use for the request.                                                  |
| `request`                                                                            | [operations.ListBalancesRequest](../../pkg/models/operations/listbalancesrequest.md) | :heavy_check_mark:                                                                   | The request object to use for the request.                                           |


### Response

**[*operations.ListBalancesResponse](../../pkg/models/operations/listbalancesresponse.md), error**
| Error Object       | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4xx-5xx            | */*                |

## ListWallets

List all wallets

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

    request := operations.ListWalletsRequest{
        Cursor: v2.String("aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ=="),
        Metadata: map[string]string{
            "admin": "true",
        },
        Name: v2.String("wallet1"),
        PageSize: v2.Int64(100),
    }
    
    ctx := context.Background()
    res, err := s.Wallets.ListWallets(ctx, request)
    if err != nil {
        log.Fatal(err)
    }
    if res.ListWalletsResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                          | Type                                                                               | Required                                                                           | Description                                                                        |
| ---------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------- |
| `ctx`                                                                              | [context.Context](https://pkg.go.dev/context#Context)                              | :heavy_check_mark:                                                                 | The context to use for the request.                                                |
| `request`                                                                          | [operations.ListWalletsRequest](../../pkg/models/operations/listwalletsrequest.md) | :heavy_check_mark:                                                                 | The request object to use for the request.                                         |


### Response

**[*operations.ListWalletsResponse](../../pkg/models/operations/listwalletsresponse.md), error**
| Error Object                   | Status Code                    | Content Type                   |
| ------------------------------ | ------------------------------ | ------------------------------ |
| sdkerrors.WalletsErrorResponse | default                        | application/json               |
| sdkerrors.SDKError             | 4xx-5xx                        | */*                            |

## UpdateWallet

Update a wallet

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

    request := operations.UpdateWalletRequest{
        ID: "<id>",
    }
    
    ctx := context.Background()
    res, err := s.Wallets.UpdateWallet(ctx, request)
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
| `request`                                                                            | [operations.UpdateWalletRequest](../../pkg/models/operations/updatewalletrequest.md) | :heavy_check_mark:                                                                   | The request object to use for the request.                                           |


### Response

**[*operations.UpdateWalletResponse](../../pkg/models/operations/updatewalletresponse.md), error**
| Error Object                   | Status Code                    | Content Type                   |
| ------------------------------ | ------------------------------ | ------------------------------ |
| sdkerrors.WalletsErrorResponse | default                        | application/json               |
| sdkerrors.SDKError             | 4xx-5xx                        | */*                            |

## VoidHold

Cancel a hold

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

    request := operations.VoidHoldRequest{
        HoldID: "<value>",
    }
    
    ctx := context.Background()
    res, err := s.Wallets.VoidHold(ctx, request)
    if err != nil {
        log.Fatal(err)
    }
    if res != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                    | Type                                                                         | Required                                                                     | Description                                                                  |
| ---------------------------------------------------------------------------- | ---------------------------------------------------------------------------- | ---------------------------------------------------------------------------- | ---------------------------------------------------------------------------- |
| `ctx`                                                                        | [context.Context](https://pkg.go.dev/context#Context)                        | :heavy_check_mark:                                                           | The context to use for the request.                                          |
| `request`                                                                    | [operations.VoidHoldRequest](../../pkg/models/operations/voidholdrequest.md) | :heavy_check_mark:                                                           | The request object to use for the request.                                   |


### Response

**[*operations.VoidHoldResponse](../../pkg/models/operations/voidholdresponse.md), error**
| Error Object                   | Status Code                    | Content Type                   |
| ------------------------------ | ------------------------------ | ------------------------------ |
| sdkerrors.WalletsErrorResponse | default                        | application/json               |
| sdkerrors.SDKError             | 4xx-5xx                        | */*                            |

## WalletsgetServerInfo

Get server info

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

**[*operations.WalletsgetServerInfoResponse](../../pkg/models/operations/walletsgetserverinforesponse.md), error**
| Error Object                   | Status Code                    | Content Type                   |
| ------------------------------ | ------------------------------ | ------------------------------ |
| sdkerrors.WalletsErrorResponse | default                        | application/json               |
| sdkerrors.SDKError             | 4xx-5xx                        | */*                            |
