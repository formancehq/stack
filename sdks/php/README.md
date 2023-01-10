# OpenAPIClient-php

Open, modular foundation for unique payments flows

# Introduction
This API is documented in **OpenAPI format**.

# Authentication
Formance Stack offers one forms of authentication:
  - OAuth2
OAuth2 - an open protocol to allow secure authorization in a simple
and standard method from web, mobile and desktop applications.
<SecurityDefinitions />


For more information, please visit [https://www.formance.com](https://www.formance.com).

## Installation & Usage

### Requirements

PHP 7.4 and later.
Should also work with PHP 8.0.

### Composer

To install the bindings via [Composer](https://getcomposer.org/), add the following to `composer.json`:

```json
{
  "repositories": [
    {
      "type": "vcs",
      "url": "https://github.com/formancehq/formance-sdk-php.git"
    }
  ],
  "require": {
    "formancehq/formance-sdk-php": "*@dev"
  }
}
```

Then run `composer install`

### Manual Installation

Download the files and include `autoload.php`:

```php
<?php
require_once('/path/to/OpenAPIClient-php/vendor/autoload.php');
```

## Getting Started

Please follow the [installation procedure](#installation--usage) and then run the following:

```php
<?php
require_once(__DIR__ . '/vendor/autoload.php');



// Configure OAuth2 access token for authorization: Authorization
$config = Formance\Configuration::getDefaultConfiguration()->setAccessToken('YOUR_ACCESS_TOKEN');


$apiInstance = new Formance\Api\AccountsApi(
    // If you want use custom http client, pass your client which implements `GuzzleHttp\ClientInterface`.
    // This is optional, `GuzzleHttp\Client` will be used as default.
    new GuzzleHttp\Client(),
    $config
);
$ledger = ledger001; // string | Name of the ledger.
$address = users:001; // string | Exact address of the account. It must match the following regular expressions pattern: ``` ^\\w+(:\\w+)*$ ```
$request_body = NULL; // array<string,mixed> | metadata

try {
    $apiInstance->addMetadataToAccount($ledger, $address, $request_body);
} catch (Exception $e) {
    echo 'Exception when calling AccountsApi->addMetadataToAccount: ', $e->getMessage(), PHP_EOL;
}

```

## API Endpoints

All URIs are relative to *http://localhost*

Class | Method | HTTP request | Description
------------ | ------------- | ------------- | -------------
*AccountsApi* | [**addMetadataToAccount**](docs/Api/AccountsApi.md#addmetadatatoaccount) | **POST** /api/ledger/{ledger}/accounts/{address}/metadata | Add metadata to an account
*AccountsApi* | [**countAccounts**](docs/Api/AccountsApi.md#countaccounts) | **HEAD** /api/ledger/{ledger}/accounts | Count the accounts from a ledger
*AccountsApi* | [**getAccount**](docs/Api/AccountsApi.md#getaccount) | **GET** /api/ledger/{ledger}/accounts/{address} | Get account by its address
*AccountsApi* | [**listAccounts**](docs/Api/AccountsApi.md#listaccounts) | **GET** /api/ledger/{ledger}/accounts | List accounts from a ledger
*BalancesApi* | [**getBalances**](docs/Api/BalancesApi.md#getbalances) | **GET** /api/ledger/{ledger}/balances | Get the balances from a ledger&#39;s account
*BalancesApi* | [**getBalancesAggregated**](docs/Api/BalancesApi.md#getbalancesaggregated) | **GET** /api/ledger/{ledger}/aggregate/balances | Get the aggregated balances from selected accounts
*ClientsApi* | [**addScopeToClient**](docs/Api/ClientsApi.md#addscopetoclient) | **PUT** /api/auth/clients/{clientId}/scopes/{scopeId} | Add scope to client
*ClientsApi* | [**createClient**](docs/Api/ClientsApi.md#createclient) | **POST** /api/auth/clients | Create client
*ClientsApi* | [**createSecret**](docs/Api/ClientsApi.md#createsecret) | **POST** /api/auth/clients/{clientId}/secrets | Add a secret to a client
*ClientsApi* | [**deleteClient**](docs/Api/ClientsApi.md#deleteclient) | **DELETE** /api/auth/clients/{clientId} | Delete client
*ClientsApi* | [**deleteScopeFromClient**](docs/Api/ClientsApi.md#deletescopefromclient) | **DELETE** /api/auth/clients/{clientId}/scopes/{scopeId} | Delete scope from client
*ClientsApi* | [**deleteSecret**](docs/Api/ClientsApi.md#deletesecret) | **DELETE** /api/auth/clients/{clientId}/secrets/{secretId} | Delete a secret from a client
*ClientsApi* | [**listClients**](docs/Api/ClientsApi.md#listclients) | **GET** /api/auth/clients | List clients
*ClientsApi* | [**readClient**](docs/Api/ClientsApi.md#readclient) | **GET** /api/auth/clients/{clientId} | Read client
*ClientsApi* | [**updateClient**](docs/Api/ClientsApi.md#updateclient) | **PUT** /api/auth/clients/{clientId} | Update client
*DefaultApi* | [**getServerInfo**](docs/Api/DefaultApi.md#getserverinfo) | **GET** /api/auth/_info | Get server info
*DefaultApi* | [**searchgetServerInfo**](docs/Api/DefaultApi.md#searchgetserverinfo) | **GET** /api/search/_info | Get server info
*LedgerApi* | [**getLedgerInfo**](docs/Api/LedgerApi.md#getledgerinfo) | **GET** /api/ledger/{ledger}/_info | Get information about a ledger
*LogsApi* | [**listLogs**](docs/Api/LogsApi.md#listlogs) | **GET** /api/ledger/{ledger}/log | List the logs from a ledger
*MappingApi* | [**getMapping**](docs/Api/MappingApi.md#getmapping) | **GET** /api/ledger/{ledger}/mapping | Get the mapping of a ledger
*MappingApi* | [**updateMapping**](docs/Api/MappingApi.md#updatemapping) | **PUT** /api/ledger/{ledger}/mapping | Update the mapping of a ledger
*PaymentsApi* | [**connectorsStripeTransfer**](docs/Api/PaymentsApi.md#connectorsstripetransfer) | **POST** /api/payments/connectors/stripe/transfer | Transfer funds between Stripe accounts
*PaymentsApi* | [**getAllConnectors**](docs/Api/PaymentsApi.md#getallconnectors) | **GET** /api/payments/connectors | Get all installed connectors
*PaymentsApi* | [**getAllConnectorsConfigs**](docs/Api/PaymentsApi.md#getallconnectorsconfigs) | **GET** /api/payments/connectors/configs | Get all available connectors configs
*PaymentsApi* | [**getConnectorTask**](docs/Api/PaymentsApi.md#getconnectortask) | **GET** /api/payments/connectors/{connector}/tasks/{taskId} | Read a specific task of the connector
*PaymentsApi* | [**getPayment**](docs/Api/PaymentsApi.md#getpayment) | **GET** /api/payments/payments/{paymentId} | Returns a payment.
*PaymentsApi* | [**installConnector**](docs/Api/PaymentsApi.md#installconnector) | **POST** /api/payments/connectors/{connector} | Install connector
*PaymentsApi* | [**listConnectorTasks**](docs/Api/PaymentsApi.md#listconnectortasks) | **GET** /api/payments/connectors/{connector}/tasks | List connector tasks
*PaymentsApi* | [**listPayments**](docs/Api/PaymentsApi.md#listpayments) | **GET** /api/payments/payments | Returns a list of payments.
*PaymentsApi* | [**paymentslistAccounts**](docs/Api/PaymentsApi.md#paymentslistaccounts) | **GET** /api/payments/accounts | Returns a list of accounts.
*PaymentsApi* | [**readConnectorConfig**](docs/Api/PaymentsApi.md#readconnectorconfig) | **GET** /api/payments/connectors/{connector}/config | Read connector config
*PaymentsApi* | [**resetConnector**](docs/Api/PaymentsApi.md#resetconnector) | **POST** /api/payments/connectors/{connector}/reset | Reset connector
*PaymentsApi* | [**uninstallConnector**](docs/Api/PaymentsApi.md#uninstallconnector) | **DELETE** /api/payments/connectors/{connector} | Uninstall connector
*ScopesApi* | [**addTransientScope**](docs/Api/ScopesApi.md#addtransientscope) | **PUT** /api/auth/scopes/{scopeId}/transient/{transientScopeId} | Add a transient scope to a scope
*ScopesApi* | [**createScope**](docs/Api/ScopesApi.md#createscope) | **POST** /api/auth/scopes | Create scope
*ScopesApi* | [**deleteScope**](docs/Api/ScopesApi.md#deletescope) | **DELETE** /api/auth/scopes/{scopeId} | Delete scope
*ScopesApi* | [**deleteTransientScope**](docs/Api/ScopesApi.md#deletetransientscope) | **DELETE** /api/auth/scopes/{scopeId}/transient/{transientScopeId} | Delete a transient scope from a scope
*ScopesApi* | [**listScopes**](docs/Api/ScopesApi.md#listscopes) | **GET** /api/auth/scopes | List scopes
*ScopesApi* | [**readScope**](docs/Api/ScopesApi.md#readscope) | **GET** /api/auth/scopes/{scopeId} | Read scope
*ScopesApi* | [**updateScope**](docs/Api/ScopesApi.md#updatescope) | **PUT** /api/auth/scopes/{scopeId} | Update scope
*ScriptApi* | [**runScript**](docs/Api/ScriptApi.md#runscript) | **POST** /api/ledger/{ledger}/script | Execute a Numscript
*SearchApi* | [**search**](docs/Api/SearchApi.md#search) | **POST** /api/search/ | Search
*ServerApi* | [**getInfo**](docs/Api/ServerApi.md#getinfo) | **GET** /api/ledger/_info | Show server information
*StatsApi* | [**readStats**](docs/Api/StatsApi.md#readstats) | **GET** /api/ledger/{ledger}/stats | Get statistics from a ledger
*TransactionsApi* | [**addMetadataOnTransaction**](docs/Api/TransactionsApi.md#addmetadataontransaction) | **POST** /api/ledger/{ledger}/transactions/{txid}/metadata | Set the metadata of a transaction by its ID
*TransactionsApi* | [**countTransactions**](docs/Api/TransactionsApi.md#counttransactions) | **HEAD** /api/ledger/{ledger}/transactions | Count the transactions from a ledger
*TransactionsApi* | [**createTransaction**](docs/Api/TransactionsApi.md#createtransaction) | **POST** /api/ledger/{ledger}/transactions | Create a new transaction to a ledger
*TransactionsApi* | [**createTransactions**](docs/Api/TransactionsApi.md#createtransactions) | **POST** /api/ledger/{ledger}/transactions/batch | Create a new batch of transactions to a ledger
*TransactionsApi* | [**getTransaction**](docs/Api/TransactionsApi.md#gettransaction) | **GET** /api/ledger/{ledger}/transactions/{txid} | Get transaction from a ledger by its ID
*TransactionsApi* | [**listTransactions**](docs/Api/TransactionsApi.md#listtransactions) | **GET** /api/ledger/{ledger}/transactions | List transactions from a ledger
*TransactionsApi* | [**revertTransaction**](docs/Api/TransactionsApi.md#reverttransaction) | **POST** /api/ledger/{ledger}/transactions/{txid}/revert | Revert a ledger transaction by its ID
*UsersApi* | [**listUsers**](docs/Api/UsersApi.md#listusers) | **GET** /api/auth/users | List users
*UsersApi* | [**readUser**](docs/Api/UsersApi.md#readuser) | **GET** /api/auth/users/{userId} | Read user
*WalletsApi* | [**confirmHold**](docs/Api/WalletsApi.md#confirmhold) | **POST** /api/wallets/holds/{hold_id}/confirm | Confirm a hold
*WalletsApi* | [**createBalance**](docs/Api/WalletsApi.md#createbalance) | **POST** /api/wallets/wallets/{id}/balances | Create a balance
*WalletsApi* | [**createWallet**](docs/Api/WalletsApi.md#createwallet) | **POST** /api/wallets/wallets | Create a new wallet
*WalletsApi* | [**creditWallet**](docs/Api/WalletsApi.md#creditwallet) | **POST** /api/wallets/wallets/{id}/credit | Credit a wallet
*WalletsApi* | [**debitWallet**](docs/Api/WalletsApi.md#debitwallet) | **POST** /api/wallets/wallets/{id}/debit | Debit a wallet
*WalletsApi* | [**getBalance**](docs/Api/WalletsApi.md#getbalance) | **GET** /api/wallets/wallets/{id}/balances/{balanceName} | Get detailed balance
*WalletsApi* | [**getHold**](docs/Api/WalletsApi.md#gethold) | **GET** /api/wallets/holds/{holdID} | Get a hold
*WalletsApi* | [**getHolds**](docs/Api/WalletsApi.md#getholds) | **GET** /api/wallets/holds | Get all holds for a wallet
*WalletsApi* | [**getTransactions**](docs/Api/WalletsApi.md#gettransactions) | **GET** /api/wallets/transactions | 
*WalletsApi* | [**getWallet**](docs/Api/WalletsApi.md#getwallet) | **GET** /api/wallets/wallets/{id} | Get a wallet
*WalletsApi* | [**listBalances**](docs/Api/WalletsApi.md#listbalances) | **GET** /api/wallets/wallets/{id}/balances | List balances of a wallet
*WalletsApi* | [**listWallets**](docs/Api/WalletsApi.md#listwallets) | **GET** /api/wallets/wallets | List all wallets
*WalletsApi* | [**updateWallet**](docs/Api/WalletsApi.md#updatewallet) | **PATCH** /api/wallets/wallets/{id} | Update a wallet
*WalletsApi* | [**voidHold**](docs/Api/WalletsApi.md#voidhold) | **POST** /api/wallets/holds/{hold_id}/void | Cancel a hold
*WalletsApi* | [**walletsgetServerInfo**](docs/Api/WalletsApi.md#walletsgetserverinfo) | **GET** /api/wallets/_info | Get server info
*WebhooksApi* | [**activateConfig**](docs/Api/WebhooksApi.md#activateconfig) | **PUT** /api/webhooks/configs/{id}/activate | Activate one config
*WebhooksApi* | [**changeConfigSecret**](docs/Api/WebhooksApi.md#changeconfigsecret) | **PUT** /api/webhooks/configs/{id}/secret/change | Change the signing secret of a config
*WebhooksApi* | [**deactivateConfig**](docs/Api/WebhooksApi.md#deactivateconfig) | **PUT** /api/webhooks/configs/{id}/deactivate | Deactivate one config
*WebhooksApi* | [**deleteConfig**](docs/Api/WebhooksApi.md#deleteconfig) | **DELETE** /api/webhooks/configs/{id} | Delete one config
*WebhooksApi* | [**getManyConfigs**](docs/Api/WebhooksApi.md#getmanyconfigs) | **GET** /api/webhooks/configs | Get many configs
*WebhooksApi* | [**insertConfig**](docs/Api/WebhooksApi.md#insertconfig) | **POST** /api/webhooks/configs | Insert a new config
*WebhooksApi* | [**testConfig**](docs/Api/WebhooksApi.md#testconfig) | **GET** /api/webhooks/configs/{id}/test | Test one config

## Models

- [Account](docs/Model/Account.md)
- [AccountResponse](docs/Model/AccountResponse.md)
- [AccountWithVolumesAndBalances](docs/Model/AccountWithVolumesAndBalances.md)
- [AccountsCursorResponse](docs/Model/AccountsCursorResponse.md)
- [AccountsCursorResponseCursor](docs/Model/AccountsCursorResponseCursor.md)
- [AggregateBalancesResponse](docs/Model/AggregateBalancesResponse.md)
- [AssetHolder](docs/Model/AssetHolder.md)
- [Attempt](docs/Model/Attempt.md)
- [AttemptResponse](docs/Model/AttemptResponse.md)
- [Balance](docs/Model/Balance.md)
- [BalanceWithAssets](docs/Model/BalanceWithAssets.md)
- [BalancesCursorResponse](docs/Model/BalancesCursorResponse.md)
- [BalancesCursorResponseCursor](docs/Model/BalancesCursorResponseCursor.md)
- [BankingCircleConfig](docs/Model/BankingCircleConfig.md)
- [Client](docs/Model/Client.md)
- [ClientAllOf](docs/Model/ClientAllOf.md)
- [ClientOptions](docs/Model/ClientOptions.md)
- [ClientSecret](docs/Model/ClientSecret.md)
- [Config](docs/Model/Config.md)
- [ConfigChangeSecret](docs/Model/ConfigChangeSecret.md)
- [ConfigInfo](docs/Model/ConfigInfo.md)
- [ConfigInfoResponse](docs/Model/ConfigInfoResponse.md)
- [ConfigResponse](docs/Model/ConfigResponse.md)
- [ConfigUser](docs/Model/ConfigUser.md)
- [ConfigsResponse](docs/Model/ConfigsResponse.md)
- [ConfigsResponseCursor](docs/Model/ConfigsResponseCursor.md)
- [ConfigsResponseCursorAllOf](docs/Model/ConfigsResponseCursorAllOf.md)
- [ConfirmHoldRequest](docs/Model/ConfirmHoldRequest.md)
- [ConnectorBaseInfo](docs/Model/ConnectorBaseInfo.md)
- [ConnectorConfig](docs/Model/ConnectorConfig.md)
- [Connectors](docs/Model/Connectors.md)
- [Contract](docs/Model/Contract.md)
- [CreateBalanceResponse](docs/Model/CreateBalanceResponse.md)
- [CreateClientResponse](docs/Model/CreateClientResponse.md)
- [CreateScopeResponse](docs/Model/CreateScopeResponse.md)
- [CreateSecretResponse](docs/Model/CreateSecretResponse.md)
- [CreateWalletRequest](docs/Model/CreateWalletRequest.md)
- [CreateWalletResponse](docs/Model/CreateWalletResponse.md)
- [CreditWalletRequest](docs/Model/CreditWalletRequest.md)
- [CurrencyCloudConfig](docs/Model/CurrencyCloudConfig.md)
- [Cursor](docs/Model/Cursor.md)
- [DebitWalletRequest](docs/Model/DebitWalletRequest.md)
- [DebitWalletResponse](docs/Model/DebitWalletResponse.md)
- [DummyPayConfig](docs/Model/DummyPayConfig.md)
- [ErrorResponse](docs/Model/ErrorResponse.md)
- [ErrorsEnum](docs/Model/ErrorsEnum.md)
- [ExpandedDebitHold](docs/Model/ExpandedDebitHold.md)
- [ExpandedDebitHoldAllOf](docs/Model/ExpandedDebitHoldAllOf.md)
- [GetBalanceResponse](docs/Model/GetBalanceResponse.md)
- [GetHoldResponse](docs/Model/GetHoldResponse.md)
- [GetHoldsResponse](docs/Model/GetHoldsResponse.md)
- [GetHoldsResponseCursor](docs/Model/GetHoldsResponseCursor.md)
- [GetHoldsResponseCursorAllOf](docs/Model/GetHoldsResponseCursorAllOf.md)
- [GetPaymentResponse](docs/Model/GetPaymentResponse.md)
- [GetTransactionsResponse](docs/Model/GetTransactionsResponse.md)
- [GetTransactionsResponseCursor](docs/Model/GetTransactionsResponseCursor.md)
- [GetTransactionsResponseCursorAllOf](docs/Model/GetTransactionsResponseCursorAllOf.md)
- [GetWalletResponse](docs/Model/GetWalletResponse.md)
- [Hold](docs/Model/Hold.md)
- [LedgerAccountSubject](docs/Model/LedgerAccountSubject.md)
- [LedgerInfo](docs/Model/LedgerInfo.md)
- [LedgerInfoResponse](docs/Model/LedgerInfoResponse.md)
- [LedgerInfoStorage](docs/Model/LedgerInfoStorage.md)
- [LedgerStorage](docs/Model/LedgerStorage.md)
- [ListAccountsResponse](docs/Model/ListAccountsResponse.md)
- [ListBalancesResponse](docs/Model/ListBalancesResponse.md)
- [ListBalancesResponseCursor](docs/Model/ListBalancesResponseCursor.md)
- [ListBalancesResponseCursorAllOf](docs/Model/ListBalancesResponseCursorAllOf.md)
- [ListClientsResponse](docs/Model/ListClientsResponse.md)
- [ListConnectorTasks200ResponseInner](docs/Model/ListConnectorTasks200ResponseInner.md)
- [ListConnectorsConfigsResponse](docs/Model/ListConnectorsConfigsResponse.md)
- [ListConnectorsConfigsResponseConnector](docs/Model/ListConnectorsConfigsResponseConnector.md)
- [ListConnectorsConfigsResponseConnectorKey](docs/Model/ListConnectorsConfigsResponseConnectorKey.md)
- [ListConnectorsResponse](docs/Model/ListConnectorsResponse.md)
- [ListPaymentsResponse](docs/Model/ListPaymentsResponse.md)
- [ListScopesResponse](docs/Model/ListScopesResponse.md)
- [ListUsersResponse](docs/Model/ListUsersResponse.md)
- [ListWalletsResponse](docs/Model/ListWalletsResponse.md)
- [ListWalletsResponseCursor](docs/Model/ListWalletsResponseCursor.md)
- [ListWalletsResponseCursorAllOf](docs/Model/ListWalletsResponseCursorAllOf.md)
- [Log](docs/Model/Log.md)
- [LogsCursorResponse](docs/Model/LogsCursorResponse.md)
- [LogsCursorResponseCursor](docs/Model/LogsCursorResponseCursor.md)
- [Mapping](docs/Model/Mapping.md)
- [MappingResponse](docs/Model/MappingResponse.md)
- [MigrationInfo](docs/Model/MigrationInfo.md)
- [ModulrConfig](docs/Model/ModulrConfig.md)
- [Monetary](docs/Model/Monetary.md)
- [Payment](docs/Model/Payment.md)
- [PaymentsAccount](docs/Model/PaymentsAccount.md)
- [PostTransaction](docs/Model/PostTransaction.md)
- [PostTransactionScript](docs/Model/PostTransactionScript.md)
- [Posting](docs/Model/Posting.md)
- [Query](docs/Model/Query.md)
- [ReadClientResponse](docs/Model/ReadClientResponse.md)
- [ReadUserResponse](docs/Model/ReadUserResponse.md)
- [Response](docs/Model/Response.md)
- [Scope](docs/Model/Scope.md)
- [ScopeAllOf](docs/Model/ScopeAllOf.md)
- [ScopeOptions](docs/Model/ScopeOptions.md)
- [Script](docs/Model/Script.md)
- [ScriptResponse](docs/Model/ScriptResponse.md)
- [Secret](docs/Model/Secret.md)
- [SecretAllOf](docs/Model/SecretAllOf.md)
- [SecretOptions](docs/Model/SecretOptions.md)
- [ServerInfo](docs/Model/ServerInfo.md)
- [Stats](docs/Model/Stats.md)
- [StatsResponse](docs/Model/StatsResponse.md)
- [StripeConfig](docs/Model/StripeConfig.md)
- [StripeTask](docs/Model/StripeTask.md)
- [StripeTransferRequest](docs/Model/StripeTransferRequest.md)
- [Subject](docs/Model/Subject.md)
- [TaskDescriptorBankingCircle](docs/Model/TaskDescriptorBankingCircle.md)
- [TaskDescriptorBankingCircleDescriptor](docs/Model/TaskDescriptorBankingCircleDescriptor.md)
- [TaskDescriptorCurrencyCloud](docs/Model/TaskDescriptorCurrencyCloud.md)
- [TaskDescriptorCurrencyCloudDescriptor](docs/Model/TaskDescriptorCurrencyCloudDescriptor.md)
- [TaskDescriptorDummyPay](docs/Model/TaskDescriptorDummyPay.md)
- [TaskDescriptorDummyPayDescriptor](docs/Model/TaskDescriptorDummyPayDescriptor.md)
- [TaskDescriptorModulr](docs/Model/TaskDescriptorModulr.md)
- [TaskDescriptorModulrDescriptor](docs/Model/TaskDescriptorModulrDescriptor.md)
- [TaskDescriptorStripe](docs/Model/TaskDescriptorStripe.md)
- [TaskDescriptorStripeDescriptor](docs/Model/TaskDescriptorStripeDescriptor.md)
- [TaskDescriptorWise](docs/Model/TaskDescriptorWise.md)
- [TaskDescriptorWiseDescriptor](docs/Model/TaskDescriptorWiseDescriptor.md)
- [Total](docs/Model/Total.md)
- [Transaction](docs/Model/Transaction.md)
- [TransactionData](docs/Model/TransactionData.md)
- [TransactionResponse](docs/Model/TransactionResponse.md)
- [Transactions](docs/Model/Transactions.md)
- [TransactionsCursorResponse](docs/Model/TransactionsCursorResponse.md)
- [TransactionsCursorResponseCursor](docs/Model/TransactionsCursorResponseCursor.md)
- [TransactionsResponse](docs/Model/TransactionsResponse.md)
- [UpdateWalletRequest](docs/Model/UpdateWalletRequest.md)
- [User](docs/Model/User.md)
- [Volume](docs/Model/Volume.md)
- [Wallet](docs/Model/Wallet.md)
- [WalletSubject](docs/Model/WalletSubject.md)
- [WalletWithBalances](docs/Model/WalletWithBalances.md)
- [WalletWithBalancesBalances](docs/Model/WalletWithBalancesBalances.md)
- [WalletsErrorResponse](docs/Model/WalletsErrorResponse.md)
- [WalletsTransaction](docs/Model/WalletsTransaction.md)
- [WalletsVolume](docs/Model/WalletsVolume.md)
- [WebhooksConfig](docs/Model/WebhooksConfig.md)
- [WiseConfig](docs/Model/WiseConfig.md)

## Authorization

### Authorization

- **Type**: `OAuth`
- **Flow**: `application`
- **Authorization URL**: ``
- **Scopes**: N/A

## Tests

To run the tests, use:

```bash
composer install
vendor/bin/phpunit
```

## Author

support@formance.com

## About this package

This PHP package is automatically generated by the [OpenAPI Generator](https://openapi-generator.tech) project:

- API version: `develop`
- Build package: `org.openapitools.codegen.languages.PhpClientCodegen`
