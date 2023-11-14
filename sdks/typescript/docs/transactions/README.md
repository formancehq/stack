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
import { ErrorsEnum } from "@formance/formance-sdk/dist/sdk/models/shared";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.transactions.createTransactions({
  transactions: {
    transactions: [
      {
        metadata: {
          "sint": "accusantium",
          "mollitia": "reiciendis",
          "mollitia": "ad",
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
        timestamp: new Date("2022-02-07T18:15:06.372Z"),
      },
      {
        metadata: {
          "nemo": "quasi",
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
        timestamp: new Date("2020-04-29T08:15:14.819Z"),
      },
      {
        metadata: {
          "maxime": "deleniti",
          "facilis": "in",
        },
        postings: [
          {
            amount: 100,
            asset: "COIN",
            destination: "users:002",
            source: "users:001",
          },
        ],
        reference: "ref:001",
        timestamp: new Date("2022-01-30T09:19:56.236Z"),
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

## addMetadataOnTransaction

Set the metadata of a transaction by its ID

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { AddMetadataOnTransactionResponse } from "@formance/formance-sdk/dist/sdk/models/operations";
import { ErrorsEnum } from "@formance/formance-sdk/dist/sdk/models/shared";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.transactions.addMetadataOnTransaction({
  requestBody: {
    "expedita": "nihil",
    "repellat": "quibusdam",
  },
  ledger: "ledger001",
  txid: 1234,
}).then((res: AddMetadataOnTransactionResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## countTransactions

Count the transactions from a ledger

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { CountTransactionsResponse } from "@formance/formance-sdk/dist/sdk/models/operations";
import { ErrorsEnum } from "@formance/formance-sdk/dist/sdk/models/shared";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.transactions.countTransactions({
  account: "users:001",
  destination: "users:001",
  endTime: new Date("2022-02-04T19:17:08.641Z"),
  ledger: "ledger001",
  metadata: {
    "accusantium": "consequuntur",
    "praesentium": "natus",
    "magni": "sunt",
    "quo": "illum",
  },
  reference: "ref:001",
  source: "users:001",
  startTime: new Date("2020-07-30T23:39:27.609Z"),
}).then((res: CountTransactionsResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## createTransaction

Create a new transaction to a ledger

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { CreateTransactionResponse } from "@formance/formance-sdk/dist/sdk/models/operations";
import { ErrorsEnum } from "@formance/formance-sdk/dist/sdk/models/shared";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.transactions.createTransaction({
  idempotencyKey: "ea",
  postTransaction: {
    metadata: {
      "odit": "ea",
      "accusantium": "ab",
      "maiores": "quidem",
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
    script: {
      plain: "vars {
    account $user
    }
    send [COIN 10] (
    	source = @world
    	destination = $user
    )
    ",
      vars: {
        "autem": "nam",
        "eaque": "pariatur",
      },
    },
    timestamp: new Date("2022-01-09T22:25:53.570Z"),
  },
  ledger: "ledger001",
  preview: true,
}).then((res: CreateTransactionResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## getTransaction

Get transaction from a ledger by its ID

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { GetTransactionResponse } from "@formance/formance-sdk/dist/sdk/models/operations";
import { ErrorsEnum } from "@formance/formance-sdk/dist/sdk/models/shared";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
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

## listTransactions

List transactions from a ledger, sorted by txid in descending order.

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { ListTransactionsResponse } from "@formance/formance-sdk/dist/sdk/models/operations";
import { ErrorsEnum } from "@formance/formance-sdk/dist/sdk/models/shared";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.transactions.listTransactions({
  account: "users:001",
  after: "perferendis",
  cursor: "aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ==",
  destination: "users:001",
  endTime: new Date("2022-04-23T08:05:27.776Z"),
  ledger: "ledger001",
  metadata: {
    "cumque": "corporis",
  },
  pageSize: 944124,
  reference: "ref:001",
  source: "users:001",
  startTime: new Date("2021-07-02T12:00:47.197Z"),
}).then((res: ListTransactionsResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## revertTransaction

Revert a ledger transaction by its ID

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { RevertTransactionResponse } from "@formance/formance-sdk/dist/sdk/models/operations";
import { ErrorsEnum } from "@formance/formance-sdk/dist/sdk/models/shared";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
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
