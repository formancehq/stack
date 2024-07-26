# FormancePaymentsV1
(*Payments.V1*)

### Available Operations

* [AddAccountToPool](#addaccounttopool) - Add an account to a pool
* [ConnectorsTransfer](#connectorstransfer) - Transfer funds between Connector accounts
* [CreateAccount](#createaccount) - Create an account
* [CreateBankAccount](#createbankaccount) - Create a BankAccount in Payments and on the PSP
* [CreatePayment](#createpayment) - Create a payment
* [CreatePool](#createpool) - Create a Pool
* [CreateTransferInitiation](#createtransferinitiation) - Create a TransferInitiation
* [DeletePool](#deletepool) - Delete a Pool
* [DeleteTransferInitiation](#deletetransferinitiation) - Delete a transfer initiation
* [ForwardBankAccount](#forwardbankaccount) - Forward a bank account to a connector
* [GetAccountBalances](#getaccountbalances) - Get account balances
* [GetBankAccount](#getbankaccount) - Get a bank account created by user on Formance
* [~~GetConnectorTask~~](#getconnectortask) - Read a specific task of the connector :warning: **Deprecated**
* [GetConnectorTaskV1](#getconnectortaskv1) - Read a specific task of the connector
* [GetPayment](#getpayment) - Get a payment
* [GetPool](#getpool) - Get a Pool
* [GetPoolBalances](#getpoolbalances) - Get pool balances
* [GetTransferInitiation](#gettransferinitiation) - Get a transfer initiation
* [InstallConnector](#installconnector) - Install a connector
* [ListAllConnectors](#listallconnectors) - List all installed connectors
* [ListBankAccounts](#listbankaccounts) - List bank accounts created by user on Formance
* [ListConfigsAvailableConnectors](#listconfigsavailableconnectors) - List the configs of each available connector
* [~~ListConnectorTasks~~](#listconnectortasks) - List tasks from a connector :warning: **Deprecated**
* [ListConnectorTasksV1](#listconnectortasksv1) - List tasks from a connector
* [ListPayments](#listpayments) - List payments
* [ListPools](#listpools) - List Pools
* [ListTransferInitiations](#listtransferinitiations) - List Transfer Initiations
* [PaymentsgetAccount](#paymentsgetaccount) - Get an account
* [PaymentsgetServerInfo](#paymentsgetserverinfo) - Get server info
* [PaymentslistAccounts](#paymentslistaccounts) - List accounts
* [~~ReadConnectorConfig~~](#readconnectorconfig) - Read the config of a connector :warning: **Deprecated**
* [ReadConnectorConfigV1](#readconnectorconfigv1) - Read the config of a connector
* [RemoveAccountFromPool](#removeaccountfrompool) - Remove an account from a pool
* [~~ResetConnector~~](#resetconnector) - Reset a connector :warning: **Deprecated**
* [ResetConnectorV1](#resetconnectorv1) - Reset a connector
* [RetryTransferInitiation](#retrytransferinitiation) - Retry a failed transfer initiation
* [ReverseTransferInitiation](#reversetransferinitiation) - Reverse a transfer initiation
* [UdpateTransferInitiationStatus](#udpatetransferinitiationstatus) - Update the status of a transfer initiation
* [~~UninstallConnector~~](#uninstallconnector) - Uninstall a connector :warning: **Deprecated**
* [UninstallConnectorV1](#uninstallconnectorv1) - Uninstall a connector
* [UpdateBankAccountMetadata](#updatebankaccountmetadata) - Update metadata of a bank account
* [UpdateConnectorConfigV1](#updateconnectorconfigv1) - Update the config of a connector
* [UpdateMetadata](#updatemetadata) - Update metadata

## AddAccountToPool

Add an account to a pool

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
    request := operations.AddAccountToPoolRequest{
        AddAccountToPoolRequest: shared.AddAccountToPoolRequest{
            AccountID: "<value>",
        },
        PoolID: "XXX",
    }
    ctx := context.Background()
    res, err := s.Payments.V1.AddAccountToPool(ctx, request)
    if err != nil {
        log.Fatal(err)
    }
    if res != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                                    | Type                                                                                         | Required                                                                                     | Description                                                                                  |
| -------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------- |
| `ctx`                                                                                        | [context.Context](https://pkg.go.dev/context#Context)                                        | :heavy_check_mark:                                                                           | The context to use for the request.                                                          |
| `request`                                                                                    | [operations.AddAccountToPoolRequest](../../pkg/models/operations/addaccounttopoolrequest.md) | :heavy_check_mark:                                                                           | The request object to use for the request.                                                   |
| `opts`                                                                                       | [][operations.Option](../../pkg/models/operations/option.md)                                 | :heavy_minus_sign:                                                                           | The options for this request.                                                                |


### Response

**[*operations.AddAccountToPoolResponse](../../pkg/models/operations/addaccounttopoolresponse.md), error**
| Error Object                    | Status Code                     | Content Type                    |
| ------------------------------- | ------------------------------- | ------------------------------- |
| sdkerrors.PaymentsErrorResponse | default                         | application/json                |
| sdkerrors.SDKError              | 4xx-5xx                         | */*                             |

## ConnectorsTransfer

Execute a transfer between two accounts.

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
    request := operations.ConnectorsTransferRequest{
        TransferRequest: shared.TransferRequest{
            Amount: big.NewInt(100),
            Asset: "USD",
            Destination: "acct_1Gqj58KZcSIg2N2q",
            Source: v2.String("acct_1Gqj58KZcSIg2N2q"),
        },
        Connector: shared.ConnectorBankingCircle,
    }
    ctx := context.Background()
    res, err := s.Payments.V1.ConnectorsTransfer(ctx, request)
    if err != nil {
        log.Fatal(err)
    }
    if res.TransferResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                                        | Type                                                                                             | Required                                                                                         | Description                                                                                      |
| ------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------ |
| `ctx`                                                                                            | [context.Context](https://pkg.go.dev/context#Context)                                            | :heavy_check_mark:                                                                               | The context to use for the request.                                                              |
| `request`                                                                                        | [operations.ConnectorsTransferRequest](../../pkg/models/operations/connectorstransferrequest.md) | :heavy_check_mark:                                                                               | The request object to use for the request.                                                       |
| `opts`                                                                                           | [][operations.Option](../../pkg/models/operations/option.md)                                     | :heavy_minus_sign:                                                                               | The options for this request.                                                                    |


### Response

**[*operations.ConnectorsTransferResponse](../../pkg/models/operations/connectorstransferresponse.md), error**
| Error Object                    | Status Code                     | Content Type                    |
| ------------------------------- | ------------------------------- | ------------------------------- |
| sdkerrors.PaymentsErrorResponse | default                         | application/json                |
| sdkerrors.SDKError              | 4xx-5xx                         | */*                             |

## CreateAccount

Create an account

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"os"
	"github.com/formancehq/formance-sdk-go/v2"
	"github.com/formancehq/formance-sdk-go/v2/pkg/types"
	"context"
	"log"
)

func main() {
    s := v2.New(
        v2.WithSecurity(shared.Security{
            Authorization: os.Getenv("AUTHORIZATION"),
        }),
    )
    request := shared.AccountRequest{
        ConnectorID: "<value>",
        CreatedAt: types.MustTimeFromString("2024-08-19T02:15:08.668Z"),
        Reference: "<value>",
        Type: shared.AccountTypeUnknown,
    }
    ctx := context.Background()
    res, err := s.Payments.V1.CreateAccount(ctx, request)
    if err != nil {
        log.Fatal(err)
    }
    if res.PaymentsAccountResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                          | Type                                                               | Required                                                           | Description                                                        |
| ------------------------------------------------------------------ | ------------------------------------------------------------------ | ------------------------------------------------------------------ | ------------------------------------------------------------------ |
| `ctx`                                                              | [context.Context](https://pkg.go.dev/context#Context)              | :heavy_check_mark:                                                 | The context to use for the request.                                |
| `request`                                                          | [shared.AccountRequest](../../pkg/models/shared/accountrequest.md) | :heavy_check_mark:                                                 | The request object to use for the request.                         |
| `opts`                                                             | [][operations.Option](../../pkg/models/operations/option.md)       | :heavy_minus_sign:                                                 | The options for this request.                                      |


### Response

**[*operations.CreateAccountResponse](../../pkg/models/operations/createaccountresponse.md), error**
| Error Object                    | Status Code                     | Content Type                    |
| ------------------------------- | ------------------------------- | ------------------------------- |
| sdkerrors.PaymentsErrorResponse | default                         | application/json                |
| sdkerrors.SDKError              | 4xx-5xx                         | */*                             |

## CreateBankAccount

Create a bank account in Payments and on the PSP.

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
    request := shared.BankAccountRequest{
        ConnectorID: "<value>",
        Country: "GB",
        Name: "My account",
    }
    ctx := context.Background()
    res, err := s.Payments.V1.CreateBankAccount(ctx, request)
    if err != nil {
        log.Fatal(err)
    }
    if res.BankAccountResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                  | Type                                                                       | Required                                                                   | Description                                                                |
| -------------------------------------------------------------------------- | -------------------------------------------------------------------------- | -------------------------------------------------------------------------- | -------------------------------------------------------------------------- |
| `ctx`                                                                      | [context.Context](https://pkg.go.dev/context#Context)                      | :heavy_check_mark:                                                         | The context to use for the request.                                        |
| `request`                                                                  | [shared.BankAccountRequest](../../pkg/models/shared/bankaccountrequest.md) | :heavy_check_mark:                                                         | The request object to use for the request.                                 |
| `opts`                                                                     | [][operations.Option](../../pkg/models/operations/option.md)               | :heavy_minus_sign:                                                         | The options for this request.                                              |


### Response

**[*operations.CreateBankAccountResponse](../../pkg/models/operations/createbankaccountresponse.md), error**
| Error Object                    | Status Code                     | Content Type                    |
| ------------------------------- | ------------------------------- | ------------------------------- |
| sdkerrors.PaymentsErrorResponse | default                         | application/json                |
| sdkerrors.SDKError              | 4xx-5xx                         | */*                             |

## CreatePayment

Create a payment

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"os"
	"github.com/formancehq/formance-sdk-go/v2"
	"math/big"
	"github.com/formancehq/formance-sdk-go/v2/pkg/types"
	"context"
	"log"
)

func main() {
    s := v2.New(
        v2.WithSecurity(shared.Security{
            Authorization: os.Getenv("AUTHORIZATION"),
        }),
    )
    request := shared.PaymentRequest{
        Amount: big.NewInt(100),
        Asset: "USD",
        ConnectorID: "<value>",
        CreatedAt: types.MustTimeFromString("2024-11-09T01:03:21.153Z"),
        Reference: "<value>",
        Scheme: shared.PaymentSchemeGooglePay,
        Status: shared.PaymentStatusDisputeWon,
        Type: shared.PaymentTypeTransfer,
    }
    ctx := context.Background()
    res, err := s.Payments.V1.CreatePayment(ctx, request)
    if err != nil {
        log.Fatal(err)
    }
    if res.PaymentResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                          | Type                                                               | Required                                                           | Description                                                        |
| ------------------------------------------------------------------ | ------------------------------------------------------------------ | ------------------------------------------------------------------ | ------------------------------------------------------------------ |
| `ctx`                                                              | [context.Context](https://pkg.go.dev/context#Context)              | :heavy_check_mark:                                                 | The context to use for the request.                                |
| `request`                                                          | [shared.PaymentRequest](../../pkg/models/shared/paymentrequest.md) | :heavy_check_mark:                                                 | The request object to use for the request.                         |
| `opts`                                                             | [][operations.Option](../../pkg/models/operations/option.md)       | :heavy_minus_sign:                                                 | The options for this request.                                      |


### Response

**[*operations.CreatePaymentResponse](../../pkg/models/operations/createpaymentresponse.md), error**
| Error Object                    | Status Code                     | Content Type                    |
| ------------------------------- | ------------------------------- | ------------------------------- |
| sdkerrors.PaymentsErrorResponse | default                         | application/json                |
| sdkerrors.SDKError              | 4xx-5xx                         | */*                             |

## CreatePool

Create a Pool

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
    request := shared.PoolRequest{
        AccountIDs: []string{
            "<value>",
        },
        Name: "<value>",
    }
    ctx := context.Background()
    res, err := s.Payments.V1.CreatePool(ctx, request)
    if err != nil {
        log.Fatal(err)
    }
    if res.PoolResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                    | Type                                                         | Required                                                     | Description                                                  |
| ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ |
| `ctx`                                                        | [context.Context](https://pkg.go.dev/context#Context)        | :heavy_check_mark:                                           | The context to use for the request.                          |
| `request`                                                    | [shared.PoolRequest](../../pkg/models/shared/poolrequest.md) | :heavy_check_mark:                                           | The request object to use for the request.                   |
| `opts`                                                       | [][operations.Option](../../pkg/models/operations/option.md) | :heavy_minus_sign:                                           | The options for this request.                                |


### Response

**[*operations.CreatePoolResponse](../../pkg/models/operations/createpoolresponse.md), error**
| Error Object                    | Status Code                     | Content Type                    |
| ------------------------------- | ------------------------------- | ------------------------------- |
| sdkerrors.PaymentsErrorResponse | default                         | application/json                |
| sdkerrors.SDKError              | 4xx-5xx                         | */*                             |

## CreateTransferInitiation

Create a transfer initiation

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"os"
	"github.com/formancehq/formance-sdk-go/v2"
	"math/big"
	"github.com/formancehq/formance-sdk-go/v2/pkg/types"
	"context"
	"log"
)

func main() {
    s := v2.New(
        v2.WithSecurity(shared.Security{
            Authorization: os.Getenv("AUTHORIZATION"),
        }),
    )
    request := shared.TransferInitiationRequest{
        Amount: big.NewInt(256698),
        Asset: "USD",
        Description: "Multi-tiered incremental methodology",
        DestinationAccountID: "<value>",
        Reference: "XXX",
        ScheduledAt: types.MustTimeFromString("2023-05-04T22:47:54.364Z"),
        SourceAccountID: "<value>",
        Type: shared.TransferInitiationRequestTypeTransfer,
        Validated: false,
    }
    ctx := context.Background()
    res, err := s.Payments.V1.CreateTransferInitiation(ctx, request)
    if err != nil {
        log.Fatal(err)
    }
    if res.TransferInitiationResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                                | Type                                                                                     | Required                                                                                 | Description                                                                              |
| ---------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------- |
| `ctx`                                                                                    | [context.Context](https://pkg.go.dev/context#Context)                                    | :heavy_check_mark:                                                                       | The context to use for the request.                                                      |
| `request`                                                                                | [shared.TransferInitiationRequest](../../pkg/models/shared/transferinitiationrequest.md) | :heavy_check_mark:                                                                       | The request object to use for the request.                                               |
| `opts`                                                                                   | [][operations.Option](../../pkg/models/operations/option.md)                             | :heavy_minus_sign:                                                                       | The options for this request.                                                            |


### Response

**[*operations.CreateTransferInitiationResponse](../../pkg/models/operations/createtransferinitiationresponse.md), error**
| Error Object                    | Status Code                     | Content Type                    |
| ------------------------------- | ------------------------------- | ------------------------------- |
| sdkerrors.PaymentsErrorResponse | default                         | application/json                |
| sdkerrors.SDKError              | 4xx-5xx                         | */*                             |

## DeletePool

Delete a pool by its id.

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
    request := operations.DeletePoolRequest{
        PoolID: "XXX",
    }
    ctx := context.Background()
    res, err := s.Payments.V1.DeletePool(ctx, request)
    if err != nil {
        log.Fatal(err)
    }
    if res != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                        | Type                                                                             | Required                                                                         | Description                                                                      |
| -------------------------------------------------------------------------------- | -------------------------------------------------------------------------------- | -------------------------------------------------------------------------------- | -------------------------------------------------------------------------------- |
| `ctx`                                                                            | [context.Context](https://pkg.go.dev/context#Context)                            | :heavy_check_mark:                                                               | The context to use for the request.                                              |
| `request`                                                                        | [operations.DeletePoolRequest](../../pkg/models/operations/deletepoolrequest.md) | :heavy_check_mark:                                                               | The request object to use for the request.                                       |
| `opts`                                                                           | [][operations.Option](../../pkg/models/operations/option.md)                     | :heavy_minus_sign:                                                               | The options for this request.                                                    |


### Response

**[*operations.DeletePoolResponse](../../pkg/models/operations/deletepoolresponse.md), error**
| Error Object                    | Status Code                     | Content Type                    |
| ------------------------------- | ------------------------------- | ------------------------------- |
| sdkerrors.PaymentsErrorResponse | default                         | application/json                |
| sdkerrors.SDKError              | 4xx-5xx                         | */*                             |

## DeleteTransferInitiation

Delete a transfer initiation by its id.

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
    request := operations.DeleteTransferInitiationRequest{
        TransferID: "XXX",
    }
    ctx := context.Background()
    res, err := s.Payments.V1.DeleteTransferInitiation(ctx, request)
    if err != nil {
        log.Fatal(err)
    }
    if res != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                                                    | Type                                                                                                         | Required                                                                                                     | Description                                                                                                  |
| ------------------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------------------ |
| `ctx`                                                                                                        | [context.Context](https://pkg.go.dev/context#Context)                                                        | :heavy_check_mark:                                                                                           | The context to use for the request.                                                                          |
| `request`                                                                                                    | [operations.DeleteTransferInitiationRequest](../../pkg/models/operations/deletetransferinitiationrequest.md) | :heavy_check_mark:                                                                                           | The request object to use for the request.                                                                   |
| `opts`                                                                                                       | [][operations.Option](../../pkg/models/operations/option.md)                                                 | :heavy_minus_sign:                                                                                           | The options for this request.                                                                                |


### Response

**[*operations.DeleteTransferInitiationResponse](../../pkg/models/operations/deletetransferinitiationresponse.md), error**
| Error Object                    | Status Code                     | Content Type                    |
| ------------------------------- | ------------------------------- | ------------------------------- |
| sdkerrors.PaymentsErrorResponse | default                         | application/json                |
| sdkerrors.SDKError              | 4xx-5xx                         | */*                             |

## ForwardBankAccount

Forward a bank account to a connector

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
    request := operations.ForwardBankAccountRequest{
        ForwardBankAccountRequest: shared.ForwardBankAccountRequest{
            ConnectorID: "<value>",
        },
        BankAccountID: "XXX",
    }
    ctx := context.Background()
    res, err := s.Payments.V1.ForwardBankAccount(ctx, request)
    if err != nil {
        log.Fatal(err)
    }
    if res.BankAccountResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                                        | Type                                                                                             | Required                                                                                         | Description                                                                                      |
| ------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------ |
| `ctx`                                                                                            | [context.Context](https://pkg.go.dev/context#Context)                                            | :heavy_check_mark:                                                                               | The context to use for the request.                                                              |
| `request`                                                                                        | [operations.ForwardBankAccountRequest](../../pkg/models/operations/forwardbankaccountrequest.md) | :heavy_check_mark:                                                                               | The request object to use for the request.                                                       |
| `opts`                                                                                           | [][operations.Option](../../pkg/models/operations/option.md)                                     | :heavy_minus_sign:                                                                               | The options for this request.                                                                    |


### Response

**[*operations.ForwardBankAccountResponse](../../pkg/models/operations/forwardbankaccountresponse.md), error**
| Error Object                    | Status Code                     | Content Type                    |
| ------------------------------- | ------------------------------- | ------------------------------- |
| sdkerrors.PaymentsErrorResponse | default                         | application/json                |
| sdkerrors.SDKError              | 4xx-5xx                         | */*                             |

## GetAccountBalances

Get account balances

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
    request := operations.GetAccountBalancesRequest{
        AccountID: "XXX",
        Cursor: v2.String("aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ=="),
        PageSize: v2.Int64(100),
        Sort: []string{
            "date:asc",
            "status:desc",
        },
    }
    ctx := context.Background()
    res, err := s.Payments.V1.GetAccountBalances(ctx, request)
    if err != nil {
        log.Fatal(err)
    }
    if res.BalancesCursor != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                                        | Type                                                                                             | Required                                                                                         | Description                                                                                      |
| ------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------ |
| `ctx`                                                                                            | [context.Context](https://pkg.go.dev/context#Context)                                            | :heavy_check_mark:                                                                               | The context to use for the request.                                                              |
| `request`                                                                                        | [operations.GetAccountBalancesRequest](../../pkg/models/operations/getaccountbalancesrequest.md) | :heavy_check_mark:                                                                               | The request object to use for the request.                                                       |
| `opts`                                                                                           | [][operations.Option](../../pkg/models/operations/option.md)                                     | :heavy_minus_sign:                                                                               | The options for this request.                                                                    |


### Response

**[*operations.GetAccountBalancesResponse](../../pkg/models/operations/getaccountbalancesresponse.md), error**
| Error Object                    | Status Code                     | Content Type                    |
| ------------------------------- | ------------------------------- | ------------------------------- |
| sdkerrors.PaymentsErrorResponse | default                         | application/json                |
| sdkerrors.SDKError              | 4xx-5xx                         | */*                             |

## GetBankAccount

Get a bank account created by user on Formance

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
    request := operations.GetBankAccountRequest{
        BankAccountID: "XXX",
    }
    ctx := context.Background()
    res, err := s.Payments.V1.GetBankAccount(ctx, request)
    if err != nil {
        log.Fatal(err)
    }
    if res.BankAccountResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                                | Type                                                                                     | Required                                                                                 | Description                                                                              |
| ---------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------- |
| `ctx`                                                                                    | [context.Context](https://pkg.go.dev/context#Context)                                    | :heavy_check_mark:                                                                       | The context to use for the request.                                                      |
| `request`                                                                                | [operations.GetBankAccountRequest](../../pkg/models/operations/getbankaccountrequest.md) | :heavy_check_mark:                                                                       | The request object to use for the request.                                               |
| `opts`                                                                                   | [][operations.Option](../../pkg/models/operations/option.md)                             | :heavy_minus_sign:                                                                       | The options for this request.                                                            |


### Response

**[*operations.GetBankAccountResponse](../../pkg/models/operations/getbankaccountresponse.md), error**
| Error Object                    | Status Code                     | Content Type                    |
| ------------------------------- | ------------------------------- | ------------------------------- |
| sdkerrors.PaymentsErrorResponse | default                         | application/json                |
| sdkerrors.SDKError              | 4xx-5xx                         | */*                             |

## ~~GetConnectorTask~~

Get a specific task associated to the connector.

> :warning: **DEPRECATED**: This will be removed in a future release, please migrate away from it as soon as possible.

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
    request := operations.GetConnectorTaskRequest{
        Connector: shared.ConnectorAdyen,
        TaskID: "task1",
    }
    ctx := context.Background()
    res, err := s.Payments.V1.GetConnectorTask(ctx, request)
    if err != nil {
        log.Fatal(err)
    }
    if res.TaskResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                                    | Type                                                                                         | Required                                                                                     | Description                                                                                  |
| -------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------- |
| `ctx`                                                                                        | [context.Context](https://pkg.go.dev/context#Context)                                        | :heavy_check_mark:                                                                           | The context to use for the request.                                                          |
| `request`                                                                                    | [operations.GetConnectorTaskRequest](../../pkg/models/operations/getconnectortaskrequest.md) | :heavy_check_mark:                                                                           | The request object to use for the request.                                                   |
| `opts`                                                                                       | [][operations.Option](../../pkg/models/operations/option.md)                                 | :heavy_minus_sign:                                                                           | The options for this request.                                                                |


### Response

**[*operations.GetConnectorTaskResponse](../../pkg/models/operations/getconnectortaskresponse.md), error**
| Error Object                    | Status Code                     | Content Type                    |
| ------------------------------- | ------------------------------- | ------------------------------- |
| sdkerrors.PaymentsErrorResponse | default                         | application/json                |
| sdkerrors.SDKError              | 4xx-5xx                         | */*                             |

## GetConnectorTaskV1

Get a specific task associated to the connector.

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
    request := operations.GetConnectorTaskV1Request{
        Connector: shared.ConnectorBankingCircle,
        ConnectorID: "XXX",
        TaskID: "task1",
    }
    ctx := context.Background()
    res, err := s.Payments.V1.GetConnectorTaskV1(ctx, request)
    if err != nil {
        log.Fatal(err)
    }
    if res.TaskResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                                        | Type                                                                                             | Required                                                                                         | Description                                                                                      |
| ------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------ |
| `ctx`                                                                                            | [context.Context](https://pkg.go.dev/context#Context)                                            | :heavy_check_mark:                                                                               | The context to use for the request.                                                              |
| `request`                                                                                        | [operations.GetConnectorTaskV1Request](../../pkg/models/operations/getconnectortaskv1request.md) | :heavy_check_mark:                                                                               | The request object to use for the request.                                                       |
| `opts`                                                                                           | [][operations.Option](../../pkg/models/operations/option.md)                                     | :heavy_minus_sign:                                                                               | The options for this request.                                                                    |


### Response

**[*operations.GetConnectorTaskV1Response](../../pkg/models/operations/getconnectortaskv1response.md), error**
| Error Object                    | Status Code                     | Content Type                    |
| ------------------------------- | ------------------------------- | ------------------------------- |
| sdkerrors.PaymentsErrorResponse | default                         | application/json                |
| sdkerrors.SDKError              | 4xx-5xx                         | */*                             |

## GetPayment

Get a payment

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
    request := operations.GetPaymentRequest{
        PaymentID: "XXX",
    }
    ctx := context.Background()
    res, err := s.Payments.V1.GetPayment(ctx, request)
    if err != nil {
        log.Fatal(err)
    }
    if res.PaymentResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                        | Type                                                                             | Required                                                                         | Description                                                                      |
| -------------------------------------------------------------------------------- | -------------------------------------------------------------------------------- | -------------------------------------------------------------------------------- | -------------------------------------------------------------------------------- |
| `ctx`                                                                            | [context.Context](https://pkg.go.dev/context#Context)                            | :heavy_check_mark:                                                               | The context to use for the request.                                              |
| `request`                                                                        | [operations.GetPaymentRequest](../../pkg/models/operations/getpaymentrequest.md) | :heavy_check_mark:                                                               | The request object to use for the request.                                       |
| `opts`                                                                           | [][operations.Option](../../pkg/models/operations/option.md)                     | :heavy_minus_sign:                                                               | The options for this request.                                                    |


### Response

**[*operations.GetPaymentResponse](../../pkg/models/operations/getpaymentresponse.md), error**
| Error Object                    | Status Code                     | Content Type                    |
| ------------------------------- | ------------------------------- | ------------------------------- |
| sdkerrors.PaymentsErrorResponse | default                         | application/json                |
| sdkerrors.SDKError              | 4xx-5xx                         | */*                             |

## GetPool

Get a Pool

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
    request := operations.GetPoolRequest{
        PoolID: "XXX",
    }
    ctx := context.Background()
    res, err := s.Payments.V1.GetPool(ctx, request)
    if err != nil {
        log.Fatal(err)
    }
    if res.PoolResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                  | Type                                                                       | Required                                                                   | Description                                                                |
| -------------------------------------------------------------------------- | -------------------------------------------------------------------------- | -------------------------------------------------------------------------- | -------------------------------------------------------------------------- |
| `ctx`                                                                      | [context.Context](https://pkg.go.dev/context#Context)                      | :heavy_check_mark:                                                         | The context to use for the request.                                        |
| `request`                                                                  | [operations.GetPoolRequest](../../pkg/models/operations/getpoolrequest.md) | :heavy_check_mark:                                                         | The request object to use for the request.                                 |
| `opts`                                                                     | [][operations.Option](../../pkg/models/operations/option.md)               | :heavy_minus_sign:                                                         | The options for this request.                                              |


### Response

**[*operations.GetPoolResponse](../../pkg/models/operations/getpoolresponse.md), error**
| Error Object                    | Status Code                     | Content Type                    |
| ------------------------------- | ------------------------------- | ------------------------------- |
| sdkerrors.PaymentsErrorResponse | default                         | application/json                |
| sdkerrors.SDKError              | 4xx-5xx                         | */*                             |

## GetPoolBalances

Get pool balances

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"os"
	"github.com/formancehq/formance-sdk-go/v2"
	"github.com/formancehq/formance-sdk-go/v2/pkg/types"
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
    request := operations.GetPoolBalancesRequest{
        At: types.MustTimeFromString("2023-05-05T06:40:23.018Z"),
        PoolID: "XXX",
    }
    ctx := context.Background()
    res, err := s.Payments.V1.GetPoolBalances(ctx, request)
    if err != nil {
        log.Fatal(err)
    }
    if res.PoolBalancesResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                                  | Type                                                                                       | Required                                                                                   | Description                                                                                |
| ------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------ |
| `ctx`                                                                                      | [context.Context](https://pkg.go.dev/context#Context)                                      | :heavy_check_mark:                                                                         | The context to use for the request.                                                        |
| `request`                                                                                  | [operations.GetPoolBalancesRequest](../../pkg/models/operations/getpoolbalancesrequest.md) | :heavy_check_mark:                                                                         | The request object to use for the request.                                                 |
| `opts`                                                                                     | [][operations.Option](../../pkg/models/operations/option.md)                               | :heavy_minus_sign:                                                                         | The options for this request.                                                              |


### Response

**[*operations.GetPoolBalancesResponse](../../pkg/models/operations/getpoolbalancesresponse.md), error**
| Error Object                    | Status Code                     | Content Type                    |
| ------------------------------- | ------------------------------- | ------------------------------- |
| sdkerrors.PaymentsErrorResponse | default                         | application/json                |
| sdkerrors.SDKError              | 4xx-5xx                         | */*                             |

## GetTransferInitiation

Get a transfer initiation

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
    request := operations.GetTransferInitiationRequest{
        TransferID: "XXX",
    }
    ctx := context.Background()
    res, err := s.Payments.V1.GetTransferInitiation(ctx, request)
    if err != nil {
        log.Fatal(err)
    }
    if res.TransferInitiationResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                                              | Type                                                                                                   | Required                                                                                               | Description                                                                                            |
| ------------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------------ |
| `ctx`                                                                                                  | [context.Context](https://pkg.go.dev/context#Context)                                                  | :heavy_check_mark:                                                                                     | The context to use for the request.                                                                    |
| `request`                                                                                              | [operations.GetTransferInitiationRequest](../../pkg/models/operations/gettransferinitiationrequest.md) | :heavy_check_mark:                                                                                     | The request object to use for the request.                                                             |
| `opts`                                                                                                 | [][operations.Option](../../pkg/models/operations/option.md)                                           | :heavy_minus_sign:                                                                                     | The options for this request.                                                                          |


### Response

**[*operations.GetTransferInitiationResponse](../../pkg/models/operations/gettransferinitiationresponse.md), error**
| Error Object                    | Status Code                     | Content Type                    |
| ------------------------------- | ------------------------------- | ------------------------------- |
| sdkerrors.PaymentsErrorResponse | default                         | application/json                |
| sdkerrors.SDKError              | 4xx-5xx                         | */*                             |

## InstallConnector

Install a connector by its name and config.

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
    request := operations.InstallConnectorRequest{
        ConnectorConfig: shared.CreateConnectorConfigWiseConfig(
                shared.WiseConfig{
                    APIKey: "XXX",
                    Name: "My Wise Account",
                    PollingPeriod: v2.String("60s"),
                },
        ),
        Connector: shared.ConnectorAdyen,
    }
    ctx := context.Background()
    res, err := s.Payments.V1.InstallConnector(ctx, request)
    if err != nil {
        log.Fatal(err)
    }
    if res.ConnectorResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                                    | Type                                                                                         | Required                                                                                     | Description                                                                                  |
| -------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------- |
| `ctx`                                                                                        | [context.Context](https://pkg.go.dev/context#Context)                                        | :heavy_check_mark:                                                                           | The context to use for the request.                                                          |
| `request`                                                                                    | [operations.InstallConnectorRequest](../../pkg/models/operations/installconnectorrequest.md) | :heavy_check_mark:                                                                           | The request object to use for the request.                                                   |
| `opts`                                                                                       | [][operations.Option](../../pkg/models/operations/option.md)                                 | :heavy_minus_sign:                                                                           | The options for this request.                                                                |


### Response

**[*operations.InstallConnectorResponse](../../pkg/models/operations/installconnectorresponse.md), error**
| Error Object                    | Status Code                     | Content Type                    |
| ------------------------------- | ------------------------------- | ------------------------------- |
| sdkerrors.PaymentsErrorResponse | default                         | application/json                |
| sdkerrors.SDKError              | 4xx-5xx                         | */*                             |

## ListAllConnectors

List all installed connectors.

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
    res, err := s.Payments.V1.ListAllConnectors(ctx)
    if err != nil {
        log.Fatal(err)
    }
    if res.ConnectorsResponse != nil {
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

**[*operations.ListAllConnectorsResponse](../../pkg/models/operations/listallconnectorsresponse.md), error**
| Error Object                    | Status Code                     | Content Type                    |
| ------------------------------- | ------------------------------- | ------------------------------- |
| sdkerrors.PaymentsErrorResponse | default                         | application/json                |
| sdkerrors.SDKError              | 4xx-5xx                         | */*                             |

## ListBankAccounts

List all bank accounts created by user on Formance.

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
    request := operations.ListBankAccountsRequest{
        Cursor: v2.String("aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ=="),
        PageSize: v2.Int64(100),
        Sort: []string{
            "date:asc",
            "status:desc",
        },
    }
    ctx := context.Background()
    res, err := s.Payments.V1.ListBankAccounts(ctx, request)
    if err != nil {
        log.Fatal(err)
    }
    if res.BankAccountsCursor != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                                    | Type                                                                                         | Required                                                                                     | Description                                                                                  |
| -------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------- |
| `ctx`                                                                                        | [context.Context](https://pkg.go.dev/context#Context)                                        | :heavy_check_mark:                                                                           | The context to use for the request.                                                          |
| `request`                                                                                    | [operations.ListBankAccountsRequest](../../pkg/models/operations/listbankaccountsrequest.md) | :heavy_check_mark:                                                                           | The request object to use for the request.                                                   |
| `opts`                                                                                       | [][operations.Option](../../pkg/models/operations/option.md)                                 | :heavy_minus_sign:                                                                           | The options for this request.                                                                |


### Response

**[*operations.ListBankAccountsResponse](../../pkg/models/operations/listbankaccountsresponse.md), error**
| Error Object                    | Status Code                     | Content Type                    |
| ------------------------------- | ------------------------------- | ------------------------------- |
| sdkerrors.PaymentsErrorResponse | default                         | application/json                |
| sdkerrors.SDKError              | 4xx-5xx                         | */*                             |

## ListConfigsAvailableConnectors

List the configs of each available connector.

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
    res, err := s.Payments.V1.ListConfigsAvailableConnectors(ctx)
    if err != nil {
        log.Fatal(err)
    }
    if res.ConnectorsConfigsResponse != nil {
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

**[*operations.ListConfigsAvailableConnectorsResponse](../../pkg/models/operations/listconfigsavailableconnectorsresponse.md), error**
| Error Object                    | Status Code                     | Content Type                    |
| ------------------------------- | ------------------------------- | ------------------------------- |
| sdkerrors.PaymentsErrorResponse | default                         | application/json                |
| sdkerrors.SDKError              | 4xx-5xx                         | */*                             |

## ~~ListConnectorTasks~~

List all tasks associated with this connector.

> :warning: **DEPRECATED**: This will be removed in a future release, please migrate away from it as soon as possible.

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
    request := operations.ListConnectorTasksRequest{
        Connector: shared.ConnectorModulr,
        Cursor: v2.String("aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ=="),
        PageSize: v2.Int64(100),
    }
    ctx := context.Background()
    res, err := s.Payments.V1.ListConnectorTasks(ctx, request)
    if err != nil {
        log.Fatal(err)
    }
    if res.TasksCursor != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                                        | Type                                                                                             | Required                                                                                         | Description                                                                                      |
| ------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------ |
| `ctx`                                                                                            | [context.Context](https://pkg.go.dev/context#Context)                                            | :heavy_check_mark:                                                                               | The context to use for the request.                                                              |
| `request`                                                                                        | [operations.ListConnectorTasksRequest](../../pkg/models/operations/listconnectortasksrequest.md) | :heavy_check_mark:                                                                               | The request object to use for the request.                                                       |
| `opts`                                                                                           | [][operations.Option](../../pkg/models/operations/option.md)                                     | :heavy_minus_sign:                                                                               | The options for this request.                                                                    |


### Response

**[*operations.ListConnectorTasksResponse](../../pkg/models/operations/listconnectortasksresponse.md), error**
| Error Object                    | Status Code                     | Content Type                    |
| ------------------------------- | ------------------------------- | ------------------------------- |
| sdkerrors.PaymentsErrorResponse | default                         | application/json                |
| sdkerrors.SDKError              | 4xx-5xx                         | */*                             |

## ListConnectorTasksV1

List all tasks associated with this connector.

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
    request := operations.ListConnectorTasksV1Request{
        Connector: shared.ConnectorBankingCircle,
        ConnectorID: "XXX",
        Cursor: v2.String("aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ=="),
        PageSize: v2.Int64(100),
    }
    ctx := context.Background()
    res, err := s.Payments.V1.ListConnectorTasksV1(ctx, request)
    if err != nil {
        log.Fatal(err)
    }
    if res.TasksCursor != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                                            | Type                                                                                                 | Required                                                                                             | Description                                                                                          |
| ---------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------- |
| `ctx`                                                                                                | [context.Context](https://pkg.go.dev/context#Context)                                                | :heavy_check_mark:                                                                                   | The context to use for the request.                                                                  |
| `request`                                                                                            | [operations.ListConnectorTasksV1Request](../../pkg/models/operations/listconnectortasksv1request.md) | :heavy_check_mark:                                                                                   | The request object to use for the request.                                                           |
| `opts`                                                                                               | [][operations.Option](../../pkg/models/operations/option.md)                                         | :heavy_minus_sign:                                                                                   | The options for this request.                                                                        |


### Response

**[*operations.ListConnectorTasksV1Response](../../pkg/models/operations/listconnectortasksv1response.md), error**
| Error Object                    | Status Code                     | Content Type                    |
| ------------------------------- | ------------------------------- | ------------------------------- |
| sdkerrors.PaymentsErrorResponse | default                         | application/json                |
| sdkerrors.SDKError              | 4xx-5xx                         | */*                             |

## ListPayments

List payments

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
    request := operations.ListPaymentsRequest{
        Cursor: v2.String("aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ=="),
        PageSize: v2.Int64(100),
        Sort: []string{
            "date:asc",
            "status:desc",
        },
    }
    ctx := context.Background()
    res, err := s.Payments.V1.ListPayments(ctx, request)
    if err != nil {
        log.Fatal(err)
    }
    if res.PaymentsCursor != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                            | Type                                                                                 | Required                                                                             | Description                                                                          |
| ------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------ |
| `ctx`                                                                                | [context.Context](https://pkg.go.dev/context#Context)                                | :heavy_check_mark:                                                                   | The context to use for the request.                                                  |
| `request`                                                                            | [operations.ListPaymentsRequest](../../pkg/models/operations/listpaymentsrequest.md) | :heavy_check_mark:                                                                   | The request object to use for the request.                                           |
| `opts`                                                                               | [][operations.Option](../../pkg/models/operations/option.md)                         | :heavy_minus_sign:                                                                   | The options for this request.                                                        |


### Response

**[*operations.ListPaymentsResponse](../../pkg/models/operations/listpaymentsresponse.md), error**
| Error Object                    | Status Code                     | Content Type                    |
| ------------------------------- | ------------------------------- | ------------------------------- |
| sdkerrors.PaymentsErrorResponse | default                         | application/json                |
| sdkerrors.SDKError              | 4xx-5xx                         | */*                             |

## ListPools

List Pools

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
    request := operations.ListPoolsRequest{
        Cursor: v2.String("aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ=="),
        PageSize: v2.Int64(100),
        Sort: []string{
            "date:asc",
            "status:desc",
        },
    }
    ctx := context.Background()
    res, err := s.Payments.V1.ListPools(ctx, request)
    if err != nil {
        log.Fatal(err)
    }
    if res.PoolsCursor != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                      | Type                                                                           | Required                                                                       | Description                                                                    |
| ------------------------------------------------------------------------------ | ------------------------------------------------------------------------------ | ------------------------------------------------------------------------------ | ------------------------------------------------------------------------------ |
| `ctx`                                                                          | [context.Context](https://pkg.go.dev/context#Context)                          | :heavy_check_mark:                                                             | The context to use for the request.                                            |
| `request`                                                                      | [operations.ListPoolsRequest](../../pkg/models/operations/listpoolsrequest.md) | :heavy_check_mark:                                                             | The request object to use for the request.                                     |
| `opts`                                                                         | [][operations.Option](../../pkg/models/operations/option.md)                   | :heavy_minus_sign:                                                             | The options for this request.                                                  |


### Response

**[*operations.ListPoolsResponse](../../pkg/models/operations/listpoolsresponse.md), error**
| Error Object                    | Status Code                     | Content Type                    |
| ------------------------------- | ------------------------------- | ------------------------------- |
| sdkerrors.PaymentsErrorResponse | default                         | application/json                |
| sdkerrors.SDKError              | 4xx-5xx                         | */*                             |

## ListTransferInitiations

List Transfer Initiations

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
    request := operations.ListTransferInitiationsRequest{
        Cursor: v2.String("aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ=="),
        PageSize: v2.Int64(100),
        Sort: []string{
            "date:asc",
            "status:desc",
        },
    }
    ctx := context.Background()
    res, err := s.Payments.V1.ListTransferInitiations(ctx, request)
    if err != nil {
        log.Fatal(err)
    }
    if res.TransferInitiationsCursor != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                                                  | Type                                                                                                       | Required                                                                                                   | Description                                                                                                |
| ---------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------- |
| `ctx`                                                                                                      | [context.Context](https://pkg.go.dev/context#Context)                                                      | :heavy_check_mark:                                                                                         | The context to use for the request.                                                                        |
| `request`                                                                                                  | [operations.ListTransferInitiationsRequest](../../pkg/models/operations/listtransferinitiationsrequest.md) | :heavy_check_mark:                                                                                         | The request object to use for the request.                                                                 |
| `opts`                                                                                                     | [][operations.Option](../../pkg/models/operations/option.md)                                               | :heavy_minus_sign:                                                                                         | The options for this request.                                                                              |


### Response

**[*operations.ListTransferInitiationsResponse](../../pkg/models/operations/listtransferinitiationsresponse.md), error**
| Error Object                    | Status Code                     | Content Type                    |
| ------------------------------- | ------------------------------- | ------------------------------- |
| sdkerrors.PaymentsErrorResponse | default                         | application/json                |
| sdkerrors.SDKError              | 4xx-5xx                         | */*                             |

## PaymentsgetAccount

Get an account

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
    request := operations.PaymentsgetAccountRequest{
        AccountID: "XXX",
    }
    ctx := context.Background()
    res, err := s.Payments.V1.PaymentsgetAccount(ctx, request)
    if err != nil {
        log.Fatal(err)
    }
    if res.PaymentsAccountResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                                        | Type                                                                                             | Required                                                                                         | Description                                                                                      |
| ------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------ |
| `ctx`                                                                                            | [context.Context](https://pkg.go.dev/context#Context)                                            | :heavy_check_mark:                                                                               | The context to use for the request.                                                              |
| `request`                                                                                        | [operations.PaymentsgetAccountRequest](../../pkg/models/operations/paymentsgetaccountrequest.md) | :heavy_check_mark:                                                                               | The request object to use for the request.                                                       |
| `opts`                                                                                           | [][operations.Option](../../pkg/models/operations/option.md)                                     | :heavy_minus_sign:                                                                               | The options for this request.                                                                    |


### Response

**[*operations.PaymentsgetAccountResponse](../../pkg/models/operations/paymentsgetaccountresponse.md), error**
| Error Object                    | Status Code                     | Content Type                    |
| ------------------------------- | ------------------------------- | ------------------------------- |
| sdkerrors.PaymentsErrorResponse | default                         | application/json                |
| sdkerrors.SDKError              | 4xx-5xx                         | */*                             |

## PaymentsgetServerInfo

Get server info

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
    res, err := s.Payments.V1.PaymentsgetServerInfo(ctx)
    if err != nil {
        log.Fatal(err)
    }
    if res.ServerInfo != nil {
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

**[*operations.PaymentsgetServerInfoResponse](../../pkg/models/operations/paymentsgetserverinforesponse.md), error**
| Error Object                    | Status Code                     | Content Type                    |
| ------------------------------- | ------------------------------- | ------------------------------- |
| sdkerrors.PaymentsErrorResponse | default                         | application/json                |
| sdkerrors.SDKError              | 4xx-5xx                         | */*                             |

## PaymentslistAccounts

List accounts

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
    request := operations.PaymentslistAccountsRequest{
        Cursor: v2.String("aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ=="),
        PageSize: v2.Int64(100),
        Sort: []string{
            "date:asc",
            "status:desc",
        },
    }
    ctx := context.Background()
    res, err := s.Payments.V1.PaymentslistAccounts(ctx, request)
    if err != nil {
        log.Fatal(err)
    }
    if res.AccountsCursor != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                                            | Type                                                                                                 | Required                                                                                             | Description                                                                                          |
| ---------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------- |
| `ctx`                                                                                                | [context.Context](https://pkg.go.dev/context#Context)                                                | :heavy_check_mark:                                                                                   | The context to use for the request.                                                                  |
| `request`                                                                                            | [operations.PaymentslistAccountsRequest](../../pkg/models/operations/paymentslistaccountsrequest.md) | :heavy_check_mark:                                                                                   | The request object to use for the request.                                                           |
| `opts`                                                                                               | [][operations.Option](../../pkg/models/operations/option.md)                                         | :heavy_minus_sign:                                                                                   | The options for this request.                                                                        |


### Response

**[*operations.PaymentslistAccountsResponse](../../pkg/models/operations/paymentslistaccountsresponse.md), error**
| Error Object                    | Status Code                     | Content Type                    |
| ------------------------------- | ------------------------------- | ------------------------------- |
| sdkerrors.PaymentsErrorResponse | default                         | application/json                |
| sdkerrors.SDKError              | 4xx-5xx                         | */*                             |

## ~~ReadConnectorConfig~~

Read connector config

> :warning: **DEPRECATED**: This will be removed in a future release, please migrate away from it as soon as possible.

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
    request := operations.ReadConnectorConfigRequest{
        Connector: shared.ConnectorGeneric,
    }
    ctx := context.Background()
    res, err := s.Payments.V1.ReadConnectorConfig(ctx, request)
    if err != nil {
        log.Fatal(err)
    }
    if res.ConnectorConfigResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                                          | Type                                                                                               | Required                                                                                           | Description                                                                                        |
| -------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------- |
| `ctx`                                                                                              | [context.Context](https://pkg.go.dev/context#Context)                                              | :heavy_check_mark:                                                                                 | The context to use for the request.                                                                |
| `request`                                                                                          | [operations.ReadConnectorConfigRequest](../../pkg/models/operations/readconnectorconfigrequest.md) | :heavy_check_mark:                                                                                 | The request object to use for the request.                                                         |
| `opts`                                                                                             | [][operations.Option](../../pkg/models/operations/option.md)                                       | :heavy_minus_sign:                                                                                 | The options for this request.                                                                      |


### Response

**[*operations.ReadConnectorConfigResponse](../../pkg/models/operations/readconnectorconfigresponse.md), error**
| Error Object                    | Status Code                     | Content Type                    |
| ------------------------------- | ------------------------------- | ------------------------------- |
| sdkerrors.PaymentsErrorResponse | default                         | application/json                |
| sdkerrors.SDKError              | 4xx-5xx                         | */*                             |

## ReadConnectorConfigV1

Read connector config

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
    request := operations.ReadConnectorConfigV1Request{
        Connector: shared.ConnectorCurrencyCloud,
        ConnectorID: "XXX",
    }
    ctx := context.Background()
    res, err := s.Payments.V1.ReadConnectorConfigV1(ctx, request)
    if err != nil {
        log.Fatal(err)
    }
    if res.ConnectorConfigResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                                              | Type                                                                                                   | Required                                                                                               | Description                                                                                            |
| ------------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------------ |
| `ctx`                                                                                                  | [context.Context](https://pkg.go.dev/context#Context)                                                  | :heavy_check_mark:                                                                                     | The context to use for the request.                                                                    |
| `request`                                                                                              | [operations.ReadConnectorConfigV1Request](../../pkg/models/operations/readconnectorconfigv1request.md) | :heavy_check_mark:                                                                                     | The request object to use for the request.                                                             |
| `opts`                                                                                                 | [][operations.Option](../../pkg/models/operations/option.md)                                           | :heavy_minus_sign:                                                                                     | The options for this request.                                                                          |


### Response

**[*operations.ReadConnectorConfigV1Response](../../pkg/models/operations/readconnectorconfigv1response.md), error**
| Error Object                    | Status Code                     | Content Type                    |
| ------------------------------- | ------------------------------- | ------------------------------- |
| sdkerrors.PaymentsErrorResponse | default                         | application/json                |
| sdkerrors.SDKError              | 4xx-5xx                         | */*                             |

## RemoveAccountFromPool

Remove an account from a pool by its id.

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
    request := operations.RemoveAccountFromPoolRequest{
        AccountID: "XXX",
        PoolID: "XXX",
    }
    ctx := context.Background()
    res, err := s.Payments.V1.RemoveAccountFromPool(ctx, request)
    if err != nil {
        log.Fatal(err)
    }
    if res != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                                              | Type                                                                                                   | Required                                                                                               | Description                                                                                            |
| ------------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------------ |
| `ctx`                                                                                                  | [context.Context](https://pkg.go.dev/context#Context)                                                  | :heavy_check_mark:                                                                                     | The context to use for the request.                                                                    |
| `request`                                                                                              | [operations.RemoveAccountFromPoolRequest](../../pkg/models/operations/removeaccountfrompoolrequest.md) | :heavy_check_mark:                                                                                     | The request object to use for the request.                                                             |
| `opts`                                                                                                 | [][operations.Option](../../pkg/models/operations/option.md)                                           | :heavy_minus_sign:                                                                                     | The options for this request.                                                                          |


### Response

**[*operations.RemoveAccountFromPoolResponse](../../pkg/models/operations/removeaccountfrompoolresponse.md), error**
| Error Object                    | Status Code                     | Content Type                    |
| ------------------------------- | ------------------------------- | ------------------------------- |
| sdkerrors.PaymentsErrorResponse | default                         | application/json                |
| sdkerrors.SDKError              | 4xx-5xx                         | */*                             |

## ~~ResetConnector~~

Reset a connector by its name.
It will remove the connector and ALL PAYMENTS generated with it.


> :warning: **DEPRECATED**: This will be removed in a future release, please migrate away from it as soon as possible.

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
    request := operations.ResetConnectorRequest{
        Connector: shared.ConnectorAtlar,
    }
    ctx := context.Background()
    res, err := s.Payments.V1.ResetConnector(ctx, request)
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
| `request`                                                                                | [operations.ResetConnectorRequest](../../pkg/models/operations/resetconnectorrequest.md) | :heavy_check_mark:                                                                       | The request object to use for the request.                                               |
| `opts`                                                                                   | [][operations.Option](../../pkg/models/operations/option.md)                             | :heavy_minus_sign:                                                                       | The options for this request.                                                            |


### Response

**[*operations.ResetConnectorResponse](../../pkg/models/operations/resetconnectorresponse.md), error**
| Error Object                    | Status Code                     | Content Type                    |
| ------------------------------- | ------------------------------- | ------------------------------- |
| sdkerrors.PaymentsErrorResponse | default                         | application/json                |
| sdkerrors.SDKError              | 4xx-5xx                         | */*                             |

## ResetConnectorV1

Reset a connector by its name.
It will remove the connector and ALL PAYMENTS generated with it.


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
    request := operations.ResetConnectorV1Request{
        Connector: shared.ConnectorGeneric,
        ConnectorID: "XXX",
    }
    ctx := context.Background()
    res, err := s.Payments.V1.ResetConnectorV1(ctx, request)
    if err != nil {
        log.Fatal(err)
    }
    if res != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                                    | Type                                                                                         | Required                                                                                     | Description                                                                                  |
| -------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------- |
| `ctx`                                                                                        | [context.Context](https://pkg.go.dev/context#Context)                                        | :heavy_check_mark:                                                                           | The context to use for the request.                                                          |
| `request`                                                                                    | [operations.ResetConnectorV1Request](../../pkg/models/operations/resetconnectorv1request.md) | :heavy_check_mark:                                                                           | The request object to use for the request.                                                   |
| `opts`                                                                                       | [][operations.Option](../../pkg/models/operations/option.md)                                 | :heavy_minus_sign:                                                                           | The options for this request.                                                                |


### Response

**[*operations.ResetConnectorV1Response](../../pkg/models/operations/resetconnectorv1response.md), error**
| Error Object                    | Status Code                     | Content Type                    |
| ------------------------------- | ------------------------------- | ------------------------------- |
| sdkerrors.PaymentsErrorResponse | default                         | application/json                |
| sdkerrors.SDKError              | 4xx-5xx                         | */*                             |

## RetryTransferInitiation

Retry a failed transfer initiation

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
    request := operations.RetryTransferInitiationRequest{
        TransferID: "XXX",
    }
    ctx := context.Background()
    res, err := s.Payments.V1.RetryTransferInitiation(ctx, request)
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
| `request`                                                                                                  | [operations.RetryTransferInitiationRequest](../../pkg/models/operations/retrytransferinitiationrequest.md) | :heavy_check_mark:                                                                                         | The request object to use for the request.                                                                 |
| `opts`                                                                                                     | [][operations.Option](../../pkg/models/operations/option.md)                                               | :heavy_minus_sign:                                                                                         | The options for this request.                                                                              |


### Response

**[*operations.RetryTransferInitiationResponse](../../pkg/models/operations/retrytransferinitiationresponse.md), error**
| Error Object                    | Status Code                     | Content Type                    |
| ------------------------------- | ------------------------------- | ------------------------------- |
| sdkerrors.PaymentsErrorResponse | default                         | application/json                |
| sdkerrors.SDKError              | 4xx-5xx                         | */*                             |

## ReverseTransferInitiation

Reverse transfer initiation

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
    request := operations.ReverseTransferInitiationRequest{
        ReverseTransferInitiationRequest: shared.ReverseTransferInitiationRequest{
            Amount: big.NewInt(327549),
            Asset: "USD",
            Description: "Streamlined high-level local area network",
            Metadata: map[string]string{
                "key": "<value>",
            },
            Reference: "XXX",
        },
        TransferID: "XXX",
    }
    ctx := context.Background()
    res, err := s.Payments.V1.ReverseTransferInitiation(ctx, request)
    if err != nil {
        log.Fatal(err)
    }
    if res != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                                                      | Type                                                                                                           | Required                                                                                                       | Description                                                                                                    |
| -------------------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------- |
| `ctx`                                                                                                          | [context.Context](https://pkg.go.dev/context#Context)                                                          | :heavy_check_mark:                                                                                             | The context to use for the request.                                                                            |
| `request`                                                                                                      | [operations.ReverseTransferInitiationRequest](../../pkg/models/operations/reversetransferinitiationrequest.md) | :heavy_check_mark:                                                                                             | The request object to use for the request.                                                                     |
| `opts`                                                                                                         | [][operations.Option](../../pkg/models/operations/option.md)                                                   | :heavy_minus_sign:                                                                                             | The options for this request.                                                                                  |


### Response

**[*operations.ReverseTransferInitiationResponse](../../pkg/models/operations/reversetransferinitiationresponse.md), error**
| Error Object                    | Status Code                     | Content Type                    |
| ------------------------------- | ------------------------------- | ------------------------------- |
| sdkerrors.PaymentsErrorResponse | default                         | application/json                |
| sdkerrors.SDKError              | 4xx-5xx                         | */*                             |

## UdpateTransferInitiationStatus

Update a transfer initiation status

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
    request := operations.UdpateTransferInitiationStatusRequest{
        UpdateTransferInitiationStatusRequest: shared.UpdateTransferInitiationStatusRequest{
            Status: shared.StatusValidated,
        },
        TransferID: "XXX",
    }
    ctx := context.Background()
    res, err := s.Payments.V1.UdpateTransferInitiationStatus(ctx, request)
    if err != nil {
        log.Fatal(err)
    }
    if res != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                                                                | Type                                                                                                                     | Required                                                                                                                 | Description                                                                                                              |
| ------------------------------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------------------------------ |
| `ctx`                                                                                                                    | [context.Context](https://pkg.go.dev/context#Context)                                                                    | :heavy_check_mark:                                                                                                       | The context to use for the request.                                                                                      |
| `request`                                                                                                                | [operations.UdpateTransferInitiationStatusRequest](../../pkg/models/operations/udpatetransferinitiationstatusrequest.md) | :heavy_check_mark:                                                                                                       | The request object to use for the request.                                                                               |
| `opts`                                                                                                                   | [][operations.Option](../../pkg/models/operations/option.md)                                                             | :heavy_minus_sign:                                                                                                       | The options for this request.                                                                                            |


### Response

**[*operations.UdpateTransferInitiationStatusResponse](../../pkg/models/operations/udpatetransferinitiationstatusresponse.md), error**
| Error Object                    | Status Code                     | Content Type                    |
| ------------------------------- | ------------------------------- | ------------------------------- |
| sdkerrors.PaymentsErrorResponse | default                         | application/json                |
| sdkerrors.SDKError              | 4xx-5xx                         | */*                             |

## ~~UninstallConnector~~

Uninstall a connector by its name.

> :warning: **DEPRECATED**: This will be removed in a future release, please migrate away from it as soon as possible.

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
    request := operations.UninstallConnectorRequest{
        Connector: shared.ConnectorModulr,
    }
    ctx := context.Background()
    res, err := s.Payments.V1.UninstallConnector(ctx, request)
    if err != nil {
        log.Fatal(err)
    }
    if res != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                                        | Type                                                                                             | Required                                                                                         | Description                                                                                      |
| ------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------ |
| `ctx`                                                                                            | [context.Context](https://pkg.go.dev/context#Context)                                            | :heavy_check_mark:                                                                               | The context to use for the request.                                                              |
| `request`                                                                                        | [operations.UninstallConnectorRequest](../../pkg/models/operations/uninstallconnectorrequest.md) | :heavy_check_mark:                                                                               | The request object to use for the request.                                                       |
| `opts`                                                                                           | [][operations.Option](../../pkg/models/operations/option.md)                                     | :heavy_minus_sign:                                                                               | The options for this request.                                                                    |


### Response

**[*operations.UninstallConnectorResponse](../../pkg/models/operations/uninstallconnectorresponse.md), error**
| Error Object                    | Status Code                     | Content Type                    |
| ------------------------------- | ------------------------------- | ------------------------------- |
| sdkerrors.PaymentsErrorResponse | default                         | application/json                |
| sdkerrors.SDKError              | 4xx-5xx                         | */*                             |

## UninstallConnectorV1

Uninstall a connector by its name.

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
    request := operations.UninstallConnectorV1Request{
        Connector: shared.ConnectorGeneric,
        ConnectorID: "XXX",
    }
    ctx := context.Background()
    res, err := s.Payments.V1.UninstallConnectorV1(ctx, request)
    if err != nil {
        log.Fatal(err)
    }
    if res != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                                            | Type                                                                                                 | Required                                                                                             | Description                                                                                          |
| ---------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------- |
| `ctx`                                                                                                | [context.Context](https://pkg.go.dev/context#Context)                                                | :heavy_check_mark:                                                                                   | The context to use for the request.                                                                  |
| `request`                                                                                            | [operations.UninstallConnectorV1Request](../../pkg/models/operations/uninstallconnectorv1request.md) | :heavy_check_mark:                                                                                   | The request object to use for the request.                                                           |
| `opts`                                                                                               | [][operations.Option](../../pkg/models/operations/option.md)                                         | :heavy_minus_sign:                                                                                   | The options for this request.                                                                        |


### Response

**[*operations.UninstallConnectorV1Response](../../pkg/models/operations/uninstallconnectorv1response.md), error**
| Error Object                    | Status Code                     | Content Type                    |
| ------------------------------- | ------------------------------- | ------------------------------- |
| sdkerrors.PaymentsErrorResponse | default                         | application/json                |
| sdkerrors.SDKError              | 4xx-5xx                         | */*                             |

## UpdateBankAccountMetadata

Update metadata of a bank account

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
    request := operations.UpdateBankAccountMetadataRequest{
        UpdateBankAccountMetadataRequest: shared.UpdateBankAccountMetadataRequest{
            Metadata: map[string]string{
                "key": "<value>",
            },
        },
        BankAccountID: "XXX",
    }
    ctx := context.Background()
    res, err := s.Payments.V1.UpdateBankAccountMetadata(ctx, request)
    if err != nil {
        log.Fatal(err)
    }
    if res != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                                                      | Type                                                                                                           | Required                                                                                                       | Description                                                                                                    |
| -------------------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------- |
| `ctx`                                                                                                          | [context.Context](https://pkg.go.dev/context#Context)                                                          | :heavy_check_mark:                                                                                             | The context to use for the request.                                                                            |
| `request`                                                                                                      | [operations.UpdateBankAccountMetadataRequest](../../pkg/models/operations/updatebankaccountmetadatarequest.md) | :heavy_check_mark:                                                                                             | The request object to use for the request.                                                                     |
| `opts`                                                                                                         | [][operations.Option](../../pkg/models/operations/option.md)                                                   | :heavy_minus_sign:                                                                                             | The options for this request.                                                                                  |


### Response

**[*operations.UpdateBankAccountMetadataResponse](../../pkg/models/operations/updatebankaccountmetadataresponse.md), error**
| Error Object                    | Status Code                     | Content Type                    |
| ------------------------------- | ------------------------------- | ------------------------------- |
| sdkerrors.PaymentsErrorResponse | default                         | application/json                |
| sdkerrors.SDKError              | 4xx-5xx                         | */*                             |

## UpdateConnectorConfigV1

Update connector config

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
    request := operations.UpdateConnectorConfigV1Request{
        ConnectorConfig: shared.CreateConnectorConfigStripeConfig(
                shared.StripeConfig{
                    APIKey: "XXX",
                    Name: "My Stripe Account",
                    PageSize: v2.Int64(50),
                    PollingPeriod: v2.String("60s"),
                },
        ),
        Connector: shared.ConnectorStripe,
        ConnectorID: "XXX",
    }
    ctx := context.Background()
    res, err := s.Payments.V1.UpdateConnectorConfigV1(ctx, request)
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
| `request`                                                                                                  | [operations.UpdateConnectorConfigV1Request](../../pkg/models/operations/updateconnectorconfigv1request.md) | :heavy_check_mark:                                                                                         | The request object to use for the request.                                                                 |
| `opts`                                                                                                     | [][operations.Option](../../pkg/models/operations/option.md)                                               | :heavy_minus_sign:                                                                                         | The options for this request.                                                                              |


### Response

**[*operations.UpdateConnectorConfigV1Response](../../pkg/models/operations/updateconnectorconfigv1response.md), error**
| Error Object                    | Status Code                     | Content Type                    |
| ------------------------------- | ------------------------------- | ------------------------------- |
| sdkerrors.PaymentsErrorResponse | default                         | application/json                |
| sdkerrors.SDKError              | 4xx-5xx                         | */*                             |

## UpdateMetadata

Update metadata

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
    request := operations.UpdateMetadataRequest{
        RequestBody: map[string]string{
            "key": "<value>",
        },
        PaymentID: "XXX",
    }
    ctx := context.Background()
    res, err := s.Payments.V1.UpdateMetadata(ctx, request)
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
| `request`                                                                                | [operations.UpdateMetadataRequest](../../pkg/models/operations/updatemetadatarequest.md) | :heavy_check_mark:                                                                       | The request object to use for the request.                                               |
| `opts`                                                                                   | [][operations.Option](../../pkg/models/operations/option.md)                             | :heavy_minus_sign:                                                                       | The options for this request.                                                            |


### Response

**[*operations.UpdateMetadataResponse](../../pkg/models/operations/updatemetadataresponse.md), error**
| Error Object                    | Status Code                     | Content Type                    |
| ------------------------------- | ------------------------------- | ------------------------------- |
| sdkerrors.PaymentsErrorResponse | default                         | application/json                |
| sdkerrors.SDKError              | 4xx-5xx                         | */*                             |
