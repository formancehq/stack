# payments

### Available Operations

* [connectorsTransfer](#connectorstransfer) - Transfer funds between Connector accounts
* [createBankAccount](#createbankaccount) - Create a BankAccount in Payments and on the PSP
* [createTransferInitiation](#createtransferinitiation) - Create a TransferInitiation
* [deleteTransferInitiation](#deletetransferinitiation) - Delete a transfer initiation
* [getAccountBalances](#getaccountbalances) - Get account balances
* [getBankAccount](#getbankaccount) - Get a bank account created by user on Formance
* [getConnectorTask](#getconnectortask) - Read a specific task of the connector
* [getPayment](#getpayment) - Get a payment
* [getTransferInitiation](#gettransferinitiation) - Get a transfer initiation
* [installConnector](#installconnector) - Install a connector
* [listAllConnectors](#listallconnectors) - List all installed connectors
* [listBankAccounts](#listbankaccounts) - List bank accounts created by user on Formance
* [listConfigsAvailableConnectors](#listconfigsavailableconnectors) - List the configs of each available connector
* [listConnectorTasks](#listconnectortasks) - List tasks from a connector
* [listConnectorsTransfers](#listconnectorstransfers) - List transfers and their statuses
* [listPayments](#listpayments) - List payments
* [listTransferInitiations](#listtransferinitiations) - List Transfer Initiations
* [paymentsgetAccount](#paymentsgetaccount) - Get an account
* [paymentsgetServerInfo](#paymentsgetserverinfo) - Get server info
* [paymentslistAccounts](#paymentslistaccounts) - List accounts
* [readConnectorConfig](#readconnectorconfig) - Read the config of a connector
* [resetConnector](#resetconnector) - Reset a connector
* [retryTransferInitiation](#retrytransferinitiation) - Retry a failed transfer initiation
* [udpateTransferInitiationStatus](#udpatetransferinitiationstatus) - Update the status of a transfer initiation
* [uninstallConnector](#uninstallconnector) - Uninstall a connector
* [updateMetadata](#updatemetadata) - Update metadata

## connectorsTransfer

Execute a transfer between two accounts.

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { ConnectorsTransferResponse } from "@formance/formance-sdk/dist/sdk/models/operations";
import { Connector } from "@formance/formance-sdk/dist/sdk/models/shared";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.payments.connectorsTransfer({
  transferRequest: {
    amount: 100,
    asset: "USD",
    destination: "acct_1Gqj58KZcSIg2N2q",
    source: "acct_1Gqj58KZcSIg2N2q",
  },
  connector: Connector.CurrencyCloud,
}).then((res: ConnectorsTransferResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## createBankAccount

Create a bank account in Payments and on the PSP.

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { CreateBankAccountResponse } from "@formance/formance-sdk/dist/sdk/models/operations";
import { Connector } from "@formance/formance-sdk/dist/sdk/models/shared";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.payments.createBankAccount({
  accountNumber: "voluptates",
  country: "GB",
  iban: "quasi",
  name: "My account",
  provider: Connector.Moneycorp,
  swiftBicCode: "sint",
}).then((res: CreateBankAccountResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## createTransferInitiation

Create a transfer initiation

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { CreateTransferInitiationResponse } from "@formance/formance-sdk/dist/sdk/models/operations";
import {
  Connector,
  PaymentStatus,
  TransferInitiationRequestType,
  TransferInitiationStatus,
  TransferInitiationType,
} from "@formance/formance-sdk/dist/sdk/models/shared";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.payments.createTransferInitiation({
  amount: 83112,
  asset: "USD",
  createdAt: new Date("2022-03-02T21:33:21.372Z"),
  description: "enim",
  destinationAccountID: "consequatur",
  provider: Connector.BankingCircle,
  reference: "XXX",
  sourceAccountID: "quibusdam",
  type: TransferInitiationRequestType.Transfer,
  validated: false,
}).then((res: CreateTransferInitiationResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## deleteTransferInitiation

Delete a transfer initiation by its id.

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { DeleteTransferInitiationResponse } from "@formance/formance-sdk/dist/sdk/models/operations";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.payments.deleteTransferInitiation({
  transferId: "deserunt",
}).then((res: DeleteTransferInitiationResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## getAccountBalances

Get account balances

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { GetAccountBalancesResponse } from "@formance/formance-sdk/dist/sdk/models/operations";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.payments.getAccountBalances({
  accountId: "distinctio",
  asset: "quibusdam",
  cursor: "aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ==",
  from: new Date("2022-09-26T08:57:48.803Z"),
  limit: 183191,
  pageSize: 397821,
  sort: [
    "quos",
    "perferendis",
    "magni",
  ],
  to: new Date("2021-11-22T01:26:35.048Z"),
}).then((res: GetAccountBalancesResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## getBankAccount

Get a bank account created by user on Formance

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { GetBankAccountResponse } from "@formance/formance-sdk/dist/sdk/models/operations";
import { Connector } from "@formance/formance-sdk/dist/sdk/models/shared";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.payments.getBankAccount({
  bankAccountId: "alias",
}).then((res: GetBankAccountResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## getConnectorTask

Get a specific task associated to the connector.

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { GetConnectorTaskResponse } from "@formance/formance-sdk/dist/sdk/models/operations";
import { Connector } from "@formance/formance-sdk/dist/sdk/models/shared";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.payments.getConnectorTask({
  connector: Connector.DummyPay,
  taskId: "dolorum",
}).then((res: GetConnectorTaskResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## getPayment

Get a payment

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { GetPaymentResponse } from "@formance/formance-sdk/dist/sdk/models/operations";
import { Connector, PaymentScheme, PaymentStatus, PaymentType } from "@formance/formance-sdk/dist/sdk/models/shared";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.payments.getPayment({
  paymentId: "excepturi",
}).then((res: GetPaymentResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## getTransferInitiation

Get a transfer initiation

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { GetTransferInitiationResponse } from "@formance/formance-sdk/dist/sdk/models/operations";
import { Connector, PaymentStatus, TransferInitiationStatus, TransferInitiationType } from "@formance/formance-sdk/dist/sdk/models/shared";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.payments.getTransferInitiation({
  transferId: "tempora",
}).then((res: GetTransferInitiationResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## installConnector

Install a connector by its name and config.

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { InstallConnectorResponse } from "@formance/formance-sdk/dist/sdk/models/operations";
import { Connector } from "@formance/formance-sdk/dist/sdk/models/shared";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.payments.installConnector({
  requestBody: {
    authorizationEndpoint: "XXX",
    endpoint: "XXX",
    password: "XXX",
    pollingPeriod: "60s",
    userCertificate: "XXX",
    userCertificateKey: "XXX",
    username: "XXX",
  },
  connector: Connector.BankingCircle,
}).then((res: InstallConnectorResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## listAllConnectors

List all installed connectors.

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { ListAllConnectorsResponse } from "@formance/formance-sdk/dist/sdk/models/operations";
import { Connector } from "@formance/formance-sdk/dist/sdk/models/shared";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.payments.listAllConnectors().then((res: ListAllConnectorsResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## listBankAccounts

List all bank accounts created by user on Formance.

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { ListBankAccountsResponse } from "@formance/formance-sdk/dist/sdk/models/operations";
import { Connector } from "@formance/formance-sdk/dist/sdk/models/shared";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.payments.listBankAccounts({
  cursor: "aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ==",
  pageSize: 288476,
  sort: [
    "eum",
    "non",
    "eligendi",
    "sint",
  ],
}).then((res: ListBankAccountsResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## listConfigsAvailableConnectors

List the configs of each available connector.

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { ListConfigsAvailableConnectorsResponse } from "@formance/formance-sdk/dist/sdk/models/operations";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.payments.listConfigsAvailableConnectors().then((res: ListConfigsAvailableConnectorsResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## listConnectorTasks

List all tasks associated with this connector.

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { ListConnectorTasksResponse } from "@formance/formance-sdk/dist/sdk/models/operations";
import { Connector } from "@formance/formance-sdk/dist/sdk/models/shared";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.payments.listConnectorTasks({
  connector: Connector.Modulr,
  cursor: "aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ==",
  pageSize: 592042,
}).then((res: ListConnectorTasksResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## listConnectorsTransfers

List transfers

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { ListConnectorsTransfersResponse } from "@formance/formance-sdk/dist/sdk/models/operations";
import { Connector } from "@formance/formance-sdk/dist/sdk/models/shared";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.payments.listConnectorsTransfers({
  connector: Connector.Moneycorp,
}).then((res: ListConnectorsTransfersResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## listPayments

List payments

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { ListPaymentsResponse } from "@formance/formance-sdk/dist/sdk/models/operations";
import { Connector, PaymentScheme, PaymentStatus, PaymentType } from "@formance/formance-sdk/dist/sdk/models/shared";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.payments.listPayments({
  cursor: "aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ==",
  pageSize: 572252,
  sort: [
    "dolor",
    "debitis",
    "a",
  ],
}).then((res: ListPaymentsResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## listTransferInitiations

List Transfer Initiations

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { ListTransferInitiationsResponse } from "@formance/formance-sdk/dist/sdk/models/operations";
import { Connector, PaymentStatus, TransferInitiationStatus, TransferInitiationType } from "@formance/formance-sdk/dist/sdk/models/shared";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.payments.listTransferInitiations({
  cursor: "aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ==",
  pageSize: 680056,
  query: "in",
  sort: [
    "illum",
    "maiores",
  ],
}).then((res: ListTransferInitiationsResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## paymentsgetAccount

Get an account

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { PaymentsgetAccountResponse } from "@formance/formance-sdk/dist/sdk/models/operations";
import { Connector } from "@formance/formance-sdk/dist/sdk/models/shared";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.payments.paymentsgetAccount({
  accountId: "rerum",
}).then((res: PaymentsgetAccountResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## paymentsgetServerInfo

Get server info

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { PaymentsgetServerInfoResponse } from "@formance/formance-sdk/dist/sdk/models/operations";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.payments.paymentsgetServerInfo().then((res: PaymentsgetServerInfoResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## paymentslistAccounts

List accounts

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { PaymentslistAccountsResponse } from "@formance/formance-sdk/dist/sdk/models/operations";
import { Connector } from "@formance/formance-sdk/dist/sdk/models/shared";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.payments.paymentslistAccounts({
  cursor: "aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ==",
  pageSize: 116202,
  sort: [
    "cumque",
    "facere",
  ],
}).then((res: PaymentslistAccountsResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## readConnectorConfig

Read connector config

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { ReadConnectorConfigResponse } from "@formance/formance-sdk/dist/sdk/models/operations";
import { Connector } from "@formance/formance-sdk/dist/sdk/models/shared";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.payments.readConnectorConfig({
  connector: Connector.Modulr,
}).then((res: ReadConnectorConfigResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## resetConnector

Reset a connector by its name.
It will remove the connector and ALL PAYMENTS generated with it.


### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { ResetConnectorResponse } from "@formance/formance-sdk/dist/sdk/models/operations";
import { Connector } from "@formance/formance-sdk/dist/sdk/models/shared";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.payments.resetConnector({
  connector: Connector.Modulr,
}).then((res: ResetConnectorResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## retryTransferInitiation

Retry a failed transfer initiation

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { RetryTransferInitiationResponse } from "@formance/formance-sdk/dist/sdk/models/operations";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.payments.retryTransferInitiation({
  transferId: "laborum",
}).then((res: RetryTransferInitiationResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## udpateTransferInitiationStatus

Update a transfer initiation status

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { UdpateTransferInitiationStatusResponse } from "@formance/formance-sdk/dist/sdk/models/operations";
import { UpdateTransferInitiationStatusRequestStatus } from "@formance/formance-sdk/dist/sdk/models/shared";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.payments.udpateTransferInitiationStatus({
  updateTransferInitiationStatusRequest: {
    status: UpdateTransferInitiationStatusRequestStatus.Validated,
  },
  transferId: "non",
}).then((res: UdpateTransferInitiationStatusResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## uninstallConnector

Uninstall a connector by its name.

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { UninstallConnectorResponse } from "@formance/formance-sdk/dist/sdk/models/operations";
import { Connector } from "@formance/formance-sdk/dist/sdk/models/shared";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.payments.uninstallConnector({
  connector: Connector.CurrencyCloud,
}).then((res: UninstallConnectorResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## updateMetadata

Update metadata

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { UpdateMetadataResponse } from "@formance/formance-sdk/dist/sdk/models/operations";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.payments.updateMetadata({
  paymentMetadata: {
    key: "enim",
  },
  paymentId: "accusamus",
}).then((res: UpdateMetadataResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```
