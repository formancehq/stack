<div align="center">
    <picture>
        <source srcset="https://user-images.githubusercontent.com/6267663/221572723-e77f55a3-5d19-4a13-94f8-e7b0b340d71e.svg" media="(prefers-color-scheme: dark)">
        <img src="https://user-images.githubusercontent.com/6267663/221572726-6982541c-d1cf-4d9f-9bbf-cd774a2713e6.svg">
    </picture>
   <h1>Formance Typescript SDK</h1>
   <p><strong>Open Source Ledger for money-moving platforms</strong></p>
   <p>Build and track custom fit money flows on a scalable financial infrastructure.</p>
   <a href="https://docs.formance.com"><img src="https://img.shields.io/static/v1?label=Docs&message=Docs&color=000&style=for-the-badge" /></a>
   <a href="https://join.slack.com/t/formance-community/shared_invite/zt-1of48xmgy-Jc6RH8gzcWf5D0qD2HBPQA"><img src="https://img.shields.io/static/v1?label=Slack&message=Join&color=7289da&style=for-the-badge" /></a>
  <a href="https://opensource.org/licenses/MIT"><img src="https://img.shields.io/badge/License-MIT-blue.svg?style=for-the-badge" /></a>
</div>

<!-- Start SDK Installation -->
## SDK Installation

### NPM

```bash
npm add <UNSET>
```

### Yarn

```bash
yarn add <UNSET>
```
<!-- End SDK Installation -->

## SDK Example Usage
<!-- Start SDK Example Usage -->
```typescript
import { SDK } from "@formance/formance-sdk";
import { GetVersionsResponse } from "@formance/formance-sdk/dist/sdk/models/operations";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.getVersions().then((res: GetVersionsResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```
<!-- End SDK Example Usage -->

<!-- Start SDK Available Operations -->
## Available Resources and Operations

### [SDK](docs/sdk/README.md)

