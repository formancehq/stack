# Payments

### Available Operations

* [ConnectorsTransfer](#connectorstransfer) - Transfer funds between Connector accounts
* [CreateBankAccount](#createbankaccount) - Create a BankAccount in Payments and on the PSP
* [CreateTransferInitiation](#createtransferinitiation) - Create a TransferInitiation
* [DeleteTransferInitiation](#deletetransferinitiation) - Delete a transfer initiation
* [GetAccountBalances](#getaccountbalances) - Get account balances
* [GetBankAccount](#getbankaccount) - Get a bank account created by user on Formance
* [~~GetConnectorTask~~](#getconnectortask) - Read a specific task of the connector :warning: **Deprecated**
* [GetConnectorTaskV1](#getconnectortaskv1) - Read a specific task of the connector
* [GetPayment](#getpayment) - Get a payment
* [GetTransferInitiation](#gettransferinitiation) - Get a transfer initiation
* [InstallConnector](#installconnector) - Install a connector
* [ListAllConnectors](#listallconnectors) - List all installed connectors
* [ListBankAccounts](#listbankaccounts) - List bank accounts created by user on Formance
* [ListConfigsAvailableConnectors](#listconfigsavailableconnectors) - List the configs of each available connector
* [~~ListConnectorTasks~~](#listconnectortasks) - List tasks from a connector :warning: **Deprecated**
* [ListConnectorTasksV1](#listconnectortasksv1) - List tasks from a connector
* [ListPayments](#listpayments) - List payments
* [ListTransferInitiations](#listtransferinitiations) - List Transfer Initiations
* [PaymentsgetAccount](#paymentsgetaccount) - Get an account
* [PaymentsgetServerInfo](#paymentsgetserverinfo) - Get server info
* [PaymentslistAccounts](#paymentslistaccounts) - List accounts
* [~~ReadConnectorConfig~~](#readconnectorconfig) - Read the config of a connector :warning: **Deprecated**
* [ReadConnectorConfigV1](#readconnectorconfigv1) - Read the config of a connector
* [~~ResetConnector~~](#resetconnector) - Reset a connector :warning: **Deprecated**
* [ResetConnectorV1](#resetconnectorv1) - Reset a connector
* [RetryTransferInitiation](#retrytransferinitiation) - Retry a failed transfer initiation
* [UdpateTransferInitiationStatus](#udpatetransferinitiationstatus) - Update the status of a transfer initiation
* [~~UninstallConnector~~](#uninstallconnector) - Uninstall a connector :warning: **Deprecated**
* [UninstallConnectorV1](#uninstallconnectorv1) - Uninstall a connector
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
    s := formance.New(
        formance.WithSecurity(shared.Security{
            Authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
        }),
    )

    ctx := context.Background()
    res, err := s.Payments.ConnectorsTransfer(ctx, operations.ConnectorsTransferRequest{
        TransferRequest: shared.TransferRequest{
            Amount: big.NewInt(100),
            Asset: "USD",
            Destination: "acct_1Gqj58KZcSIg2N2q",
            Source: formance.String("acct_1Gqj58KZcSIg2N2q"),
        },
        Connector: shared.ConnectorWise,
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.TransferResponse != nil {
        // handle response
    }
}
```

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
    s := formance.New(
        formance.WithSecurity(shared.Security{
            Authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
        }),
    )

    ctx := context.Background()
    res, err := s.Payments.CreateBankAccount(ctx, shared.BankAccountRequest{
        AccountNumber: formance.String("consequatur"),
        ConnectorID: "est",
        Country: "GB",
        Iban: formance.String("quibusdam"),
        Name: "My account",
        SwiftBicCode: formance.String("explicabo"),
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.BankAccountResponse != nil {
        // handle response
    }
}
```

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
    s := formance.New(
        formance.WithSecurity(shared.Security{
            Authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
        }),
    )

    ctx := context.Background()
    res, err := s.Payments.CreateTransferInitiation(ctx, shared.TransferInitiationRequest{
        Amount: big.NewInt(647174),
        Asset: "USD",
        ConnectorID: formance.String("distinctio"),
        Description: "quibusdam",
        DestinationAccountID: "labore",
        Provider: shared.ConnectorWise.ToPointer(),
        Reference: "XXX",
        ScheduledAt: types.MustTimeFromString("2022-08-08T19:05:24.174Z"),
        SourceAccountID: "cupiditate",
        Type: shared.TransferInitiationRequestTypePayout,
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
    s := formance.New(
        formance.WithSecurity(shared.Security{
            Authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
        }),
    )

    ctx := context.Background()
    res, err := s.Payments.DeleteTransferInitiation(ctx, operations.DeleteTransferInitiationRequest{
        TransferID: "perferendis",
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.StatusCode == http.StatusOK {
        // handle response
    }
}
```

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
    s := formance.New(
        formance.WithSecurity(shared.Security{
            Authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
        }),
    )

    ctx := context.Background()
    res, err := s.Payments.GetAccountBalances(ctx, operations.GetAccountBalancesRequest{
        AccountID: "magni",
        Asset: formance.String("assumenda"),
        Cursor: formance.String("aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ=="),
        From: types.MustTimeFromString("2022-12-30T06:52:02.282Z"),
        Limit: formance.Int64(146441),
        PageSize: formance.Int64(677817),
        Sort: []string{
            "tempora",
            "facilis",
            "tempore",
        },
        To: types.MustTimeFromString("2022-01-14T19:13:42.009Z"),
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.BalancesCursor != nil {
        // handle response
    }
}
```

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
    s := formance.New(
        formance.WithSecurity(shared.Security{
            Authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
        }),
    )

    ctx := context.Background()
    res, err := s.Payments.GetBankAccount(ctx, operations.GetBankAccountRequest{
        BankAccountID: "eum",
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.BankAccountResponse != nil {
        // handle response
    }
}
```

## ~~GetConnectorTask~~

Get a specific task associated to the connector.

> :warning: **DEPRECATED**: this method will be removed in a future release, please migrate away from it as soon as possible.

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
    res, err := s.Payments.GetConnectorTask(ctx, operations.GetConnectorTaskRequest{
        Connector: shared.ConnectorDummyPay,
        TaskID: "eligendi",
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.TaskResponse != nil {
        // handle response
    }
}
```

## GetConnectorTaskV1

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
    s := formance.New(
        formance.WithSecurity(shared.Security{
            Authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
        }),
    )

    ctx := context.Background()
    res, err := s.Payments.GetConnectorTaskV1(ctx, operations.GetConnectorTaskV1Request{
        Connector: shared.ConnectorCurrencyCloud,
        ConnectorID: "aliquid",
        TaskID: "provident",
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.TaskResponse != nil {
        // handle response
    }
}
```

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
    s := formance.New(
        formance.WithSecurity(shared.Security{
            Authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
        }),
    )

    ctx := context.Background()
    res, err := s.Payments.GetPayment(ctx, operations.GetPaymentRequest{
        PaymentID: "necessitatibus",
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.PaymentResponse != nil {
        // handle response
    }
}
```

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
    s := formance.New(
        formance.WithSecurity(shared.Security{
            Authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
        }),
    )

    ctx := context.Background()
    res, err := s.Payments.GetTransferInitiation(ctx, operations.GetTransferInitiationRequest{
        TransferID: "sint",
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.TransferInitiationResponse != nil {
        // handle response
    }
}
```

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
    s := formance.New(
        formance.WithSecurity(shared.Security{
            Authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
        }),
    )

    ctx := context.Background()
    res, err := s.Payments.InstallConnector(ctx, operations.InstallConnectorRequest{
        RequestBody: shared.BankingCircleConfig{
            AuthorizationEndpoint: "XXX",
            Endpoint: "XXX",
            Name: "My Banking Circle Account",
            Password: "XXX",
            PollingPeriod: formance.String("60s"),
            UserCertificate: "XXX",
            UserCertificateKey: "XXX",
            Username: "XXX",
        },
        Connector: shared.ConnectorDummyPay,
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.ConnectorResponse != nil {
        // handle response
    }
}
```

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
    s := formance.New(
        formance.WithSecurity(shared.Security{
            Authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
        }),
    )

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
    s := formance.New(
        formance.WithSecurity(shared.Security{
            Authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
        }),
    )

    ctx := context.Background()
    res, err := s.Payments.ListBankAccounts(ctx, operations.ListBankAccountsRequest{
        Cursor: formance.String("aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ=="),
        PageSize: formance.Int64(891555),
        Sort: []string{
            "dolorum",
            "in",
            "in",
            "illum",
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
    s := formance.New(
        formance.WithSecurity(shared.Security{
            Authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
        }),
    )

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

## ~~ListConnectorTasks~~

List all tasks associated with this connector.

> :warning: **DEPRECATED**: this method will be removed in a future release, please migrate away from it as soon as possible.

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
    res, err := s.Payments.ListConnectorTasks(ctx, operations.ListConnectorTasksRequest{
        Connector: shared.ConnectorMoneycorp,
        Cursor: formance.String("aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ=="),
        PageSize: formance.Int64(699479),
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.TasksCursor != nil {
        // handle response
    }
}
```

## ListConnectorTasksV1

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
    s := formance.New(
        formance.WithSecurity(shared.Security{
            Authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
        }),
    )

    ctx := context.Background()
    res, err := s.Payments.ListConnectorTasksV1(ctx, operations.ListConnectorTasksV1Request{
        Connector: shared.ConnectorStripe,
        ConnectorID: "magnam",
        Cursor: formance.String("aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ=="),
        PageSize: formance.Int64(767024),
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.TasksCursor != nil {
        // handle response
    }
}
```

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
    s := formance.New(
        formance.WithSecurity(shared.Security{
            Authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
        }),
    )

    ctx := context.Background()
    res, err := s.Payments.ListPayments(ctx, operations.ListPaymentsRequest{
        Cursor: formance.String("aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ=="),
        PageSize: formance.Int64(813798),
        Sort: []string{
            "aliquid",
            "laborum",
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
    s := formance.New(
        formance.WithSecurity(shared.Security{
            Authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
        }),
    )

    ctx := context.Background()
    res, err := s.Payments.ListTransferInitiations(ctx, operations.ListTransferInitiationsRequest{
        Cursor: formance.String("aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ=="),
        PageSize: formance.Int64(881104),
        Query: formance.String("non"),
        Sort: []string{
            "enim",
            "accusamus",
            "delectus",
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
    s := formance.New(
        formance.WithSecurity(shared.Security{
            Authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
        }),
    )

    ctx := context.Background()
    res, err := s.Payments.PaymentsgetAccount(ctx, operations.PaymentsgetAccountRequest{
        AccountID: "quidem",
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.PaymentsAccountResponse != nil {
        // handle response
    }
}
```

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
    s := formance.New(
        formance.WithSecurity(shared.Security{
            Authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
        }),
    )

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
    s := formance.New(
        formance.WithSecurity(shared.Security{
            Authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
        }),
    )

    ctx := context.Background()
    res, err := s.Payments.PaymentslistAccounts(ctx, operations.PaymentslistAccountsRequest{
        Cursor: formance.String("aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ=="),
        PageSize: formance.Int64(588465),
        Sort: []string{
            "id",
            "blanditiis",
            "deleniti",
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

## ~~ReadConnectorConfig~~

Read connector config

> :warning: **DEPRECATED**: this method will be removed in a future release, please migrate away from it as soon as possible.

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
    res, err := s.Payments.ReadConnectorConfig(ctx, operations.ReadConnectorConfigRequest{
        Connector: shared.ConnectorMoneycorp,
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.ConnectorConfigResponse != nil {
        // handle response
    }
}
```

## ReadConnectorConfigV1

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
    s := formance.New(
        formance.WithSecurity(shared.Security{
            Authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
        }),
    )

    ctx := context.Background()
    res, err := s.Payments.ReadConnectorConfigV1(ctx, operations.ReadConnectorConfigV1Request{
        Connector: shared.ConnectorDummyPay,
        ConnectorID: "deserunt",
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.ConnectorConfigResponse != nil {
        // handle response
    }
}
```

## ~~ResetConnector~~

Reset a connector by its name.
It will remove the connector and ALL PAYMENTS generated with it.


> :warning: **DEPRECATED**: this method will be removed in a future release, please migrate away from it as soon as possible.

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
    res, err := s.Payments.ResetConnector(ctx, operations.ResetConnectorRequest{
        Connector: shared.ConnectorModulr,
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.StatusCode == http.StatusOK {
        // handle response
    }
}
```

## ResetConnectorV1

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
    s := formance.New(
        formance.WithSecurity(shared.Security{
            Authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
        }),
    )

    ctx := context.Background()
    res, err := s.Payments.ResetConnectorV1(ctx, operations.ResetConnectorV1Request{
        Connector: shared.ConnectorModulr,
        ConnectorID: "natus",
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.StatusCode == http.StatusOK {
        // handle response
    }
}
```

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
    s := formance.New(
        formance.WithSecurity(shared.Security{
            Authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
        }),
    )

    ctx := context.Background()
    res, err := s.Payments.RetryTransferInitiation(ctx, operations.RetryTransferInitiationRequest{
        TransferID: "omnis",
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.StatusCode == http.StatusOK {
        // handle response
    }
}
```

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
    s := formance.New(
        formance.WithSecurity(shared.Security{
            Authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
        }),
    )

    ctx := context.Background()
    res, err := s.Payments.UdpateTransferInitiationStatus(ctx, operations.UdpateTransferInitiationStatusRequest{
        UpdateTransferInitiationStatusRequest: shared.UpdateTransferInitiationStatusRequest{
            Status: shared.UpdateTransferInitiationStatusRequestStatusProcessed,
        },
        TransferID: "perferendis",
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.StatusCode == http.StatusOK {
        // handle response
    }
}
```

## ~~UninstallConnector~~

Uninstall a connector by its name.

> :warning: **DEPRECATED**: this method will be removed in a future release, please migrate away from it as soon as possible.

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
    res, err := s.Payments.UninstallConnector(ctx, operations.UninstallConnectorRequest{
        Connector: shared.ConnectorModulr,
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.StatusCode == http.StatusOK {
        // handle response
    }
}
```

## UninstallConnectorV1

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
    s := formance.New(
        formance.WithSecurity(shared.Security{
            Authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
        }),
    )

    ctx := context.Background()
    res, err := s.Payments.UninstallConnectorV1(ctx, operations.UninstallConnectorV1Request{
        Connector: shared.ConnectorWise,
        ConnectorID: "distinctio",
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.StatusCode == http.StatusOK {
        // handle response
    }
}
```

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
    s := formance.New(
        formance.WithSecurity(shared.Security{
            Authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
        }),
    )

    ctx := context.Background()
    res, err := s.Payments.UpdateMetadata(ctx, operations.UpdateMetadataRequest{
        PaymentMetadata: shared.PaymentMetadata{
            Key: formance.String("id"),
        },
        PaymentID: "labore",
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.StatusCode == http.StatusOK {
        // handle response
    }
}
```
