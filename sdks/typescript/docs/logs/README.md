# logs

### Available Operations

* [listLogs](#listlogs) - List the logs from a ledger

## listLogs

List the logs from a ledger, sorted by ID in descending order.

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { ListLogsResponse } from "@formance/formance-sdk/dist/sdk/models/operations";
import { ErrorsEnum, LogType } from "@formance/formance-sdk/dist/sdk/models/shared";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.logs.listLogs({
  cursor: "aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ==",
  endTime: new Date("2022-09-04T08:35:09.957Z"),
  ledger: "ledger001",
  pageSize: 570197,
  startTime: new Date("2022-07-24T21:51:02.112Z"),
}).then((res: ListLogsResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```
