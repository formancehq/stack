# script

### Available Operations

* [~~runScript~~](#runscript) - Execute a Numscript :warning: **Deprecated**

## ~~runScript~~

This route is deprecated, and has been merged into `POST /{ledger}/transactions`.


> :warning: **DEPRECATED**: this method will be removed in a future release, please migrate away from it as soon as possible.

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { RunScriptResponse } from "@formance/formance-sdk/dist/sdk/models/operations";
import { ErrorsEnum } from "@formance/formance-sdk/dist/sdk/models/shared";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.script.runScript({
  script: {
    metadata: {
      "nisi": "vel",
      "natus": "omnis",
      "molestiae": "perferendis",
    },
    plain: "vars {
  account $user
  }
  send [COIN 10] (
  	source = @world
  	destination = $user
  )
  ",
    reference: "order_1234",
    vars: {
      "magnam": "distinctio",
      "id": "labore",
    },
  },
  ledger: "ledger001",
  preview: true,
}).then((res: RunScriptResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```
