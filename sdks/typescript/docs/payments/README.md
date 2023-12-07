# payments

### Available Operations

* [connectorsTransfer](#connectorstransfer) - Transfer funds between Connector accounts
* [createBankAccount](#createbankaccount) - Create a BankAccount in Payments and on the PSP
* [createTransferInitiation](#createtransferinitiation) - Create a TransferInitiation
* [deleteTransferInitiation](#deletetransferinitiation) - Delete a transfer initiation
* [getAccountBalances](#getaccountbalances) - Get account balances
* [getBankAccount](#getbankaccount) - Get a bank account created by user on Formance
* [~~getConnectorTask~~](#getconnectortask) - Read a specific task of the connector :warning: **Deprecated**
* [getConnectorTaskV1](#getconnectortaskv1) - Read a specific task of the connector
* [getPayment](#getpayment) - Get a payment
* [getTransferInitiation](#gettransferinitiation) - Get a transfer initiation
* [installConnector](#installconnector) - Install a connector
* [listAllConnectors](#listallconnectors) - List all installed connectors
* [listBankAccounts](#listbankaccounts) - List bank accounts created by user on Formance
* [listConfigsAvailableConnectors](#listconfigsavailableconnectors) - List the configs of each available connector
* [~~listConnectorTasks~~](#listconnectortasks) - List tasks from a connector :warning: **Deprecated**
* [listConnectorTasksV1](#listconnectortasksv1) - List tasks from a connector
* [listPayments](#listpayments) - List payments
* [listTransferInitiations](#listtransferinitiations) - List Transfer Initiations
* [paymentsgetAccount](#paymentsgetaccount) - Get an account
* [paymentsgetServerInfo](#paymentsgetserverinfo) - Get server info
* [paymentslistAccounts](#paymentslistaccounts) - List accounts
* [~~readConnectorConfig~~](#readconnectorconfig) - Read the config of a connector :warning: **Deprecated**
* [readConnectorConfigV1](#readconnectorconfigv1) - Read the config of a connector
* [~~resetConnector~~](#resetconnector) - Reset a connector :warning: **Deprecated**
* [resetConnectorV1](#resetconnectorv1) - Reset a connector
* [retryTransferInitiation](#retrytransferinitiation) - Retry a failed transfer initiation
* [udpateTransferInitiationStatus](#udpatetransferinitiationstatus) - Update the status of a transfer initiation
* [~~uninstallConnector~~](#uninstallconnector) - Uninstall a connector :warning: **Deprecated**
* [uninstallConnectorV1](#uninstallconnectorv1) - Uninstall a connector
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
  connector: Connector.Wise,
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

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.payments.createBankAccount({
  accountNumber: "consequatur",
  connectorID: "est",
  country: "GB",
  iban: "quibusdam",
  name: "My account",
  swiftBicCode: "explicabo",
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
  amount: 647174,
  asset: "USD",
  connectorID: "distinctio",
  description: "quibusdam",
  destinationAccountID: "labore",
  provider: Connector.Wise,
  reference: "XXX",
  scheduledAt: new Date("2022-08-08T19:05:24.174Z"),
  sourceAccountID: "cupiditate",
  type: TransferInitiationRequestType.Payout,
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
  transferId: "perferendis",
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
  accountId: "magni",
  asset: "assumenda",
  cursor: "aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ==",
  from: new Date("2022-12-30T06:52:02.282Z"),
  limit: 146441,
  pageSize: 677817,
  sort: [
    "tempora",
    "facilis",
    "tempore",
  ],
  to: new Date("2022-01-14T19:13:42.009Z"),
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

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.payments.getBankAccount({
  bankAccountId: "eum",
}).then((res: GetBankAccountResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## ~~getConnectorTask~~

Get a specific task associated to the connector.

> :warning: **DEPRECATED**: this method will be removed in a future release, please migrate away from it as soon as possible.

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
  taskId: "eligendi",
}).then((res: GetConnectorTaskResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## getConnectorTaskV1

Get a specific task associated to the connector.

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { GetConnectorTaskV1Response } from "@formance/formance-sdk/dist/sdk/models/operations";
import { Connector } from "@formance/formance-sdk/dist/sdk/models/shared";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.payments.getConnectorTaskV1({
  connector: Connector.CurrencyCloud,
  connectorId: "aliquid",
  taskId: "provident",
}).then((res: GetConnectorTaskV1Response) => {
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
  paymentId: "necessitatibus",
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
import { PaymentStatus, TransferInitiationStatus, TransferInitiationType } from "@formance/formance-sdk/dist/sdk/models/shared";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.payments.getTransferInitiation({
  transferId: "sint",
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
    name: "My Banking Circle Account",
    password: "XXX",
    pollingPeriod: "60s",
    userCertificate: "XXX",
    userCertificateKey: "XXX",
    username: "XXX",
  },
  connector: Connector.DummyPay,
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

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.payments.listBankAccounts({
  cursor: "aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ==",
  pageSize: 891555,
  sort: [
    "dolorum",
    "in",
    "in",
    "illum",
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

## ~~listConnectorTasks~~

List all tasks associated with this connector.

> :warning: **DEPRECATED**: this method will be removed in a future release, please migrate away from it as soon as possible.

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
  connector: Connector.Moneycorp,
  cursor: "aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ==",
  pageSize: 699479,
}).then((res: ListConnectorTasksResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## listConnectorTasksV1

List all tasks associated with this connector.

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { ListConnectorTasksV1Response } from "@formance/formance-sdk/dist/sdk/models/operations";
import { Connector } from "@formance/formance-sdk/dist/sdk/models/shared";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.payments.listConnectorTasksV1({
  connector: Connector.Stripe,
  connectorId: "magnam",
  cursor: "aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ==",
  pageSize: 767024,
}).then((res: ListConnectorTasksV1Response) => {
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
  pageSize: 813798,
  sort: [
    "aliquid",
    "laborum",
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
import { PaymentStatus, TransferInitiationStatus, TransferInitiationType } from "@formance/formance-sdk/dist/sdk/models/shared";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.payments.listTransferInitiations({
  cursor: "aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ==",
  pageSize: 881104,
  query: "non",
  sort: [
    "enim",
    "accusamus",
    "delectus",
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

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.payments.paymentsgetAccount({
  accountId: "quidem",
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

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.payments.paymentslistAccounts({
  cursor: "aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ==",
  pageSize: 588465,
  sort: [
    "id",
    "blanditiis",
    "deleniti",
  ],
}).then((res: PaymentslistAccountsResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## ~~readConnectorConfig~~

Read connector config

> :warning: **DEPRECATED**: this method will be removed in a future release, please migrate away from it as soon as possible.

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
  connector: Connector.Moneycorp,
}).then((res: ReadConnectorConfigResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## readConnectorConfigV1

Read connector config

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { ReadConnectorConfigV1Response } from "@formance/formance-sdk/dist/sdk/models/operations";
import { Connector } from "@formance/formance-sdk/dist/sdk/models/shared";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.payments.readConnectorConfigV1({
  connector: Connector.DummyPay,
  connectorId: "deserunt",
}).then((res: ReadConnectorConfigV1Response) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## ~~resetConnector~~

Reset a connector by its name.
It will remove the connector and ALL PAYMENTS generated with it.


> :warning: **DEPRECATED**: this method will be removed in a future release, please migrate away from it as soon as possible.

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

## resetConnectorV1

Reset a connector by its name.
It will remove the connector and ALL PAYMENTS generated with it.


### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { ResetConnectorV1Response } from "@formance/formance-sdk/dist/sdk/models/operations";
import { Connector } from "@formance/formance-sdk/dist/sdk/models/shared";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.payments.resetConnectorV1({
  connector: Connector.Modulr,
  connectorId: "natus",
}).then((res: ResetConnectorV1Response) => {
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
  transferId: "omnis",
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
    status: UpdateTransferInitiationStatusRequestStatus.Processed,
  },
  transferId: "perferendis",
}).then((res: UdpateTransferInitiationStatusResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## ~~uninstallConnector~~

Uninstall a connector by its name.

> :warning: **DEPRECATED**: this method will be removed in a future release, please migrate away from it as soon as possible.

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
  connector: Connector.Modulr,
}).then((res: UninstallConnectorResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## uninstallConnectorV1

Uninstall a connector by its name.

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { UninstallConnectorV1Response } from "@formance/formance-sdk/dist/sdk/models/operations";
import { Connector } from "@formance/formance-sdk/dist/sdk/models/shared";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.payments.uninstallConnectorV1({
  connector: Connector.Wise,
  connectorId: "distinctio",
}).then((res: UninstallConnectorV1Response) => {
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
    key: "id",
  },
  paymentId: "labore",
}).then((res: UpdateMetadataResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```
