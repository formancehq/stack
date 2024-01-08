# github.com/formancehq/formance-sdk-go

<div align="left">
    <a href="https://speakeasyapi.dev/"><img src="https://custom-icon-badges.demolab.com/badge/-Built%20By%20Speakeasy-212015?style=for-the-badge&logoColor=FBE331&logo=speakeasy&labelColor=545454" /></a>
    <a href="https://opensource.org/licenses/MIT">
        <img src="https://img.shields.io/badge/License-MIT-blue.svg" style="width: 100px; height: 28px;" />
    </a>
</div>


## 🏗 **Welcome to your new SDK!** 🏗

It has been generated successfully based on your OpenAPI spec. However, it is not yet ready for production use. Here are some next steps:
- [ ] 🛠 Make your SDK feel handcrafted by [customizing it](https://www.speakeasyapi.dev/docs/customize-sdks)
- [ ] ♻️ Refine your SDK quickly by iterating locally with the [Speakeasy CLI](https://github.com/speakeasy-api/speakeasy)
- [ ] 🎁 Publish your SDK to package managers by [configuring automatic publishing](https://www.speakeasyapi.dev/docs/productionize-sdks/publish-sdks)
- [ ] ✨ When ready to productionize, delete this section from the README

<!-- Start SDK Installation [installation] -->
## SDK Installation

```bash
go get github.com/formancehq/formance-sdk-go
```
<!-- End SDK Installation [installation] -->

<!-- Start SDK Example Usage [usage] -->
## SDK Example Usage

### Example

```go
package main

import (
	"context"
	formancesdkgo "github.com/formancehq/formance-sdk-go"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"log"
)

func main() {
	s := formancesdkgo.New(
		formancesdkgo.WithSecurity(shared.Security{
			Authorization: "Bearer <YOUR_ACCESS_TOKEN_HERE>",
		}),
	)

	ctx := context.Background()
	res, err := s.GetVersions(ctx)
	if err != nil {
		log.Fatal(err)
	}

	if res.GetVersionsResponse != nil {
		// handle response
	}
}

```
<!-- End SDK Example Usage [usage] -->

<!-- Start Available Resources and Operations [operations] -->
## Available Resources and Operations

### [Formance SDK](docs/sdks/formance/README.md)

* [GetVersions](docs/sdks/formance/README.md#getversions) - Show stack version information
* [GetAPIAuthWellKnownOpenidConfiguration](docs/sdks/formance/README.md#getapiauthwellknownopenidconfiguration)

### [Auth](docs/sdks/auth/README.md)

* [CreateClient](docs/sdks/auth/README.md#createclient) - Create client
* [CreateSecret](docs/sdks/auth/README.md#createsecret) - Add a secret to a client
* [DeleteClient](docs/sdks/auth/README.md#deleteclient) - Delete client
* [DeleteSecret](docs/sdks/auth/README.md#deletesecret) - Delete a secret from a client
* [GetServerInfo](docs/sdks/auth/README.md#getserverinfo) - Get server info
* [ListClients](docs/sdks/auth/README.md#listclients) - List clients
* [ListUsers](docs/sdks/auth/README.md#listusers) - List users
* [ReadClient](docs/sdks/auth/README.md#readclient) - Read client
* [ReadUser](docs/sdks/auth/README.md#readuser) - Read user
* [UpdateClient](docs/sdks/auth/README.md#updateclient) - Update client

### [Ledger](docs/sdks/ledger/README.md)

* [CreateTransactions](docs/sdks/ledger/README.md#createtransactions) - Create a new batch of transactions to a ledger
* [AddMetadataOnTransaction](docs/sdks/ledger/README.md#addmetadataontransaction) - Set the metadata of a transaction by its ID
* [AddMetadataToAccount](docs/sdks/ledger/README.md#addmetadatatoaccount) - Add metadata to an account
* [CountAccounts](docs/sdks/ledger/README.md#countaccounts) - Count the accounts from a ledger
* [CountTransactions](docs/sdks/ledger/README.md#counttransactions) - Count the transactions from a ledger
* [CreateTransaction](docs/sdks/ledger/README.md#createtransaction) - Create a new transaction to a ledger
* [GetAccount](docs/sdks/ledger/README.md#getaccount) - Get account by its address
* [GetBalances](docs/sdks/ledger/README.md#getbalances) - Get the balances from a ledger's account
* [GetBalancesAggregated](docs/sdks/ledger/README.md#getbalancesaggregated) - Get the aggregated balances from selected accounts
* [GetInfo](docs/sdks/ledger/README.md#getinfo) - Show server information
* [GetLedgerInfo](docs/sdks/ledger/README.md#getledgerinfo) - Get information about a ledger
* [GetMapping](docs/sdks/ledger/README.md#getmapping) - Get the mapping of a ledger
* [GetTransaction](docs/sdks/ledger/README.md#gettransaction) - Get transaction from a ledger by its ID
* [ListAccounts](docs/sdks/ledger/README.md#listaccounts) - List accounts from a ledger
* [ListLogs](docs/sdks/ledger/README.md#listlogs) - List the logs from a ledger
* [ListTransactions](docs/sdks/ledger/README.md#listtransactions) - List transactions from a ledger
* [ReadStats](docs/sdks/ledger/README.md#readstats) - Get statistics from a ledger
* [RevertTransaction](docs/sdks/ledger/README.md#reverttransaction) - Revert a ledger transaction by its ID
* [~~RunScript~~](docs/sdks/ledger/README.md#runscript) - Execute a Numscript :warning: **Deprecated**
* [UpdateMapping](docs/sdks/ledger/README.md#updatemapping) - Update the mapping of a ledger
* [V2AddMetadataOnTransaction](docs/sdks/ledger/README.md#v2addmetadataontransaction) - Set the metadata of a transaction by its ID
* [V2AddMetadataToAccount](docs/sdks/ledger/README.md#v2addmetadatatoaccount) - Add metadata to an account
* [V2CountAccounts](docs/sdks/ledger/README.md#v2countaccounts) - Count the accounts from a ledger
* [V2CountTransactions](docs/sdks/ledger/README.md#v2counttransactions) - Count the transactions from a ledger
* [V2CreateBulk](docs/sdks/ledger/README.md#v2createbulk) - Bulk request
* [V2CreateLedger](docs/sdks/ledger/README.md#v2createledger) - Create a ledger
* [V2CreateTransaction](docs/sdks/ledger/README.md#v2createtransaction) - Create a new transaction to a ledger
* [V2DeleteAccountMetadata](docs/sdks/ledger/README.md#v2deleteaccountmetadata) - Delete metadata by key
* [V2DeleteTransactionMetadata](docs/sdks/ledger/README.md#v2deletetransactionmetadata) - Delete metadata by key
* [V2GetAccount](docs/sdks/ledger/README.md#v2getaccount) - Get account by its address
* [V2GetBalancesAggregated](docs/sdks/ledger/README.md#v2getbalancesaggregated) - Get the aggregated balances from selected accounts
* [V2GetInfo](docs/sdks/ledger/README.md#v2getinfo) - Show server information
* [V2GetLedger](docs/sdks/ledger/README.md#v2getledger) - Get a ledger
* [V2GetLedgerInfo](docs/sdks/ledger/README.md#v2getledgerinfo) - Get information about a ledger
* [V2GetTransaction](docs/sdks/ledger/README.md#v2gettransaction) - Get transaction from a ledger by its ID
* [V2ListAccounts](docs/sdks/ledger/README.md#v2listaccounts) - List accounts from a ledger
* [V2ListLedgers](docs/sdks/ledger/README.md#v2listledgers) - List ledgers
* [V2ListLogs](docs/sdks/ledger/README.md#v2listlogs) - List the logs from a ledger
* [V2ListTransactions](docs/sdks/ledger/README.md#v2listtransactions) - List transactions from a ledger
* [V2ReadStats](docs/sdks/ledger/README.md#v2readstats) - Get statistics from a ledger
* [V2RevertTransaction](docs/sdks/ledger/README.md#v2reverttransaction) - Revert a ledger transaction by its ID

### [Orchestration](docs/sdks/orchestration/README.md)

* [CancelEvent](docs/sdks/orchestration/README.md#cancelevent) - Cancel a running workflow
* [CreateTrigger](docs/sdks/orchestration/README.md#createtrigger) - Create trigger
* [CreateWorkflow](docs/sdks/orchestration/README.md#createworkflow) - Create workflow
* [DeleteTrigger](docs/sdks/orchestration/README.md#deletetrigger) - Delete trigger
* [DeleteWorkflow](docs/sdks/orchestration/README.md#deleteworkflow) - Delete a flow by id
* [GetInstance](docs/sdks/orchestration/README.md#getinstance) - Get a workflow instance by id
* [GetInstanceHistory](docs/sdks/orchestration/README.md#getinstancehistory) - Get a workflow instance history by id
* [GetInstanceStageHistory](docs/sdks/orchestration/README.md#getinstancestagehistory) - Get a workflow instance stage history
* [GetWorkflow](docs/sdks/orchestration/README.md#getworkflow) - Get a flow by id
* [ListInstances](docs/sdks/orchestration/README.md#listinstances) - List instances of a workflow
* [ListTriggers](docs/sdks/orchestration/README.md#listtriggers) - List triggers
* [ListTriggersOccurrences](docs/sdks/orchestration/README.md#listtriggersoccurrences) - List triggers occurrences
* [ListWorkflows](docs/sdks/orchestration/README.md#listworkflows) - List registered workflows
* [OrchestrationgetServerInfo](docs/sdks/orchestration/README.md#orchestrationgetserverinfo) - Get server info
* [ReadTrigger](docs/sdks/orchestration/README.md#readtrigger) - Read trigger
* [RunWorkflow](docs/sdks/orchestration/README.md#runworkflow) - Run workflow
* [SendEvent](docs/sdks/orchestration/README.md#sendevent) - Send an event to a running workflow
* [TestTrigger](docs/sdks/orchestration/README.md#testtrigger) - Test trigger
* [V2CancelEvent](docs/sdks/orchestration/README.md#v2cancelevent) - Cancel a running workflow
* [V2CreateTrigger](docs/sdks/orchestration/README.md#v2createtrigger) - Create trigger
* [V2CreateWorkflow](docs/sdks/orchestration/README.md#v2createworkflow) - Create workflow
* [V2DeleteTrigger](docs/sdks/orchestration/README.md#v2deletetrigger) - Delete trigger
* [V2DeleteWorkflow](docs/sdks/orchestration/README.md#v2deleteworkflow) - Delete a flow by id
* [V2GetInstance](docs/sdks/orchestration/README.md#v2getinstance) - Get a workflow instance by id
* [V2GetInstanceHistory](docs/sdks/orchestration/README.md#v2getinstancehistory) - Get a workflow instance history by id
* [V2GetInstanceStageHistory](docs/sdks/orchestration/README.md#v2getinstancestagehistory) - Get a workflow instance stage history
* [V2GetServerInfo](docs/sdks/orchestration/README.md#v2getserverinfo) - Get server info
* [V2GetWorkflow](docs/sdks/orchestration/README.md#v2getworkflow) - Get a flow by id
* [V2ListInstances](docs/sdks/orchestration/README.md#v2listinstances) - List instances of a workflow
* [V2ListTriggers](docs/sdks/orchestration/README.md#v2listtriggers) - List triggers
* [V2ListTriggersOccurrences](docs/sdks/orchestration/README.md#v2listtriggersoccurrences) - List triggers occurrences
* [V2ListWorkflows](docs/sdks/orchestration/README.md#v2listworkflows) - List registered workflows
* [V2ReadTrigger](docs/sdks/orchestration/README.md#v2readtrigger) - Read trigger
* [V2RunWorkflow](docs/sdks/orchestration/README.md#v2runworkflow) - Run workflow
* [V2SendEvent](docs/sdks/orchestration/README.md#v2sendevent) - Send an event to a running workflow

### [Payments](docs/sdks/payments/README.md)

* [AddAccountToPool](docs/sdks/payments/README.md#addaccounttopool) - Add an account to a pool
* [ConnectorsTransfer](docs/sdks/payments/README.md#connectorstransfer) - Transfer funds between Connector accounts
* [CreateBankAccount](docs/sdks/payments/README.md#createbankaccount) - Create a BankAccount in Payments and on the PSP
* [CreatePayment](docs/sdks/payments/README.md#createpayment) - Create a payment
* [CreatePool](docs/sdks/payments/README.md#createpool) - Create a Pool
* [CreateTransferInitiation](docs/sdks/payments/README.md#createtransferinitiation) - Create a TransferInitiation
* [DeletePool](docs/sdks/payments/README.md#deletepool) - Delete a Pool
* [DeleteTransferInitiation](docs/sdks/payments/README.md#deletetransferinitiation) - Delete a transfer initiation
* [GetAccountBalances](docs/sdks/payments/README.md#getaccountbalances) - Get account balances
* [GetBankAccount](docs/sdks/payments/README.md#getbankaccount) - Get a bank account created by user on Formance
* [~~GetConnectorTask~~](docs/sdks/payments/README.md#getconnectortask) - Read a specific task of the connector :warning: **Deprecated**
* [GetConnectorTaskV1](docs/sdks/payments/README.md#getconnectortaskv1) - Read a specific task of the connector
* [GetPayment](docs/sdks/payments/README.md#getpayment) - Get a payment
* [GetPool](docs/sdks/payments/README.md#getpool) - Get a Pool
* [GetPoolBalances](docs/sdks/payments/README.md#getpoolbalances) - Get pool balances
* [GetTransferInitiation](docs/sdks/payments/README.md#gettransferinitiation) - Get a transfer initiation
* [InstallConnector](docs/sdks/payments/README.md#installconnector) - Install a connector
* [ListAllConnectors](docs/sdks/payments/README.md#listallconnectors) - List all installed connectors
* [ListBankAccounts](docs/sdks/payments/README.md#listbankaccounts) - List bank accounts created by user on Formance
* [ListConfigsAvailableConnectors](docs/sdks/payments/README.md#listconfigsavailableconnectors) - List the configs of each available connector
* [~~ListConnectorTasks~~](docs/sdks/payments/README.md#listconnectortasks) - List tasks from a connector :warning: **Deprecated**
* [ListConnectorTasksV1](docs/sdks/payments/README.md#listconnectortasksv1) - List tasks from a connector
* [ListPayments](docs/sdks/payments/README.md#listpayments) - List payments
* [ListPools](docs/sdks/payments/README.md#listpools) - List Pools
* [ListTransferInitiations](docs/sdks/payments/README.md#listtransferinitiations) - List Transfer Initiations
* [PaymentsgetAccount](docs/sdks/payments/README.md#paymentsgetaccount) - Get an account
* [PaymentsgetServerInfo](docs/sdks/payments/README.md#paymentsgetserverinfo) - Get server info
* [PaymentslistAccounts](docs/sdks/payments/README.md#paymentslistaccounts) - List accounts
* [~~ReadConnectorConfig~~](docs/sdks/payments/README.md#readconnectorconfig) - Read the config of a connector :warning: **Deprecated**
* [ReadConnectorConfigV1](docs/sdks/payments/README.md#readconnectorconfigv1) - Read the config of a connector
* [RemoveAccountFromPool](docs/sdks/payments/README.md#removeaccountfrompool) - Remove an account from a pool
* [~~ResetConnector~~](docs/sdks/payments/README.md#resetconnector) - Reset a connector :warning: **Deprecated**
* [ResetConnectorV1](docs/sdks/payments/README.md#resetconnectorv1) - Reset a connector
* [RetryTransferInitiation](docs/sdks/payments/README.md#retrytransferinitiation) - Retry a failed transfer initiation
* [UdpateTransferInitiationStatus](docs/sdks/payments/README.md#udpatetransferinitiationstatus) - Update the status of a transfer initiation
* [~~UninstallConnector~~](docs/sdks/payments/README.md#uninstallconnector) - Uninstall a connector :warning: **Deprecated**
* [UninstallConnectorV1](docs/sdks/payments/README.md#uninstallconnectorv1) - Uninstall a connector
* [UpdateConnectorConfigV1](docs/sdks/payments/README.md#updateconnectorconfigv1) - Update the config of a connector
* [UpdateMetadata](docs/sdks/payments/README.md#updatemetadata) - Update metadata

### [Reconciliation](docs/sdks/reconciliation/README.md)

* [CreatePolicy](docs/sdks/reconciliation/README.md#createpolicy) - Create a policy
* [DeletePolicy](docs/sdks/reconciliation/README.md#deletepolicy) - Delete a policy
* [GetPolicy](docs/sdks/reconciliation/README.md#getpolicy) - Get a policy
* [GetReconciliation](docs/sdks/reconciliation/README.md#getreconciliation) - Get a reconciliation
* [ListPolicies](docs/sdks/reconciliation/README.md#listpolicies) - List policies
* [ListReconciliations](docs/sdks/reconciliation/README.md#listreconciliations) - List reconciliations
* [Reconcile](docs/sdks/reconciliation/README.md#reconcile) - Reconcile using a policy
* [ReconciliationgetServerInfo](docs/sdks/reconciliation/README.md#reconciliationgetserverinfo) - Get server info

### [Search](docs/sdks/search/README.md)

* [Search](docs/sdks/search/README.md#search) - Search
* [SearchgetServerInfo](docs/sdks/search/README.md#searchgetserverinfo) - Get server info

### [Wallets](docs/sdks/wallets/README.md)

* [ConfirmHold](docs/sdks/wallets/README.md#confirmhold) - Confirm a hold
* [CreateBalance](docs/sdks/wallets/README.md#createbalance) - Create a balance
* [CreateWallet](docs/sdks/wallets/README.md#createwallet) - Create a new wallet
* [CreditWallet](docs/sdks/wallets/README.md#creditwallet) - Credit a wallet
* [DebitWallet](docs/sdks/wallets/README.md#debitwallet) - Debit a wallet
* [GetBalance](docs/sdks/wallets/README.md#getbalance) - Get detailed balance
* [GetHold](docs/sdks/wallets/README.md#gethold) - Get a hold
* [GetHolds](docs/sdks/wallets/README.md#getholds) - Get all holds for a wallet
* [GetTransactions](docs/sdks/wallets/README.md#gettransactions)
* [GetWallet](docs/sdks/wallets/README.md#getwallet) - Get a wallet
* [GetWalletSummary](docs/sdks/wallets/README.md#getwalletsummary) - Get wallet summary
* [ListBalances](docs/sdks/wallets/README.md#listbalances) - List balances of a wallet
* [ListWallets](docs/sdks/wallets/README.md#listwallets) - List all wallets
* [UpdateWallet](docs/sdks/wallets/README.md#updatewallet) - Update a wallet
* [VoidHold](docs/sdks/wallets/README.md#voidhold) - Cancel a hold
* [WalletsgetServerInfo](docs/sdks/wallets/README.md#walletsgetserverinfo) - Get server info

### [Webhooks](docs/sdks/webhooks/README.md)

* [ActivateConfig](docs/sdks/webhooks/README.md#activateconfig) - Activate one config
* [ChangeConfigSecret](docs/sdks/webhooks/README.md#changeconfigsecret) - Change the signing secret of a config
* [DeactivateConfig](docs/sdks/webhooks/README.md#deactivateconfig) - Deactivate one config
* [DeleteConfig](docs/sdks/webhooks/README.md#deleteconfig) - Delete one config
* [GetManyConfigs](docs/sdks/webhooks/README.md#getmanyconfigs) - Get many configs
* [InsertConfig](docs/sdks/webhooks/README.md#insertconfig) - Insert a new config
* [TestConfig](docs/sdks/webhooks/README.md#testconfig) - Test one config
<!-- End Available Resources and Operations [operations] -->

<!-- Start Error Handling [errors] -->
## Error Handling

Handling errors in this SDK should largely match your expectations.  All operations return a response object or an error, they will never return both.  When specified by the OpenAPI spec document, the SDK will return the appropriate subclass.

| Error Object            | Status Code             | Content Type            |
| ----------------------- | ----------------------- | ----------------------- |
| sdkerrors.ErrorResponse | 400,404                 | application/json        |
| sdkerrors.SDKError      | 400-600                 | */*                     |

### Example

```go
package main

import (
	"context"
	"errors"
	formancesdkgo "github.com/formancehq/formance-sdk-go"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/sdkerrors"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"log"
)

func main() {
	s := formancesdkgo.New(
		formancesdkgo.WithSecurity(shared.Security{
			Authorization: "Bearer <YOUR_ACCESS_TOKEN_HERE>",
		}),
	)

	ctx := context.Background()
	res, err := s.Ledger.AddMetadataToAccount(ctx, operations.AddMetadataToAccountRequest{
		RequestBody: map[string]interface{}{
			"key": "string",
		},
		Address: "users:001",
		Ledger:  "ledger001",
	})
	if err != nil {

		var e *sdkerrors.ErrorResponse
		if errors.As(err, &e) {
			// handle error
			log.Fatal(e.Error())
		}

		var e *sdkerrors.SDKError
		if errors.As(err, &e) {
			// handle error
			log.Fatal(e.Error())
		}
	}
}

```
<!-- End Error Handling [errors] -->

<!-- Start Server Selection [server] -->
## Server Selection

### Select Server by Index

You can override the default server globally using the `WithServerIndex` option when initializing the SDK client instance. The selected server will then be used as the default on the operations that use it. This table lists the indexes associated with the available servers:

| # | Server | Variables |
| - | ------ | --------- |
| 0 | `http://localhost` | None |

#### Example

```go
package main

import (
	"context"
	formancesdkgo "github.com/formancehq/formance-sdk-go"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"log"
)

func main() {
	s := formancesdkgo.New(
		formancesdkgo.WithServerIndex(0),
		formancesdkgo.WithSecurity(shared.Security{
			Authorization: "Bearer <YOUR_ACCESS_TOKEN_HERE>",
		}),
	)

	ctx := context.Background()
	res, err := s.GetVersions(ctx)
	if err != nil {
		log.Fatal(err)
	}

	if res.GetVersionsResponse != nil {
		// handle response
	}
}

```


### Override Server URL Per-Client

The default server can also be overridden globally using the `WithServerURL` option when initializing the SDK client instance. For example:
```go
package main

import (
	"context"
	formancesdkgo "github.com/formancehq/formance-sdk-go"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"log"
)

func main() {
	s := formancesdkgo.New(
		formancesdkgo.WithServerURL("http://localhost"),
		formancesdkgo.WithSecurity(shared.Security{
			Authorization: "Bearer <YOUR_ACCESS_TOKEN_HERE>",
		}),
	)

	ctx := context.Background()
	res, err := s.GetVersions(ctx)
	if err != nil {
		log.Fatal(err)
	}

	if res.GetVersionsResponse != nil {
		// handle response
	}
}

```
<!-- End Server Selection [server] -->

<!-- Start Custom HTTP Client [http-client] -->
## Custom HTTP Client

The Go SDK makes API calls that wrap an internal HTTP client. The requirements for the HTTP client are very simple. It must match this interface:

```go
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}
```

The built-in `net/http` client satisfies this interface and a default client based on the built-in is provided by default. To replace this default with a client of your own, you can implement this interface yourself or provide your own client configured as desired. Here's a simple example, which adds a client with a 30 second timeout.

```go
import (
	"net/http"
	"time"
	"github.com/myorg/your-go-sdk"
)

var (
	httpClient = &http.Client{Timeout: 30 * time.Second}
	sdkClient  = sdk.New(sdk.WithClient(httpClient))
)
```

This can be a convenient way to configure timeouts, cookies, proxies, custom headers, and other low-level configuration.
<!-- End Custom HTTP Client [http-client] -->

<!-- Start Authentication [security] -->
## Authentication

### Per-Client Security Schemes

This SDK supports the following security scheme globally:

| Name            | Type            | Scheme          |
| --------------- | --------------- | --------------- |
| `Authorization` | oauth2          | OAuth2 token    |

You can configure it using the `WithSecurity` option when initializing the SDK client instance. For example:
```go
package main

import (
	"context"
	formancesdkgo "github.com/formancehq/formance-sdk-go"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"log"
)

func main() {
	s := formancesdkgo.New(
		formancesdkgo.WithSecurity(shared.Security{
			Authorization: "Bearer <YOUR_ACCESS_TOKEN_HERE>",
		}),
	)

	ctx := context.Background()
	res, err := s.GetVersions(ctx)
	if err != nil {
		log.Fatal(err)
	}

	if res.GetVersionsResponse != nil {
		// handle response
	}
}

```
<!-- End Authentication [security] -->

<!-- Start Special Types [types] -->
## Special Types


<!-- End Special Types [types] -->

<!-- Placeholder for Future Speakeasy SDK Sections -->

# Development

## Maturity

This SDK is in beta, and there may be breaking changes between versions without a major version update. Therefore, we recommend pinning usage
to a specific package version. This way, you can install the same version each time without breaking changes unless you are intentionally
looking for the latest version.

## Contributions

While we value open-source contributions to this SDK, this library is generated programmatically.
Feel free to open a PR or a Github issue as a proof of concept and we'll do our best to include it in a future release!

### SDK Created by [Speakeasy](https://docs.speakeasyapi.dev/docs/using-speakeasy/client-sdks)