* [getVersions](docs/sdk/README.md#getversions) - Show stack version information

### [accounts](docs/accounts/README.md)

* [addMetadataToAccount](docs/accounts/README.md#addmetadatatoaccount) - Add metadata to an account
* [countAccounts](docs/accounts/README.md#countaccounts) - Count the accounts from a ledger
* [getAccount](docs/accounts/README.md#getaccount) - Get account by its address
* [listAccounts](docs/accounts/README.md#listaccounts) - List accounts from a ledger

### [auth](docs/auth/README.md)

* [addScopeToClient](docs/auth/README.md#addscopetoclient) - Add scope to client
* [addTransientScope](docs/auth/README.md#addtransientscope) - Add a transient scope to a scope
* [createClient](docs/auth/README.md#createclient) - Create client
* [createScope](docs/auth/README.md#createscope) - Create scope
* [createSecret](docs/auth/README.md#createsecret) - Add a secret to a client
* [deleteClient](docs/auth/README.md#deleteclient) - Delete client
* [deleteScope](docs/auth/README.md#deletescope) - Delete scope
* [deleteScopeFromClient](docs/auth/README.md#deletescopefromclient) - Delete scope from client
* [deleteSecret](docs/auth/README.md#deletesecret) - Delete a secret from a client
* [deleteTransientScope](docs/auth/README.md#deletetransientscope) - Delete a transient scope from a scope
* [getServerInfo](docs/auth/README.md#getserverinfo) - Get server info
* [listClients](docs/auth/README.md#listclients) - List clients
* [listScopes](docs/auth/README.md#listscopes) - List scopes
* [listUsers](docs/auth/README.md#listusers) - List users
* [readClient](docs/auth/README.md#readclient) - Read client
* [readScope](docs/auth/README.md#readscope) - Read scope
* [readUser](docs/auth/README.md#readuser) - Read user
* [updateClient](docs/auth/README.md#updateclient) - Update client
* [updateScope](docs/auth/README.md#updatescope) - Update scope

### [balances](docs/balances/README.md)

* [getBalances](docs/balances/README.md#getbalances) - Get the balances from a ledger's account
* [getBalancesAggregated](docs/balances/README.md#getbalancesaggregated) - Get the aggregated balances from selected accounts

### [ledger](docs/ledger/README.md)

* [getLedgerInfo](docs/ledger/README.md#getledgerinfo) - Get information about a ledger

### [logs](docs/logs/README.md)

* [listLogs](docs/logs/README.md#listlogs) - List the logs from a ledger

### [mapping](docs/mapping/README.md)

* [getMapping](docs/mapping/README.md#getmapping) - Get the mapping of a ledger
* [updateMapping](docs/mapping/README.md#updatemapping) - Update the mapping of a ledger

### [orchestration](docs/orchestration/README.md)

* [cancelEvent](docs/orchestration/README.md#cancelevent) - Cancel a running workflow
* [createWorkflow](docs/orchestration/README.md#createworkflow) - Create workflow
* [deleteWorkflow](docs/orchestration/README.md#deleteworkflow) - Delete a flow by id
* [getInstance](docs/orchestration/README.md#getinstance) - Get a workflow instance by id
* [getInstanceHistory](docs/orchestration/README.md#getinstancehistory) - Get a workflow instance history by id
* [getInstanceStageHistory](docs/orchestration/README.md#getinstancestagehistory) - Get a workflow instance stage history
* [getWorkflow](docs/orchestration/README.md#getworkflow) - Get a flow by id
* [listInstances](docs/orchestration/README.md#listinstances) - List instances of a workflow
* [listWorkflows](docs/orchestration/README.md#listworkflows) - List registered workflows
* [orchestrationgetServerInfo](docs/orchestration/README.md#orchestrationgetserverinfo) - Get server info
* [runWorkflow](docs/orchestration/README.md#runworkflow) - Run workflow
* [sendEvent](docs/orchestration/README.md#sendevent) - Send an event to a running workflow

### [payments](docs/payments/README.md)

* [connectorsTransfer](docs/payments/README.md#connectorstransfer) - Transfer funds between Connector accounts
* [createBankAccount](docs/payments/README.md#createbankaccount) - Create a BankAccount in Payments and on the PSP
* [createTransferInitiation](docs/payments/README.md#createtransferinitiation) - Create a TransferInitiation
* [deleteTransferInitiation](docs/payments/README.md#deletetransferinitiation) - Delete a transfer initiation
* [getAccountBalances](docs/payments/README.md#getaccountbalances) - Get account balances
* [getBankAccount](docs/payments/README.md#getbankaccount) - Get a bank account created by user on Formance
* [getConnectorTask](docs/payments/README.md#getconnectortask) - Read a specific task of the connector
* [getPayment](docs/payments/README.md#getpayment) - Get a payment
* [getTransferInitiation](docs/payments/README.md#gettransferinitiation) - Get a transfer initiation
* [installConnector](docs/payments/README.md#installconnector) - Install a connector
* [listAllConnectors](docs/payments/README.md#listallconnectors) - List all installed connectors
* [listBankAccounts](docs/payments/README.md#listbankaccounts) - List bank accounts created by user on Formance
* [listConfigsAvailableConnectors](docs/payments/README.md#listconfigsavailableconnectors) - List the configs of each available connector
* [listConnectorTasks](docs/payments/README.md#listconnectortasks) - List tasks from a connector
* [listConnectorsTransfers](docs/payments/README.md#listconnectorstransfers) - List transfers and their statuses
* [listPayments](docs/payments/README.md#listpayments) - List payments
* [listTransferInitiations](docs/payments/README.md#listtransferinitiations) - List Transfer Initiations
* [paymentsgetAccount](docs/payments/README.md#paymentsgetaccount) - Get an account
* [paymentsgetServerInfo](docs/payments/README.md#paymentsgetserverinfo) - Get server info
* [paymentslistAccounts](docs/payments/README.md#paymentslistaccounts) - List accounts
* [readConnectorConfig](docs/payments/README.md#readconnectorconfig) - Read the config of a connector
* [resetConnector](docs/payments/README.md#resetconnector) - Reset a connector
* [retryTransferInitiation](docs/payments/README.md#retrytransferinitiation) - Retry a failed transfer initiation
* [udpateTransferInitiationStatus](docs/payments/README.md#udpatetransferinitiationstatus) - Update the status of a transfer initiation
* [uninstallConnector](docs/payments/README.md#uninstallconnector) - Uninstall a connector
* [updateMetadata](docs/payments/README.md#updatemetadata) - Update metadata

### [script](docs/script/README.md)

* [~~runScript~~](docs/script/README.md#runscript) - Execute a Numscript :warning: **Deprecated**

### [search](docs/search/README.md)

* [search](docs/search/README.md#search) - Search
* [searchgetServerInfo](docs/search/README.md#searchgetserverinfo) - Get server info

### [server](docs/server/README.md)

* [getInfo](docs/server/README.md#getinfo) - Show server information

### [stats](docs/stats/README.md)

* [readStats](docs/stats/README.md#readstats) - Get statistics from a ledger

### [transactions](docs/transactions/README.md)

* [createTransactions](docs/transactions/README.md#createtransactions) - Create a new batch of transactions to a ledger
* [addMetadataOnTransaction](docs/transactions/README.md#addmetadataontransaction) - Set the metadata of a transaction by its ID
* [countTransactions](docs/transactions/README.md#counttransactions) - Count the transactions from a ledger
* [createTransaction](docs/transactions/README.md#createtransaction) - Create a new transaction to a ledger
* [getTransaction](docs/transactions/README.md#gettransaction) - Get transaction from a ledger by its ID
* [listTransactions](docs/transactions/README.md#listtransactions) - List transactions from a ledger
* [revertTransaction](docs/transactions/README.md#reverttransaction) - Revert a ledger transaction by its ID

### [wallets](docs/wallets/README.md)

* [confirmHold](docs/wallets/README.md#confirmhold) - Confirm a hold
* [createBalance](docs/wallets/README.md#createbalance) - Create a balance
* [createWallet](docs/wallets/README.md#createwallet) - Create a new wallet
* [creditWallet](docs/wallets/README.md#creditwallet) - Credit a wallet
* [debitWallet](docs/wallets/README.md#debitwallet) - Debit a wallet
* [getBalance](docs/wallets/README.md#getbalance) - Get detailed balance
* [getHold](docs/wallets/README.md#gethold) - Get a hold
* [getHolds](docs/wallets/README.md#getholds) - Get all holds for a wallet
* [getTransactions](docs/wallets/README.md#gettransactions)
* [getWallet](docs/wallets/README.md#getwallet) - Get a wallet
* [getWalletSummary](docs/wallets/README.md#getwalletsummary) - Get wallet summary
* [listBalances](docs/wallets/README.md#listbalances) - List balances of a wallet
* [listWallets](docs/wallets/README.md#listwallets) - List all wallets
* [updateWallet](docs/wallets/README.md#updatewallet) - Update a wallet
* [voidHold](docs/wallets/README.md#voidhold) - Cancel a hold
* [walletsgetServerInfo](docs/wallets/README.md#walletsgetserverinfo) - Get server info

### [webhooks](docs/webhooks/README.md)

* [activateConfig](docs/webhooks/README.md#activateconfig) - Activate one config
* [changeConfigSecret](docs/webhooks/README.md#changeconfigsecret) - Change the signing secret of a config
* [deactivateConfig](docs/webhooks/README.md#deactivateconfig) - Deactivate one config
* [deleteConfig](docs/webhooks/README.md#deleteconfig) - Delete one config
* [getManyConfigs](docs/webhooks/README.md#getmanyconfigs) - Get many configs
* [insertConfig](docs/webhooks/README.md#insertconfig) - Insert a new config
* [testConfig](docs/webhooks/README.md#testconfig) - Test one config
<!-- End SDK Available Operations -->

### SDK Generated by [Speakeasy](https://docs.speakeasyapi.dev/docs/using-speakeasy/client-sdks)
