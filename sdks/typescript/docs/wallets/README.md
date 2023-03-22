# wallets

### Available Operations

* [confirmHold](#confirmhold) - Confirm a hold
* [createBalance](#createbalance) - Create a balance
* [createWallet](#createwallet) - Create a new wallet
* [creditWallet](#creditwallet) - Credit a wallet
* [debitWallet](#debitwallet) - Debit a wallet
* [getBalance](#getbalance) - Get detailed balance
* [getHold](#gethold) - Get a hold
* [getHolds](#getholds) - Get all holds for a wallet
* [getTransactions](#gettransactions)
* [getWallet](#getwallet) - Get a wallet
* [getWalletSummary](#getwalletsummary) - Get wallet summary
* [listBalances](#listbalances) - List balances of a wallet
* [listWallets](#listwallets) - List all wallets
* [updateWallet](#updatewallet) - Update a wallet
* [voidHold](#voidhold) - Cancel a hold
* [walletsgetServerInfo](#walletsgetserverinfo) - Get server info

## confirmHold

Confirm a hold

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { ConfirmHoldResponse } from "@formance/formance-sdk/dist/sdk/models/operations";
import { WalletsErrorResponseErrorCode } from "@formance/formance-sdk/dist/sdk/models/shared";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.wallets.confirmHold({
  confirmHoldRequest: {
    amount: 100,
    final: true,
  },
  holdId: "rerum",
}).then((res: ConfirmHoldResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## createBalance

Create a balance

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { CreateBalanceResponse } from "@formance/formance-sdk/dist/sdk/models/operations";
import { WalletsErrorResponseErrorCode } from "@formance/formance-sdk/dist/sdk/models/shared";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.wallets.createBalance({
  createBalanceRequest: {
    expiresAt: new Date("2022-09-14T10:27:07.590Z"),
    name: "Irving Jenkins",
    priority: 881104,
  },
  id: "395efb9b-a88f-43a6-a997-074ba4469b6e",
}).then((res: CreateBalanceResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## createWallet

Create a new wallet

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { CreateWalletResponse } from "@formance/formance-sdk/dist/sdk/models/operations";
import { WalletsErrorResponseErrorCode } from "@formance/formance-sdk/dist/sdk/models/shared";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.wallets.createWallet({
  metadata: {
    "architecto": "magnam",
  },
  name: "Miriam Hermann",
}).then((res: CreateWalletResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## creditWallet

Credit a wallet

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { CreditWalletResponse } from "@formance/formance-sdk/dist/sdk/models/operations";
import { WalletsErrorResponseErrorCode } from "@formance/formance-sdk/dist/sdk/models/shared";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.wallets.creditWallet({
  creditWalletRequest: {
    amount: {
      amount: 574325,
      asset: "accusantium",
    },
    balance: "mollitia",
    metadata: {
      "mollitia": "ad",
      "eum": "dolor",
      "necessitatibus": "odit",
      "nemo": "quasi",
    },
    reference: "iure",
    sources: [
      {
        balance: "eius",
        identifier: "maxime",
        type: "deleniti",
      },
      {
        balance: "in",
        identifier: "architecto",
        type: "architecto",
      },
      {
        balance: "ullam",
        identifier: "expedita",
        type: "nihil",
      },
      {
        balance: "quibusdam",
        identifier: "sed",
        type: "saepe",
      },
    ],
  },
  id: "d028921c-ddc6-4926-81fb-576b0d5f0d30",
}).then((res: CreditWalletResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## debitWallet

Debit a wallet

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { DebitWalletResponse } from "@formance/formance-sdk/dist/sdk/models/operations";
import { WalletsErrorResponseErrorCode } from "@formance/formance-sdk/dist/sdk/models/shared";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.wallets.debitWallet({
  debitWalletRequest: {
    amount: {
      amount: 764912,
      asset: "corporis",
    },
    balances: [
      "libero",
      "nobis",
      "dolores",
      "quis",
    ],
    description: "totam",
    destination: {
      identifier: "eaque",
      type: "quis",
    },
    metadata: {
      "eos": "perferendis",
    },
    pending: false,
  },
  id: "2c73d5fe-9b90-4c28-909b-3fe49a8d9cbf",
}).then((res: DebitWalletResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## getBalance

Get detailed balance

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { GetBalanceResponse } from "@formance/formance-sdk/dist/sdk/models/operations";
import { WalletsErrorResponseErrorCode } from "@formance/formance-sdk/dist/sdk/models/shared";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.wallets.getBalance({
  balanceName: "quaerat",
  id: "8633323f-9b77-4f3a-8100-674ebf69280d",
}).then((res: GetBalanceResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## getHold

Get a hold

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { GetHoldResponse } from "@formance/formance-sdk/dist/sdk/models/operations";
import { WalletsErrorResponseErrorCode } from "@formance/formance-sdk/dist/sdk/models/shared";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.wallets.getHold({
  holdID: "ab",
}).then((res: GetHoldResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## getHolds

Get all holds for a wallet

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { GetHoldsResponse } from "@formance/formance-sdk/dist/sdk/models/operations";
import { WalletsErrorResponseErrorCode } from "@formance/formance-sdk/dist/sdk/models/shared";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.wallets.getHolds({
  cursor: "aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ==",
  metadata: {
    "dolorum": "iusto",
    "voluptate": "dolorum",
    "deleniti": "omnis",
  },
  pageSize: 896672,
  walletID: "distinctio",
}).then((res: GetHoldsResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## getTransactions

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { GetTransactionsResponse } from "@formance/formance-sdk/dist/sdk/models/operations";
import { WalletsErrorResponseErrorCode } from "@formance/formance-sdk/dist/sdk/models/shared";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.wallets.getTransactions({
  cursor: "aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ==",
  pageSize: 990339,
  walletID: "nihil",
}).then((res: GetTransactionsResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## getWallet

Get a wallet

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { GetWalletResponse } from "@formance/formance-sdk/dist/sdk/models/operations";
import { WalletsErrorResponseErrorCode } from "@formance/formance-sdk/dist/sdk/models/shared";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.wallets.getWallet({
  id: "37ae4203-ce5e-46a9-9d8a-0d446ce2af7a",
}).then((res: GetWalletResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## getWalletSummary

Get wallet summary

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { GetWalletSummaryResponse } from "@formance/formance-sdk/dist/sdk/models/operations";
import { WalletsErrorResponseErrorCode } from "@formance/formance-sdk/dist/sdk/models/shared";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.wallets.getWalletSummary({
  id: "73cf3be4-53f8-470b-b26b-5a73429cdb1a",
}).then((res: GetWalletSummaryResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## listBalances

List balances of a wallet

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { ListBalancesResponse } from "@formance/formance-sdk/dist/sdk/models/operations";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.wallets.listBalances({
  id: "8422bb67-9d23-4227-95bf-0cbb1e31b8b9",
}).then((res: ListBalancesResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## listWallets

List all wallets

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { ListWalletsResponse } from "@formance/formance-sdk/dist/sdk/models/operations";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.wallets.listWallets({
  cursor: "aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ==",
  metadata: {
    "delectus": "dolorem",
  },
  name: "Clara Fisher Jr.",
  pageSize: 16429,
}).then((res: ListWalletsResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## updateWallet

Update a wallet

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { UpdateWalletResponse } from "@formance/formance-sdk/dist/sdk/models/operations";
import { WalletsErrorResponseErrorCode } from "@formance/formance-sdk/dist/sdk/models/shared";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.wallets.updateWallet({
  requestBody: {
    metadata: {
      "itaque": "consequatur",
      "est": "repellendus",
      "porro": "doloribus",
    },
  },
  id: "4b921879-fce9-453f-b3ef-7fbc7abd74dd",
}).then((res: UpdateWalletResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## voidHold

Cancel a hold

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { VoidHoldResponse } from "@formance/formance-sdk/dist/sdk/models/operations";
import { WalletsErrorResponseErrorCode } from "@formance/formance-sdk/dist/sdk/models/shared";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.wallets.voidHold({
  holdId: "sequi",
}).then((res: VoidHoldResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## walletsgetServerInfo

Get server info

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { WalletsgetServerInfoResponse } from "@formance/formance-sdk/dist/sdk/models/operations";
import { WalletsErrorResponseErrorCode } from "@formance/formance-sdk/dist/sdk/models/shared";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.wallets.walletsgetServerInfo().then((res: WalletsgetServerInfoResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```
