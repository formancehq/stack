# Payments

### Available Operations

* [ConnectorsTransfer](#connectorstransfer) - Transfer funds between Connector accounts
* [CreateBankAccount](#createbankaccount) - Create a BankAccount in Payments and on the PSP
* [CreateTransferInitiation](#createtransferinitiation) - Create a TransferInitiation
* [DeleteTransferInitiation](#deletetransferinitiation) - Delete a transfer initiation
* [GetAccountBalances](#getaccountbalances) - Get account balances
* [GetBankAccount](#getbankaccount) - Get a bank account created by user on Formance
* [GetConnectorTask](#getconnectortask) - Read a specific task of the connector
* [GetPayment](#getpayment) - Get a payment
* [GetTransferInitiation](#gettransferinitiation) - Get a transfer initiation
* [InstallConnector](#installconnector) - Install a connector
* [ListAllConnectors](#listallconnectors) - List all installed connectors
* [ListBankAccounts](#listbankaccounts) - List bank accounts created by user on Formance
* [ListConfigsAvailableConnectors](#listconfigsavailableconnectors) - List the configs of each available connector
* [ListConnectorTasks](#listconnectortasks) - List tasks from a connector
* [ListConnectorsTransfers](#listconnectorstransfers) - List transfers and their statuses
* [ListPayments](#listpayments) - List payments
* [ListTransferInitiations](#listtransferinitiations) - List Transfer Initiations
* [PaymentsgetAccount](#paymentsgetaccount) - Get an account
* [PaymentsgetServerInfo](#paymentsgetserverinfo) - Get server info
* [PaymentslistAccounts](#paymentslistaccounts) - List accounts
* [ReadConnectorConfig](#readconnectorconfig) - Read the config of a connector
* [ResetConnector](#resetconnector) - Reset a connector
* [RetryTransferInitiation](#retrytransferinitiation) - Retry a failed transfer initiation
* [UdpateTransferInitiationStatus](#udpatetransferinitiationstatus) - Update the status of a transfer initiation
* [UninstallConnector](#uninstallconnector) - Uninstall a connector
* [UpdateMetadata](#updatemetadata) - Update metadata

## ConnectorsTransfer

Execute a transfer between two accounts.

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
    s := formance.New()

    ctx := context.Background()
    res, err := s.Payments.ConnectorsTransfer(ctx, operations.ConnectorsTransferRequest{
        TransferRequest: shared.TransferRequest{
            Amount: big.NewInt(100),
            Asset: "USD",
            Destination: "acct_1Gqj58KZcSIg2N2q",
            Source: formance.String("acct_1Gqj58KZcSIg2N2q"),
        },
        Connector: shared.ConnectorMoneycorp,
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.TransferResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                                    | Type                                                                                         | Required                                                                                     | Description                                                                                  |
| -------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------- |
| `ctx`                                                                                        | [context.Context](https://pkg.go.dev/context#Context)                                        | :heavy_check_mark:                                                                           | The context to use for the request.                                                          |
| `request`                                                                                    | [operations.ConnectorsTransferRequest](../../models/operations/connectorstransferrequest.md) | :heavy_check_mark:                                                                           | The request object to use for the request.                                                   |


### Response

**[*operations.ConnectorsTransferResponse](../../models/operations/connectorstransferresponse.md), error**


## CreateBankAccount

Create a bank account in Payments and on the PSP.

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
    s := formance.New()

    ctx := context.Background()
    res, err := s.Payments.CreateBankAccount(ctx, shared.BankAccountRequest{
        AccountNumber: formance.String("dolorum"),
        Country: "GB",
        Iban: formance.String("in"),
        Name: "My account",
        Provider: shared.ConnectorModulr,
        SwiftBicCode: formance.String("illum"),
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.BankAccountResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                              | Type                                                                   | Required                                                               | Description                                                            |
| ---------------------------------------------------------------------- | ---------------------------------------------------------------------- | ---------------------------------------------------------------------- | ---------------------------------------------------------------------- |
| `ctx`                                                                  | [context.Context](https://pkg.go.dev/context#Context)                  | :heavy_check_mark:                                                     | The context to use for the request.                                    |
| `request`                                                              | [shared.BankAccountRequest](../../models/shared/bankaccountrequest.md) | :heavy_check_mark:                                                     | The request object to use for the request.                             |


### Response

**[*operations.CreateBankAccountResponse](../../models/operations/createbankaccountresponse.md), error**


## CreateTransferInitiation

Create a transfer initiation

### Example Usage

```go
package main

import(
	"context"
	"log"
	"github.com/formancehq/formance-sdk-go"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"math/big"
	"github.com/formancehq/formance-sdk-go/pkg/types"
)

func main() {
    s := formance.New()

    ctx := context.Background()
    res, err := s.Payments.CreateTransferInitiation(ctx, shared.TransferInitiationRequest{
        Amount: big.NewInt(978571),
        Asset: "USD",
        CreatedAt: types.MustTimeFromString("2022-10-08T04:08:32.252Z"),
        Description: "magnam",
        DestinationAccountID: "cumque",
        Provider: shared.ConnectorMangopay,
        Reference: "XXX",
        SourceAccountID: "ea",
        Type: shared.TransferInitiationRequestTypeTransfer,
        Validated: false,
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.TransferInitiationResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                            | Type                                                                                 | Required                                                                             | Description                                                                          |
| ------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------ |
| `ctx`                                                                                | [context.Context](https://pkg.go.dev/context#Context)                                | :heavy_check_mark:                                                                   | The context to use for the request.                                                  |
| `request`                                                                            | [shared.TransferInitiationRequest](../../models/shared/transferinitiationrequest.md) | :heavy_check_mark:                                                                   | The request object to use for the request.                                           |


### Response

**[*operations.CreateTransferInitiationResponse](../../models/operations/createtransferinitiationresponse.md), error**


## DeleteTransferInitiation

Delete a transfer initiation by its id.

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
    s := formance.New()

    ctx := context.Background()
    res, err := s.Payments.DeleteTransferInitiation(ctx, operations.DeleteTransferInitiationRequest{
        TransferID: "laborum",
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
| `request`                                                                                                | [operations.DeleteTransferInitiationRequest](../../models/operations/deletetransferinitiationrequest.md) | :heavy_check_mark:                                                                                       | The request object to use for the request.                                                               |


### Response

**[*operations.DeleteTransferInitiationResponse](../../models/operations/deletetransferinitiationresponse.md), error**


## GetAccountBalances

Get account balances

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
    s := formance.New()

    ctx := context.Background()
    res, err := s.Payments.GetAccountBalances(ctx, operations.GetAccountBalancesRequest{
        AccountID: "accusamus",
        Asset: formance.String("non"),
        Cursor: formance.String("aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ=="),
        From: types.MustTimeFromString("2022-05-17T08:24:52.669Z"),
        Limit: formance.Int64(881736),
        PageSize: formance.Int64(965417),
        Sort: []string{
            "provident",
            "nam",
            "id",
        },
        To: types.MustTimeFromString("2021-12-07T18:13:34.827Z"),
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.BalancesCursor != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                                    | Type                                                                                         | Required                                                                                     | Description                                                                                  |
| -------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------- |
| `ctx`                                                                                        | [context.Context](https://pkg.go.dev/context#Context)                                        | :heavy_check_mark:                                                                           | The context to use for the request.                                                          |
| `request`                                                                                    | [operations.GetAccountBalancesRequest](../../models/operations/getaccountbalancesrequest.md) | :heavy_check_mark:                                                                           | The request object to use for the request.                                                   |


### Response

**[*operations.GetAccountBalancesResponse](../../models/operations/getaccountbalancesresponse.md), error**


## GetBankAccount

Get a bank account created by user on Formance

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
    s := formance.New()

    ctx := context.Background()
    res, err := s.Payments.GetBankAccount(ctx, operations.GetBankAccountRequest{
        BankAccountID: "sapiente",
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.BankAccountResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                            | Type                                                                                 | Required                                                                             | Description                                                                          |
| ------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------ |
| `ctx`                                                                                | [context.Context](https://pkg.go.dev/context#Context)                                | :heavy_check_mark:                                                                   | The context to use for the request.                                                  |
| `request`                                                                            | [operations.GetBankAccountRequest](../../models/operations/getbankaccountrequest.md) | :heavy_check_mark:                                                                   | The request object to use for the request.                                           |


### Response

**[*operations.GetBankAccountResponse](../../models/operations/getbankaccountresponse.md), error**


## GetConnectorTask

Get a specific task associated to the connector.

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
    s := formance.New()

    ctx := context.Background()
    res, err := s.Payments.GetConnectorTask(ctx, operations.GetConnectorTaskRequest{
        Connector: shared.ConnectorDummyPay,
        TaskID: "deserunt",
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.TaskResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                                | Type                                                                                     | Required                                                                                 | Description                                                                              |
| ---------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------- |
| `ctx`                                                                                    | [context.Context](https://pkg.go.dev/context#Context)                                    | :heavy_check_mark:                                                                       | The context to use for the request.                                                      |
| `request`                                                                                | [operations.GetConnectorTaskRequest](../../models/operations/getconnectortaskrequest.md) | :heavy_check_mark:                                                                       | The request object to use for the request.                                               |


### Response

**[*operations.GetConnectorTaskResponse](../../models/operations/getconnectortaskresponse.md), error**


## GetPayment

Get a payment

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
    s := formance.New()

    ctx := context.Background()
    res, err := s.Payments.GetPayment(ctx, operations.GetPaymentRequest{
        PaymentID: "nisi",
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.PaymentResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                    | Type                                                                         | Required                                                                     | Description                                                                  |
| ---------------------------------------------------------------------------- | ---------------------------------------------------------------------------- | ---------------------------------------------------------------------------- | ---------------------------------------------------------------------------- |
| `ctx`                                                                        | [context.Context](https://pkg.go.dev/context#Context)                        | :heavy_check_mark:                                                           | The context to use for the request.                                          |
| `request`                                                                    | [operations.GetPaymentRequest](../../models/operations/getpaymentrequest.md) | :heavy_check_mark:                                                           | The request object to use for the request.                                   |


### Response

**[*operations.GetPaymentResponse](../../models/operations/getpaymentresponse.md), error**


## GetTransferInitiation

Get a transfer initiation

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
    s := formance.New()

    ctx := context.Background()
    res, err := s.Payments.GetTransferInitiation(ctx, operations.GetTransferInitiationRequest{
        TransferID: "vel",
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.TransferInitiationResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                                          | Type                                                                                               | Required                                                                                           | Description                                                                                        |
| -------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------- |
| `ctx`                                                                                              | [context.Context](https://pkg.go.dev/context#Context)                                              | :heavy_check_mark:                                                                                 | The context to use for the request.                                                                |
| `request`                                                                                          | [operations.GetTransferInitiationRequest](../../models/operations/gettransferinitiationrequest.md) | :heavy_check_mark:                                                                                 | The request object to use for the request.                                                         |


### Response

**[*operations.GetTransferInitiationResponse](../../models/operations/gettransferinitiationresponse.md), error**


## InstallConnector

Install a connector by its name and config.

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
    s := formance.New()

    ctx := context.Background()
    res, err := s.Payments.InstallConnector(ctx, operations.InstallConnectorRequest{
        RequestBody: shared.CurrencyCloudConfig{
            APIKey: "XXX",
            Endpoint: formance.String("XXX"),
            LoginID: "XXX",
            PollingPeriod: formance.String("60s"),
        },
        Connector: shared.ConnectorCurrencyCloud,
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
| `request`                                                                                | [operations.InstallConnectorRequest](../../models/operations/installconnectorrequest.md) | :heavy_check_mark:                                                                       | The request object to use for the request.                                               |


### Response

**[*operations.InstallConnectorResponse](../../models/operations/installconnectorresponse.md), error**


## ListAllConnectors

List all installed connectors.

### Example Usage

```go
package main

import(
	"context"
	"log"
	"github.com/formancehq/formance-sdk-go"
)

func main() {
    s := formance.New()

    ctx := context.Background()
    res, err := s.Payments.ListAllConnectors(ctx)
    if err != nil {
        log.Fatal(err)
    }

    if res.ConnectorsResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                             | Type                                                  | Required                                              | Description                                           |
| ----------------------------------------------------- | ----------------------------------------------------- | ----------------------------------------------------- | ----------------------------------------------------- |
| `ctx`                                                 | [context.Context](https://pkg.go.dev/context#Context) | :heavy_check_mark:                                    | The context to use for the request.                   |


### Response

**[*operations.ListAllConnectorsResponse](../../models/operations/listallconnectorsresponse.md), error**


## ListBankAccounts

List all bank accounts created by user on Formance.

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
    s := formance.New()

    ctx := context.Background()
    res, err := s.Payments.ListBankAccounts(ctx, operations.ListBankAccountsRequest{
        Cursor: formance.String("aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ=="),
        PageSize: formance.Int64(474867),
        Sort: []string{
            "nihil",
        },
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.BankAccountsCursor != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                                | Type                                                                                     | Required                                                                                 | Description                                                                              |
| ---------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------- |
| `ctx`                                                                                    | [context.Context](https://pkg.go.dev/context#Context)                                    | :heavy_check_mark:                                                                       | The context to use for the request.                                                      |
| `request`                                                                                | [operations.ListBankAccountsRequest](../../models/operations/listbankaccountsrequest.md) | :heavy_check_mark:                                                                       | The request object to use for the request.                                               |


### Response

**[*operations.ListBankAccountsResponse](../../models/operations/listbankaccountsresponse.md), error**


## ListConfigsAvailableConnectors

List the configs of each available connector.

### Example Usage

```go
package main

import(
	"context"
	"log"
	"github.com/formancehq/formance-sdk-go"
)

func main() {
    s := formance.New()

    ctx := context.Background()
    res, err := s.Payments.ListConfigsAvailableConnectors(ctx)
    if err != nil {
        log.Fatal(err)
    }

    if res.ConnectorsConfigsResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                             | Type                                                  | Required                                              | Description                                           |
| ----------------------------------------------------- | ----------------------------------------------------- | ----------------------------------------------------- | ----------------------------------------------------- |
| `ctx`                                                 | [context.Context](https://pkg.go.dev/context#Context) | :heavy_check_mark:                                    | The context to use for the request.                   |


### Response

**[*operations.ListConfigsAvailableConnectorsResponse](../../models/operations/listconfigsavailableconnectorsresponse.md), error**


## ListConnectorTasks

List all tasks associated with this connector.

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
    s := formance.New()

    ctx := context.Background()
    res, err := s.Payments.ListConnectorTasks(ctx, operations.ListConnectorTasksRequest{
        Connector: shared.ConnectorWise,
        Cursor: formance.String("aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ=="),
        PageSize: formance.Int64(716075),
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.TasksCursor != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                                    | Type                                                                                         | Required                                                                                     | Description                                                                                  |
| -------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------- |
| `ctx`                                                                                        | [context.Context](https://pkg.go.dev/context#Context)                                        | :heavy_check_mark:                                                                           | The context to use for the request.                                                          |
| `request`                                                                                    | [operations.ListConnectorTasksRequest](../../models/operations/listconnectortasksrequest.md) | :heavy_check_mark:                                                                           | The request object to use for the request.                                                   |


### Response

**[*operations.ListConnectorTasksResponse](../../models/operations/listconnectortasksresponse.md), error**


## ListConnectorsTransfers

List transfers

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
    s := formance.New()

    ctx := context.Background()
    res, err := s.Payments.ListConnectorsTransfers(ctx, operations.ListConnectorsTransfersRequest{
        Connector: shared.ConnectorBankingCircle,
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.TransfersResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                                              | Type                                                                                                   | Required                                                                                               | Description                                                                                            |
| ------------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------------ |
| `ctx`                                                                                                  | [context.Context](https://pkg.go.dev/context#Context)                                                  | :heavy_check_mark:                                                                                     | The context to use for the request.                                                                    |
| `request`                                                                                              | [operations.ListConnectorsTransfersRequest](../../models/operations/listconnectorstransfersrequest.md) | :heavy_check_mark:                                                                                     | The request object to use for the request.                                                             |


### Response

**[*operations.ListConnectorsTransfersResponse](../../models/operations/listconnectorstransfersresponse.md), error**


## ListPayments

List payments

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
    s := formance.New()

    ctx := context.Background()
    res, err := s.Payments.ListPayments(ctx, operations.ListPaymentsRequest{
        Cursor: formance.String("aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ=="),
        PageSize: formance.Int64(287991),
        Sort: []string{
            "suscipit",
            "natus",
        },
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.PaymentsCursor != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                        | Type                                                                             | Required                                                                         | Description                                                                      |
| -------------------------------------------------------------------------------- | -------------------------------------------------------------------------------- | -------------------------------------------------------------------------------- | -------------------------------------------------------------------------------- |
| `ctx`                                                                            | [context.Context](https://pkg.go.dev/context#Context)                            | :heavy_check_mark:                                                               | The context to use for the request.                                              |
| `request`                                                                        | [operations.ListPaymentsRequest](../../models/operations/listpaymentsrequest.md) | :heavy_check_mark:                                                               | The request object to use for the request.                                       |


### Response

**[*operations.ListPaymentsResponse](../../models/operations/listpaymentsresponse.md), error**


## ListTransferInitiations

List Transfer Initiations

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
    s := formance.New()

    ctx := context.Background()
    res, err := s.Payments.ListTransferInitiations(ctx, operations.ListTransferInitiationsRequest{
        Cursor: formance.String("aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ=="),
        PageSize: formance.Int64(749170),
        Query: formance.String("eum"),
        Sort: []string{
            "aspernatur",
            "architecto",
            "magnam",
            "et",
        },
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.TransferInitiationsCursor != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                                              | Type                                                                                                   | Required                                                                                               | Description                                                                                            |
| ------------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------------ |
| `ctx`                                                                                                  | [context.Context](https://pkg.go.dev/context#Context)                                                  | :heavy_check_mark:                                                                                     | The context to use for the request.                                                                    |
| `request`                                                                                              | [operations.ListTransferInitiationsRequest](../../models/operations/listtransferinitiationsrequest.md) | :heavy_check_mark:                                                                                     | The request object to use for the request.                                                             |


### Response

**[*operations.ListTransferInitiationsResponse](../../models/operations/listtransferinitiationsresponse.md), error**


## PaymentsgetAccount

Get an account

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
    s := formance.New()

    ctx := context.Background()
    res, err := s.Payments.PaymentsgetAccount(ctx, operations.PaymentsgetAccountRequest{
        AccountID: "excepturi",
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.PaymentsAccountResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                                    | Type                                                                                         | Required                                                                                     | Description                                                                                  |
| -------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------- |
| `ctx`                                                                                        | [context.Context](https://pkg.go.dev/context#Context)                                        | :heavy_check_mark:                                                                           | The context to use for the request.                                                          |
| `request`                                                                                    | [operations.PaymentsgetAccountRequest](../../models/operations/paymentsgetaccountrequest.md) | :heavy_check_mark:                                                                           | The request object to use for the request.                                                   |


### Response

**[*operations.PaymentsgetAccountResponse](../../models/operations/paymentsgetaccountresponse.md), error**


## PaymentsgetServerInfo

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
    s := formance.New()

    ctx := context.Background()
    res, err := s.Payments.PaymentsgetServerInfo(ctx)
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

**[*operations.PaymentsgetServerInfoResponse](../../models/operations/paymentsgetserverinforesponse.md), error**


## PaymentslistAccounts

List accounts

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
    s := formance.New()

    ctx := context.Background()
    res, err := s.Payments.PaymentslistAccounts(ctx, operations.PaymentslistAccountsRequest{
        Cursor: formance.String("aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ=="),
        PageSize: formance.Int64(354047),
        Sort: []string{
            "quos",
            "sint",
            "accusantium",
        },
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.AccountsCursor != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                                        | Type                                                                                             | Required                                                                                         | Description                                                                                      |
| ------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------ |
| `ctx`                                                                                            | [context.Context](https://pkg.go.dev/context#Context)                                            | :heavy_check_mark:                                                                               | The context to use for the request.                                                              |
| `request`                                                                                        | [operations.PaymentslistAccountsRequest](../../models/operations/paymentslistaccountsrequest.md) | :heavy_check_mark:                                                                               | The request object to use for the request.                                                       |


### Response

**[*operations.PaymentslistAccountsResponse](../../models/operations/paymentslistaccountsresponse.md), error**


## ReadConnectorConfig

Read connector config

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
    s := formance.New()

    ctx := context.Background()
    res, err := s.Payments.ReadConnectorConfig(ctx, operations.ReadConnectorConfigRequest{
        Connector: shared.ConnectorBankingCircle,
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.ConnectorConfigResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                                      | Type                                                                                           | Required                                                                                       | Description                                                                                    |
| ---------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------- |
| `ctx`                                                                                          | [context.Context](https://pkg.go.dev/context#Context)                                          | :heavy_check_mark:                                                                             | The context to use for the request.                                                            |
| `request`                                                                                      | [operations.ReadConnectorConfigRequest](../../models/operations/readconnectorconfigrequest.md) | :heavy_check_mark:                                                                             | The request object to use for the request.                                                     |


### Response

**[*operations.ReadConnectorConfigResponse](../../models/operations/readconnectorconfigresponse.md), error**


## ResetConnector

Reset a connector by its name.
It will remove the connector and ALL PAYMENTS generated with it.


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
    s := formance.New()

    ctx := context.Background()
    res, err := s.Payments.ResetConnector(ctx, operations.ResetConnectorRequest{
        Connector: shared.ConnectorMoneycorp,
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

| Parameter                                                                            | Type                                                                                 | Required                                                                             | Description                                                                          |
| ------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------ |
| `ctx`                                                                                | [context.Context](https://pkg.go.dev/context#Context)                                | :heavy_check_mark:                                                                   | The context to use for the request.                                                  |
| `request`                                                                            | [operations.ResetConnectorRequest](../../models/operations/resetconnectorrequest.md) | :heavy_check_mark:                                                                   | The request object to use for the request.                                           |


### Response

**[*operations.ResetConnectorResponse](../../models/operations/resetconnectorresponse.md), error**


## RetryTransferInitiation

Retry a failed transfer initiation

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
    s := formance.New()

    ctx := context.Background()
    res, err := s.Payments.RetryTransferInitiation(ctx, operations.RetryTransferInitiationRequest{
        TransferID: "mollitia",
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

| Parameter                                                                                              | Type                                                                                                   | Required                                                                                               | Description                                                                                            |
| ------------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------------ |
| `ctx`                                                                                                  | [context.Context](https://pkg.go.dev/context#Context)                                                  | :heavy_check_mark:                                                                                     | The context to use for the request.                                                                    |
| `request`                                                                                              | [operations.RetryTransferInitiationRequest](../../models/operations/retrytransferinitiationrequest.md) | :heavy_check_mark:                                                                                     | The request object to use for the request.                                                             |


### Response

**[*operations.RetryTransferInitiationResponse](../../models/operations/retrytransferinitiationresponse.md), error**


## UdpateTransferInitiationStatus

Update a transfer initiation status

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
    s := formance.New()

    ctx := context.Background()
    res, err := s.Payments.UdpateTransferInitiationStatus(ctx, operations.UdpateTransferInitiationStatusRequest{
        UpdateTransferInitiationStatusRequest: shared.UpdateTransferInitiationStatusRequest{
            Status: shared.UpdateTransferInitiationStatusRequestStatusProcessing,
        },
        TransferID: "eum",
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

| Parameter                                                                                                            | Type                                                                                                                 | Required                                                                                                             | Description                                                                                                          |
| -------------------------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------- |
| `ctx`                                                                                                                | [context.Context](https://pkg.go.dev/context#Context)                                                                | :heavy_check_mark:                                                                                                   | The context to use for the request.                                                                                  |
| `request`                                                                                                            | [operations.UdpateTransferInitiationStatusRequest](../../models/operations/udpatetransferinitiationstatusrequest.md) | :heavy_check_mark:                                                                                                   | The request object to use for the request.                                                                           |


### Response

**[*operations.UdpateTransferInitiationStatusResponse](../../models/operations/udpatetransferinitiationstatusresponse.md), error**


## UninstallConnector

Uninstall a connector by its name.

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
    s := formance.New()

    ctx := context.Background()
    res, err := s.Payments.UninstallConnector(ctx, operations.UninstallConnectorRequest{
        Connector: shared.ConnectorDummyPay,
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

| Parameter                                                                                    | Type                                                                                         | Required                                                                                     | Description                                                                                  |
| -------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------- |
| `ctx`                                                                                        | [context.Context](https://pkg.go.dev/context#Context)                                        | :heavy_check_mark:                                                                           | The context to use for the request.                                                          |
| `request`                                                                                    | [operations.UninstallConnectorRequest](../../models/operations/uninstallconnectorrequest.md) | :heavy_check_mark:                                                                           | The request object to use for the request.                                                   |


### Response

**[*operations.UninstallConnectorResponse](../../models/operations/uninstallconnectorresponse.md), error**


## UpdateMetadata

Update metadata

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
    s := formance.New()

    ctx := context.Background()
    res, err := s.Payments.UpdateMetadata(ctx, operations.UpdateMetadataRequest{
        PaymentMetadata: shared.PaymentMetadata{
            Key: formance.String("necessitatibus"),
        },
        PaymentID: "odit",
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

| Parameter                                                                            | Type                                                                                 | Required                                                                             | Description                                                                          |
| ------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------ |
| `ctx`                                                                                | [context.Context](https://pkg.go.dev/context#Context)                                | :heavy_check_mark:                                                                   | The context to use for the request.                                                  |
| `request`                                                                            | [operations.UpdateMetadataRequest](../../models/operations/updatemetadatarequest.md) | :heavy_check_mark:                                                                   | The request object to use for the request.                                           |


### Response

**[*operations.UpdateMetadataResponse](../../models/operations/updatemetadataresponse.md), error**

