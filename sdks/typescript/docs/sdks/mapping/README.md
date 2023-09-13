# mapping

### Available Operations

* [getMapping](#getmapping) - Get the mapping of a ledger
* [updateMapping](#updatemapping) - Update the mapping of a ledger

## getMapping

Get the mapping of a ledger

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { GetMappingResponse } from "@formance/formance-sdk/dist/sdk/models/operations";

const sdk = new SDK({
  security: {
    authorization: "",
  },
});

sdk.mapping.getMapping({
  ledger: "ledger001",
}).then((res: GetMappingResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

### Parameters

| Parameter                                                                    | Type                                                                         | Required                                                                     | Description                                                                  |
| ---------------------------------------------------------------------------- | ---------------------------------------------------------------------------- | ---------------------------------------------------------------------------- | ---------------------------------------------------------------------------- |
| `request`                                                                    | [operations.GetMappingRequest](../../models/operations/getmappingrequest.md) | :heavy_check_mark:                                                           | The request object to use for the request.                                   |
| `config`                                                                     | [AxiosRequestConfig](https://axios-http.com/docs/req_config)                 | :heavy_minus_sign:                                                           | Available config options for making requests.                                |


### Response

**Promise<[operations.GetMappingResponse](../../models/operations/getmappingresponse.md)>**


## updateMapping

Update the mapping of a ledger

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { UpdateMappingResponse } from "@formance/formance-sdk/dist/sdk/models/operations";

const sdk = new SDK({
  security: {
    authorization: "",
  },
});

sdk.mapping.updateMapping({
  mapping: {
    contracts: [
      {
        account: "users:001",
        expr: {},
      },
      {
        account: "users:001",
        expr: {},
      },
    ],
  },
  ledger: "ledger001",
}).then((res: UpdateMappingResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

### Parameters

| Parameter                                                                          | Type                                                                               | Required                                                                           | Description                                                                        |
| ---------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------- |
| `request`                                                                          | [operations.UpdateMappingRequest](../../models/operations/updatemappingrequest.md) | :heavy_check_mark:                                                                 | The request object to use for the request.                                         |
| `config`                                                                           | [AxiosRequestConfig](https://axios-http.com/docs/req_config)                       | :heavy_minus_sign:                                                                 | Available config options for making requests.                                      |


### Response

**Promise<[operations.UpdateMappingResponse](../../models/operations/updatemappingresponse.md)>**

