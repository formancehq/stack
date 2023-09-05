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

const sdk = new SDK({
  security: {
    authorization: "",
  },
});

sdk.balances.getBalances({
  address: "users:001",
  after: "users:003",
  cursor: "aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ==",
  ledger: "ledger001",
  pageSize: 210382,
}).then((res: GetBalancesResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

### Parameters

| Parameter                                                                      | Type                                                                           | Required                                                                       | Description                                                                    |
| ------------------------------------------------------------------------------ | ------------------------------------------------------------------------------ | ------------------------------------------------------------------------------ | ------------------------------------------------------------------------------ |
| `request`                                                                      | [operations.GetBalancesRequest](../../models/operations/getbalancesrequest.md) | :heavy_check_mark:                                                             | The request object to use for the request.                                     |
| `config`                                                                       | [AxiosRequestConfig](https://axios-http.com/docs/req_config)                   | :heavy_minus_sign:                                                             | Available config options for making requests.                                  |


### Response

**Promise<[operations.GetBalancesResponse](../../models/operations/getbalancesresponse.md)>**


## getBalancesAggregated

Get the aggregated balances from selected accounts

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { GetBalancesAggregatedResponse } from "@formance/formance-sdk/dist/sdk/models/operations";

const sdk = new SDK({
  security: {
    authorization: "",
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

### Parameters

| Parameter                                                                                          | Type                                                                                               | Required                                                                                           | Description                                                                                        |
| -------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------- |
| `request`                                                                                          | [operations.GetBalancesAggregatedRequest](../../models/operations/getbalancesaggregatedrequest.md) | :heavy_check_mark:                                                                                 | The request object to use for the request.                                                         |
| `config`                                                                                           | [AxiosRequestConfig](https://axios-http.com/docs/req_config)                                       | :heavy_minus_sign:                                                                                 | Available config options for making requests.                                                      |


### Response

**Promise<[operations.GetBalancesAggregatedResponse](../../models/operations/getbalancesaggregatedresponse.md)>**

