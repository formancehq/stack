# github.com/formancehq/sdk-go

<!-- Start SDK Installation -->
## SDK Installation

```bash
go get github.com/formancehq/formance-sdk-go
```
<!-- End SDK Installation -->

## SDK Example Usage
<!-- Start SDK Example Usage -->
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
    res, err := s.GetVersions(ctx)
    if err != nil {
        log.Fatal(err)
    }

    if res.GetVersionsResponse != nil {
        // handle response
    }
}
```
<!-- End SDK Example Usage -->

<!-- Start SDK Available Operations -->
## Available Resources and Operations

### [Formance SDK](docs/formance/README.md)

* [GetVersions](docs/formance/README.md#getversions) - Show stack version information

### [Auth](docs/auth/README.md)

* [AddScopeToClient](docs/auth/README.md#addscopetoclient) - Add scope to client
* [AddTransientScope](docs/auth/README.md#addtransientscope) - Add a transient scope to a scope
* [CreateClient](docs/auth/README.md#createclient) - Create client
* [CreateScope](docs/auth/README.md#createscope) - Create scope
* [CreateSecret](docs/auth/README.md#createsecret) - Add a secret to a client
* [DeleteClient](docs/auth/README.md#deleteclient) - Delete client
* [DeleteScope](docs/auth/README.md#deletescope) - Delete scope
* [DeleteScopeFromClient](docs/auth/README.md#deletescopefromclient) - Delete scope from client
* [DeleteSecret](docs/auth/README.md#deletesecret) - Delete a secret from a client
* [DeleteTransientScope](docs/auth/README.md#deletetransientscope) - Delete a transient scope from a scope
* [GetServerInfo](docs/auth/README.md#getserverinfo) - Get server info
* [ListClients](docs/auth/README.md#listclients) - List clients
* [ListScopes](docs/auth/README.md#listscopes) - List scopes
* [ListUsers](docs/auth/README.md#listusers) - List users
* [ReadClient](docs/auth/README.md#readclient) - Read client
* [ReadScope](docs/auth/README.md#readscope) - Read scope
* [ReadUser](docs/auth/README.md#readuser) - Read user
* [UpdateClient](docs/auth/README.md#updateclient) - Update client
* [UpdateScope](docs/auth/README.md#updatescope) - Update scope

### [Ledger](docs/ledger/README.md)

* [AddMetadataOnTransaction](docs/ledger/README.md#addmetadataontransaction) - Set the metadata of a transaction by its ID
* [AddMetadataToAccount](docs/ledger/README.md#addmetadatatoaccount) - Add metadata to an account
* [CountAccounts](docs/ledger/README.md#countaccounts) - Count the accounts from a ledger
* [CountTransactions](docs/ledger/README.md#counttransactions) - Count the transactions from a ledger
* [CreateTransaction](docs/ledger/README.md#createtransaction) - Create a new transaction to a ledger
* [GetAccount](docs/ledger/README.md#getaccount) - Get account by its address
* [GetBalances](docs/ledger/README.md#getbalances) - Get the balances from a ledger's account
* [GetBalancesAggregated](docs/ledger/README.md#getbalancesaggregated) - Get the aggregated balances from selected accounts
* [GetInfo](docs/ledger/README.md#getinfo) - Show server information
* [GetLedgerInfo](docs/ledger/README.md#getledgerinfo) - Get information about a ledger
* [GetTransaction](docs/ledger/README.md#gettransaction) - Get transaction from a ledger by its ID
* [ListAccounts](docs/ledger/README.md#listaccounts) - List accounts from a ledger
* [ListLogs](docs/ledger/README.md#listlogs) - List the logs from a ledger
* [ListTransactions](docs/ledger/README.md#listtransactions) - List transactions from a ledger
* [ReadStats](docs/ledger/README.md#readstats) - Get statistics from a ledger
* [RevertTransaction](docs/ledger/README.md#reverttransaction) - Revert a ledger transaction by its ID

### [Orchestration](docs/orchestration/README.md)

* [CancelEvent](docs/orchestration/README.md#cancelevent) - Cancel a running workflow
* [CreateWorkflow](docs/orchestration/README.md#createworkflow) - Create workflow
* [GetInstance](docs/orchestration/README.md#getinstance) - Get a workflow instance by id
* [GetInstanceHistory](docs/orchestration/README.md#getinstancehistory) - Get a workflow instance history by id
* [GetInstanceStageHistory](docs/orchestration/README.md#getinstancestagehistory) - Get a workflow instance stage history
* [GetWorkflow](docs/orchestration/README.md#getworkflow) - Get a flow by id
* [ListInstances](docs/orchestration/README.md#listinstances) - List instances of a workflow
* [ListWorkflows](docs/orchestration/README.md#listworkflows) - List registered workflows
* [OrchestrationgetServerInfo](docs/orchestration/README.md#orchestrationgetserverinfo) - Get server info
* [RunWorkflow](docs/orchestration/README.md#runworkflow) - Run workflow
* [SendEvent](docs/orchestration/README.md#sendevent) - Send an event to a running workflow

### [Payments](docs/payments/README.md)

* [ConnectorsStripeTransfer](docs/payments/README.md#connectorsstripetransfer) - Transfer funds between Stripe accounts
* [ConnectorsTransfer](docs/payments/README.md#connectorstransfer) - Transfer funds between Connector accounts
* [GetConnectorTask](docs/payments/README.md#getconnectortask) - Read a specific task of the connector
* [GetPayment](docs/payments/README.md#getpayment) - Get a payment
* [InstallConnector](docs/payments/README.md#installconnector) - Install a connector
* [ListAllConnectors](docs/payments/README.md#listallconnectors) - List all installed connectors
* [ListConfigsAvailableConnectors](docs/payments/README.md#listconfigsavailableconnectors) - List the configs of each available connector
* [ListConnectorTasks](docs/payments/README.md#listconnectortasks) - List tasks from a connector
* [ListConnectorsTransfers](docs/payments/README.md#listconnectorstransfers) - List transfers and their statuses
* [ListPayments](docs/payments/README.md#listpayments) - List payments
* [PaymentsgetServerInfo](docs/payments/README.md#paymentsgetserverinfo) - Get server info
* [PaymentslistAccounts](docs/payments/README.md#paymentslistaccounts) - List accounts
* [ReadConnectorConfig](docs/payments/README.md#readconnectorconfig) - Read the config of a connector
* [ResetConnector](docs/payments/README.md#resetconnector) - Reset a connector
* [UninstallConnector](docs/payments/README.md#uninstallconnector) - Uninstall a connector
* [UpdateMetadata](docs/payments/README.md#updatemetadata) - Update metadata

### [Search](docs/search/README.md)

* [Search](docs/search/README.md#search) - Search
* [SearchgetServerInfo](docs/search/README.md#searchgetserverinfo) - Get server info

### [Wallets](docs/wallets/README.md)

* [ConfirmHold](docs/wallets/README.md#confirmhold) - Confirm a hold
* [CreateBalance](docs/wallets/README.md#createbalance) - Create a balance
* [CreateWallet](docs/wallets/README.md#createwallet) - Create a new wallet
* [CreditWallet](docs/wallets/README.md#creditwallet) - Credit a wallet
* [DebitWallet](docs/wallets/README.md#debitwallet) - Debit a wallet
* [GetBalance](docs/wallets/README.md#getbalance) - Get detailed balance
* [GetHold](docs/wallets/README.md#gethold) - Get a hold
* [GetHolds](docs/wallets/README.md#getholds) - Get all holds for a wallet
* [GetTransactions](docs/wallets/README.md#gettransactions)
* [GetWallet](docs/wallets/README.md#getwallet) - Get a wallet
* [GetWalletSummary](docs/wallets/README.md#getwalletsummary) - Get wallet summary
* [ListBalances](docs/wallets/README.md#listbalances) - List balances of a wallet
* [ListWallets](docs/wallets/README.md#listwallets) - List all wallets
* [UpdateWallet](docs/wallets/README.md#updatewallet) - Update a wallet
* [VoidHold](docs/wallets/README.md#voidhold) - Cancel a hold
* [WalletsgetServerInfo](docs/wallets/README.md#walletsgetserverinfo) - Get server info

### [Webhooks](docs/webhooks/README.md)

* [ActivateConfig](docs/webhooks/README.md#activateconfig) - Activate one config
* [ChangeConfigSecret](docs/webhooks/README.md#changeconfigsecret) - Change the signing secret of a config
* [DeactivateConfig](docs/webhooks/README.md#deactivateconfig) - Deactivate one config
* [DeleteConfig](docs/webhooks/README.md#deleteconfig) - Delete one config
* [GetManyConfigs](docs/webhooks/README.md#getmanyconfigs) - Get many configs
* [InsertConfig](docs/webhooks/README.md#insertconfig) - Insert a new config
* [TestConfig](docs/webhooks/README.md#testconfig) - Test one config
<!-- End SDK Available Operations -->

### Maturity

This SDK is in beta and therefore, we recommend pinning usage to a specific package version.
This way, you can install the same version each time without breaking changes unless you are intentionally
looking for the latest version.

### Contributions

While we value open-source contributions to this SDK, this library is generated and maintained programmatically.
Feel free to open a PR or a Github issue as a proof of concept and we'll do our best to include it in a future release !

### SDK Created by [Speakeasy](https://docs.speakeasyapi.dev/docs/using-speakeasy/client-sdks)
