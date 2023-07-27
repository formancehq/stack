# payments

### Available Operations

* [connectorsStripeTransfer](#connectorsstripetransfer) - Transfer funds between Stripe accounts
* [connectorsTransfer](#connectorstransfer) - Transfer funds between Connector accounts
* [getAccountBalances](#getaccountbalances) - Get account balances
* [getConnectorTask](#getconnectortask) - Read a specific task of the connector
* [getPayment](#getpayment) - Get a payment
* [installConnector](#installconnector) - Install a connector
* [listAllConnectors](#listallconnectors) - List all installed connectors
* [listConfigsAvailableConnectors](#listconfigsavailableconnectors) - List the configs of each available connector
* [listConnectorTasks](#listconnectortasks) - List tasks from a connector
* [listConnectorsTransfers](#listconnectorstransfers) - List transfers and their statuses
* [listPayments](#listpayments) - List payments
* [paymentsgetServerInfo](#paymentsgetserverinfo) - Get server info
* [paymentslistAccounts](#paymentslistaccounts) - List accounts
* [readConnectorConfig](#readconnectorconfig) - Read the config of a connector
* [resetConnector](#resetconnector) - Reset a connector
* [uninstallConnector](#uninstallconnector) - Uninstall a connector
* [updateMetadata](#updatemetadata) - Update metadata

## connectorsStripeTransfer

Execute a transfer between two Stripe accounts.

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { ConnectorsStripeTransferResponse } from "@formance/formance-sdk/dist/sdk/models/operations";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.payments.connectorsStripeTransfer({
  amount: 100,
  asset: "USD",
  destination: "acct_1Gqj58KZcSIg2N2q",
  metadata: {
    "labore": "delectus",
    "eum": "non",
    "eligendi": "sint",
  },
}).then((res: ConnectorsStripeTransferResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

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
  connector: Connector.Modulr,
}).then((res: ConnectorsTransferResponse) => {
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
  accountID: "provident",
  currency: "necessitatibus",
  cursor: "aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ==",
  from: new Date("2021-09-21T14:06:09.271Z"),
  limit: 223081,
  pageSize: 891555,
  sort: [
    "dolorum",
    "in",
    "in",
    "illum",
  ],
  to: new Date("2020-11-26T01:41:04.216Z"),
}).then((res: GetAccountBalancesResponse) => {
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
  connector: Connector.Stripe,
  taskId: "magnam",
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
  paymentId: "cumque",
}).then((res: GetPaymentResponse) => {
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
    apiKey: "XXX",
    clientID: "XXX",
    endpoint: "XXX",
    pollingPeriod: "60s",
  },
  connector: Connector.Modulr,
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
  pageSize: 675439,
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
  pageSize: 249796,
  sort: [
    "enim",
    "accusamus",
    "delectus",
  ],
}).then((res: ListPaymentsResponse) => {
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
  pageSize: 692532,
  sort: [
    "nam",
    "id",
    "blanditiis",
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
  connector: Connector.CurrencyCloud,
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
  connector: Connector.Moneycorp,
}).then((res: ResetConnectorResponse) => {
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
  connector: Connector.DummyPay,
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
    key: "deserunt",
  },
  paymentId: "nisi",
}).then((res: UpdateMetadataResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```
