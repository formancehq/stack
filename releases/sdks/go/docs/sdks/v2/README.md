# V2
(*Ledger.V2*)

### Available Operations

* [AddMetadataOnTransaction](#addmetadataontransaction) - Set the metadata of a transaction by its ID
* [AddMetadataToAccount](#addmetadatatoaccount) - Add metadata to an account
* [CountAccounts](#countaccounts) - Count the accounts from a ledger
* [CountTransactions](#counttransactions) - Count the transactions from a ledger
* [CreateBulk](#createbulk) - Bulk request
* [CreateLedger](#createledger) - Create a ledger
* [CreateTransaction](#createtransaction) - Create a new transaction to a ledger
* [DeleteAccountMetadata](#deleteaccountmetadata) - Delete metadata by key
* [DeleteLedgerMetadata](#deleteledgermetadata) - Delete ledger metadata by key
* [DeleteTransactionMetadata](#deletetransactionmetadata) - Delete metadata by key
* [ExportLogs](#exportlogs) - Export logs
* [GetAccount](#getaccount) - Get account by its address
* [GetBalancesAggregated](#getbalancesaggregated) - Get the aggregated balances from selected accounts
* [GetInfo](#getinfo) - Show server information
* [GetLedger](#getledger) - Get a ledger
* [GetLedgerInfo](#getledgerinfo) - Get information about a ledger
* [GetTransaction](#gettransaction) - Get transaction from a ledger by its ID
* [GetVolumesWithBalances](#getvolumeswithbalances) - Get list of volumes with balances for (account/asset)
* [ImportLogs](#importlogs)
* [ListAccounts](#listaccounts) - List accounts from a ledger
* [ListLedgers](#listledgers) - List ledgers
* [ListLogs](#listlogs) - List the logs from a ledger
* [ListTransactions](#listtransactions) - List transactions from a ledger
* [ReadStats](#readstats) - Get statistics from a ledger
* [RevertTransaction](#reverttransaction) - Revert a ledger transaction by its ID
* [UpdateLedgerMetadata](#updateledgermetadata) - Update ledger metadata

## AddMetadataOnTransaction

Set the metadata of a transaction by its ID

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"os"
	"github.com/formancehq/formance-sdk-go/v2"
	"math/big"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"context"
	"log"
)

func main() {
    s := v2.New(
        v2.WithSecurity(shared.Security{
            Authorization: os.Getenv("AUTHORIZATION"),
        }),
    )
    request := operations.V2AddMetadataOnTransactionRequest{
        RequestBody: map[string]string{
            "admin": "true",
        },
        DryRun: v2.Bool(true),
        ID: big.NewInt(1234),
        Ledger: "ledger001",
    }
    ctx := context.Background()
    res, err := s.Ledger.V2.AddMetadataOnTransaction(ctx, request)
    if err != nil {
        log.Fatal(err)
    }
    if res != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                                                        | Type                                                                                                             | Required                                                                                                         | Description                                                                                                      |
| ---------------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------- |
| `ctx`                                                                                                            | [context.Context](https://pkg.go.dev/context#Context)                                                            | :heavy_check_mark:                                                                                               | The context to use for the request.                                                                              |
| `request`                                                                                                        | [operations.V2AddMetadataOnTransactionRequest](../../pkg/models/operations/v2addmetadataontransactionrequest.md) | :heavy_check_mark:                                                                                               | The request object to use for the request.                                                                       |
| `opts`                                                                                                           | [][operations.Option](../../pkg/models/operations/option.md)                                                     | :heavy_minus_sign:                                                                                               | The options for this request.                                                                                    |


### Response

**[*operations.V2AddMetadataOnTransactionResponse](../../pkg/models/operations/v2addmetadataontransactionresponse.md), error**
| Error Object              | Status Code               | Content Type              |
| ------------------------- | ------------------------- | ------------------------- |
| sdkerrors.V2ErrorResponse | default                   | application/json          |
| sdkerrors.SDKError        | 4xx-5xx                   | */*                       |

## AddMetadataToAccount

Add metadata to an account

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"os"
	"github.com/formancehq/formance-sdk-go/v2"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"context"
	"log"
)

func main() {
    s := v2.New(
        v2.WithSecurity(shared.Security{
            Authorization: os.Getenv("AUTHORIZATION"),
        }),
    )
    request := operations.V2AddMetadataToAccountRequest{
        RequestBody: map[string]string{
            "admin": "true",
        },
        Address: "users:001",
        DryRun: v2.Bool(true),
        Ledger: "ledger001",
    }
    ctx := context.Background()
    res, err := s.Ledger.V2.AddMetadataToAccount(ctx, request)
    if err != nil {
        log.Fatal(err)
    }
    if res != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                                                | Type                                                                                                     | Required                                                                                                 | Description                                                                                              |
| -------------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------- |
| `ctx`                                                                                                    | [context.Context](https://pkg.go.dev/context#Context)                                                    | :heavy_check_mark:                                                                                       | The context to use for the request.                                                                      |
| `request`                                                                                                | [operations.V2AddMetadataToAccountRequest](../../pkg/models/operations/v2addmetadatatoaccountrequest.md) | :heavy_check_mark:                                                                                       | The request object to use for the request.                                                               |
| `opts`                                                                                                   | [][operations.Option](../../pkg/models/operations/option.md)                                             | :heavy_minus_sign:                                                                                       | The options for this request.                                                                            |


### Response

**[*operations.V2AddMetadataToAccountResponse](../../pkg/models/operations/v2addmetadatatoaccountresponse.md), error**
| Error Object              | Status Code               | Content Type              |
| ------------------------- | ------------------------- | ------------------------- |
| sdkerrors.V2ErrorResponse | default                   | application/json          |
| sdkerrors.SDKError        | 4xx-5xx                   | */*                       |

## CountAccounts

Count the accounts from a ledger

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"os"
	"github.com/formancehq/formance-sdk-go/v2"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"context"
	"log"
)

func main() {
    s := v2.New(
        v2.WithSecurity(shared.Security{
            Authorization: os.Getenv("AUTHORIZATION"),
        }),
    )
    request := operations.V2CountAccountsRequest{
        Ledger: "ledger001",
    }
    ctx := context.Background()
    res, err := s.Ledger.V2.CountAccounts(ctx, request)
    if err != nil {
        log.Fatal(err)
    }
    if res != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                                  | Type                                                                                       | Required                                                                                   | Description                                                                                |
| ------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------ |
| `ctx`                                                                                      | [context.Context](https://pkg.go.dev/context#Context)                                      | :heavy_check_mark:                                                                         | The context to use for the request.                                                        |
| `request`                                                                                  | [operations.V2CountAccountsRequest](../../pkg/models/operations/v2countaccountsrequest.md) | :heavy_check_mark:                                                                         | The request object to use for the request.                                                 |
| `opts`                                                                                     | [][operations.Option](../../pkg/models/operations/option.md)                               | :heavy_minus_sign:                                                                         | The options for this request.                                                              |


### Response

**[*operations.V2CountAccountsResponse](../../pkg/models/operations/v2countaccountsresponse.md), error**
| Error Object              | Status Code               | Content Type              |
| ------------------------- | ------------------------- | ------------------------- |
| sdkerrors.V2ErrorResponse | default                   | application/json          |
| sdkerrors.SDKError        | 4xx-5xx                   | */*                       |

## CountTransactions

Count the transactions from a ledger

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"os"
	"github.com/formancehq/formance-sdk-go/v2"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"context"
	"log"
)

func main() {
    s := v2.New(
        v2.WithSecurity(shared.Security{
            Authorization: os.Getenv("AUTHORIZATION"),
        }),
    )
    request := operations.V2CountTransactionsRequest{
        Ledger: "ledger001",
    }
    ctx := context.Background()
    res, err := s.Ledger.V2.CountTransactions(ctx, request)
    if err != nil {
        log.Fatal(err)
    }
    if res != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                                          | Type                                                                                               | Required                                                                                           | Description                                                                                        |
| -------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------- |
| `ctx`                                                                                              | [context.Context](https://pkg.go.dev/context#Context)                                              | :heavy_check_mark:                                                                                 | The context to use for the request.                                                                |
| `request`                                                                                          | [operations.V2CountTransactionsRequest](../../pkg/models/operations/v2counttransactionsrequest.md) | :heavy_check_mark:                                                                                 | The request object to use for the request.                                                         |
| `opts`                                                                                             | [][operations.Option](../../pkg/models/operations/option.md)                                       | :heavy_minus_sign:                                                                                 | The options for this request.                                                                      |


### Response

**[*operations.V2CountTransactionsResponse](../../pkg/models/operations/v2counttransactionsresponse.md), error**
| Error Object              | Status Code               | Content Type              |
| ------------------------- | ------------------------- | ------------------------- |
| sdkerrors.V2ErrorResponse | default                   | application/json          |
| sdkerrors.SDKError        | 4xx-5xx                   | */*                       |

## CreateBulk

Bulk request

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"os"
	"github.com/formancehq/formance-sdk-go/v2"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"context"
	"log"
)

func main() {
    s := v2.New(
        v2.WithSecurity(shared.Security{
            Authorization: os.Getenv("AUTHORIZATION"),
        }),
    )
    request := operations.V2CreateBulkRequest{
        RequestBody: []shared.V2BulkElement{
            shared.CreateV2BulkElementV2BulkElementAddMetadata(
                shared.V2BulkElementAddMetadata{
                    Action: "<value>",
                },
            ),
        },
        Ledger: "ledger001",
    }
    ctx := context.Background()
    res, err := s.Ledger.V2.CreateBulk(ctx, request)
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
| `opts`                                                                               | [][operations.Option](../../pkg/models/operations/option.md)                         | :heavy_minus_sign:                                                                   | The options for this request.                                                        |


### Response

**[*operations.V2CreateBulkResponse](../../pkg/models/operations/v2createbulkresponse.md), error**
| Error Object              | Status Code               | Content Type              |
| ------------------------- | ------------------------- | ------------------------- |
| sdkerrors.V2ErrorResponse | default                   | application/json          |
| sdkerrors.SDKError        | 4xx-5xx                   | */*                       |

## CreateLedger

Create a ledger

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"os"
	"github.com/formancehq/formance-sdk-go/v2"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"context"
	"log"
)

func main() {
    s := v2.New(
        v2.WithSecurity(shared.Security{
            Authorization: os.Getenv("AUTHORIZATION"),
        }),
    )
    request := operations.V2CreateLedgerRequest{
        V2CreateLedgerRequest: &shared.V2CreateLedgerRequest{
            Metadata: map[string]string{
                "admin": "true",
            },
        },
        Ledger: "ledger001",
    }
    ctx := context.Background()
    res, err := s.Ledger.V2.CreateLedger(ctx, request)
    if err != nil {
        log.Fatal(err)
    }
    if res != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                                | Type                                                                                     | Required                                                                                 | Description                                                                              |
| ---------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------- |
| `ctx`                                                                                    | [context.Context](https://pkg.go.dev/context#Context)                                    | :heavy_check_mark:                                                                       | The context to use for the request.                                                      |
| `request`                                                                                | [operations.V2CreateLedgerRequest](../../pkg/models/operations/v2createledgerrequest.md) | :heavy_check_mark:                                                                       | The request object to use for the request.                                               |
| `opts`                                                                                   | [][operations.Option](../../pkg/models/operations/option.md)                             | :heavy_minus_sign:                                                                       | The options for this request.                                                            |


### Response

**[*operations.V2CreateLedgerResponse](../../pkg/models/operations/v2createledgerresponse.md), error**
| Error Object              | Status Code               | Content Type              |
| ------------------------- | ------------------------- | ------------------------- |
| sdkerrors.V2ErrorResponse | default                   | application/json          |
| sdkerrors.SDKError        | 4xx-5xx                   | */*                       |

## CreateTransaction

Create a new transaction to a ledger

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"os"
	"github.com/formancehq/formance-sdk-go/v2"
	"math/big"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"context"
	"log"
)

func main() {
    s := v2.New(
        v2.WithSecurity(shared.Security{
            Authorization: os.Getenv("AUTHORIZATION"),
        }),
    )
    request := operations.V2CreateTransactionRequest{
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
                Vars: map[string]any{
                    "user": "users:042",
                },
            },
        },
        DryRun: v2.Bool(true),
        Ledger: "ledger001",
    }
    ctx := context.Background()
    res, err := s.Ledger.V2.CreateTransaction(ctx, request)
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
| `opts`                                                                                             | [][operations.Option](../../pkg/models/operations/option.md)                                       | :heavy_minus_sign:                                                                                 | The options for this request.                                                                      |


### Response

**[*operations.V2CreateTransactionResponse](../../pkg/models/operations/v2createtransactionresponse.md), error**
| Error Object              | Status Code               | Content Type              |
| ------------------------- | ------------------------- | ------------------------- |
| sdkerrors.V2ErrorResponse | default                   | application/json          |
| sdkerrors.SDKError        | 4xx-5xx                   | */*                       |

## DeleteAccountMetadata

Delete metadata by key

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"os"
	"github.com/formancehq/formance-sdk-go/v2"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"context"
	"log"
)

func main() {
    s := v2.New(
        v2.WithSecurity(shared.Security{
            Authorization: os.Getenv("AUTHORIZATION"),
        }),
    )
    request := operations.V2DeleteAccountMetadataRequest{
        Address: "69266 Krajcik Bypass",
        Key: "foo",
        Ledger: "ledger001",
    }
    ctx := context.Background()
    res, err := s.Ledger.V2.DeleteAccountMetadata(ctx, request)
    if err != nil {
        log.Fatal(err)
    }
    if res != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                                                  | Type                                                                                                       | Required                                                                                                   | Description                                                                                                |
| ---------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------- |
| `ctx`                                                                                                      | [context.Context](https://pkg.go.dev/context#Context)                                                      | :heavy_check_mark:                                                                                         | The context to use for the request.                                                                        |
| `request`                                                                                                  | [operations.V2DeleteAccountMetadataRequest](../../pkg/models/operations/v2deleteaccountmetadatarequest.md) | :heavy_check_mark:                                                                                         | The request object to use for the request.                                                                 |
| `opts`                                                                                                     | [][operations.Option](../../pkg/models/operations/option.md)                                               | :heavy_minus_sign:                                                                                         | The options for this request.                                                                              |


### Response

**[*operations.V2DeleteAccountMetadataResponse](../../pkg/models/operations/v2deleteaccountmetadataresponse.md), error**
| Error Object       | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4xx-5xx            | */*                |

## DeleteLedgerMetadata

Delete ledger metadata by key

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"os"
	"github.com/formancehq/formance-sdk-go/v2"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"context"
	"log"
)

func main() {
    s := v2.New(
        v2.WithSecurity(shared.Security{
            Authorization: os.Getenv("AUTHORIZATION"),
        }),
    )
    request := operations.V2DeleteLedgerMetadataRequest{
        Key: "foo",
        Ledger: "ledger001",
    }
    ctx := context.Background()
    res, err := s.Ledger.V2.DeleteLedgerMetadata(ctx, request)
    if err != nil {
        log.Fatal(err)
    }
    if res != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                                                | Type                                                                                                     | Required                                                                                                 | Description                                                                                              |
| -------------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------- |
| `ctx`                                                                                                    | [context.Context](https://pkg.go.dev/context#Context)                                                    | :heavy_check_mark:                                                                                       | The context to use for the request.                                                                      |
| `request`                                                                                                | [operations.V2DeleteLedgerMetadataRequest](../../pkg/models/operations/v2deleteledgermetadatarequest.md) | :heavy_check_mark:                                                                                       | The request object to use for the request.                                                               |
| `opts`                                                                                                   | [][operations.Option](../../pkg/models/operations/option.md)                                             | :heavy_minus_sign:                                                                                       | The options for this request.                                                                            |


### Response

**[*operations.V2DeleteLedgerMetadataResponse](../../pkg/models/operations/v2deleteledgermetadataresponse.md), error**
| Error Object              | Status Code               | Content Type              |
| ------------------------- | ------------------------- | ------------------------- |
| sdkerrors.V2ErrorResponse | default                   | application/json          |
| sdkerrors.SDKError        | 4xx-5xx                   | */*                       |

## DeleteTransactionMetadata

Delete metadata by key

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"os"
	"github.com/formancehq/formance-sdk-go/v2"
	"math/big"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"context"
	"log"
)

func main() {
    s := v2.New(
        v2.WithSecurity(shared.Security{
            Authorization: os.Getenv("AUTHORIZATION"),
        }),
    )
    request := operations.V2DeleteTransactionMetadataRequest{
        ID: big.NewInt(1234),
        Key: "foo",
        Ledger: "ledger001",
    }
    ctx := context.Background()
    res, err := s.Ledger.V2.DeleteTransactionMetadata(ctx, request)
    if err != nil {
        log.Fatal(err)
    }
    if res != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                                                          | Type                                                                                                               | Required                                                                                                           | Description                                                                                                        |
| ------------------------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------------------------ |
| `ctx`                                                                                                              | [context.Context](https://pkg.go.dev/context#Context)                                                              | :heavy_check_mark:                                                                                                 | The context to use for the request.                                                                                |
| `request`                                                                                                          | [operations.V2DeleteTransactionMetadataRequest](../../pkg/models/operations/v2deletetransactionmetadatarequest.md) | :heavy_check_mark:                                                                                                 | The request object to use for the request.                                                                         |
| `opts`                                                                                                             | [][operations.Option](../../pkg/models/operations/option.md)                                                       | :heavy_minus_sign:                                                                                                 | The options for this request.                                                                                      |


### Response

**[*operations.V2DeleteTransactionMetadataResponse](../../pkg/models/operations/v2deletetransactionmetadataresponse.md), error**
| Error Object              | Status Code               | Content Type              |
| ------------------------- | ------------------------- | ------------------------- |
| sdkerrors.V2ErrorResponse | default                   | application/json          |
| sdkerrors.SDKError        | 4xx-5xx                   | */*                       |

## ExportLogs

Export logs

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"os"
	"github.com/formancehq/formance-sdk-go/v2"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"context"
	"log"
)

func main() {
    s := v2.New(
        v2.WithSecurity(shared.Security{
            Authorization: os.Getenv("AUTHORIZATION"),
        }),
    )
    request := operations.V2ExportLogsRequest{
        Ledger: "ledger001",
    }
    ctx := context.Background()
    res, err := s.Ledger.V2.ExportLogs(ctx, request)
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
| `request`                                                                            | [operations.V2ExportLogsRequest](../../pkg/models/operations/v2exportlogsrequest.md) | :heavy_check_mark:                                                                   | The request object to use for the request.                                           |
| `opts`                                                                               | [][operations.Option](../../pkg/models/operations/option.md)                         | :heavy_minus_sign:                                                                   | The options for this request.                                                        |


### Response

**[*operations.V2ExportLogsResponse](../../pkg/models/operations/v2exportlogsresponse.md), error**
| Error Object       | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4xx-5xx            | */*                |

## GetAccount

Get account by its address

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"os"
	"github.com/formancehq/formance-sdk-go/v2"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"context"
	"log"
)

func main() {
    s := v2.New(
        v2.WithSecurity(shared.Security{
            Authorization: os.Getenv("AUTHORIZATION"),
        }),
    )
    request := operations.V2GetAccountRequest{
        Address: "users:001",
        Ledger: "ledger001",
    }
    ctx := context.Background()
    res, err := s.Ledger.V2.GetAccount(ctx, request)
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
| `opts`                                                                               | [][operations.Option](../../pkg/models/operations/option.md)                         | :heavy_minus_sign:                                                                   | The options for this request.                                                        |


### Response

**[*operations.V2GetAccountResponse](../../pkg/models/operations/v2getaccountresponse.md), error**
| Error Object              | Status Code               | Content Type              |
| ------------------------- | ------------------------- | ------------------------- |
| sdkerrors.V2ErrorResponse | default                   | application/json          |
| sdkerrors.SDKError        | 4xx-5xx                   | */*                       |

## GetBalancesAggregated

Get the aggregated balances from selected accounts

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"os"
	"github.com/formancehq/formance-sdk-go/v2"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"context"
	"log"
)

func main() {
    s := v2.New(
        v2.WithSecurity(shared.Security{
            Authorization: os.Getenv("AUTHORIZATION"),
        }),
    )
    request := operations.V2GetBalancesAggregatedRequest{
        Ledger: "ledger001",
    }
    ctx := context.Background()
    res, err := s.Ledger.V2.GetBalancesAggregated(ctx, request)
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
| `opts`                                                                                                     | [][operations.Option](../../pkg/models/operations/option.md)                                               | :heavy_minus_sign:                                                                                         | The options for this request.                                                                              |


### Response

**[*operations.V2GetBalancesAggregatedResponse](../../pkg/models/operations/v2getbalancesaggregatedresponse.md), error**
| Error Object              | Status Code               | Content Type              |
| ------------------------- | ------------------------- | ------------------------- |
| sdkerrors.V2ErrorResponse | default                   | application/json          |
| sdkerrors.SDKError        | 4xx-5xx                   | */*                       |

## GetInfo

Show server information

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"os"
	"github.com/formancehq/formance-sdk-go/v2"
	"context"
	"log"
)

func main() {
    s := v2.New(
        v2.WithSecurity(shared.Security{
            Authorization: os.Getenv("AUTHORIZATION"),
        }),
    )

    ctx := context.Background()
    res, err := s.Ledger.V2.GetInfo(ctx)
    if err != nil {
        log.Fatal(err)
    }
    if res.V2ConfigInfoResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                    | Type                                                         | Required                                                     | Description                                                  |
| ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ |
| `ctx`                                                        | [context.Context](https://pkg.go.dev/context#Context)        | :heavy_check_mark:                                           | The context to use for the request.                          |
| `opts`                                                       | [][operations.Option](../../pkg/models/operations/option.md) | :heavy_minus_sign:                                           | The options for this request.                                |


### Response

**[*operations.V2GetInfoResponse](../../pkg/models/operations/v2getinforesponse.md), error**
| Error Object              | Status Code               | Content Type              |
| ------------------------- | ------------------------- | ------------------------- |
| sdkerrors.V2ErrorResponse | default                   | application/json          |
| sdkerrors.SDKError        | 4xx-5xx                   | */*                       |

## GetLedger

Get a ledger

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"os"
	"github.com/formancehq/formance-sdk-go/v2"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"context"
	"log"
)

func main() {
    s := v2.New(
        v2.WithSecurity(shared.Security{
            Authorization: os.Getenv("AUTHORIZATION"),
        }),
    )
    request := operations.V2GetLedgerRequest{
        Ledger: "ledger001",
    }
    ctx := context.Background()
    res, err := s.Ledger.V2.GetLedger(ctx, request)
    if err != nil {
        log.Fatal(err)
    }
    if res.V2GetLedgerResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                          | Type                                                                               | Required                                                                           | Description                                                                        |
| ---------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------- |
| `ctx`                                                                              | [context.Context](https://pkg.go.dev/context#Context)                              | :heavy_check_mark:                                                                 | The context to use for the request.                                                |
| `request`                                                                          | [operations.V2GetLedgerRequest](../../pkg/models/operations/v2getledgerrequest.md) | :heavy_check_mark:                                                                 | The request object to use for the request.                                         |
| `opts`                                                                             | [][operations.Option](../../pkg/models/operations/option.md)                       | :heavy_minus_sign:                                                                 | The options for this request.                                                      |


### Response

**[*operations.V2GetLedgerResponse](../../pkg/models/operations/v2getledgerresponse.md), error**
| Error Object              | Status Code               | Content Type              |
| ------------------------- | ------------------------- | ------------------------- |
| sdkerrors.V2ErrorResponse | default                   | application/json          |
| sdkerrors.SDKError        | 4xx-5xx                   | */*                       |

## GetLedgerInfo

Get information about a ledger

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"os"
	"github.com/formancehq/formance-sdk-go/v2"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"context"
	"log"
)

func main() {
    s := v2.New(
        v2.WithSecurity(shared.Security{
            Authorization: os.Getenv("AUTHORIZATION"),
        }),
    )
    request := operations.V2GetLedgerInfoRequest{
        Ledger: "ledger001",
    }
    ctx := context.Background()
    res, err := s.Ledger.V2.GetLedgerInfo(ctx, request)
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
| `opts`                                                                                     | [][operations.Option](../../pkg/models/operations/option.md)                               | :heavy_minus_sign:                                                                         | The options for this request.                                                              |


### Response

**[*operations.V2GetLedgerInfoResponse](../../pkg/models/operations/v2getledgerinforesponse.md), error**
| Error Object              | Status Code               | Content Type              |
| ------------------------- | ------------------------- | ------------------------- |
| sdkerrors.V2ErrorResponse | default                   | application/json          |
| sdkerrors.SDKError        | 4xx-5xx                   | */*                       |

## GetTransaction

Get transaction from a ledger by its ID

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"os"
	"github.com/formancehq/formance-sdk-go/v2"
	"math/big"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"context"
	"log"
)

func main() {
    s := v2.New(
        v2.WithSecurity(shared.Security{
            Authorization: os.Getenv("AUTHORIZATION"),
        }),
    )
    request := operations.V2GetTransactionRequest{
        ID: big.NewInt(1234),
        Ledger: "ledger001",
    }
    ctx := context.Background()
    res, err := s.Ledger.V2.GetTransaction(ctx, request)
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
| `opts`                                                                                       | [][operations.Option](../../pkg/models/operations/option.md)                                 | :heavy_minus_sign:                                                                           | The options for this request.                                                                |


### Response

**[*operations.V2GetTransactionResponse](../../pkg/models/operations/v2gettransactionresponse.md), error**
| Error Object              | Status Code               | Content Type              |
| ------------------------- | ------------------------- | ------------------------- |
| sdkerrors.V2ErrorResponse | default                   | application/json          |
| sdkerrors.SDKError        | 4xx-5xx                   | */*                       |

## GetVolumesWithBalances

Get list of volumes with balances for (account/asset)

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"os"
	"github.com/formancehq/formance-sdk-go/v2"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"context"
	"log"
)

func main() {
    s := v2.New(
        v2.WithSecurity(shared.Security{
            Authorization: os.Getenv("AUTHORIZATION"),
        }),
    )
    request := operations.V2GetVolumesWithBalancesRequest{
        Cursor: v2.String("aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ=="),
        GroupBy: v2.Int64(3),
        Ledger: "ledger001",
        PageSize: v2.Int64(100),
    }
    ctx := context.Background()
    res, err := s.Ledger.V2.GetVolumesWithBalances(ctx, request)
    if err != nil {
        log.Fatal(err)
    }
    if res.V2VolumesWithBalanceCursorResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                                                    | Type                                                                                                         | Required                                                                                                     | Description                                                                                                  |
| ------------------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------------------ |
| `ctx`                                                                                                        | [context.Context](https://pkg.go.dev/context#Context)                                                        | :heavy_check_mark:                                                                                           | The context to use for the request.                                                                          |
| `request`                                                                                                    | [operations.V2GetVolumesWithBalancesRequest](../../pkg/models/operations/v2getvolumeswithbalancesrequest.md) | :heavy_check_mark:                                                                                           | The request object to use for the request.                                                                   |
| `opts`                                                                                                       | [][operations.Option](../../pkg/models/operations/option.md)                                                 | :heavy_minus_sign:                                                                                           | The options for this request.                                                                                |


### Response

**[*operations.V2GetVolumesWithBalancesResponse](../../pkg/models/operations/v2getvolumeswithbalancesresponse.md), error**
| Error Object              | Status Code               | Content Type              |
| ------------------------- | ------------------------- | ------------------------- |
| sdkerrors.V2ErrorResponse | default                   | application/json          |
| sdkerrors.SDKError        | 4xx-5xx                   | */*                       |

## ImportLogs

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"os"
	"github.com/formancehq/formance-sdk-go/v2"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"context"
	"log"
)

func main() {
    s := v2.New(
        v2.WithSecurity(shared.Security{
            Authorization: os.Getenv("AUTHORIZATION"),
        }),
    )
    request := operations.V2ImportLogsRequest{
        Ledger: "ledger001",
    }
    ctx := context.Background()
    res, err := s.Ledger.V2.ImportLogs(ctx, request)
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
| `request`                                                                            | [operations.V2ImportLogsRequest](../../pkg/models/operations/v2importlogsrequest.md) | :heavy_check_mark:                                                                   | The request object to use for the request.                                           |
| `opts`                                                                               | [][operations.Option](../../pkg/models/operations/option.md)                         | :heavy_minus_sign:                                                                   | The options for this request.                                                        |


### Response

**[*operations.V2ImportLogsResponse](../../pkg/models/operations/v2importlogsresponse.md), error**
| Error Object              | Status Code               | Content Type              |
| ------------------------- | ------------------------- | ------------------------- |
| sdkerrors.V2ErrorResponse | default                   | application/json          |
| sdkerrors.SDKError        | 4xx-5xx                   | */*                       |

## ListAccounts

List accounts from a ledger, sorted by address in descending order.

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"os"
	"github.com/formancehq/formance-sdk-go/v2"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"context"
	"log"
)

func main() {
    s := v2.New(
        v2.WithSecurity(shared.Security{
            Authorization: os.Getenv("AUTHORIZATION"),
        }),
    )
    request := operations.V2ListAccountsRequest{
        Cursor: v2.String("aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ=="),
        Ledger: "ledger001",
        PageSize: v2.Int64(100),
    }
    ctx := context.Background()
    res, err := s.Ledger.V2.ListAccounts(ctx, request)
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
| `opts`                                                                                   | [][operations.Option](../../pkg/models/operations/option.md)                             | :heavy_minus_sign:                                                                       | The options for this request.                                                            |


### Response

**[*operations.V2ListAccountsResponse](../../pkg/models/operations/v2listaccountsresponse.md), error**
| Error Object              | Status Code               | Content Type              |
| ------------------------- | ------------------------- | ------------------------- |
| sdkerrors.V2ErrorResponse | default                   | application/json          |
| sdkerrors.SDKError        | 4xx-5xx                   | */*                       |

## ListLedgers

List ledgers

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"os"
	"github.com/formancehq/formance-sdk-go/v2"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"context"
	"log"
)

func main() {
    s := v2.New(
        v2.WithSecurity(shared.Security{
            Authorization: os.Getenv("AUTHORIZATION"),
        }),
    )
    request := operations.V2ListLedgersRequest{
        Cursor: v2.String("aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ=="),
        PageSize: v2.Int64(100),
    }
    ctx := context.Background()
    res, err := s.Ledger.V2.ListLedgers(ctx, request)
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
| `opts`                                                                                 | [][operations.Option](../../pkg/models/operations/option.md)                           | :heavy_minus_sign:                                                                     | The options for this request.                                                          |


### Response

**[*operations.V2ListLedgersResponse](../../pkg/models/operations/v2listledgersresponse.md), error**
| Error Object              | Status Code               | Content Type              |
| ------------------------- | ------------------------- | ------------------------- |
| sdkerrors.V2ErrorResponse | default                   | application/json          |
| sdkerrors.SDKError        | 4xx-5xx                   | */*                       |

## ListLogs

List the logs from a ledger, sorted by ID in descending order.

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"os"
	"github.com/formancehq/formance-sdk-go/v2"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"context"
	"log"
)

func main() {
    s := v2.New(
        v2.WithSecurity(shared.Security{
            Authorization: os.Getenv("AUTHORIZATION"),
        }),
    )
    request := operations.V2ListLogsRequest{
        Cursor: v2.String("aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ=="),
        Ledger: "ledger001",
        PageSize: v2.Int64(100),
    }
    ctx := context.Background()
    res, err := s.Ledger.V2.ListLogs(ctx, request)
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
| `opts`                                                                           | [][operations.Option](../../pkg/models/operations/option.md)                     | :heavy_minus_sign:                                                               | The options for this request.                                                    |


### Response

**[*operations.V2ListLogsResponse](../../pkg/models/operations/v2listlogsresponse.md), error**
| Error Object              | Status Code               | Content Type              |
| ------------------------- | ------------------------- | ------------------------- |
| sdkerrors.V2ErrorResponse | default                   | application/json          |
| sdkerrors.SDKError        | 4xx-5xx                   | */*                       |

## ListTransactions

List transactions from a ledger, sorted by id in descending order.

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"os"
	"github.com/formancehq/formance-sdk-go/v2"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"context"
	"log"
)

func main() {
    s := v2.New(
        v2.WithSecurity(shared.Security{
            Authorization: os.Getenv("AUTHORIZATION"),
        }),
    )
    request := operations.V2ListTransactionsRequest{
        Cursor: v2.String("aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ=="),
        Ledger: "ledger001",
        PageSize: v2.Int64(100),
    }
    ctx := context.Background()
    res, err := s.Ledger.V2.ListTransactions(ctx, request)
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
| `opts`                                                                                           | [][operations.Option](../../pkg/models/operations/option.md)                                     | :heavy_minus_sign:                                                                               | The options for this request.                                                                    |


### Response

**[*operations.V2ListTransactionsResponse](../../pkg/models/operations/v2listtransactionsresponse.md), error**
| Error Object              | Status Code               | Content Type              |
| ------------------------- | ------------------------- | ------------------------- |
| sdkerrors.V2ErrorResponse | default                   | application/json          |
| sdkerrors.SDKError        | 4xx-5xx                   | */*                       |

## ReadStats

Get statistics from a ledger. (aggregate metrics on accounts and transactions)


### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"os"
	"github.com/formancehq/formance-sdk-go/v2"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"context"
	"log"
)

func main() {
    s := v2.New(
        v2.WithSecurity(shared.Security{
            Authorization: os.Getenv("AUTHORIZATION"),
        }),
    )
    request := operations.V2ReadStatsRequest{
        Ledger: "ledger001",
    }
    ctx := context.Background()
    res, err := s.Ledger.V2.ReadStats(ctx, request)
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
| `opts`                                                                             | [][operations.Option](../../pkg/models/operations/option.md)                       | :heavy_minus_sign:                                                                 | The options for this request.                                                      |


### Response

**[*operations.V2ReadStatsResponse](../../pkg/models/operations/v2readstatsresponse.md), error**
| Error Object              | Status Code               | Content Type              |
| ------------------------- | ------------------------- | ------------------------- |
| sdkerrors.V2ErrorResponse | default                   | application/json          |
| sdkerrors.SDKError        | 4xx-5xx                   | */*                       |

## RevertTransaction

Revert a ledger transaction by its ID

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"os"
	"github.com/formancehq/formance-sdk-go/v2"
	"math/big"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"context"
	"log"
)

func main() {
    s := v2.New(
        v2.WithSecurity(shared.Security{
            Authorization: os.Getenv("AUTHORIZATION"),
        }),
    )
    request := operations.V2RevertTransactionRequest{
        ID: big.NewInt(1234),
        Ledger: "ledger001",
    }
    ctx := context.Background()
    res, err := s.Ledger.V2.RevertTransaction(ctx, request)
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
| `opts`                                                                                             | [][operations.Option](../../pkg/models/operations/option.md)                                       | :heavy_minus_sign:                                                                                 | The options for this request.                                                                      |


### Response

**[*operations.V2RevertTransactionResponse](../../pkg/models/operations/v2reverttransactionresponse.md), error**
| Error Object              | Status Code               | Content Type              |
| ------------------------- | ------------------------- | ------------------------- |
| sdkerrors.V2ErrorResponse | default                   | application/json          |
| sdkerrors.SDKError        | 4xx-5xx                   | */*                       |

## UpdateLedgerMetadata

Update ledger metadata

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"os"
	"github.com/formancehq/formance-sdk-go/v2"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"context"
	"log"
)

func main() {
    s := v2.New(
        v2.WithSecurity(shared.Security{
            Authorization: os.Getenv("AUTHORIZATION"),
        }),
    )
    request := operations.V2UpdateLedgerMetadataRequest{
        RequestBody: map[string]string{
            "admin": "true",
        },
        Ledger: "ledger001",
    }
    ctx := context.Background()
    res, err := s.Ledger.V2.UpdateLedgerMetadata(ctx, request)
    if err != nil {
        log.Fatal(err)
    }
    if res != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                                                | Type                                                                                                     | Required                                                                                                 | Description                                                                                              |
| -------------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------- |
| `ctx`                                                                                                    | [context.Context](https://pkg.go.dev/context#Context)                                                    | :heavy_check_mark:                                                                                       | The context to use for the request.                                                                      |
| `request`                                                                                                | [operations.V2UpdateLedgerMetadataRequest](../../pkg/models/operations/v2updateledgermetadatarequest.md) | :heavy_check_mark:                                                                                       | The request object to use for the request.                                                               |
| `opts`                                                                                                   | [][operations.Option](../../pkg/models/operations/option.md)                                             | :heavy_minus_sign:                                                                                       | The options for this request.                                                                            |


### Response

**[*operations.V2UpdateLedgerMetadataResponse](../../pkg/models/operations/v2updateledgermetadataresponse.md), error**
| Error Object              | Status Code               | Content Type              |
| ------------------------- | ------------------------- | ------------------------- |
| sdkerrors.V2ErrorResponse | default                   | application/json          |
| sdkerrors.SDKError        | 4xx-5xx                   | */*                       |
