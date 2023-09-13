<div align="center">
    <picture>
        <source srcset="https://user-images.githubusercontent.com/6267663/221572723-e77f55a3-5d19-4a13-94f8-e7b0b340d71e.svg" media="(prefers-color-scheme: dark)">
        <img src="https://user-images.githubusercontent.com/6267663/221572726-6982541c-d1cf-4d9f-9bbf-cd774a2713e6.svg">
    </picture>
   <h1>Formance PHP SDK</h1>
   <p><strong>Open Source Ledger for money-moving platforms</strong></p>
   <p>Build and track custom fit money flows on a scalable financial infrastructure.</p>
   <a href="https://docs.formance.com"><img src="https://img.shields.io/static/v1?label=Docs&message=Docs&color=000&style=for-the-badge" /></a>
   <a href="https://join.slack.com/t/formance-community/shared_invite/zt-1of48xmgy-Jc6RH8gzcWf5D0qD2HBPQA"><img src="https://img.shields.io/static/v1?label=Slack&message=Join&color=7289da&style=for-the-badge" /></a>
  <a href="https://opensource.org/licenses/MIT"><img src="https://img.shields.io/badge/License-MIT-blue.svg?style=for-the-badge" /></a>
</div>

<!-- Start SDK Installation -->
## SDK Installation

### Composer

To install the SDK first add the below to your `composer.json` file:

```json
{
    "repositories": [
        {
            "type": "github",
            "url": "<UNSET>.git"
        }
    ],
    "require": {
        "formance-sdk-php": "*"
    }
}
```

Then run the following command:

```bash
composer update
```
<!-- End SDK Installation -->
## SDK Example Usage
<!-- Start SDK Example Usage -->


```php
<?php

declare(strict_types=1);
require_once 'vendor/autoload.php';

use \formance\stack\SDK;
use \formance\stack\Models\Shared\Security;

$sdk = SDK::builder()
    ->build();

try {
    $response = $sdk->getVersions();

    if ($response->getVersionsResponse !== null) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```
<!-- End SDK Example Usage -->

<!-- Start SDK Available Operations -->
## Available Resources and Operations

### [SDK](docs/sdks/sdk/README.md)

