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
import { ErrorsEnum } from "@formance/formance-sdk/dist/sdk/models/shared";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
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

## updateMapping

Update the mapping of a ledger

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { UpdateMappingResponse } from "@formance/formance-sdk/dist/sdk/models/operations";
import { ErrorsEnum } from "@formance/formance-sdk/dist/sdk/models/shared";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.mapping.updateMapping({
  mapping: {
    contracts: [
      {
        account: "users:001",
        expr: {
          "numquam": "commodi",
          "quam": "molestiae",
          "velit": "error",
        },
      },
      {
        account: "users:001",
        expr: {
          "quis": "vitae",
        },
      },
      {
        account: "users:001",
        expr: {
          "animi": "enim",
          "odit": "quo",
          "sequi": "tenetur",
        },
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
