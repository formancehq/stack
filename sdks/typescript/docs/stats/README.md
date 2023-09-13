# stats

### Available Operations

* [readStats](#readstats) - Get statistics from a ledger

## readStats

Get statistics from a ledger. (aggregate metrics on accounts and transactions)


### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { ReadStatsResponse } from "@formance/formance-sdk/dist/sdk/models/operations";
import { ErrorsEnum } from "@formance/formance-sdk/dist/sdk/models/shared";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.stats.readStats({
  ledger: "ledger001",
}).then((res: ReadStatsResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```
