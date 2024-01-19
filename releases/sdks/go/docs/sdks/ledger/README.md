# Ledger
(*Ledger*)

### Available Operations

* [CreateTransactions](#createtransactions) - Create a new batch of transactions to a ledger
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
* [GetMapping](#getmapping) - Get the mapping of a ledger
* [GetTransaction](#gettransaction) - Get transaction from a ledger by its ID
* [ListAccounts](#listaccounts) - List accounts from a ledger
* [ListLogs](#listlogs) - List the logs from a ledger
* [ListTransactions](#listtransactions) - List transactions from a ledger
* [ReadStats](#readstats) - Get statistics from a ledger
* [RevertTransaction](#reverttransaction) - Revert a ledger transaction by its ID
* [~~RunScript~~](#runscript) - Execute a Numscript :warning: **Deprecated**
* [UpdateMapping](#updatemapping) - Update the mapping of a ledger
* [V2AddMetadataOnTransaction](#v2addmetadataontransaction) - Set the metadata of a transaction by its ID
* [V2AddMetadataToAccount](#v2addmetadatatoaccount) - Add metadata to an account
* [V2CountAccounts](#v2countaccounts) - Count the accounts from a ledger
* [V2CountTransactions](#v2counttransactions) - Count the transactions from a ledger
* [V2CreateBulk](#v2createbulk) - Bulk request
* [V2CreateLedger](#v2createledger) - Create a ledger
* [V2CreateTransaction](#v2createtransaction) - Create a new transaction to a ledger
* [V2DeleteAccountMetadata](#v2deleteaccountmetadata) - Delete metadata by key
* [V2DeleteTransactionMetadata](#v2deletetransactionmetadata) - Delete metadata by key
* [V2GetAccount](#v2getaccount) - Get account by its address
* [V2GetBalancesAggregated](#v2getbalancesaggregated) - Get the aggregated balances from selected accounts
* [V2GetInfo](#v2getinfo) - Show server information
* [V2GetLedger](#v2getledger) - Get a ledger
* [V2GetLedgerInfo](#v2getledgerinfo) - Get information about a ledger
* [V2GetTransaction](#v2gettransaction) - Get transaction from a ledger by its ID
* [V2ListAccounts](#v2listaccounts) - List accounts from a ledger
* [V2ListLedgers](#v2listledgers) - List ledgers
* [V2ListLogs](#v2listlogs) - List the logs from a ledger
* [V2ListTransactions](#v2listtransactions) - List transactions from a ledger
* [V2ReadStats](#v2readstats) - Get statistics from a ledger
* [V2RevertTransaction](#v2reverttransaction) - Revert a ledger transaction by its ID

## CreateTransactions

Create a new batch of transactions to a ledger

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/formancehq/formance-sdk-go/v2"
	"context"
	"math/big"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"log"
)

func main() {
    s := v2.New(
        v2.WithSecurity(shared.Security{
            Authorization: "Bearer <YOUR_ACCESS_TOKEN_HERE>",
        }),
    )

    ctx := context.Background()
    res, err := s.Ledger.CreateTransactions(ctx, operations.CreateTransactionsRequest{
        Transactions: shared.Transactions{
            Transactions: []shared.TransactionData{
                shared.TransactionData{
                    Metadata: map[string]interface{}{
                        "key": "string",
                    },
                    Postings: []shared.Posting{
                        shared.Posting{
                            Amount: big.NewInt(100),
                            Asset: "COIN",
                            Destination: "users:002",
                            Source: "users:001",
                        },
                    },
                    Reference: v2.String("ref:001"),
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

### Parameters

| Parameter                                                                                        | Type                                                                                             | Required                                                                                         | Description                                                                                      |
| ------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------ |
| `ctx`                                                                                            | [context.Context](https://pkg.go.dev/context#Context)                                            | :heavy_check_mark:                                                                               | The context to use for the request.                                                              |
| `request`                                                                                        | [operations.CreateTransactionsRequest](../../pkg/models/operations/createtransactionsrequest.md) | :heavy_check_mark:                                                                               | The request object to use for the request.                                                       |


### Response

**[*operations.CreateTransactionsResponse](../../pkg/models/operations/createtransactionsresponse.md), error**
| Error Object       | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4xx-5xx            | */*                |

## AddMetadataOnTransaction

Set the metadata of a transaction by its ID

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/formancehq/formance-sdk-go/v2"
	"context"
	"math/big"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"log"
	"net/http"
)

func main() {
    s := v2.New(
        v2.WithSecurity(shared.Security{
            Authorization: "Bearer <YOUR_ACCESS_TOKEN_HERE>",
        }),
    )

    ctx := context.Background()
    res, err := s.Ledger.AddMetadataOnTransaction(ctx, operations.AddMetadataOnTransactionRequest{
        RequestBody: map[string]interface{}{
            "key": "string",
        },
        Ledger: "ledger001",
        Txid: big.NewInt(1234),
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

| Parameter                                                                                                    | Type                                                                                                         | Required                                                                                                     | Description                                                                                                  |
| ------------------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------------------ |
| `ctx`                                                                                                        | [context.Context](https://pkg.go.dev/context#Context)                                                        | :heavy_check_mark:                                                                                           | The context to use for the request.                                                                          |
| `request`                                                                                                    | [operations.AddMetadataOnTransactionRequest](../../pkg/models/operations/addmetadataontransactionrequest.md) | :heavy_check_mark:                                                                                           | The request object to use for the request.                                                                   |


### Response

**[*operations.AddMetadataOnTransactionResponse](../../pkg/models/operations/addmetadataontransactionresponse.md), error**
| Error Object       | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4xx-5xx            | */*                |

## AddMetadataToAccount

Add metadata to an account

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/formancehq/formance-sdk-go/v2"
	"context"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"log"
	"net/http"
)

func main() {
    s := v2.New(
        v2.WithSecurity(shared.Security{
            Authorization: "Bearer <YOUR_ACCESS_TOKEN_HERE>",
        }),
    )

    ctx := context.Background()
    res, err := s.Ledger.AddMetadataToAccount(ctx, operations.AddMetadataToAccountRequest{
        RequestBody: map[string]interface{}{
            "key": "string",
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

### Parameters

| Parameter                                                                                            | Type                                                                                                 | Required                                                                                             | Description                                                                                          |
| ---------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------- |
| `ctx`                                                                                                | [context.Context](https://pkg.go.dev/context#Context)                                                | :heavy_check_mark:                                                                                   | The context to use for the request.                                                                  |
| `request`                                                                                            | [operations.AddMetadataToAccountRequest](../../pkg/models/operations/addmetadatatoaccountrequest.md) | :heavy_check_mark:                                                                                   | The request object to use for the request.                                                           |


### Response

**[*operations.AddMetadataToAccountResponse](../../pkg/models/operations/addmetadatatoaccountresponse.md), error**
| Error Object            | Status Code             | Content Type            |
| ----------------------- | ----------------------- | ----------------------- |
| sdkerrors.ErrorResponse | 400,404                 | application/json        |
| sdkerrors.SDKError      | 4xx-5xx                 | */*                     |

## CountAccounts

Count the accounts from a ledger

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/formancehq/formance-sdk-go/v2"
	"context"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"log"
	"net/http"
)

func main() {
    s := v2.New(
        v2.WithSecurity(shared.Security{
            Authorization: "Bearer <YOUR_ACCESS_TOKEN_HERE>",
        }),
    )

    ctx := context.Background()
    res, err := s.Ledger.CountAccounts(ctx, operations.CountAccountsRequest{
        Address: v2.String("users:.+"),
        Ledger: "ledger001",
        Metadata: map[string]interface{}{
            "key": "string",
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

### Parameters

| Parameter                                                                              | Type                                                                                   | Required                                                                               | Description                                                                            |
| -------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------- |
| `ctx`                                                                                  | [context.Context](https://pkg.go.dev/context#Context)                                  | :heavy_check_mark:                                                                     | The context to use for the request.                                                    |
| `request`                                                                              | [operations.CountAccountsRequest](../../pkg/models/operations/countaccountsrequest.md) | :heavy_check_mark:                                                                     | The request object to use for the request.                                             |


### Response

**[*operations.CountAccountsResponse](../../pkg/models/operations/countaccountsresponse.md), error**
| Error Object       | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4xx-5xx            | */*                |

## CountTransactions

Count the transactions from a ledger

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/formancehq/formance-sdk-go/v2"
	"context"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"log"
	"net/http"
)

func main() {
    s := v2.New(
        v2.WithSecurity(shared.Security{
            Authorization: "Bearer <YOUR_ACCESS_TOKEN_HERE>",
        }),
    )

    ctx := context.Background()
    res, err := s.Ledger.CountTransactions(ctx, operations.CountTransactionsRequest{
        Account: v2.String("users:001"),
        Destination: v2.String("users:001"),
        Ledger: "ledger001",
        Metadata: &operations.Metadata{},
        Reference: v2.String("ref:001"),
        Source: v2.String("users:001"),
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

| Parameter                                                                                      | Type                                                                                           | Required                                                                                       | Description                                                                                    |
| ---------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------- |
| `ctx`                                                                                          | [context.Context](https://pkg.go.dev/context#Context)                                          | :heavy_check_mark:                                                                             | The context to use for the request.                                                            |
| `request`                                                                                      | [operations.CountTransactionsRequest](../../pkg/models/operations/counttransactionsrequest.md) | :heavy_check_mark:                                                                             | The request object to use for the request.                                                     |


### Response

**[*operations.CountTransactionsResponse](../../pkg/models/operations/counttransactionsresponse.md), error**
| Error Object       | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4xx-5xx            | */*                |

## CreateTransaction

Create a new transaction to a ledger

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/formancehq/formance-sdk-go/v2"
	"context"
	"math/big"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"log"
)

func main() {
    s := v2.New(
        v2.WithSecurity(shared.Security{
            Authorization: "Bearer <YOUR_ACCESS_TOKEN_HERE>",
        }),
    )

    ctx := context.Background()
    res, err := s.Ledger.CreateTransaction(ctx, operations.CreateTransactionRequest{
        PostTransaction: shared.PostTransaction{
            Metadata: map[string]interface{}{
                "key": "string",
            },
            Postings: []shared.Posting{
                shared.Posting{
                    Amount: big.NewInt(100),
                    Asset: "COIN",
                    Destination: "users:002",
                    Source: "users:001",
                },
            },
            Reference: v2.String("ref:001"),
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
                    "user": "string",
                },
            },
        },
        Ledger: "ledger001",
        Preview: v2.Bool(true),
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.TransactionsResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                                      | Type                                                                                           | Required                                                                                       | Description                                                                                    |
| ---------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------- |
| `ctx`                                                                                          | [context.Context](https://pkg.go.dev/context#Context)                                          | :heavy_check_mark:                                                                             | The context to use for the request.                                                            |
| `request`                                                                                      | [operations.CreateTransactionRequest](../../pkg/models/operations/createtransactionrequest.md) | :heavy_check_mark:                                                                             | The request object to use for the request.                                                     |


### Response

**[*operations.CreateTransactionResponse](../../pkg/models/operations/createtransactionresponse.md), error**
| Error Object            | Status Code             | Content Type            |
| ----------------------- | ----------------------- | ----------------------- |
| sdkerrors.ErrorResponse | 400                     | application/json        |
| sdkerrors.SDKError      | 4xx-5xx                 | */*                     |

## GetAccount

Get account by its address

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/formancehq/formance-sdk-go/v2"
	"context"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"log"
)

func main() {
    s := v2.New(
        v2.WithSecurity(shared.Security{
            Authorization: "Bearer <YOUR_ACCESS_TOKEN_HERE>",
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

### Parameters

| Parameter                                                                        | Type                                                                             | Required                                                                         | Description                                                                      |
| -------------------------------------------------------------------------------- | -------------------------------------------------------------------------------- | -------------------------------------------------------------------------------- | -------------------------------------------------------------------------------- |
| `ctx`                                                                            | [context.Context](https://pkg.go.dev/context#Context)                            | :heavy_check_mark:                                                               | The context to use for the request.                                              |
| `request`                                                                        | [operations.GetAccountRequest](../../pkg/models/operations/getaccountrequest.md) | :heavy_check_mark:                                                               | The request object to use for the request.                                       |


### Response

**[*operations.GetAccountResponse](../../pkg/models/operations/getaccountresponse.md), error**
| Error Object       | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4xx-5xx            | */*                |

## GetBalances

Get the balances from a ledger's account

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/formancehq/formance-sdk-go/v2"
	"context"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"log"
)

func main() {
    s := v2.New(
        v2.WithSecurity(shared.Security{
            Authorization: "Bearer <YOUR_ACCESS_TOKEN_HERE>",
        }),
    )

    ctx := context.Background()
    res, err := s.Ledger.GetBalances(ctx, operations.GetBalancesRequest{
        Address: v2.String("users:001"),
        After: v2.String("users:003"),
        Cursor: v2.String("aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ=="),
        Ledger: "ledger001",
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.BalancesCursorResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                          | Type                                                                               | Required                                                                           | Description                                                                        |
| ---------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------- |
| `ctx`                                                                              | [context.Context](https://pkg.go.dev/context#Context)                              | :heavy_check_mark:                                                                 | The context to use for the request.                                                |
| `request`                                                                          | [operations.GetBalancesRequest](../../pkg/models/operations/getbalancesrequest.md) | :heavy_check_mark:                                                                 | The request object to use for the request.                                         |


### Response

**[*operations.GetBalancesResponse](../../pkg/models/operations/getbalancesresponse.md), error**
| Error Object       | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4xx-5xx            | */*                |

## GetBalancesAggregated

Get the aggregated balances from selected accounts

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/formancehq/formance-sdk-go/v2"
	"context"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"log"
)

func main() {
    s := v2.New(
        v2.WithSecurity(shared.Security{
            Authorization: "Bearer <YOUR_ACCESS_TOKEN_HERE>",
        }),
    )

    ctx := context.Background()
    res, err := s.Ledger.GetBalancesAggregated(ctx, operations.GetBalancesAggregatedRequest{
        Address: v2.String("users:001"),
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

### Parameters

| Parameter                                                                                              | Type                                                                                                   | Required                                                                                               | Description                                                                                            |
| ------------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------------ |
| `ctx`                                                                                                  | [context.Context](https://pkg.go.dev/context#Context)                                                  | :heavy_check_mark:                                                                                     | The context to use for the request.                                                                    |
| `request`                                                                                              | [operations.GetBalancesAggregatedRequest](../../pkg/models/operations/getbalancesaggregatedrequest.md) | :heavy_check_mark:                                                                                     | The request object to use for the request.                                                             |


### Response

**[*operations.GetBalancesAggregatedResponse](../../pkg/models/operations/getbalancesaggregatedresponse.md), error**
| Error Object       | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4xx-5xx            | */*                |

## GetInfo

Show server information

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
            Authorization: "Bearer <YOUR_ACCESS_TOKEN_HERE>",
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

### Parameters

| Parameter                                             | Type                                                  | Required                                              | Description                                           |
| ----------------------------------------------------- | ----------------------------------------------------- | ----------------------------------------------------- | ----------------------------------------------------- |
| `ctx`                                                 | [context.Context](https://pkg.go.dev/context#Context) | :heavy_check_mark:                                    | The context to use for the request.                   |


### Response

**[*operations.GetInfoResponse](../../pkg/models/operations/getinforesponse.md), error**
| Error Object       | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4xx-5xx            | */*                |

## GetLedgerInfo

Get information about a ledger

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/formancehq/formance-sdk-go/v2"
	"context"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"log"
)

func main() {
    s := v2.New(
        v2.WithSecurity(shared.Security{
            Authorization: "Bearer <YOUR_ACCESS_TOKEN_HERE>",
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

### Parameters

| Parameter                                                                              | Type                                                                                   | Required                                                                               | Description                                                                            |
| -------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------- |
| `ctx`                                                                                  | [context.Context](https://pkg.go.dev/context#Context)                                  | :heavy_check_mark:                                                                     | The context to use for the request.                                                    |
| `request`                                                                              | [operations.GetLedgerInfoRequest](../../pkg/models/operations/getledgerinforequest.md) | :heavy_check_mark:                                                                     | The request object to use for the request.                                             |


### Response

**[*operations.GetLedgerInfoResponse](../../pkg/models/operations/getledgerinforesponse.md), error**
| Error Object       | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4xx-5xx            | */*                |

## GetMapping

Get the mapping of a ledger

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/formancehq/formance-sdk-go/v2"
	"context"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"log"
)

func main() {
    s := v2.New(
        v2.WithSecurity(shared.Security{
            Authorization: "Bearer <YOUR_ACCESS_TOKEN_HERE>",
        }),
    )

    ctx := context.Background()
    res, err := s.Ledger.GetMapping(ctx, operations.GetMappingRequest{
        Ledger: "ledger001",
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.MappingResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                        | Type                                                                             | Required                                                                         | Description                                                                      |
| -------------------------------------------------------------------------------- | -------------------------------------------------------------------------------- | -------------------------------------------------------------------------------- | -------------------------------------------------------------------------------- |
| `ctx`                                                                            | [context.Context](https://pkg.go.dev/context#Context)                            | :heavy_check_mark:                                                               | The context to use for the request.                                              |
| `request`                                                                        | [operations.GetMappingRequest](../../pkg/models/operations/getmappingrequest.md) | :heavy_check_mark:                                                               | The request object to use for the request.                                       |


### Response

**[*operations.GetMappingResponse](../../pkg/models/operations/getmappingresponse.md), error**
| Error Object       | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4xx-5xx            | */*                |

## GetTransaction

Get transaction from a ledger by its ID

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/formancehq/formance-sdk-go/v2"
	"context"
	"math/big"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"log"
)

func main() {
    s := v2.New(
        v2.WithSecurity(shared.Security{
            Authorization: "Bearer <YOUR_ACCESS_TOKEN_HERE>",
        }),
    )

    ctx := context.Background()
    res, err := s.Ledger.GetTransaction(ctx, operations.GetTransactionRequest{
        Ledger: "ledger001",
        Txid: big.NewInt(1234),
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.TransactionResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                                | Type                                                                                     | Required                                                                                 | Description                                                                              |
| ---------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------- |
| `ctx`                                                                                    | [context.Context](https://pkg.go.dev/context#Context)                                    | :heavy_check_mark:                                                                       | The context to use for the request.                                                      |
| `request`                                                                                | [operations.GetTransactionRequest](../../pkg/models/operations/gettransactionrequest.md) | :heavy_check_mark:                                                                       | The request object to use for the request.                                               |


### Response

**[*operations.GetTransactionResponse](../../pkg/models/operations/gettransactionresponse.md), error**
| Error Object       | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4xx-5xx            | */*                |

## ListAccounts

List accounts from a ledger, sorted by address in descending order.

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/formancehq/formance-sdk-go/v2"
	"context"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"log"
)

func main() {
    s := v2.New(
        v2.WithSecurity(shared.Security{
            Authorization: "Bearer <YOUR_ACCESS_TOKEN_HERE>",
        }),
    )

    ctx := context.Background()
    res, err := s.Ledger.ListAccounts(ctx, operations.ListAccountsRequest{
        Address: v2.String("users:.+"),
        After: v2.String("users:003"),
        Balance: v2.Int64(2400),
        Cursor: v2.String("aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ=="),
        Ledger: "ledger001",
        Metadata: map[string]interface{}{
            "key": "string",
        },
        PaginationToken: v2.String("aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ=="),
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.AccountsCursorResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                            | Type                                                                                 | Required                                                                             | Description                                                                          |
| ------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------ |
| `ctx`                                                                                | [context.Context](https://pkg.go.dev/context#Context)                                | :heavy_check_mark:                                                                   | The context to use for the request.                                                  |
| `request`                                                                            | [operations.ListAccountsRequest](../../pkg/models/operations/listaccountsrequest.md) | :heavy_check_mark:                                                                   | The request object to use for the request.                                           |


### Response

**[*operations.ListAccountsResponse](../../pkg/models/operations/listaccountsresponse.md), error**
| Error Object       | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4xx-5xx            | */*                |

## ListLogs

List the logs from a ledger, sorted by ID in descending order.

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/formancehq/formance-sdk-go/v2"
	"context"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"log"
)

func main() {
    s := v2.New(
        v2.WithSecurity(shared.Security{
            Authorization: "Bearer <YOUR_ACCESS_TOKEN_HERE>",
        }),
    )

    ctx := context.Background()
    res, err := s.Ledger.ListLogs(ctx, operations.ListLogsRequest{
        After: v2.String("1234"),
        Cursor: v2.String("aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ=="),
        Ledger: "ledger001",
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.LogsCursorResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                    | Type                                                                         | Required                                                                     | Description                                                                  |
| ---------------------------------------------------------------------------- | ---------------------------------------------------------------------------- | ---------------------------------------------------------------------------- | ---------------------------------------------------------------------------- |
| `ctx`                                                                        | [context.Context](https://pkg.go.dev/context#Context)                        | :heavy_check_mark:                                                           | The context to use for the request.                                          |
| `request`                                                                    | [operations.ListLogsRequest](../../pkg/models/operations/listlogsrequest.md) | :heavy_check_mark:                                                           | The request object to use for the request.                                   |


### Response

**[*operations.ListLogsResponse](../../pkg/models/operations/listlogsresponse.md), error**
| Error Object       | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4xx-5xx            | */*                |

## ListTransactions

List transactions from a ledger, sorted by txid in descending order.

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/formancehq/formance-sdk-go/v2"
	"context"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"log"
)

func main() {
    s := v2.New(
        v2.WithSecurity(shared.Security{
            Authorization: "Bearer <YOUR_ACCESS_TOKEN_HERE>",
        }),
    )

    ctx := context.Background()
    res, err := s.Ledger.ListTransactions(ctx, operations.ListTransactionsRequest{
        Account: v2.String("users:001"),
        After: v2.String("1234"),
        Cursor: v2.String("aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ=="),
        Destination: v2.String("users:001"),
        Ledger: "ledger001",
        Metadata: map[string]interface{}{
            "key": "string",
        },
        Reference: v2.String("ref:001"),
        Source: v2.String("users:001"),
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.TransactionsCursorResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                                    | Type                                                                                         | Required                                                                                     | Description                                                                                  |
| -------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------- |
| `ctx`                                                                                        | [context.Context](https://pkg.go.dev/context#Context)                                        | :heavy_check_mark:                                                                           | The context to use for the request.                                                          |
| `request`                                                                                    | [operations.ListTransactionsRequest](../../pkg/models/operations/listtransactionsrequest.md) | :heavy_check_mark:                                                                           | The request object to use for the request.                                                   |


### Response

**[*operations.ListTransactionsResponse](../../pkg/models/operations/listtransactionsresponse.md), error**
| Error Object       | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4xx-5xx            | */*                |

## ReadStats

Get statistics from a ledger. (aggregate metrics on accounts and transactions)


### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/formancehq/formance-sdk-go/v2"
	"context"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"log"
)

func main() {
    s := v2.New(
        v2.WithSecurity(shared.Security{
            Authorization: "Bearer <YOUR_ACCESS_TOKEN_HERE>",
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

### Parameters

| Parameter                                                                      | Type                                                                           | Required                                                                       | Description                                                                    |
| ------------------------------------------------------------------------------ | ------------------------------------------------------------------------------ | ------------------------------------------------------------------------------ | ------------------------------------------------------------------------------ |
| `ctx`                                                                          | [context.Context](https://pkg.go.dev/context#Context)                          | :heavy_check_mark:                                                             | The context to use for the request.                                            |
| `request`                                                                      | [operations.ReadStatsRequest](../../pkg/models/operations/readstatsrequest.md) | :heavy_check_mark:                                                             | The request object to use for the request.                                     |


### Response

**[*operations.ReadStatsResponse](../../pkg/models/operations/readstatsresponse.md), error**
| Error Object       | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4xx-5xx            | */*                |

## RevertTransaction

Revert a ledger transaction by its ID

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/formancehq/formance-sdk-go/v2"
	"context"
	"math/big"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"log"
)

func main() {
    s := v2.New(
        v2.WithSecurity(shared.Security{
            Authorization: "Bearer <YOUR_ACCESS_TOKEN_HERE>",
        }),
    )

    ctx := context.Background()
    res, err := s.Ledger.RevertTransaction(ctx, operations.RevertTransactionRequest{
        Ledger: "ledger001",
        Txid: big.NewInt(1234),
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.TransactionResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                                      | Type                                                                                           | Required                                                                                       | Description                                                                                    |
| ---------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------- |
| `ctx`                                                                                          | [context.Context](https://pkg.go.dev/context#Context)                                          | :heavy_check_mark:                                                                             | The context to use for the request.                                                            |
| `request`                                                                                      | [operations.RevertTransactionRequest](../../pkg/models/operations/reverttransactionrequest.md) | :heavy_check_mark:                                                                             | The request object to use for the request.                                                     |


### Response

**[*operations.RevertTransactionResponse](../../pkg/models/operations/reverttransactionresponse.md), error**
| Error Object       | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4xx-5xx            | */*                |

## ~~RunScript~~

This route is deprecated, and has been merged into `POST /{ledger}/transactions`.


> :warning: **DEPRECATED**: This will be removed in a future release, please migrate away from it as soon as possible.

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/formancehq/formance-sdk-go/v2"
	"context"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"log"
)

func main() {
    s := v2.New(
        v2.WithSecurity(shared.Security{
            Authorization: "Bearer <YOUR_ACCESS_TOKEN_HERE>",
        }),
    )

    ctx := context.Background()
    res, err := s.Ledger.RunScript(ctx, operations.RunScriptRequest{
        Script: shared.Script{
            Metadata: map[string]interface{}{
                "key": "string",
            },
            Plain: "vars {
        account $user
        }
        send [COIN 10] (
        	source = @world
        	destination = $user
        )
        ",
            Reference: v2.String("order_1234"),
            Vars: map[string]interface{}{
                "user": "string",
            },
        },
        Ledger: "ledger001",
        Preview: v2.Bool(true),
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.ScriptResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                      | Type                                                                           | Required                                                                       | Description                                                                    |
| ------------------------------------------------------------------------------ | ------------------------------------------------------------------------------ | ------------------------------------------------------------------------------ | ------------------------------------------------------------------------------ |
| `ctx`                                                                          | [context.Context](https://pkg.go.dev/context#Context)                          | :heavy_check_mark:                                                             | The context to use for the request.                                            |
| `request`                                                                      | [operations.RunScriptRequest](../../pkg/models/operations/runscriptrequest.md) | :heavy_check_mark:                                                             | The request object to use for the request.                                     |


### Response

**[*operations.RunScriptResponse](../../pkg/models/operations/runscriptresponse.md), error**
| Error Object       | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4xx-5xx            | */*                |

## UpdateMapping

Update the mapping of a ledger

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/formancehq/formance-sdk-go/v2"
	"context"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"log"
)

func main() {
    s := v2.New(
        v2.WithSecurity(shared.Security{
            Authorization: "Bearer <YOUR_ACCESS_TOKEN_HERE>",
        }),
    )

    ctx := context.Background()
    res, err := s.Ledger.UpdateMapping(ctx, operations.UpdateMappingRequest{
        Mapping: &shared.Mapping{
            Contracts: []shared.Contract{
                shared.Contract{
                    Account: v2.String("users:001"),
                    Expr: shared.Expr{},
                },
            },
        },
        Ledger: "ledger001",
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.MappingResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                              | Type                                                                                   | Required                                                                               | Description                                                                            |
| -------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------- |
| `ctx`                                                                                  | [context.Context](https://pkg.go.dev/context#Context)                                  | :heavy_check_mark:                                                                     | The context to use for the request.                                                    |
| `request`                                                                              | [operations.UpdateMappingRequest](../../pkg/models/operations/updatemappingrequest.md) | :heavy_check_mark:                                                                     | The request object to use for the request.                                             |


### Response

**[*operations.UpdateMappingResponse](../../pkg/models/operations/updatemappingresponse.md), error**
| Error Object       | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4xx-5xx            | */*                |

## V2AddMetadataOnTransaction

Set the metadata of a transaction by its ID

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/formancehq/formance-sdk-go/v2"
	"context"
	"math/big"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"log"
	"net/http"
)

func main() {
    s := v2.New(
        v2.WithSecurity(shared.Security{
            Authorization: "Bearer <YOUR_ACCESS_TOKEN_HERE>",
        }),
    )

    ctx := context.Background()
    res, err := s.Ledger.V2AddMetadataOnTransaction(ctx, operations.V2AddMetadataOnTransactionRequest{
        RequestBody: map[string]string{
            "admin": "true",
        },
        DryRun: v2.Bool(true),
        ID: big.NewInt(1234),
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

### Parameters

| Parameter                                                                                                        | Type                                                                                                             | Required                                                                                                         | Description                                                                                                      |
| ---------------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------- |
| `ctx`                                                                                                            | [context.Context](https://pkg.go.dev/context#Context)                                                            | :heavy_check_mark:                                                                                               | The context to use for the request.                                                                              |
| `request`                                                                                                        | [operations.V2AddMetadataOnTransactionRequest](../../pkg/models/operations/v2addmetadataontransactionrequest.md) | :heavy_check_mark:                                                                                               | The request object to use for the request.                                                                       |


### Response

**[*operations.V2AddMetadataOnTransactionResponse](../../pkg/models/operations/v2addmetadataontransactionresponse.md), error**
| Error Object              | Status Code               | Content Type              |
| ------------------------- | ------------------------- | ------------------------- |
| sdkerrors.V2ErrorResponse | 400,404                   | application/json          |
| sdkerrors.SDKError        | 4xx-5xx                   | */*                       |

## V2AddMetadataToAccount

Add metadata to an account

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/formancehq/formance-sdk-go/v2"
	"context"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"log"
	"net/http"
)

func main() {
    s := v2.New(
        v2.WithSecurity(shared.Security{
            Authorization: "Bearer <YOUR_ACCESS_TOKEN_HERE>",
        }),
    )

    ctx := context.Background()
    res, err := s.Ledger.V2AddMetadataToAccount(ctx, operations.V2AddMetadataToAccountRequest{
        RequestBody: map[string]string{
            "admin": "true",
        },
        Address: "users:001",
        DryRun: v2.Bool(true),
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

### Parameters

| Parameter                                                                                                | Type                                                                                                     | Required                                                                                                 | Description                                                                                              |
| -------------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------- |
| `ctx`                                                                                                    | [context.Context](https://pkg.go.dev/context#Context)                                                    | :heavy_check_mark:                                                                                       | The context to use for the request.                                                                      |
| `request`                                                                                                | [operations.V2AddMetadataToAccountRequest](../../pkg/models/operations/v2addmetadatatoaccountrequest.md) | :heavy_check_mark:                                                                                       | The request object to use for the request.                                                               |


### Response

**[*operations.V2AddMetadataToAccountResponse](../../pkg/models/operations/v2addmetadatatoaccountresponse.md), error**
| Error Object              | Status Code               | Content Type              |
| ------------------------- | ------------------------- | ------------------------- |
| sdkerrors.V2ErrorResponse | 400,404                   | application/json          |
| sdkerrors.SDKError        | 4xx-5xx                   | */*                       |

## V2CountAccounts

Count the accounts from a ledger

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/formancehq/formance-sdk-go/v2"
	"context"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"log"
	"net/http"
)

func main() {
    s := v2.New(
        v2.WithSecurity(shared.Security{
            Authorization: "Bearer <YOUR_ACCESS_TOKEN_HERE>",
        }),
    )

    ctx := context.Background()
    res, err := s.Ledger.V2CountAccounts(ctx, operations.V2CountAccountsRequest{
        RequestBody: map[string]interface{}{
            "key": "string",
        },
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

### Parameters

| Parameter                                                                                  | Type                                                                                       | Required                                                                                   | Description                                                                                |
| ------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------ |
| `ctx`                                                                                      | [context.Context](https://pkg.go.dev/context#Context)                                      | :heavy_check_mark:                                                                         | The context to use for the request.                                                        |
| `request`                                                                                  | [operations.V2CountAccountsRequest](../../pkg/models/operations/v2countaccountsrequest.md) | :heavy_check_mark:                                                                         | The request object to use for the request.                                                 |


### Response

**[*operations.V2CountAccountsResponse](../../pkg/models/operations/v2countaccountsresponse.md), error**
| Error Object       | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4xx-5xx            | */*                |

## V2CountTransactions

Count the transactions from a ledger

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/formancehq/formance-sdk-go/v2"
	"context"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"log"
	"net/http"
)

func main() {
    s := v2.New(
        v2.WithSecurity(shared.Security{
            Authorization: "Bearer <YOUR_ACCESS_TOKEN_HERE>",
        }),
    )

    ctx := context.Background()
    res, err := s.Ledger.V2CountTransactions(ctx, operations.V2CountTransactionsRequest{
        RequestBody: map[string]interface{}{
            "key": "string",
        },
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

### Parameters

| Parameter                                                                                          | Type                                                                                               | Required                                                                                           | Description                                                                                        |
| -------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------- |
| `ctx`                                                                                              | [context.Context](https://pkg.go.dev/context#Context)                                              | :heavy_check_mark:                                                                                 | The context to use for the request.                                                                |
| `request`                                                                                          | [operations.V2CountTransactionsRequest](../../pkg/models/operations/v2counttransactionsrequest.md) | :heavy_check_mark:                                                                                 | The request object to use for the request.                                                         |


### Response

**[*operations.V2CountTransactionsResponse](../../pkg/models/operations/v2counttransactionsresponse.md), error**
| Error Object       | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4xx-5xx            | */*                |

## V2CreateBulk

Bulk request

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/formancehq/formance-sdk-go/v2"
	"context"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"log"
)

func main() {
    s := v2.New(
        v2.WithSecurity(shared.Security{
            Authorization: "Bearer <YOUR_ACCESS_TOKEN_HERE>",
        }),
    )

    ctx := context.Background()
    res, err := s.Ledger.V2CreateBulk(ctx, operations.V2CreateBulkRequest{
        RequestBody: []shared.V2BulkElement{
            shared.CreateV2BulkElementV2BulkElementAddMetadata(
                shared.V2BulkElementAddMetadata{
                    Action: "string",
                    Data: &shared.V2BulkElementAddMetadataData{
                        Metadata: map[string]string{
                            "key": "string",
                        },
                        TargetID: shared.CreateV2TargetIDStr(
                        "string",
                        ),
                        TargetType: shared.V2TargetTypeAccount,
                    },
                },
            ),
        },
        Ledger: "ledger001",
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.V2BulkResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                            | Type                                                                                 | Required                                                                             | Description                                                                          |
| ------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------ |
| `ctx`                                                                                | [context.Context](https://pkg.go.dev/context#Context)                                | :heavy_check_mark:                                                                   | The context to use for the request.                                                  |
| `request`                                                                            | [operations.V2CreateBulkRequest](../../pkg/models/operations/v2createbulkrequest.md) | :heavy_check_mark:                                                                   | The request object to use for the request.                                           |


### Response

**[*operations.V2CreateBulkResponse](../../pkg/models/operations/v2createbulkresponse.md), error**
| Error Object       | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4xx-5xx            | */*                |

## V2CreateLedger

Create a ledger

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/formancehq/formance-sdk-go/v2"
	"context"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"log"
	"net/http"
)

func main() {
    s := v2.New(
        v2.WithSecurity(shared.Security{
            Authorization: "Bearer <YOUR_ACCESS_TOKEN_HERE>",
        }),
    )

    ctx := context.Background()
    res, err := s.Ledger.V2CreateLedger(ctx, operations.V2CreateLedgerRequest{
        V2CreateLedgerRequest: &shared.V2CreateLedgerRequest{},
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

### Parameters

| Parameter                                                                                | Type                                                                                     | Required                                                                                 | Description                                                                              |
| ---------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------- |
| `ctx`                                                                                    | [context.Context](https://pkg.go.dev/context#Context)                                    | :heavy_check_mark:                                                                       | The context to use for the request.                                                      |
| `request`                                                                                | [operations.V2CreateLedgerRequest](../../pkg/models/operations/v2createledgerrequest.md) | :heavy_check_mark:                                                                       | The request object to use for the request.                                               |


### Response

**[*operations.V2CreateLedgerResponse](../../pkg/models/operations/v2createledgerresponse.md), error**
| Error Object              | Status Code               | Content Type              |
| ------------------------- | ------------------------- | ------------------------- |
| sdkerrors.V2ErrorResponse | 400                       | application/json          |
| sdkerrors.SDKError        | 4xx-5xx                   | */*                       |

## V2CreateTransaction

Create a new transaction to a ledger

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/formancehq/formance-sdk-go/v2"
	"context"
	"math/big"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"log"
)

func main() {
    s := v2.New(
        v2.WithSecurity(shared.Security{
            Authorization: "Bearer <YOUR_ACCESS_TOKEN_HERE>",
        }),
    )

    ctx := context.Background()
    res, err := s.Ledger.V2CreateTransaction(ctx, operations.V2CreateTransactionRequest{
        V2PostTransaction: shared.V2PostTransaction{
            Metadata: map[string]string{
                "admin": "true",
            },
            Postings: []shared.V2Posting{
                shared.V2Posting{
                    Amount: big.NewInt(100),
                    Asset: "COIN",
                    Destination: "users:002",
                    Source: "users:001",
                },
            },
            Reference: v2.String("ref:001"),
            Script: &shared.V2PostTransactionScript{
                Plain: "vars {
            account $user
            }
            send [COIN 10] (
            	source = @world
            	destination = $user
            )
            ",
                Vars: map[string]interface{}{
                    "user": "string",
                },
            },
        },
        DryRun: v2.Bool(true),
        Ledger: "ledger001",
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.V2CreateTransactionResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                                          | Type                                                                                               | Required                                                                                           | Description                                                                                        |
| -------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------- |
| `ctx`                                                                                              | [context.Context](https://pkg.go.dev/context#Context)                                              | :heavy_check_mark:                                                                                 | The context to use for the request.                                                                |
| `request`                                                                                          | [operations.V2CreateTransactionRequest](../../pkg/models/operations/v2createtransactionrequest.md) | :heavy_check_mark:                                                                                 | The request object to use for the request.                                                         |


### Response

**[*operations.V2CreateTransactionResponse](../../pkg/models/operations/v2createtransactionresponse.md), error**
| Error Object              | Status Code               | Content Type              |
| ------------------------- | ------------------------- | ------------------------- |
| sdkerrors.V2ErrorResponse | 400                       | application/json          |
| sdkerrors.SDKError        | 4xx-5xx                   | */*                       |

## V2DeleteAccountMetadata

Delete metadata by key

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/formancehq/formance-sdk-go/v2"
	"context"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"log"
	"net/http"
)

func main() {
    s := v2.New(
        v2.WithSecurity(shared.Security{
            Authorization: "Bearer <YOUR_ACCESS_TOKEN_HERE>",
        }),
    )

    ctx := context.Background()
    res, err := s.Ledger.V2DeleteAccountMetadata(ctx, operations.V2DeleteAccountMetadataRequest{
        Address: "3680 Emile Grove",
        Key: "foo",
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

### Parameters

| Parameter                                                                                                  | Type                                                                                                       | Required                                                                                                   | Description                                                                                                |
| ---------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------- |
| `ctx`                                                                                                      | [context.Context](https://pkg.go.dev/context#Context)                                                      | :heavy_check_mark:                                                                                         | The context to use for the request.                                                                        |
| `request`                                                                                                  | [operations.V2DeleteAccountMetadataRequest](../../pkg/models/operations/v2deleteaccountmetadatarequest.md) | :heavy_check_mark:                                                                                         | The request object to use for the request.                                                                 |


### Response

**[*operations.V2DeleteAccountMetadataResponse](../../pkg/models/operations/v2deleteaccountmetadataresponse.md), error**
| Error Object       | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4xx-5xx            | */*                |

## V2DeleteTransactionMetadata

Delete metadata by key

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/formancehq/formance-sdk-go/v2"
	"context"
	"math/big"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"log"
	"net/http"
)

func main() {
    s := v2.New(
        v2.WithSecurity(shared.Security{
            Authorization: "Bearer <YOUR_ACCESS_TOKEN_HERE>",
        }),
    )

    ctx := context.Background()
    res, err := s.Ledger.V2DeleteTransactionMetadata(ctx, operations.V2DeleteTransactionMetadataRequest{
        ID: big.NewInt(1234),
        Key: "foo",
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

### Parameters

| Parameter                                                                                                          | Type                                                                                                               | Required                                                                                                           | Description                                                                                                        |
| ------------------------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------------------------ |
| `ctx`                                                                                                              | [context.Context](https://pkg.go.dev/context#Context)                                                              | :heavy_check_mark:                                                                                                 | The context to use for the request.                                                                                |
| `request`                                                                                                          | [operations.V2DeleteTransactionMetadataRequest](../../pkg/models/operations/v2deletetransactionmetadatarequest.md) | :heavy_check_mark:                                                                                                 | The request object to use for the request.                                                                         |


### Response

**[*operations.V2DeleteTransactionMetadataResponse](../../pkg/models/operations/v2deletetransactionmetadataresponse.md), error**
| Error Object              | Status Code               | Content Type              |
| ------------------------- | ------------------------- | ------------------------- |
| sdkerrors.V2ErrorResponse | 400                       | application/json          |
| sdkerrors.SDKError        | 4xx-5xx                   | */*                       |

## V2GetAccount

Get account by its address

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/formancehq/formance-sdk-go/v2"
	"context"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"log"
)

func main() {
    s := v2.New(
        v2.WithSecurity(shared.Security{
            Authorization: "Bearer <YOUR_ACCESS_TOKEN_HERE>",
        }),
    )

    ctx := context.Background()
    res, err := s.Ledger.V2GetAccount(ctx, operations.V2GetAccountRequest{
        Address: "users:001",
        Ledger: "ledger001",
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.V2AccountResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                            | Type                                                                                 | Required                                                                             | Description                                                                          |
| ------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------ |
| `ctx`                                                                                | [context.Context](https://pkg.go.dev/context#Context)                                | :heavy_check_mark:                                                                   | The context to use for the request.                                                  |
| `request`                                                                            | [operations.V2GetAccountRequest](../../pkg/models/operations/v2getaccountrequest.md) | :heavy_check_mark:                                                                   | The request object to use for the request.                                           |


### Response

**[*operations.V2GetAccountResponse](../../pkg/models/operations/v2getaccountresponse.md), error**
| Error Object       | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4xx-5xx            | */*                |

## V2GetBalancesAggregated

Get the aggregated balances from selected accounts

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/formancehq/formance-sdk-go/v2"
	"context"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"log"
)

func main() {
    s := v2.New(
        v2.WithSecurity(shared.Security{
            Authorization: "Bearer <YOUR_ACCESS_TOKEN_HERE>",
        }),
    )

    ctx := context.Background()
    res, err := s.Ledger.V2GetBalancesAggregated(ctx, operations.V2GetBalancesAggregatedRequest{
        RequestBody: map[string]interface{}{
            "key": "string",
        },
        Ledger: "ledger001",
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.V2AggregateBalancesResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                                                  | Type                                                                                                       | Required                                                                                                   | Description                                                                                                |
| ---------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------- |
| `ctx`                                                                                                      | [context.Context](https://pkg.go.dev/context#Context)                                                      | :heavy_check_mark:                                                                                         | The context to use for the request.                                                                        |
| `request`                                                                                                  | [operations.V2GetBalancesAggregatedRequest](../../pkg/models/operations/v2getbalancesaggregatedrequest.md) | :heavy_check_mark:                                                                                         | The request object to use for the request.                                                                 |


### Response

**[*operations.V2GetBalancesAggregatedResponse](../../pkg/models/operations/v2getbalancesaggregatedresponse.md), error**
| Error Object       | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4xx-5xx            | */*                |

## V2GetInfo

Show server information

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
            Authorization: "Bearer <YOUR_ACCESS_TOKEN_HERE>",
        }),
    )

    ctx := context.Background()
    res, err := s.Ledger.V2GetInfo(ctx)
    if err != nil {
        log.Fatal(err)
    }

    if res.V2ConfigInfoResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                             | Type                                                  | Required                                              | Description                                           |
| ----------------------------------------------------- | ----------------------------------------------------- | ----------------------------------------------------- | ----------------------------------------------------- |
| `ctx`                                                 | [context.Context](https://pkg.go.dev/context#Context) | :heavy_check_mark:                                    | The context to use for the request.                   |


### Response

**[*operations.V2GetInfoResponse](../../pkg/models/operations/v2getinforesponse.md), error**
| Error Object       | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4xx-5xx            | */*                |

## V2GetLedger

Get a ledger

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/formancehq/formance-sdk-go/v2"
	"context"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"log"
)

func main() {
    s := v2.New(
        v2.WithSecurity(shared.Security{
            Authorization: "Bearer <YOUR_ACCESS_TOKEN_HERE>",
        }),
    )

    ctx := context.Background()
    res, err := s.Ledger.V2GetLedger(ctx, operations.V2GetLedgerRequest{
        Ledger: "ledger001",
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.V2Ledger != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                          | Type                                                                               | Required                                                                           | Description                                                                        |
| ---------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------- |
| `ctx`                                                                              | [context.Context](https://pkg.go.dev/context#Context)                              | :heavy_check_mark:                                                                 | The context to use for the request.                                                |
| `request`                                                                          | [operations.V2GetLedgerRequest](../../pkg/models/operations/v2getledgerrequest.md) | :heavy_check_mark:                                                                 | The request object to use for the request.                                         |


### Response

**[*operations.V2GetLedgerResponse](../../pkg/models/operations/v2getledgerresponse.md), error**
| Error Object       | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4xx-5xx            | */*                |

## V2GetLedgerInfo

Get information about a ledger

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/formancehq/formance-sdk-go/v2"
	"context"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"log"
)

func main() {
    s := v2.New(
        v2.WithSecurity(shared.Security{
            Authorization: "Bearer <YOUR_ACCESS_TOKEN_HERE>",
        }),
    )

    ctx := context.Background()
    res, err := s.Ledger.V2GetLedgerInfo(ctx, operations.V2GetLedgerInfoRequest{
        Ledger: "ledger001",
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.V2LedgerInfoResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                                  | Type                                                                                       | Required                                                                                   | Description                                                                                |
| ------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------ |
| `ctx`                                                                                      | [context.Context](https://pkg.go.dev/context#Context)                                      | :heavy_check_mark:                                                                         | The context to use for the request.                                                        |
| `request`                                                                                  | [operations.V2GetLedgerInfoRequest](../../pkg/models/operations/v2getledgerinforequest.md) | :heavy_check_mark:                                                                         | The request object to use for the request.                                                 |


### Response

**[*operations.V2GetLedgerInfoResponse](../../pkg/models/operations/v2getledgerinforesponse.md), error**
| Error Object       | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4xx-5xx            | */*                |

## V2GetTransaction

Get transaction from a ledger by its ID

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/formancehq/formance-sdk-go/v2"
	"context"
	"math/big"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"log"
)

func main() {
    s := v2.New(
        v2.WithSecurity(shared.Security{
            Authorization: "Bearer <YOUR_ACCESS_TOKEN_HERE>",
        }),
    )

    ctx := context.Background()
    res, err := s.Ledger.V2GetTransaction(ctx, operations.V2GetTransactionRequest{
        ID: big.NewInt(1234),
        Ledger: "ledger001",
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.V2GetTransactionResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                                    | Type                                                                                         | Required                                                                                     | Description                                                                                  |
| -------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------- |
| `ctx`                                                                                        | [context.Context](https://pkg.go.dev/context#Context)                                        | :heavy_check_mark:                                                                           | The context to use for the request.                                                          |
| `request`                                                                                    | [operations.V2GetTransactionRequest](../../pkg/models/operations/v2gettransactionrequest.md) | :heavy_check_mark:                                                                           | The request object to use for the request.                                                   |


### Response

**[*operations.V2GetTransactionResponse](../../pkg/models/operations/v2gettransactionresponse.md), error**
| Error Object              | Status Code               | Content Type              |
| ------------------------- | ------------------------- | ------------------------- |
| sdkerrors.V2ErrorResponse | 404                       | application/json          |
| sdkerrors.SDKError        | 4xx-5xx                   | */*                       |

## V2ListAccounts

List accounts from a ledger, sorted by address in descending order.

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/formancehq/formance-sdk-go/v2"
	"context"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"log"
)

func main() {
    s := v2.New(
        v2.WithSecurity(shared.Security{
            Authorization: "Bearer <YOUR_ACCESS_TOKEN_HERE>",
        }),
    )

    ctx := context.Background()
    res, err := s.Ledger.V2ListAccounts(ctx, operations.V2ListAccountsRequest{
        RequestBody: map[string]interface{}{
            "key": "string",
        },
        Cursor: v2.String("aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ=="),
        Ledger: "ledger001",
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.V2AccountsCursorResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                                | Type                                                                                     | Required                                                                                 | Description                                                                              |
| ---------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------- |
| `ctx`                                                                                    | [context.Context](https://pkg.go.dev/context#Context)                                    | :heavy_check_mark:                                                                       | The context to use for the request.                                                      |
| `request`                                                                                | [operations.V2ListAccountsRequest](../../pkg/models/operations/v2listaccountsrequest.md) | :heavy_check_mark:                                                                       | The request object to use for the request.                                               |


### Response

**[*operations.V2ListAccountsResponse](../../pkg/models/operations/v2listaccountsresponse.md), error**
| Error Object              | Status Code               | Content Type              |
| ------------------------- | ------------------------- | ------------------------- |
| sdkerrors.V2ErrorResponse | 400                       | application/json          |
| sdkerrors.SDKError        | 4xx-5xx                   | */*                       |

## V2ListLedgers

List ledgers

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/formancehq/formance-sdk-go/v2"
	"context"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"log"
)

func main() {
    s := v2.New(
        v2.WithSecurity(shared.Security{
            Authorization: "Bearer <YOUR_ACCESS_TOKEN_HERE>",
        }),
    )

    ctx := context.Background()
    res, err := s.Ledger.V2ListLedgers(ctx, operations.V2ListLedgersRequest{
        Cursor: v2.String("aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ=="),
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.V2LedgerListResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                              | Type                                                                                   | Required                                                                               | Description                                                                            |
| -------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------- |
| `ctx`                                                                                  | [context.Context](https://pkg.go.dev/context#Context)                                  | :heavy_check_mark:                                                                     | The context to use for the request.                                                    |
| `request`                                                                              | [operations.V2ListLedgersRequest](../../pkg/models/operations/v2listledgersrequest.md) | :heavy_check_mark:                                                                     | The request object to use for the request.                                             |


### Response

**[*operations.V2ListLedgersResponse](../../pkg/models/operations/v2listledgersresponse.md), error**
| Error Object       | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4xx-5xx            | */*                |

## V2ListLogs

List the logs from a ledger, sorted by ID in descending order.

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/formancehq/formance-sdk-go/v2"
	"context"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"log"
)

func main() {
    s := v2.New(
        v2.WithSecurity(shared.Security{
            Authorization: "Bearer <YOUR_ACCESS_TOKEN_HERE>",
        }),
    )

    ctx := context.Background()
    res, err := s.Ledger.V2ListLogs(ctx, operations.V2ListLogsRequest{
        RequestBody: map[string]interface{}{
            "key": "string",
        },
        Cursor: v2.String("aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ=="),
        Ledger: "ledger001",
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.V2LogsCursorResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                        | Type                                                                             | Required                                                                         | Description                                                                      |
| -------------------------------------------------------------------------------- | -------------------------------------------------------------------------------- | -------------------------------------------------------------------------------- | -------------------------------------------------------------------------------- |
| `ctx`                                                                            | [context.Context](https://pkg.go.dev/context#Context)                            | :heavy_check_mark:                                                               | The context to use for the request.                                              |
| `request`                                                                        | [operations.V2ListLogsRequest](../../pkg/models/operations/v2listlogsrequest.md) | :heavy_check_mark:                                                               | The request object to use for the request.                                       |


### Response

**[*operations.V2ListLogsResponse](../../pkg/models/operations/v2listlogsresponse.md), error**
| Error Object       | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4xx-5xx            | */*                |

## V2ListTransactions

List transactions from a ledger, sorted by id in descending order.

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/formancehq/formance-sdk-go/v2"
	"context"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"log"
)

func main() {
    s := v2.New(
        v2.WithSecurity(shared.Security{
            Authorization: "Bearer <YOUR_ACCESS_TOKEN_HERE>",
        }),
    )

    ctx := context.Background()
    res, err := s.Ledger.V2ListTransactions(ctx, operations.V2ListTransactionsRequest{
        RequestBody: map[string]interface{}{
            "key": "string",
        },
        Cursor: v2.String("aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ=="),
        Ledger: "ledger001",
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.V2TransactionsCursorResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                                        | Type                                                                                             | Required                                                                                         | Description                                                                                      |
| ------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------ |
| `ctx`                                                                                            | [context.Context](https://pkg.go.dev/context#Context)                                            | :heavy_check_mark:                                                                               | The context to use for the request.                                                              |
| `request`                                                                                        | [operations.V2ListTransactionsRequest](../../pkg/models/operations/v2listtransactionsrequest.md) | :heavy_check_mark:                                                                               | The request object to use for the request.                                                       |


### Response

**[*operations.V2ListTransactionsResponse](../../pkg/models/operations/v2listtransactionsresponse.md), error**
| Error Object              | Status Code               | Content Type              |
| ------------------------- | ------------------------- | ------------------------- |
| sdkerrors.V2ErrorResponse | 400,404                   | application/json          |
| sdkerrors.SDKError        | 4xx-5xx                   | */*                       |

## V2ReadStats

Get statistics from a ledger. (aggregate metrics on accounts and transactions)


### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/formancehq/formance-sdk-go/v2"
	"context"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"log"
)

func main() {
    s := v2.New(
        v2.WithSecurity(shared.Security{
            Authorization: "Bearer <YOUR_ACCESS_TOKEN_HERE>",
        }),
    )

    ctx := context.Background()
    res, err := s.Ledger.V2ReadStats(ctx, operations.V2ReadStatsRequest{
        Ledger: "ledger001",
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.V2StatsResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                          | Type                                                                               | Required                                                                           | Description                                                                        |
| ---------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------- |
| `ctx`                                                                              | [context.Context](https://pkg.go.dev/context#Context)                              | :heavy_check_mark:                                                                 | The context to use for the request.                                                |
| `request`                                                                          | [operations.V2ReadStatsRequest](../../pkg/models/operations/v2readstatsrequest.md) | :heavy_check_mark:                                                                 | The request object to use for the request.                                         |


### Response

**[*operations.V2ReadStatsResponse](../../pkg/models/operations/v2readstatsresponse.md), error**
| Error Object       | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4xx-5xx            | */*                |

## V2RevertTransaction

Revert a ledger transaction by its ID

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/formancehq/formance-sdk-go/v2"
	"context"
	"math/big"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"log"
)

func main() {
    s := v2.New(
        v2.WithSecurity(shared.Security{
            Authorization: "Bearer <YOUR_ACCESS_TOKEN_HERE>",
        }),
    )

    ctx := context.Background()
    res, err := s.Ledger.V2RevertTransaction(ctx, operations.V2RevertTransactionRequest{
        ID: big.NewInt(1234),
        Ledger: "ledger001",
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.V2RevertTransactionResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                                          | Type                                                                                               | Required                                                                                           | Description                                                                                        |
| -------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------- |
| `ctx`                                                                                              | [context.Context](https://pkg.go.dev/context#Context)                                              | :heavy_check_mark:                                                                                 | The context to use for the request.                                                                |
| `request`                                                                                          | [operations.V2RevertTransactionRequest](../../pkg/models/operations/v2reverttransactionrequest.md) | :heavy_check_mark:                                                                                 | The request object to use for the request.                                                         |


### Response

**[*operations.V2RevertTransactionResponse](../../pkg/models/operations/v2reverttransactionresponse.md), error**
| Error Object              | Status Code               | Content Type              |
| ------------------------- | ------------------------- | ------------------------- |
| sdkerrors.V2ErrorResponse | 400                       | application/json          |
| sdkerrors.SDKError        | 4xx-5xx                   | */*                       |
