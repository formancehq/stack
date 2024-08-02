# github.com/formancehq/formance-sdk-go/v2

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
go get github.com/formancehq/formance-sdk-go/v2
```
<!-- End SDK Installation [installation] -->

<!-- Start SDK Example Usage [usage] -->
## SDK Example Usage

### Example

```go
package main

import (
	"context"
	"github.com/formancehq/formance-sdk-go/v2"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"log"
	"os"
)

func main() {
	s := v2.New(
		v2.WithSecurity(shared.Security{
			Authorization: os.Getenv("AUTHORIZATION"),
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


### [Auth.V1](docs/sdks/v1/README.md)

* [CreateClient](docs/sdks/v1/README.md#createclient) - Create client
* [CreateSecret](docs/sdks/v1/README.md#createsecret) - Add a secret to a client
* [DeleteClient](docs/sdks/v1/README.md#deleteclient) - Delete client
* [DeleteSecret](docs/sdks/v1/README.md#deletesecret) - Delete a secret from a client
* [GetOIDCWellKnowns](docs/sdks/v1/README.md#getoidcwellknowns) - Retrieve OpenID connect well-knowns.
* [GetServerInfo](docs/sdks/v1/README.md#getserverinfo) - Get server info
* [ListClients](docs/sdks/v1/README.md#listclients) - List clients
* [ListUsers](docs/sdks/v1/README.md#listusers) - List users
* [ReadClient](docs/sdks/v1/README.md#readclient) - Read client
* [ReadUser](docs/sdks/v1/README.md#readuser) - Read user
* [UpdateClient](docs/sdks/v1/README.md#updateclient) - Update client


### [Ledger.V1](docs/sdks/formancev1/README.md)

* [CreateTransactions](docs/sdks/formancev1/README.md#createtransactions) - Create a new batch of transactions to a ledger
* [AddMetadataOnTransaction](docs/sdks/formancev1/README.md#addmetadataontransaction) - Set the metadata of a transaction by its ID
* [AddMetadataToAccount](docs/sdks/formancev1/README.md#addmetadatatoaccount) - Add metadata to an account
* [CountAccounts](docs/sdks/formancev1/README.md#countaccounts) - Count the accounts from a ledger
* [CountTransactions](docs/sdks/formancev1/README.md#counttransactions) - Count the transactions from a ledger
* [CreateTransaction](docs/sdks/formancev1/README.md#createtransaction) - Create a new transaction to a ledger
* [GetAccount](docs/sdks/formancev1/README.md#getaccount) - Get account by its address
* [GetBalances](docs/sdks/formancev1/README.md#getbalances) - Get the balances from a ledger's account
* [GetBalancesAggregated](docs/sdks/formancev1/README.md#getbalancesaggregated) - Get the aggregated balances from selected accounts
* [GetInfo](docs/sdks/formancev1/README.md#getinfo) - Show server information
* [GetLedgerInfo](docs/sdks/formancev1/README.md#getledgerinfo) - Get information about a ledger
* [GetMapping](docs/sdks/formancev1/README.md#getmapping) - Get the mapping of a ledger
* [GetTransaction](docs/sdks/formancev1/README.md#gettransaction) - Get transaction from a ledger by its ID
* [ListAccounts](docs/sdks/formancev1/README.md#listaccounts) - List accounts from a ledger
* [ListLogs](docs/sdks/formancev1/README.md#listlogs) - List the logs from a ledger
* [ListTransactions](docs/sdks/formancev1/README.md#listtransactions) - List transactions from a ledger
* [ReadStats](docs/sdks/formancev1/README.md#readstats) - Get statistics from a ledger
* [RevertTransaction](docs/sdks/formancev1/README.md#reverttransaction) - Revert a ledger transaction by its ID
* [~~RunScript~~](docs/sdks/formancev1/README.md#runscript) - Execute a Numscript :warning: **Deprecated**
* [UpdateMapping](docs/sdks/formancev1/README.md#updatemapping) - Update the mapping of a ledger

### [Ledger.V2](docs/sdks/v2/README.md)

* [AddMetadataOnTransaction](docs/sdks/v2/README.md#addmetadataontransaction) - Set the metadata of a transaction by its ID
* [AddMetadataToAccount](docs/sdks/v2/README.md#addmetadatatoaccount) - Add metadata to an account
* [CountAccounts](docs/sdks/v2/README.md#countaccounts) - Count the accounts from a ledger
* [CountTransactions](docs/sdks/v2/README.md#counttransactions) - Count the transactions from a ledger
* [CreateBulk](docs/sdks/v2/README.md#createbulk) - Bulk request
* [CreateLedger](docs/sdks/v2/README.md#createledger) - Create a ledger
* [CreateTransaction](docs/sdks/v2/README.md#createtransaction) - Create a new transaction to a ledger
* [DeleteAccountMetadata](docs/sdks/v2/README.md#deleteaccountmetadata) - Delete metadata by key
* [DeleteLedgerMetadata](docs/sdks/v2/README.md#deleteledgermetadata) - Delete ledger metadata by key
* [DeleteTransactionMetadata](docs/sdks/v2/README.md#deletetransactionmetadata) - Delete metadata by key
* [ExportLogs](docs/sdks/v2/README.md#exportlogs) - Export logs
* [GetAccount](docs/sdks/v2/README.md#getaccount) - Get account by its address
* [GetBalancesAggregated](docs/sdks/v2/README.md#getbalancesaggregated) - Get the aggregated balances from selected accounts
* [GetInfo](docs/sdks/v2/README.md#getinfo) - Show server information
* [GetLedger](docs/sdks/v2/README.md#getledger) - Get a ledger
* [GetLedgerInfo](docs/sdks/v2/README.md#getledgerinfo) - Get information about a ledger
* [GetTransaction](docs/sdks/v2/README.md#gettransaction) - Get transaction from a ledger by its ID
* [GetVolumesWithBalances](docs/sdks/v2/README.md#getvolumeswithbalances) - Get list of volumes with balances for (account/asset)
* [ImportLogs](docs/sdks/v2/README.md#importlogs)
* [ListAccounts](docs/sdks/v2/README.md#listaccounts) - List accounts from a ledger
* [ListLedgers](docs/sdks/v2/README.md#listledgers) - List ledgers
* [ListLogs](docs/sdks/v2/README.md#listlogs) - List the logs from a ledger
* [ListTransactions](docs/sdks/v2/README.md#listtransactions) - List transactions from a ledger
* [ReadStats](docs/sdks/v2/README.md#readstats) - Get statistics from a ledger
* [RevertTransaction](docs/sdks/v2/README.md#reverttransaction) - Revert a ledger transaction by its ID
* [UpdateLedgerMetadata](docs/sdks/v2/README.md#updateledgermetadata) - Update ledger metadata


### [Orchestration.V1](docs/sdks/formanceorchestrationv1/README.md)

* [CancelEvent](docs/sdks/formanceorchestrationv1/README.md#cancelevent) - Cancel a running workflow
* [CreateTrigger](docs/sdks/formanceorchestrationv1/README.md#createtrigger) - Create trigger
* [CreateWorkflow](docs/sdks/formanceorchestrationv1/README.md#createworkflow) - Create workflow
* [DeleteTrigger](docs/sdks/formanceorchestrationv1/README.md#deletetrigger) - Delete trigger
* [DeleteWorkflow](docs/sdks/formanceorchestrationv1/README.md#deleteworkflow) - Delete a flow by id
* [GetInstance](docs/sdks/formanceorchestrationv1/README.md#getinstance) - Get a workflow instance by id
* [GetInstanceHistory](docs/sdks/formanceorchestrationv1/README.md#getinstancehistory) - Get a workflow instance history by id
* [GetInstanceStageHistory](docs/sdks/formanceorchestrationv1/README.md#getinstancestagehistory) - Get a workflow instance stage history
* [GetWorkflow](docs/sdks/formanceorchestrationv1/README.md#getworkflow) - Get a flow by id
* [ListInstances](docs/sdks/formanceorchestrationv1/README.md#listinstances) - List instances of a workflow
* [ListTriggers](docs/sdks/formanceorchestrationv1/README.md#listtriggers) - List triggers
* [ListTriggersOccurrences](docs/sdks/formanceorchestrationv1/README.md#listtriggersoccurrences) - List triggers occurrences
* [ListWorkflows](docs/sdks/formanceorchestrationv1/README.md#listworkflows) - List registered workflows
* [OrchestrationgetServerInfo](docs/sdks/formanceorchestrationv1/README.md#orchestrationgetserverinfo) - Get server info
* [ReadTrigger](docs/sdks/formanceorchestrationv1/README.md#readtrigger) - Read trigger
* [RunWorkflow](docs/sdks/formanceorchestrationv1/README.md#runworkflow) - Run workflow
* [SendEvent](docs/sdks/formanceorchestrationv1/README.md#sendevent) - Send an event to a running workflow

### [Orchestration.V2](docs/sdks/formancev2/README.md)

* [CancelEvent](docs/sdks/formancev2/README.md#cancelevent) - Cancel a running workflow
* [CreateTrigger](docs/sdks/formancev2/README.md#createtrigger) - Create trigger
* [CreateWorkflow](docs/sdks/formancev2/README.md#createworkflow) - Create workflow
* [DeleteTrigger](docs/sdks/formancev2/README.md#deletetrigger) - Delete trigger
* [DeleteWorkflow](docs/sdks/formancev2/README.md#deleteworkflow) - Delete a flow by id
* [GetInstance](docs/sdks/formancev2/README.md#getinstance) - Get a workflow instance by id
* [GetInstanceHistory](docs/sdks/formancev2/README.md#getinstancehistory) - Get a workflow instance history by id
* [GetInstanceStageHistory](docs/sdks/formancev2/README.md#getinstancestagehistory) - Get a workflow instance stage history
* [GetServerInfo](docs/sdks/formancev2/README.md#getserverinfo) - Get server info
* [GetWorkflow](docs/sdks/formancev2/README.md#getworkflow) - Get a flow by id
* [ListInstances](docs/sdks/formancev2/README.md#listinstances) - List instances of a workflow
* [ListTriggers](docs/sdks/formancev2/README.md#listtriggers) - List triggers
* [ListTriggersOccurrences](docs/sdks/formancev2/README.md#listtriggersoccurrences) - List triggers occurrences
* [ListWorkflows](docs/sdks/formancev2/README.md#listworkflows) - List registered workflows
* [ReadTrigger](docs/sdks/formancev2/README.md#readtrigger) - Read trigger
* [RunWorkflow](docs/sdks/formancev2/README.md#runworkflow) - Run workflow
* [SendEvent](docs/sdks/formancev2/README.md#sendevent) - Send an event to a running workflow
* [TestTrigger](docs/sdks/formancev2/README.md#testtrigger) - Test trigger


### [Payments.V1](docs/sdks/formancepaymentsv1/README.md)

* [AddAccountToPool](docs/sdks/formancepaymentsv1/README.md#addaccounttopool) - Add an account to a pool
* [ConnectorsTransfer](docs/sdks/formancepaymentsv1/README.md#connectorstransfer) - Transfer funds between Connector accounts
* [CreateAccount](docs/sdks/formancepaymentsv1/README.md#createaccount) - Create an account
* [CreateBankAccount](docs/sdks/formancepaymentsv1/README.md#createbankaccount) - Create a BankAccount in Payments and on the PSP
* [CreatePayment](docs/sdks/formancepaymentsv1/README.md#createpayment) - Create a payment
* [CreatePool](docs/sdks/formancepaymentsv1/README.md#createpool) - Create a Pool
* [CreateTransferInitiation](docs/sdks/formancepaymentsv1/README.md#createtransferinitiation) - Create a TransferInitiation
* [DeletePool](docs/sdks/formancepaymentsv1/README.md#deletepool) - Delete a Pool
* [DeleteTransferInitiation](docs/sdks/formancepaymentsv1/README.md#deletetransferinitiation) - Delete a transfer initiation
* [ForwardBankAccount](docs/sdks/formancepaymentsv1/README.md#forwardbankaccount) - Forward a bank account to a connector
* [GetAccountBalances](docs/sdks/formancepaymentsv1/README.md#getaccountbalances) - Get account balances
* [GetBankAccount](docs/sdks/formancepaymentsv1/README.md#getbankaccount) - Get a bank account created by user on Formance
* [~~GetConnectorTask~~](docs/sdks/formancepaymentsv1/README.md#getconnectortask) - Read a specific task of the connector :warning: **Deprecated**
* [GetConnectorTaskV1](docs/sdks/formancepaymentsv1/README.md#getconnectortaskv1) - Read a specific task of the connector
* [GetPayment](docs/sdks/formancepaymentsv1/README.md#getpayment) - Get a payment
* [GetPool](docs/sdks/formancepaymentsv1/README.md#getpool) - Get a Pool
* [GetPoolBalances](docs/sdks/formancepaymentsv1/README.md#getpoolbalances) - Get pool balances
* [GetTransferInitiation](docs/sdks/formancepaymentsv1/README.md#gettransferinitiation) - Get a transfer initiation
* [InstallConnector](docs/sdks/formancepaymentsv1/README.md#installconnector) - Install a connector
* [ListAllConnectors](docs/sdks/formancepaymentsv1/README.md#listallconnectors) - List all installed connectors
* [ListBankAccounts](docs/sdks/formancepaymentsv1/README.md#listbankaccounts) - List bank accounts created by user on Formance
* [ListConfigsAvailableConnectors](docs/sdks/formancepaymentsv1/README.md#listconfigsavailableconnectors) - List the configs of each available connector
* [~~ListConnectorTasks~~](docs/sdks/formancepaymentsv1/README.md#listconnectortasks) - List tasks from a connector :warning: **Deprecated**
* [ListConnectorTasksV1](docs/sdks/formancepaymentsv1/README.md#listconnectortasksv1) - List tasks from a connector
* [ListPayments](docs/sdks/formancepaymentsv1/README.md#listpayments) - List payments
* [ListPools](docs/sdks/formancepaymentsv1/README.md#listpools) - List Pools
* [ListTransferInitiations](docs/sdks/formancepaymentsv1/README.md#listtransferinitiations) - List Transfer Initiations
* [PaymentsgetAccount](docs/sdks/formancepaymentsv1/README.md#paymentsgetaccount) - Get an account
* [PaymentsgetServerInfo](docs/sdks/formancepaymentsv1/README.md#paymentsgetserverinfo) - Get server info
* [PaymentslistAccounts](docs/sdks/formancepaymentsv1/README.md#paymentslistaccounts) - List accounts
* [~~ReadConnectorConfig~~](docs/sdks/formancepaymentsv1/README.md#readconnectorconfig) - Read the config of a connector :warning: **Deprecated**
* [ReadConnectorConfigV1](docs/sdks/formancepaymentsv1/README.md#readconnectorconfigv1) - Read the config of a connector
* [RemoveAccountFromPool](docs/sdks/formancepaymentsv1/README.md#removeaccountfrompool) - Remove an account from a pool
* [~~ResetConnector~~](docs/sdks/formancepaymentsv1/README.md#resetconnector) - Reset a connector :warning: **Deprecated**
* [ResetConnectorV1](docs/sdks/formancepaymentsv1/README.md#resetconnectorv1) - Reset a connector
* [RetryTransferInitiation](docs/sdks/formancepaymentsv1/README.md#retrytransferinitiation) - Retry a failed transfer initiation
* [ReverseTransferInitiation](docs/sdks/formancepaymentsv1/README.md#reversetransferinitiation) - Reverse a transfer initiation
* [UdpateTransferInitiationStatus](docs/sdks/formancepaymentsv1/README.md#udpatetransferinitiationstatus) - Update the status of a transfer initiation
* [~~UninstallConnector~~](docs/sdks/formancepaymentsv1/README.md#uninstallconnector) - Uninstall a connector :warning: **Deprecated**
* [UninstallConnectorV1](docs/sdks/formancepaymentsv1/README.md#uninstallconnectorv1) - Uninstall a connector
* [UpdateBankAccountMetadata](docs/sdks/formancepaymentsv1/README.md#updatebankaccountmetadata) - Update metadata of a bank account
* [UpdateConnectorConfigV1](docs/sdks/formancepaymentsv1/README.md#updateconnectorconfigv1) - Update the config of a connector
* [UpdateMetadata](docs/sdks/formancepaymentsv1/README.md#updatemetadata) - Update metadata


### [Reconciliation.V1](docs/sdks/formancereconciliationv1/README.md)

* [CreatePolicy](docs/sdks/formancereconciliationv1/README.md#createpolicy) - Create a policy
* [DeletePolicy](docs/sdks/formancereconciliationv1/README.md#deletepolicy) - Delete a policy
* [GetPolicy](docs/sdks/formancereconciliationv1/README.md#getpolicy) - Get a policy
* [GetReconciliation](docs/sdks/formancereconciliationv1/README.md#getreconciliation) - Get a reconciliation
* [ListPolicies](docs/sdks/formancereconciliationv1/README.md#listpolicies) - List policies
* [ListReconciliations](docs/sdks/formancereconciliationv1/README.md#listreconciliations) - List reconciliations
* [Reconcile](docs/sdks/formancereconciliationv1/README.md#reconcile) - Reconcile using a policy
* [ReconciliationgetServerInfo](docs/sdks/formancereconciliationv1/README.md#reconciliationgetserverinfo) - Get server info


### [Search.V1](docs/sdks/formancesearchv1/README.md)

* [Search](docs/sdks/formancesearchv1/README.md#search) - search.v1
* [SearchgetServerInfo](docs/sdks/formancesearchv1/README.md#searchgetserverinfo) - Get server info


### [Wallets.V1](docs/sdks/formancewalletsv1/README.md)

* [ConfirmHold](docs/sdks/formancewalletsv1/README.md#confirmhold) - Confirm a hold
* [CreateBalance](docs/sdks/formancewalletsv1/README.md#createbalance) - Create a balance
* [CreateWallet](docs/sdks/formancewalletsv1/README.md#createwallet) - Create a new wallet
* [CreditWallet](docs/sdks/formancewalletsv1/README.md#creditwallet) - Credit a wallet
* [DebitWallet](docs/sdks/formancewalletsv1/README.md#debitwallet) - Debit a wallet
* [GetBalance](docs/sdks/formancewalletsv1/README.md#getbalance) - Get detailed balance
* [GetHold](docs/sdks/formancewalletsv1/README.md#gethold) - Get a hold
* [GetHolds](docs/sdks/formancewalletsv1/README.md#getholds) - Get all holds for a wallet
* [GetTransactions](docs/sdks/formancewalletsv1/README.md#gettransactions)
* [GetWallet](docs/sdks/formancewalletsv1/README.md#getwallet) - Get a wallet
* [GetWalletSummary](docs/sdks/formancewalletsv1/README.md#getwalletsummary) - Get wallet summary
* [ListBalances](docs/sdks/formancewalletsv1/README.md#listbalances) - List balances of a wallet
* [ListWallets](docs/sdks/formancewalletsv1/README.md#listwallets) - List all wallets
* [UpdateWallet](docs/sdks/formancewalletsv1/README.md#updatewallet) - Update a wallet
* [VoidHold](docs/sdks/formancewalletsv1/README.md#voidhold) - Cancel a hold
* [WalletsgetServerInfo](docs/sdks/formancewalletsv1/README.md#walletsgetserverinfo) - Get server info


### [Webhooks.V1](docs/sdks/formancewebhooksv1/README.md)

* [ActivateConfig](docs/sdks/formancewebhooksv1/README.md#activateconfig) - Activate one config
* [ChangeConfigSecret](docs/sdks/formancewebhooksv1/README.md#changeconfigsecret) - Change the signing secret of a config
* [DeactivateConfig](docs/sdks/formancewebhooksv1/README.md#deactivateconfig) - Deactivate one config
* [DeleteConfig](docs/sdks/formancewebhooksv1/README.md#deleteconfig) - Delete one config
* [GetManyConfigs](docs/sdks/formancewebhooksv1/README.md#getmanyconfigs) - Get many configs
* [InsertConfig](docs/sdks/formancewebhooksv1/README.md#insertconfig) - Insert a new config
* [TestConfig](docs/sdks/formancewebhooksv1/README.md#testconfig) - Test one config
<!-- End Available Resources and Operations [operations] -->

<!-- Start Error Handling [errors] -->
## Error Handling

Handling errors in this SDK should largely match your expectations.  All operations return a response object or an error, they will never return both.  When specified by the OpenAPI spec document, the SDK will return the appropriate subclass.

| Error Object            | Status Code             | Content Type            |
| ----------------------- | ----------------------- | ----------------------- |
| sdkerrors.ErrorResponse | default                 | application/json        |
| sdkerrors.SDKError      | 4xx-5xx                 | */*                     |

### Example

```go
package main

import (
	"context"
	"errors"
	"github.com/formancehq/formance-sdk-go/v2"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/sdkerrors"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"log"
	"math/big"
	"os"
)

func main() {
	s := v2.New(
		v2.WithSecurity(shared.Security{
			Authorization: os.Getenv("AUTHORIZATION"),
		}),
	)
	request := operations.CreateTransactionsRequest{
		Transactions: shared.Transactions{
			Transactions: []shared.TransactionData{
				shared.TransactionData{
					Postings: []shared.Posting{
						shared.Posting{
							Amount:      big.NewInt(100),
							Asset:       "COIN",
							Destination: "users:002",
							Source:      "users:001",
						},
					},
					Reference: v2.String("ref:001"),
				},
			},
		},
		Ledger: "ledger001",
	}
	ctx := context.Background()
	res, err := s.Ledger.V1.CreateTransactions(ctx, request)
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
	"github.com/formancehq/formance-sdk-go/v2"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"log"
	"os"
)

func main() {
	s := v2.New(
		v2.WithServerIndex(0),
		v2.WithSecurity(shared.Security{
			Authorization: os.Getenv("AUTHORIZATION"),
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
	"github.com/formancehq/formance-sdk-go/v2"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"log"
	"os"
)

func main() {
	s := v2.New(
		v2.WithServerURL("http://localhost"),
		v2.WithSecurity(shared.Security{
			Authorization: os.Getenv("AUTHORIZATION"),
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
	"github.com/formancehq/formance-sdk-go/v2"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"log"
	"os"
)

func main() {
	s := v2.New(
		v2.WithSecurity(shared.Security{
			Authorization: os.Getenv("AUTHORIZATION"),
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

<!-- Start Retries [retries] -->
## Retries

Some of the endpoints in this SDK support retries. If you use the SDK without any configuration, it will fall back to the default retry strategy provided by the API. However, the default retry strategy can be overridden on a per-operation basis, or across the entire SDK.

To change the default retry strategy for a single API call, simply provide a `retry.Config` object to the call by using the `WithRetries` option:
```go
package main

import (
	"context"
	"github.com/formancehq/formance-sdk-go/v2"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/formancehq/formance-sdk-go/v2/pkg/retry"
	"log"
	"os"
	"pkg/models/operations"
)

func main() {
	s := v2.New(
		v2.WithSecurity(shared.Security{
			Authorization: os.Getenv("AUTHORIZATION"),
		}),
	)

	ctx := context.Background()
	res, err := s.GetVersions(ctx, operations.WithRetries(
		retry.Config{
			Strategy: "backoff",
			Backoff: &retry.BackoffStrategy{
				InitialInterval: 1,
				MaxInterval:     50,
				Exponent:        1.1,
				MaxElapsedTime:  100,
			},
			RetryConnectionErrors: false,
		}))
	if err != nil {
		log.Fatal(err)
	}
	if res.GetVersionsResponse != nil {
		// handle response
	}
}

```

If you'd like to override the default retry strategy for all operations that support retries, you can use the `WithRetryConfig` option at SDK initialization:
```go
package main

import (
	"context"
	"github.com/formancehq/formance-sdk-go/v2"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/formancehq/formance-sdk-go/v2/pkg/retry"
	"log"
	"os"
)

func main() {
	s := v2.New(
		v2.WithRetryConfig(
			retry.Config{
				Strategy: "backoff",
				Backoff: &retry.BackoffStrategy{
					InitialInterval: 1,
					MaxInterval:     50,
					Exponent:        1.1,
					MaxElapsedTime:  100,
				},
				RetryConnectionErrors: false,
			}),
		v2.WithSecurity(shared.Security{
			Authorization: os.Getenv("AUTHORIZATION"),
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
<!-- End Retries [retries] -->

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
