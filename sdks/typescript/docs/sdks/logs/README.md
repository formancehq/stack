# logs

### Available Operations

* [listLogs](#listlogs) - List the logs from a ledger

## listLogs

List the logs from a ledger, sorted by ID in descending order.

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { ListLogsResponse } from "@formance/formance-sdk/dist/sdk/models/operations";

const sdk = new SDK({
  security: {
    authorization: "",
  },
});

sdk.logs.listLogs({
  after: "corporis",
  cursor: "aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ==",
  endTime: new Date("2022-04-01T23:59:21.675Z"),
  ledger: "ledger001",
  pageSize: 315428,
  startTime: new Date("2022-04-10T11:47:13.463Z"),
}).then((res: ListLogsResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

### Parameters

| Parameter                                                                | Type                                                                     | Required                                                                 | Description                                                              |
| ------------------------------------------------------------------------ | ------------------------------------------------------------------------ | ------------------------------------------------------------------------ | ------------------------------------------------------------------------ |
| `request`                                                                | [operations.ListLogsRequest](../../models/operations/listlogsrequest.md) | :heavy_check_mark:                                                       | The request object to use for the request.                               |
| `config`                                                                 | [AxiosRequestConfig](https://axios-http.com/docs/req_config)             | :heavy_minus_sign:                                                       | Available config options for making requests.                            |


### Response

**Promise<[operations.ListLogsResponse](../../models/operations/listlogsresponse.md)>**

