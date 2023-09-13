# balances

### Available Operations

* [getBalances](#getbalances) - Get the balances from a ledger's account
* [getBalancesAggregated](#getbalancesaggregated) - Get the aggregated balances from selected accounts

## getBalances

Get the balances from a ledger's account

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { GetBalancesResponse } from "@formance/formance-sdk/dist/sdk/models/operations";
import { ErrorsEnum } from "@formance/formance-sdk/dist/sdk/models/shared";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.balances.getBalances({
  address: "users:001",
  after: "users:003",
  cursor: "aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ==",
  ledger: "ledger001",
  pageSize: 958950,
}).then((res: GetBalancesResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## getBalancesAggregated

Get the aggregated balances from selected accounts

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { GetBalancesAggregatedResponse } from "@formance/formance-sdk/dist/sdk/models/operations";
import { ErrorsEnum } from "@formance/formance-sdk/dist/sdk/models/shared";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.balances.getBalancesAggregated({
  address: "users:001",
  ledger: "ledger001",
}).then((res: GetBalancesAggregatedResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```
