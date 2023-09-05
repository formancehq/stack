# transactions

### Available Operations

* [createTransactions](#createtransactions) - Create a new batch of transactions to a ledger
* [addMetadataOnTransaction](#addmetadataontransaction) - Set the metadata of a transaction by its ID
* [countTransactions](#counttransactions) - Count the transactions from a ledger
* [createTransaction](#createtransaction) - Create a new transaction to a ledger
* [getTransaction](#gettransaction) - Get transaction from a ledger by its ID
* [listTransactions](#listtransactions) - List transactions from a ledger
* [revertTransaction](#reverttransaction) - Revert a ledger transaction by its ID

## createTransactions

Create a new batch of transactions to a ledger

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { CreateTransactionsResponse } from "@formance/formance-sdk/dist/sdk/models/operations";

const sdk = new SDK({
  security: {
    authorization: "",
  },
});

sdk.transactions.createTransactions({
  transactions: {
    transactions: [
      {
        metadata: {
          "deserunt": "distinctio",
        },
        postings: [
          {
            amount: 100,
            asset: "COIN",
            destination: "users:002",
            source: "users:001",
          },
          {
            amount: 100,
            asset: "COIN",
            destination: "users:002",
            source: "users:001",
          },
          {
            amount: 100,
            asset: "COIN",
            destination: "users:002",
            source: "users:001",
          },
          {
            amount: 100,
            asset: "COIN",
            destination: "users:002",
            source: "users:001",
          },
        ],
        reference: "ref:001",
        timestamp: new Date("2022-09-26T08:57:48.803Z"),
      },
      {
        metadata: {
          "aliquid": "cupiditate",
        },
        postings: [
          {
            amount: 100,
            asset: "COIN",
            destination: "users:002",
            source: "users:001",
          },
          {
            amount: 100,
            asset: "COIN",
            destination: "users:002",
            source: "users:001",
          },
          {
            amount: 100,
            asset: "COIN",
            destination: "users:002",
            source: "users:001",
          },
        ],
        reference: "ref:001",
        timestamp: new Date("2022-11-01T19:07:16.800Z"),
      },
      {
        metadata: {
          "ipsam": "alias",
          "fugit": "dolorum",
          "excepturi": "tempora",
          "facilis": "tempore",
        },
        postings: [
          {
            amount: 100,
            asset: "COIN",
            destination: "users:002",
            source: "users:001",
          },
          {
            amount: 100,
            asset: "COIN",
            destination: "users:002",
            source: "users:001",
          },
        ],
        reference: "ref:001",
        timestamp: new Date("2021-09-13T13:11:22.288Z"),
      },
      {
        metadata: {
          "eligendi": "sint",
        },
        postings: [
          {
            amount: 100,
            asset: "COIN",
            destination: "users:002",
            source: "users:001",
          },
          {
            amount: 100,
            asset: "COIN",
            destination: "users:002",
            source: "users:001",
          },
        ],
        reference: "ref:001",
        timestamp: new Date("2021-03-17T21:24:26.606Z"),
      },
    ],
  },
  ledger: "ledger001",
}).then((res: CreateTransactionsResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

### Parameters

| Parameter                                                                                    | Type                                                                                         | Required                                                                                     | Description                                                                                  |
| -------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------- |
| `request`                                                                                    | [operations.CreateTransactionsRequest](../../models/operations/createtransactionsrequest.md) | :heavy_check_mark:                                                                           | The request object to use for the request.                                                   |
| `config`                                                                                     | [AxiosRequestConfig](https://axios-http.com/docs/req_config)                                 | :heavy_minus_sign:                                                                           | Available config options for making requests.                                                |


### Response

**Promise<[operations.CreateTransactionsResponse](../../models/operations/createtransactionsresponse.md)>**


## addMetadataOnTransaction

Set the metadata of a transaction by its ID

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { AddMetadataOnTransactionResponse } from "@formance/formance-sdk/dist/sdk/models/operations";

const sdk = new SDK({
  security: {
    authorization: "",
  },
});

sdk.transactions.addMetadataOnTransaction({
  requestBody: {
    "officia": "dolor",
    "debitis": "a",
    "dolorum": "in",
  },
  ledger: "ledger001",
  txid: 1234,
}).then((res: AddMetadataOnTransactionResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

### Parameters

| Parameter                                                                                                | Type                                                                                                     | Required                                                                                                 | Description                                                                                              |
| -------------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------- |
| `request`                                                                                                | [operations.AddMetadataOnTransactionRequest](../../models/operations/addmetadataontransactionrequest.md) | :heavy_check_mark:                                                                                       | The request object to use for the request.                                                               |
| `config`                                                                                                 | [AxiosRequestConfig](https://axios-http.com/docs/req_config)                                             | :heavy_minus_sign:                                                                                       | Available config options for making requests.                                                            |


### Response

**Promise<[operations.AddMetadataOnTransactionResponse](../../models/operations/addmetadataontransactionresponse.md)>**


## countTransactions

Count the transactions from a ledger

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { CountTransactionsResponse } from "@formance/formance-sdk/dist/sdk/models/operations";

const sdk = new SDK({
  security: {
    authorization: "",
  },
});

sdk.transactions.countTransactions({
  account: "users:001",
  destination: "users:001",
  endTime: new Date("2022-02-26T01:27:36.152Z"),
  ledger: "ledger001",
  metadata: {},
  reference: "ref:001",
  source: "users:001",
  startTime: new Date("2020-11-26T01:41:04.216Z"),
}).then((res: CountTransactionsResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

### Parameters

| Parameter                                                                                  | Type                                                                                       | Required                                                                                   | Description                                                                                |
| ------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------ |
| `request`                                                                                  | [operations.CountTransactionsRequest](../../models/operations/counttransactionsrequest.md) | :heavy_check_mark:                                                                         | The request object to use for the request.                                                 |
| `config`                                                                                   | [AxiosRequestConfig](https://axios-http.com/docs/req_config)                               | :heavy_minus_sign:                                                                         | Available config options for making requests.                                              |


### Response

**Promise<[operations.CountTransactionsResponse](../../models/operations/counttransactionsresponse.md)>**


## createTransaction

Create a new transaction to a ledger

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { CreateTransactionResponse } from "@formance/formance-sdk/dist/sdk/models/operations";

const sdk = new SDK({
  security: {
    authorization: "",
  },
});

sdk.transactions.createTransaction({
  idempotencyKey: "dicta",
  postTransaction: {
    metadata: {
      "cumque": "facere",
      "ea": "aliquid",
    },
    postings: [
      {
        amount: 100,
        asset: "COIN",
        destination: "users:002",
        source: "users:001",
      },
      {
        amount: 100,
        asset: "COIN",
        destination: "users:002",
        source: "users:001",
      },
      {
        amount: 100,
        asset: "COIN",
        destination: "users:002",
        source: "users:001",
      },
    ],
    reference: "ref:001",
    script: {
      plain: "vars {
    account $user
    }
    send [COIN 10] (
    	source = @world
    	destination = $user
    )
    ",
      vars: {},
    },
    timestamp: new Date("2022-04-02T11:21:13.260Z"),
  },
  ledger: "ledger001",
  preview: true,
}).then((res: CreateTransactionResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

### Parameters

| Parameter                                                                                  | Type                                                                                       | Required                                                                                   | Description                                                                                |
| ------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------ |
| `request`                                                                                  | [operations.CreateTransactionRequest](../../models/operations/createtransactionrequest.md) | :heavy_check_mark:                                                                         | The request object to use for the request.                                                 |
| `config`                                                                                   | [AxiosRequestConfig](https://axios-http.com/docs/req_config)                               | :heavy_minus_sign:                                                                         | Available config options for making requests.                                              |


### Response

**Promise<[operations.CreateTransactionResponse](../../models/operations/createtransactionresponse.md)>**


## getTransaction

Get transaction from a ledger by its ID

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { GetTransactionResponse } from "@formance/formance-sdk/dist/sdk/models/operations";

const sdk = new SDK({
  security: {
    authorization: "",
  },
});

sdk.transactions.getTransaction({
  ledger: "ledger001",
  txid: 1234,
}).then((res: GetTransactionResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

### Parameters

| Parameter                                                                            | Type                                                                                 | Required                                                                             | Description                                                                          |
| ------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------ |
| `request`                                                                            | [operations.GetTransactionRequest](../../models/operations/gettransactionrequest.md) | :heavy_check_mark:                                                                   | The request object to use for the request.                                           |
| `config`                                                                             | [AxiosRequestConfig](https://axios-http.com/docs/req_config)                         | :heavy_minus_sign:                                                                   | Available config options for making requests.                                        |


### Response

**Promise<[operations.GetTransactionResponse](../../models/operations/gettransactionresponse.md)>**


## listTransactions

List transactions from a ledger, sorted by txid in descending order.

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { ListTransactionsResponse } from "@formance/formance-sdk/dist/sdk/models/operations";

const sdk = new SDK({
  security: {
    authorization: "",
  },
});

sdk.transactions.listTransactions({
  account: "users:001",
  after: "occaecati",
  cursor: "aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ==",
  destination: "users:001",
  endTime: new Date("2022-02-13T03:59:53.583Z"),
  ledger: "ledger001",
  metadata: {},
  pageSize: 965417,
  reference: "ref:001",
  source: "users:001",
  startTime: new Date("2021-10-28T10:05:29.849Z"),
}).then((res: ListTransactionsResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

### Parameters

| Parameter                                                                                | Type                                                                                     | Required                                                                                 | Description                                                                              |
| ---------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------- |
| `request`                                                                                | [operations.ListTransactionsRequest](../../models/operations/listtransactionsrequest.md) | :heavy_check_mark:                                                                       | The request object to use for the request.                                               |
| `config`                                                                                 | [AxiosRequestConfig](https://axios-http.com/docs/req_config)                             | :heavy_minus_sign:                                                                       | Available config options for making requests.                                            |


### Response

**Promise<[operations.ListTransactionsResponse](../../models/operations/listtransactionsresponse.md)>**


## revertTransaction

Revert a ledger transaction by its ID

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { RevertTransactionResponse } from "@formance/formance-sdk/dist/sdk/models/operations";

const sdk = new SDK({
  security: {
    authorization: "",
  },
});

sdk.transactions.revertTransaction({
  ledger: "ledger001",
  txid: 1234,
}).then((res: RevertTransactionResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

### Parameters

| Parameter                                                                                  | Type                                                                                       | Required                                                                                   | Description                                                                                |
| ------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------ |
| `request`                                                                                  | [operations.RevertTransactionRequest](../../models/operations/reverttransactionrequest.md) | :heavy_check_mark:                                                                         | The request object to use for the request.                                                 |
| `config`                                                                                   | [AxiosRequestConfig](https://axios-http.com/docs/req_config)                               | :heavy_minus_sign:                                                                         | Available config options for making requests.                                              |


### Response

**Promise<[operations.RevertTransactionResponse](../../models/operations/reverttransactionresponse.md)>**