* [getVersions](docs/sdks/sdk/README.md#getversions) - Show stack version information

### [accounts](docs/sdks/accounts/README.md)

* [addMetadataToAccount](docs/sdks/accounts/README.md#addmetadatatoaccount) - Add metadata to an account
* [countAccounts](docs/sdks/accounts/README.md#countaccounts) - Count the accounts from a ledger
* [getAccount](docs/sdks/accounts/README.md#getaccount) - Get account by its address
* [listAccounts](docs/sdks/accounts/README.md#listaccounts) - List accounts from a ledger

### [auth](docs/sdks/auth/README.md)

* [addScopeToClient](docs/sdks/auth/README.md#addscopetoclient) - Add scope to client
* [addTransientScope](docs/sdks/auth/README.md#addtransientscope) - Add a transient scope to a scope
* [createClient](docs/sdks/auth/README.md#createclient) - Create client
* [createScope](docs/sdks/auth/README.md#createscope) - Create scope
* [createSecret](docs/sdks/auth/README.md#createsecret) - Add a secret to a client
* [deleteClient](docs/sdks/auth/README.md#deleteclient) - Delete client
* [deleteScope](docs/sdks/auth/README.md#deletescope) - Delete scope
* [deleteScopeFromClient](docs/sdks/auth/README.md#deletescopefromclient) - Delete scope from client
* [deleteSecret](docs/sdks/auth/README.md#deletesecret) - Delete a secret from a client
* [deleteTransientScope](docs/sdks/auth/README.md#deletetransientscope) - Delete a transient scope from a scope
* [getServerInfo](docs/sdks/auth/README.md#getserverinfo) - Get server info
* [listClients](docs/sdks/auth/README.md#listclients) - List clients
* [listScopes](docs/sdks/auth/README.md#listscopes) - List scopes
* [listUsers](docs/sdks/auth/README.md#listusers) - List users
* [readClient](docs/sdks/auth/README.md#readclient) - Read client
* [readScope](docs/sdks/auth/README.md#readscope) - Read scope
* [readUser](docs/sdks/auth/README.md#readuser) - Read user
* [updateClient](docs/sdks/auth/README.md#updateclient) - Update client
* [updateScope](docs/sdks/auth/README.md#updatescope) - Update scope

### [balances](docs/sdks/balances/README.md)

* [getBalances](docs/sdks/balances/README.md#getbalances) - Get the balances from a ledger's account
* [getBalancesAggregated](docs/sdks/balances/README.md#getbalancesaggregated) - Get the aggregated balances from selected accounts

### [ledger](docs/sdks/ledger/README.md)

* [getLedgerInfo](docs/sdks/ledger/README.md#getledgerinfo) - Get information about a ledger

### [logs](docs/sdks/logs/README.md)

* [listLogs](docs/sdks/logs/README.md#listlogs) - List the logs from a ledger

### [mapping](docs/sdks/mapping/README.md)

* [getMapping](docs/sdks/mapping/README.md#getmapping) - Get the mapping of a ledger
* [updateMapping](docs/sdks/mapping/README.md#updatemapping) - Update the mapping of a ledger

### [orchestration](docs/sdks/orchestration/README.md)

* [cancelEvent](docs/sdks/orchestration/README.md#cancelevent) - Cancel a running workflow
* [createWorkflow](docs/sdks/orchestration/README.md#createworkflow) - Create workflow
* [deleteWorkflow](docs/sdks/orchestration/README.md#deleteworkflow) - Delete a flow by id
* [getInstance](docs/sdks/orchestration/README.md#getinstance) - Get a workflow instance by id
* [getInstanceHistory](docs/sdks/orchestration/README.md#getinstancehistory) - Get a workflow instance history by id
* [getInstanceStageHistory](docs/sdks/orchestration/README.md#getinstancestagehistory) - Get a workflow instance stage history
* [getWorkflow](docs/sdks/orchestration/README.md#getworkflow) - Get a flow by id
* [listInstances](docs/sdks/orchestration/README.md#listinstances) - List instances of a workflow
* [listWorkflows](docs/sdks/orchestration/README.md#listworkflows) - List registered workflows
* [orchestrationgetServerInfo](docs/sdks/orchestration/README.md#orchestrationgetserverinfo) - Get server info
* [runWorkflow](docs/sdks/orchestration/README.md#runworkflow) - Run workflow
* [sendEvent](docs/sdks/orchestration/README.md#sendevent) - Send an event to a running workflow

### [payments](docs/sdks/payments/README.md)

* [connectorsStripeTransfer](docs/sdks/payments/README.md#connectorsstripetransfer) - Transfer funds between Stripe accounts
* [connectorsTransfer](docs/sdks/payments/README.md#connectorstransfer) - Transfer funds between Connector accounts
* [getAccountBalances](docs/sdks/payments/README.md#getaccountbalances) - Get account balances
* [getConnectorTask](docs/sdks/payments/README.md#getconnectortask) - Read a specific task of the connector
* [getPayment](docs/sdks/payments/README.md#getpayment) - Get a payment
* [installConnector](docs/sdks/payments/README.md#installconnector) - Install a connector
* [listAllConnectors](docs/sdks/payments/README.md#listallconnectors) - List all installed connectors
* [listConfigsAvailableConnectors](docs/sdks/payments/README.md#listconfigsavailableconnectors) - List the configs of each available connector
* [listConnectorTasks](docs/sdks/payments/README.md#listconnectortasks) - List tasks from a connector
* [listConnectorsTransfers](docs/sdks/payments/README.md#listconnectorstransfers) - List transfers and their statuses
* [listPayments](docs/sdks/payments/README.md#listpayments) - List payments
* [paymentsgetAccount](docs/sdks/payments/README.md#paymentsgetaccount) - Get an account
* [paymentsgetServerInfo](docs/sdks/payments/README.md#paymentsgetserverinfo) - Get server info
* [paymentslistAccounts](docs/sdks/payments/README.md#paymentslistaccounts) - List accounts
* [readConnectorConfig](docs/sdks/payments/README.md#readconnectorconfig) - Read the config of a connector
* [resetConnector](docs/sdks/payments/README.md#resetconnector) - Reset a connector
* [uninstallConnector](docs/sdks/payments/README.md#uninstallconnector) - Uninstall a connector
* [updateMetadata](docs/sdks/payments/README.md#updatemetadata) - Update metadata

### [script](docs/sdks/script/README.md)

* [~~runScript~~](docs/sdks/script/README.md#runscript) - Execute a Numscript :warning: **Deprecated**

### [search](docs/sdks/search/README.md)

* [search](docs/sdks/search/README.md#search) - Search
* [searchgetServerInfo](docs/sdks/search/README.md#searchgetserverinfo) - Get server info

### [server](docs/sdks/server/README.md)

* [getInfo](docs/sdks/server/README.md#getinfo) - Show server information

### [stats](docs/sdks/stats/README.md)

* [readStats](docs/sdks/stats/README.md#readstats) - Get statistics from a ledger

### [transactions](docs/sdks/transactions/README.md)

* [createTransactions](docs/sdks/transactions/README.md#createtransactions) - Create a new batch of transactions to a ledger
* [addMetadataOnTransaction](docs/sdks/transactions/README.md#addmetadataontransaction) - Set the metadata of a transaction by its ID
* [countTransactions](docs/sdks/transactions/README.md#counttransactions) - Count the transactions from a ledger
* [createTransaction](docs/sdks/transactions/README.md#createtransaction) - Create a new transaction to a ledger
* [getTransaction](docs/sdks/transactions/README.md#gettransaction) - Get transaction from a ledger by its ID
* [listTransactions](docs/sdks/transactions/README.md#listtransactions) - List transactions from a ledger
* [revertTransaction](docs/sdks/transactions/README.md#reverttransaction) - Revert a ledger transaction by its ID

### [wallets](docs/sdks/wallets/README.md)

* [confirmHold](docs/sdks/wallets/README.md#confirmhold) - Confirm a hold
* [createBalance](docs/sdks/wallets/README.md#createbalance) - Create a balance
* [createWallet](docs/sdks/wallets/README.md#createwallet) - Create a new wallet
* [creditWallet](docs/sdks/wallets/README.md#creditwallet) - Credit a wallet
* [debitWallet](docs/sdks/wallets/README.md#debitwallet) - Debit a wallet
* [getBalance](docs/sdks/wallets/README.md#getbalance) - Get detailed balance
* [getHold](docs/sdks/wallets/README.md#gethold) - Get a hold
* [getHolds](docs/sdks/wallets/README.md#getholds) - Get all holds for a wallet
* [getTransactions](docs/sdks/wallets/README.md#gettransactions)
* [getWallet](docs/sdks/wallets/README.md#getwallet) - Get a wallet
* [getWalletSummary](docs/sdks/wallets/README.md#getwalletsummary) - Get wallet summary
* [listBalances](docs/sdks/wallets/README.md#listbalances) - List balances of a wallet
* [listWallets](docs/sdks/wallets/README.md#listwallets) - List all wallets
* [updateWallet](docs/sdks/wallets/README.md#updatewallet) - Update a wallet
* [voidHold](docs/sdks/wallets/README.md#voidhold) - Cancel a hold
* [walletsgetServerInfo](docs/sdks/wallets/README.md#walletsgetserverinfo) - Get server info

### [webhooks](docs/sdks/webhooks/README.md)

* [activateConfig](docs/sdks/webhooks/README.md#activateconfig) - Activate one config
* [changeConfigSecret](docs/sdks/webhooks/README.md#changeconfigsecret) - Change the signing secret of a config
* [deactivateConfig](docs/sdks/webhooks/README.md#deactivateconfig) - Deactivate one config
* [deleteConfig](docs/sdks/webhooks/README.md#deleteconfig) - Delete one config
* [getManyConfigs](docs/sdks/webhooks/README.md#getmanyconfigs) - Get many configs
* [insertConfig](docs/sdks/webhooks/README.md#insertconfig) - Insert a new config
* [testConfig](docs/sdks/webhooks/README.md#testconfig) - Test one config
<!-- End SDK Available Operations -->

### SDK Generated by [Speakeasy](https://docs.speakeasyapi.dev/docs/using-speakeasy/client-sdks)
